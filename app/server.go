package app

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/libsql/libsql-client-go/libsql"
	"github.com/richgrov/constructify/app/jobs"
	"go.uber.org/zap"
)

type constructorService interface {
	build(id string, prompt string) (string, error)
}

type app struct {
	logger      *zap.Logger
	constructor constructorService
	jobsManager jobs.JobService
}

func (app *app) index(c *gin.Context) {
	c.HTML(http.StatusOK, "index/index.tmpl", gin.H{})
}

func (app *app) view(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	job, err := app.jobsManager.GetJob(id)

	if job == nil {
		c.HTML(http.StatusNotFound, "view/index.tmpl", gin.H{
			"object":        "assets/schem/404.glb",
			"prompt":        "Sorry, couldn't find that build.",
			"skipAnimation": true,
		})
		return
	}

	if err != nil {
		app.logger.Error("query failed: ", zap.String("id", id), zap.Error(err))
		c.HTML(http.StatusInternalServerError, "view/index.tmpl", gin.H{
			"object":        "assets/schem/error.glb",
			"prompt":        "Sorry, that's our fault. Please try again later.",
			"skipAnimation": true,
		})
		return
	}

	switch job.Status {
	case jobs.StatusSucceess, jobs.StatusInProgress:
		c.HTML(http.StatusOK, "view/index.tmpl", gin.H{
			"prompt":        job.Prompt,
			"object":        job.Result,
			"id":            id,
			"skipAnimation": true,
		})
	case jobs.StatusFailed:
		c.HTML(http.StatusNotImplemented, "view/index.tmpl", gin.H{
			"prompt":        "Sorry, couldn't make that. Try again later?",
			"object":        "assets/schem/x.glb",
			"skipAnimation": true,
		})
	}
}

func (app *app) generate(c *gin.Context) {
	id, err := uuid.NewRandom()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index/prompt.tmpl", gin.H{
			"error": "Internal server error",
		})
		return
	}

	prompt := c.PostForm("prompt")
	if len(strings.TrimSpace(prompt)) == 0 {
		c.HTML(http.StatusBadRequest, "index/prompt.tmpl", gin.H{
			"error": "Invalid prompt",
		})
		return
	}

	err = app.jobsManager.StartJob(id.String(), prompt)
	if err != nil {
		app.logger.Log(zap.ErrorLevel, "exec failed", zap.Error(err))
		c.HTML(http.StatusInternalServerError, "index/prompt.tmpl", gin.H{
			"error": "Internal server error",
		})
		return
	}

	go app.constructAndStoreResult(id, prompt)

	c.Header("HX-Redirect", "/view?id="+id.String())
	c.Data(http.StatusOK, gin.MIMEHTML, nil)
}

func (app *app) object(c *gin.Context) {
	c.Header("HX-Trigger-After-Settle", "displayMesh")

	id := c.Query("id")
	if len(id) == 0 {
		c.HTML(http.StatusBadRequest, "view/object.tmpl", gin.H{
			"object": "assets/schem/error.glb",
			"prompt": "Build ID is missing",
		})
		return
	}

	job, err := app.jobsManager.WaitForCompletion(id, 2*time.Minute)
	if err == jobs.ErrTimeout {
		c.HTML(http.StatusGatewayTimeout, "view/object.tmpl", gin.H{
			"prompt": "Took too long to generate. Try something simpler?",
			"object": "assets/schem/x.glb",
		})
		return
	} else if err != nil {
		app.logger.Error("query failed: ", zap.String("id", id), zap.Error(err))
		c.HTML(http.StatusInternalServerError, "view/object.tmpl", gin.H{
			"object": "assets/schem/error.glb",
			"prompt": "Sorry, that's our fault. Please try again later.",
		})
		return
	}

	switch job.Status {
	case jobs.StatusSucceess:
		c.HTML(http.StatusOK, "view/object.tmpl", gin.H{
			"object": job.Result,
			"prompt": job.Prompt,
		})
	case jobs.StatusFailed:
		c.HTML(http.StatusNotImplemented, "view/object.tmpl", gin.H{
			"prompt": "Sorry, couldn't make that. Try again later?",
			"object": "assets/schem/x.glb",
		})
	}
}

func (app *app) constructAndStoreResult(id uuid.UUID, prompt string) {
	objectPath, err := app.constructor.build(id.String(), prompt)

	var status jobs.Status
	var result string
	if err != nil {
		app.logger.Log(zap.ErrorLevel, "failed to generate build",
			zap.String("id", id.String()),
			zap.String("prompt", prompt),
			zap.Error(err),
		)
		status = jobs.StatusFailed
		result = err.Error()
	} else {
		status = jobs.StatusSucceess
		result = objectPath
	}

	err = app.jobsManager.UpdateJobStatus(id.String(), status, result)
	if err != nil {
		app.logger.Log(zap.ErrorLevel, "exec failed", zap.Error(err))
	}
}

func Run(dbUrl string, logger *zap.Logger) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*/*.tmpl")
	router.Static("/assets", "assets/")
	router.Static("/dist", "dist/")

	var constructor constructorService
	serviceEnv := os.Getenv("CONSTRUCTOR_SERVICE")
	if serviceEnv == "mock" {
		constructor = &mockConstructorService{}
	} else if serviceEnv != "" {
		constructor = &httpConstructorService{
			url: serviceEnv,
		}
	} else {
		panic("variable CONSTRUCTOR_SERVICE not set")
	}

	jobService, err := jobs.NewSqlServce(dbUrl)
	if err != nil {
		panic(err)
	}

	app := &app{
		logger:      logger,
		constructor: constructor,
		jobsManager: jobService,
	}

	router.GET("/", app.index)
	router.GET("/view", app.view)
	router.POST("/generate", app.generate)
	router.GET("/object", app.object)

	address := os.Getenv("BIND")
	if address == "" {
		address = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(address + ":" + port)
}

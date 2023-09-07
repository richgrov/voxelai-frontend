package main

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/libsql/libsql-client-go/libsql"
	"go.uber.org/zap"
)

type app struct {
	db     *sql.DB
	logger *zap.Logger
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

	var prompt string
	var status sql.NullString
	var result sql.NullString
	err := app.db.QueryRow("SELECT prompt, status, result FROM jobs WHERE id=? LIMIT 1", id).Scan(&prompt, &status, &result)

	if err == sql.ErrNoRows {
		c.HTML(http.StatusNotFound, "view/index.tmpl", gin.H{
			"error":  "Not found",
			"prompt": "Sorry, couldn't find that build.",
		})
		return
	}

	if err != nil {
		app.logger.Error("query failed: ", zap.String("id", id), zap.Error(err))
		c.HTML(http.StatusInternalServerError, "view/index.tmpl", gin.H{
			"error":  "Internal Error",
			"prompt": "Sorry, that's our fault. Please try again later.",
		})
		return
	}

	if status.String == "FAILED" {
		c.HTML(http.StatusNotImplemented, "view/index.tmpl", gin.H{
			"prompt": "Sorry, couldn't make that",
			"error":  "Not completable",
		})
		return
	}

	c.HTML(http.StatusOK, "view/index.tmpl", gin.H{
		"prompt": prompt,
		"object": result.String,
		"id":     id,
	})
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

	_, err = app.db.Exec("INSERT INTO jobs (id, prompt) VALUES (?, ?)", id, prompt)
	if err != nil {
		app.logger.Log(zap.ErrorLevel, "exec failed", zap.Error(err))
		c.HTML(http.StatusInternalServerError, "index/prompt.tmpl", gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.Header("HX-Redirect", "/view?id="+id.String())
	c.Data(http.StatusOK, gin.MIMEHTML, nil)
}

func (app *app) object(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.HTML(http.StatusBadRequest, "view/error.tmpl", "Build ID is missing")
		return
	}

	maxPollTime := time.Now().Add(2 * time.Minute)
	for {
		now := time.Now()
		if now.Compare(maxPollTime) > -1 {
			break
		}

		var status sql.NullString
		var result sql.NullString
		err := app.db.QueryRow("SELECT status, result FROM jobs WHERE id=? LIMIT 1", id).Scan(&status, &result)

		if err != nil {
			// No need for custom message if the row wasn't found. That is handled by
			// /view. If the row somehow got deleted between the page load and this
			// request, it's an internal mishap
			app.logger.Error("query failed: ", zap.String("id", id), zap.Error(err))
			c.HTML(http.StatusInternalServerError, "view/object.tmpl", gin.H{
				"error": "Internal server error- Please try again later.",
			})
			return
		}

		switch status.String {
		case "COMPLETED":
			c.HTML(http.StatusOK, "view/object.tmpl", result.String)
			return
		case "FAILED":
			c.HTML(http.StatusNotImplemented, "view/error.tmpl", gin.H{
				"error": "Couldn't figure out how to make that!",
			})
			return
		}

		time.Sleep(2 * time.Second)
	}

	c.HTML(http.StatusGatewayTimeout, "view/error.tmpl", gin.H{
		"error": "Build took too long to generate. Try again?",
	})
}
func runApp(db *sql.DB, logger *zap.Logger) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*/*.tmpl")
	router.Static("/assets", "assets/")
	router.Static("/dist", "dist/")

	app := &app{
		db:     db,
		logger: logger,
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

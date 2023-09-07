package main

import (
	"database/sql"
	"net/http"
	"strings"

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
	var object sql.NullString
	switch err := app.db.QueryRow("SELECT prompt, object FROM jobs WHERE id=? LIMIT 1", id).Scan(&prompt, &object); err {
	case nil:
		c.HTML(http.StatusOK, "view/index.tmpl", gin.H{
			"prompt": prompt,
		})

	case sql.ErrNoRows:
		c.HTML(http.StatusNotFound, "view/index.tmpl", gin.H{})

	default:
		app.logger.Error("query failed: ", zap.String("id", id), zap.Error(err))
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

	_, err = app.db.Exec("INSERT INTO jobs (id, prompt, object) VALUES (?, ?, ?)", id, prompt, nil)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index/prompt.tmpl", gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.Header("HX-Redirect", "/view?id="+id.String())
	c.Data(http.StatusOK, gin.MIMEHTML, nil)
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

package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/richgrov/constructify/app"
	"go.uber.org/zap"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	if *debug {
		if err := godotenv.Load(".env.local"); err != nil {
			panic(err)
		}
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("DB_URL not set")
	}

	db, err := sql.Open("libsql", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	if *debug {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS jobs (id text NOT NULL PRIMARY KEY, prompt text NOT NULL, status text, result text)")
		if err != nil {
			panic(err)
		}
	}

	logger, err := createLogger(*debug)
	if err != nil {
		panic(err)
	}

	app.Run(db, logger)
}

func createLogger(debug bool) (*zap.Logger, error) {
	if debug {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}

package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/joho/godotenv"
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
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS jobs (id text, prompt text, object text)")
		if err != nil {
			panic(err)
		}
	}

	logger, err := createLogger(*debug)
	if err != nil {
		panic(err)
	}

	runApp(db, logger)
}

func createLogger(debug bool) (*zap.Logger, error) {
	if debug {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}

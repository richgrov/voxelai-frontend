package main

import (
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

	logger, err := createLogger(*debug)
	if err != nil {
		panic(err)
	}

	app.Run(dbUrl, logger)
}

func createLogger(debug bool) (*zap.Logger, error) {
	if debug {
		return zap.NewDevelopment()
	} else {
		return zap.NewProduction()
	}
}

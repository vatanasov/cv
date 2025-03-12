package main

import (
	"autobiography/internal/database"
	"autobiography/internal/models"
	"flag"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
)

type config struct {
	httpPort int
	debug    bool
}

type application struct {
	config config
	logger *slog.Logger
	wg     *sync.WaitGroup
	models models.Models
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.New("./db.sqlite")
	if err != nil {
		logger.Error("error connecting to database", "err", err)
		os.Exit(1)
	}

	err = run(logger, db)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

}

func run(logger *slog.Logger, db database.Db) error {
	var cfg config

	flag.IntVar(&cfg.httpPort, "http-port", 4444, "port to listen on for HTTP requests")
	flag.BoolVar(&cfg.debug, "debug", false, "enable debug mode")

	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	return app.serveHTTP()
}

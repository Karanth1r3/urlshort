package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/Karanth1r3/url-short-learn/internal/config"
	"github.com/Karanth1r3/url-short-learn/internal/util/logger/slg"

	_ "github.com/mattn/go-sqlite3"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Initializing config (parsing from file): cleanenv
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	cfg := config.MustLoad()

	fmt.Println(cfg)

	// Initializing slog logger
	log := setupLogger(cfg.Env)
	// Showing some info about current setup through logger
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enables")

	//Initializing storage

	storage, err := sql.Open("sqlite3", "./storage/storage.db")
	//sqlite.New(cfg.StoragePath)
	if err != nil {
		//fmt.Println(err)
		log.Error("failed to init storage", slg.Err(err))
		os.Exit(1) // Terminating program with exitcode 1
	}
	defer storage.Close()
	//defer storage.DB.Close()
	_ = storage
	//fmt.Println(storage)

	// TODO: init router: chi, chi render

	// TODO" run server:
}

// Logger setup func depends on environment because on local machine text logs shall be enough, in other cases different option may be preferrable
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		// Initializing logger with simple text output & debug level for local
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			// Initializing logger with json handler & debug level for dev
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		// Initializing logger with Info level for prod
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

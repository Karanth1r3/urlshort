package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Karanth1r3/url-short-learn/internal/config"
	"github.com/Karanth1r3/url-short-learn/internal/httpapi/mw"
	"github.com/Karanth1r3/url-short-learn/internal/httpapi/server/handlers/url/deleter"
	"github.com/Karanth1r3/url-short-learn/internal/httpapi/server/handlers/url/redirect"
	"github.com/Karanth1r3/url-short-learn/internal/httpapi/server/handlers/url/save"
	"github.com/Karanth1r3/url-short-learn/internal/storage/pg"
	"github.com/Karanth1r3/url-short-learn/internal/util"
	"github.com/Karanth1r3/url-short-learn/internal/util/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	db, err := util.ConnectDB(cfg.DB)
	if err != nil {
		log.Error("failed to init storage", slg.Err(err))
		os.Exit(1) // Terminating program with exitcode 1
	}

	defer db.Close()

	storage := pg.New(db)

	//_ = storage

	router := chi.NewRouter()

	// middleware
	// mw for getting req id for tracing
	//router.Use(middleware.RequestID)
	// mw for getting ips of requesters
	//router.Use(middleware.RealIP)
	// mw for request loggers (internal logger of chi)
	router.Use(middleware.Logger)
	router.Use(mw.New(log))
	//mw for recovering from panic
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// urls with /url will be handleded by this subrouter func
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		// this Handlers will be available only for authorized /url is not required as it is already present in group (above)
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", deleter.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting sever", slog.String("address: ", cfg.HTTPServer.Address))

	// Configuring server before launch
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// ListenAndServe is blocking. If execution passed it's position => server has stopped
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

	/*
		err = storage.SaveURL("asdf", "asdf")
		if err != nil {
			log.Error("failed to save url", slg.Err(err))
			os.Exit(1)
		}
	*/
	/*
		log.Info("saved url")

		err = storage.DeleteURL("asd")
		if err != nil {
			log.Error("failed to get url", slg.Err(err))
			os.Exit(1)
		}
		log.Info("deleted")
		//fmt.Println(storage)
	*/
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

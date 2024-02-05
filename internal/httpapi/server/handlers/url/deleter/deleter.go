package deleter

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/Karanth1r3/url-short-learn/internal/httpapi/response"
	"github.com/Karanth1r3/url-short-learn/internal/storage"
	"github.com/Karanth1r3/url-short-learn/internal/util/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const blockDesc = "handlers.url.delete.New()"
		// Setting up local logger
		log = log.With(
			slog.String("op", blockDesc),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		//Trying to parse alias from request body
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}
		err := urlDeleter.DeleteURL(alias)
		// If alias not found => notify client & exit
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias: ", alias)

			render.JSON(w, r, resp.Error("alias not found"))

			return
		}
		if err != nil {
			// Logging error
			log.Info("could not delete url", slg.Err(err))
			// Sending info to client
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))
	}
}

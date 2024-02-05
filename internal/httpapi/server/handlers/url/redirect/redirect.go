package redirect

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

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const blockDesc = "handlers.url.redirect.New()"

		log = log.With(
			slog.String("op", blockDesc),
			slog.String("requiest_id", middleware.GetReqID(r.Context())),
		)
		// Hard link to chi package, not recommended
		alias := chi.URLParam(r, "alias")
		// If alias is empty => Send client according message & return
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias: ", alias)

			render.JSON(w, r, resp.Error("alias not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", slg.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		// If alias found => log it
		log.Info("got url", slog.String("url", resURL))
		//Redirect line
		http.Redirect(w, r, resURL, http.StatusFound)

	}
}

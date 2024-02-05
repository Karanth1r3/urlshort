package save

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"

	"github.com/Karanth1r3/url-short-learn/internal/httpapi/random"
	resp "github.com/Karanth1r3/url-short-learn/internal/httpapi/response"
	"github.com/Karanth1r3/url-short-learn/internal/storage"
	"github.com/Karanth1r3/url-short-learn/internal/util/logger/slg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// Response will contain json, this thing will be used to parse it and get info in required state (Request will be parsed to this)
type Request struct {
	URL   string `json:"url" validate:"required,url"` // validate is tag for validator. required => the field is necessary. url - type for check
	Alias string `json:"alias,omitempty"`
}

// For parsing responses
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"` // If alias is not present in request (omitepmty) => handler will generate it. That's why it's here
}

// Declaring shorter storage interface variant at the place of usage
type URLSaver interface {
	SaveURL(urlToSave, alias string) error
}

// Generatable alias length (for cases with empty alias in request)
// TODO - probably will be moved to config
const (
	aliasLength = 6
)

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const blockDesc = "handlers.url.save.New()"

		log = log.With(
			slog.String("op", blockDesc),
			slog.String("requiest_id", middleware.GetReqID(r.Context())), // For request tracing

		)
		var req Request
		// Unmarshalling with chi internal render addition package
		err := render.DecodeJSON(r.Body, &req)
		buf := make([]byte, 4096)
		buf, _ = httputil.DumpRequest(r, true)

		if err != nil {
			log.Error("failed to decode request body", slg.Err(err)) // slg.Err - wrapping errors for slog (in (util.slg) package )
			fmt.Println(string(buf))

			render.JSON(w, r, resp.Error("failed to decode request"))
			// render.JSON(w, r, Response {
			//	Status: StatusError,
			//		Error: msg,
			// }) // If without created Response type from httpapi.response - it would have been supposed to look like that instead of last uncommented line

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors) // setting type to validation error type
			log.Error("invalid request", slg.Err(err))      // In log there will be unformatted error

			//render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr)) // Forma			// render.JSON(w, r, Response {
			//	Status: StatusError,
			//		Error: msg,
			// }) // If without created Response type from httpapi.response - it would have been supposed to look like that instead of last uncommented lineting errors to readable condition with (httpapi.response internal package)

			return
		}
		// alias is not a neccessary field, so it's check is here
		alias := req.Alias
		if alias == "" {
			// TODO - handle cases when generated alias already exists in storage .........................................................
			// If alias is empty => generate it on server side
			alias = random.NewRandomString(aliasLength)
		}

		err = urlSaver.SaveURL(req.URL, alias)
		// If url exists => send according message
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", slg.Err(err)) // Log can contain used db techs (for example, pg). Could be unwanted behaviour

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}
		// If request was handleded correctly => write it to log
		log.Info("url added")

		responseOK(w, r, alias)
	}
}

// Forming successfull response to save req
func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

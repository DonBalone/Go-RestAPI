package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	resp "url-shortener/internal/http-server/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

// здесь будет храниться функция конструктор для хендлера
type Request struct {
	URL   string `json:"url" validate:"required,url"` // трек validate дает понять валидатору, что это обязательное поле и что это должен быт url
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failsed to decode requesr"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		//валидация это првоерка чего-либо на
		//соответствие заданным условиям
		if err := validator.New().Struct(req); err != nil { // проверка на валидацию
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request body", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" { // если пользователь не указал алиас
			alias = random.NewRandomString(aliasLength)
		}

		// сохранение urlSaver в id
		id, err := urlSaver.SaveURL(req.URL, alias)
		// если url уже существует
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("URL already exists"))

			return
		}
		// ошибка при добавлении url
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}
		// url точно есть, теперь пишем, что он есть
		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}

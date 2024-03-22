package redirect

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/http-server/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

// URLGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

// функция - конструктор хендлера

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	// ф-я которая является самим хендлером
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"
		// стандартный логгер
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		// получаем alias из router
		alias := chi.URLParam(r, "alias")
		// если alias нет
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}
		// обращение к нашему геттеру
		resURL, err := urlGetter.GetURL(alias) // запросили url по alias
		// обрабатываем ошибку из файла storage.go
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		// говорим, что url получили
		log.Info("got url", slog.String("url", resURL))

		// redirect to found url
		// redirect - перенаправление пользователя с
		// одной страницы на другую
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}

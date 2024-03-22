package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)
			// пример запроса
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)
			// создание нового запроса
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			// проверка, что здесь не должно быть ошибки
			require.NoError(t, err) // require нужен, когда после него все будет сломано
			// assert же нужен когда мы хотим проверить несколько кейсов, которые не зависят друг от друга
			// т.е. продолжить затем работу кода

			//тестирование http сервера
			rr := httptest.NewRecorder()
			// запускаем наш запрос
			handler.ServeHTTP(rr, req)
			// смотрим, какой статус вернул recorder
			require.Equal(t, rr.Code, http.StatusOK)
			// что было записано в тело
			body := rr.Body.String()

			var resp save.Response
			// проверка на то, норм ли файлы в json'e,
			// unmarshal разбирает json данные  в указанный объект,
			// где body - входные данные в json формате
			// &resp - куда нужно разобрать данные json
			// возвращает внутрянка(json...) err, которые потом уже require
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			// смотрим, совпадает ли ошибка, которую вернул handler
			// совпадает с ошибкой из testa
			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}
}

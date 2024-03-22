package api

import (
	"errors"
	"fmt"
	"net/http"
)

// если код не 302, то вернет ошибку
var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

// GetRedirect returns the final URL after redirection.
func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"
	// http клиент с параметром check redirect
	client := &http.Client{
		// check нужен для того, чтоыб понимать, нужно ли
		// выполнять перенаправление и как обрабатывать
		// перенаправление
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// возвращаем особую ошибку после которой наша функция не будет
			// следовать нашим редиректам
			return http.ErrUseLastResponse // stop after 1st redirect
		},
	}
	// делаем запрос
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	// закрываем тело ответа, после проработки всех функций
	defer func() { _ = resp.Body.Close() }()
	// проверяем статус ответа
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}
	// возвращаем тот url, на который происходит редирект
	return resp.Header.Get("Location"), nil
}

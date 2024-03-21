package random

import (
	"math/rand" // генерирует псевдо случайные числа
	"time"
)

func NewRandomString(length int) string {
	// установка сида для нашего рандомайзера
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	// слайс из символов
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// до слайс - буфер для символов
	b := make([]rune, length)
	//заполнения слайса b случайным набором символов
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}

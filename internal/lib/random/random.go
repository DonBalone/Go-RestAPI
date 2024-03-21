package random

import (
	"math/rand" // генерирует псевдо случайные числа
	"time"
)

func NewRandomString(length int) string {
	// установка сида для нашего рандомайзера
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}

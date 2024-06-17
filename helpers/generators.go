package helpers

import "math/rand"

func GenerateRandomString(size int, source string) string {
	var letter = []byte(source)
	b := make([]byte, size)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func GenerateWebSaveToken() string {
	return GenerateRandomString(128, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@!-_~")
}

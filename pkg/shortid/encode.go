package shortid

import (
	"math/rand"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func Encode(length int) string {
	result := make([]byte, length)
	for i := range length {
		result[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(result)
}

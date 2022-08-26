package utils

import (
	"FunNow/url-shortener/constants"
	"math/rand"
)

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = constants.LetterBytes[rand.Int63()%int64(len(constants.LetterBytes))]
	}
	return string(b)
}

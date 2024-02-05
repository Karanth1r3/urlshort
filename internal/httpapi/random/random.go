package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	// Random based on unix-time seed
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))] // Getting random char from "chars" var
	}

	return string(b)
}

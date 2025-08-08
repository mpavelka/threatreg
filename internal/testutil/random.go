package testutil

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func AddRandSuffix(value string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	suffix := make([]byte, 8)
	for i := range suffix {
		suffix[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return value + string(suffix)
}

func RandString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

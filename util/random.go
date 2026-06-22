package util

import (
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	k := len(letters)
	for i := range b {
		b[i] = letters[rnd.Intn(k)]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)
	return currencies[rnd.Intn(n)]
}

func RandomEmail() string {
	return RandomString(6) + "@example.com"
}

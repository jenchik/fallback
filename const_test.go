package fallback_test

import (
	"math/rand"
	"time"
)

const (
	chars  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	items  = 1000
	lenKey = 10
	steps  = 500

	benchWorks = 1000
	testWorks  = 100

	lazyLoaderTimeout = time.Millisecond
)

func randString(n int) string {
	buf := make([]byte, n)
	l := len(chars)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < n; i++ {
		buf[i] = chars[rand.Intn(l)]
	}
	return string(buf)
}

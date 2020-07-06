package tests

import (
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

var (
	// Rand is the shared random number generator.
	Rand *rand.Rand
)

func init() {
	Rand = rand.New(rand.NewSource(time.Now().Unix()))
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[Rand.Intn(len(letters))]
	}
	return string(b)
}

func randomBytes(length int64) ([]byte, error) {
	return ioutil.ReadAll(randomReader(length))
}

func randomReader(length int64) io.Reader {
	// Can't use the global Rand here because Rand.Read is not concurrency safe.
	return io.LimitReader(rand.New(rand.NewSource(time.Now().Unix())), length)
}

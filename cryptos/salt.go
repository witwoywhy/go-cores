package cryptos

import (
	"time"

	"math/rand"
)

const (
	CharsetWithoutSpecialChar = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charset                   = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" + "~!@#$%^&*()_"
)

var GetSaltKey func(lenght int64, charsets ...string) string = GenerateSaltKey

func GenerateSaltKey(lenght int64, charsets ...string) string {
	var char string
	if len(charsets) > 0 {
		char = charsets[0]
	} else {
		char = charset
	}

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, lenght)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(char))]
	}

	return string(b)
}

func SetMockSaltKeyFunc() {
	GetSaltKey = func(lenght int64, charsets ...string) string {
		return "mockedSaltKey"
	}
}

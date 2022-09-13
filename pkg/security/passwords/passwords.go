package passwords

import (
	"crypto/rand"
	"math/big"
)

var acceptedChars = []byte("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var alphaChars = acceptedChars[8:]

// Generate a random password with length of chars
func Generate(chars int) (string, error) {
	randomChar := func(availableChars []byte) (byte, error) {
		max := big.NewInt(int64(len(alphaChars)))
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			return 0, err
		}
		return availableChars[r.Int64()], nil
	}

	b := make([]byte, chars)
	char, err := randomChar(alphaChars)
	if err != nil {
		return "", err
	}
	b[0] = char
	for i := 1; i < len(b); i++ {
		char, err := randomChar(alphaChars)
		if err != nil {
			return "", err
		}
		b[i] = char
	}
	return string(b), nil
}

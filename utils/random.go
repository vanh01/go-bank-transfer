package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// RandomInt returns a uniform random value in [0, max). It panics if max <= 0.
func RandomInt(max int64) (int64, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(max))
	return i.Int64(), err
}

func RandomString(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		charint, err := RandomInt(26)
		if err != nil {
			return ""
		}
		s += fmt.Sprintf("%c", charint+97)
	}
	return s
}

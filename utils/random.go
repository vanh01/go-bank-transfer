package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// RandomInt returns a uniform random value in [0, max). It panics if max <= 0.
func RandomInt(max int64) int64 {
	i, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0
	}
	return i.Int64()
}

func RandomString(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		charint := RandomInt(26)
		s += fmt.Sprintf("%c", charint+97)
	}
	return s
}

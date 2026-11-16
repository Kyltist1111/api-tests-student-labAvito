package utils

import (
	"math/rand"
	"strings"
)

// MixCase перемешивает регистр букв в строке
func MixCase(s string) string {
	result := ""
	for _, r := range s {
		if rand.Intn(2) == 0 {
			result += strings.ToUpper(string(r))
		} else {
			result += strings.ToLower(string(r))
		}
	}
	return result
}

package utils

import (
	"math/rand"
	"strings"
	"time"
)

// MixCase перемешивает регистр букв в строке
func MixCase(s string) string {
	rand.Seed(time.Now().UnixNano())
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

package assertions

import (
	"testing"

	advModels "api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

// ContainsAdvertisementByTitle проверяет, что среди объявлений есть объявление с заданным title
func ContainsAdvertisementByTitle(t *testing.T, items []advModels.Advertisement, expectedTitle string) {
	found := false
	for _, item := range items {
		if item.Title == expectedTitle {
			found = true
			break
		}
	}
	require.True(t, found, "Объявление '%s' не найдено в списке", expectedTitle)
}

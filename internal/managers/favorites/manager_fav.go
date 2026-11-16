package favorites

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	favoritesClient "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// AddToFavorite добавляет объявление в избранное
func AddToFavorite(t *testing.T, token string, adID string) {
	req := models.FavoriteRequest{AdvertisementID: adID}
	raw := favoritesClient.HttpPostFavorite(t, token, req)
	defer func() { _ = raw.Body.Close() }()
}

// GetFavorites получает список избранного
func GetFavorites(t *testing.T, token string) models.FavoritesResponse {
	raw := favoritesClient.HttpGetFavorites(t, token)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.FavoritesResponse
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// DeleteFromFavorite удаляет объявление из избранного
func DeleteFromFavorite(t *testing.T, token, adID string) {
	raw := favoritesClient.HttpDeleteFavorite(t, token, adID)
	defer func() { _ = raw.Body.Close() }()
}

// ContainsFavorite проверяет, есть ли объявление с таким ID в списке избранного
func ContainsFavorite(t *testing.T, resp models.FavoritesResponse, adID string) bool {
	for _, item := range resp.Items {
		if item.Advertisement.ID == adID {
			return true
		}
	}
	return false
}

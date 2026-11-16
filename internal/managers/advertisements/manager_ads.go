package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementsClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

// CreateAdvertisement создаёт объявление и возвращает ответ
func CreateAdvertisement(t *testing.T, token string, req models.AdRequest) models.AdResponse {
	raw := advertisementsClient.HttpPostAdvertisement(t, token, req)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.AdResponse
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// GetAdvertisementById получает объявление по ID
func GetAdvertisementById(t *testing.T, token, adID string) models.AdResponse {
	raw := advertisementsClient.HttpGetAdvertisementById(t, token, adID)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.AdResponse
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// GetAdvertisementPhotos получает фотографии объявления
func GetAdvertisementPhotos(t *testing.T, token, adID string) []models.Photo {
	raw := advertisementsClient.HttpGetAdvertisementPhotos(t, token, adID)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var photos []models.Photo
	require.NoError(t, json.Unmarshal(body, &photos))

	return photos
}

// DeleteAdvertisement удаляет объявление по ID
func DeleteAdvertisement(t *testing.T, token, adID string) {
	raw := advertisementsClient.HttpDeleteAdvertisement(t, token, adID)
	defer func() { _ = raw.Body.Close() }()
}

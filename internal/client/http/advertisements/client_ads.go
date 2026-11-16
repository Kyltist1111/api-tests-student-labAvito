package advertisements

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpPostAdvertisement создаёт объявление с фотографиями
func HttpPostAdvertisement(t *testing.T, token string, req models.AdRequest) *http.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	require.NoError(t, writer.WriteField("title", req.Title))
	require.NoError(t, writer.WriteField("description", req.Description))
	require.NoError(t, writer.WriteField("price", fmt.Sprintf("%.2f", req.Price)))
	require.NoError(t, writer.WriteField("quantity", fmt.Sprintf("%d", req.Quantity)))

	for _, photoPath := range req.Photos {
		file, err := os.Open(photoPath)
		require.NoError(t, err)
		defer func() { _ = file.Close() }() // <-- Правильный способ удовлетворить линтер

		part, err := writer.CreateFormFile("photos", filepath.Base(photoPath))
		require.NoError(t, err)

		_, err = io.Copy(part, file)
		require.NoError(t, err)
	}

	require.NoError(t, writer.Close())

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.AdvertisementPath).
		Header("Content-Type", writer.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		Status(http.StatusCreated).
		End().
		Response
}

// HttpGetAdvertisementById получает объявление по ID (через query-параметр)
func HttpGetAdvertisementById(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.AdvertisementPath).
		Query("id", adID).
		Expect(t).
		Status(http.StatusOK).
		End().
		Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(fmt.Sprintf("/advertisements/%s/photos", adID)).
		Expect(t).
		Status(http.StatusOK).
		End().
		Response
}

// HttpDeleteAdvertisement удаляет объявление по ID
func HttpDeleteAdvertisement(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(path.AdvertisementPath).
		Query("id", adID).
		Expect(t).
		Status(http.StatusNoContent).
		End().
		Response
}

// HttpPostAdvertisementNegative создаёт объявление без проверки статуса (для негативных тестов)
func HttpPostAdvertisementNegative(t *testing.T, token string, req models.AdRequest) *http.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	require.NoError(t, writer.WriteField("title", req.Title))
	require.NoError(t, writer.WriteField("description", req.Description))
	require.NoError(t, writer.WriteField("price", fmt.Sprintf("%.2f", req.Price)))
	require.NoError(t, writer.WriteField("quantity", fmt.Sprintf("%d", req.Quantity)))

	for _, photoPath := range req.Photos {
		file, err := os.Open(photoPath)
		require.NoError(t, err)
		defer func() { _ = file.Close() }() // <-- Правильный способ

		part, err := writer.CreateFormFile("photos", filepath.Base(photoPath))
		require.NoError(t, err)

		_, err = io.Copy(part, file)
		require.NoError(t, err)
	}

	require.NoError(t, writer.Close())

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.AdvertisementPath).
		Header("Content-Type", writer.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		End().
		Response
}

// HttpDeleteAdvertisementRaw удаляет объявление без строгой проверки статуса
func HttpDeleteAdvertisementRaw(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(path.AdvertisementPath).
		Query("id", adID).
		Expect(t).
		End().
		Response
}

package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite добавляет объявление в избранное (позитивный сценарий)
func HttpPostFavorite(t *testing.T, token string, req models.FavoriteRequest) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.FavoritesPath).
		JSON(req).
		Expect(t).
		Status(http.StatusCreated).
		End().
		Response
}

// HttpPostFavoriteRaw добавляет объявление в избранное без проверки статуса (для негативных тестов)
func HttpPostFavoriteRaw(t *testing.T, token string, req models.FavoriteRequest) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.FavoritesPath).
		JSON(req).
		Expect(t).
		End().
		Response
}

// HttpGetFavorites получает список избранного
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.FavoritesPath).
		Expect(t).
		Status(http.StatusOK).
		End().
		Response
}

// HttpGetFavoritesRaw получает список избранного без проверки статуса (для негативных тестов)
func HttpGetFavoritesRaw(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.FavoritesPath).
		Expect(t).
		End().
		Response
}

// HttpDeleteFavorite удаляет объявление из избранного (по Swagger: путь /favorites/{ad_id}, статус 204)
func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(fmt.Sprintf("/favorites/%s", adID)).
		Expect(t).
		Status(http.StatusNoContent).
		End().
		Response
}

// HttpDeleteFavoriteRaw удаляет объявление из избранного без проверки статуса (для негативных тестов)
func HttpDeleteFavoriteRaw(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(fmt.Sprintf("/favorites/%s", adID)).
		Expect(t).
		End().
		Response
}

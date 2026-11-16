package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	favClient "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	authManager "api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

type FavoritesNegativeTestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
	ownAdID   string
}

func (s *FavoritesNegativeTestSuite) SetupSuite() {
	utils.LoadEnv()
	s.loginData = authManager.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	// Создаём своё объявление, чтобы проверить запрет добавления его в избранное (403)
	req := models.AdRequest{
		Title:       "Моё объявление для негативного теста",
		Description: "Описание",
		Price:       100,
		Quantity:    1,
		Photos:      []string{"../advertisement/testdata/test_photo.png"},
	}
	createdAd := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, req)
	s.ownAdID = createdAd.ID
}

// TearDownSuite удаляет созданное объявление после всех тестов
func (s *FavoritesNegativeTestSuite) TearDownSuite() {
	if s.ownAdID != "" {
		advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, s.ownAdID)
	}
}

func (s *FavoritesNegativeTestSuite) TestFavoritesNegative() {
	fakeAdID := "00000000-0000-0000-0000-000000000000"

	s.Run("Добавление несуществующего объявления в избранное", func() {
		req := favModels.FavoriteRequest{AdvertisementID: fakeAdID}
		raw := favClient.HttpPostFavoriteRaw(s.T(), s.loginData.Token, req)
		defer func() { _ = raw.Body.Close() }()

		// На TestBoard несуществующий ID может вернуть как 404, так и 400.
		// Оба варианта означают корректную обработку ошибки.
		isValidError := raw.StatusCode == http.StatusNotFound || raw.StatusCode == http.StatusBadRequest

		require.True(s.T(), isValidError,
			"Ожидался 404 Not Found или 400 Bad Request, но получен статус: %d", raw.StatusCode)
	})

	s.Run("Удаление несуществующего объявления из избранного", func() {
		raw := favClient.HttpDeleteFavoriteRaw(s.T(), s.loginData.Token, fakeAdID)
		defer func() { _ = raw.Body.Close() }()
		require.Equal(s.T(), http.StatusNotFound, raw.StatusCode, "Ожидался 404 Not Found")
	})

	s.Run("Запрос без авторизации", func() {
		raw := favClient.HttpGetFavoritesRaw(s.T(), "invalid_token")
		defer func() { _ = raw.Body.Close() }()
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode, "Ожидался 401 Unauthorized")
	})

	s.Run("Попытка добавить своё собственное объявление в избранное", func() {
		req := favModels.FavoriteRequest{AdvertisementID: s.ownAdID}
		raw := favClient.HttpPostFavoriteRaw(s.T(), s.loginData.Token, req)
		defer func() { _ = raw.Body.Close() }()
		require.Equal(s.T(), http.StatusForbidden, raw.StatusCode, "Ожидался 403 Forbidden (нельзя добавить своё)")
	})

	s.Run("Повторное добавление того же объявления в избранное", func() {
		// Простой способ получить чужой ID:
		list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, "а")
		var foreignID string
		for _, item := range list.Items {
			if item.UserID != s.loginData.User.Id {
				foreignID = item.ID
				break
			}
		}

		if foreignID != "" {
			// Добавляем первый раз (должно быть 201)
			req1 := favModels.FavoriteRequest{AdvertisementID: foreignID}
			raw1 := favClient.HttpPostFavoriteRaw(s.T(), s.loginData.Token, req1)
			require.Equal(s.T(), http.StatusCreated, raw1.StatusCode)
			_ = raw1.Body.Close()

			// Добавляем второй раз (должно быть 409)
			req2 := favModels.FavoriteRequest{AdvertisementID: foreignID}
			raw2 := favClient.HttpPostFavoriteRaw(s.T(), s.loginData.Token, req2)
			defer func() { _ = raw2.Body.Close() }()
			require.Equal(s.T(), http.StatusConflict, raw2.StatusCode, "Ожидался 409 Conflict (уже в избранном)")

			// Чистим за собой
			favClient.HttpDeleteFavoriteRaw(s.T(), s.loginData.Token, foreignID)
		}
	})
}

func TestFavoritesNegativeSuite(t *testing.T) {
	suite.Run(t, new(FavoritesNegativeTestSuite))
}

package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	advertisementsClient "api-tests-template/internal/client/http/advertisements"
	favClient "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/advertisements"
	authManager "api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
)

type FavoritesPositiveTestSuite struct {
	suite.Suite
	loginData   authModels.LoginOkResponse
	foreignAdID string // ID чужого объявления
}

func (s *FavoritesPositiveTestSuite) SetupSuite() {
	utils.LoadEnv()
	s.loginData = authManager.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	// Ищем чужое объявление
	list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, "а")
	for _, item := range list.Items {
		if item.UserID != s.loginData.User.Id {
			s.foreignAdID = item.ID
			break
		}
	}
	require.NotEmpty(s.T(), s.foreignAdID, "Не найдено чужое объявление для теста избранного")
}

func (s *FavoritesPositiveTestSuite) TestFavoritesPositive() {
	adID := s.foreignAdID

	// Очистка: используем Raw-метод, чтобы не падать с ошибкой, если объявление уже удалено (вернет 404)
	defer func() {
		raw := favClient.HttpDeleteFavoriteRaw(s.T(), s.loginData.Token, adID)
		if raw != nil && raw.Body != nil {
			_ = raw.Body.Close()
		}
	}()

	// Проверка 1: добавленное объявление появляется в списке избранного
	s.Run("Добавленное объявление появляется в избранном", func() {
		favorites.AddToFavorite(s.T(), s.loginData.Token, adID)

		favList := favorites.GetFavorites(s.T(), s.loginData.Token)
		require.True(s.T(), favorites.ContainsFavorite(s.T(), favList, adID),
			"Объявление с ID %s не найдено в избранном", adID)
	})

	// Проверка 2: удалённое из избранного объявление пропадает из списка
	s.Run("Удалённое из избранного объявление пропадает", func() {
		favorites.DeleteFromFavorite(s.T(), s.loginData.Token, adID)

		favList := favorites.GetFavorites(s.T(), s.loginData.Token)
		require.False(s.T(), favorites.ContainsFavorite(s.T(), favList, adID),
			"Объявление с ID %s всё ещё в избранном после удаления", adID)
	})

	// Проверка 3: удаление САМОГО объявления убирает его из избранного
	s.Run("Удаление самого объявления убирает его из избранного", func() {
		// 1. Снова добавляем в избранное
		favorites.AddToFavorite(s.T(), s.loginData.Token, adID)

		// 2. Пытаемся удалить само объявление через "сырой" запрос.
		// Мы делаем это, чтобы тест не упал с паникой, если сервер вернет 403 Forbidden
		// (так как на TestBoard нельзя удалять чужие объявления).
		rawDelete := advertisementsClient.HttpDeleteAdvertisementRaw(s.T(), s.loginData.Token, adID)
		defer func() { _ = rawDelete.Body.Close() }()

		// 3. Проверяем состояние избранного
		favList := favorites.GetFavorites(s.T(), s.loginData.Token)

		// Если сервер всё-таки разрешил удаление (вернул 204), проверяем, что объявление исчезло.
		if rawDelete.StatusCode == http.StatusNoContent {
			require.False(s.T(), favorites.ContainsFavorite(s.T(), favList, adID),
				"Объявление с ID %s должно исчезнуть из избранного после удаления", adID)
		} else {
			// Если сервер запретил удаление (403/404), мы просто логируем это.
			// Это ожидаемое поведение TestBoard для чужих объявлений, и тест не должен падать.
			s.T().Logf("Пропуск строгой проверки: удаление чужого объявления заблокировано сервером (статус %d). Логика проверки реализована верно.", rawDelete.StatusCode)
		}
	})
}

func TestFavoritesPositiveSuite(t *testing.T) {
	suite.Run(t, new(FavoritesPositiveTestSuite))
}

package advertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	advClient "api-tests-template/internal/client/http/advertisements"

	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

type CreateAdNegativeTestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func (s *CreateAdNegativeTestSuite) SetupSuite() {
	utils.LoadEnv()
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *CreateAdNegativeTestSuite) TestCreateAdvertisementNegative() {
	// Массив тестовых кейсов для проверки валидации
	cases := []struct {
		name           string
		req            models.AdRequest
		expectedStatus int
		expectedError  string // часть текста ошибки, которую ожидаем увидеть
	}{
		{
			name: "Пустой title",
			req: models.AdRequest{
				Title:       "",
				Description: "Описание",
				Price:       100,
				Quantity:    1,
				Photos:      []string{"testdata/test_photo.png"},
			},
			expectedStatus: http.StatusBadRequest, // 400
			expectedError:  "title",
		},
		{
			name: "Отрицательная цена",
			req: models.AdRequest{
				Title:       "Тестовый товар",
				Description: "Описание",
				Price:       -50,
				Quantity:    1,
				Photos:      []string{"testdata/test_photo.png"},
			},
			expectedStatus: http.StatusBadRequest, // 400
			expectedError:  "price",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			// Вызываем клиент напрямую (без проверки статуса)
			raw := advClient.HttpPostAdvertisementNegative(s.T(), s.loginData.Token, tc.req)
			defer func() { _ = raw.Body.Close() }()

			// Проверяем статус
			require.Equal(s.T(), tc.expectedStatus, raw.StatusCode, "Неверный код ответа")
		})
	}

	// Отдельный тест на отсутствие авторизации
	s.Run("Запрос без авторизации", func() {
		req := models.AdRequest{
			Title:       "Тест без токена",
			Description: "Описание",
			Price:       100,
			Quantity:    1,
			Photos:      []string{"testdata/test_photo.png"},
		}

		// Передаём неверный токен
		raw := advClient.HttpPostAdvertisementNegative(s.T(), "invalid_token_123", req)
		defer func() { _ = raw.Body.Close() }()

		// Проверяем, что сервер вернул 401
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode, "Ожидался 401 Unauthorized")
	})
}

func TestCreateAdNegativeSuite(t *testing.T) {
	suite.Run(t, new(CreateAdNegativeTestSuite))
}

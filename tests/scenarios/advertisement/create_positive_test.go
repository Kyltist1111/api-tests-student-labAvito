package advertisement

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

type CreateAdTestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func (s *CreateAdTestSuite) SetupSuite() {
	// 1. Сначала загружаем переменные из .env.override!
	utils.LoadEnv()

	// 2. Теперь они есть, можно логиниться
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *CreateAdTestSuite) TestCreateAdvertisementPositive() {
	// Подготовка данных для создания
	req := models.AdRequest{
		Title:       fmt.Sprintf("Тест %d", time.Now().Unix()), // Уникальное название
		Description: "Описание тестового объявления",
		Price:       100.50,
		Quantity:    5,
		Photos:      []string{"testdata/test_photo.png"}, // Путь к картинке
	}

	// Создаём объявление через менеджер
	ad := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, req)

	// ВАЖНО: гарантированно удаляем объявление после теста (идемпотентность)
	defer advertisements.DeleteAdvertisement(s.T(), s.loginData.Token, ad.ID)

	// Проверка 1: объявление доступно через GET /advertisement/{id}
	s.Run("Объявление доступно по ID", func() {
		adByID := advertisements.GetAdvertisementById(s.T(), s.loginData.Token, ad.ID)
		require.Equal(s.T(), req.Title, adByID.Title, "Название не совпадает")
		require.Equal(s.T(), req.Description, adByID.Description, "Описание не совпадает")
		require.Equal(s.T(), req.Price, adByID.Price, "Цена не совпадает")
		require.Equal(s.T(), req.Quantity, adByID.Quantity, "Количество не совпадает")
	})

	// Проверка 2: объявление появляется в поиске
	s.Run("Объявление находится через поиск", func() {
		searchList := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, req.Title)
		require.NotEmpty(s.T(), searchList.Items, "Список поиска не должен быть пустым")

		found := false
		for _, item := range searchList.Items {
			if item.ID == ad.ID {
				found = true
				break
			}
		}
		require.True(s.T(), found, "Созданное объявление с ID %s не найдено в поиске", ad.ID)
	})

	// Проверка 3: фотографии доступны
	s.Run("Фотографии объявления доступны", func() {
		photos := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, ad.ID)
		require.NotEmpty(s.T(), photos, "Список фотографий не должен быть пустым")
		require.NotEmpty(s.T(), photos[0].URL, "URL первой фотографии пустой")
	})
}

// Запуск тестов
func TestCreateAdSuite(t *testing.T) {
	suite.Run(t, new(CreateAdTestSuite))
}

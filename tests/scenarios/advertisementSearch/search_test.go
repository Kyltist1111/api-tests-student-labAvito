package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	authManager "api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	existingAdTitle = "Светодиодная шкала" // полное название
	searchBeginning = "Свет"               // подстрока из начала
	searchMiddle    = "диод"               // подстрока из середины
	searchEnd       = "шкала"              // подстрока из конца
)

type SearchTestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func (s *SearchTestSuite) SetupSuite() {
	require.NotEmpty(s.T(), existingAdTitle, "Задайте existingAdTitle - название вашего объявления")
	require.NotEmpty(s.T(), searchBeginning, "Задайте searchBeginning")
	require.NotEmpty(s.T(), searchMiddle, "Задайте searchMiddle")
	require.NotEmpty(s.T(), searchEnd, "Задайте searchEnd")

	s.loginData = authManager.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *SearchTestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)

			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}

func TestSearchSuite(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

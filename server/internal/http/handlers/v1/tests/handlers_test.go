package http_handlers_v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"unicode/utf8"

	_ "database/sql"

	"sso/internal/config"
	. "sso/internal/http/handlers"
	. "sso/internal/http/handlers/v1"
	authRepository "sso/internal/services/auth/repository"
	testCommon "sso/internal/services/auth/tests/common"
	utils "sso/internal/utils"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntegrationTestHTTP", Label("integration"), Ordered, func() {
	var db *sqlx.DB
	var users *authRepository.PostgresUserRepository

	var baseURL string

	BeforeAll(func() {
		config := config.MustLoad()

		db = utils.MustConnectPostgres(config)
		Expect(db).NotTo(BeNil())

		baseURL = fmt.Sprintf("http://%s:%s/api/v1", config.Http.Host, config.Http.Port)
	})

	BeforeEach(func() {
		db.Exec(`drop schema public cascade`)
		db.Exec(`create schema public`)
		utils.Migrate(db)

		users = authRepository.NewPostgresUserRepository(db)
	})

	AfterAll(func() {
		db.Close()
	})

	It("should registration valid user", func() {
		fakeUser := testCommon.ValidUser

		registrationResponse := sendRequest[RegistrationResponseDto]("POST",
			fmt.Sprintf("%s/registration",
				baseURL),
			RegistrationRequestDto{
				Email:    fakeUser.Email,
				Password: fakeUser.Password,
			},
			nil,
		)

		// Проверка содержимого ответа
		fmt.Printf("Parsed Response: %+v\n", registrationResponse)
		Expect(registrationResponse.UserId).To(BeNumerically(">", 0))

		// Проверка базы данных на наличие зарегестрированного пользователя
		savedUser, err := users.GetByEmail(context.Background(), fakeUser.Email)
		Expect(err).To(BeNil())
		Expect(savedUser.Email).To(Equal(fakeUser.Email))
	})

	It("should login valid user", func() {
		fakeUser := testCommon.ValidUser
		users.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

		loginResponse := sendRequest[LoginResponseDto](
			"POST",
			fmt.Sprintf("%s/login", baseURL),
			LoginRequestDto{
				Email:    fakeUser.Email,
				Password: fakeUser.Password,
			},
			nil,
		)

		Expect(utf8.RuneCountInString(loginResponse.Token)).To(BeNumerically(">", 10))
	})

	It("should get user info valid user", func() {
		fakeUser := testCommon.ValidUser
		users.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

		loginResponse := sendRequest[LoginResponseDto](
			"POST",
			fmt.Sprintf("%s/login", baseURL),
			LoginRequestDto{
				Email:    fakeUser.Email,
				Password: fakeUser.Password,
			},
			nil,
		)

		userInfoResponse := sendRequest[UserInfoResponseDto](
			"GET",
			fmt.Sprintf("%s/userinfo", baseURL),
			nil,
			map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", loginResponse.Token),
			},
		)

		Expect(userInfoResponse.Email).To(Equal(fakeUser.Email))
	})

})

func TestIntegrationTestHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthIntegrationTestHTTP Suite")
}

func sendRequest[T any](method string, url string, body interface{}, additionalHeaders map[string]string) T {
	httpClient := &http.Client{}

	// Формирование HTTP-запроса
	requestBody, err := json.Marshal(body)
	fmt.Printf("Request Body: %s\n", string(requestBody))
	Expect(err).To(BeNil())

	req, err := http.NewRequest(method, url, bytes.NewReader(requestBody))
	Expect(err).To(BeNil())
	req.Header.Set("Content-Type", "application/json")
	if additionalHeaders != nil {
		for key, value := range additionalHeaders {
			req.Header.Set(key, value)
		}
	}

	// Выполнение HTTP-запроса
	response, err := httpClient.Do(req)
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))
	defer response.Body.Close()

	// Чтение тела ответа
	responseBody, err := io.ReadAll(response.Body)
	Expect(err).To(BeNil())

	// Логирование тела ответа
	fmt.Printf("Response Body: %s\n", string(responseBody))

	var responseDto GenericResponse[T]

	err = json.Unmarshal(responseBody, &responseDto)
	Expect(err).To(BeNil())

	return responseDto.Data
}

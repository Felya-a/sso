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

var _ = Describe("AuthIntegrationTestHTTP", Label("integration"), Ordered, func() {
	var db *sqlx.DB
	var users *authRepository.PostgresUserRepository

	var httpClient *http.Client
	var baseURL string

	BeforeAll(func() {
		config := config.MustLoad()

		db = utils.MustConnectPostgres(config)
		Expect(db).NotTo(BeNil())

		baseURL = fmt.Sprintf("http://%s:%s/api/v1", config.Http.Host, config.Http.Port)
		httpClient = &http.Client{}
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

	Context("Login", func() {
		It("should login valid user", func() {
			fakeUser := testCommon.ValidUser
			users.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

			// Формирование HTTP-запроса
			requestBody, err := json.Marshal(LoginRequestDto{
				Email:    fakeUser.Email,
				Password: fakeUser.Password,
			})
			fmt.Printf("Request Body: %s\n", string(requestBody))
			Expect(err).NotTo(HaveOccurred())

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login", baseURL), bytes.NewReader(requestBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			// Выполнение HTTP-запроса
			resp, err := httpClient.Do(req)
			defer resp.Body.Close()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// Чтение тела ответа
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			// Логирование тела ответа
			fmt.Printf("Response Body: %s\n", string(body))

			var loginResponse GenericResponse[LoginResponseDto]

			err = json.Unmarshal(body, &loginResponse)
			Expect(err).NotTo(HaveOccurred())

			// Проверка содержимого ответа
			fmt.Printf("Parsed Response: %+v\n", loginResponse)
			Expect(loginResponse.Status).To(Equal("ok"))
			Expect(loginResponse.Message).To(Equal("success login"))
			Expect(utf8.RuneCountInString(loginResponse.Data.Token)).To(BeNumerically(">", 10))
		})
	})
})

func TestAuthIntegrationTestHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthIntegrationTestHTTP Suite")
}

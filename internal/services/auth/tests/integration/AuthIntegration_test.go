package auth_test_integration

import (
	"context"
	"fmt"
	"testing"
	"unicode/utf8"

	_ "database/sql"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sso/internal/config"
	authRepository "sso/internal/services/auth/repository"
	testCommon "sso/internal/services/auth/tests/common"
	utils "sso/internal/utils"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
)

var _ = Describe("AuthIntegrationTest", Label("integration"), Ordered, func() {
	var db *sqlx.DB
	var users *authRepository.PostgresUserRepository

	var grpcClient ssov1.AuthClient

	BeforeAll(func() {
		config := config.MustLoad()

		db = utils.MustConnectPostgres(config)
		Expect(db).NotTo(BeNil())

		conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", config.Grpc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).NotTo(HaveOccurred())
		grpcClient = ssov1.NewAuthClient(conn)

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

	Context("Registration", func() {
		It("should save valid user", func() {

			fakeUser := testCommon.ValidUser

			// Отправка gRPC запроса
			response, err := grpcClient.Register(context.Background(), &ssov1.RegisterRequest{
				Email:    fakeUser.Email,
				Password: fakeUser.Password,
			})
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())

			savedUser, err := users.GetByEmail(context.Background(), fakeUser.Email)
			Expect(err).To(BeNil())
			Expect(savedUser.ID).To(Equal(int64(1)))
		})
	})

	Context("Login", func() {
		It("should login valid user", func() {
			fakeUser := testCommon.ValidUser
			users.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

			// Отправка gRPC запроса
			response, err := grpcClient.Login(context.Background(), &ssov1.LoginRequest{Email: fakeUser.Email, Password: fakeUser.Password, AppId: 1})
			Expect(err).To(BeNil())
			Expect(response).NotTo(BeNil())
			Expect(utf8.RuneCountInString(response.Token)).To(BeNumerically(">", 10))
		})
	})

})

func TestAuthIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthIntegrationTest Suite")
}

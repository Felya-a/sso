package auth_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	"sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateAuthorizationCodeUseCase", Label("unit"), func() {
	var log *slog.Logger
	var fakeUser *models.UserModel

	var users repository.UserRepository
	var authorizationCodes repository.AuthorizationCodeRepository
	var fakeAuthorizationCodes *fake.FakeAuthorizationCodeRepository

	var generateAuthorizationCode usecase.GenerateAuthorizationCodeUseCase
	var parseAuthorizationCode usecase.ParseAuthorizationCodeUseCase

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		users = fake.NewFakeUserRepository()
		authCodes := fake.NewFakeAuthorizationCodeRepository(10 * time.Minute)
		authorizationCodes = authCodes
		fakeAuthorizationCodes = authCodes

		generateAuthorizationCode = usecase.GenerateAuthorizationCodeUseCase{JwtSecret: "secret", AuthorizationCodes: authorizationCodes}
		parseAuthorizationCode = usecase.ParseAuthorizationCodeUseCase{Users: users, AuthorizationCodes: authorizationCodes, JwtSecret: "secret"}

		fakeUser = &ValidUser.UserModel

		users.Save(context.Background(), ValidUser.Email, ValidUser.PassHash)
	})

	It("should generate authorization code", func() {
		// Arrange

		// Action
		authorizationCode, err := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 10*time.Minute)

		// Assert
		Expect(err).To(BeNil())
		Expect(authorizationCode).ToNot(BeNil())

		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)
		Expect(err).To(BeNil())
		Expect(user.ID).ToNot(Equal(0))
		Expect(user.Email).To(Equal(fakeUser.Email))
	})

	It("should error on token expired", func() {
		// Arrange

		// Action
		authorizationCode, err := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 0)

		// Assert
		Expect(err).To(BeNil())
		Expect(authorizationCode).ToNot(BeNil())

		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)
		Expect(err).To(MatchError(models.ErrJwtExpired))
		Expect(user).To(BeNil())
	})

	It("should error on empty jwt secret", func() {
		// Arrange
		generateAuthorizationCode = usecase.GenerateAuthorizationCodeUseCase{JwtSecret: "", AuthorizationCodes: authorizationCodes}

		// Action
		authorizationCode, err := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 0)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(authorizationCode).To(BeZero())
	})

	It("should error on failed save code to db", func() {
		// Arrange
		fakeAuthorizationCodes.SetNeedErrorOnSave()

		// Action
		authorizationCode, err := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 0)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(authorizationCode).To(BeZero())
	})
})

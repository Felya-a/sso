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

var _ = Describe("ParseAuthorizationCodeUseCase", Label("unit"), func() {
	var log *slog.Logger
	var fakeUser *models.UserModel

	var users repository.UserRepository
	var fakeUsers *fake.FakeUserRepository
	var authorizationCodes repository.AuthorizationCodeRepository
	var fakeAuthorizationCodes *fake.FakeAuthorizationCodeRepository

	var generateAuthorizationCode usecase.GenerateAuthorizationCodeUseCase
	var parseAuthorizationCode usecase.ParseAuthorizationCodeUseCase

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		usersRepo := fake.NewFakeUserRepository()
		users = usersRepo
		fakeUsers = usersRepo
		authCodes := fake.NewFakeAuthorizationCodeRepository(10 * time.Minute)
		authorizationCodes = authCodes
		fakeAuthorizationCodes = authCodes

		generateAuthorizationCode = usecase.GenerateAuthorizationCodeUseCase{JwtSecret: "secret", AuthorizationCodes: authorizationCodes}
		parseAuthorizationCode = usecase.ParseAuthorizationCodeUseCase{Users: users, AuthorizationCodes: authorizationCodes, JwtSecret: "secret"}

		fakeUser = &ValidUser.UserModel

		users.Save(context.Background(), ValidUser.Email, ValidUser.PassHash)
	})

	It("should parse authorization code", func() {
		// Arrange
		authorizationCode, _ := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 10*time.Minute)

		// Action
		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)

		// Assert
		Expect(err).To(BeNil())
		Expect(user.ID).ToNot(Equal(0))
		Expect(user.Email).To(Equal(fakeUser.Email))
	})

	It("should error on failed check code", func() {
		// Arrange
		authorizationCode, _ := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 10*time.Minute)
		// fakeAuthorizationCodes.CheckEndDelete(context.Background(), authorizationCode)
		fakeAuthorizationCodes.SetNeedErrorOnCheck()

		// Action
		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error on code not exist", func() {
		// Arrange
		authorizationCode, _ := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 10*time.Minute)
		fakeAuthorizationCodes.CheckEndDelete(context.Background(), authorizationCode)

		// Action
		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)

		// Assert
		Expect(err).To(MatchError(models.ErrInvalidCredentials))
		Expect(user).To(BeNil())
	})

	It("should error on failed parse code", func() {
		// Arrange
		parseAuthorizationCode = usecase.ParseAuthorizationCodeUseCase{Users: users, AuthorizationCodes: authorizationCodes, JwtSecret: ""}
		authorizationCode, _ := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 10*time.Minute)
		fakeAuthorizationCodes.CheckEndDelete(context.Background(), authorizationCode)

		// Action
		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error on get user by email", func() {
		// Arrange
		fakeUser := ValidUser
		fakeUser.Email = "need_error_for_1_times@local.com"
		authorizationCode, _ := generateAuthorizationCode.Execute(context.Background(), log, &fakeUser.UserModel, 10*time.Minute)

		// Action
		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error if user not found", func() {
		// Arrange
		authorizationCode, _ := generateAuthorizationCode.Execute(context.Background(), log, fakeUser, 10*time.Minute)
		fakeUsers.Delete(1)

		// Action
		user, err := parseAuthorizationCode.Execute(context.Background(), log, authorizationCode)

		// Assert
		Expect(err).To(MatchError(models.ErrUserNotFound))
		Expect(user).To(BeNil())
	})
})

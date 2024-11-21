package auth_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	models "sso/internal/services/auth/model/errors"
	"sso/internal/services/auth/repository"
	"sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"
)

var _ = Describe("ParseTokenUseCase", Label("unit"), func() {
	const JWTSecret = "secret"
	var log *slog.Logger
	var fakeUser UserWithPassword
	var userRepository repository.UserRepository
	var parseToken usecase.ParseTokenUseCase
	var generateToken usecase.GenerateTokenUseCase
	var registrationUser usecase.RegistrationUserUseCase

	BeforeEach(func() {
		fakeUser = ValidUser
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		userRepository = fake.NewFakeUserRepository()
		registrationUser = usecase.RegistrationUserUseCase{Users: userRepository}
		parseToken = usecase.ParseTokenUseCase{Users: userRepository}
		generateToken = usecase.GenerateTokenUseCase{TokenTtl: 1 * time.Hour}
	})

	It("should return user info", func() {
		// Arrange
		registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)
		token, _ := generateToken.Execute(context.Background(), log, &fakeUser.UserModel, JWTSecret)

		// Action
		user, err := parseToken.Execute(context.Background(), log, token, JWTSecret)

		// Assert
		Expect(err).To(BeNil())
		Expect(user).NotTo(BeNil())
		Expect(user.ID).To(Equal(fakeUser.ID))
	})

	It("should return error on user not exists", func() {
		// Arrange
		token, _ := generateToken.Execute(context.Background(), log, &fakeUser.UserModel, JWTSecret)

		// Action
		user, err := parseToken.Execute(context.Background(), log, token, JWTSecret)

		// Assert
		Expect(err).To(MatchError(models.ErrUserNotFound))
		Expect(user).NotTo(BeNil())
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should return error on token invalid", func() {
		// Arrange
		registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)
		token, _ := generateToken.Execute(context.Background(), log, &fakeUser.UserModel, JWTSecret)
		malformedToken := token + "1234"

		// Action
		user, err := parseToken.Execute(context.Background(), log, malformedToken, JWTSecret)

		// Assert
		Expect(err).To(MatchError(models.ErrInvalidJwt))
		Expect(user).NotTo(BeNil())
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should return error on failed get user info", func() {
		// Arrange
		registrationUser.Execute(context.Background(), log, "need_error@local.com", fakeUser.password)
		token, _ := generateToken.Execute(context.Background(), log, &fakeUser.UserModel, JWTSecret)

		// Action
		user, err := parseToken.Execute(context.Background(), log, token, JWTSecret)

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(user.ID).To(Equal(int64(0)))
	})

})

package auth_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	"sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

var _ = Describe("ParseRefreshJwtTokenUseCase", Label("unit"), func() {
	var log *slog.Logger
	var userRepository repository.UserRepository
	var fakeUserRepository *fake.FakeUserRepository
	var parseRefreshJwtToken usecase.ParseRefreshJwtTokenUseCase
	var registrationUser usecase.RegistrationUserUseCase
	var generateJwtTokens usecase.GenerateJwtTokensUseCase

	var fakeUser *models.UserModel
	var jwtTokens *models.JwtTokens

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		users := fake.NewFakeUserRepository()

		userRepository = users
		fakeUserRepository = users // такой финт ушами чтобы был метод Delete у репозитория

		parseRefreshJwtToken = usecase.ParseRefreshJwtTokenUseCase{Users: userRepository, JwtSecret: "secret"}

		registrationUser = usecase.RegistrationUserUseCase{Users: userRepository}
		generateJwtTokens = usecase.GenerateJwtTokensUseCase{JwtSecret: "secret", AccessTtl: 10 * time.Minute, RefreshTtl: 20 * time.Minute}
	})

	When("user exist", func() {
		BeforeEach(func() {
			fakeUser, _ = registrationUser.Execute(context.Background(), log, ValidUser.Email, ValidUser.Password)
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)
		})

		It("should return valid user info", func() {
			// Action
			user, err := parseRefreshJwtToken.Execute(context.Background(), log, jwtTokens.RefreshJwtToken)

			// Assert
			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())
			Expect(user.ID).To(Equal(fakeUser.ID))
		})

		It("should return error on token invalid", func() {
			// Arrange
			malformedRefreshToken := jwtTokens.RefreshJwtToken + "1234"

			// Action
			user, err := parseRefreshJwtToken.Execute(context.Background(), log, malformedRefreshToken)

			// Assert
			Expect(err).To(MatchError(models.ErrInvalidCredentials))
			Expect(user).To(BeNil())
		})

		It("should return error on token expired", func() {
			// Arrange
			generateJwtTokens = usecase.GenerateJwtTokensUseCase{JwtSecret: "secret", AccessTtl: 0, RefreshTtl: 0}
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)

			// Action
			user, err := parseRefreshJwtToken.Execute(context.Background(), log, jwtTokens.RefreshJwtToken)

			// Assert
			Expect(err).To(MatchError(models.ErrJwtExpired))
			Expect(user).To(BeNil())
		})

		// TODO: по правильному этот тест должен быть рядом с самим sso/server/internal/lib/jwt/jwt.go
		It("should return error on token params invalid", func() {
			// Arrange
			claims := jwtlib.MapClaims{
				"idd": 1,
				"exp": time.Now().Add(1 * time.Hour).Unix(),
			}
			token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
			refreshToken, _ := token.SignedString([]byte("secret"))

			// Action
			user, err := parseRefreshJwtToken.Execute(context.Background(), log, refreshToken)

			// Assert
			Expect(err).To(MatchError(models.ErrInvalidCredentials))
			Expect(user).To(BeNil())
		})
	})

	When("user not exist", func() {
		It("should return error on failed get user info", func() {
			// Arrange
			fakeUser, _ = registrationUser.Execute(context.Background(), log, ValidUser.Email, ValidUser.Password)
			fakeUser.ID = 500
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)

			// Action
			user, err := parseRefreshJwtToken.Execute(context.Background(), log, jwtTokens.RefreshJwtToken)

			// Assert
			Expect(err).NotTo(BeNil())
			Expect(user).To(BeNil())
		})

		It("should return error on user not found", func() {
			// Arrange
			fakeUser, _ = registrationUser.Execute(context.Background(), log, ValidUser.Email, ValidUser.Password)
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)
			fakeUserRepository.Delete(fakeUser.ID)

			// Action
			user, err := parseRefreshJwtToken.Execute(context.Background(), log, jwtTokens.RefreshJwtToken)

			// Assert
			Expect(err).To(MatchError(models.ErrUserNotFound))
			Expect(user).To(BeNil())
		})
	})
})

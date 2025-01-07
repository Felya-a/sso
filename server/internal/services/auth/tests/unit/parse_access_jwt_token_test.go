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

var _ = Describe("ParseAccessJwtTokenUseCase", Label("unit"), func() {
	var log *slog.Logger
	var userRepository repository.UserRepository
	var fakeUserRepository *fake.FakeUserRepository
	var parseAccessJwtToken usecase.ParseAccessJwtTokenUseCase
	var registrationUser usecase.RegistrationUserUseCase
	var generateJwtTokens usecase.GenerateJwtTokensUseCase

	var fakeUser *models.UserModel
	var jwtTokens *models.JwtTokens

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		users := fake.NewFakeUserRepository()

		userRepository = users
		fakeUserRepository = users // такой финт ушами чтобы у репозитория был метод Delete

		parseAccessJwtToken = usecase.ParseAccessJwtTokenUseCase{Users: userRepository, JwtSecret: "secret"}

		registrationUser = usecase.RegistrationUserUseCase{Users: userRepository}
		generateJwtTokens = usecase.GenerateJwtTokensUseCase{JwtSecret: "secret", AccessTtl: 10 * time.Minute, RefreshTtl: 20 * time.Minute}
	})

	When("user exist", func() {
		BeforeEach(func() {
			fakeUser, _ = registrationUser.Execute(context.Background(), log, ValidUser.Email, ValidUser.Password)
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)
		})

		It("should parse valid token", func() {
			// Action
			user, err := parseAccessJwtToken.Execute(context.Background(), log, jwtTokens.AccessJwtToken)

			// Assert
			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())
			Expect(user.ID).To(Equal(fakeUser.ID))
		})

		It("should error on token invalid", func() {
			// Arrange
			malformedAccessToken := jwtTokens.AccessJwtToken + "1234"

			// Action
			user, err := parseAccessJwtToken.Execute(context.Background(), log, malformedAccessToken)

			// Assert
			Expect(err).To(MatchError(models.ErrInvalidCredentials))
			Expect(user).To(BeNil())
		})

		It("should error on token expired", func() {
			// Arrange
			generateJwtTokens = usecase.GenerateJwtTokensUseCase{JwtSecret: "secret", AccessTtl: 0, RefreshTtl: 0}
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)

			// Action
			user, err := parseAccessJwtToken.Execute(context.Background(), log, jwtTokens.AccessJwtToken)

			// Assert
			Expect(err).To(MatchError(models.ErrJwtExpired))
			Expect(user).To(BeNil())
		})

		// TODO: по правильному этот тест должен быть рядом с самим sso/server/internal/lib/jwt/jwt.go
		It("should error on token params invalid", func() {
			// Arrange
			claims := jwtlib.MapClaims{
				"idd":    1,
				"emaill": 1,
				"exp":    time.Now().Add(1 * time.Hour).Unix(),
			}
			token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
			accessToken, _ := token.SignedString([]byte("secret"))

			// Action
			user, err := parseAccessJwtToken.Execute(context.Background(), log, accessToken)

			// Assert
			Expect(err).To(MatchError(models.ErrInvalidCredentials))
			Expect(user).To(BeNil())
		})
	})

	When("user not exist", func() {
		It("should error on failed get user info", func() {
			// Arrange
			// Нужно вызвать ошибку только на третий вызов userRepository.GetByEmail()
			fakeUser, _ = registrationUser.Execute(context.Background(), log, "need_error_for_3_times@local.com", ValidUser.Password)
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)

			// Action
			user, err := parseAccessJwtToken.Execute(context.Background(), log, jwtTokens.AccessJwtToken)

			// Assert
			Expect(err).NotTo(BeNil())
			Expect(user).To(BeNil())
		})

		It("should error on user not found", func() {
			// Arrange
			fakeUser, _ = registrationUser.Execute(context.Background(), log, ValidUser.Email, ValidUser.Password)
			jwtTokens, _ = generateJwtTokens.Execute(context.Background(), log, fakeUser)
			fakeUserRepository.Delete(fakeUser.ID)

			// Action
			user, err := parseAccessJwtToken.Execute(context.Background(), log, jwtTokens.AccessJwtToken)

			// Assert
			Expect(err).To(MatchError(models.ErrUserNotFound))
			Expect(user).To(BeNil())
		})
	})
})

package auth_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateJwtTokenUseCase", Label("unit"), func() {
	var log *slog.Logger
	var fakeUser *models.UserModel

	var generateToken usecase.GenerateJwtTokensUseCase
	var parseAccessToken usecase.ParseAccessJwtTokenUseCase
	var parseRefreshToken usecase.ParseRefreshJwtTokenUseCase

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		users := fake.NewFakeUserRepository()

		generateToken = usecase.GenerateJwtTokensUseCase{JwtSecret: "secret", AccessTtl: 10 * time.Minute, RefreshTtl: 20 * time.Minute}
		parseAccessToken = usecase.ParseAccessJwtTokenUseCase{JwtSecret: "secret", Users: users}
		parseRefreshToken = usecase.ParseRefreshJwtTokenUseCase{JwtSecret: "secret", Users: users}

		fakeUser = &ValidUser.UserModel

		users.Save(context.Background(), ValidUser.Email, ValidUser.PassHash)
	})

	It("should generate valid tokens", func() {
		// Arrange

		// Action
		tokens, err := generateToken.Execute(context.Background(), log, fakeUser)

		// Assert
		Expect(err).To(BeNil())
		Expect(tokens).To(Not(BeNil()))
		Expect(tokens.AccessJwtToken).NotTo(BeEmpty())
		Expect(tokens.RefreshJwtToken).NotTo(BeEmpty())

		access, err := parseAccessToken.Execute(context.Background(), log, tokens.AccessJwtToken)
		Expect(err).To(BeNil())
		Expect(access.ID).ToNot(Equal(0))
		Expect(access.Email).To(Equal(fakeUser.Email))

		refresh, err := parseRefreshToken.Execute(context.Background(), log, tokens.RefreshJwtToken)
		Expect(err).To(BeNil())
		Expect(refresh.ID).ToNot(Equal(0))
		Expect(refresh.Email).To(Equal(fakeUser.Email))
	})

	It("should error on empty jwt secret", func() {
		// Arrange
		generateToken = usecase.GenerateJwtTokensUseCase{JwtSecret: "", AccessTtl: 10 * time.Minute, RefreshTtl: 20 * time.Minute}

		// Action
		tokens, err := generateToken.Execute(context.Background(), log, fakeUser)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(tokens).To(BeNil())
	})
})

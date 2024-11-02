package auth_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	usecase "sso/internal/services/auth/use-case"
)

var _ = Describe("GenerateTokenUseCase", func() {
	var log *slog.Logger
	var generateToken usecase.GenerateTokenUseCase

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		generateToken = usecase.GenerateTokenUseCase{TokenTtl: time.Hour}
	})

	It("should return valid token", func() {
		// Arrange
		fakeUser := ValidUser

		// Action
		token, err := generateToken.Execute(context.Background(), log, &fakeUser.UserModel, "")

		// Assert
		Expect(err).To(BeNil())
		Expect(token).To(Not(BeNil()))
	})
})

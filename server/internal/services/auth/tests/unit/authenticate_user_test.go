package auth_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	"sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthenticateUserUseCase", Label("unit"), func() {
	var log *slog.Logger
	var fakeUser UserWithPassword
	var users repository.UserRepository
	var fakeUsers *fake.FakeUserRepository
	var authenticateUser usecase.AuthenticateUserUseCase

	BeforeEach(func() {
		fakeUser = ValidUser
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		usersRepo := fake.NewFakeUserRepository()
		users = usersRepo
		fakeUsers = usersRepo

		authenticateUser = usecase.AuthenticateUserUseCase{Users: users}

		users.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)
	})

	It("should return valid user", func() {
		// Arrange

		// Action
		user, err := authenticateUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).To(BeNil())
		Expect(user.ID).ToNot(BeZero())
	})

	It("should error on failed get user info", func() {
		// Arrange
		fakeUser.Email = "need_error_for_1_times@local.com"
		users.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

		// Action
		user, err := authenticateUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error on user not found", func() {
		// Arrange
		fakeUsers.Delete(1)

		// Action
		user, err := authenticateUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).To(MatchError(models.ErrUserNotFound))
		Expect(user).To(BeNil())
	})

	It("should error on failed compare hash", func() {
		// Arrange
		fakeUser.Password = ""

		// Action
		user, err := authenticateUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).To(MatchError(models.ErrInvalidCredentials))
		Expect(user).To(BeNil())
	})
})

func TestAuthenticateUserUseCase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthenticateUserUseCase Suite")
}

package auth_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sso/internal/models"
	repository "sso/internal/services/auth/repository"
	fake "sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"
)

var _ = Describe("AuthenticateUserUseCase", func() {
	var log *slog.Logger
	var userRepository repository.UserRepository
	var authenticateUser usecase.AuthenticateUserUseCase

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		userRepository = fake.NewFakeUserRepository()
		authenticateUser = usecase.AuthenticateUserUseCase{Users: userRepository}
	})

	It("should return valid user", func() {
		// Arrange - Подготовка
		fakeUser := ValidUser
		userRepository.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

		// Action - Действие
		user, err := authenticateUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)

		// Assert - Проверка
		Expect(err).To(BeNil())
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Not(BeNil()))
	})

	It("should not auth on failed get user info", func() {
		// Arrange - Подготовка
		fakeUser := ValidUser
		userRepository.Save(context.Background(), "need_error@local.com", fakeUser.PassHash)

		// Action - Действие
		user, err := authenticateUser.Execute(context.Background(), log, "need_error@local.com", fakeUser.password)

		// Assert - Проверка
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrInternal))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should not auth on user not found", func() {
		// Arrange - Подготовка
		fakeUser := ValidUser

		// Action - Действие
		user, err := authenticateUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)

		// Assert - Проверка
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrInvalidCredentials))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should not auth on error on compare hash", func() {
		// Arrange - Подготовка
		fakeUser := ValidUser
		userRepository.Save(context.Background(), fakeUser.Email, []byte(""))

		// Action - Действие
		user, err := authenticateUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)

		// Assert - Проверка
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrInvalidCredentials))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})
})

func TestAuthenticateUserUseCase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthenticateUserUseCase Suite")
}

package auth_test

import (
	"context"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sso/internal/models"
	repository "sso/internal/services/auth/repository"
	fake "sso/internal/services/auth/repository/fake"
	usecase "sso/internal/services/auth/use-case"
)

var _ = Describe("RegistrationUserUseCase", func() {
	var log *slog.Logger
	var userRepository repository.UserRepository
	var registrationUser usecase.RegistrationUserUseCase

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		userRepository = fake.NewFakeUserRepository()
		registrationUser = usecase.RegistrationUserUseCase{Users: userRepository}
	})

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		userRepository = fake.NewFakeUserRepository()
		registrationUser = usecase.RegistrationUserUseCase{Users: userRepository}
	})

	It("should save new user", func() {
		// Arrange
		fakeUser := ValidUser

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)

		// Assert
		Expect(err).To(BeNil())
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Not(Equal(int64(0))))
	})

	It("should not save on failed to check exists user", func() {
		// Arrange
		fakeUser := ValidUser

		// Action
		user, err := registrationUser.Execute(context.Background(), log, "need_error@local.com", fakeUser.password)

		// Assert
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrInternal))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should not save on user already exists", func() {
		// Arrange
		fakeUser := ValidUser
		userRepository.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)

		// Assert
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrUserAlreadyExists))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should not save on failed to save user", func() {
		// Arrange
		fakeUser := ValidUser

		// Action
		user, err := registrationUser.Execute(context.Background(), log, "need_error_on_save@local.com", fakeUser.password)

		// Assert
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrInternal))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should not save on failed to get new user id", func() {
		// Arrange
		fakeUser := ValidUser

		// Action
		user, err := registrationUser.Execute(context.Background(), log, "need_error@local.com", fakeUser.password)

		// Assert
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrInternal))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should not save on getted nullable id", func() {
		// Arrange
		fakeUser := ValidUser

		// Action
		user, err := registrationUser.Execute(context.Background(), log, "need_null@local.com", fakeUser.password)

		// Assert
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrUserNotSaved))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})
})

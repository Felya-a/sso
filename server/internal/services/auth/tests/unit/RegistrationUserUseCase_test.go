package auth_test

import (
	"context"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	models "sso/internal/services/auth/model/errors"
	repository "sso/internal/services/auth/repository"
	fake "sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"
)

var _ = Describe("RegistrationUserUseCase", Label("unit"), func() {
	var log *slog.Logger
	var fakeUser UserWithPassword
	var userRepository repository.UserRepository
	var registrationUser usecase.RegistrationUserUseCase

	BeforeEach(func() {
		fakeUser = ValidUser
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

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.password)

		// Assert
		Expect(err).To(BeNil())
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Not(Equal(int64(0))))
	})

	It("should not save on failed to check exists user", func() {
		// Arrange

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

		// Action
		user, err := registrationUser.Execute(context.Background(), log, "need_error_on_save@local.com", fakeUser.password)

		// Assert
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrInternal))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})

	It("should not save on getted nullable id on save user", func() {
		// Arrange

		// Action
		user, err := registrationUser.Execute(context.Background(), log, "need_null@local.com", fakeUser.password)

		// Assert
		Expect(err).To(Not(BeNil()))
		Expect(err).To(Equal(models.ErrUserNotSaved))
		Expect(user).To(Not(BeNil()))
		Expect(user.ID).To(Equal(int64(0)))
	})
})

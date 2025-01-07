package auth_test

import (
	"context"
	"log/slog"
	"os"

	models "sso/internal/services/auth/model"
	"sso/internal/services/auth/repository"
	"sso/internal/services/auth/tests/unit/fake"
	usecase "sso/internal/services/auth/use-case"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegistrationUserUseCase", Label("unit"), func() {
	var log *slog.Logger
	var fakeUser UserWithPassword
	var users repository.UserRepository
	var fakeUsers *fake.FakeUserRepository

	var registrationUser usecase.RegistrationUserUseCase

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		usersRepo := fake.NewFakeUserRepository()
		users = usersRepo
		fakeUsers = usersRepo

		fakeUser = ValidUser

		registrationUser = usecase.RegistrationUserUseCase{Users: users}
	})

	It("should register user", func() {
		// Arrange

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).To(BeNil())
		Expect(user.ID).ToNot(BeZero())
		Expect(user.Email).To(Equal(fakeUser.Email))
	})

	It("should error on failed get user by email", func() {
		// Arrange
		fakeUser.Email = "need_error_for_1_times@local.com"

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error if user already exist", func() {
		// Arrange
		users.Save(context.Background(), fakeUser.Email, fakeUser.PassHash)

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).To(MatchError(models.ErrUserAlreadyExists))
		Expect(user).To(BeNil())
	})

	It("should error on failed hash password", func() {
		// Arrange
		fakeUser.Password = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error on failed save user", func() {
		// Arrange
		fakeUser.Email = "need_error_on_save@local.com"

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error on failed get user by email after save", func() {
		// Arrange
		fakeUser.Email = "need_error_for_2_times@local.com"

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).ToNot(BeNil())
		Expect(user).To(BeNil())
	})

	It("should error on user not saved", func() {
		// Arrange
		fakeUsers.SetCounter(int64(-1))

		// Action
		user, err := registrationUser.Execute(context.Background(), log, fakeUser.Email, fakeUser.Password)

		// Assert
		Expect(err).To(MatchError(models.ErrUserNotFound))
		Expect(user).To(BeNil())
	})
})

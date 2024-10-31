package fake_repository

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/models"
)

type FakeUserRepository struct {
	users []user
	log   *slog.Logger
}

type user struct {
	ID       int64
	Email    string
	Password []byte
}

func NewFakeUserRepository(log *slog.Logger) *FakeUserRepository {
	return &FakeUserRepository{log: log}
}

var counter int64

func (r FakeUserRepository) Save(
	ctx context.Context,
	email string,
	passHash []byte,
) (err error) {
	fmt.Println("counter1: ", counter)
	counter++
	r.users = append(r.users, user{ID: counter, Email: email, Password: passHash})
	fmt.Println(r.users)
	fmt.Println("len(r.users)1: ", len(r.users))
	return nil
}

func (r FakeUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (models.UserModel, error) {
	fmt.Println("counter2: ", counter)
	var user models.UserModel

	fmt.Println("len(r.users)2: ", len(r.users))
	for _, u := range r.users {
		fmt.Println(u.Email, "==", email)
		if u.Email == email {
			user = models.UserModel{ID: u.ID, Email: u.Email, PassHash: u.Password}
			break
		}
	}

	return user, nil
}

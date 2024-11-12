package fake

import (
	"context"
	"errors"
	auth "sso/internal/services/auth/model"
)

type FakeUserRepository struct {
	users   []user
	counter int64
}

type user struct {
	ID       int64
	Email    string
	Password []byte
}

func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{}
}

func (r *FakeUserRepository) Save(
	ctx context.Context,
	email string,
	passHash []byte,
) (err error) {
	// For test only
	if email == "need_error_on_save@local.com" {
		return errors.New("error for test")
	}

	r.counter++
	r.users = append(r.users, user{ID: r.counter, Email: email, Password: passHash})
	return nil
}

func (r *FakeUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*auth.UserModel, error) {
	var user auth.UserModel

	// For test only
	if email == "need_error@local.com" {
		return &user, errors.New("error for test")
	}

	// For test only
	if email == "need_null@local.com" {
		return &user, nil
	}

	for _, u := range r.users {
		if u.Email == email {
			user = auth.UserModel{ID: u.ID, Email: u.Email, PassHash: u.Password}
			break
		}
	}

	return &user, nil
}

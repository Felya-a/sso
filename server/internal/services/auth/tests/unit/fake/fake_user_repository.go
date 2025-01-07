package fake

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	auth "sso/internal/services/auth/model"
	"strconv"
	"strings"
)

type FakeUserRepository struct {
	Users        []user
	counter      int64
	errorCounter map[int]int
}

type user struct {
	ID       int64
	Email    string
	Password []byte
}

func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{errorCounter: make(map[int]int)}
}

func (r *FakeUserRepository) GetById(
	ctx context.Context,
	id int64,
) (*auth.UserModel, error) {
	var user auth.UserModel

	// For test only
	if id == 500 {
		return nil, errors.New("error from fake UserRepository")
	}

	for _, u := range r.Users {
		if u.ID == id {
			user = auth.UserModel{ID: u.ID, Email: u.Email, PassHash: u.Password}
			break
		}
	}

	return &user, nil
}

func (r *FakeUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*auth.UserModel, error) {
	var user auth.UserModel

	// For test only
	// Давать ошибку только на конкретную попытку
	if strings.HasPrefix(email, "need_error_for_") {
		re := regexp.MustCompile(`need_error_for_(\d+)_times@local\.com`)
		// Находим все совпадения
		match := re.FindStringSubmatch(email)
		// Сохраняет результат
		count, _ := strconv.Atoi(match[1])
		r.errorCounter[count] = r.errorCounter[count] + 1

		if count == r.errorCounter[count] {
			return nil, errors.New("error from fake UserRepository")
		}
	}

	// For test only
	if email == "need_null@local.com" {
		return &user, nil
	}

	for _, u := range r.Users {
		if u.Email == email {
			user = auth.UserModel{ID: u.ID, Email: u.Email, PassHash: u.Password}
			break
		}
	}

	return &user, nil
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
	r.Users = append(r.Users, user{ID: r.counter, Email: email, Password: passHash})
	return nil
}

// Fake only
func (r *FakeUserRepository) Delete(
	id int64,
) {
	// Поиск индекса элемента
	var index int = -1

	for i, v := range r.Users {
		if v.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("[FakeUserRepository] пользователь не найден")
		return
	}

	r.Users = append(r.Users[:index], r.Users[index+1:]...)
}

func (r *FakeUserRepository) SetCounter(
	counter int64,
) {
	r.counter = counter
}

package fake

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type FakeAuthorizationCodeRepository struct {
	Codes            []string
	ttlCode          time.Duration
	needErrorOnSave  bool
	needErrorOnCheck bool
}

func NewFakeAuthorizationCodeRepository(ttlCode time.Duration) *FakeAuthorizationCodeRepository {
	return &FakeAuthorizationCodeRepository{ttlCode: ttlCode}
}

func (r *FakeAuthorizationCodeRepository) CheckEndDelete(
	ctx context.Context,
	code string,
) (bool, error) {
	if r.needErrorOnCheck {
		return false, errors.New("fake error")
	}

	var index int = -1
	for i, v := range r.Codes {
		if v == code {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("[FakeAuthorizationCodeRepository] код не найден")
		return false, nil
	}

	r.Codes = append(r.Codes[:index], r.Codes[index+1:]...)
	return true, nil
}

func (r *FakeAuthorizationCodeRepository) Save(
	ctx context.Context,
	code string,
) error {
	if r.needErrorOnSave {
		return errors.New("fake error")
	}

	r.Codes = append(r.Codes, code)

	return nil
}

/* FOR TEST ONLY */

func (r *FakeAuthorizationCodeRepository) SetNeedErrorOnSave() {
	r.needErrorOnSave = true
}

func (r *FakeAuthorizationCodeRepository) SetNeedErrorOnCheck() {
	r.needErrorOnCheck = true
}

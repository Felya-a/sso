package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	models "sso/internal/services/auth/model"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r PostgresUserRepository) GetById(
	ctx context.Context,
	id int64,
) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `
		select
			id,
			email,
			password
		from public.user
		where id = $1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.UserModel{}, nil
		}
		return &models.UserModel{}, err
	}

	return &user, nil
}

func (r PostgresUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `
		select
			id,
			email,
			password
		from public.user
		where email = $1
	`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.UserModel{}, nil
		}
		return &models.UserModel{}, err
	}

	return &user, nil
}

func (r PostgresUserRepository) Save(
	ctx context.Context,
	email string,
	passHash []byte,
) (err error) {
	_, err = r.db.Exec(`
		insert into public.user (
			email,
			password
		) values (
			$1,
			$2
		)
	`, email, passHash)

	if err != nil {
		// TODO: refactor log
		fmt.Println(err.Error())
		return err
	}

	return nil
}

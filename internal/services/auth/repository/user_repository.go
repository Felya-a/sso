package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	authModels "sso/internal/services/auth/model"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPostgresUserRepository(db *sqlx.DB, log *slog.Logger) *PostgresUserRepository {
	return &PostgresUserRepository{db: db, log: log}
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
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (r PostgresUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*authModels.UserModel, error) {
	var user authModels.UserModel
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
			return &authModels.UserModel{}, nil
		}
		return &authModels.UserModel{}, err
	}

	return &user, nil
}

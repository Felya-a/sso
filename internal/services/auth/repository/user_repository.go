package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sso/internal/domain/models"

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
) (uid int64, err error) {
	res := r.db.MustExec(`
		insert into public.user (
			email,
			password
		) values (
			$1,
			$2
		)
	`, email, passHash)
	id, _ := res.LastInsertId()
	return id, nil
}

func (r PostgresUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	var user models.User
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
			return models.User{}, nil
		}
		r.log.Debug(err.Error())
		return models.User{}, err
	}

	return user, nil
}

func (r PostgresUserRepository) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	panic("not implemented")
}

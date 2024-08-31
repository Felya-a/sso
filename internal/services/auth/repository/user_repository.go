package repository

import (
	"context"
	"sso/internal/domain/models"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r PostgresUserRepository) Save(
	ctx context.Context,
	email string,
	passHash []byte,
) (uid int64, err error) {
	panic("not implemented")
}

func (r PostgresUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	panic("not implemented")
}

func (r PostgresUserRepository) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	panic("not implemented")
}

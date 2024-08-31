package repository

import (
	"context"
	"sso/internal/domain/models"

	"github.com/jmoiron/sqlx"
)

type PostgresAppRepository struct {
	db *sqlx.DB
}

func NewPostgresAppRepository(db *sqlx.DB) *PostgresAppRepository {
	return &PostgresAppRepository{db: db}
}

func (r PostgresAppRepository) App(ctx context.Context, appID int) (models.App, error) {
	panic("not implemented")
}

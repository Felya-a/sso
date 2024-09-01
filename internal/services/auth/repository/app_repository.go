package repository

import (
	"context"
	"log/slog"
	"sso/internal/domain/models"

	"github.com/jmoiron/sqlx"
)

type PostgresAppRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPostgresAppRepository(db *sqlx.DB, log *slog.Logger) *PostgresAppRepository {
	return &PostgresAppRepository{db: db, log: log}
}

func (r PostgresAppRepository) App(ctx context.Context, appID int) (models.App, error) {
	return models.App{ID: 1}, nil
}

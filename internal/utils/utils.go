package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sso/internal/config"

	_ "database/sql"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func MustConnectPostgres(config config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", config.GetPostgresConnectionString())
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func Migrate(db *sqlx.DB) {
	fullMigrationsPath := path.Join(GetWdPath(), "./migrations")

	// Применение миграций
	if err := goose.Up(db.DB, fullMigrationsPath); err != nil {
		if errors.Is(err, goose.ErrNoMigrations) {
			fmt.Println("Все миграции выполнены")
		}
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
}

func GetWdPath() string {
	wdFromEnv := os.Getenv("WORKDIR_PATH")
	wdFromOs, _ := os.Getwd()

	if wdFromEnv != "" {
		return wdFromEnv
	}

	return wdFromOs
}

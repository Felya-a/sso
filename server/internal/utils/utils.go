package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sso/internal/config"
	"strconv"

	_ "database/sql"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func MustConnectPostgres(config config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", GetPostgresConnectionString(config))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func MustConnectRedis(config config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		// Addr: "localhost:6379", // Укажите адрес вашего Redis
		Addr: config.Redis.Address,
	})
	return rdb
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

func GetPostgresUrl(config config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=disable",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
	)
}

func GetPostgresConnectionString(config config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Postgres.Host,
		strconv.Itoa(config.Postgres.Port),
		config.Postgres.User,
		config.Postgres.Database,
		config.Postgres.Password,
		"disable",
	)
}

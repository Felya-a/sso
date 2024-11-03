package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config структура, содержащая всю конфигурацию
type Config struct {
	Env       string         `yaml:"env" env-required:"true"`
	TokenTtl  time.Duration  `yaml:"token_ttl"`
	JWTSecret string         `yaml:"jwt_secret"`
	Grpc      GrpcConfig     `yaml:"grpc"`
	Postgres  PostgresConfig `yaml:"postgres"`
}

// PostgresConfig структура, содержащая настройки для подключения к Postgresql
type PostgresConfig struct {
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

// GrpcConfig структура, содержащая настройки для gRPC
type GrpcConfig struct {
	Port    string `yaml:"port"`
	Timeout string `yaml:"timeout"`
}

var config Config

// Get возвращает копию текущей конфигурации
func Get() Config {
	return config
}

func (c *Config) GetPostgresUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=disable",
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database,
	)
}

func (c *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Postgres.Host,
		strconv.Itoa(c.Postgres.Port),
		c.Postgres.User,
		c.Postgres.Database,
		c.Postgres.Password,
		"disable",
	)
}

// MustLoad загружает конфигурацию из файла и возвращает её
func MustLoad() Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty. you need to specify --config=<file_path>")
	}

	fullConfigPath := getAbsoluteConfigPath(configPath)

	if _, err := os.Stat(fullConfigPath); os.IsNotExist(err) {
		panic("config file does not exist: " + fullConfigPath)
	}

	if err := cleanenv.ReadConfig(fullConfigPath, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return config
}

// fetchConfigPath возвращает путь к конфигурационному файлу из аргументов командной строки или переменной окружения
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func getAbsoluteConfigPath(configPath string) string {
	wdFromEnv := os.Getenv("WORKDIR_PATH")
	wdFromOs, _ := os.Getwd()

	if wdFromEnv != "" {
		return path.Join(wdFromEnv, configPath)
	}
	return path.Join(wdFromOs, configPath)
}

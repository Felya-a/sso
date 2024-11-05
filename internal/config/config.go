package config

import (
	"flag"
	"os"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config структура, содержащая всю конфигурацию
type Config struct {
	Env       string        `env:"ENV" env-required:"true"`
	TokenTtl  time.Duration `env:"TOKEN_TTL"`
	JWTSecret string        `env:"JWT_SECRET"`
	Grpc      GrpcConfig
	Postgres  PostgresConfig
}

// PostgresConfig структура, содержащая настройки для подключения к Postgresql
type PostgresConfig struct {
	User     string `env:"POSTGRES_USER"`
	Database string `env:"POSTGRES_DATABASE"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
}

// GrpcConfig структура, содержащая настройки для gRPC
type GrpcConfig struct {
	Port    string `env:"GRPC_PORT"`
	Timeout string `env:"GRPC_TIMEOUT"`
}

var config Config

// Get возвращает копию текущей конфигурации
func Get() Config {
	return config
}

// MustLoad загружает конфигурацию из файла и возвращает её
func MustLoad() Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty. you need to specify --config=<file_path> or environment CONFIG_PATH")
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

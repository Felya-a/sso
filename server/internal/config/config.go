package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config структура, содержащая всю конфигурацию
type Config struct {
	Env      string `env:"ENV" env-required:"true"`
	Jwt      JwtConfig
	Grpc     GrpcConfig
	Http     HttpConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type JwtConfig struct {
	Secret     string        `env:"JWT_SECRET" env-required:"true"`
	AccessTtl  time.Duration `env:"JWT_ACCESS_TTL" env-required:"true"`
	RefreshTtl time.Duration `env:"JWT_REFRESH_TTL" env-required:"true"`
}

// PostgresConfig структура, содержащая настройки для подключения к Postgresql
type PostgresConfig struct {
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Database string `env:"POSTGRES_DATABASE" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
}

// RedisConfig структура, содержащая настройки для подключения к Redis
type RedisConfig struct {
	Address string `env:"REDIS_ADDRESS" env-required:"true"`
}

// GrpcConfig структура, содержащая настройки для gRPC
type GrpcConfig struct {
	Host    string `env:"GRPC_HOST" env-required:"true" env-description:"gRPC server host for tests"`
	Port    string `env:"GRPC_PORT" env-required:"true"`
	Timeout string `env:"GRPC_TIMEOUT"`
}

// HttpConfig структура, содержащая настройки для http
type HttpConfig struct {
	Host string `env:"HTTP_HOST" env-required:"true" env-description:"http server host for tests"`
	Port string `env:"HTTP_PORT" env-required:"true"`
}

var config Config

// Get возвращает копию текущей конфигурации
func Get() Config {
	return config
}

// MustLoad загружает конфигурацию из файла и возвращает её
func MustLoad() Config {
	// Чтение переменных из окружения
	err := cleanenv.ReadEnv(&config)
	if err == nil {
		return config
	}
	fmt.Println("error on read raw env: " + err.Error())

	// Чтение переменных из конфигурационного файла
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty. you need to specify --config=<file_path> or environment CONFIG_PATH")
	}

	fullConfigPath := getAbsoluteConfigPath(configPath)

	fmt.Println(fullConfigPath)

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

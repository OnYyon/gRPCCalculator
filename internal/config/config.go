package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

type AppConfig struct {
	Name    string `yaml:"name" env:"APP_NAME" env-default:"gRPCCalculator"`
	Version string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
}

type ServerConfig struct {
	Host                 string `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port                 string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	ComputingPower       int    `yaml:"computing_power" env:"COMPUTING_POWER" env-default:"13"`
	TimeAdditionMS       int64  `yaml:"TIME_ADDITION_MS" env:"TIME_ADDITION_MS" env-default:"100"`
	TimeSubtractionMS    int64  `yaml:"TIME_SUBTRACTION_MS" env:"TIME_SUBTRACTION_MS" env-default:"100"`
	TimeMultiplicationMS int64  `yaml:"TIME_MULTIPLICATIONS_MS" env:"TIME_MULTIPLICATIONS_MS" env-default:"300"`
	TimeDivisionMS       int64  `yaml:"TIME_DIVISIONS_MS" env:"TIME_DIVISIONS_MS" env-default:"400"`
}

type DatabaseConfig struct {
	DBPath         string `yaml:"dbPath" env:"DB_PATH" env-default:"./storage/calculator.db"`
	MigrationsPath string `yaml:"migrationsPath" env:"MIGRATIONS_PATH" env-default:"./migrations"`
}

type AuthConfig struct {
	JWTSecret    string        `yaml:"jwt_secret" env:"JWT_SECRET" env-default:"yandex-go"`
	TokenExpires time.Duration `yaml:"token_expires" env:"TOKEN_EXPIRES" env-default:"24h"`
}

func Load(configPath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Name    string `yaml:"name" env:"APP_NAME" env-default:"gRPCCalculator"`
		Version string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	} `yaml:"app"`

	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
		Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	} `yaml:"server"`

	Database struct {
		DBPath         string `yaml:"dbPath" env:"DB_PATH" env-default:"./storage/calculator.db"`
		MigrationsPath string `yaml:"migrationsPath" env:"MIGRATIONS_PATH" env-default:"./migrations"`
	} `yaml:"database"`

	Auth struct {
		JWTSecret    string        `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
		TokenExpires time.Duration `yaml:"token_expires" env:"TOKEN_EXPIRES" env-default:"24h"`
	} `yaml:"auth"`
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

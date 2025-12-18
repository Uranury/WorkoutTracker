package config

import "fmt"

type Config struct {
	DBConfig
	RedisAddr      string `yaml:"redis_addr" env:"REDIS_ADDR" env-default:":6379"`
	MigrationsPath string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
	ListenAddr     string `yaml:"listen_addr" env:"LISTEN_ADDR" env-default:":8080"`
}

type DBConfig struct {
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
}

func (cfg DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}

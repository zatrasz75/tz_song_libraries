package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Server struct {
		AddrHost           string        `yaml:"host" env:"APP_IP" env-default:"0.0.0.0"`
		AddrPort           string        `yaml:"port" env:"APP_PORT" env-default:"8586"`
		ReadTimeout        time.Duration `yaml:"read-timeout" env:"READ_TIMEOUT" env-default:"3s"`
		WriteTimeout       time.Duration `yaml:"write-timeout" env:"WRITE_TIMEOUT" env-default:"3s"`
		IdleTimeout        time.Duration `yaml:"idle-timeout" env:"IDLE_TIMEOUT" env-default:"6s"`
		ShutdownTime       time.Duration `yaml:"shutdown-timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
		CORSAllowedOrigins []string      `yaml:"cors-allowed-origins" env:"CORS_ALLOWED_ORIGINS" env-default:"localhost"`
	} `yaml:"server"`
	DataBase struct {
		ConnStr string `env:"DB_CONNECTION_STRING" env-description:"db string"`

		Host     string `yaml:"host" env:"HOST_DB"`
		User     string `yaml:"username" env:"POSTGRES_USER"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
		Url      string `yaml:"db-url" env:"URL_DB"`
		Name     string `yaml:"db-name" env:"POSTGRES_DB"`
		Port     string `yaml:"port" env:"PORT_DB"`

		PoolMax      int           `yaml:"pool-max" env:"PG_POOL_MAX" env-default:"2"`
		ConnAttempts int           `yaml:"conn-attempts" env:"PG_CONN_ATTEMPTS" env-default:"5"`
		ConnTimeout  time.Duration `yaml:"conn-timeout" env:"PG_TIMEOUT" env-default:"2s"`
	} `yaml:"database"`
	External struct {
		Url string `yaml:"url" env:"URL_EXTERNAL"`
	} `yaml:"external"`
}

func NewConfig(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	cfg.DataBase.ConnStr = initDB(cfg)

	return &cfg, nil
}

func initDB(cfg Config) string {
	if cfg.DataBase.ConnStr != "" {
		return cfg.DataBase.ConnStr
	}
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DataBase.Host,
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.Url,
		cfg.DataBase.Port,
		cfg.DataBase.Name,
	)
}

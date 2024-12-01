package config

import (
	"fmt"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort      int    `env:"HTTP_PORT" env-default:"8080"`
	SongDetailAPI string `env:"SONG_DETAIL_API" env-required:"true"`
	Mode          string `env:"MODE" env-default:"development"`

	DB DBConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

func (c *DBConfig) ToDSN() string {
	q := url.Values{}
	q.Add("sslmode", "disable")

	u := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Path:     c.Name,
		RawQuery: q.Encode(),
	}

	return u.String()
}

func ParseConfig() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

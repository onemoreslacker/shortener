package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"

	"github.com/caarlos0/env"
	"github.com/spf13/viper"
)

type (
	Database struct {
		Username string `env:"MONGO_INITDB_ROOT_USERNAME" yaml:"username" envDefault:"username"`
		Password string `env:"MONGO_INITDB_ROOT_PASSWORD" yaml:"password" envDefault:"password"`
		Host     string `env:"MONGO_HOST" yaml:"host" envDefault:"mongo"`
		Port     string `env:"MONGO_PORT" yaml:"port" envDefault:"27017"`
		Name     string `env:"MONGO_INITDB_DATABASE" yaml:"name" envDefault:"db"`
	}

	Serving struct {
		UIHTTP         string `env:"UI_HTTP" yaml:"uiHTTP" envDefault:"8080"`
		RedirectorHTTP string `env:"REDIRECTOR_HTTP" yaml:"redirectorHTTP" envDefault:"8081"`
		RedirectorGRPC string `env:"REDIRECTOR_GRPC" yaml:"redirectorGRPC" envDefault:"9091"`
		ShortenerHTTP  string `env:"SHORTENER_HTTP" yaml:"shortenerHTTP" envDefault:"8082"`
		ShortenerGRPC  string `env:"SHORTENER_GRPC" yaml:"shortenerGRPC" envDefault:"9092"`
	}

	Collections struct {
		Links string `env:"LINKS_COLLECTION" yaml:"links" envDefault:"links"`
	}

	ShorteningPolicy struct {
		Length  int    `env:"SHORTEN_LENGTH" yaml:"length" envDefault:"6"`
		BaseURL string `env:"BASE_URL" yaml:"baseURL" envDefault:"http://localhost"`
	}
)

type Config struct {
	Database         Database         `yaml:"database"`
	Serving          Serving          `yaml:"serving"`
	Collections      Collections      `yaml:"collections"`
	ShorteningPolicy ShorteningPolicy `yaml:"shorteningPolicy"`
}

func New(path string) (*Config, error) {
	cfg, err := NewConfigFromFile(path)
	if err != nil {
		return nil, err
	}

	if err = NewConfigFromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewConfigFromFile(name string) (*Config, error) {
	cfg := &Config{}

	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigFile(name)

	if err := v.ReadConfig(bytes.NewBuffer(configBytes)); err != nil {
		return nil, fmt.Errorf("config: failed to read: %w", err)
	}

	if err := v.MergeInConfig(); err != nil {
		if errors.Is(err, &viper.ConfigParseError{}) {
			return nil, fmt.Errorf("config: failed to merge: %w", err)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal: %w", err)
	}

	return cfg, nil
}

func NewConfigFromEnv(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("config: failed to parse env: %w", err)
	}

	return nil
}

func (d *Database) ToDSN() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

//go:embed config.yaml
var configBytes []byte

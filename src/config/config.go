package config

import (
	"errors"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
	DBConfig
	Environment
	JWT
	TenantID            uuid.UUID `env:"KOSVI_TENANT_ID,notEmpty,required"`
	MasterEncryptionKey string    `env:"MASTER_ENCRYPTION_KEY,notEmpty,required"`
}

func (c *Config) IsDev() bool {
	return c.Environment.Name == "development"
}

type JWT struct {
	TTL                int    `env:"JWT_TTL" envDefault:"3600"`
	ECPrivateKeyBase64 string `env:"JWT_EC_PRIVATE_KEY_BASE64,notEmpty,required"`
	ECPublicKeyBase64  string `env:"JWT_EC_PUBLIC_KEY_BASE64,notEmpty,required"`
}

type Environment struct {
	Name string `env:"ENV" envDefault:"development"`
}

type DBConfig struct {
	DSN     string `env:"KOSVI_DB_DSN,notEmpty,required"`
	DBHost  string `env:"DB_HOST,notEmpty,required"`
	DBPort  int64  `env:"DB_PORT,notEmpty,required"`
	SSLMode string `env:"SSL_MODE,notEmpty,required"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() []error {
	err := env.Parse(c)
	if err == nil {
		return nil
	}

	if errors.Is(err, env.EmptyVarError{}) {
		return []error{errors.New("missing environment variables")}
	}

	aggErr := env.AggregateError{}
	if ok := errors.As(err, &aggErr); ok {
		errs := []error{}
		for _, err := range aggErr.Errors {
			errs = append(errs, errors.New(err.Error()))
		}
		return errs
	}

	return nil
}

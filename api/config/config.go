package config

import (
	"bytes"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"reflect"
	"sync"
)

type (
	Config struct {
		App        App
		HTTP       HTTP
		Log        Log
		PostgreSQL PostgreSQL
	}

	// App - represent application configuration.
	App struct {
		BaseURL string `env:"BASE_URL"    env-default:"http://localhost:8082"`
	}

	// HTTP - represents http configuration.
	HTTP struct {
		Port string `env:"HTTP_PORT" env-default:"8082"`
	}

	// Log - represents logger configuration.
	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug"`
	}

	// PostgreSQL - represents PostgreSQL database configuration.
	PostgreSQL struct {
		User     string `env:"POSTGRESQL_USER"     env-default:"postgres"`
		Password string `env:"POSTGRESQL_PASSWORD" env-default:"postgres"`
		Host     string `env:"POSTGRESQL_HOST"     env-default:"127.0.0.1"`
		Database string `env:"POSTGRESQL_DATABASE" env-default:"api"`
	}

	// JWT - represents jwt configuration.
	JWT struct {
		SignKey string `env:"JWT_SIGN_KEY"     env-default:"sajkdjk1ndansdnan"`
	}
)

// Replace is used to replace values in static files with populated values from config.
// It replaces config variables (e.g. {{BaseURL}}) with real values from config
// and serves updated files.
func (config *Config) Replace(input []byte) []byte {
	// get config keys
	v := reflect.ValueOf(config.App)
	typeOfS := v.Type()

	// replace all config variables
	output := input
	for i := 0; i < v.NumField(); i++ {
		output = bytes.Replace(
			output,
			[]byte(fmt.Sprintf("{{%s}}", typeOfS.Field(i).Name)),
			[]byte(v.Field(i).String()),
			-1,
		)
	}

	return output
}

var (
	config Config
	once   sync.Once
)

// Get returns config
func Get() *Config {
	once.Do(func() {
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			log.Fatal("failed to read env", err)
		}
	})

	return &config
}

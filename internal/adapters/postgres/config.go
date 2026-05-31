// Конфигурация адаптера для PostgreSQL.
package postgres

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// PostgreSQL Connection Config.
type Config struct {
	User     string
	Password string
	Port     string
	Host     string
	DBName   string
}

func (c Config) getConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.DBName)
}

var (
	// ErrPostgresEnvVarNotFound.
	ErrPostgresEnvVarNotFound = errors.New("postgres env var are not found")
	// ErrPostgresEnvVarAreEmpty.
	ErrPostgresEnvVarAreEmpty = errors.New("postgres env var are empty")
)

// Create New Postgres config. Panic if error happened.
func MustNewConfig() *Config {
	var (
		notFound = make([]int, 0, 5)
		params   = make([]string, 0, 5)
	)
	pgUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		notFound = append(notFound, 0)
	}
	pgPass, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		notFound = append(notFound, 0)
	}
	pgPort, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		notFound = append(notFound, 0)
	}
	pgHost, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		notFound = append(notFound, 0)
	}
	pgDB, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		notFound = append(notFound, 0)
	}
	if slices.Contains(notFound, 0) {
		panic(ErrPostgresEnvVarNotFound)
	}
	params = append(params, pgUser, pgPass, pgPort, pgHost, pgDB)
	for _, v := range params {
		clean := strings.TrimSpace(v)
		if len(clean) == 0 {
			panic(ErrPostgresEnvVarAreEmpty)
		}
	}
	return &Config{
		User:     pgUser,
		Password: pgPass,
		Host:     pgHost,
		DBName:   pgDB,
	}
}

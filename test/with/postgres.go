//go:build integration
// +build integration

package with

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbConfig struct {
	user, password, db string
	port, exposed      nat.Port
}

func Getenv(key, fallback string) string {
	if env := os.Getenv(key); len(env) != 0 {
		return env
	}
	return fallback
}

func Postgres(t *testing.T, ctx context.Context) *dbConfig {
	t.Helper()

	config := dbConfig{
		user:     Getenv("DB_USER", "test"),
		password: Getenv("DB_PASSWORD", ""),
		db:       Getenv("DB_NAME", "test"),
		port:     nat.Port("5432"),
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{config.port.Port()},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		Env: map[string]string{
			"POSTGRES_USER":             config.user,
			"POSTGRES_PASSWORD":         config.password,
			"POSTGRES_DB":               config.db,
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
	}

	psql, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	config.exposed, err = psql.MappedPort(ctx, config.port)
	assert.NoError(t, err)

	t.Cleanup(func() { _ = psql.Terminate(ctx) })

	return &config
}

func (cfg *dbConfig) Connect(t *testing.T, ctx context.Context) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.user,
		cfg.password,
		cfg.db,
		cfg.exposed.Port(),
	)

	conn := func() (*gorm.DB, error) {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	t.Logf("config: %s", dsn)
	t.Log("connecting...")
	defer t.Log("connecting...done")

	db, err := conn()
	for err != nil {
		time.Sleep(time.Millisecond * 500)
		db, err = conn()
	}

	return db
}

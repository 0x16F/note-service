package tester

import (
	"context"
	"notes-manager/src/pkg/fsscanner"
	"notes-manager/src/pkg/migrate"
	"notes-manager/src/usecase/dcontainer"
	"notes-manager/src/usecase/repository/pgconnector"
	"notes-manager/src/usecase/repository/rsconnector"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func SetupPostgresql(t *testing.T) *gorm.DB {
	ctx := context.Background()

	postgresqlC, err := dcontainer.NewPostgres(ctx)
	require.NoError(t, err)

	// Cleanup function to terminate the container after the test
	t.Cleanup(func() {
		if err := postgresqlC.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	startDir, err := os.Getwd()
	require.NoError(t, err)

	migrationsPath, err := fsscanner.FindDirectory(startDir, "migrations")
	require.NoError(t, err)
	require.NotEqual(t, migrationsPath, "")

	host, err := postgresqlC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresqlC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := pgconnector.Config{
		Host:     host,
		User:     "postgres",
		Password: "note-service",
		Port:     uint16(port.Int()),
		DB:       "note_service",
	}

	db, err := pgconnector.Connect(&cfg)
	require.NoError(t, err)

	err = migrate.ApplyMigrations(&cfg, false, "file://"+migrationsPath)
	require.NoError(t, err)

	return db
}

func SetupRedis(t *testing.T) *redis.Client {
	ctx := context.Background()

	redisC, err := dcontainer.NewRedis(ctx)
	require.NoError(t, err)

	// Cleanup function to terminate the container after the test
	t.Cleanup(func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	startDir, err := os.Getwd()
	require.NoError(t, err)

	migrationsPath, err := fsscanner.FindDirectory(startDir, "migrations")
	require.NoError(t, err)
	require.NotEqual(t, migrationsPath, "")

	endpoint, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)

	rCfg := rsconnector.Config{
		Host:     endpoint,
		Password: "",
		DB:       0,
	}

	client, err := rsconnector.Connect(&rCfg)
	require.NoError(t, err)

	return client
}

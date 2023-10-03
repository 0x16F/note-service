package user

import (
	"context"
	"notes-manager/src/pkg/fsscanner"
	"notes-manager/src/pkg/migrate"
	"notes-manager/src/usecase/dcontainer"
	"notes-manager/src/usecase/repository/pgconnector"
	"notes-manager/src/usecase/repository/rsconnector"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	ctx := context.Background()

	postgresqlC, err := dcontainer.NewPostgres(ctx)
	require.NoError(t, err)

	redisC, err := dcontainer.NewRedis(ctx)
	require.NoError(t, err)

	defer func() {
		if err := postgresqlC.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	startDir, err := os.Getwd()
	require.NoError(t, err)

	migrationsPath, err := fsscanner.FindDirectory(startDir, "migrations")
	require.NoError(t, err)
	require.NotEqual(t, migrationsPath, "")

	host, err := postgresqlC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresqlC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	pCfg := pgconnector.Config{
		Host:     host,
		User:     "postgres",
		Password: "note-service",
		Port:     uint16(port.Int()),
		DB:       "note_service",
	}

	db, err := pgconnector.Connect(&pCfg)
	require.NoError(t, err)

	err = migrate.ApplyMigrations(&pCfg, false, "file://"+migrationsPath)
	require.NoError(t, err)

	defer migrate.ApplyMigrations(&pCfg, true, "file://"+migrationsPath)

	endpoint, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)

	rCfg := rsconnector.Config{
		Host:     endpoint,
		Password: "",
		DB:       0,
	}

	client, err := rsconnector.Connect(&rCfg)
	require.NoError(t, err)

	repo := NewRepo(db, client)

	u := New("login", "password")

	fetchedUser, err := repo.Fetch(ctx, u.Id)
	require.Error(t, err)
	require.Nil(t, fetchedUser)

	err = repo.Create(ctx, u)
	require.NoError(t, err)

	fetchedUser, err = repo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)

	fetchedUser, err = repo.FetchLogin(ctx, u.Login)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)

	u.Login = "new login"

	err = repo.Update(ctx, u)
	require.NoError(t, err)

	fetchedUser, err = repo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)

	err = repo.Delete(ctx, u.Id)
	require.NoError(t, err)

	fetchedUser, err = repo.Fetch(ctx, u.Id)
	require.Error(t, err)
	require.Nil(t, fetchedUser)
}

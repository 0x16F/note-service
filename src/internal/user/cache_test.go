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

var cRepo Repository

func setupC(t *testing.T) context.Context {
	ctx := context.Background()

	postgresqlC, err := dcontainer.NewPostgres(ctx)
	require.NoError(t, err)

	redisC, err := dcontainer.NewRedis(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := postgresqlC.Terminate(ctx); err != nil {
			t.Fatal(err)
		}

		if err := redisC.Terminate(ctx); err != nil {
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

	t.Cleanup(func() {
		migrate.ApplyMigrations(&pCfg, true, "file://"+migrationsPath)
	})

	endpoint, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)

	rCfg := rsconnector.Config{
		Host:     endpoint,
		Password: "",
		DB:       0,
	}

	client, err := rsconnector.Connect(&rCfg)
	require.NoError(t, err)

	cRepo = NewRepo(db, client)

	return ctx
}

func TestCreateUserWithRepo(t *testing.T) {
	ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = cRepo.Create(ctx, u)
	require.NoError(t, err)
}

func TestFetchUserWithRepo(t *testing.T) {
	ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = cRepo.Create(ctx, u)
	require.NoError(t, err)

	fetchedUser, err := cRepo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestFetchUserByLoginWithRepo(t *testing.T) {
	ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	// Ensure user is created first
	err = cRepo.Create(ctx, u)
	require.NoError(t, err)

	// Fetch the user by login and compare
	fetchedUser, err := cRepo.FetchLogin(ctx, u.Login)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestUpdateUserWithRepo(t *testing.T) {
	ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = cRepo.Create(ctx, u)
	require.NoError(t, err)

	// Update user login
	u.Login = "new login"
	err = cRepo.Update(ctx, u)
	require.NoError(t, err)

	// Fetch updated user and compare
	fetchedUser, err := cRepo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestDeleteUserWithRepo(t *testing.T) {
	ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = cRepo.Create(ctx, u)
	require.NoError(t, err)

	// Delete the user
	err = cRepo.Delete(ctx, u.Id)
	require.NoError(t, err)

	// Ensure user is deleted
	fetchedUser, err := cRepo.Fetch(ctx, u.Id)
	require.Error(t, err)
	require.Nil(t, fetchedUser)
}

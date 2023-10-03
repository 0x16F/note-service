package user

import (
	"context"
	"notes-manager/src/pkg/fsscanner"
	"notes-manager/src/pkg/migrate"
	"notes-manager/src/usecase/dcontainer"
	"notes-manager/src/usecase/repository/pgconnector"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// global variables for shared resources, such as the database connection and repository
var rDb *gorm.DB
var rRepo Repository

func setupR(t *testing.T) context.Context {
	ctx := context.Background()

	postgresqlC, err := dcontainer.NewPostgres(ctx)
	require.NoError(t, err)

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

	rDb, err = pgconnector.Connect(&cfg)
	require.NoError(t, err)

	err = migrate.ApplyMigrations(&cfg, false, "file://"+migrationsPath)
	require.NoError(t, err)

	t.Cleanup(func() {
		migrate.ApplyMigrations(&cfg, true, "file://"+migrationsPath)
	})

	rRepo = NewDatabaseRepo(rDb)

	return ctx
}

func TestCreateUser(t *testing.T) {
	ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = rRepo.Create(ctx, u)
	require.NoError(t, err)
}

func TestFetchUser(t *testing.T) {
	ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	// Ensure user is created first
	err = rRepo.Create(ctx, u)
	require.NoError(t, err)

	// Fetch the user and compare
	fetchedUser, err := rRepo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestFetchUserByLogin(t *testing.T) {
	ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	// Ensure user is created first
	err = rRepo.Create(ctx, u)
	require.NoError(t, err)

	// Fetch the user by login and compare
	fetchedUser, err := rRepo.FetchLogin(ctx, u.Login)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestUpdateUser(t *testing.T) {
	ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = rRepo.Create(ctx, u)
	require.NoError(t, err)

	u.Login = "new login"
	err = rRepo.Update(ctx, u)
	require.NoError(t, err)

	fetchedUser, err := rRepo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestDeleteUser(t *testing.T) {
	ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = rRepo.Create(ctx, u)
	require.NoError(t, err)

	err = rRepo.Delete(ctx, u.Id)
	require.NoError(t, err)

	fetchedUser, err := rRepo.Fetch(ctx, u.Id)
	require.Error(t, err)
	require.Nil(t, fetchedUser)
}

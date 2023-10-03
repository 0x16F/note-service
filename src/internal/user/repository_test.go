package user

import (
	"context"
	"notes-manager/src/usecase/tester"
	"testing"

	"github.com/stretchr/testify/require"
)

// setupR initializes the necessary components for an integration test.
// This function creates a new PostgreSQL container, applies migrations, and connects to the services.
func setupR(t *testing.T) (Repository, context.Context) {
	db := tester.SetupPostgresql(t)
	ctx := context.Background()

	return NewDatabaseRepo(db), ctx
}

func TestCreateUser(t *testing.T) {
	repo, ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, u)
	require.NoError(t, err)
}

func TestFetchUser(t *testing.T) {
	repo, ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	// Ensure user is created first
	err = repo.Create(ctx, u)
	require.NoError(t, err)

	// Fetch the user and compare
	fetchedUser, err := repo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestFetchUserByLogin(t *testing.T) {
	repo, ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	// Ensure user is created first
	err = repo.Create(ctx, u)
	require.NoError(t, err)

	// Fetch the user by login and compare
	fetchedUser, err := repo.FetchLogin(ctx, u.Login)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestUpdateUser(t *testing.T) {
	repo, ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, u)
	require.NoError(t, err)

	u.Login = "new login"
	err = repo.Update(ctx, u)
	require.NoError(t, err)

	fetchedUser, err := repo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestDeleteUser(t *testing.T) {
	repo, ctx := setupR(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, u)
	require.NoError(t, err)

	err = repo.Delete(ctx, u.Id)
	require.NoError(t, err)

	fetchedUser, err := repo.Fetch(ctx, u.Id)
	require.Error(t, err)
	require.Nil(t, fetchedUser)
}

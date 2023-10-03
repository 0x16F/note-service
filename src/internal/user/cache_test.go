package user

import (
	"context"
	"notes-manager/src/usecase/tester"
	"testing"

	"github.com/stretchr/testify/require"
)

// setupC initializes the necessary components for an integration test.
// This function creates a new PostgreSQL and Redis container, applies migrations, and connects to the services.
func setupC(t *testing.T) (Repository, context.Context) {
	ctx := context.Background()

	db := tester.SetupPostgresql(t)
	client := tester.SetupRedis(t)

	return NewRepo(db, client), ctx
}

func TestCreateUserWithRepo(t *testing.T) {
	repo, ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, u)
	require.NoError(t, err)
}

func TestFetchUserWithRepo(t *testing.T) {
	repo, ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, u)
	require.NoError(t, err)

	fetchedUser, err := repo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestFetchUserByLoginWithRepo(t *testing.T) {
	repo, ctx := setupC(t)

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

func TestUpdateUserWithRepo(t *testing.T) {
	repo, ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, u)
	require.NoError(t, err)

	// Update user login
	u.Login = "new login"
	err = repo.Update(ctx, u)
	require.NoError(t, err)

	// Fetch updated user and compare
	fetchedUser, err := repo.Fetch(ctx, u.Id)
	require.NoError(t, err)
	require.Equal(t, u, fetchedUser)
}

func TestDeleteUserWithRepo(t *testing.T) {
	repo, ctx := setupC(t)

	u, err := New("login", "password")
	require.NoError(t, err)

	err = repo.Create(ctx, u)
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(ctx, u.Id)
	require.NoError(t, err)

	// Ensure user is deleted
	fetchedUser, err := repo.Fetch(ctx, u.Id)
	require.Error(t, err)
	require.Nil(t, fetchedUser)
}

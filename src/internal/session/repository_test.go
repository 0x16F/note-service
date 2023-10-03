package session

import (
	"context"
	"notes-manager/src/usecase/tester"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupC(t *testing.T) (Repository, context.Context) {
	client := tester.SetupRedis(t)
	ctx := context.Background()

	return NewRepo(client), ctx
}

func TestFetchNonexistentSession(t *testing.T) {
	repo, ctx := setupC(t)

	s := New(uuid.New(), "user")
	fetchedSession, err := repo.Fetch(ctx, s.Id)

	require.NotNil(t, err)
	require.Nil(t, fetchedSession)
}

func TestCreateSession(t *testing.T) {
	repo, ctx := setupC(t)

	s := New(uuid.New(), "user")
	err := repo.Create(ctx, s)

	require.NoError(t, err)

	fetchedSession, err := repo.Fetch(ctx, s.Id)
	require.NoError(t, err)
	require.Equal(t, s, fetchedSession)
}

func TestFetchAllSessions(t *testing.T) {
	repo, ctx := setupC(t)

	s := New(uuid.New(), "user")
	err := repo.Create(ctx, s)
	require.NoError(t, err)

	fetchedSessions, err := repo.FetchAll(ctx, s.UserId)
	require.NoError(t, err)
	require.Equal(t, []*Session{s}, fetchedSessions)
}

func TestUpdateSessionActivity(t *testing.T) {
	repo, ctx := setupC(t)

	s := New(uuid.New(), "user")
	err := repo.Create(ctx, s)
	require.NoError(t, err)

	s.UpdateActivity()
	err = repo.Update(ctx, s)
	require.NoError(t, err)

	fetchedSession, err := repo.Fetch(ctx, s.Id)
	require.NoError(t, err)
	require.Equal(t, s, fetchedSession)
}

func TestDeleteSession(t *testing.T) {
	repo, ctx := setupC(t)

	s := New(uuid.New(), "user")
	err := repo.Create(ctx, s)
	require.NoError(t, err)

	err = repo.Delete(ctx, s.Id)
	require.NoError(t, err)

	fetchedSession, err := repo.Fetch(ctx, s.Id)
	require.NotNil(t, err)
	require.Nil(t, fetchedSession)
}

func TestFetchAllAfterDelete(t *testing.T) {
	repo, ctx := setupC(t)

	s := New(uuid.New(), "user")
	err := repo.Create(ctx, s)
	require.NoError(t, err)

	err = repo.Delete(ctx, s.Id)
	require.NoError(t, err)

	fetchedSessions, err := repo.FetchAll(ctx, s.UserId)
	require.NoError(t, err)
	require.Equal(t, 0, len(fetchedSessions))
}

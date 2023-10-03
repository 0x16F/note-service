package session

import (
	"context"
	"notes-manager/src/usecase/dcontainer"
	"notes-manager/src/usecase/repository/rsconnector"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	ctx := context.Background()

	redisC, err := dcontainer.NewRedis(ctx)
	require.NoError(t, err)

	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	endpoint, err := redisC.Endpoint(ctx, "")
	require.NoError(t, err)

	cfg := &rsconnector.Config{
		Host:     endpoint,
		Password: "",
		DB:       0,
	}

	client, err := rsconnector.Connect(cfg)
	require.NoError(t, err)

	repo := NewRepo(client)

	s := New(uuid.New(), "user")

	fetchedSession, err := repo.Fetch(ctx, s.Id)
	require.NotNil(t, err)
	require.Nil(t, fetchedSession)

	fetchedSessions, err := repo.FetchAll(ctx, s.UserId)
	require.NoError(t, err)
	require.Equal(t, 0, len(fetchedSessions))

	err = repo.Create(ctx, s)
	require.NoError(t, err)

	fetchedSession, err = repo.Fetch(ctx, s.Id)
	require.NoError(t, err)
	require.Equal(t, s, fetchedSession)

	fetchedSessions, err = repo.FetchAll(ctx, s.UserId)
	require.NoError(t, err)
	require.Equal(t, []*Session{s}, fetchedSessions)

	s.UpdateActivity()

	err = repo.Update(ctx, s)
	require.NoError(t, err)

	fetchedSession, err = repo.Fetch(ctx, s.Id)
	require.NoError(t, err)
	require.Equal(t, s, fetchedSession)

	err = repo.Delete(ctx, s.Id)
	require.NoError(t, err)

	fetchedSession, err = repo.Fetch(ctx, s.Id)
	require.NotNil(t, err)
	require.Nil(t, fetchedSession)

	fetchedSessions, err = repo.FetchAll(ctx, s.UserId)
	require.NoError(t, err)
	require.Equal(t, 0, len(fetchedSessions))
}

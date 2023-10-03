package note

import (
	"context"
	"notes-manager/src/internal/user"
	"notes-manager/src/usecase/tester"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupC(t *testing.T) (Repository, *Note, context.Context) {
	ctx := context.Background()

	db := tester.SetupPostgresql(t)
	client := tester.SetupRedis(t)

	repo := NewRepo(db, client)
	uRepo := user.NewRepo(db, client)

	u, err := user.New("login", "password")
	require.NoError(t, err)

	require.NoError(t, uRepo.Create(ctx, u))

	return repo, New(u.Id, "my note", "my content"), ctx
}

func TestCreateFetch(t *testing.T) {
	repo, note, ctx := setupC(t)

	// Test fetching a non-existing note
	fetchedNote, err := repo.Fetch(ctx, note.Id)
	require.Nil(t, fetchedNote)
	require.NotNil(t, err)

	// Test creating a new note
	require.NoError(t, repo.Create(ctx, note))

	// Test fetching the created note
	fetchedNote, err = repo.Fetch(ctx, note.Id)
	require.NoError(t, err)
	require.Equal(t, note, fetchedNote)
}

func TestUpdate(t *testing.T) {
	repo, note, ctx := setupC(t)

	require.NoError(t, repo.Create(ctx, note))

	// Test updating the note
	dto := NoteDTO{
		Id:        note.Id,
		Title:     "new title",
		Content:   "new content",
		UpdatedAt: time.Now().UTC(),
	}

	require.NoError(t, repo.Update(ctx, &dto))

	// Test fetching the updated note
	fetchedNote, err := repo.Fetch(ctx, note.Id)
	require.NoError(t, err)
	require.Equal(t, dto.Content, fetchedNote.Content)
}

func TestDelete(t *testing.T) {
	repo, note, ctx := setupC(t)

	require.NoError(t, repo.Create(ctx, note))

	// Test deleting the note
	require.NoError(t, repo.Delete(ctx, note.Id))

	// Test fetching the deleted note
	fetchedNote, err := repo.Fetch(ctx, note.Id)
	require.Nil(t, fetchedNote)
	require.NotNil(t, err)
}

func TestFetchAll(t *testing.T) {
	repo, note, ctx := setupC(t)

	require.NoError(t, repo.Create(ctx, note))

	// Test fetching all notes for the user
	fetchedNotes, err := repo.FetchAll(ctx, note.AuthorId)
	require.NoError(t, err)
	require.Equal(t, len(fetchedNotes), 1)
}

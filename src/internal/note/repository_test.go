package note

import (
	"context"
	"notes-manager/src/internal/user"
	"notes-manager/src/usecase/tester"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupR(t *testing.T) (Repository, *Note, context.Context) {
	ctx := context.Background()

	db := tester.SetupPostgresql(t)

	repo := NewDatabaseRepo(db)
	uRepo := user.NewDatabaseRepo(db)

	u, err := user.New("login", "password")
	require.NoError(t, err)

	err = uRepo.Create(ctx, u)
	require.NoError(t, err)

	return repo, New(u.Id, "my note", "my content", true), ctx
}

func TestNoteFetchingNonExisting(t *testing.T) {
	repo, note, ctx := setupR(t)

	fetchedNote, err := repo.Fetch(ctx, note.Id)
	require.Nil(t, fetchedNote)
	require.NotNil(t, err)
}

func TestNoteCreationAndFetch(t *testing.T) {
	repo, note, ctx := setupR(t)

	require.NoError(t, repo.Create(ctx, note))

	fetchedNote, err := repo.Fetch(ctx, note.Id)
	require.NoError(t, err)
	require.Equal(t, note, fetchedNote)
}

func TestNoteUpdate(t *testing.T) {
	repo, note, ctx := setupR(t)

	require.NoError(t, repo.Create(ctx, note))

	dto := NoteDTO{
		Id:        note.Id,
		Title:     "new title",
		Content:   "new content",
		UpdatedAt: time.Now().UTC(),
	}

	require.NoError(t, repo.Update(ctx, &dto))

	fetchedNote, err := repo.Fetch(ctx, note.Id)
	require.NoError(t, err)
	require.Equal(t, dto.Content, fetchedNote.Content)
}

func TestNoteFetchAll(t *testing.T) {
	repo, note, ctx := setupR(t)

	require.NoError(t, repo.Create(ctx, note))

	fetchedNotes, err := repo.FetchAll(ctx, note.AuthorId)
	require.NoError(t, err)
	require.Equal(t, len(fetchedNotes), 1)
}

func TestNoteDelete(t *testing.T) {
	repo, n, ctx := setupR(t)

	require.NoError(t, repo.Create(ctx, n))

	require.NoError(t, repo.Delete(ctx, n.Id))

	fetchedNote, err := repo.Fetch(ctx, n.Id)
	require.Nil(t, fetchedNote)
	require.NotNil(t, err)
}

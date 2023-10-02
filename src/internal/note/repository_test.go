package note

import (
	"context"
	"notes-manager/src/pkg/fsscanner"
	"notes-manager/src/pkg/migrate"
	"notes-manager/src/usecase/dcontainer"
	"notes-manager/src/usecase/repository/pgconnector"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_newDatabaseRepo(t *testing.T) {
	ctx := context.Background()

	postgresqlC, err := dcontainer.NewPostgres(ctx)
	require.NoError(t, err)

	defer func() {
		if err := postgresqlC.Terminate(ctx); err != nil {
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

	defer migrate.ApplyMigrations(&cfg, true, "file://"+migrationsPath)

	repo := NewDatabaseRepo(db)

	n := New(uuid.MustParse("07f3c5a1-70ea-4e3f-b9b5-110d29891673"), "my note", "my content")

	fetchedNote, err := repo.Fetch(ctx, n.Id)
	require.Nil(t, fetchedNote)
	require.NotNil(t, err)

	err = repo.Create(ctx, n)
	require.NoError(t, err)

	fetchedNote, err = repo.Fetch(ctx, n.Id)
	require.NoError(t, err)
	require.Equal(t, n, fetchedNote)

	dto := NoteDTO{
		Id:        n.Id,
		Title:     "new title",
		Content:   "new content",
		UpdatedAt: time.Now().UTC(),
	}

	err = repo.Update(ctx, &dto)
	require.NoError(t, err)

	fetchedNote, err = repo.Fetch(ctx, n.Id)
	require.NoError(t, err)
	require.Equal(t, dto.Content, fetchedNote.Content)

	fetchedNotes, err := repo.FetchAll(ctx, n.AuthorId)
	require.NoError(t, err)
	require.Equal(t, len(fetchedNotes), 1)

	err = repo.Delete(ctx, n.Id)
	require.NoError(t, err)

	fetchedNote, err = repo.Fetch(ctx, n.Id)
	require.Nil(t, fetchedNote)
	require.NotNil(t, err)

	fetchedNotes, err = repo.FetchAll(ctx, n.AuthorId)
	require.NoError(t, err)
	require.Equal(t, len(fetchedNotes), 0)
}

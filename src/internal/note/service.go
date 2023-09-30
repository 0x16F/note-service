package note

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func getNoteKey(noteId uuid.UUID) string {
	return fmt.Sprintf("ns:notes:%s", noteId)
}

func getNotesKey(userId uuid.UUID) string {
	return fmt.Sprintf("ns:notes:users:%d", userId)
}

func New(authorId uuid.UUID, title string, content string) *Note {
	currentTime := time.Now().UTC()

	return &Note{
		Id:        uuid.New(),
		AuthorId:  authorId,
		Title:     title,
		Content:   content,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func (Note) TableName() string {
	return "ns_notes"
}

func (NoteDTO) TableName() string {
	return Note{}.TableName()
}

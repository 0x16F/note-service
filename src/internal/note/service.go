package note

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// getNoteKey generates a unique cache key for a specific note based on its ID.
func getNoteKey(noteId uuid.UUID) string {
	return fmt.Sprintf("%s:%s", notesKeyBase, noteId)
}

// getNotesKey generates a unique cache key for all notes of a specific user based on their ID.
func getNotesKey(userId uuid.UUID) string {
	return fmt.Sprintf("%s:%s", notesUserIdKeyBase, userId)
}

// New initializes a new Note instance with the provided details and sets the creation and update times to the current UTC time.
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

// TableName specifies the table name for the Note struct when used with an ORM like Gorm.
func (Note) TableName() string {
	return "ns_notes"
}

// TableName specifies the table name for the NoteDTO struct by reusing the Note struct's table name.
func (NoteDTO) TableName() string {
	return Note{}.TableName()
}

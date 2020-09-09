package notes

import (
	"github.com/google/uuid"
)

type Note struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type NotesRepository interface {
	GetNotes() ([]Note, error)
	AddNote(note *Note) error
	RemoveNote(note *Note) error
}

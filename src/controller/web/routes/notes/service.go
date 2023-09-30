package notes

import (
	"net/http"
	"notes-manager/src/controller/web/headers"
	"notes-manager/src/controller/web/responses"
	"notes-manager/src/internal/note"
	"notes-manager/src/usecase/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func New(repo *repository.Repository) *Router {
	return &Router{
		repo:      repo,
		validator: validator.New(),
	}
}

func (r *Router) FetchAll(c *fiber.Ctx) error {
	s := headers.GetSession(c)

	notes, err := r.repo.Notes.FetchAll(c.Context(), s.UserId)
	if err != nil {
		logrus.Error(err)
		return responses.System("failed to fetch notes", err.Error())
	}

	return c.JSON(fiber.Map{
		"notes": notes,
	})
}

func (r *Router) Fetch(c *fiber.Ctx) error {
	noteId, err := uuid.Parse(c.Params("note_id"))
	if err != nil {
		return responses.BadRequest("failed to parse note id, type must be uuid", err.Error())
	}

	n, err := r.repo.Notes.Fetch(c.Context(), noteId)
	if err != nil {
		if err == note.ErrNoteIsNotExists {
			return responses.New(http.StatusNotFound, "note is not exists", err.Error())
		}

		logrus.Error(err)
		return responses.System("failed to fetch note", err.Error())
	}

	s := headers.GetSession(c)

	if s.UserId != n.AuthorId {
		return responses.New(http.StatusForbidden, "you don't have enough permissions to see this note", nil)
	}

	return c.JSON(fiber.Map{
		"note": n,
	})
}

func (r *Router) Create(c *fiber.Ctx) error {
	request := CreateNoteRequest{}
	if err := c.BodyParser(&request); err != nil {
		return responses.BadRequest("bad json", err.Error())
	}

	if err := r.validator.StructCtx(c.Context(), &request); err != nil {
		return responses.BadRequest("failed to validate some fields", err.Error())
	}

	s := headers.GetSession(c)
	n := note.New(s.UserId, request.Title, request.Content)

	if err := r.repo.Notes.Create(c.Context(), n); err != nil {
		return responses.System("failed to create new note", err.Error())
	}

	return c.JSON(fiber.Map{
		"note": n,
	})
}

func (r *Router) Delete(c *fiber.Ctx) error {
	noteId, err := uuid.Parse(c.Params("note_id"))
	if err != nil {
		return responses.BadRequest("failed to parse note id, type must be uuid", err.Error())
	}

	n, err := r.repo.Notes.Fetch(c.Context(), noteId)
	if err != nil {
		if err == note.ErrNoteIsNotExists {
			return responses.New(http.StatusNotFound, "note is not exists", err.Error())
		}

		logrus.Error(err)
		return responses.System("failed to fetch note", err.Error())
	}

	s := headers.GetSession(c)

	if s.UserId != n.AuthorId {
		return responses.New(http.StatusForbidden, "you don't have enough permissions to delete this note", nil)
	}

	if err := r.repo.Notes.Delete(c.Context(), noteId); err != nil {
		logrus.Error(err)
		return responses.System("failed to delete note", err.Error())
	}

	return nil
}

func (r *Router) Update(c *fiber.Ctx) error {
	request := UpdateNoteRequest{}
	if err := c.BodyParser(&request); err != nil {
		return responses.BadRequest("bad json", err.Error())
	}

	if err := r.validator.StructCtx(c.Context(), &request); err != nil {
		return responses.BadRequest("failed to validate some fields", err.Error())
	}

	n, err := r.repo.Notes.Fetch(c.Context(), request.NoteId)
	if err != nil {
		if err == note.ErrNoteIsNotExists {
			return responses.New(http.StatusNotFound, "note is not exists", err.Error())
		}

		logrus.Error(err)
		return responses.System("failed to fetch note", err.Error())
	}

	s := headers.GetSession(c)

	if s.UserId != n.AuthorId {
		return responses.New(http.StatusForbidden, "you don't have enough permissions to update this note", nil)
	}

	dto := note.NoteDTO{
		Id:        request.NoteId,
		Title:     request.Title,
		Content:   request.Content,
		UpdatedAt: time.Now().UTC(),
	}

	if err := r.repo.Notes.Update(c.Context(), &dto); err != nil {
		logrus.Error(err)
		return responses.System("failed to update note", err.Error())
	}

	return nil
}
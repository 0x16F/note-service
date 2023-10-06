package notes

import (
	"math"
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
		headers:   headers.New(),
	}
}

// @Summary      Fetch all notes
// @Description  returning all notes of current user
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param 		 public query bool false "Fetch all public notes" default(false)
// @Param 		 page query int false "Page" minimum(1) default(1)
// @Param 		 per_page query int false "Notes per page" minimum(1) maximum(50) default(50)
// @Success      200  {object}  UserNotesResponse
// @Success      400  {object}  UserNotesResponse
// @Failure      500  {object}  responses.Error
// @Router       /v0/notes [get]
func (r *Router) FetchAll(c *fiber.Ctx) error {
	params := FetchAllParams{}
	if err := c.QueryParser(&params); err != nil {
		return responses.BadRequest("failed to parse queries", err.Error())
	}

	if params.Public {
		if err := r.validator.StructCtx(c.Context(), &params); err != nil {
			return responses.BadRequest("some fields is invalid", err.Error())
		}

		notes, count, err := r.repo.Notes.FetchPublic(c.Context(), params.Page, params.PerPage)
		if err != nil {
			logrus.Error(err)
			return responses.System("failed to fetch public notes", err.Error())
		}

		pages := int64(math.Ceil(float64(count) / float64(params.PerPage)))

		return c.JSON(&UserNotesResponse{
			Notes: notes,
			Pages: pages,
		})
	}

	s := r.headers.GetSession(c)

	notes, err := r.repo.Notes.FetchAll(c.Context(), s.UserId)
	if err != nil {
		logrus.Error(err)
		return responses.System("failed to fetch notes", err.Error())
	}

	return c.JSON(&UserNotesResponse{
		Notes: notes,
		Pages: 0,
	})
}

// @Summary      Fetch note
// @Description  returning specific note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param		 id   path		string true "Note ID"
// @Success      200  {object}  NoteResponse
// @Failure      403  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Failure      500  {object}  responses.Error
// @Router       /v0/notes/{id} [get]
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

	s := r.headers.GetSession(c)

	if s.UserId != n.AuthorId {
		return responses.New(http.StatusForbidden, "you don't have enough permissions to see this note", nil)
	}

	return c.JSON(&NoteResponse{
		Note: n,
	})
}

// @Summary      Create note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        request	body	CreateNoteRequest  true  "create note params"
// @Success      200  {object}  NoteResponse
// @Failure      400  {object}  responses.Error
// @Failure      500  {object}  responses.Error
// @Router       /v0/notes [post]
func (r *Router) Create(c *fiber.Ctx) error {
	request := CreateNoteRequest{}
	if err := c.BodyParser(&request); err != nil {
		return responses.BadRequest("bad json", err.Error())
	}

	if err := r.validator.StructCtx(c.Context(), &request); err != nil {
		return responses.BadRequest("failed to validate some fields", err.Error())
	}

	s := r.headers.GetSession(c)
	n := note.New(s.UserId, request.Title, request.Content, true)

	if err := r.repo.Notes.Create(c.Context(), n); err != nil {
		logrus.Error(err)
		return responses.System("failed to create new note", err.Error())
	}

	return c.JSON(&NoteResponse{
		Note: n,
	})
}

// @Summary      Delete note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param		 id   path		string true "Note ID"
// @Success      200
// @Failure      400  {object}  responses.Error
// @Failure      403  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Failure      500  {object}  responses.Error
// @Router       /v0/notes/{id} [delete]
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

	s := r.headers.GetSession(c)

	if s.UserId != n.AuthorId {
		return responses.New(http.StatusForbidden, "you don't have enough permissions to delete this note", nil)
	}

	if err := r.repo.Notes.Delete(c.Context(), noteId); err != nil {
		logrus.Error(err)
		return responses.System("failed to delete note", err.Error())
	}

	return nil
}

// @Summary      Update note
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        request	body	UpdateNoteRequest  true  "update note params"
// @Success      200  {object}  NoteResponse
// @Failure      400  {object}  responses.Error
// @Failure      403  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Failure      500  {object}  responses.Error
// @Router       /v0/notes [patch]
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

	s := r.headers.GetSession(c)

	if s.UserId != n.AuthorId {
		return responses.New(http.StatusForbidden, "you don't have enough permissions to update this note", nil)
	}

	dto := note.NoteDTO{
		Id:        request.NoteId,
		Title:     request.Title,
		Content:   request.Content,
		UpdatedAt: time.Now().UTC(),
		IsPrivate: request.IsPrivate,
	}

	if err := r.repo.Notes.Update(c.Context(), &dto); err != nil {
		logrus.Error(err)
		return responses.System("failed to update note", err.Error())
	}

	n, err = r.repo.Notes.Fetch(c.Context(), request.NoteId)
	if err != nil {
		if err == note.ErrNoteIsNotExists {
			return responses.New(http.StatusNotFound, "note is not exists", err.Error())
		}

		logrus.Error(err)
		return responses.System("failed to fetch note", err.Error())
	}

	return c.JSON(&NoteResponse{
		Note: n,
	})
}

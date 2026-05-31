// Имплементация сервиса NoteService.
package usecase

import (
	"context"
	"errors"
	"log/slog"

	auth "webnote/internal/auth/models"
	"webnote/internal/note/models"
	"webnote/internal/note/repository"
)

// Структура NoteService с репозиторием PostgresNoteRepository.
type NoteService struct {
	repo *repository.PostgresNoteRepository
}

// Создает новый экземпляр NoteService.
func NewNoteService(r *repository.PostgresNoteRepository) *NoteService {
	return &NoteService{repo: r}
}

// Создание записки.
func (s *NoteService) CreateNote(ctx context.Context, note CreateNoteReq) (int64, error) {
	const op = "note.usecase.service.CreateNote"
	currentUser, ok := auth.CurrentUserFromContext(ctx)
	if !ok {
		slog.Debug("context key 'currentUser' are empty", "op", op)
		return 0, errors.New("not authorized")
	}

	in := models.Note{
		UserID: currentUser.ID,
		Name:   note.Name,
		Text:   note.Text,
	}
	return s.repo.CreateNote(ctx, in)
}

func (s *NoteService) GetNote(ctx context.Context, noteID int64) (NoteRes, error) {
	const op = ""
	var out NoteRes
	currentUser, ok := auth.CurrentUserFromContext(ctx)
	if !ok {
		slog.Debug("context key 'currentUser' are empty", "op", op)
		return NoteRes{}, errors.New("not authorized")
	}
	note, err := s.repo.GetNote(ctx, noteID, currentUser.ID)
	if err != nil {
		return NoteRes{}, err
	}
	out = NoteRes{
		ID:        note.ID,
		UserID:    note.UserID,
		Name:      note.Name,
		Text:      note.Text,
		CreatedAt: note.CreatedAt,
	}
	return out, nil
}

// Получение всех записок.
func (s *NoteService) GetMyNotes(ctx context.Context) ([]NoteRes, error) {
	const op = "note.usecase.service.GetMyNotes"
	currentUser, ok := auth.CurrentUserFromContext(ctx)
	if !ok {
		slog.Debug("context key 'currentUser' are empty", "op", op)
		return nil, errors.New("not authorized")
	}
	notes, err := s.repo.GetMyNotes(ctx, currentUser.ID)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return nil, err
	}
	var out []NoteRes
	for _, n := range notes {
		tmp := NoteRes{
			ID:        n.ID,
			UserID:    n.UserID,
			Name:      n.Name,
			Text:      n.Text,
			CreatedAt: n.CreatedAt,
		}
		out = append(out, tmp)
	}
	return out, nil
}

func (s *NoteService) UpdateNoteByID(ctx context.Context, noteID int64, newNote UpdateNoteReq) (NoteRes, error) {
	const op = "note.usecase.service.UpdateNoteByID"
	currentUser, ok := auth.CurrentUserFromContext(ctx)
	if !ok {
		slog.Debug("context key 'currentUser' are empty", "op", op)
		return NoteRes{}, errors.New("not authorized")
	}
	in := models.Note{
		Name: newNote.Name,
		Text: newNote.Text,
	}
	note, err := s.repo.UpdateNote(ctx, noteID, currentUser.ID, in)
	if err != nil {
		slog.Debug("context key 'currentUser' are empty", "op", op)
		return NoteRes{}, errors.New("not authorized")
	}
	out := NoteRes{
		ID:        note.ID,
		UserID:    note.UserID,
		Name:      note.Name,
		Text:      note.Text,
		CreatedAt: note.CreatedAt,
	}
	return out, nil
}

func (s *NoteService) DeleteNoteByID(ctx context.Context, noteID int64) (NoteIDRes, error) {
	const op = "note.usecase.service.DeleteNoteByID"
	currentUser, ok := auth.CurrentUserFromContext(ctx)
	if !ok {
		slog.Debug("context key 'currentUser' are empty", "op", op)
		return NoteIDRes{0}, errors.New("not authorized")
	}
	id, err := s.repo.DeleteNoteByID(ctx, noteID, currentUser.ID)
	if err != nil {
		return NoteIDRes{0}, err
	}
	return NoteIDRes{id}, nil
}

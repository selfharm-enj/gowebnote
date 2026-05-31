// Postgresql имплементация NoteRepository.
package repository

import (
	"context"
	"errors"
	"log/slog"
	"webnote/internal/note/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Структура реализующая NoteRepository.
type PostgresNoteRepository struct {
	db *pgxpool.Pool
}

// Создает новый экземпляр PostgresNoteRepository.
func NewPostgresNoteRepository(db *pgxpool.Pool) *PostgresNoteRepository {
	return &PostgresNoteRepository{db: db}
}

// Создание новой записки.
func (r *PostgresNoteRepository) CreateNote(ctx context.Context, note models.Note) (int64, error) {
	const op = "note.repository.postgres.CreateNote"
	stmt := `
	INSERT INTO
		notes(user_id, name, text)
	VALUES
		($1, $2, $3)
	RETURNING
		id
	`
	var out int64
	row := r.db.QueryRow(ctx, stmt, note.UserID, note.Name, note.Text)
	if err := row.Scan(&out); err != nil {
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				slog.Debug("unique constraint violation", "op", op)
			case pgerrcode.NotNullViolation:
				slog.Debug("not null violation", "op", op)
			default:
				slog.Debug("uncatched postgres error", "err", err.Error(), "op", op)
			}
		}
		return 0, err
	}
	slog.Debug("note created", "note_id", out, "op", op)
	return out, nil
}

// Получение записки юзера 'userID' по 'noteID'.
func (r *PostgresNoteRepository) GetNote(ctx context.Context, noteID int64, userID int64) (models.Note, error) {
	const op = "note.repository.postgres.GetNoteByID"
	stmt := `
	SELECT
		id, user_id, name, text, created_at
	FROM
		notes
	WHERE
		id = $1
		AND
		user_id = $2
	`
	var out models.Note
	row := r.db.QueryRow(ctx, stmt, noteID, userID)
	if err := row.Scan(&out.ID, &out.UserID, &out.Name, &out.Text, &out.CreatedAt); err != nil {
		slog.Debug(err.Error(), "op", op)
		return models.Note{}, err
	}
	return out, nil
}

// Получение всех записок принадлежащих пользователю с 'userID'.
func (r *PostgresNoteRepository) GetMyNotes(ctx context.Context, userID int64) ([]models.Note, error) {
	const op = "note.repository.postgres.GetMyNotes"
	stmt := `
	SELECT
		id, user_id, name, text, created_at
	FROM
		notes
	WHERE
		user_id = $1
	`
	rows, err := r.db.Query(ctx, stmt, userID)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return nil, err
	}
	defer rows.Close()
	var notes []models.Note
	notes, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.Note])
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return nil, err
	}
	return notes, nil
}

// Обновление содержимого записки.
func (r *PostgresNoteRepository) UpdateNote(ctx context.Context, noteID int64,
	userID int64, newNote models.Note) (models.Note, error) {
	const op = "note.repository.postgres.UpdateNoteByID"
	stmt := `
	UPDATE
		notes
	SET
		name = $1,
		text = $2
	WHERE
		id = $3
		AND
		user_id = $4
	RETURNING
		id, user_id, name, text, created_at
	`
	var out models.Note
	row := r.db.QueryRow(ctx, stmt, newNote.Name, newNote.Text, noteID, userID)
	err := row.Scan(&out.ID, &out.UserID, &out.Name, &out.Text, &out.CreatedAt)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return models.Note{}, err
	}
	return out, nil
}

// Удаление записки по 'noteID'.
func (r *PostgresNoteRepository) DeleteNoteByID(ctx context.Context, noteID int64, userID int64) (int64, error) {
	const op = "note.repository.postgres.DeleteNoteByID"
	stmt := `
	DELETE FROM
		notes
	WHERE
		id = $1
		AND
		user_id = $2
	RETURNING
		id
	`
	var out int64
	row := r.db.QueryRow(ctx, stmt, noteID, userID)
	if err := row.Scan(&out); err != nil {
		slog.Debug(err.Error(), "op", op)
		return 0, err
	}
	return out, nil
}

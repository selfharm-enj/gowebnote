// Postgresql имплементация UserRepository.
package repository

import (
	"context"
	"errors"
	"log/slog"
	"webnote/internal/user/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Структура реализующая UserRepository.
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// Создает новый экземпляр PostgresUserRepository.
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Создание нового пользователя.
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user models.User) (int64, error) {
	const op = "user.repository.postgres.CreateUser"
	stmt := `
	INSERT INTO
		users(email, password)
	VALUES
		($1, $2)
	RETURNING
		id
	`
	var out int64
	row := r.db.QueryRow(ctx, stmt, user.Email, user.Password)
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
	slog.Debug("user created", "user_id", out, "op", op)
	return out, nil
}

var ErrUserNotExists = errors.New("user with passed email doesn't exists")

// Получение пользователя по email.
func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	stmt := `
	SELECT
		id, email, password, created_at
	FROM
		users
	WHERE
		email = $1
	`
	var out models.User
	row := r.db.QueryRow(ctx, stmt, email)
	if err := row.Scan(&out.ID, &out.Email, &out.Password, &out.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotExists
		}
		return models.User{}, err
	}
	return out, nil
}

// Получение пользователя по ID.
func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	const op = "user.repository.postgres.GetUserByID"
	stmt := `
	SELECT
		id, email, password, created_at
	FROM
		users
	WHERE
		id = $1
	`
	var out models.User
	row := r.db.QueryRow(ctx, stmt, id)
	if err := row.Scan(&out.ID, &out.Email, &out.Password, &out.CreatedAt); err != nil {
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
		return models.User{}, err
	}
	return out, nil
}

// Получение сессии по 'session_token'.
func (r *PostgresUserRepository) GetSessionByToken(ctx context.Context, sessionToken string) (models.Session, error) {
	const op = "user.repository.postgres.GetSessionByToken"
	stmt := `
	SELECT
		id, user_id, session_token, csrf_token, max_age, expire_at
	FROM
		sessions
	WHERE
		session_token = $1
	`
	var out models.Session
	row := r.db.QueryRow(ctx, stmt, sessionToken)
	err := row.Scan(&out.ID, &out.UserID, &out.SessionToken, &out.CSRFToken, &out.MaxAge, &out.ExpireAt)
	if err != nil {
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
		return models.Session{}, err
	}
	return out, err
}

// Создание сессии.
func (r *PostgresUserRepository) CreateSession(ctx context.Context, session models.Session) (models.Session, error) {
	const op = "user.repository.postgres.CreateSession"
	stmt := `
	INSERT INTO
		sessions(user_id, session_token, csrf_token, max_age, expire_at)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING
		id, user_id, session_token, csrf_token, max_age, expire_at
	`
	var out models.Session
	row := r.db.QueryRow(ctx, stmt,
		session.UserID,
		session.SessionToken,
		session.CSRFToken,
		session.MaxAge,
		session.ExpireAt)
	err := row.Scan(&out.ID, &out.UserID, &out.SessionToken, &out.CSRFToken, &out.MaxAge, &out.ExpireAt)
	if err != nil {
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
		return models.Session{}, err
	}
	return out, nil
}

// Удаление сессии по 'userID'.
func (r *PostgresUserRepository) DeleteSessionByUserID(ctx context.Context, userID int64) (int64, error) {
	const op = "user.repository.postgres.DeleteSessionByUserID"
	stmt := `
	DELETE FROM
		sessions
	WHERE
		user_id = $1
	RETURNING
		id
	`
	var out int64
	row := r.db.QueryRow(ctx, stmt, userID)
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
	return out, nil
}

// Интерфейс User и Session.
package repository

import (
	"context"

	"webnote/internal/user/models"
)

// Интерфейс для работы с User и Session.
type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetSessionByToken(ctx context.Context, sessionToken string) (models.Session, error)
	CreateSession(ctx context.Context, session models.Session) (models.Session, error)
	DeleteSessionByUserID(ctx context.Context, userID int64) (int64, error)
}

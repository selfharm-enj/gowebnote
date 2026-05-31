// Имплементация сервиса UserService.
package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"time"

	auth "webnote/internal/auth/models"
	"webnote/internal/user/models"
	"webnote/internal/user/repository"
	"webnote/utils"
)

// Структура UserService с репозиторием PostgresUserRepository.
type UserService struct {
	repo   *repository.PostgresUserRepository
	hasher *utils.BcryptHasher
}

// Создает новый экземпляр UserService.
func NewUserService(r *repository.PostgresUserRepository, h *utils.BcryptHasher) *UserService {
	return &UserService{repo: r, hasher: h}
}

// Создание пользователя.
func (s *UserService) CreateUser(ctx context.Context, in UserCredentialsReq) (int64, error) {
	const op = "user.usecase.service.CreateUser"
	user, err := s.repo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotExists) {
			slog.Debug("err while getting user by email", "err", err.Error(), "op", op)
			return 0, err
		}
	}
	if user.ID != 0 {
		slog.Debug("email already taken", "op", op)
		return 0, errors.New("email already taken")
	}
	hash, err := s.hasher.HashPassword(in.Password)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return 0, err
	}
	newUser := models.User{
		Email:    in.Email,
		Password: hash,
	}
	id, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
	}
	return id, err
}

// Вход в УЗ пользователя.
func (s *UserService) Login(ctx context.Context, in UserCredentialsReq) (UserSessionRes, error) {
	const op = "user.usecase.service.Login"
	user, err := s.repo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return UserSessionRes{}, err
	}
	if user.ID == 0 {
		slog.Debug("user with target email doesn't exists", "op", op)
		return UserSessionRes{}, errors.New("user with target email doesn't exists")
	}
	if err := s.hasher.Compare(user.Password, in.Password); err != nil {
		slog.Debug(err.Error(), "op", op)
		return UserSessionRes{}, errors.New("wrong password")
	}
	sessionToken := rand.Text()
	csrfToken := rand.Text()
	newSess := models.Session{
		UserID:       user.ID,
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		MaxAge:       int(time.Hour.Seconds()),
		ExpireAt:     time.Now().Add(time.Hour),
	}
	userSess, err := s.repo.CreateSession(ctx, newSess)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return UserSessionRes{}, err
	}
	out := UserSessionRes{
		UserID:       userSess.UserID,
		SessionToken: userSess.SessionToken,
		CSRFToken:    userSess.CSRFToken,
		MaxAge:       userSess.MaxAge,
		ExpireAt:     userSess.ExpireAt,
	}
	return out, nil
}

// Выход из УЗ пользователя.
func (s *UserService) Logout(ctx context.Context, userID int64) (int64, error) {
	const op = "user.usecase.service.Logout"
	out, err := s.repo.DeleteSessionByUserID(ctx, userID)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return 0, err
	}
	return out, nil
}

// Получение сессии по токену.
func (s *UserService) GetSessionByToken(ctx context.Context, sessionToken string) (UserSessionRes, error) {
	const op = "user.usecase.service.GetSessionByToken"
	sess, err := s.repo.GetSessionByToken(ctx, sessionToken)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return UserSessionRes{}, err
	}
	out := UserSessionRes{
		UserID:       sess.UserID,
		SessionToken: sess.SessionToken,
		CSRFToken:    sess.CSRFToken,
		MaxAge:       sess.MaxAge,
		ExpireAt:     sess.ExpireAt,
	}
	return out, nil
}

// Валидация сессии по 'sessionToken'.
func (s *UserService) ValidateSession(ctx context.Context, sessionToken string) (auth.CurrentUser, error) {
	const op = "user.usecase.service.ValidateSession"
	session, err := s.repo.GetSessionByToken(ctx, sessionToken)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return auth.CurrentUser{}, err
	}
	if session.ExpireAt.Before(time.Now()) {
		slog.Debug("user session already expired", "op", op)
		return auth.CurrentUser{}, nil
	}
	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		slog.Debug(err.Error(), "op", op)
		return auth.CurrentUser{}, err
	}
	out := auth.CurrentUser{
		ID:           user.ID,
		Email:        user.Email,
		SessionID:    session.ID,
		SessionToken: session.SessionToken,
		CSRFToken:    session.CSRFToken,
	}
	return out, nil
}

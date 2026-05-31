package usecase

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrBadUserID          = errors.New("user_id should be more than 0")
	ErrEmailFieldEmpty    = errors.New("email field are empty")
	ErrInvalidEmail       = errors.New("invalid email passed")
	ErrPasswordFieldEmpty = errors.New("password field are empty")
)

type UserCredentialsReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserCredentialsReq) Validate() error {
	email := strings.TrimSpace(r.Email)
	password := strings.TrimSpace(r.Password)
	if len(email) == 0 {
		return ErrEmailFieldEmpty
	}
	if len(password) == 0 {
		return ErrPasswordFieldEmpty
	}
	if !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}
	return nil
}

type UserIDRes struct {
	ID int64 `json:"id"`
}

type UserSessionRes struct {
	UserID       int64     `json:"user_id"`
	SessionToken string    `json:"session_token"`
	CSRFToken    string    `json:"csrf_token"`
	MaxAge       int       `json:"max_age"`
	ExpireAt     time.Time `json:"expire_at"`
}

type UserAuthResultRes struct {
	Status string `json:"auth_status"`
}

type DeleteSessionReq struct {
	UserID int64 `json:"user_id"`
}

func (r *DeleteSessionReq) Validate() error {
	if r.UserID == 0 {
		return ErrBadUserID
	}
	return nil
}

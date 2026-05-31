// HTTP handlers для работы с сервисом UserService.
package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"webnote/internal/user/usecase"
	"webnote/utils"
)

// Handler имплементирующий UserService.
type UserHandler struct {
	service *usecase.UserService
}

// Создание нового экземпляра 'UserHandler'.
func NewUserHandler(s *usecase.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Создание пользователя.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const op = "user.transport.http.handler.CreateUser"
	var req usecase.UserCredentialsReq
	if err := utils.ReadJSON(r, &req); err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	if err := req.Validate(); err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	id, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, usecase.UserIDRes{ID: id})
}

// Вход в УЗ пользователя.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "user.transport.http.handler.Login"
	var req usecase.UserCredentialsReq
	if err := utils.ReadJSON(r, &req); err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	if err := req.Validate(); err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	session, err := h.service.Login(r.Context(), req)
	if err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	// #nosec G124
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(int(session.UserID)),
		MaxAge:   session.MaxAge,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
		Path:     "/",
	})
	// #nosec G124
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		MaxAge:   session.MaxAge,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Path:     "/",
	})
	// #nosec G124
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    session.CSRFToken,
		MaxAge:   session.MaxAge,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
		Path:     "/",
	})
	utils.WriteJSON(w, http.StatusOK, usecase.UserAuthResultRes{Status: "ok"})
}

// Выход из УЗ пользователя.
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "user.transport.http.handler.Logout"
	var req usecase.DeleteSessionReq
	if err := utils.ReadJSON(r, &req); err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	if err := req.Validate(); err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorBadRequest(w)
		return
	}
	id, err := h.service.Logout(r.Context(), req.UserID)
	if err != nil {
		slog.Debug("", "err", err.Error(), "op", op)
		utils.WriteJSONErrorInternalServer(w)
		return
	}
	// #nosec G124
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
	})
	// #nosec G124
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})
	// #nosec G124
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
	})
	utils.WriteJSON(w, http.StatusOK, usecase.UserIDRes{ID: id})
}

// Имплементация middleware для аутентификации пользователей.
package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	authmodels "webnote/internal/auth/models"
	useruc "webnote/internal/user/usecase"
)

// Интерфейс для валидации сессий.
type SessionValidator interface {
	ValidateSession(ctx context.Context, sessionToken string) (authmodels.CurrentUser, error)
}

// Middleware управления доступа к ресурсам для которых требуетс  аутентификация.
func RequireAuth(next http.HandlerFunc, userService *useruc.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessToken, err := r.Cookie("session_token")
		if err != nil || sessToken.Value == "" {
			slog.Debug("auth.middleware.RequireAuth")
			fmt.Fprintf(w, "%d %s", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			// utils.WriteJSONError(w, http.StatusUnauthorized, "not authorized")
			return
		}
		csrfToken, err := r.Cookie("csrf_token")
		if err != nil || csrfToken.Value == "" {
			slog.Debug("auth.middleware.RequireAuth")
			fmt.Fprintf(w, "%d %s", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			// utils.WriteJSONError(w, http.StatusUnauthorized, "not authorized")
			return
		}
		currentUser, err := userService.ValidateSession(r.Context(), sessToken.Value)
		if err != nil {
			slog.Debug("user unathorized", "func", "auth.middleware.RequireAuth")
			fmt.Fprintf(w, "%d %s", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			// utils.WriteJSONError(w, http.StatusUnauthorized, "not authorized")
			return
		}
		ctx := authmodels.WithCurrentUser(r.Context(), currentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

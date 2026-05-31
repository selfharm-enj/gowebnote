// Создание контекста для пакета auth.
package auth

import "context"

type ctxKeyType string

const currentUserKey ctxKeyType = "currentUser"

// Возвращает контекст с значением 'CurrentUser'.
func WithCurrentUser(ctx context.Context, u CurrentUser) context.Context {
	return context.WithValue(ctx, currentUserKey, u)
}

// Возвращает 'CurrentUser' из контекста.
func CurrentUserFromContext(ctx context.Context) (CurrentUser, bool) {
	u, ok := ctx.Value(currentUserKey).(CurrentUser)
	return u, ok
}

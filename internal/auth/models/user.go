// Имплементация структуры CurrentUser.
package auth

// CurrentUser.
type CurrentUser struct {
	ID           int64
	Email        string
	SessionID    int64
	SessionToken string
	CSRFToken    string
}

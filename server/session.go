package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Checks whether the user has a valid session.
func CheckSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405, http.StatusText(http.StatusMethodNotAllowed), "Only GET method is allowed!", nil)
	}

	w.Header().Set("Content-Type", "application/json")

	user, err := GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"loggedIn": false}`)
		return
	}
	fmt.Fprintf(w, `{"loggedIn": true, "username": %q}`, user.Username)
}

// Get the user from the current session using cookies.
func GetUser(r *http.Request) (*User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, fmt.Errorf("no session token provided")
	}

	token := cookie.Value
	var session struct {
		UserID    int
		ExpiresAt time.Time
	}

	// Select user_id, expires_at from DB based on token value.
	err = DB.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token = ?`, token).Scan(&session.UserID, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired session token")
	}

	// Check if the session is expired.
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	// Fetch the user associated with the session from DB
	var user User
	err = DB.QueryRow(`SELECT id, email, username FROM users WHERE id = ?`, session.UserID).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

// Create session token (cookie) and insert it into DB.
func CreateSession(w http.ResponseWriter, user *User) error {
	token := uuid.New().String()

	// Set token expiration time.
	expiresAt := time.Now().Add(24 * time.Hour)

	// Limit concurrent sessions to only one per user
	_, err := DB.Exec(`DELETE FROM sessions WHERE user_id = ?`, user.ID)
	if err != nil {
		return err
	}

	// Insert session into DB.
	_, err = DB.Exec(`INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)`, user.ID, token, expiresAt)
	if err != nil {
		return err
	}

	// Set the token in a cookie
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Moderate CSRF protection, security & usability balance
		Secure:   true,                 // Required for HTTPS
		HttpOnly: true,                 // Protects against XSS, blocks access to document.cookie
	}

	http.SetCookie(w, cookie)
	return nil
}

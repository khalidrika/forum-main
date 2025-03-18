package server

import (
	"net/http"
	"time"
)

// Delete session and expire cookie
func clearSession(w http.ResponseWriter, token string) error {
	// Remove session from DB
	_, err := DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	if err != nil {
		return err
	}

	// Expire the cookie, immediate invalidation
	expiredCookie := &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, expiredCookie)
	return nil
}

// Handle log out functionality
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		JsonError(w, "No session token provided", http.StatusUnauthorized, err)
		return
	}

	token := cookie.Value
	if err := clearSession(w, token); err != nil {
		JsonError(w, "Failed to remove session", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

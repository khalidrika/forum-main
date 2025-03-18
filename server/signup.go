package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Signing up a new user.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed, nil)
		return
	}

	// Limit the size of the request body to 30 KB
	r.Body = http.MaxBytesReader(w, r.Body, 30000)

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		JsonError(w, "Invalid request payload or size exceeded", http.StatusBadRequest, err)
		return
	}

	if err := ValidateSignUp(user.Email, user.Username, user.Password); err != nil {
		JsonError(w, err.Error(), http.StatusNotAcceptable, err)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		JsonError(w, "unexpected error, try again later", http.StatusInternalServerError, err)
		return
	}

	insertUser := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`
	_, err = DB.Exec(insertUser, user.Email, user.Username, hashedPassword)
	if err != nil {
		JsonError(w, "unexpected error, try again later", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// Validate signup payload.
func ValidateSignUp(email, username, password string) error {
	if email == "" || username == "" || password == "" {
		return fmt.Errorf("email, username, and password are required")
	}

	// Validate Email.
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,20}$`
	if match, _ := regexp.MatchString(emailRegex, email); !match {
		return fmt.Errorf("invalid email format")
	}
	if len(email) > 200 {
		return fmt.Errorf("email is too long")
	}

	// Validate username.
	if len(username) < 3 {
		return fmt.Errorf("username is too short")
	}
	if len(username) > 16 {
		return fmt.Errorf("username is too long")
	}
	usernameRegex := `^[a-zA-Z0-9_.-]+$`
	if match, _ := regexp.MatchString(usernameRegex, username); !match {
		return fmt.Errorf("username can only contain letters, digits, underscores, dots, and hyphens")
	}

	// Validate password.
	if len(password) < 6 {
		return fmt.Errorf("password is too short")
	}
	if len(password) > 64 {
		return fmt.Errorf("password is too long")
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)
	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return fmt.Errorf("password must contain at least one lowercase letter, one uppercase letter, one digit, and one special character")
	}

	// Check if already used email/username.
	exists := false
	err := DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)`, email, username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("unexpected error, try again later")
	}
	if exists {
		return fmt.Errorf("email or username already exists")
	}

	return nil
}

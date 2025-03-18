package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Expected JSON structure for adding comments
type CommentPayload struct {
	PostID  int    `json:"id"`
	Content string `json:"content"`
}

// Handles adding a comment to DB
func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonError(w, "Method not allowed", 405, nil)
		return
	}
	// Limit the size of the request body to 8 KB
	r.Body = http.MaxBytesReader(w, r.Body, 8000)

	user, err := GetUser(r)
	if err != nil {
		JsonError(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized, err)
		return
	}

	var payload CommentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		JsonError(w, "Invalid request payload or size exceeded", http.StatusBadRequest, err)
		return
	}

	err = processComment(&payload, user)
	if err != nil {
		JsonError(w, err.Error(), http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Comment added successfully"))
}

// Validates, sanitizes and inserts the comment
func processComment(payload *CommentPayload, user *User) error {
	// Trim spaces and escape HTML
	payload.Content = strings.TrimSpace(payload.Content)
	payload.Content = html.EscapeString(payload.Content)

	// Remove more than three consecutive new lines to just two
	re := regexp.MustCompile(`(\r\n|\r|\n){3,}`)
	payload.Content = re.ReplaceAllString(payload.Content, "\n\n")

	// Check minimum length after sanitization
	if len(payload.Content) < 5 {
		return errors.New("comment is too short")
	}

	if payload.Content == "" {
		return errors.New("write something in your comment")
	}

	// Check if post exists
	var exists bool
	err := DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)`, payload.PostID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to verify post existence: %w", err)
	}
	if !exists {
		return errors.New("post does not exist")
	}

	// Save to database
	_, err = DB.Exec(`
        INSERT INTO comments (post_id, user_id, content)
        VALUES (?, ?, ?)
    `, payload.PostID, user.ID, payload.Content)
	if err != nil {
		return fmt.Errorf("failed to add comments: %w", err)
	}

	return nil
}

// Fetch all comments
func GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		JsonError(w, "No post id provided", http.StatusBadRequest, nil)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		JsonError(w, "Wrong offset", http.StatusBadRequest, err)
		return
	}

	rows, err := DB.Query(`
		SELECT c.id, c.content, c.created_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
		LIMIT 20 OFFSET ?
	`, postID, offset)
	if err != nil {
		JsonError(w, "Failed to query comments", http.StatusInternalServerError, err)
		return
	}

	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Username); err != nil {
			fmt.Fprintln(w, err.Error())
			continue
		}
		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// Fetch just the number of comments
func GetCommentsCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		JsonError(w, "Missing post_id", http.StatusBadRequest, nil)
		return
	}

	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id = ?`, postID).Scan(&count)
	if err != nil {
		JsonError(w, "Error querying comment count: "+err.Error(), http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}

package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Serve just the html of the post.
func ServePostHtml(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405, http.StatusText(http.StatusMethodNotAllowed), "Only GET method is allowed!", nil)
		return
	}
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		ErrorHandler(w, 400, "Bad Request", "Missing post id!", nil)
		return
	}

	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		ErrorHandler(w, 500, "Internal Server Error", "Database error occurred while validating post_id.", err)
		return
	}

	if !exists {
		ErrorHandler(w, 404, "Not Found", "Post with the specified ID does not exist.", nil)
		return
	}
	ParseAndExecute(w, "", "static/templates/post.html")
}

// Serve the Post json
func SinglePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		JsonError(w, "Missing post_id", http.StatusBadRequest, nil)
		return
	}
	var post Post

	// SELECT the columns that match post ID
	err := DB.QueryRow(`
        SELECT p.id, p.title, p.content, u.username, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?`,
		postID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Username,
		&post.CreatedAt,
	)
	if err != nil {
		JsonError(w, "Post not found or DB error: "+err.Error(), http.StatusNotFound, err)
		return
	}

	// Fetch post categories/tags
	pid, _ := strconv.Atoi(postID)
	cats, err := FetchPostCategories(pid)
	if err == nil {
		post.Categories = cats
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

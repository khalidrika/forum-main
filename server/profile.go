package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Num of posts on each scroll load in profile.
var ProfileLimit = 6

// Serve the profile html
func ServeProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed", "You can only use GET method", nil)
		return
	}

	user, err := GetUser(r)
	if err != nil {
		ErrorHandler(w, http.StatusUnauthorized, "Login or signup to view profile", "", nil)
		return
	}

	username := r.URL.Query().Get("user")
	if username == "" {
		ErrorHandler(w, http.StatusUnauthorized, "", "Invalid username", nil)
		return
	}

	if user.Username != username {
		ErrorHandler(w, http.StatusUnauthorized, "You can only view your profile", "Are you a stalker", nil)
		return
	}
	usr := struct {
		Username string
	}{
		Username: username,
	}

	ParseAndExecute(w, usr, "./static/templates/profile.html")
}

// Handle fetching user data
func ProfileInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	user, err := GetUser(r)
	if err != nil {
		JsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	var created time.Time
	DB.QueryRow(`SELECT created_at FROM users WHERE id = ?`, user.ID).Scan(&created)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"email":      user.Email,
		"username":   user.Username,
		"created_at": created,
	})
}

// Handler to Get the User's Own Posts
func UserPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	user, err := GetUser(r)
	if err != nil {
		JsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	// Read "offset" query param (default 0 if missing)
	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		JsonError(w, "Bad request", http.StatusBadRequest, err)
		return
	}

	rows, err := DB.Query(`
      	SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
      	FROM posts p
      	JOIN users u ON p.user_id = u.id
      	WHERE p.user_id = ?
      	ORDER BY p.id DESC
      	LIMIT ?
      	OFFSET ?
    `, user.ID, ProfileLimit, offset)
	if err != nil {
		JsonError(w, "Failed to get user's posts", http.StatusInternalServerError, err)
		return
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Username); err != nil {
			JsonError(w, "Failed to scan post", http.StatusInternalServerError, err)
			return
		}
		cats, _ := FetchPostCategories(p.ID)
		p.Categories = cats
		// If you store or join a profilePic somewhere, add it here
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}


// Handler to Get the User's *Liked* Posts
func LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	user, err := GetUser(r)
	if err != nil {
		JsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	// Read "offset" query param (default 0 if missing)
	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		JsonError(w, "Bad request", http.StatusBadRequest, err)
		return
	}

	rows, err := DB.Query(`
      	SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
      	FROM post_reactions pr
      	JOIN posts p ON pr.post_id = p.id
      	JOIN users u ON p.user_id = u.id
      	WHERE pr.user_id = ? AND pr.reaction_type = 'like'
      	ORDER BY p.id DESC
      	LIMIT ?
      	OFFSET ?
    `, user.ID, ProfileLimit, offset)
    if err != nil {
        JsonError(w, "Failed to get liked posts", http.StatusInternalServerError, err)
        return
    }

	defer rows.Close()

	var posts []Post
    for rows.Next() {
        var p Post
        if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Username); err != nil {
            JsonError(w, "Failed to scan liked posts", http.StatusInternalServerError, err)
            return
        }
		cats, _ := FetchPostCategories(p.ID)
		p.Categories = cats
        posts = append(posts, p)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var LikeMap = map[string]bool{
	"like":    true,
	"dislike": true,
}

// Handle adding a new post/comment reaction.
func AddReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	user, err := GetUser(r)
	if err != nil {
		JsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	// Decide if it’s a post or comment from query param
	typeParam := r.URL.Query().Get("type")
	var tableName, idColumn string
	switch typeParam {
	case "comment":
		tableName = "comment_reactions"
		idColumn = "comment_id"
	case "post":
		tableName = "post_reactions"
		idColumn = "post_id"
	}

	var payload struct {
		ID           int    `json:"id"`
		ReactionType string `json:"reaction_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		JsonError(w, "Inavlid payload", http.StatusBadRequest, err)
		return
	}

	if !LikeMap[payload.ReactionType] {
		JsonError(w, "Bad request", http.StatusBadRequest, nil)
		return
	}

	// Check if post or comment exists
	var exists bool
	table := typeParam + "s" // Ensures it becomes "posts" or "comments"
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = ?)`, table)
	err = DB.QueryRow(query, payload.ID).Scan(&exists)
	if err != nil {
		JsonError(w, "Failed to verify post or comment existence", http.StatusInternalServerError, err)
		return
	}
	if !exists {
		JsonError(w, "Post or comment does not exist", http.StatusBadRequest, nil)
		return
	}

	existingReaction := ""
	querySelect := `SELECT reaction_type FROM ` + tableName + ` WHERE user_id = ? AND ` + idColumn + ` = ?`
	DB.QueryRow(querySelect, user.ID, payload.ID).Scan(&existingReaction)

	if existingReaction != payload.ReactionType {
		// Upsert reaction (INSERT or UPDATE)
		queryUpsert := `
		 INSERT INTO ` + tableName + ` (` + idColumn + `, user_id, reaction_type)
		 VALUES (?, ?, ?)
		 ON CONFLICT(` + idColumn + `, user_id) DO UPDATE SET reaction_type = excluded.reaction_type
	 `
		_, err = DB.Exec(queryUpsert, payload.ID, user.ID, payload.ReactionType)
		if err != nil {
			JsonError(w, "Failed to add/update reaction", http.StatusInternalServerError, err)
			return
		}
	} else {
		// If the same, remove it => "un-toggle"
		queryDelete := `DELETE FROM ` + tableName + ` WHERE ` + idColumn + ` = ? AND user_id = ?`
		_, err = DB.Exec(queryDelete, payload.ID, user.ID)
		if err != nil {
			JsonError(w, "Failed to remove reaction", http.StatusInternalServerError, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reaction added/updated successfully"))
}

// Get a post/comments reactions.
func GetReactions(w http.ResponseWriter, r *http.Request) {
	typeParam := "post"
	ID := r.URL.Query().Get("post_id")
	if ID == "" {
		ID = r.URL.Query().Get("comment_id")
		typeParam = "comment"
	}

	if ID == "" {
		JsonError(w, "Missing id", http.StatusBadRequest, nil)
		return
	}

	var tableName, idColumn string
	switch typeParam {
	case "comment":
		tableName = "comment_reactions"
		idColumn = "comment_id"
	case "post":
		tableName = "post_reactions"
		idColumn = "post_id"
	}

	var likes, dislikes int

	err := DB.QueryRow(`
	SELECT
		COALESCE(SUM(CASE WHEN reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
	FROM `+tableName+`
	WHERE `+idColumn+` = ?`, ID).Scan(&likes, &dislikes)
	if err != nil {
		JsonError(w, "Failed to get reactions", http.StatusInternalServerError, err)
		return
	}

	// Fetch the current user's reaction if logged in
	userReaction := ""
	if user, err := GetUser(r); err == nil {
		_ = DB.QueryRow(`
            SELECT reaction_type
            FROM `+tableName+`
            WHERE `+idColumn+` = ? AND user_id = ?`,
			ID, user.ID).Scan(&userReaction)
	}

	response := map[string]interface{}{
		"likes":        likes,
		"dislikes":     dislikes,
		"userReaction": userReaction,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

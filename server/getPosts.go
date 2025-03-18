package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Num of posts on each scroll load in home.
var HomeLimit = 10

// Handle fetching all Posts with query offset and query tags.
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		JsonError(w, "Bad request", http.StatusBadRequest, err)
		return
	}

	tagsParam := r.URL.Query().Get("tags")

	var posts []Post

	if tagsParam == "" {
		// No filter => return all posts
		posts, err = FetchAllPosts(offset)
	} else {
		rawTags := strings.Split(tagsParam, ",")
		var tags []string
		for _, t := range rawTags {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, strings.ToLower(trimmed))
			}
		}
		if len(tags) == 0 {
			posts, err = FetchAllPosts(offset)
		} else {
			posts, err = FetchPostsByTags(offset, tags)
		}
	}

	if err != nil {
		JsonError(w, "Failed to query posts: "+err.Error(), http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Returns all posts (10 limit, offset)
func FetchAllPosts(offset int) ([]Post, error) {
	rows, err := DB.Query(`
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
        LIMIT ? OFFSET ?`, HomeLimit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanRows(rows)
}

// Return only posts that have all given tags
func FetchPostsByTags(offset int, tags []string) ([]Post, error) {
	placeholders := make([]string, len(tags))
	for i := range tags {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	args := make([]interface{}, 0, len(tags)+2)
	for _, t := range tags {
		args = append(args, t)
	}
	// Next param is the count for HAVING COUNT
	args = append(args, len(tags))
	// Append LIMIT (before-last param)
	args = append(args, HomeLimit)
	// Last param is offset
	args = append(args, offset)

	query := fmt.Sprintf(`
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
        FROM posts p
        JOIN users u ON p.user_id = u.id
        JOIN post_categories pc ON p.id = pc.post_id
        JOIN categories c ON pc.category_id = c.id
        WHERE LOWER(c.name) IN (%s)
        GROUP BY p.id
        HAVING COUNT(DISTINCT LOWER(c.name)) = ?
        ORDER BY p.created_at DESC
        LIMIT ? OFFSET ?`, inClause)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanRows(rows)
}

// Helper function to fetch categories for a post
func FetchPostCategories(postID int) ([]Category, error) {
	rows, err := DB.Query(`
        SELECT c.id, c.name
        FROM post_categories pc
        JOIN categories c ON pc.category_id = c.id
        WHERE pc.post_id = ?`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

// Helper function that scan rows and return a bunch of posts.
func ScanRows(rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var pa Post
		if err := rows.Scan(
			&pa.ID,
			&pa.UserID,
			&pa.Title,
			&pa.Content,
			&pa.CreatedAt,
			&pa.Username,
		); err != nil {
			return nil, err
		}
		cats, _ := FetchPostCategories(pa.ID)
		pa.Categories = cats
		posts = append(posts, pa)
	}
	return posts, rows.Err()
}

package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
)

const (
	// To avoid conflicts (post payload allowed in front and denied in back)
	// backend max should be bigger or equal than frontend max.
	maxTitleSize      = 700
	maxContentSize    = 10000
	maxCategoriesSize = 1000
)

// Handle adding a new post to DB.
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	safeTitle, safeContent, categories, quit := LimitRequestBody(w, r)
	if quit {
		return
	}

	user, err := GetUser(r)
	if err != nil {
		JsonError(w, "Login to add a post", http.StatusUnauthorized, err)
		return
	}

	// Insert post
	res, err := DB.Exec(`
	INSERT INTO posts (user_id, title, content)
	VALUES (?, ?, ?)`,
		user.ID, safeTitle, safeContent,
	)
	if err != nil {
		JsonError(w, "Failed to create post", http.StatusInternalServerError, err)
		return
	}
	postID, err := res.LastInsertId()
	if err != nil {
		JsonError(w, "Failed to retrieve post ID", http.StatusInternalServerError, err)
		return
	}

	// Insert categories in join table.
	if quit := InsertCategories(w, postID, categories); quit {
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post created successfully"))
}

// Check form values and Limit their readers one by one.
func LimitRequestBody(w http.ResponseWriter, r *http.Request) (string, string, []string, bool) {
	mr, err := r.MultipartReader()
	if err != nil {
		JsonError(w, "Invalid form values", http.StatusBadRequest, err)
		return "", "", nil, true
	}

	var titleB, contentB, catJson []byte
	var categories []string

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			JsonError(w, "Error reading form part", http.StatusInternalServerError, err)
			return "", "", nil, true
		}

		switch part.FormName() {
		case "title":
			titleB, err = LimitRead(part, maxTitleSize)
			if err != nil {
				JsonError(w, "Title is too big", http.StatusBadRequest, err)
				return "", "", nil, true
			}
		case "content":
			contentB, err = LimitRead(part, maxContentSize)
			if err != nil {
				JsonError(w, fmt.Sprintf("Content exceeded max length of %d characters", maxContentSize), http.StatusBadRequest, err)
				return "", "", nil, true
			}
		case "categories":
			catJson, err = LimitRead(part, maxCategoriesSize)
			if err != nil {
				JsonError(w, fmt.Sprintf("Categories exceed max length of %d", maxCategoriesSize), http.StatusBadRequest, err)
				return "", "", nil, true
			}
			if len(catJson) > 0 {
				err = json.Unmarshal([]byte(catJson), &categories)
				if err != nil {
					JsonError(w, "Invalid categories format", http.StatusBadRequest, err)
					return "", "", nil, true
				}
			}
		}
	}
	title, content := string(titleB), string(contentB)

	if title == "" || content == "" {
		JsonError(w, "Title and content are required", http.StatusBadRequest, nil)
		return "", "", nil, true
	}
	if len(title) < 4 {
		JsonError(w, "Title is too short", http.StatusBadRequest, nil)
		return "", "", nil, true
	}
	if len(content) < 6 {
		JsonError(w, "Post content is too short", http.StatusBadRequest, nil)
		return "", "", nil, true
	}

	// remove consecutive more than three consecutive new lines to just two.
	re := regexp.MustCompile(`(\r\n|\r|\n){3,}`)
	content = re.ReplaceAllString(content, "\n\n")

	return html.EscapeString(title), html.EscapeString(content), categories, false
}

// Read to a buffer and check if the maxSize limit is exeeded after the Read.
func LimitRead(part io.Reader, maxSize int) ([]byte, error) {
	buf := make([]byte, maxSize)
	n, err := part.Read(buf)
	if n == maxSize && err == nil {
		return nil, fmt.Errorf("input exceeds max allowed size of %d bytes", maxSize)
	}
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf[:n], nil
}

// Check if the categories in the post payload are present in categories Table.
// And Insert them into post_categories join table
func InsertCategories(w http.ResponseWriter, postID int64, categories []string) bool {
	if len(categories) > 3 {
		JsonError(w, "You can select only up to 3 categories", http.StatusBadRequest, nil)
		return true
	}

	for _, category := range categories {
		var categoryID int64

		// Check if category exists, and retrieve its ID
		err := DB.QueryRow(`
            SELECT id 
            FROM categories 
            WHERE name = ?
        `, category).Scan(&categoryID)

		// If the category does not exist
		if err == sql.ErrNoRows {
			JsonError(w, fmt.Sprintf("Category %s not found.", category), http.StatusBadRequest, err)
			return true
		} else if err != nil {
			JsonError(w, "Failed to find category", http.StatusInternalServerError, err)
			return true
		}

		// Insert into post_categories join table
		_, err = DB.Exec(`
            INSERT INTO post_categories (post_id, category_id) 
            VALUES (?, ?) 
        `, postID, categoryID)
		if err != nil {
			JsonError(w, "failed to link category.", http.StatusInternalServerError, err)
			return true
		}
	}
	return false
}

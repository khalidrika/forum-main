package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Return categories from DB.
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	rows, err := DB.Query(`SELECT name FROM categories`)
	if err != nil {
		JsonError(w, "Database error", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			JsonError(w, "Database error", http.StatusInternalServerError, err)
			return
		}
		categories = append(categories, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Add new categories 
// (No UI created for this you need to use a request using KinoTan_n token).
func AddCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JsonError(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	user, err := GetUser(r)
	if err != nil {
		JsonError(w, "No valid token found", http.StatusUnauthorized, err)
		return
	}

	if user.Username != "KinoTan_n" {
		JsonError(w, "Forbidden", http.StatusForbidden, nil)
		return
	}

	var reqBody struct {
		Categories []string `json:"categories"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		JsonError(w, "Invalid JSON body", http.StatusBadRequest, err)
		return
	}

	if len(reqBody.Categories) == 0 {
		JsonError(w, "No categories provided", http.StatusBadRequest, nil)
		return
	}

	for _, cat := range reqBody.Categories {
		if cat == "" {
			continue
		}

		_, err := DB.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, cat)
		if err != nil {
			JsonError(w, fmt.Sprintf("Failed inserting category '%s'", cat), http.StatusInternalServerError, err)
			return
		}

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"msg":"Categories added (duplicates ignored)","user":"%s"}`, user.Username)
}

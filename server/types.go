package server

import "time"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Post represents a forum post
type Post struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	Username   string     `json:"username"`
	Categories []Category `json:"categories,omitempty"`
}

// Post categories
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

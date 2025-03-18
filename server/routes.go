package server

import (
	"net/http"
	"time"
)

// Application routes.
func Routes() http.Handler {
	mux := http.NewServeMux()

	// Rate limiters, prevent spam and DoS attacks.
	// Allow 1 request per 20(x) microsecond
	rl := NewRateLimiter(20 * time.Microsecond)

	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/css/", FilesHandler)
	mux.HandleFunc("/js/", FilesHandler)
	mux.HandleFunc("/img/", FilesHandler)
	mux.HandleFunc("/cooldown", Cooldown)
	mux.HandleFunc("/post", ServePostHtml)
	mux.HandleFunc("/profile", ServeProfile)
	mux.HandleFunc("/api/check-session", CheckSession)
	mux.HandleFunc("/api/get-posts", GetPostsHandler)
	mux.HandleFunc("/api/get-categories", GetCategoriesHandler)
	mux.HandleFunc("/api/get-reactions", GetReactions)
	mux.HandleFunc("/api/get-singlePost", SinglePostHandler)
	mux.HandleFunc("/api/get-comments", GetComments)
	mux.HandleFunc("/api/comments-count", GetCommentsCount)
	mux.HandleFunc("/api/user-info", ProfileInfoHandler)
	mux.HandleFunc("/api/user-posts", UserPosts)
	mux.HandleFunc("/api/user-liked-posts", LikedPosts)

	// Rate Limiting (Auth & Content Creation)
	mux.Handle("/api/login", rl.Middleware(http.HandlerFunc(LoginHandler)))
	mux.Handle("/api/signup", rl.Middleware(http.HandlerFunc(SignUpHandler)))
	mux.Handle("/api/logout", rl.Middleware(http.HandlerFunc(LogoutHandler)))
	mux.Handle("/api/create-post", rl.Middleware(http.HandlerFunc(CreatePostHandler)))
	mux.Handle("/api/add-categories", rl.Middleware(http.HandlerFunc(AddCategoriesHandler)))
	mux.Handle("/api/add-reaction", rl.Middleware(http.HandlerFunc(AddReaction)))
	mux.Handle("/api/add-comment", rl.Middleware(http.HandlerFunc(AddComment)))

	return cspMiddleware(mux)
}

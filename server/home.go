package server

import "net/http"

// Handle index web page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, 404, "Looks like you're lost!", "The page you are looking for is not available!", nil)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405, http.StatusText(http.StatusMethodNotAllowed), "Only GET method is allowed!", nil)
		return
	}
	ParseAndExecute(w, "", "static/templates/home.html")
}

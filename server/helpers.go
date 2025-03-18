package server

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
)

// Parse the html files and execute them after checking for errors.
func ParseAndExecute(w http.ResponseWriter, data any, filename string) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		if strings.HasSuffix(filename, "error.html") {
			ServeCloudError(w, data.(ErrorData), err)
			return
		}
		ErrorHandler(w, http.StatusInternalServerError, "Something seems wrong, try again later!", "Internal Server Error!", err)
		return
	}

	// Write to a temporary buffer instead of writing directly to w.
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		if strings.HasSuffix(filename, "error.html") {
			ServeCloudError(w, data.(ErrorData), err)
			return
		}
		ErrorHandler(w, http.StatusInternalServerError, "Something seems wrong, try again later!", "Internal Server Error!", err)
		return
	}
	// If successful, write the buffer content to the ResponseWriter .
	buf.WriteTo(w)
}

func Cooldown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405, http.StatusText(http.StatusMethodNotAllowed), "Only GET method is allowed!", nil)
		return
	}
	ParseAndExecute(w, "", "static/templates/cooldown.html")
}

// Content-Security-Policy header, prevents XSS
func cspMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "script-src 'self';")
		next.ServeHTTP(w, r)
	})
}

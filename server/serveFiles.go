package server

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"time"
)

// Handle serving static content.
func FilesHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "static" + r.URL.Path

	filesBytes, err := os.ReadFile(filePath)

	// Prevent directory traversal attacks, ex: http://127.0.0.1:8080/css/..%2F..%2Fmain.go
	if err != nil || strings.Contains(filePath, "..") {
		ErrorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), "You don't have permission to access this link!", err)
		return
	}

	http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(filesBytes))
}

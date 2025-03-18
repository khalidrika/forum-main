package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ErrorData struct {
	Msg1       string
	Msg2       string
	StatusCode int
}

// Parse and execute error.html page depending on error type.
func ErrorHandler(w http.ResponseWriter, statusCode int, msg1, msg2 string, err error) {
	// Log the error in case of internal server error.
	if err != nil && statusCode == 500 {
		log.Println(err)
	}

	Error := ErrorData{
		Msg1:       msg1,
		Msg2:       msg2,
		StatusCode: statusCode,
	}

	w.WriteHeader(statusCode)
	ParseAndExecute(w, Error, "static/templates/error.html")
}

// Serve the error page from the json file error link.
func ServeCloudError(w http.ResponseWriter, error ErrorData, err error) {
	log.Println(err)
	errBody, err := GetErrorPage()
	if err != nil {
		http.Error(w, http.StatusText(error.StatusCode), error.StatusCode)
		log.Println(err)
		return
	}
	// Replace placeholders in the error page with dynamic messages
	errBody = strings.ReplaceAll(errBody, "{{.Msg1}}", error.Msg1)
	errBody = strings.ReplaceAll(errBody, "{{.Msg2}}", error.Msg2)
	errBody = strings.ReplaceAll(errBody, "{{.StatusCode}}", strconv.Itoa(error.StatusCode))
	if error.StatusCode == 500 {
		errBody = strings.ReplaceAll(errBody, `<a href="/"><button class="submit">Go to Home</button>`, "")
	}

	// Write the error page
	w.Write([]byte(errBody))
}

// GET the error page from Cloud (CloudLinks.ErrorPage).
func GetErrorPage() (string, error) {
	url := CloudLinks.ErrorPage

	// Make a GET request to the URL.
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch error page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch error page: received status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read error page content: %v", err)
	}
	return string(body), nil
}

// Return error response to JS fetches.
func JsonError(w http.ResponseWriter, msg string, statusCode int, err error) {
	// Log the error in case of internal server error.
	if err != nil && statusCode == 500 {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := struct {
		Msg string `json:"msg"`
	}{
		Msg: msg,
	}
	json.NewEncoder(w).Encode(resp)
}

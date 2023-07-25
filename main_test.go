package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check if the requested path is "/robots.txt"
	if r.URL.Path != "/robots.txt" {
		http.Error(w, "Invalid request path", http.StatusNotFound)
		return
	}

	// Write the robots.txt content to the response
	robotsContent := `User-agent: *
Disallow:`

	_, err := fmt.Fprint(w, robotsContent)
	if err != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
}

func TestRobotsEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/robots.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Pass a reference to the handler function you want to test
	// Assume that you have a function called RobotsHandler in your main package
	RobotsHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `User-agent: *
Disallow:`
	actual := strings.TrimSpace(rr.Body.String()) // Trimming to avoid issues with line breaks
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			actual, expected)
	}
}

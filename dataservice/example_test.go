package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHandleHealthcheck validates the health check endpoint.
func TestHandleHealthcheck(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleHealthcheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse response JSON: %v", err)
	}

	if val, ok := response["status"]; !ok || val != "OK" {
		t.Errorf("Expected response status OK, got %v", response["status"])
	}
}

// TestHandleRequest_BadRequest checks invalid method handling.
func TestHandleRequest_BadRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/posts", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleRequest(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest (400), got %v", rr.Code)
	}
}

// TestHandleCreatePost_InvalidJSON validates invalid JSON handling.
func TestHandleCreatePost_InvalidJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	invalidJSON := strings.NewReader("{invalid json}") // malformed JSON
	req, err := http.NewRequest("POST", "/posts", invalidJSON)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleCreatePost(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest (400) for invalid JSON, got %v", rr.Code)
	}
}

// TestHandleGetComments_InvalidPostID tests handling of invalid post IDs.
func TestHandleGetComments_InvalidPostID(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/posts/invalid/comments", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleGetComments(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest (400) for invalid post ID, got %v", rr.Code)
	}
}
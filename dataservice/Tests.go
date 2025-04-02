package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHealthcheck(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleHealthcheck(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", rr.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse response: %v", err)
	}

	if response["status"] != "OK" {
		t.Errorf("expected status OK, got %v", response["status"])
	}
}

func TestHandleRequest_BadRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("DELETE", "/posts", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleRequest(rr, r)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rr.Code)
	}
}

func TestHandleCreatePost_InvalidJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	invalidJSON := strings.NewReader("{invalid json}")
	r, err := http.NewRequest("POST", "/posts", invalidJSON)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleCreatePost(rr, r)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rr.Code)
	}
}

func TestHandleGetComments_InvalidPostID(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/posts/invalid/comments", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	HandleGetComments(rr, r)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest, got %v", rr.Code)
	}
}

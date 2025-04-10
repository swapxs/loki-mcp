package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := HealthHandler("1.0.0")

	// Call the handler directly and pass in our request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response HealthResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	// Check the response fields
	if response.Status != "ok" {
		t.Errorf("handler returned unexpected status: got %v want %v",
			response.Status, "ok")
	}

	if response.Version != "1.0.0" {
		t.Errorf("handler returned unexpected version: got %v want %v",
			response.Version, "1.0.0")
	}

	// Check that timestamp is recent (within the last minute)
	if time.Since(response.Timestamp) > time.Minute {
		t.Errorf("handler returned unexpected timestamp: %v is too old",
			response.Timestamp)
	}
}

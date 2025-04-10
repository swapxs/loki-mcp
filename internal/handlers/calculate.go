package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// CalculateRequest represents a calculation request
type CalculateRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// CalculateResponse represents a calculation response
type CalculateResponse struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Result    float64 `json:"result"`
}

// CalculateHandler returns a handler for calculation requests
func CalculateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body
		var req CalculateRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Perform the calculation
		var result float64
		var err error

		switch req.Operation {
		case "add":
			result = req.A + req.B
		case "subtract":
			result = req.A - req.B
		case "multiply":
			result = req.A * req.B
		case "divide":
			if req.B == 0 {
				http.Error(w, "Division by zero", http.StatusBadRequest)
				return
			}
			result = req.A / req.B
		default:
			http.Error(w, fmt.Sprintf("Unsupported operation: %s", req.Operation), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create the response
		response := CalculateResponse{
			Operation: req.Operation,
			A:         req.A,
			B:         req.B,
			Result:    result,
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(response)
	}
}

// CalculateQueryHandler handles calculation requests via query parameters
func CalculateQueryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		operation := r.URL.Query().Get("operation")
		aStr := r.URL.Query().Get("a")
		bStr := r.URL.Query().Get("b")

		// Validate parameters
		if operation == "" || aStr == "" || bStr == "" {
			http.Error(w, "Missing required parameters: operation, a, b", http.StatusBadRequest)
			return
		}

		// Parse numbers
		a, err := strconv.ParseFloat(aStr, 64)
		if err != nil {
			http.Error(w, "Invalid parameter 'a': must be a number", http.StatusBadRequest)
			return
		}

		b, err := strconv.ParseFloat(bStr, 64)
		if err != nil {
			http.Error(w, "Invalid parameter 'b': must be a number", http.StatusBadRequest)
			return
		}

		// Create request and use the same logic as the POST handler
		req := CalculateRequest{
			Operation: operation,
			A:         a,
			B:         b,
		}

		// Perform the calculation
		var result float64

		switch req.Operation {
		case "add":
			result = req.A + req.B
		case "subtract":
			result = req.A - req.B
		case "multiply":
			result = req.A * req.B
		case "divide":
			if req.B == 0 {
				http.Error(w, "Division by zero", http.StatusBadRequest)
				return
			}
			result = req.A / req.B
		default:
			http.Error(w, fmt.Sprintf("Unsupported operation: %s", req.Operation), http.StatusBadRequest)
			return
		}

		// Create the response
		response := CalculateResponse{
			Operation: req.Operation,
			A:         req.A,
			B:         req.B,
			Result:    result,
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(response)
	}
}

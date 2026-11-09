package handlers

import (
	"net/http"

	"github.com/goccy/go-json"
)

// APIResponse represents the standard JSON envelope for all API responses.
// It follows a consistent structure to simplify frontend parsing and error handling.
type APIResponse struct {
	Success bool        `json:"success"`         // Indicates if the request was successful
	Data    interface{} `json:"data,omitempty"`  // Holds the main response payload
	Error   *APIError   `json:"error,omitempty"` // Contains error details if success is false
	Meta    *Meta       `json:"meta,omitempty"`  // Holds pagination or auxiliary information
}

// APIError contains structured error information for the client.
type APIError struct {
	Code    string `json:"code"`              // Technical error identifier (e.g., "NOT_FOUND")
	Message string `json:"message"`           // Human-readable error message
	Details string `json:"details,omitempty"` // Detailed technical information for debugging
}

// Meta contains auxiliary metadata about the response, such as status codes or pagination.
type Meta struct {
	Code    string `json:"code"`              // Technical status code or identifier
	Message string `json:"message"`           // Contextual message about the operation
	Details string `json:"details,omitempty"` // Additional metadata details
}

// writeJSON is an internal helper that serializes the payload to JSON and
// sets the appropriate HTTP headers and status codes.
func writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// Success sends a 200 OK response with the provided data payload.
func Success(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta sends a 200 OK response including both a data payload and metadata.
func SuccessWithMeta(w http.ResponseWriter, data interface{}, meta *Meta) {
	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 Created response, typically used after a successful POST request.
func Created(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response, indicating a successful operation without a body.
func NoContent(w http.ResponseWriter) {
	writeJSON(w, http.StatusNoContent, APIResponse{
		Success: true,
	})
}

// Error sends a failure response with a specific HTTP status code, error code, and message.
func Error(w http.ResponseWriter, statusCode int, code, message string) {
	writeJSON(w, statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorWithDetails sends a failure response including extra technical details.
func ErrorWithDetails(w http.ResponseWriter, statusCode int, code, message, details string) {
	writeJSON(w, statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest sends a 400 Bad Request error.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

// NotFound sends a 404 Not Found error.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict sends a 409 Conflict error, often used for duplicate resources.
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "CONFLICT", message)
}

// InternalServerError sends a 500 Internal Server Error.
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

// Unauthorized sends a 401 Unauthorized error.
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden sends a 403 Forbidden error.
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, "FORBIDDEN", message)
}

// ValidationError sends a 400 Bad Request specifically for validation failures, including details.
func ValidationError(w http.ResponseWriter, details string) {
	ErrorWithDetails(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}

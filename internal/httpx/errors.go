package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	CodeInvalidID        Code = "invalid_id"
	CodeListingNotFound  Code = "listing_not_found"
	CodeInternalError    Code = "internal_error"
	CodeMalformedJson    Code = "malformed_json"
	CodeValidationFailed Code = "validation_failed"
	CodeUnauthenticated  Code = "unauthenticated"
	CodeForbidden        Code = "forbidden"
	CodeConflict         Code = "conflict"
	CodeRateLimited      Code = "rate_limited"
)

type errorPayload struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorPayload `json:"error"`
}

func Error(w http.ResponseWriter, status int, message string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorResponse{Error: errorPayload{Code: code, Message: message}})
}

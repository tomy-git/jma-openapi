// SPDX-FileCopyrightText: 2026 The jma-openapi contributors
//
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/tomy-git/jma-openapi/internal/gen"
)

type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Details    map[string]any
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(statusCode int, code, message string, details map[string]any, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
		Err:        err,
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil, err)
	}

	requestID := middleware.GetReqID(r.Context())
	if requestID == "" {
		requestID = "unknown"
	}

	payload := gen.ErrorEnvelope{
		Error: gen.ErrorModel{
			Code:      appErr.Code,
			Message:   appErr.Message,
			RequestID: requestID,
			Details:   appErr.Details,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

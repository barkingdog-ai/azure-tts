package model

import "fmt"

// APIError represents an error that occured on an API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

// APIErrorResponse is the full error respnose that has been returned by an API.
type APIErrorResponse struct {
	Error APIError `json:"error"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%d:%s] %s", e.StatusCode, e.Type, e.Message)
}

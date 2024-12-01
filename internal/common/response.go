package common

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// swagger:response Response
type Response struct {
	Message string `json:"message"`
}

type PaginationResponse[T any] struct {
	PaginationMetadata PaginationMetadata `json:"metadata"`
	Message            string             `json:"message"`
	Body               []T                `json:"body"`
}

// swagger:response BodyResponse
type BodyResponse struct {
	Message string      `json:"message"`
	Body    interface{} `json:"body,omitempty"`
}

// swagger:response ErrorResponse
type ErrorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func FormatErrorResponse(message string, error error) ErrorResponse {
	switch err := error.(type) {
	case validator.ValidationErrors:
		errorMessages := make([]string, 0)
		for _, e := range err {
			const errorFormat string = "%s: validation failed on the tag '%s'"
			errorMessages = append(errorMessages, fmt.Sprintf(errorFormat, e.Field(), e.Tag()))
		}

		return ErrorResponse{
			Message: message,
			Errors:  errorMessages,
		}
	default:
		return ErrorResponse{
			Message: message,
			Errors:  []string{error.Error()},
		}
	}
}

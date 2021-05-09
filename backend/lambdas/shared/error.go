package shared

import "fmt"

type APIError struct {
	StatusCode int
	Inner      error
	Message    string
}

func NewAPIError(message string, statusCode int) error {
	return &APIError{Message: message, StatusCode: statusCode}
}

func (e *APIError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s, %s, code: %d", e.Inner.Error(), e.Message, e.StatusCode)
	}
	return fmt.Sprintf("%s, status code: %d", e.Message, e.StatusCode)
}

func (e *APIError) Is(target error) bool {
	_, ok := target.(*APIError)
	return ok
}

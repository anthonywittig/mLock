package shared

import "fmt"

type APIError struct {
	ResponseCode int
	Inner        error
	Message      string
}

func NewAPIError(message string, responseCode int) error {
	return &APIError{Message: message, ResponseCode: responseCode}
}

func (e *APIError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s, %s, code: %d", e.Inner.Error(), e.Message, e.ResponseCode)
	}
	return fmt.Sprintf("%s, code: %d", e.Message, e.ResponseCode)
}

func (e *APIError) Is(target error) bool {
	_, ok := target.(*APIError)
	return ok
}

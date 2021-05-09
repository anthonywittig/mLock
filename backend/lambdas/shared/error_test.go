package shared

import (
	"errors"
	"testing"
)

func TestAPIError_Is_Simple(t *testing.T) {
	err := NewAPIError("hi", 400)
	if !errors.Is(err, &APIError{}) {
		t.Fatal("expected true but was false")
	}
}

func TestAPIError_As(t *testing.T) {
	err := NewAPIError("hi", 400)
	var apiErr *APIError
	ok := errors.As(err, &apiErr)

	if !ok {
		t.Fatal("expected ok but was !ok")
	}

	if apiErr.Message != "hi" || apiErr.StatusCode != 400 {
		t.Fatalf("unexpected: %+v", apiErr)
	}
}

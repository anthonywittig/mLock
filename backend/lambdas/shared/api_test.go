package shared

import (
	"testing"
)

func TestAPIResponse_AddCookie_Simple(t *testing.T) {
	resp, err := NewAPIResponse(200, "hi")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	resp.AddCookie("A", "1")

	val, exists := resp.Proxy.Headers[SetCookieHeaderName]
	if !exists {
		t.Fatal("cookie doesn't exist")
	}

	expectedVal := "A=1; Path=/; Max-Age=86400; HttpOnly; Secure; SameSite=Strict"
	if val != expectedVal {
		t.Fatalf("unexpected cookie value; expected \"%s\", was \"%s\"", expectedVal, val)
	}
}

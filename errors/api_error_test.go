package errors

import "testing"

func TestAPIErrorErrorIsNotEmpty(t *testing.T) {
	err := (&APIError{Code: 400, Description: "bad request"}).Error()
	if err == "" {
		t.Fatal("expected non-empty error string")
	}
}

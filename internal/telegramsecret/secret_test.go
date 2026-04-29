package telegramsecret

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	for _, token := range []string{"", "abc", "ABC_123-def", strings.Repeat("a", 256)} {
		if err := Validate(token); err != nil {
			t.Fatalf("Validate(%q) unexpected error: %v", token, err)
		}
	}

	for _, token := range []string{"bad token", "bad.token", strings.Repeat("a", 257)} {
		err := Validate(token)
		if err == nil {
			t.Fatalf("Validate(%q) expected error", token)
		}
		if strings.Contains(err.Error(), token) {
			t.Fatalf("error leaked secret: %q", err.Error())
		}
	}
}

func TestConstantTimeEqual(t *testing.T) {
	if !ConstantTimeEqual("secret", "secret") {
		t.Fatal("expected match")
	}
	if ConstantTimeEqual("secret", "wrong") {
		t.Fatal("expected mismatch")
	}
	if ConstantTimeEqual("secret", "secret-long") {
		t.Fatal("expected length mismatch")
	}
}

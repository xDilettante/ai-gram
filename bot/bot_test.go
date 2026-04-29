package bot

import "testing"

func TestNewRejectsEmptyToken(t *testing.T) {
	bot, err := New(BotConfig{})
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if bot != nil {
		t.Fatal("expected nil bot for empty token")
	}
}

func TestNewCreatesBotWithToken(t *testing.T) {
	bot, err := New(BotConfig{Token: "123:abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bot == nil {
		t.Fatal("expected bot")
	}
	if got := bot.Token(); got != "123:abc" {
		t.Fatalf("unexpected token: %q", got)
	}
}

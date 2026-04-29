// Package main prints a token-safe bot identity for smoke scripts.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	me, err := b.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("getMe: %w", err)
	}
	if me.Username == "" {
		fmt.Println("bot username: <empty>")
		return nil
	}
	fmt.Printf("bot username: @%s\n", me.Username)
	return nil
}

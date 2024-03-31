package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	ms "github.com/mmcnicol/message-store"
)

type MockableMessageStore interface {
	SaveEntry(topic string, entry ms.Entry) (int64, error)
	ReadEntry(topic string, offset int64) (*ms.Entry, error)
	PollForNextEntry(topic string, offset int64, pollDuration time.Duration) (*ms.Entry, error)
}

type MessageStoreClient struct {
	MessageStore MockableMessageStore
}

func main() {

	// Create a context that cancels when the application terminates
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a signal channel to capture termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	backend := NewBackend()

	// Start a goroutine to poll topic "system.audit.event"
	go backend.pollSystemAuditEvent(ctx)
	// Start a goroutine to poll topic "user.login.attempt"
	go backend.pollUserLoginAttempt(ctx)
	// Start a goroutine to poll topic "user.login.attempt.outcome"
	go backend.pollUserLoginAttemptOutcome(ctx)

	// Start a goroutine to generate user login attempts to topic "user.login.attempt"
	go backend.generateUserLoginAttempts(ctx)

	// Wait for termination signal
	<-sigCh
	fmt.Println("Termination signal received, canceling polling")
	cancel() // Cancel the context
}

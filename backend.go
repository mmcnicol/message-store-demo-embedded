package main

import (
	ms "github.com/mmcnicol/message-store"
)

// Backend represents a server side application
type Backend struct {
	MessageStore MockableMessageStore
}

// NewBackend creates a new instance of Backend
func NewBackend() *Backend {

	msgStore := ms.NewMessageStore()
	return &Backend{
		MessageStore: msgStore,
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ms "github.com/mmcnicol/message-store"
)

// Define the structure of SystemAuditEvent
type SystemAuditEvent struct {
	UserName          string `json:"userName"`
	SubjectIdentifier string `json:"subjectIdentifier"`
	AuditEvent        string `json:"auditEvent"`
}

// NewSystemAuditEvent creates a new instance of SystemAuditEvent
func NewSystemAuditEvent(userName, auditEvent string) *SystemAuditEvent {

	return &SystemAuditEvent{
		UserName:   userName,
		AuditEvent: auditEvent,
	}
}

// sendSystemAuditEvent sends a system audit event to a topic
func (b *Backend) sendSystemAuditEvent(systemAuditEvent SystemAuditEvent) {

	topic := SYSTEM_AUDIT_EVENT

	eventJSON, err := json.Marshal(systemAuditEvent)
	if err != nil {
		fmt.Printf("failed to marshal event: %v", err)
	}

	messageStoreEntry := ms.Entry{}
	messageStoreEntry.Key = []byte(systemAuditEvent.UserName)
	messageStoreEntry.Value = []byte(eventJSON)

	offset, err := b.MessageStore.SaveEntry(topic, messageStoreEntry)
	if err != nil {
		fmt.Printf("SaveEntry for topic '%s' failed: %v", topic, err)
	}
	fmt.Println("saved SystemAuditEvent to topic ", topic, " at offset ", offset)
}

// pollSystemAuditEvent polls the topic for system audit events
func (b *Backend) pollSystemAuditEvent(ctx context.Context) {

	topic := SYSTEM_AUDIT_EVENT
	offset := int64(-1)
	pollDuration := 100 * time.Millisecond

	// Loop until the context is canceled
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Polling canceled")
			return
		default:
			//entry, err := b.MessageStore.ReadEntry(topic, offset)
			entry, err := b.MessageStore.PollForNextEntry(topic, offset, pollDuration)
			if err != nil {
				fmt.Println("PollForNextEntry for topic '", topic, "', offset ", offset, ", returned error: ", err)
				// Wait for a short duration
				time.Sleep(500 * time.Millisecond)
			}
			if entry == nil {
				// do nothing - as this just means there are no unread entries in the topic

				//fmt.Println("PollForNextEntry for topic '", topic, "', offset ", offset, ", returned nil entry.Value")
				// Wait for a short duration
				time.Sleep(500 * time.Millisecond)
				continue // Continue to the next iteration of the loop
			}
			fmt.Println("PollForNextEntry for topic '", topic, "', offset ", offset, ", returned an entry")
			offset++
			b.processEntryFromPollSystemAuditEvent(*entry)
		}
	}
}

// processEntryFromPollSystemAuditEvent processes an entry from polling a topic
func (b *Backend) processEntryFromPollSystemAuditEvent(entry ms.Entry) {

	var systemAuditEvent SystemAuditEvent
	if err := json.Unmarshal(entry.Value, &systemAuditEvent); err != nil {
		fmt.Printf("failed to unmarshal systemAuditEvent: %v", err)
	}
	fmt.Printf("received systemAuditEvent: %v\n", systemAuditEvent)
}

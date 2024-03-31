package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	ms "github.com/mmcnicol/message-store"
)

// Define the structure of UserSubjectAccessAttempt
type UserSubjectAccessAttempt struct {
	UserName          string `json:"userName"`
	SubjectIdentifier string `json:"subjectIdentifier"`
}

// NewUserSubjectAccessAttempt creates a new instance of UserSubjectAccessAttempt
func NewUserSubjectAccessAttempt(userName, subjectIdentifier string) *UserSubjectAccessAttempt {

	return &UserSubjectAccessAttempt{
		UserName:          userName,
		SubjectIdentifier: subjectIdentifier,
	}
}

// generateUserSubjectAccessAttempt generates a user subject access attempt
func (b *Backend) generateUserSubjectAccessAttempt(userName string) {

	userSubjectAccessAttempt := NewUserSubjectAccessAttempt(userName, b.generateRandomSubjectIdentifier(10))
	b.sendUserSubjectAccessAttempt(*userSubjectAccessAttempt)
	systemAuditEvent := NewSystemAuditEventWithSubject(userSubjectAccessAttempt.UserName, userSubjectAccessAttempt.SubjectIdentifier, "login attempt")
	b.sendSystemAuditEvent(*systemAuditEvent)
}

// generateRandomSubjectIdentifier generates a random subject identifier consisting of digits 0 to 9
func (b *Backend) generateRandomSubjectIdentifier(subjectIdentifierLength int) string {

	// Define the character set
	charset := "0123456789"
	subjectIdentifier := make([]byte, subjectIdentifierLength)

	// Generate each character of the subject identifier randomly from the character set
	for i := 0; i < subjectIdentifierLength; i++ {
		randomIndex := rand.Intn(len(charset))
		subjectIdentifier[i] = charset[randomIndex]
	}

	return string(subjectIdentifier)
}

// sendUserSubjectAccessAttempt sends a user subject access attempt to a topic
func (b *Backend) sendUserSubjectAccessAttempt(userSubjectAccessAttempt UserSubjectAccessAttempt) {

	topic := USER_SUBJECT_ACCESS_ATTEMPT_TOPIC

	eventJSON, err := json.Marshal(userSubjectAccessAttempt)
	if err != nil {
		fmt.Printf("failed to marshal event: %v", err)
	}

	messageStoreEntry := ms.Entry{}
	messageStoreEntry.Key = []byte(userSubjectAccessAttempt.UserName)
	messageStoreEntry.Value = []byte(eventJSON)

	offset, err := b.MessageStore.SaveEntry(topic, messageStoreEntry)
	if err != nil {
		fmt.Printf("SaveEntry for topic '%s' failed: %v", topic, err)
	}
	fmt.Println("saved userSubjectAccessAttempt to topic ", topic, " at offset ", offset)
}

// pollUserSubjectAccessAttempt polls the topic for user subject access attempts
func (b *Backend) pollUserSubjectAccessAttempt(ctx context.Context) {

	topic := USER_SUBJECT_ACCESS_ATTEMPT_TOPIC
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
				// do nothing - as this just means there are no unread entries in the topic
				//fmt.Println("PollForNextEntry for topic '", topic, "', offset ", offset, ", returned error: ", err)
				// Wait for a short duration
				time.Sleep(500 * time.Millisecond)
			}
			if entry == nil {
				//fmt.Println("PollForNextEntry for topic '", topic, "', offset ", offset, ", returned nil entry.Value")
				// Wait for a short duration
				time.Sleep(500 * time.Millisecond)
				continue // Continue to the next iteration of the loop
			}
			fmt.Println("PollForNextEntry for topic '", topic, "', offset ", offset, ", returned an entry")
			offset++
			b.processEntryFromPollUserSubjectAccessAttempt(*entry)
		}
	}
}

// processEntryFromPollUserSubjectAccessAttempt processes an entry from polling a topic
func (b *Backend) processEntryFromPollUserSubjectAccessAttempt(entry ms.Entry) {

	var userSubjectAccessAttempt UserSubjectAccessAttempt
	if err := json.Unmarshal(entry.Value, &userSubjectAccessAttempt); err != nil {
		fmt.Printf("failed to unmarshal userSubjectAccessAttempt: %v", err)
	}
	fmt.Printf("Received userSubjectAccessAttempt: %v\n", userSubjectAccessAttempt)

	b.processUserSubjectAccessAttempt(userSubjectAccessAttempt)
}

// processUserSubjectAccessAttempt processes a user subject access attempt
func (b *Backend) processUserSubjectAccessAttempt(userSubjectAccessAttempt UserSubjectAccessAttempt) {

	userSubjectAccessAttemptOutcome := NewUserSubjectAccessAttemptOutcome(userSubjectAccessAttempt.UserName, userSubjectAccessAttempt.SubjectIdentifier, b.getUserSubjectAccessAttemptOutcome())
	b.sendUserSubjectAccessAttemptOutcome(*userSubjectAccessAttemptOutcome)

	var auditEvent string
	if userSubjectAccessAttemptOutcome.Outcome {
		auditEvent = "user subject access attempt successful"
	} else {
		auditEvent = "user subject access attempt failed"
	}
	systemAuditEvent := NewSystemAuditEventWithSubject(userSubjectAccessAttempt.UserName, userSubjectAccessAttempt.SubjectIdentifier, auditEvent)
	b.sendSystemAuditEvent(*systemAuditEvent)
}

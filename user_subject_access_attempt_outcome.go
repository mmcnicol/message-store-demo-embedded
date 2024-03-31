package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	ms "github.com/mmcnicol/message-store"
)

// Define the structure of UserSubjectAccessAttemptOutcome
type UserSubjectAccessAttemptOutcome struct {
	UserName          string `json:"userName"`
	SubjectIdentifier string `json:"subjectIdentifier"`
	Outcome           bool   `json:"outcome"`
}

// NewUserSubjectAccessAttemptOutcome creates a new instance of UserSubjectAccessAttemptOutcome
func NewUserSubjectAccessAttemptOutcome(userName string, subjectIdentifier string, outcome bool) *UserSubjectAccessAttemptOutcome {

	return &UserSubjectAccessAttemptOutcome{
		UserName:          userName,
		SubjectIdentifier: subjectIdentifier,
		Outcome:           outcome,
	}
}

// getUserSubjectAccessAttemptOutcome returns a user subject access attempt outcome
func (b *Backend) getUserSubjectAccessAttemptOutcome() bool {

	// Generate a random boolean value
	randomBool := rand.Intn(2) == 0 // Returns true with 50% probability, false with 50% probability
	return randomBool
}

// sendUserSubjectAccessAttemptOutcome sends a user subject access attempt outcome to a topic
func (b *Backend) sendUserSubjectAccessAttemptOutcome(userSubjectAccessAttemptOutcome UserSubjectAccessAttemptOutcome) {

	topic := USER_SUBJECT_ACCESS_ATTEMPT_TOPIC_OUTCOME

	eventJSON, err := json.Marshal(userSubjectAccessAttemptOutcome)
	if err != nil {
		fmt.Printf("failed to marshal event: %v", err)
	}

	messageStoreEntry := ms.Entry{}
	messageStoreEntry.Key = []byte(userSubjectAccessAttemptOutcome.UserName)
	messageStoreEntry.Value = []byte(eventJSON)

	offset, err := b.MessageStore.SaveEntry(topic, messageStoreEntry)
	if err != nil {
		fmt.Printf("SaveEntry for topic '%s' failed: %v", topic, err)
	}
	fmt.Println("saved UserSubjectAccessAttemptOutcome to topic ", topic, " at offset ", offset)
}

// pollUserSubjectAccessAttemptOutcome polls the topic for user subject access attempt outcomes
func (b *Backend) pollUserSubjectAccessAttemptOutcome(ctx context.Context) {

	topic := USER_SUBJECT_ACCESS_ATTEMPT_TOPIC_OUTCOME
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
			b.processEntryFromPollUserSubjectAccessAttemptOutcome(*entry)
		}
	}
}

// processEntryFromPollUserSubjectAccessAttemptOutcome processes an entry from polling a topic
func (b *Backend) processEntryFromPollUserSubjectAccessAttemptOutcome(entry ms.Entry) {

	var userSubjectAccessAttemptOutcome UserSubjectAccessAttemptOutcome
	if err := json.Unmarshal(entry.Value, &userSubjectAccessAttemptOutcome); err != nil {
		fmt.Printf("failed to unmarshal userSubjectAccessAttemptOutcome: %v", err)
	}
	fmt.Printf("received userSubjectAccessAttemptOutcome: %v\n", userSubjectAccessAttemptOutcome)
}

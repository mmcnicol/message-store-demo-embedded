package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	ms "github.com/mmcnicol/message-store"
)

// Define the structure of UserLoginAttemptOutcome
type UserLoginAttemptOutcome struct {
	UserName string `json:"userName"`
	Outcome  bool   `json:"outcome"`
}

// NewUserLoginAttemptOutcome creates a new instance of UserLoginAttemptOutcome
func NewUserLoginAttemptOutcome(userName string, outcome bool) *UserLoginAttemptOutcome {

	return &UserLoginAttemptOutcome{
		UserName: userName,
		Outcome:  outcome,
	}
}

// getUserLoginAttemptOutcome returns a user login attempt outcome
func (b *Backend) getUserLoginAttemptOutcome() bool {

	// Generate a random boolean value
	randomBool := rand.Intn(2) == 0 // Returns true with 50% probability, false with 50% probability
	return randomBool
}

// sendUserLoginAttemptOutcome sends a user login attempt outcome to a topic
func (b *Backend) sendUserLoginAttemptOutcome(userLoginAttemptOutcome UserLoginAttemptOutcome) {

	topic := USER_LOGIN_ATTEMPT_OUTCOME

	eventJSON, err := json.Marshal(userLoginAttemptOutcome)
	if err != nil {
		fmt.Printf("failed to marshal event: %v", err)
	}

	messageStoreEntry := ms.Entry{}
	messageStoreEntry.Key = []byte(userLoginAttemptOutcome.UserName)
	messageStoreEntry.Value = []byte(eventJSON)

	offset, err := b.MessageStore.SaveEntry(topic, messageStoreEntry)
	if err != nil {
		fmt.Printf("SaveEntry for topic '%s' failed: %v", topic, err)
	}
	fmt.Println("saved UserLoginAttemptOutcome to topic ", topic, " at offset ", offset)
}

// pollUserLoginAttemptOutcome polls the topic for user login attempt outcomes
func (b *Backend) pollUserLoginAttemptOutcome(ctx context.Context) {

	topic := USER_LOGIN_ATTEMPT_OUTCOME
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
			b.processEntryFromPollUserLoginAttemptOutcome(*entry)
		}
	}
}

// processEntryFromPollUserLoginAttemptOutcome processes an entry from polling a topic
func (b *Backend) processEntryFromPollUserLoginAttemptOutcome(entry ms.Entry) {

	var userLoginAttemptOutcome UserLoginAttemptOutcome
	if err := json.Unmarshal(entry.Value, &userLoginAttemptOutcome); err != nil {
		fmt.Printf("failed to unmarshal userLoginAttemptOutcome: %v", err)
	}
	fmt.Printf("received userLoginAttemptOutcome: %v\n", userLoginAttemptOutcome)
}

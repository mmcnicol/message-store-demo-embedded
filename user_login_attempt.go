package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	ms "github.com/mmcnicol/message-store"
)

// Define the structure of UserLoginAttempt
type UserLoginAttempt struct {
	UserName     string `json:"userName"`
	UserPassword string `json:"userPassword"`
}

// NewUserLoginAttempt creates a new instance of UserLoginAttempt
func NewUserLoginAttempt(userName, userPassword string) *UserLoginAttempt {

	return &UserLoginAttempt{
		UserName:     userName,
		UserPassword: userPassword,
	}
}

// generateUserLoginAttempts generates user login events
func (b *Backend) generateUserLoginAttempts(ctx context.Context) {

	// Loop until the context is canceled
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Polling canceled")
			return
		default:
			b.generateUserLoginAttempt()
			// Wait for a short duration
			time.Sleep(10 * time.Second)
		}
	}
}

// generateUserLoginAttempt generates a user login attempt
func (b *Backend) generateUserLoginAttempt() {

	userLoginAttempt := NewUserLoginAttempt(b.generateRandomUserName(), b.generateRandomPassword(8))
	b.sendUserLoginAttempt(*userLoginAttempt)
	systemAuditEvent := NewSystemAuditEvent(userLoginAttempt.UserName, "login attempt")
	b.sendSystemAuditEvent(*systemAuditEvent)
}

// generateRandomUserName generates a random userName
func (b *Backend) generateRandomUserName() string {

	// Lists of first names and surnames
	firstNames := []string{"Andrew", "David", "John", "Kate", "Brigitte", "Paula"}
	surnames := []string{"White", "Brown", "MacDonald", "Green", "Blue", "Pink"}

	// Randomly select a first name and a surname
	firstName := firstNames[rand.Intn(len(firstNames))]
	surname := surnames[rand.Intn(len(surnames))]

	// Concatenate the first character of the first name with the full surname
	userName := strings.ToLower(string(firstName[0])) + strings.ToLower(surname)

	return userName
}

// generateRandomPassword generates a random user password consisting of digits 0 to 9
func (b *Backend) generateRandomPassword(userPasswordLength int) string {

	// Define the character set
	charset := "0123456789"
	password := make([]byte, userPasswordLength)

	// Generate each character of the password randomly from the character set
	for i := 0; i < userPasswordLength; i++ {
		randomIndex := rand.Intn(len(charset))
		password[i] = charset[randomIndex]
	}

	return string(password)
}

// sendUserLoginAttempt sends a user login attempt to a topic
func (b *Backend) sendUserLoginAttempt(userLoginAttempt UserLoginAttempt) {

	topic := USER_LOGIN_ATTEMPT

	eventJSON, err := json.Marshal(userLoginAttempt)
	if err != nil {
		fmt.Printf("failed to marshal event: %v", err)
	}

	messageStoreEntry := ms.Entry{}
	messageStoreEntry.Key = []byte(userLoginAttempt.UserName)
	messageStoreEntry.Value = []byte(eventJSON)

	offset, err := b.MessageStore.SaveEntry(topic, messageStoreEntry)
	if err != nil {
		fmt.Printf("SaveEntry for topic '%s' failed: %v", topic, err)
	}
	fmt.Println("saved userLoginAttempt to topic ", topic, " at offset ", offset)
}

// pollUserLoginAttempt polls the topic for user login attempts
func (b *Backend) pollUserLoginAttempt(ctx context.Context) {

	topic := USER_LOGIN_ATTEMPT
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
			b.processEntryFromPollUserLoginAttempt(*entry)
		}
	}
}

// processEntryFromPollUserLoginAttempt processes an entry from polling a topic
func (b *Backend) processEntryFromPollUserLoginAttempt(entry ms.Entry) {

	var userLoginAttempt UserLoginAttempt
	if err := json.Unmarshal(entry.Value, &userLoginAttempt); err != nil {
		fmt.Printf("failed to unmarshal userLoginAttempt: %v", err)
	}
	fmt.Printf("Received userLoginAttempt: %v\n", userLoginAttempt)

	b.processUserLoginAttempt(userLoginAttempt)
}

// processUserLoginAttempt processes a user login attempt
func (b *Backend) processUserLoginAttempt(userLoginAttempt UserLoginAttempt) {

	userLoginAttemptOutcome := NewUserLoginAttemptOutcome(userLoginAttempt.UserName, b.getUserLoginAttemptOutcome())
	b.sendUserLoginAttemptOutcome(*userLoginAttemptOutcome)

	var auditEvent string
	if userLoginAttemptOutcome.Outcome {
		auditEvent = "login successful"
	} else {
		auditEvent = "login failed"
	}
	systemAuditEvent := NewSystemAuditEvent(userLoginAttempt.UserName, auditEvent)
	b.sendSystemAuditEvent(*systemAuditEvent)
}

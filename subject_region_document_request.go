package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	ms "github.com/mmcnicol/message-store"
)

// Define the structure of SubjectRegionDocumentRequest
type SubjectRegionDocumentRequest struct {
	SubjectIdentifier string `json:"subjectIdentifier"`
	Region            int    `json:"region"`
	UserName          string `json:"userName"`
}

// NewSubjectRegionDocumentRequest creates a new instance of SubjectRegionDocumentRequest
func NewSubjectRegionDocumentRequest(subjectIdentifier string, region int, userName string) *SubjectRegionDocumentRequest {

	return &SubjectRegionDocumentRequest{
		SubjectIdentifier: subjectIdentifier,
		Region:            region,
		UserName:          userName,
	}
}

// sendSubjectRegionDocumentRequest sends a subject region document request to a topic
func (b *Backend) sendSubjectRegionDocumentRequest(subjectRegionDocumentRequest SubjectRegionDocumentRequest) {

	topic := SUBJECT_REGION_DOCUMENT_REQUEST_TOPIC

	eventJSON, err := json.Marshal(subjectRegionDocumentRequest)
	if err != nil {
		fmt.Printf("failed to marshal event: %v", err)
	}

	messageStoreEntry := ms.Entry{}
	messageStoreEntry.Key = []byte(subjectRegionDocumentRequest.UserName)
	messageStoreEntry.Value = []byte(eventJSON)

	offset, err := b.MessageStore.SaveEntry(topic, messageStoreEntry)
	if err != nil {
		fmt.Printf("SaveEntry for topic '%s' failed: %v", topic, err)
	}
	fmt.Println("saved subjectRegionDocumentRequest to topic ", topic, " at offset ", offset)
}

// pollSubjectRegionDocumentRequest polls the topic for subject region document requests
func (b *Backend) pollSubjectRegionDocumentRequest(ctx context.Context) {

	topic := SUBJECT_REGION_DOCUMENT_REQUEST_TOPIC
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
			b.processEntryFromPollSubjectRegionDocumentRequest(*entry)
		}
	}
}

// processEntryFromPollSubjectRegionDocumentRequest processes an entry from polling a topic
func (b *Backend) processEntryFromPollSubjectRegionDocumentRequest(entry ms.Entry) {

	var subjectRegionDocumentRequest SubjectRegionDocumentRequest
	if err := json.Unmarshal(entry.Value, &subjectRegionDocumentRequest); err != nil {
		fmt.Printf("failed to unmarshal subjectRegionDocumentRequest: %v", err)
	}
	fmt.Printf("Received subjectRegionDocumentRequest: %v\n", subjectRegionDocumentRequest)

	b.processSubjectRegionDocumentRequest(subjectRegionDocumentRequest)
}

// processSubjectRegionDocumentRequest processes a subject region document request
func (b *Backend) processSubjectRegionDocumentRequest(subjectRegionDocumentRequest SubjectRegionDocumentRequest) {
	// Generate a random number between 0 and 99
	randomNumber := rand.Intn(100)

	// Determine if the response should be successful or an error
	var subjectRegionDocumentResponse *SubjectRegionDocumentResponse
	if randomNumber < 80 {
		// Response is successful
		documents := b.generateRandomSubjectRegionDocuments(subjectRegionDocumentRequest.Region)
		subjectRegionDocumentResponse = NewSubjectRegionDocumentResponse(subjectRegionDocumentRequest.SubjectIdentifier,
			documents,
			subjectRegionDocumentRequest.UserName,
			subjectRegionDocumentRequest.Region,
			nil)
	} else {
		// Response is an error
		err := generateRandomError()
		subjectRegionDocumentResponse = NewSubjectRegionDocumentResponse(subjectRegionDocumentRequest.SubjectIdentifier,
			nil,
			subjectRegionDocumentRequest.UserName,
			subjectRegionDocumentRequest.Region,
			err)
	}

	// Send the response
	b.sendSubjectRegionDocumentResponse(*subjectRegionDocumentResponse)
}

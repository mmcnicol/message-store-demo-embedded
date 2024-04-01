package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	ms "github.com/mmcnicol/message-store"
)

// Define the structure of SubjectRegionDocument
type SubjectRegionDocument struct {
	DocumentIdentifier    string    `json:"documentIdentifier"`
	DocumentDate          time.Time `json:"documentDate"`
	DocumentCategoryCode  string    `json:"documentCategoryCode"`
	DocumentCategory      string    `json:"documentCategory"`
	DocumentSpecialtyCode string    `json:"documentSpecialtyCode"`
	DocumentSpecialty     string    `json:"documentSpecialty"`
	Region                int       `json:"region"`
}

// Define the structure of SubjectRegionDocumentResponse
type SubjectRegionDocumentResponse struct {
	SubjectIdentifier string                  `json:"subjectIdentifier"`
	Documents         []SubjectRegionDocument `json:"documents"`
	UserName          string                  `json:"userName"`
	Region            int                     `json:"region"`
	Err               error                   `json:"error"`
}

// NewSubjectRegionDocumentResponse creates a new instance of SubjectRegionDocumentResponse
func NewSubjectRegionDocumentResponse(subjectIdentifier string, documents []SubjectRegionDocument, userName string, region int, err error) *SubjectRegionDocumentResponse {

	return &SubjectRegionDocumentResponse{
		SubjectIdentifier: subjectIdentifier,
		Documents:         documents,
		UserName:          userName,
		Region:            region,
		Err:               err,
	}
}

// generateRandomSubjectRegionDocuments generates a random number of SubjectRegionDocument structs
func (b *Backend) generateRandomSubjectRegionDocuments(region int) []SubjectRegionDocument {
	var documents []SubjectRegionDocument
	numDocuments := rand.Intn(6) // Generate a random number between 0 and 5
	for i := 0; i < numDocuments; i++ {
		documents = append(documents, SubjectRegionDocument{
			DocumentIdentifier:    b.generateRandomDocumentIdentifier(12),
			DocumentDate:          b.generateRandomDocumentDate(),
			DocumentCategoryCode:  "CategoryCode" + strconv.Itoa(i),
			DocumentCategory:      "Category" + strconv.Itoa(i),
			DocumentSpecialtyCode: "SpecialtyCode" + strconv.Itoa(i),
			DocumentSpecialty:     "Specialty" + strconv.Itoa(i),
			Region:                region,
		})
	}
	return documents
}

// generateRandomDocumentIdentifier generates a random document identifier consisting of digits 0 to 9
func (b *Backend) generateRandomDocumentIdentifier(documentIdentifierLength int) string {

	// Define the character set
	charset := "0123456789"
	documentIdentifier := make([]byte, documentIdentifierLength)

	// Generate each character of the document identifier randomly from the character set
	for i := 0; i < documentIdentifierLength; i++ {
		randomIndex := rand.Intn(len(charset))
		documentIdentifier[i] = charset[randomIndex]
	}

	return string(documentIdentifier)
}

// generateRandomDocumentDate generates a random document date
func (b *Backend) generateRandomDocumentDate() time.Time {
	min := time.Date(1990, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// generateRandomError generates a random error
func generateRandomError() error {

	// List of possible errors
	errors := []error{
		fmt.Errorf("http timeout"),
		fmt.Errorf("system unavailable"),
	}

	// Randomly select an error from the list
	return errors[rand.Intn(len(errors))]
}

// sendSubjectRegionDocumentResponse sends a user subject region document response to a topic
func (b *Backend) sendSubjectRegionDocumentResponse(subjectRegionDocumentResponse SubjectRegionDocumentResponse) {

	topic := SUBJECT_REGION_DOCUMENT_RESPONSE_TOPIC

	eventJSON, err := json.Marshal(subjectRegionDocumentResponse)
	if err != nil {
		fmt.Printf("failed to marshal event: %v", err)
	}

	messageStoreEntry := ms.Entry{}
	messageStoreEntry.Key = []byte(subjectRegionDocumentResponse.UserName)
	messageStoreEntry.Value = []byte(eventJSON)

	offset, err := b.MessageStore.SaveEntry(topic, messageStoreEntry)
	if err != nil {
		fmt.Printf("SaveEntry for topic '%s' failed: %v", topic, err)
	}
	fmt.Println("saved SubjectRegionDocumentResponse to topic ", topic, " at offset ", offset)
}

// pollSubjectRegionDocumentResponse polls the topic for subject region document responses
func (b *Backend) pollSubjectRegionDocumentResponse(ctx context.Context) {

	topic := SUBJECT_REGION_DOCUMENT_RESPONSE_TOPIC
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
				time.Sleep(100 * time.Millisecond)
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
			b.processEntryFromPollSubjectRegionDocumentResponse(*entry)
		}
	}
}

// processEntryFromPollSubjectRegionDocumentResponse processes an entry from polling a topic
func (b *Backend) processEntryFromPollSubjectRegionDocumentResponse(entry ms.Entry) {

	var subjectRegionDocumentResponse SubjectRegionDocumentResponse
	if err := json.Unmarshal(entry.Value, &subjectRegionDocumentResponse); err != nil {
		fmt.Printf("failed to unmarshal subjectRegionDocumentResponse: %v", err)
	}
	fmt.Printf("received subjectRegionDocumentResponse: %v\n", subjectRegionDocumentResponse)
}

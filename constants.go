package main

// Define constants for topic names
const (
	SYSTEM_AUDIT_EVENT_TOPIC                  = "system.audit.event"
	USER_LOGIN_ATTEMPT_TOPIC                  = "user.login.attempt"
	USER_LOGIN_ATTEMPT_OUTCOME_TOPIC          = "user.login.attempt.outcome"
	USER_SUBJECT_ACCESS_ATTEMPT_TOPIC         = "user.subject.access.attempt"
	USER_SUBJECT_ACCESS_ATTEMPT_OUTCOME_TOPIC = "user.subject.access.attempt.outcome"
	SUBJECT_REGION_DOCUMENT_REQUEST_TOPIC     = "subject.region.document.request"
	SUBJECT_REGION_DOCUMENT_RESPONSE_TOPIC    = "subject.region.document.response"
)

// Define constants for region names
const (
	AYRSHIRE_AND_ARRAN_REGION int = iota
	BORDERS_REGION
	DUMFRIES_AND_GALLOWAY_REGION
	FIFE_REGION
	FORTH_VALLEY_REGION
	GRAMPIAN_REGION
	GREATER_GLASGOW_AND_CYLDE_REGION
	HIGHLAND_REGION
	LOTHIAN_REGION
	LANARKSHIRE_REGION
	ORKNEY_REGION
	SHETLAND_REGION
	TAYSIDE_REGION
	WESTERN_ISLES_REGION
)

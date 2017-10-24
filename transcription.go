package messagebird

import "time"

// Transcription represents a Transcription object
type Transcription struct {
	ID          string
	RecordingID string
	Error       string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	Links       map[string]*string `json:"_links"`
}

// TranscriptionList represents a list of Transcriptions
type TranscriptionList struct {
	Data       []Transcription
	Links      map[string]*string `json:"_links"`
	Pagination map[string]int
}

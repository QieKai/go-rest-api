package messagebird

import "time"

// Recording represent a Messagebird recording object.
type Recording struct {
	ID        string
	Format    string
	LegID     string
	Status    string
	Duration  int
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Links     map[string]*string `json:"_links"`
}

// RecordingList represents a list of Recordings.
type RecordingList struct {
	Data       []Recording
	Links      map[string]*string `json:"_links"`
	Pagination map[string]int
}

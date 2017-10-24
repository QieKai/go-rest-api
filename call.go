package messagebird

import (
	"time"
)

// Call represents a voice call
type Call struct {
	ID          string
	NumberID    string
	Status      string
	Source      string
	Destination string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	EndedAt     *time.Time
	Links       map[string]*string `json:"_links"`
}

// CallList represents a list of voice calls.
type CallList struct {
	Data       []Call
	Links      map[string]*string `json:"_links"`
	Pagination map[string]int
}

// CallParams represent the parameters when creating a new voice call.
type CallParams struct {
	Source      string         `json:"source"`
	Destination string         `json:"destination"`
	CallFlow    CallFlowParams `json:"callFlow"`
}

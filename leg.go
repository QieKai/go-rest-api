package messagebird

import (
	"time"
)

// Leg represents a Messagebird Leg resource.
type Leg struct {
	ID          string
	CallID      string
	Source      string
	Destination string
	Status      string
	Direction   string
	Cost        float32
	Currency    string
	Duration    int
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	EndedAt     *time.Time
}

// LegList represent a list of Legs.
type LegList struct {
	Data       []Leg
	Links      map[string]*string `json:"_links"`
	Pagination map[string]int
}

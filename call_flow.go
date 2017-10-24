package messagebird

import (
	"time"
)

// CallFlow represents a CallFlow object.
type CallFlow struct {
	ID        string
	Title     string
	Steps     []Step
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// CallFlowParams contains the necessery parameters when creating a new CallFlow.
type CallFlowParams struct {
	Title string       `json:"title"`
	Steps []StepParams `json:"steps"`
}

// CallFlowList represents a list of CallFlows.
type CallFlowList struct {
	Links      map[string]string `json:"_links"`
	Pagination map[string]int
	Data       []CallFlow
}

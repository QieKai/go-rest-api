package messagebird

// Step represents a single CallFlow step.
type Step struct {
	ID      string
	Action  string
	Options map[string]string
}

// StepParams represents a single StepRequest step.
type StepParams struct {
	Action  string            `json:"action"`
	Options map[string]string `json:"options"`
}

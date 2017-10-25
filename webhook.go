package messagebird

import "time"

// Webhook represents a Webhook.
type Webhook struct {
	ID        string
	URL       string
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Links     map[string]*string `json:"_links"`
}

// WebhookList represents a list of Webhooks.
type WebhookList struct {
	Data       []Webhook
	Links      map[string]*string `json:"_links"`
	Pagination map[string]int
}

// WebhookParams represent the parameters used when creating a Webhook.
type WebhookParams struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

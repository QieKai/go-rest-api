package messagebird

import "time"

type Webhook struct {
	ID        string
	URL       string
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Links     map[string]*string `json:"_links"`
}

type WebhookList struct {
	Data       []Webhook
	Links      map[string]*string `json:"_links"`
	Pagination map[string]int
}

type WebhookParams struct {
	Title string `json:"_title"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

package client

import "strings"

const (
	WebhookPath    = "/repositories/%s/%s/hooks"
	WebhookPathGet = "/repositories/%s/%s/hooks/%s"
)

type Webhook struct {
	RepositoryId string   `json:"-"`
	Uuid         string   `json:"uuid,omitempty"`
	Url          string   `json:"url,omitempty"`
	Description  string   `json:"description,omitempty"`
	Events       []string `json:"events,omitempty"`
	Active       bool     `json:"active"`
	UseExisting  bool     `json:"-"`
}
type WebhookCollection struct {
	Values []Webhook `json:"values"`
}

func (wh *Webhook) WebhookEncodeId() string {
	return wh.RepositoryId + IdSeparator + wh.Uuid
}

func WebhookDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}

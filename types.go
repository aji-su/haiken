package main

type Message struct {
	Event   string   `json:"event"`
	Payload string   `json:"payload"`
	Stream  []string `json:"stream"`
}

type Status struct {
	ID          string                 `json:"id"`
	Account     *Account               `json:"account"`
	Content     string                 `json:"content"`
	Visibility  string                 `json:"visibility"`
	InReplyToID string                 `json:"in_reply_to_id"`
	Reblog      map[string]interface{} `json:"reblog"`
	SpoilerText string                 `json:"spoiler_text"`
	Tags        []*Tag                 `json:"tags"`
}

type Account struct {
	ID   string `json:"id"`
	Acct string `json:"acct"`
}

type Tag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Notification struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Account *Account `json:"account"`
}

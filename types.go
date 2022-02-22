package main

import "encoding/json"

type Message struct {
	Body struct {
		ID   string          `json:"id"`
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	} `json:"body"`
}

type Status struct {
	ID         string   `json:"id"`
	CreatedAt  string   `json:"createdAt"`
	Text       string   `json:"text"`
	Cw         *string  `json:"cw"`
	Account    Account  `json:"user"`
	UserID     string   `json:"userId"`
	Visibility string   `json:"visibility"` // enum: public home followers specified
	LocalOnly  bool     `json:"localOnly"`
	Renote     *Status  `json:"renote"`
	Mentions   []string `json:"mentions"`
	Tags       []string `json:"tags"`
}

type Account struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Host     *string `json:"host"`
}

type Tag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

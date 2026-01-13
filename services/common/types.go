package common

import "encoding/json"

type Project struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Event struct {
	Type string          `json:"type binding=required"`
	Data json.RawMessage `json:"data"`
}

type WelcomeData struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Username string `json:"username"`
}

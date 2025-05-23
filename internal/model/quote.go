package model

type Quote struct {
	ID     int    `json:"id,omitempty"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

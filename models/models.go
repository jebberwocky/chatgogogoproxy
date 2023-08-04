package models

type ChatRequest struct {
	Input string `json:"input"`
	ATag  string `json:"atag"`
}

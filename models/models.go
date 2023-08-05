package models

type ChatRequest struct {
	Input string `json:"input"`
	ATag  struct {
		Pk string `json:"pk"`
		Mk string `json:"mk"`
		Mh string `json:"mh"`
	} `json:"atag"`
}

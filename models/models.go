package models

type ChatRequest struct {
	Input string `json:"input"`
	ATag  struct {
		Pk string `json:"pk"`
		Mk string `json:"mk"`
		Mh string `json:"mh"`
	} `json:"atag"`
}

type AppContext struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Key      string `json:"key"`
	Path     string `json:"path"`
	Role     struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"role"`
}

type Apps struct {
	Apps []AppContext `json:"apps"`
}

package models

type ChatContent struct {
	Role    string `json:"role"`
	Content []struct {
		Type   string `json:"type"`
		Text   string `json:"text"`
		ImgUrl struct {
			Url string `json:"url"`
		} `json:"image_url"`
	} `json:"content"`
}

type ChatRequest struct {
	Input string `json:"input"`
	ATag  struct {
		Pk string `json:"pk"`
		Mk string `json:"mk"`
		Mh string `json:"mh"`
	} `json:"atag"`
	ChatContents []ChatContent `json:"contents"`
}

type AppContext struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Hostname     string `json:"hostname"`
	Key          string `json:"key"`
	Path         string `json:"path"`
	OpenaiApikey string `json:"openai_key"`
	OrgId        string `json:"org_id"`
	ImgProxyUrl  string `json:"img_proxy_url"`
	Role         struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"role"`
}

type Apps struct {
	Apps []AppContext `json:"apps"`
}

package chatgpt

import (
	"chatproxy/models"
	responses "chatproxy/response"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	apiKey          = "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"
	apiEndpoint_old = "https://api.openai.com/v1/completions"
	apiEndpoint     = "https://api.openai.com/v1/chat/completions"
	max_token       = 2048
	temperature     = 0.9
	top_p           = 1
	n               = 1
)

func generateResponse(d models.ChatRequest) (responses.Response, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").SetBody(
		map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []interface{}{map[string]interface{}{
				"role":    "system",
				"content": "Hi can you tell me what is the factorial of 10?"}},
			"max_tokens": 50,
		}).
		Post(apiEndpoint)

	if err != nil {
		return responses.Response{}, err
	} else {
		body := resp.Body()
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error while decoding JSON response:", err)
			return responses.Response{}, err
		}
		// Extract the content from the JSON response
		content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		return responses.Response{
			Message:         "success",
			ChatbotResponse: content,
		}, nil
	}

}

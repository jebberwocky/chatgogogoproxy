package chatgpt

import (
	"chatproxy/chatHistory"
	"chatproxy/models"
	responses "chatproxy/response"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	apiKey             = "***REMOVED***"
	apiEndpoint_legacy = "https://api.openai.com/v1/completions"
	apiEndpoint        = "https://api.openai.com/v1/chat/completions"
	max_token          = 2048
	temperature        = 0.9
	top_p              = 1
	n                  = 1
	Model_davinci      = "text-davinci-003"
	Model_v3           = "gpt-3.5-turbo"
	Model_v4           = "gpt-4"
)

func GenerateResponseLegacy(d models.ChatRequest) (responses.Response, error) {
	client := resty.New()
	prompt := d.Input
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").SetBody(
		map[string]interface{}{
			"model":       Model_davinci,
			"prompt":      prompt,
			"max_tokens":  max_token,
			"temperature": temperature,
			"top_p":       top_p,
			"n":           n,
		}).
		Post(apiEndpoint_legacy)

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

func GenerateResponse(d models.ChatRequest, model string) (responses.Response, error) {
	pk := d.ATag.Pk
	if model == Model_v4 {
		chatHistory.SetRule(pk, map[string]interface{}{
			"role":    "system",
			"content": "Hi can you tell me what is the factorial of 10?"})
	}
	messages := chatHistory.Set(pk, map[string]interface{}{
		"role":    "user",
		"content": d.Input})
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").SetBody(
		map[string]interface{}{
			"model":       model,
			"messages":    messages,
			"max_tokens":  max_token,
			"temperature": temperature,
			"top_p":       top_p,
			"n":           n,
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
		chatHistory.Set(pk, map[string]interface{}{
			"role":    "assistant",
			"content": content})
		return responses.Response{
			Message:         "assistant",
			ChatbotResponse: content,
		}, nil
	}
}

package chatgpt

import (
	"chatproxy/chatHistory"
	"chatproxy/middlewares"
	"chatproxy/models"
	responses "chatproxy/response"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	//apiKey             = "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"
	apiEndpoint_legacy  = "https://api.openai.com/v1/completions"
	apiEndpoint         = "https://api.openai.com/v1/chat/completions"
	max_token           = 2048
	temperature         = 0.9
	top_p               = 1
	n                   = 1
	Model_davinci       = "text-davinci-003"
	Model_v3            = "gpt-3.5-turbo"
	Model_v4            = "gpt-4"
	Model_self_training = "model_self_training"
)

func getSelfTrainModel(app models.AppContext) string {
	return "ft:gpt-3.5-turbo-0613:speaknoevil::7rMcCRYr"
}

func GenerateResponseLegacy(d models.ChatRequest, app models.AppContext) (responses.Response, error) {
	client := resty.New()
	client.OnBeforeRequest(middlewares.RestyOnBeforeRequest)
	client.OnAfterResponse(middlewares.RestyOnAfterResponse)
	prompt := d.Input
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+app.OpenaiApikey).
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
		fmt.Println(resp)
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error while decoding JSON response:", err)
			return responses.Response{}, err
		}
		// Extract the content from the JSON response
		if _, ok := data["error"]; ok {
			return responses.Response{
				Message:         "fail",
				ChatbotResponse: data["error"].(map[string]interface{})["message"].(string),
			}, nil
		}
		content := data["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
		return responses.Response{
			Message:         "success",
			ChatbotResponse: content,
		}, nil
	}
}

func GenerateResponse(d models.ChatRequest, app models.AppContext, model string) (responses.Response, error) {
	pk := d.ATag.Pk
	if model == Model_v4 {
		chatHistory.SetRule(pk, map[string]interface{}{
			"role":    app.Role.Role,
			"content": app.Role.Content})
	}
	if model == Model_self_training {
		chatHistory.SetRule(pk, map[string]interface{}{
			"role":    app.Role.Role,
			"content": app.Role.Content})
		model = getSelfTrainModel(app)
	}
	messages := chatHistory.Set(pk, map[string]interface{}{
		"role":    "user",
		"content": d.Input})
	client := resty.New()
	client.OnBeforeRequest(middlewares.RestyOnBeforeRequest)
	client.OnAfterResponse(middlewares.RestyOnAfterResponse)
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+app.OpenaiApikey).
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
		if _, ok := data["error"]; ok {
			return responses.Response{
				Message:         "fail",
				ChatbotResponse: data["error"].(map[string]interface{})["message"].(string),
			}, nil
		}
		// Extract the content from the JSON response
		content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		chatHistory.Set(pk, map[string]interface{}{
			"role":    "assistant",
			"content": content})
		return responses.Response{
			Message:         "success",
			ChatbotResponse: content,
		}, nil
	}
}

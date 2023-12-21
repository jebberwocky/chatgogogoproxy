package chatgpt

import (
	"chatproxy/chatHistory"
	"chatproxy/middlewares"
	"chatproxy/models"
	responses "chatproxy/response"
	"chatproxy/util"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"net/url"
)

const (
	//apiKey             = "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"
	apiEndpoint_legacy          = "https://api.openai.com/v1/completions"
	apiEndpoint                 = "https://api.openai.com/v1/chat/completions"
	apiEndpoint_Image           = "https://api.openai.com/v1/images/generations" //api.openai.com/v1/images/generations
	apiEndpoint_Image_Variation = "https://api.openai.com/v1/images/variations"
	max_token                   = 2048
	temperature                 = 0.9
	top_p                       = 1
	n                           = 1
	Model_davinci               = "text-davinci-003"
	Model_v3                    = "gpt-3.5-turbo"
	Model_v4                    = "gpt-4"
	Model_self_training         = "model_self_training"
	Model_dall_e_3              = "dall-e-3"
	Model_v4_vision             = "gpt-4-vision-preview"
	image_size                  = "1024x1024"
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

func GenerateImage(d models.ChatRequest, app models.AppContext, model string) (responses.Response, error) {
	client := resty.New()
	client.OnBeforeRequest(middlewares.RestyOnBeforeRequest)
	client.OnAfterResponse(middlewares.RestyOnAfterResponse)
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+app.OpenaiApikey).
		SetHeader("Content-Type", "application/json").SetBody(
		map[string]interface{}{
			"model":  model,
			"prompt": d.Input,
			"n":      n,
			"size":   image_size,
		}).
		Post(apiEndpoint_Image)

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
		_url := data["data"].([]interface{})[0].(map[string]interface{})["url"].(string)
		//revised_prompt := data["data"].([]interface{})[0].(map[string]interface{})["revised_prompt"].(string)
		content := fmt.Sprintf("<img src='%s'/>", _url)
		if len(app.ImgProxyUrl) > 0 {
			content = fmt.Sprintf("<img src='%s?url=%s&mh=%s'/>", app.ImgProxyUrl, url.QueryEscape(_url), d.ATag.Mh)
		}
		return responses.Response{
			Message:         "success",
			Data:            data,
			ChatbotResponse: content,
		}, nil
	}
}

func GenerateImageVariations(d models.ChatRequest, app models.AppContext) (responses.Response, error) {
	client := resty.New()
	client.OnBeforeRequest(middlewares.RestyOnBeforeRequest)
	client.OnAfterResponse(middlewares.RestyOnAfterResponse)
	var userImgDir = viper.GetString(util.UserUploadImgDir)
	var filename = d.Input
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+app.OpenaiApikey).
		SetFile("image", userImgDir+"/"+filename).
		Post(apiEndpoint_Image_Variation)

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
		var arr = data["data"].([]interface{})
		var resultcontent = ""
		for _, el := range arr {
			print(el)
			_url := el.(map[string]interface{})["url"].(string)
			if len(app.ImgProxyUrl) > 0 {
				resultcontent = fmt.Sprintf("<img src='%s?url=%s&mh=%s'/>", app.ImgProxyUrl, url.QueryEscape(_url), d.ATag.Mh)
			} else {
				resultcontent += fmt.Sprintf("<img src='%s'/>", _url)
			}
		}
		return responses.Response{
			Message:         "success",
			Data:            data,
			ChatbotResponse: resultcontent,
		}, nil
	}

}

func GenerateVision(d models.ChatRequest, app models.AppContext) (responses.Response, error) {
	client := resty.New()
	client.OnBeforeRequest(middlewares.RestyOnBeforeRequest)
	client.OnAfterResponse(middlewares.RestyOnAfterResponse)
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+app.OpenaiApikey).
		SetHeader("Content-Type", "application/json").SetBody(
		map[string]interface{}{
			"model": Model_v4_vision,
			"messages": []interface{}{
				map[string]interface{}{
					"role": "user",
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": util.FindInChatContent(d.ChatContents[0], "text"),
						},
						map[string]interface{}{
							"type": "image_url",
							"image_url": map[string]interface{}{
								"url": util.FindInChatContent(d.ChatContents[0], "image_url"),
							},
						},
					},
				},
			},
			"max_tokens": 300,
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
		return responses.Response{
			Message:         "success",
			ChatbotResponse: content,
		}, nil
	}
}

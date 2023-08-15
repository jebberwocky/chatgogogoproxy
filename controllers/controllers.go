package controllers

import (
	"chatproxy/chatgpt"
	"chatproxy/models"
	responses "chatproxy/response"
)

func DefaultHandle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	resp, err := chatgpt.GenerateResponseLegacy(d, a)
	return resp, err
}

func V3Handle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	resp, err := chatgpt.GenerateResponse(d, a, chatgpt.Model_v3)
	return resp, err
}

func V4Handle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	resp, err := chatgpt.GenerateResponse(d, a, chatgpt.Model_v4)
	return resp, err
}

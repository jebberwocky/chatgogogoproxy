package controllers

import (
	"chatproxy/chatgpt"
	"chatproxy/models"
	responses "chatproxy/response"
)

func DefaultHandle(d models.ChatRequest) (responses.Response, error) {
	resp, err := chatgpt.GenerateResponseLegacy(d)
	return resp, err
}

func V3Handle(d models.ChatRequest) (responses.Response, error) {
	resp, err := chatgpt.GenerateResponse(d, chatgpt.Model_v3)
	return resp, err
}

func V4Handle(d models.ChatRequest) (responses.Response, error) {
	resp, err := chatgpt.GenerateResponse(d, chatgpt.Model_v4)
	return resp, err
}

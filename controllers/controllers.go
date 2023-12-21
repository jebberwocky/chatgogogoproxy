package controllers

import (
	"chatproxy/chatgpt"
	"chatproxy/models"
	responses "chatproxy/response"
)

func DefaultHandle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	//resp, err := chatgpt.GenerateResponseLegacy(d, a)
	//08/25/2023 using finetune over td3
	resp, err := chatgpt.GenerateResponse(d, a, chatgpt.Model_self_training)
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

func VFineHandle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	resp, err := chatgpt.GenerateResponse(d, a, chatgpt.Model_self_training)
	return resp, err
}

func DalleHandle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	resp, err := chatgpt.GenerateImage(d, a, chatgpt.Model_dall_e_3)
	return resp, err
}

func VisionHandle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	resp, err := chatgpt.GenerateVision(d, a)
	return resp, err
}

func VariationHandle(d models.ChatRequest, a models.AppContext) (responses.Response, error) {
	resp, err := chatgpt.GenerateImageVariations(d, a)
	return resp, err
}

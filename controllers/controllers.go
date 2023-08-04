package controllers

import (
	"chatproxy/models"
	responses "chatproxy/response"
	"fmt"
)

func DefaultHandle(d models.ChatRequest) (responses.Response, error) {
	return responses.Response{}, nil
}

func V3Handle(d models.ChatRequest) (responses.Response, error) {
	return responses.Response{}, nil
}

func V4Handle(d models.ChatRequest) (responses.Response, error) {
	fmt.Println("v4")
	return responses.Response{}, nil
}

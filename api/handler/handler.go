package handler

import (
	"bot/service"
)

type HTTPHandler struct {
	service *service.Service
}

func NewHTTPHandler(service *service.Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

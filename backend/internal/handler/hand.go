package handler

import (
	"fmt"
	"stonkx/internal/domain"
)

type Handler struct {
	QuotsService domain.QuotsService
}

func New(s domain.QuotsService) *Handler {
	return &Handler{QuotsService: s}
}

func (h *Handler) Start() {
	fmt.Println("Lesssss goooooo")
}

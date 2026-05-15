package handler

import (
	"context"
	"fmt"
)

func (h *Handler) Initiate(ctx context.Context) error {
	if err := h.QuotsService.Serve(ctx); err != nil {
		return fmt.Errorf("BD not ready: %w", err)
	}
	if err := h.QuotsService.Initiate(ctx); err != nil {
		return fmt.Errorf("Can't initiate DB: %w", err)
	}
	return nil
}

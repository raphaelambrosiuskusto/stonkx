package service

import (
	"context"
	"fmt"
)

func (s *QuotsService) Serve(ctx context.Context) error {
	return s.Reader.Ready(ctx)
}

func (s *QuotsService) Initiate(ctx context.Context) error {
	if err := s.Reader.Initiate(ctx, ""); err != nil {
		return fmt.Errorf("creating table: %w", err)
	}

	if err := s.Writer.InsertMultipleSeries(ctx, GenerateDefaultMultiple(quotsInitialAmount)); err != nil {
		return fmt.Errorf("adding initial series to the table: %w", err)
	}

	return nil
}

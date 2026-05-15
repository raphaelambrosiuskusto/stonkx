package service

import (
	"context"
	"errors"
	"fmt"
	rand "math/rand/v2"
	"stonkx/internal/domain"
)

const newSeriesPerGenerationDefault = 10

type QuotsService struct {
	Reader domain.QuotsReader
	Writer domain.QuotsWriter
}

func NewQuots(reader domain.QuotsReader, writer domain.QuotsWriter) *QuotsService {
	return &QuotsService{
		Reader: reader,
		Writer: writer,
	}
}

func (s *QuotsService) GetByID(ctx context.Context, ID int) (domain.SeriesEx, error) {
	series, err := s.Reader.GetSeries(ctx, ID)
	if err != nil {
		return domain.SeriesEx{}, fmt.Errorf("service: GetByID: %w", err)
	}
	return series, nil
}

func (s *QuotsService) GetRandom(ctx context.Context) (domain.SeriesEx, error) {
	ids, err := s.Reader.ListTickerIDs(ctx)
	if err != nil {
		return domain.SeriesEx{}, fmt.Errorf("getrandom: %w", err)
	}

	series, err := s.Reader.GetSeries(ctx, ids[rand.IntN(len(ids))])
	if err != nil {
		return domain.SeriesEx{}, fmt.Errorf("reading series: %w", err)
	}
	return series, nil

}

func (s *QuotsService) Insert(ctx context.Context, chunk domain.SeriesEx) error {
	return s.Writer.InsertSeries(ctx, chunk)
}

func (s *QuotsService) InsertMultiple(ctx context.Context, chunks []domain.SeriesEx) error {
	return s.Writer.InsertMultipleSeries(ctx, chunks)
}

func (s *QuotsService) AddNew(ctx context.Context, amount int) error {
	switch {
	case amount == 1:
		return s.Writer.InsertSeries(ctx, GenerateDefault())
	case amount > 1:
		var chunks []domain.SeriesEx = make([]domain.SeriesEx, 0, amount)
		for range amount {
			chunks = append(chunks, GenerateDefault())
		}
		return s.Writer.InsertMultipleSeries(ctx, chunks)
	default:
		return errors.New("incorrect amount of series reqested for generation")
	}
}

func (s *QuotsService) AddNewDefault(ctx context.Context) error {
	return s.AddNew(ctx, newSeriesPerGenerationDefault)
}

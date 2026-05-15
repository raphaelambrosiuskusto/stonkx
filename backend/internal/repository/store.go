package repository

import (
	"context"
	"fmt"
	"stonkx/internal/domain"
	"stonkx/internal/shared"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{Pool: pool}
}

func (s *Store) GetSeries(ctx context.Context, tickerID int) (domain.SeriesEx, error) {
	var seriesEx domain.SeriesEx
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return seriesEx, fmt.Errorf("getseries: conn: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return seriesEx, fmt.Errorf("getseries: tx: %w", err)

	}

	rows, err := tx.Query(ctx, queryTicker, tickerID)
	if err != nil {
		return seriesEx, fmt.Errorf("getseries: query: %w", err)
	}
	defer rows.Close()

	var seriesIm SeriesIm
	for rows.Next() {
		var (
			time  int
			high  int
			low   int
			open  int
			close int
		)
		if err := rows.Scan(&time, &high, &low, &open, &close); err != nil {
			tx.Rollback(ctx)
			break
		}
		seriesIm.AddRow(time, high, low, open, close)
	}

	if err := tx.Commit(ctx); err != nil {
		return seriesEx, fmt.Errorf("getseries: commit: %w", err)
	}

	return seriesIm.toSeriesEx(), nil
}

func (s *Store) InsertSeries(ctx context.Context, chunk domain.SeriesEx) error {
	tickerID, err := s.GetMaxTickerID(ctx)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	defer tx.Rollback(ctx)

	for i, q := range chunk.Quots {
		_, err := tx.Exec(
			ctx,
			insertQuots,
			tickerID+1, q.Time, q.High, q.Low, q.Open, q.Close,
		)
		if err != nil {
			return fmt.Errorf("inserting a quot %v: %w", i, err)
		}
	}

	fmt.Printf("Inserted %d quotes for %d\n", len(chunk.Quots), tickerID)
	return tx.Commit(ctx)

}

func (s *Store) InsertMultipleSeries(ctx context.Context, chunks []domain.SeriesEx) error {
	tickerID, err := s.GetMaxTickerID(ctx)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	defer tx.Rollback(ctx)

	for j, chunk := range chunks {
		for i, q := range chunk.Quots {
			_, err := tx.Exec(
				ctx,
				insertQuots,
				tickerID+1+j, q.Time, q.High, q.Low, q.Open, q.Close,
			)
			if err != nil {
				return fmt.Errorf("inserting a quot %v of ticker %v: %w", i, tickerID+1+j, err)
			}
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) GetMaxTickerID(ctx context.Context) (int, error) {

	ans, err := s.Pool.Query(ctx, selectMaxTickerID)
	if err != nil {
		return -1, fmt.Errorf("selecting max ticker_id from db: %w", err)
	}
	var res int
	for ans.Next() {
		err := ans.Scan(&res)
		if err != nil {
			return -1, fmt.Errorf("reading response from db: %w", err)
		}
	}

	return res, nil
}

func (s *Store) WipeSeries(ctx context.Context, tickerID int) error {

	rows, err := s.Pool.Query(ctx, deleteQuots, tickerID)
	if err != nil {
		return fmt.Errorf("wiping ticker %v: %w", tickerID, err)
	}
	defer rows.Close()

	return nil
}

func (s *Store) WipeAllSeries(ctx context.Context) error {
	rows, err := s.Pool.Query(ctx, deleteAllQuots)
	if err != nil {
		return fmt.Errorf("truncating quots table: %w", err)
	}
	defer rows.Close()
	return nil
}

func (s *Store) ListTickerIDs(ctx context.Context) ([]int, error) {
	rows, err := s.Pool.Query(ctx, listTickersAsc)
	if err != nil {
		return nil, fmt.Errorf("listing tickerID's: %w", err)
	}
	defer rows.Close()

	ids := make([]string, 0, 16)
	for rows.Next() {
		var currentID string
		if err := rows.Scan(&currentID); err != nil {
			return nil, fmt.Errorf("reading ID from DB: %w", err)
		}
		ids = append(ids, currentID)
	}

	idlist, err := shared.Convert(ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids(string) to idlist(int): %w", err)
	}
	return idlist, nil
}

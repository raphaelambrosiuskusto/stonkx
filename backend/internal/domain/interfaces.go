package domain

import (
	"context"
)

type QuotsService interface {
	//checks if Service is ready to serve
	Serve(context.Context) error

	Initiate(ctx context.Context) error

	GetByID(context.Context, int) (SeriesEx, error)
	GetRandom(context.Context) (SeriesEx, error)
}

// DB-reading handler
type QuotsReader interface {
	//pings an underlying DB; is intended to return nil in case of successful connection
	Ready(context.Context) error
	CreateTableByFile(context.Context, string) error

	GetSeries(context.Context, int) (SeriesEx, error)
	GetMaxTickerID(context.Context) (int, error)
	ListTickerIDs(context.Context) ([]int, error)
}

// DB-writing handler
type QuotsWriter interface {
	InsertSeries(context.Context, SeriesEx) error
	InsertMultipleSeries(context.Context, []SeriesEx) error
	WipeSeries(context.Context, int) error
	WipeAllSeries(context.Context) error
}

// any type used for importing a row of Quot's (as defined in repository) from DB
type QuotsImport interface {
	AddRow(time, high, low, open, close int)
}

type QuotsExport interface {
	Values() SeriesEx
}

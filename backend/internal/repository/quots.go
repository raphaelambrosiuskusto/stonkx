package repository

import "stonkx/internal/domain"

type QuotIm struct {
	Time  int `json:"time"`
	High  int `json:"high"`
	Low   int `json:"low"`
	Open  int `json:"open"`
	Close int `json:"close"`
}

type SeriesIm struct {
	Quots []QuotIm
}

type IDList []int

func (s *SeriesIm) AddRow(time, high, low, open, close int) {
	s.Quots = append(s.Quots, QuotIm{time, high, low, open, close})
}

func (im SeriesIm) toSeriesEx() domain.SeriesEx {
	var ex domain.SeriesEx
	for _, q := range im.Quots {
		ex.Quots = append(ex.Quots, q.toQuotEx())
	}
	return ex
}

func (q QuotIm) toQuotEx() domain.QuotEx {
	return domain.QuotEx{
		Time:  q.Time,
		High:  q.High,
		Low:   q.Low,
		Open:  q.Open,
		Close: q.Close,
	}
}

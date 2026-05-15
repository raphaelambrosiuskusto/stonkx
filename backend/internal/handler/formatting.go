package handler

import "stonkx/internal/domain"

func ToRespFormat(s domain.SeriesEx) QuotResponse {
	var r QuotResponse
	for _, q := range s.Quots {
		r.Time = append(r.Time, q.Time)
		r.High = append(r.High, q.High)
		r.Low = append(r.Low, q.Low)
		r.Open = append(r.Open, q.Open)
		r.Close = append(r.Close, q.Close)
	}
	return r
}

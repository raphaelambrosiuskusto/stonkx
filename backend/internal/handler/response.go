package handler

type QuotResponse struct {
	Time  []int `json:"time"`
	High  []int `json:"high"`
	Low   []int `json:"low"`
	Open  []int `json:"open"`
	Close []int `json:"close"`
}

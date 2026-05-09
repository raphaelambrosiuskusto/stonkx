package domain

type QuotEx struct {
	Time  int
	High  int
	Low   int
	Open  int
	Close int
}

type SeriesEx struct {
	Quots []QuotEx
}

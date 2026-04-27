package handler

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}
type Quots struct {
	Time  []int `json:"time"`
	High  []int `json:"high"`
	Low   []int `json:"low"`
	Open  []int `json:"open"`
	Close []int `json:"close"`
}

func (h *Handler) GetData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	internalError := http.StatusInternalServerError

	pool := h.db
	conn, err := pool.Acquire(ctx)
	if err != nil {
		http.Error(w, err.Error(), internalError)
		return
	}

	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		http.Error(w, err.Error(), internalError)
		return
	}
	defer tx.Rollback(ctx)

	query := `select time, high, low, open, close from quotsint 
	where ticker_id = $1 and time < $2
	order by time ASC;`
	var ch chan int
	go func() {
		ch <- rand.IntN(10)
	}()
	rows, err := tx.Query(ctx, query, 1, 256)
	if err != nil {
		http.Error(w, err.Error(), internalError)
		return
	}
	defer rows.Close()

	q := Quots{}
	for rows.Next() {
		var time int
		var high int
		var low int
		var open int
		var close int

		err := rows.Scan(&time, &high, &low, &open, &close)

		if err != nil {
			tx.Rollback(ctx)
			break
		}

		q.AddRow(time, high, low, open, close)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(q)
	if err != nil {
		http.Error(w, err.Error(), internalError)
		return
	}
}

func New(pool *pgxpool.Pool) *Handler {
	return &Handler{
		db: pool,
	}
}

func (q *Quots) AddRow(time, high, low, open, close int) {
	q.Time = append(q.Time, time)
	q.High = append(q.High, high)
	q.Low = append(q.Low, low)
	q.Open = append(q.Open, open)
	q.Close = append(q.Close, close)
}

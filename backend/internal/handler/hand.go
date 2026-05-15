package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"stonkx/internal/domain"
)

type Handler struct {
	QuotsService domain.QuotsService
}

func New(s domain.QuotsService) *Handler {
	return &Handler{QuotsService: s}
}

func (h *Handler) GetQuots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := h.QuotsService.GetRandom(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ToRespFormat(resp))
}

func (h *Handler) Start() {
	mux := http.NewServeMux()
	h.registerRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("api/getcharts", h.GetQuots)
}

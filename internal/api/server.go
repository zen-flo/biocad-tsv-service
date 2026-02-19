package api

import (
	"biocad-tsv-service/internal/models"
	"biocad-tsv-service/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Server holds the dependencies for the API
type Server struct {
	MsgRepo *repository.MessageRepo
}

type MessageResponse struct {
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Total int              `json:"total"`
	Data  []models.Message `json:"data"`
}

// NewServer creates a new API server instance
func NewServer(msgRepo *repository.MessageRepo) *Server {
	return &Server{MsgRepo: msgRepo}
}

// Start starts the HTTP server on the given port
func (s *Server) Start(ctx context.Context, port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/messages", s.handleGetMessages)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("[api] starting server on %s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[api] server failed: %v", err)
		}
	}()

	// graceful shutdown
	go func() {
		<-ctx.Done()
		log.Println("[api] shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
}

// handleGetMessages handles GET /messages?unit_guid=...&page=...&limit=...
func (s *Server) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	unitGUIDStr := query.Get("unit_guid")
	if unitGUIDStr == "" {
		http.Error(w, "unit_guid is required", http.StatusBadRequest)
		return
	}

	unitGUID, err := uuid.Parse(unitGUIDStr)
	if err != nil {
		http.Error(w, "invalid unit_guid", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	total, err := s.MsgRepo.CountByUnitGUID(ctx, unitGUID)
	if err != nil {
		http.Error(w, "failed to count messages", http.StatusInternalServerError)
		return
	}

	messages, err := s.MsgRepo.GetByUnitGUIDPaginated(ctx, unitGUID, limit, offset)
	if err != nil {
		http.Error(w, "failed to query messages", http.StatusInternalServerError)
		return
	}

	resp := MessageResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Data:  messages,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

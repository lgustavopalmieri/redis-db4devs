package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type QueueItem struct {
	ID      string      `json:"id"`
	Payload interface{} `json:"payload"`
}

type QueueRepositoryInterface interface {
	Enqueue(ctx context.Context, key string, item QueueItem) error
	Dequeue(ctx context.Context, key string) (*QueueItem, error)
	QueueLength(ctx context.Context, key string) (int64, error)
}


type QueueHandler struct {
	Repo QueueRepositoryInterface
}

func NewQueueHandler(repo QueueRepositoryInterface) *QueueHandler {
	return &QueueHandler{Repo: repo}
}

func (h *QueueHandler) Enqueue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	var item QueueItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Enqueue(r.Context(), key, item); err != nil {
		http.Error(w, "Failed to enqueue item", http.StatusInternalServerError)
		log.Printf("Enqueue error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Item enqueued successfully"))
}

func (h *QueueHandler) Dequeue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	item, err := h.Repo.Dequeue(r.Context(), key)
	if err != nil {
		http.Error(w, "Failed to dequeue item", http.StatusInternalServerError)
		log.Printf("Dequeue error: %v", err)
		return
	}

	if item == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Demo struct {
	Topic   string                 `json:"topic"`
	Payload map[string]interface{} `json:"payload"`
}

type DemoRepositoryInterface interface {
	HashSave(ctx context.Context, demo *Demo) error
	HashGet(ctx context.Context, topic string) (*Demo, error)
	HashUpdate(ctx context.Context, demo *Demo) error
	HashDelete(ctx context.Context, topic string) error
}

type DemoHandler struct {
	DemoRepo DemoRepositoryInterface
}

func NewDemoHandler(repo DemoRepositoryInterface) *DemoHandler {
	return &DemoHandler{DemoRepo: repo}
}

func (h *DemoHandler) HashSave(w http.ResponseWriter, r *http.Request) {
	var demo Demo
	if err := json.NewDecoder(r.Body).Decode(&demo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.DemoRepo.HashSave(r.Context(), &demo); err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		log.Printf("Save error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Demo saved successfully"))
}

func (h *DemoHandler) HashGet(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "Missing topic parameter", http.StatusBadRequest)
		return
	}

	demo, err := h.DemoRepo.HashGet(r.Context(), topic)
	if err != nil {
		http.Error(w, "Failed to get data", http.StatusInternalServerError)
		log.Printf("Get error: %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(demo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *DemoHandler) HashUpdate(w http.ResponseWriter, r *http.Request) {
	var demo Demo
	if err := json.NewDecoder(r.Body).Decode(&demo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.DemoRepo.HashUpdate(r.Context(), &demo); err != nil {
		http.Error(w, "Failed to update data", http.StatusInternalServerError)
		log.Printf("Update error: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DemoHandler) HashDelete(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "Missing topic parameter", http.StatusBadRequest)
		return
	}

	if err := h.DemoRepo.HashDelete(r.Context(), topic); err != nil {
		http.Error(w, "Failed to delete data", http.StatusInternalServerError)
		log.Printf("Delete error: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

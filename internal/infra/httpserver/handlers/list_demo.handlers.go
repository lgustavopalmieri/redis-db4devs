package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ListHandler struct {
	Repo ListRepositoryInterface
}

// Define um objeto genérico a ser armazenado na lista
type ListItem struct {
	ID    string                 `json:"id"`
	Value map[string]interface{} `json:"value"`
}

type ListRepositoryInterface interface {
	ListPush(ctx context.Context, key string, items []ListItem) error
	ListPop(ctx context.Context, key string, fromStart bool) (*ListItem, error)
	ListGetAll(ctx context.Context, key string) ([]ListItem, error)
	ListRemove(ctx context.Context, key string, itemID string, count int) error
	ListLength(ctx context.Context, key string) (int64, error)
}

func NewListHandler(repo ListRepositoryInterface) *ListHandler {
	return &ListHandler{Repo: repo}
}

// Adiciona objetos à lista
func (h *ListHandler) Push(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Key   string     `json:"key"`
		Items []ListItem `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Repo.ListPush(r.Context(), data.Key, data.Items); err != nil {
		http.Error(w, "Failed to push items", http.StatusInternalServerError)
		log.Printf("Push error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Items pushed successfully"))
}

// Remove e retorna o primeiro ou último objeto da lista
func (h *ListHandler) Pop(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	fromStart := r.URL.Query().Get("fromStart") == "true"

	item, err := h.Repo.ListPop(r.Context(), key, fromStart)
	if err != nil {
		http.Error(w, "Failed to pop item", http.StatusInternalServerError)
		log.Printf("Pop error: %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Retorna todos os objetos da lista
func (h *ListHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	items, err := h.Repo.ListGetAll(r.Context(), key)
	if err != nil {
		http.Error(w, "Failed to get list", http.StatusInternalServerError)
		log.Printf("GetAll error: %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Remove um objeto específico da lista com base no ID
func (h *ListHandler) Remove(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Key    string `json:"key"`
		ItemID string `json:"item_id"`
		Count  int    `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Repo.ListRemove(r.Context(), data.Key, data.ItemID, data.Count); err != nil {
		http.Error(w, "Failed to remove item", http.StatusInternalServerError)
		log.Printf("Remove error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item removed successfully"))
}

// Retorna o tamanho da lista
func (h *ListHandler) GetLength(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	length, err := h.Repo.ListLength(r.Context(), key)
	if err != nil {
		http.Error(w, "Failed to get list length", http.StatusInternalServerError)
		log.Printf("GetLength error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"length": length})
}

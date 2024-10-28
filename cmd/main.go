package main

import (
	"context"
	"log"
	"net/http"

	"github.com/lgustavopalmieri/redis-db4devs/internal/infra/database/redis_repository"
	"github.com/lgustavopalmieri/redis-db4devs/internal/infra/httpserver/handlers"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	redisClient := NewRedisClient()
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	repo := redis_repository.NewRedisRepository(redisClient)
	handler := handlers.NewDemoHandler(repo)

	http.HandleFunc("/hash-save", handler.HashSave)
	http.HandleFunc("/hash-get", handler.HashGet)
	http.HandleFunc("/hash-update", handler.HashUpdate)
	http.HandleFunc("/hash-delete", handler.HashDelete)

	log.Println("Server running on :6000")
	if err := http.ListenAndServe(":6000", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

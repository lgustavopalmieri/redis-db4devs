package redis_repository

import (
	"context"
	"encoding/json"

	"github.com/lgustavopalmieri/redis-db4devs/internal/infra/httpserver/handlers"
	"github.com/redis/go-redis/v9"
)

type RedisQueueRepository struct {
	client *redis.Client
}

func NewRedisQueueRepository(client *redis.Client) *RedisQueueRepository {
	return &RedisQueueRepository{client: client}
}

func (r *RedisQueueRepository) Enqueue(ctx context.Context, key string, item handlers.QueueItem) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return r.client.RPush(ctx, key, data).Err() // Adiciona ao final da lista (fila)
}

func (r *RedisQueueRepository) Dequeue(ctx context.Context, key string) (*handlers.QueueItem, error) {
	data, err := r.client.LPop(ctx, key).Result() // Remove do início da lista (fila)
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Lista vazia
		}
		return nil, err
	}

	var item handlers.QueueItem
	if err := json.Unmarshal([]byte(data), &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *RedisQueueRepository) QueueLength(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result() // Tamanho da fila
}

package redis_repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgustavopalmieri/redis-db4devs/internal/infra/httpserver/handlers"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) HashSave(ctx context.Context, demo *handlers.Demo) error {
	key := fmt.Sprintf("demo:%s", demo.Topic)
	payloadJSON, err := json.Marshal(demo.Payload)
	if err != nil {
		return err
	}

	return r.client.HSet(ctx, key, map[string]interface{}{
		"topic":   demo.Topic,
		"payload": payloadJSON,
	}).Err()
}

func (r *RedisRepository) HashGet(ctx context.Context, topic string) (*handlers.Demo, error) {
	key := fmt.Sprintf("demo:%s", topic)

	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no data found for topic: %s", topic)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(data["payload"]), &payload); err != nil {
		return nil, err
	}

	return &handlers.Demo{
		Topic:   data["topic"],
		Payload: payload,
	}, nil
}

func (r *RedisRepository) HashUpdate(ctx context.Context, demo *handlers.Demo) error {
	key := fmt.Sprintf("demo:%s", demo.Topic)
	payloadJSON, err := json.Marshal(demo.Payload)
	if err != nil {
		return err
	}

	return r.client.HSet(ctx, key, map[string]interface{}{
		"topic":   demo.Topic,
		"payload": payloadJSON,
	}).Err()
}

func (r *RedisRepository) HashDelete(ctx context.Context, topic string) error {
	key := fmt.Sprintf("demo:%s", topic)
	return r.client.Del(ctx, key).Err()
}

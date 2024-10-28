package redis_repository

import (
	"context"
	"encoding/json"

	"github.com/lgustavopalmieri/redis-db4devs/internal/infra/httpserver/handlers"
)

// Adiciona objetos à lista como JSON
func (r *RedisRepository) ListPush(ctx context.Context, key string, items []handlers.ListItem) error {
	for _, item := range items {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		if err := r.client.LPush(ctx, key, data).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Remove e retorna o primeiro ou último objeto da lista
func (r *RedisRepository) ListPop(ctx context.Context, key string, fromStart bool) (*handlers.ListItem, error) {
	var result string
	var err error

	if fromStart {
		result, err = r.client.LPop(ctx, key).Result()
	} else {
		result, err = r.client.RPop(ctx, key).Result()
	}
	if err != nil {
		return nil, err
	}

	var item handlers.ListItem
	if err := json.Unmarshal([]byte(result), &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// Retorna todos os objetos da lista como uma slice de `ListItem`
func (r *RedisRepository) ListGetAll(ctx context.Context, key string) ([]handlers.ListItem, error) {
	results, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var items []handlers.ListItem
	for _, result := range results {
		var item handlers.ListItem
		if err := json.Unmarshal([]byte(result), &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Remove objetos específicos pelo ID
func (r *RedisRepository) ListRemove(ctx context.Context, key string, itemID string, count int) error {
	// Procura o item pelo ID e remove `count` vezes
	itemToRemove := handlers.ListItem{ID: itemID}
	data, err := json.Marshal(itemToRemove)
	if err != nil {
		return err
	}

	return r.client.LRem(ctx, key, int64(count), string(data)).Err()
}

// Retorna o tamanho total da lista
func (r *RedisRepository) ListLength(ctx context.Context, key string) (int64, error) {
	length, err := r.client.LLen(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return length, nil
}

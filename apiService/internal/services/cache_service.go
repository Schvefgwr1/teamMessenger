package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheService struct {
	redis *redis.Client
}

func NewCacheService(redisClient *redis.Client) *CacheService {
	return &CacheService{redis: redisClient}
}

// Set сохраняет данные в кеш
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	return c.redis.Set(ctx, key, data, ttl).Err()
}

// Get получает данные из кеша
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("cache miss")
		}
		return fmt.Errorf("failed to get cache data: %w", err)
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete удаляет данные из кеша
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}

// DeleteByPattern удаляет все ключи по паттерну
func (c *CacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.redis.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.redis.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists проверяет существование ключа
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.redis.Exists(ctx, key).Result()
	return result > 0, err
}

// SetTTL устанавливает TTL для существующего ключа
func (c *CacheService) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return c.redis.Expire(ctx, key, ttl).Err()
}

// GetTTL получает оставшееся время жизни ключа
func (c *CacheService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return c.redis.TTL(ctx, key).Result()
}

// Константы для типов кеша
const (
	UserCachePrefix     = "user:"
	MessageCachePrefix  = "messages:"
	ChatListCachePrefix = "chat_list:"
)

// Генераторы ключей для разных типов данных

func (c *CacheService) UserCacheKey(userID string) string {
	return fmt.Sprintf("%s%s", UserCachePrefix, userID)
}

func (c *CacheService) ChatMessagesCacheKey(chatID string) string {
	return fmt.Sprintf("%s%s", MessageCachePrefix, chatID)
}

func (c *CacheService) UserChatListCacheKey(userID string) string {
	return fmt.Sprintf("%s%s", ChatListCachePrefix, userID)
}

// Специализированные методы для пользователей

func (c *CacheService) SetUserCache(ctx context.Context, userID string, userData interface{}) error {
	key := c.UserCacheKey(userID)
	return c.Set(ctx, key, userData, 30*time.Minute) // 30 минут TTL для пользователей
}

func (c *CacheService) GetUserCache(ctx context.Context, userID string, dest interface{}) error {
	key := c.UserCacheKey(userID)
	return c.Get(ctx, key, dest)
}

func (c *CacheService) DeleteUserCache(ctx context.Context, userID string) error {
	key := c.UserCacheKey(userID)
	return c.Delete(ctx, key)
}

// Специализированные методы для сообщений чата

func (c *CacheService) SetChatMessagesCache(ctx context.Context, chatID string, messages interface{}) error {
	key := c.ChatMessagesCacheKey(chatID)
	return c.Set(ctx, key, messages, 10*time.Minute) // 10 минут TTL для сообщений
}

func (c *CacheService) GetChatMessagesCache(ctx context.Context, chatID string, dest interface{}) error {
	key := c.ChatMessagesCacheKey(chatID)
	return c.Get(ctx, key, dest)
}

func (c *CacheService) DeleteChatMessagesCache(ctx context.Context, chatID string) error {
	key := c.ChatMessagesCacheKey(chatID)
	return c.Delete(ctx, key)
}

// Специализированные методы для списка чатов пользователя

func (c *CacheService) SetUserChatListCache(ctx context.Context, userID string, chats interface{}) error {
	key := c.UserChatListCacheKey(userID)
	return c.Set(ctx, key, chats, 15*time.Minute) // 15 минут TTL для списка чатов
}

func (c *CacheService) GetUserChatListCache(ctx context.Context, userID string, dest interface{}) error {
	key := c.UserChatListCacheKey(userID)
	return c.Get(ctx, key, dest)
}

func (c *CacheService) DeleteUserChatListCache(ctx context.Context, userID string) error {
	key := c.UserChatListCacheKey(userID)
	return c.Delete(ctx, key)
}

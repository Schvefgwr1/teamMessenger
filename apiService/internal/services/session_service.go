package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type SessionStatus string

const (
	SessionActive  SessionStatus = "active"
	SessionRevoked SessionStatus = "revoked"
	SessionExpired SessionStatus = "expired"
)

type JWTSession struct {
	UserID    uuid.UUID     `json:"user_id"`
	Status    SessionStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	ExpiresAt time.Time     `json:"expires_at"`
}

type SessionService struct {
	redis *redis.Client
}

func NewSessionService(redisClient *redis.Client) *SessionService {
	return &SessionService{redis: redisClient}
}

// CreateSession создает новую сессию в Redis
func (s *SessionService) CreateSession(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	session := JWTSession{
		UserID:    userID,
		Status:    SessionActive,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Используем составной ключ: session:userID:tokenHash
	key := fmt.Sprintf("session:%s:%s", userID.String(), s.hashToken(token))

	// TTL равен времени истечения токена
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = 24 * time.Hour // Fallback на 24 часа
	}

	return s.redis.Set(ctx, key, sessionData, ttl).Err()
}

// GetSession получает сессию из Redis
func (s *SessionService) GetSession(ctx context.Context, userID uuid.UUID, token string) (*JWTSession, error) {
	key := fmt.Sprintf("session:%s:%s", userID.String(), s.hashToken(token))

	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session JWTSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// RevokeSession отзывает сессию
func (s *SessionService) RevokeSession(ctx context.Context, userID uuid.UUID, token string) error {
	session, err := s.GetSession(ctx, userID, token)
	if err != nil {
		return err
	}

	session.Status = SessionRevoked

	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("session:%s:%s", userID.String(), s.hashToken(token))

	// Оставляем TTL как есть
	currentTTL := s.redis.TTL(ctx, key).Val()
	return s.redis.Set(ctx, key, sessionData, currentTTL).Err()
}

// RevokeAllUserSessions отзывает все сессии пользователя
func (s *SessionService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("session:%s:*", userID.String())

	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		// Получаем сессию
		data, err := s.redis.Get(ctx, key).Result()
		if err != nil {
			continue // Пропускаем ошибочные ключи
		}

		var session JWTSession
		if err := json.Unmarshal([]byte(data), &session); err != nil {
			continue
		}

		// Обновляем статус
		session.Status = SessionRevoked
		sessionData, err := json.Marshal(session)
		if err != nil {
			continue
		}

		// Сохраняем обновленную сессию
		currentTTL := s.redis.TTL(ctx, key).Val()
		s.redis.Set(ctx, key, sessionData, currentTTL)
	}

	return iter.Err()
}

// IsSessionValid проверяет валидность сессии
func (s *SessionService) IsSessionValid(ctx context.Context, userID uuid.UUID, token string) (bool, error) {
	session, err := s.GetSession(ctx, userID, token)
	if err != nil {
		return false, err
	}

	// Проверяем статус
	if session.Status != SessionActive {
		return false, nil
	}

	// Проверяем срок действия
	if time.Now().After(session.ExpiresAt) {
		// Обновляем статус на expired
		session.Status = SessionExpired
		sessionData, _ := json.Marshal(session)
		key := fmt.Sprintf("session:%s:%s", userID.String(), s.hashToken(token))
		s.redis.Set(ctx, key, sessionData, time.Minute) // Короткий TTL для expired
		return false, nil
	}

	return true, nil
}

// hashToken создает SHA256 хеш токена для использования в ключе Redis
func (s *SessionService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

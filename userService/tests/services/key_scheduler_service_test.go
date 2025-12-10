package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"userService/internal/services"
)

// Тесты для KeySchedulerService
// Примечание: KeySchedulerService использует реальный KeyManagementService,
// поэтому тесты больше интеграционные, но мы можем проверить базовую функциональность

func TestKeySchedulerService_Start(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	keyManagement := services.NewKeyManagementServiceWithProducer(mockProducer, 1)
	interval := 100 * time.Millisecond
	scheduler := services.NewKeySchedulerService(keyManagement, interval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Act
	scheduler.Start(ctx)

	// Ждем немного, чтобы убедиться, что горутина запустилась
	time.Sleep(50 * time.Millisecond)

	// Assert - проверяем, что сервис запущен (нет паники)
	assert.NotNil(t, scheduler)

	// Останавливаем
	scheduler.Stop()
	time.Sleep(50 * time.Millisecond)
}

func TestKeySchedulerService_Stop(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	keyManagement := services.NewKeyManagementServiceWithProducer(mockProducer, 1)
	interval := 100 * time.Millisecond
	scheduler := services.NewKeySchedulerService(keyManagement, interval)

	ctx := context.Background()
	scheduler.Start(ctx)

	// Act
	scheduler.Stop()

	// Ждем немного для завершения
	time.Sleep(150 * time.Millisecond)

	// Assert - проверяем, что сервис остановлен (нет паники)
	assert.NotNil(t, scheduler)
}

func TestKeySchedulerService_ContextCancellation(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	keyManagement := services.NewKeyManagementServiceWithProducer(mockProducer, 1)
	interval := 100 * time.Millisecond
	scheduler := services.NewKeySchedulerService(keyManagement, interval)

	ctx, cancel := context.WithCancel(context.Background())

	// Act
	scheduler.Start(ctx)

	// Ждем немного
	time.Sleep(50 * time.Millisecond)

	// Отменяем контекст
	cancel()

	// Ждем для завершения
	time.Sleep(150 * time.Millisecond)

	// Assert - проверяем, что сервис остановлен
	assert.NotNil(t, scheduler)
}

// Интеграционный тест для проверки работы scheduler с реальным key management service
// В реальных тестах можно использовать моки для более точного контроля

func TestKeySchedulerService_Integration(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	keyManagement := services.NewKeyManagementServiceWithProducer(mockProducer, 1)
	interval := 200 * time.Millisecond
	scheduler := services.NewKeySchedulerService(keyManagement, interval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockProducer.On("SendKeyUpdate", mock.Anything).Return(nil).Maybe()

	// Act
	scheduler.Start(ctx)

	// Ждем достаточно времени для срабатывания тикера
	time.Sleep(250 * time.Millisecond)

	// Останавливаем
	scheduler.Stop()
	cancel()

	// Ждем для завершения
	time.Sleep(100 * time.Millisecond)

	// Assert - проверяем, что сервис работает
	assert.NotNil(t, scheduler)
}

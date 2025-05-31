package services

import (
	"context"
	"log"
	"time"
)

type KeySchedulerService struct {
	keyManagement *KeyManagementService
	interval      time.Duration
	ticker        *time.Ticker
	done          chan bool
}

func NewKeySchedulerService(keyManagement *KeyManagementService, interval time.Duration) *KeySchedulerService {
	return &KeySchedulerService{
		keyManagement: keyManagement,
		interval:      interval,
		done:          make(chan bool),
	}
}

func (kss *KeySchedulerService) Start(ctx context.Context) {
	log.Printf("Starting key scheduler with interval: %v", kss.interval)

	kss.ticker = time.NewTicker(kss.interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Key scheduler context cancelled")
				return
			case <-kss.done:
				log.Println("Key scheduler stopped")
				return
			case <-kss.ticker.C:
				log.Println("Scheduled key regeneration triggered")
				if err := kss.keyManagement.RegenerateKeys(); err != nil {
					log.Printf("Error during scheduled key regeneration: %v", err)
				} else {
					log.Println("Scheduled key regeneration completed successfully")
				}
			}
		}
	}()
}

func (kss *KeySchedulerService) Stop() {
	log.Println("Stopping key scheduler...")

	if kss.ticker != nil {
		kss.ticker.Stop()
	}

	select {
	case kss.done <- true:
	default:
	}
}

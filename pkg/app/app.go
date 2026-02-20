package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"pingo/pkg/config"
	"pingo/pkg/ping"
	"pingo/pkg/storage"
)

type Application struct {
	config  *config.Config
	storage *storage.ClickHouseStorage
}

func NewApplication(cfg *config.Config) (*Application, error) {
	store, err := storage.NewClickHouseStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %v", err)
	}

	// Initialize database
	ctx := context.Background()
	if err := store.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %v", err)
	}

	return &Application{
		config:  cfg,
		storage: store,
	}, nil
}

func (a *Application) Run() error {
	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Start ping workers for each target
	for _, item := range a.config.Items {
		wg.Add(1)
		go a.pingWorker(ctx, &wg, item)
	}

	// Wait for termination signal
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v. Shutting down...", sig)
		cancel()
	}

	wg.Wait()

	// Close storage connection
	if err := a.storage.Close(); err != nil {
		log.Printf("Error closing storage: %v", err)
	}

	return nil
}

func (a *Application) pingWorker(ctx context.Context, wg *sync.WaitGroup, cfg config.PingConfig) {
	defer wg.Done()

	pingInterval := config.GetPingInterval(a.config.PingInterval)
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result, err := ping.Ping(ctx, cfg.Tool, cfg.Target, cfg.Timeout)
			if err != nil {
				log.Printf("Ping failed for %s: %v", cfg.Target, err)
			}

			// Store metric
			metric := &storage.PingMetric{
				Name:       cfg.Name,
				Target:     result.Target,
				Success:    result.Success,
				DurationMs: float64(result.Duration.Milliseconds()),
				Error:      "",
				Timestamp:  result.Timestamp,
			}

			if result.Error != nil {
				metric.Error = result.Error.Error()
			}

			if err := a.storage.StoreMetric(ctx, metric); err != nil {
				log.Printf("Failed to store metric for %s: %v", cfg.Target, err)
			} else {
				log.Printf("Ping to %s: success=%t, duration=%v", cfg.Target, result.Success, result.Duration)
			}
		}
	}
}

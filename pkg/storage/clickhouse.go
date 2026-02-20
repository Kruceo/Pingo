package storage

import (
	"context"
	"pingo/pkg/config"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type PingMetric struct {
	Name       string
	Target     string
	Success    bool
	DurationMs float64
	Error      string
	Timestamp  time.Time
}

type ClickHouseStorage struct {
	conn driver.Conn
}

func NewClickHouseStorage() (*ClickHouseStorage, error) {
	chConfig := config.LoadClickHouseConfig()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{chConfig.DSN},
		Auth: clickhouse.Auth{
			Database: chConfig.Database,
			Username: chConfig.Username,
			Password: chConfig.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	return &ClickHouseStorage{conn: conn}, nil
}

func (s *ClickHouseStorage) Initialize(ctx context.Context) error {
	return s.conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS ping_metrics (
			name String,
			target String,
			success UInt8,
			duration_ms Float64,
			error String,
			timestamp DateTime
		) ENGINE = MergeTree()
		ORDER BY timestamp
	`)
}

func (s *ClickHouseStorage) StoreMetric(ctx context.Context, metric *PingMetric) error {
	return s.conn.AsyncInsert(ctx, `
		INSERT INTO ping_metrics (name, target, success, duration_ms, error, timestamp)
		VALUES (?, ?, ?, ?, ?, ?)
	`, false, metric.Name, metric.Target, metric.Success, metric.DurationMs, metric.Error, metric.Timestamp)
}

func (s *ClickHouseStorage) Close() error {
	return s.conn.Close()
}

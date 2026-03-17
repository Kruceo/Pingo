package ping

import (
	"context"
	"time"

	"github.com/go-ping/ping"
)

type PingResult struct {
	Target    string
	Success   bool
	Duration  time.Duration
	Error     error
	Timestamp time.Time
}

func PingV4(ctx context.Context, target string, timeout time.Duration) (*PingResult, error) {
	start := time.Now()

	// Criar pinger
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return &PingResult{
			Target:    target,
			Success:   false,
			Duration:  time.Since(start),
			Error:     err,
			Timestamp: time.Now(),
		}, err
	}

	// Configurar para modo não privilegiado (não requer root)
	pinger.SetPrivileged(false)

	// Configurar timeout
	pinger.Timeout = timeout

	// Configurar para IPv4
	pinger.SetNetwork("udp4")

	// Realizar apenas 1 ping
	pinger.Count = 1

	// Executar o ping
	err = pinger.Run()
	if err != nil {
		return &PingResult{
			Target:    target,
			Success:   false,
			Duration:  time.Since(start),
			Error:     err,
			Timestamp: time.Now(),
		}, err
	}

	// Obter estatísticas
	stats := pinger.Statistics()

	// Verificar se houve resposta
	success := stats.PacketsRecv > 0

	return &PingResult{
		Target:    target,
		Success:   success,
		Duration:  time.Since(start),
		Error:     nil,
		Timestamp: time.Now(),
	}, nil
}

package ping

import (
	"context"
	"time"

	"github.com/go-ping/ping"
)

func PingV6(ctx context.Context, target string, timeout time.Duration) (*PingResult, error) {
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

	// Configurar para IPv6
	pinger.SetNetwork("udp6")

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

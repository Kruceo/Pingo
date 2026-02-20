package ping

import (
	"context"
	"net"
	"time"
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
	
	conn, err := net.DialTimeout("ip4:icmp", target, timeout)
	if err != nil {
		return &PingResult{
			Target:    target,
			Success:   false,
			Duration:  time.Since(start),
			Error:     err,
			Timestamp: time.Now(),
		}, err
	}
	defer conn.Close()
	
	// Set deadline for the connection
	conn.SetDeadline(time.Now().Add(timeout))
	
	// Send ICMP echo request
	msg := make([]byte, 64)
	msg[0] = 8 // ICMP echo request
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	
	// Calculate checksum
	cs := checksum(msg)
	msg[2] = byte(cs >> 8)
	msg[3] = byte(cs & 0xff)
	
	_, err = conn.Write(msg)
	if err != nil {
		return &PingResult{
			Target:    target,
			Success:   false,
			Duration:  time.Since(start),
			Error:     err,
			Timestamp: time.Now(),
		}, err
	}
	
	// Read response
	recv := make([]byte, 1024)
	_, err = conn.Read(recv)
	if err != nil {
		return &PingResult{
			Target:    target,
			Success:   false,
			Duration:  time.Since(start),
			Error:     err,
			Timestamp: time.Now(),
		}, err
	}
	
	return &PingResult{
		Target:    target,
		Success:   true,
		Duration:  time.Since(start),
		Error:     nil,
		Timestamp: time.Now(),
	}, nil
}

func checksum(msg []byte) uint16 {
	sum := 0
	for n := 0; n < len(msg); n += 2 {
		sum += int(msg[n])<<8 | int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	return uint16(^sum)
}
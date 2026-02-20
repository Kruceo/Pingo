package ping

import (
	"context"
	"net"
	"time"
)

func PingV6(ctx context.Context, target string, timeout time.Duration) (*PingResult, error) {
	start := time.Now()

	conn, err := net.DialTimeout("ip6:ipv6-icmp", target, timeout)
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

	// Send ICMPv6 echo request
	msg := make([]byte, 64)
	msg[0] = 128 // ICMPv6 echo request
	msg[1] = 0   // Code
	msg[2] = 0   // Checksum (high byte)
	msg[3] = 0   // Checksum (low byte)
	msg[4] = 0   // Identifier (high byte)
	msg[5] = 1   // Identifier (low byte)
	msg[6] = 0   // Sequence number (high byte)
	msg[7] = 1   // Sequence number (low byte)

	// Fill the rest with some data (56 bytes to make 64 total)
	for i := 8; i < 64; i++ {
		msg[i] = byte(i)
	}

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

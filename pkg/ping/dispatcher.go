package ping

import (
	"context"
	"fmt"
	"time"
)

// Ping dispatches to the appropriate ping function based on the tool type
func Ping(ctx context.Context, tool, target string, timeout time.Duration) (*PingResult, error) {
	switch tool {
	case "pingv4":
		return PingV4(ctx, target, timeout)
	case "pingv6":
		return PingV6(ctx, target, timeout)
	default:
		return nil, fmt.Errorf("unsupported ping tool: %s", tool)
	}
}
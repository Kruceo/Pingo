package ping

import (
	"context"
	"fmt"
	"time"
)

var supportedTools = map[string]struct{}{
	"pingv4": {},
	"pingv6": {},
}

func IsSupportedTool(tool string) bool {
	_, ok := supportedTools[tool]
	return ok
}

func SupportedTools() []string {
	out := make([]string, 0, len(supportedTools))
	for tool := range supportedTools {
		out = append(out, tool)
	}
	return out
}

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

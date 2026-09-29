package metrics

import (
	"fmt"
	"strings"
	"time"
)

// DefaultWindow is the window flai stats aggregates over when none is given.
const DefaultWindow = "30d"

// ParseWindow reads a window for completed items: days (30d), weeks (12w),
// or a Go duration (720h).
func ParseWindow(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "d") || strings.HasSuffix(s, "w") {
		var n int
		unit := s[len(s)-1]
		if _, err := fmt.Sscanf(s[:len(s)-1], "%d", &n); err != nil || n <= 0 {
			return 0, fmt.Errorf("--since expects a window like 30d, 12w, or 720h, got %q", s)
		}
		if unit == 'w' {
			n *= 7
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("--since expects a window like 30d, 12w, or 720h, got %q", s)
	}
	return d, nil
}

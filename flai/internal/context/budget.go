package context

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultBudget is the size a context pack fits when neither --budget nor
// the manifest's prime.budget says otherwise: 80 KB, about 20k tokens, so
// that the pack fits an MCP tool result with room for the agent's own turn
// (ADR-0049).
const DefaultBudget = 80 * 1024

// PartLimit is the most bytes of JSON one part of a pack encodes to: a
// fifth under the 50,000 characters over which Claude Code saves a tool
// result to a file instead of showing it, for the result's envelope and a
// harness that counts differently (ADR-0104).
const PartLimit = 40000

// ParseSize reads a budget: a whole number of bytes, or one followed by K,
// KB, or KiB (1024 bytes), or M, MB, or MiB, in any case, with or without a
// space before the unit.
func ParseSize(s string) (int, error) {
	in := strings.TrimSpace(s)
	num := strings.TrimRightFunc(in, func(r rune) bool { return r < '0' || r > '9' })
	unit := strings.ToUpper(strings.TrimSpace(in[len(num):]))
	mult := map[string]int{"": 1, "B": 1, "K": 1024, "KB": 1024, "KIB": 1024, "M": 1024 * 1024, "MB": 1024 * 1024, "MIB": 1024 * 1024}[unit]
	n, err := strconv.Atoi(num)
	if err != nil || mult == 0 || n <= 0 {
		return 0, fmt.Errorf("budget %q is not a size; give bytes, such as 81920, or KB or MB, such as 80KB", s)
	}
	return n * mult, nil
}

// Budget is the budget a pack is built to: flag when it is set, else the
// project's prime.budget, else DefaultBudget.
func Budget(flag, project string) (int, error) {
	for _, s := range []string{flag, project} {
		if s != "" {
			return ParseSize(s)
		}
	}
	return DefaultBudget, nil
}

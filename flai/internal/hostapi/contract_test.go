package hostapi

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The dashboard lists the methods it asks for (REQUIRED_METHODS in
// flaiover's agent.ts) and tells the designer when the connected flai lacks
// one. That list and this table are one contract kept in two languages; this
// test fails when they drift either way.
func TestTheDashboardAsksForExactlyWhatFlaiOffers(t *testing.T) {
	src, err := os.ReadFile("../../../flaiover/src/lib/server/agent.ts")
	if err != nil {
		t.Skip("flaiover is not beside flai in this checkout: " + err.Error())
	}
	block := regexp.MustCompile(`(?s)REQUIRED_METHODS = \[(.*?)\];`).FindSubmatch(src)
	if block == nil {
		t.Fatal("REQUIRED_METHODS not found in agent.ts")
	}
	asked := map[string]bool{}
	for _, m := range regexp.MustCompile(`'([a-z.]+)'`).FindAllSubmatch(block[1], -1) {
		asked[string(m[1])] = true
	}
	offered := Methods("test", nil)
	var missing, unasked []string
	for name := range asked {
		if _, ok := offered[name]; !ok {
			missing = append(missing, name)
		}
	}
	for name := range offered {
		if !asked[name] {
			unasked = append(unasked, name)
		}
	}
	sort.Strings(missing)
	sort.Strings(unasked)
	if len(missing) > 0 {
		t.Errorf("the dashboard asks for methods flai does not offer: %s", strings.Join(missing, ", "))
	}
	if len(unasked) > 0 {
		t.Errorf("flai offers methods the dashboard does not list in REQUIRED_METHODS: %s", strings.Join(unasked, ", "))
	}
}

package hostapi

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"slices"
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
	for _, m := range regexp.MustCompile(`'([a-z._]+)'`).FindAllSubmatch(block[1], -1) {
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

// S-0420: agent.status's state carries each role under roles, by name, as a
// RoleState, and the dashboard types that entry as RoleState in flaiover's
// roles.ts. This test fails when the fields the two name drift either way,
// and when a field is not in the JSON flai answers; it skips while the
// dashboard does not type the entry yet.
func TestTheDashboardReadsTheRoleEntryFlaiAnswers(t *testing.T) {
	type run struct{}
	var answered map[string]any
	raw, err := json.Marshal(RoleState[run]{})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &answered); err != nil {
		t.Fatal(err)
	}
	var fields []string
	rt := reflect.TypeFor[RoleState[run]]()
	for i := range rt.NumField() {
		name, _, _ := strings.Cut(rt.Field(i).Tag.Get("json"), ",")
		if _, ok := answered[name]; !ok {
			t.Errorf("the field %s is left out of the JSON when empty; the dashboard reads every field", name)
		}
		fields = append(fields, name)
	}
	if want := []string{"name", "role", "path", "item", "trigger", "action", "run", "running", "runs", "held"}; !slices.Equal(fields, want) {
		t.Errorf("a role entry's fields are %v, want %v: change flaiover's RoleState with them", fields, want)
	}
	src, err := os.ReadFile("../../../flaiover/src/lib/roles.ts")
	if err != nil {
		t.Skip("flaiover's roles.ts is not beside flai in this checkout: " + err.Error())
	}
	block := regexp.MustCompile(`(?s)export (?:type RoleState(?:<[^>]*>)? = \{|interface RoleState(?:<[^>]*>)? \{)(.*?)\n\};?`).FindSubmatch(src)
	if block == nil {
		t.Skip("flaiover's roles.ts does not type RoleState yet")
	}
	var typed []string
	for _, m := range regexp.MustCompile(`(?m)^\t([a-z_]+)\??:`).FindAllSubmatch(block[1], -1) {
		typed = append(typed, string(m[1]))
	}
	sort.Strings(typed)
	sort.Strings(fields)
	if !slices.Equal(typed, fields) {
		t.Errorf("flaiover's RoleState types %v; flai answers %v", typed, fields)
	}
}

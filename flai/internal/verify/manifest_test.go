package verify

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// checkout writes a manifest with the given tests block, and the files
// named, into a temporary checkout.
func checkout(t *testing.T, tests string, files ...string) string {
	t.Helper()
	root := t.TempDir()
	m := "version: 1\nname: fx\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" + tests
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		p := filepath.Join(root, f)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// S-0273: the manifest's tiers reach a run as declared, the format
// defaulted to plain; with no tests key the default tier runs
// scripts/test.sh; tiers that are not valid are refused with what to fix.
func TestCheckoutTiers(t *testing.T) {
	root := checkout(t, "tests:\n  - name: unit\n    command: [go, test, \"{packages}\"]\n    dir: flai\n    paths: [\"flai/**\", \"!flai/testdata/**\"]\n    format: go-test-json\n  - name: smoke\n    command: [scripts/smoke.sh]\n    all_only: true\n")
	got, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []Tier{
		{Name: "unit", Command: []string{"go", "test", "{packages}"}, Dir: "flai", Paths: []string{"flai/**", "!flai/testdata/**"}, Format: FormatGoTestJSON},
		{Name: "smoke", Command: []string{"scripts/smoke.sh"}, Format: FormatPlain, AllOnly: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tiers:\n got %+v\nwant %+v", got, want)
	}

	got, err = CheckoutTiers(checkout(t, "", "scripts/test.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if want := []Tier{{Name: "test", Command: []string{"scripts/test.sh"}, Paths: []string{"**"}, Format: FormatPlain}}; !reflect.DeepEqual(got, want) {
		t.Errorf("default tier: got %+v, want %+v", got, want)
	}

	_, err = CheckoutTiers(checkout(t, "tests:\n  - name: unit\n    command: [go, test]\n    paths: [\"flai/**\"]\n    format: junit\n"))
	if err == nil || !strings.Contains(err.Error(), "tests[0].format") || !strings.Contains(err.Error(), "flai manifest set tests=") {
		t.Errorf("a bad tier: %v, want it refused naming tests[0].format and how to fix it", err)
	}
}

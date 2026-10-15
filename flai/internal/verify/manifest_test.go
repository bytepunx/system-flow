package verify

import (
	"errors"
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
// defaulted to plain and covers copied (S-0342); with no tests key the default tier runs
// scripts/test.sh; tiers that are not valid are refused with what to fix.
func TestCheckoutTiers(t *testing.T) {
	root := checkout(t, "tests:\n  - name: unit\n    command: [go, test, \"{packages}\"]\n    dir: flai\n    paths: [\"flai/**\", \"!flai/testdata/**\"]\n    format: go-test-json\n  - name: smoke\n    command: [scripts/smoke.sh]\n    all_only: true\n    covers: [unit]\n")
	got, err := CheckoutTiers(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []Tier{
		{Name: "unit", Command: []string{"go", "test", "{packages}"}, Dir: "flai", Paths: []string{"flai/**", "!flai/testdata/**"}, Format: FormatGoTestJSON},
		{Name: "smoke", Command: []string{"scripts/smoke.sh"}, Format: FormatPlain, AllOnly: true, Covers: []string{"unit"}},
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

// tagGit answers git tag with the flai release tags it holds, as git does in
// a project whose history has them, and refuses everything else.
type tagGit struct{ tags string }

func (g tagGit) Run(_, name string, args ...string) (string, error) {
	if name == "git" && strings.Join(args, " ") == "tag --merged HEAD --list flai/v*" {
		return g.tags, nil
	}
	return "", errors.New("fatal: not a git repository")
}

func (g tagGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (tagGit) LookPath(name string) (string, error) { return name, nil }

// I-0141, ADR-0158: tiers a newer flai would take, here a placeholder this
// flai does not fill, are refused as an older flai's when the checkout's
// history has a newer flai release, naming both versions and the upgrade
// rather than the manifest's fix; with no newer release, no runner, or a
// running flai that is not a release, the refusal names the fix as before.
func TestCheckoutTiersForSaysTheFlaiIsOlderThanTheProject(t *testing.T) {
	// each its own checkout, as FlaiOutdated keeps an answer per root for a minute
	bad := func() string {
		return checkout(t, "tests:\n  - name: vet\n    command: [go, vet, \"{dependents}\"]\n    paths: [\"flai/**\"]\n")
	}
	const problem = `tests[0].command (tier "vet") argument {dependents} is not a placeholder flai fills`

	_, err := CheckoutTiersFor(bad(), tagGit{tags: "flai/v1.40.1\nflai/v1.49.0\n"}, "1.40.1")
	if err == nil {
		t.Fatal("tiers that are not valid were taken")
	}
	for _, want := range []string{problem, "this flai is 1.40.1, older than flai 1.49.0 in the project's history", "upgrade it with flai host upgrade"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("behind a newer release: %q, want it to say %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "flai manifest set tests=") {
		t.Errorf("behind a newer release: %q tells the agent to fix the manifest", err)
	}

	for name, c := range map[string]struct {
		git     tagGit
		running string
	}{
		"no newer release":    {tagGit{tags: "flai/v1.40.1\n"}, "1.40.1"},
		"a dev build":         {tagGit{tags: "flai/v1.49.0\n"}, "dev"},
		"git cannot say":      {tagGit{}, "1.40.1"},
		"a newer running one": {tagGit{tags: "flai/v1.49.0\n"}, "1.50.0"},
	} {
		_, err := CheckoutTiersFor(bad(), c.git, c.running)
		if err == nil || !strings.Contains(err.Error(), problem) || !strings.Contains(err.Error(), "flai manifest set tests=") || strings.Contains(err.Error(), "upgrade it") {
			t.Errorf("%s: %v, want the manifest's fix", name, err)
		}
	}
	if _, err := CheckoutTiersFor(bad(), nil, "1.40.1"); err == nil || !strings.Contains(err.Error(), "flai manifest set tests=") {
		t.Errorf("no runner: %v, want the manifest's fix", err)
	}
}

package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

func TestLoadAndFind(t *testing.T) {
	root := t.TempDir()
	body := "version: 1\nname: demo\nkey: d\nlayout:\n  design: arch\n  docs: docs\n  wip: wip\nprojects:\n  - name: a\n    path: a\n    kind: go\n"
	if err := os.WriteFile(filepath.Join(root, File), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "a", "b")
	_ = os.MkdirAll(nested, 0o755)
	p, err := Find(nested)
	if err != nil || p != filepath.Join(root, File) {
		t.Fatalf("Find = %q, %v", p, err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "demo" || m.Dir(root, "design") != filepath.Join(root, "arch") || len(m.Projects) != 1 {
		t.Fatalf("unexpected: %+v", m)
	}
	if _, err := Find(t.TempDir()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	_ = os.WriteFile(p, []byte("version: 1\nname: x\nlayout:\n  design: d\n"), 0o644)
	if _, err := Load(p); err == nil {
		t.Fatal("expected missing layout keys error")
	}
}

// S-0082: checks: is optional and round-trips as a list of named argument
// lists, the operator's choice of where to name them when the host's own
// configuration names none.
func TestChecksRoundTrip(t *testing.T) {
	root := t.TempDir()
	body := "version: 1\nname: demo\nkey: d\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n" +
		"checks:\n  - name: flai\n    command: [scripts/flai-test.sh]\n  - name: flaiover\n    command: [bash, -c, \"cd flaiover; pnpm test\"]\n"
	p := filepath.Join(root, File)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Checks) != 2 || m.Checks[0].Name != "flai" || len(m.Checks[0].Command) != 1 ||
		m.Checks[1].Name != "flaiover" || len(m.Checks[1].Command) != 3 {
		t.Fatalf("checks: %+v", m.Checks)
	}
	// absent is a nil slice, not an error and not an empty-but-present list
	m2, err := Load(func() string {
		body := "version: 1\nname: demo2\nkey: d2\nlayout:\n  design: d\n  docs: docs\n  wip: wip\n"
		p2 := filepath.Join(t.TempDir(), File)
		_ = os.WriteFile(p2, []byte(body), 0o644)
		return p2
	}())
	if err != nil || m2.Checks != nil {
		t.Errorf("no checks: named: %v %+v", err, m2.Checks)
	}
}

// S-0181: a flai below the manifest's minimum says so, naming the version
// needed and the upgrade, before anything else is read; a dev build and a
// release at or above it read the project.
func TestMinimumFlai(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	base := "version: 1\nname: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(path, []byte(base+"flai:\n  minimum: 1.27.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	was := buildinfo.Version
	t.Cleanup(func() { buildinfo.Version = was })

	buildinfo.Version = "1.26.4"
	_, err := Load(path)
	var old *TooOldError
	if !errors.As(err, &old) || old.Minimum != "1.27.0" || old.Running != "1.26.4" ||
		!strings.Contains(err.Error(), "needs flai 1.27.0 or newer, and this is flai 1.26.4") || !strings.Contains(err.Error(), "flai host upgrade") {
		t.Fatalf("a flai below the minimum: %v", err)
	}
	for _, v := range []string{"1.27.0", "1.30.1", "dev"} {
		buildinfo.Version = v
		if m, err := Load(path); err != nil || m.Flai.Minimum != "1.27.0" {
			t.Errorf("flai %s: %v %+v", v, err, m.Flai)
		}
	}
	_ = os.WriteFile(path, []byte(base+"flai:\n  minimum: soon\n"), 0o644)
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), `flai.minimum "soon" is not a release version`) {
		t.Errorf("a minimum that is not a version: %v", err)
	}
}

func TestSetMinimum(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	base := "version: 1\nname: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	for _, c := range []struct{ before, after string }{
		{base, base + "flai:\n  minimum: 1.27.0\n"},
		{base + "flai:\n  minimum: 1.26.0\n", base + "flai:\n  minimum: 1.27.0\n"},
		{"flai:\n  other: x\n" + base, "flai:\n  minimum: 1.27.0\n  other: x\n" + base},
	} {
		_ = os.WriteFile(path, []byte(c.before), 0o644)
		if err := SetMinimum(path, "1.27.0"); err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != c.after {
			t.Errorf("from:\n%s\ngot:\n%s\nwant:\n%s", c.before, got, c.after)
		}
		if m, err := Load(path); err != nil || m.Flai.Minimum != "1.27.0" {
			t.Errorf("reads back: %v %+v", err, m.Flai)
		}
	}
}

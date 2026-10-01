package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

func TestUpgradeCommand(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dest := filepath.Join(t.TempDir(), "proj")
	if _, errOut, code := runIn(t, ".", "new", dest, "--template", miniTemplate, "--defaults", "--no-git"); code != 0 {
		t.Fatalf("new: %s", errOut)
	}
	if _, err := os.Stat(filepath.Join(dest, "system-flow.lock.yaml")); err != nil {
		t.Fatal("flai new must write the lock")
	}
	// a newer template
	nt := t.TempDir()
	_ = filepath.WalkDir(miniTemplate, func(path string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(miniTemplate, path)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(nt, rel), 0o755)
		}
		b, _ := os.ReadFile(path)
		return os.WriteFile(filepath.Join(nt, rel), b, 0o644)
	})
	y, _ := os.ReadFile(filepath.Join(nt, "template.yaml"))
	_ = os.WriteFile(filepath.Join(nt, "template.yaml"), []byte(strings.Replace(string(y), "version: 9.9.9", "version: 10.0.0", 1)), 0o644)
	_ = os.WriteFile(filepath.Join(nt, "root", "plain.txt"), []byte("verbatim v2\n"), 0o644)
	_ = os.WriteFile(filepath.Join(nt, "root", "new.txt"), []byte("new\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dest, "plain.txt"), []byte("edited\n"), 0o644)

	out, _, code := runIn(t, dest, "upgrade", "--template", miniTemplate)
	if code != 0 || !strings.Contains(out, "already at template 9.9.9") {
		t.Fatalf("same version: %d %s", code, out)
	}
	out, _, code = runIn(t, dest, "upgrade", "--template", nt, "--dry-run")
	if code != 0 || !strings.Contains(out, "9.9.9 -> 10.0.0") || !strings.Contains(out, "conflict  plain.txt") || !strings.Contains(out, "add       new.txt") || !strings.Contains(out, "dry run") {
		t.Fatalf("dry run: %d\n%s", code, out)
	}
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); err == nil {
		t.Fatal("dry run wrote files")
	}
	_, errOut, code := runIn(t, dest, "upgrade", "--template", nt)
	if code != 1 || !strings.Contains(errOut, "conflict(s) need a decision") {
		t.Fatalf("non-interactive without policy must refuse: %d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); err == nil {
		t.Fatal("refused upgrade wrote files")
	}
	out, errOut, code = runIn(t, dest, "upgrade", "--template", nt, "--keep-all")
	if code != 0 || !strings.Contains(out, "1 conflicts kept") {
		t.Fatalf("keep-all: %d %s %s", code, out, errOut)
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "plain.txt")); string(b) != "edited\n" {
		t.Error("kept file overwritten")
	}
	if _, err := os.Stat(filepath.Join(dest, "new.txt")); err != nil {
		t.Error("new file not added")
	}
	mf, _ := os.ReadFile(filepath.Join(dest, "system-flow.yaml"))
	if !strings.Contains(string(mf), `version: "10.0.0"`) {
		t.Errorf("manifest version not updated:\n%s", mf)
	}
	out, _, code = runIn(t, dest, "upgrade", "--template", nt)
	if code != 0 || !strings.Contains(out, "already at template 10.0.0") {
		t.Errorf("after upgrade: %d %s", code, out)
	}
	// dirty tree guard with git
	if _, err := exec.LookPath("git"); err == nil {
		r := execx.System{}
		for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "x"}} {
			_, _ = r.Run(dest, "git", args...)
		}
		_ = os.WriteFile(filepath.Join(dest, "dirty.txt"), []byte("x"), 0o644)
		_, errOut, code = runIn(t, dest, "upgrade", "--template", nt, "--force", "--keep-all")
		if code != 0 {
			t.Errorf("--force should allow a dirty tree: %s", errOut)
		}
		_, _ = r.Run(dest, "git", "add", "-A")
		_, _ = r.Run(dest, "git", "commit", "-q", "-m", "y")
		_ = os.WriteFile(filepath.Join(dest, "dirty2.txt"), []byte("x"), 0o644)
		_ = os.WriteFile(filepath.Join(nt, "template.yaml"), []byte(strings.Replace(string(y), "version: 9.9.9", "version: 11.0.0", 1)), 0o644)
		_, errOut, code = runIn(t, dest, "upgrade", "--template", nt, "--keep-all")
		if code == 0 || !strings.Contains(errOut, "uncommitted changes") {
			t.Errorf("dirty tree must be refused: %d %s", code, errOut)
		}
	}
	// relock on a project without a lock
	_ = os.Remove(filepath.Join(dest, "system-flow.lock.yaml"))
	out, _, code = runIn(t, dest, "upgrade", "--template", nt, "--relock", "--force")
	if code != 0 || !strings.Contains(out, "relocked") {
		t.Errorf("relock: %d %s", code, out)
	}
}

// copyTemplate copies a template directory to a temporary one a test may change.
func copyTemplate(t *testing.T, from string) string {
	t.Helper()
	to := t.TempDir()
	err := filepath.WalkDir(from, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(from, path)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(to, rel), 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(to, rel), b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return to
}

// forkVersion rewrites a fork's version and appends variables to its
// manifest, and writes team.txt.tmpl with the given body.
func forkVersion(t *testing.T, tpl, version, addVars, body string) {
	t.Helper()
	y, err := os.ReadFile(filepath.Join(tpl, "template.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	s := regexp.MustCompile(`(?m)^version: .*$`).ReplaceAllString(string(y), "version: "+version)
	s = strings.Replace(s, "layout:\n", addVars+"layout:\n", 1)
	if err := os.WriteFile(filepath.Join(tpl, "template.yaml"), []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tpl, "root", "team.txt.tmpl"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A fork's own variable is recorded by flai new and rendered with by flai
// upgrade; a variable added since gets its default, --var changes one, and a
// required one with no value is refused before anything changes (I-0040).
func TestUpgradeRendersAForksOwnVariables(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	tpl := copyTemplate(t, miniTemplate)
	forkVersion(t, tpl, "1.0.0", "  - name: team\n    required: true\n", "team={{ .team }}\n")
	dest := filepath.Join(t.TempDir(), "proj")
	if _, errOut, code := runIn(t, ".", "new", dest, "--template", tpl, "--defaults", "--no-git", "--var", "team=core"); code != 0 {
		t.Fatalf("new: %s", errOut)
	}
	lk, err := lock.Load(dest)
	if err != nil || lk == nil || lk.Vars["team"] != "core" || lk.Vars["project_name"] != "proj" {
		t.Fatalf("flai new must record the variables in the lock: %+v %v", lk, err)
	}
	readTeam := func() string {
		b, _ := os.ReadFile(filepath.Join(dest, "team.txt"))
		return string(b)
	}

	// 2.0.0 adds region with a default: the recorded team and the default render.
	forkVersion(t, tpl, "2.0.0", "  - name: region\n    default: \"{{ .team }}-eu\"\n", "team={{ .team }} region={{ .region }}\n")
	out, errOut, code := runIn(t, dest, "upgrade", "--template", tpl)
	if code != 0 {
		t.Fatalf("upgrade with the recorded value and a default: %s", errOut)
	}
	if !strings.Contains(out, `region has no recorded value: rendering with its default "core-eu"`) {
		t.Errorf("the upgrade does not name the variable that took its default:\n%s", out)
	}
	if got := readTeam(); got != "team=core region=core-eu\n" {
		t.Errorf("team.txt after 2.0.0 = %q", got)
	}
	if lk, _ := lock.Load(dest); lk.Vars["region"] != "core-eu" || lk.Vars["team"] != "core" {
		t.Errorf("upgrade must record what it rendered with: %v", lk.Vars)
	}

	// --var changes a recorded value, and the lock keeps the new one.
	forkVersion(t, tpl, "3.0.0", "", "team={{ .team }} region={{ .region }} v3\n")
	if _, errOut, code := runIn(t, dest, "upgrade", "--template", tpl, "--var", "team=billing"); code != 0 {
		t.Fatalf("upgrade with --var: %s", errOut)
	}
	if got := readTeam(); got != "team=billing region=core-eu v3\n" {
		t.Errorf("team.txt after --var = %q", got)
	}
	if lk, _ := lock.Load(dest); lk.Vars["team"] != "billing" {
		t.Errorf("--var not recorded: %v", lk.Vars)
	}

	// At the version the project is at, --var re-applies it without --force.
	if _, errOut, code := runIn(t, dest, "upgrade", "--template", tpl, "--var", "team=ops"); code != 0 {
		t.Fatalf("--var at the same version: %s", errOut)
	}
	if got := readTeam(); got != "team=ops region=core-eu v3\n" {
		t.Errorf("team.txt after --var at the same version = %q", got)
	}

	// --var for an unknown variable or one the manifest holds is refused.
	if _, errOut, code := runIn(t, dest, "upgrade", "--template", tpl, "--var", "bogus=1"); code == 0 || !strings.Contains(errOut, "unknown variable") {
		t.Errorf("unknown --var: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, dest, "upgrade", "--template", tpl, "--var", "project_name=x"); code == 0 || !strings.Contains(errOut, "project_name is read from name in system-flow.yaml") {
		t.Errorf("--var for a manifest value: %d %s", code, errOut)
	}

	// 4.0.0 adds a required variable with no default: refused, nothing written.
	forkVersion(t, tpl, "4.0.0", "  - name: cost_centre\n    required: true\n", "{{ .cost_centre }}\n")
	before := readTeam()
	lockBefore, _ := os.ReadFile(filepath.Join(dest, "system-flow.lock.yaml"))
	for _, args := range [][]string{{"upgrade", "--template", tpl}, {"upgrade", "--template", tpl, "--dry-run"}} {
		_, errOut, code = runIn(t, dest, args...)
		if code == 0 || !strings.Contains(errOut, "required variables have no value: cost_centre") {
			t.Fatalf("%v with a required variable with no value: %d %s", args, code, errOut)
		}
	}
	if readTeam() != before {
		t.Error("a refused upgrade changed team.txt")
	}
	if b, _ := os.ReadFile(filepath.Join(dest, "system-flow.lock.yaml")); string(b) != string(lockBefore) {
		t.Error("a refused upgrade changed the lock")
	}
	if mf, _ := os.ReadFile(filepath.Join(dest, "system-flow.yaml")); !strings.Contains(string(mf), `version: "3.0.0"`) {
		t.Errorf("a refused upgrade changed the manifest:\n%s", mf)
	}
	if _, errOut, code := runIn(t, dest, "upgrade", "--template", tpl, "--var", "cost_centre=42"); code != 0 {
		t.Fatalf("upgrade with the required --var: %s", errOut)
	}
	if got := readTeam(); got != "42\n" {
		t.Errorf("team.txt after 4.0.0 = %q", got)
	}

	// --relock with no lock records the values it rendered with.
	_ = os.Remove(filepath.Join(dest, "system-flow.lock.yaml"))
	if _, errOut, code := runIn(t, dest, "upgrade", "--template", tpl, "--relock", "--var", "team=ops", "--var", "cost_centre=7"); code != 0 {
		t.Fatalf("relock: %s", errOut)
	}
	if lk, _ := lock.Load(dest); lk == nil || lk.Vars["team"] != "ops" || lk.Vars["cost_centre"] != "7" || lk.Vars["region"] != "ops-eu" {
		t.Errorf("relock must record the variables: %+v", lk)
	}
}

// A required variable recorded empty takes the default a newer template
// gives it, and the values the lock records survive any characters.
func TestUpgradeVarsAndTheLockRoundTrip(t *testing.T) {
	m := template.Manifest{Variables: []template.Variable{
		{Name: "project_name", Required: true},
		{Name: "zone", Required: true, Default: "z1"},
		{Name: "note"},
	}}
	odd := "a \"quoted\": value\\ with\nnewline # and - more"
	lk := &lock.Lock{Vars: map[string]string{"zone": "", "note": odd, "gone": "x"}}
	vars, defaulted, err := upgradeVars(m, manifest.Manifest{Name: "p"}, lk, nil)
	if err != nil {
		t.Fatal(err)
	}
	if vars["zone"] != "z1" || vars["note"] != odd || vars["project_name"] != "p" || strings.Join(defaulted, ",") != "zone" {
		t.Errorf("vars %v defaulted %v", vars, defaulted)
	}
	if _, ok := vars["gone"]; ok {
		t.Error("a variable the template no longer defines is rendered with")
	}
	dir := t.TempDir()
	if err := lock.Save(dir, &lock.Lock{Vars: lockVars(vars)}); err != nil {
		t.Fatal(err)
	}
	back, err := lock.Load(dir)
	if err != nil || back.Vars["note"] != odd || back.Vars["zone"] != "z1" {
		t.Errorf("lock round trip: %+v %v", back, err)
	}
}

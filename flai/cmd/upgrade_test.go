package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/prompt"
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

	// --relock with an unreadable lock, or none, records the values it rendered with.
	if err := os.WriteFile(filepath.Join(dest, "system-flow.lock.yaml"), []byte("<<<<<<< HEAD\nfiles: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, errOut, code := runIn(t, dest, "upgrade", "--template", tpl); code != 0 || !strings.Contains(out, "already at template 4.0.0") {
		t.Errorf("an unreadable lock must not stop a project already at the version: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, dest, "upgrade", "--template", tpl, "--relock", "--var", "team=ops", "--var", "cost_centre=7"); code != 0 {
		t.Fatalf("relock: %s", errOut)
	}
	if lk, _ := lock.Load(dest); lk == nil || lk.Vars["team"] != "ops" || lk.Vars["cost_centre"] != "7" || lk.Vars["region"] != "ops-eu" {
		t.Errorf("relock must record the variables: %+v", lk)
	}
}

// --ref given alone wins over the manifest's template.ref: the upgrade renders
// the tag it names and records that ref and its version in system-flow.yaml
// and the lock, and --relock with --ref records it the same way.
func TestUpgradeRefWinsOverTheManifestsRef(t *testing.T) {
	tpl, url, commitVersion := gitTemplate(t)
	dest := newFromGit(t, url, "main")
	commitVersion("10.0.0", "v10.0.0")
	commitVersion("11.0.0", "")

	out, errOut, code := runIn(t, dest, "upgrade", "--ref", "v10.0.0")
	if code != 0 || !strings.Contains(out, "upgraded to template 10.0.0") || !strings.Contains(out, "template v10.0.0 at 10.0.0: "+reasonRef) {
		t.Fatalf("upgrade --ref v10.0.0: %d %s %s", code, out, errOut)
	}
	if got := readTeamTxt(t, dest); got != "v10.0.0\n" {
		t.Errorf("team.txt = %q, want the tagged render", got)
	}
	recordedTemplate(t, dest, url, "v10.0.0", "10.0.0")

	// A --ref naming the version the project is at changes nothing.
	if out, _, code := runIn(t, dest, "upgrade", "--ref", "v10.0.0"); code != 0 || !strings.Contains(out, "already at template 10.0.0") {
		t.Errorf("same version by --ref: %d %s", code, out)
	}

	// --ref 11.0.0 names the tag v11.0.0, and the tag is what is recorded.
	gitIn(t, tpl, "tag", "v11.0.0")
	if out, errOut, code := runIn(t, dest, "upgrade", "--relock", "--ref", "11.0.0"); code != 0 || !strings.Contains(out, "relocked") {
		t.Fatalf("relock --ref 11.0.0: %d %s %s", code, out, errOut)
	}
	recordedTemplate(t, dest, url, "v11.0.0", "11.0.0")
}

// gitTemplate makes a git template repository in a temporary directory, at
// version 9.9.9 on main, and returns its path, its file:// URL, and a func
// that commits a new version on the checked-out branch, tagged when tag is
// not empty. team.txt renders as v<version>.
func gitTemplate(t *testing.T) (tpl, url string, commitVersion func(version, tag string)) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	tpl = copyTemplate(t, miniTemplate)
	forkVersion(t, tpl, "9.9.9", "", "v9.9.9\n")
	gitIn(t, tpl, "init", "-q", "-b", "main")
	gitIn(t, tpl, "add", "-A")
	gitIn(t, tpl, "commit", "-q", "-m", "9.9.9")
	commitVersion = func(version, tag string) {
		t.Helper()
		forkVersion(t, tpl, version, "", "v"+version+"\n")
		gitIn(t, tpl, "commit", "-q", "-am", version)
		if tag != "" {
			gitIn(t, tpl, "tag", tag)
		}
	}
	return tpl, "file://" + tpl, commitVersion
}

// newFromGit makes a project with flai new from a git template at ref.
func newFromGit(t *testing.T, url, ref string) string {
	t.Helper()
	dest := filepath.Join(t.TempDir(), "proj")
	if _, errOut, code := runIn(t, ".", "new", dest, "--template", url, "--ref", ref, "--defaults", "--no-git"); code != 0 {
		t.Fatalf("new: %s", errOut)
	}
	return dest
}

// recordedTemplate checks that system-flow.yaml and the lock record the template at
// url, ref, and version.
func recordedTemplate(t *testing.T, dest, url, ref, version string) {
	t.Helper()
	mf, err := manifest.Load(filepath.Join(dest, manifest.File))
	if err != nil || mf.Template.Repo != url || mf.Template.Ref != ref || mf.Template.Version != version {
		t.Errorf("system-flow.yaml must record %s at %s from %s: %+v %v", ref, version, url, mf.Template, err)
	}
	lk, err := lock.Load(dest)
	if err != nil || lk == nil || lk.Template.Ref != ref || lk.Template.Version != version || lk.Template.Repo != url {
		t.Errorf("lock must record %s at %s from %s: %+v %v", ref, version, url, lk, err)
	}
}

func readTeamTxt(t *testing.T, dest string) string {
	t.Helper()
	b, _ := os.ReadFile(filepath.Join(dest, "team.txt"))
	return string(b)
}

// editTemplateField sets one field of the template block of system-flow.yaml,
// as an operator editing it by hand would.
func editTemplateField(t *testing.T, dest, field, value string) {
	t.Helper()
	path := filepath.Join(dest, manifest.File)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	in, done := false, false
	for i, l := range lines {
		switch {
		case strings.HasPrefix(l, "template:"):
			in = true
		case in && strings.HasPrefix(l, "  "+field+":"):
			lines[i], done = "  "+field+": "+value, true
		case in && !strings.HasPrefix(l, "  "):
			in = false
		}
	}
	if !done {
		t.Fatalf("no template.%s in %s", field, b)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
}

// upgradeState reads the files an upgrade that changes nothing must leave alone.
func upgradeState(t *testing.T, dest string) string {
	t.Helper()
	var b strings.Builder
	for _, f := range []string{manifest.File, lock.File, "team.txt"} {
		c, err := os.ReadFile(filepath.Join(dest, f))
		if err != nil {
			t.Fatal(err)
		}
		b.WriteString(f + "\n" + string(c) + "\n")
	}
	return b.String()
}

// upgradeIn runs flai in dir with the template cache in cache, kept from one
// run to the next, on a terminal that answers reply, or with no terminal
// when reply is empty.
func upgradeIn(t *testing.T, dir, cache, reply string, args ...string) (string, string, int) {
	t.Helper()
	t.Setenv("FLAI_CACHE_DIR", cache)
	var out, errOut bytes.Buffer
	tty := reply != ""
	a := &app{out: &out, errOut: &errOut, cwd: dir, stdinIsTerminal: &tty,
		clock:    func() time.Time { return time.Date(2026, 10, 6, 21, 0, 0, 0, time.UTC) },
		prompter: prompt.New(strings.NewReader(reply), &out)}
	root := newRootCmdWith(a)
	root.SetArgs(args)
	code := 0
	if err := root.Execute(); err != nil {
		code = 1
		if ee := (*exitError)(nil); errors.As(err, &ee) {
			code = ee.code
		}
		a.fail(err)
	}
	return out.String(), errOut.String(), code
}

// With no --ref, a project whose template.ref is main goes to the newest
// version tag, not to the head of main, and records the tag; the tag it
// recorded is no pin, so the next release is taken the same way (ADR-0103).
func TestUpgradeTakesTheNewestRelease(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	_, url, commitVersion := gitTemplate(t)
	dest := newFromGit(t, url, "main")
	commitVersion("10.0.0", "v10.0.0")
	commitVersion("11.0.0", "v11.0.0")
	commitVersion("12.0.0", "")
	cache := t.TempDir()

	before := upgradeState(t, dest)
	out, errOut, code := upgradeIn(t, dest, cache, "", "upgrade", "--dry-run")
	if code != 0 || !strings.Contains(out, "template v11.0.0 at 11.0.0: "+reasonNewest) || !strings.Contains(out, "9.9.9 -> 11.0.0") {
		t.Fatalf("dry run: %d %s %s", code, out, errOut)
	}
	if upgradeState(t, dest) != before {
		t.Error("a dry run changed the project")
	}
	out, errOut, code = upgradeIn(t, dest, cache, "", "upgrade")
	if code != 0 || !strings.Contains(out, "upgraded to template 11.0.0") {
		t.Fatalf("upgrade: %d %s %s", code, out, errOut)
	}
	if got := readTeamTxt(t, dest); got != "v11.0.0\n" {
		t.Errorf("team.txt = %q, want the newest release's", got)
	}
	recordedTemplate(t, dest, url, "v11.0.0", "11.0.0")
	if out, _, code := upgradeIn(t, dest, cache, "", "upgrade"); code != 0 || !strings.Contains(out, "already at template 11.0.0") {
		t.Errorf("at the newest release: %d %s", code, out)
	}

	commitVersion("13.0.0", "v13.0.0")
	out, errOut, code = upgradeIn(t, dest, cache, "", "upgrade", "--json")
	if code != 0 {
		t.Fatalf("upgrade to the next release: %d %s %s", code, out, errOut)
	}
	var res struct{ Target upgradeTarget }
	if err := json.Unmarshal([]byte(out), &res); err != nil || res.Target != (upgradeTarget{Ref: "v13.0.0", Version: "13.0.0", Reason: reasonNewest}) {
		t.Errorf("--json target: %+v %v\n%s", res.Target, err, out)
	}
	recordedTemplate(t, dest, url, "v13.0.0", "13.0.0")
}

// The reported case: the lock records main at 9.9.9 and system-flow.yaml's
// version was set by hand to the newest release. The upgrade applies it,
// rather than saying it is already there or going back to the lock's.
func TestUpgradeAppliesAVersionSetByHand(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	_, url, commitVersion := gitTemplate(t)
	dest := newFromGit(t, url, "main")
	commitVersion("10.0.0", "v10.0.0")
	commitVersion("11.0.0", "v11.0.0")
	editTemplateField(t, dest, "version", `"11.0.0"`)

	out, errOut, code := upgradeIn(t, dest, t.TempDir(), "", "upgrade")
	if code != 0 || strings.Contains(out, "already at") || !strings.Contains(out, "upgraded to template 11.0.0") || !strings.Contains(out, reasonManifest+", the newest release") {
		t.Fatalf("upgrade to the version set by hand: %d %s %s", code, out, errOut)
	}
	if got := readTeamTxt(t, dest); got != "v11.0.0\n" {
		t.Errorf("team.txt = %q, want 11.0.0's", got)
	}
	recordedTemplate(t, dest, url, "v11.0.0", "11.0.0")
}

// system-flow.yaml set by hand to an older release than the newest: without
// a terminal, or with --yes, the run names both and changes nothing; a dry
// run says it would ask; on a terminal each answer is applied and recorded.
func TestUpgradeAsksWhenSystemFlowYamlNamesAnOlderRelease(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	_, url, commitVersion := gitTemplate(t)
	projects := map[string]string{}
	for _, name := range []string{"refused", "named", "newest", "nothing"} {
		projects[name] = newFromGit(t, url, "main")
	}
	commitVersion("10.0.0", "v10.0.0")
	commitVersion("11.0.0", "v11.0.0")
	for _, dest := range projects {
		editTemplateField(t, dest, "version", `"10.0.0"`)
	}
	cache := t.TempDir()

	dest := projects["refused"]
	before := upgradeState(t, dest)
	for _, args := range [][]string{{"upgrade"}, {"upgrade", "--yes"}} {
		_, errOut, code := upgradeIn(t, dest, cache, "", args...)
		if code != 1 || !strings.Contains(errOut, "flai upgrade --ref v10.0.0") || !strings.Contains(errOut, "flai upgrade --ref v11.0.0") || !strings.Contains(errOut, "nothing was changed") {
			t.Errorf("%v without a choice: %d %s", args, code, errOut)
		}
	}
	out, errOut, code := upgradeIn(t, dest, cache, "", "upgrade", "--dry-run")
	if code != 0 || !strings.Contains(out, "names template v10.0.0 (10.0.0), but the newest release is v11.0.0 (11.0.0)") || !strings.Contains(out, "would ask") {
		t.Errorf("dry run: %d %s %s", code, out, errOut)
	}
	if upgradeState(t, dest) != before {
		t.Error("a refused upgrade changed the project")
	}

	out, errOut, code = upgradeIn(t, projects["named"], cache, "1\n", "upgrade")
	if code != 0 || !strings.Contains(out, "Which version") || !strings.Contains(out, "template v10.0.0 at 10.0.0: "+reasonChosen) || !strings.Contains(out, "upgraded to template 10.0.0") {
		t.Fatalf("choose the version system-flow.yaml names: %d %s %s", code, out, errOut)
	}
	recordedTemplate(t, projects["named"], url, "v10.0.0", "10.0.0")
	// The tag it recorded is no pin: the next upgrade takes the newest, unasked.
	if out, errOut, code := upgradeIn(t, projects["named"], cache, "", "upgrade"); code != 0 || !strings.Contains(out, "upgraded to template 11.0.0") {
		t.Errorf("after the choice, the next upgrade: %d %s %s", code, out, errOut)
	}
	recordedTemplate(t, projects["named"], url, "v11.0.0", "11.0.0")

	out, errOut, code = upgradeIn(t, projects["newest"], cache, "2\n", "upgrade")
	if code != 0 || !strings.Contains(out, "upgraded to template 11.0.0") {
		t.Fatalf("choose the newest: %d %s %s", code, out, errOut)
	}
	recordedTemplate(t, projects["newest"], url, "v11.0.0", "11.0.0")
	// Choosing the newest when already at it records that, so it is not asked again.
	editTemplateField(t, projects["newest"], "version", `"10.0.0"`)
	out, errOut, code = upgradeIn(t, projects["newest"], cache, "2\n", "upgrade")
	if code != 0 || !strings.Contains(out, "already at template 11.0.0") || !strings.Contains(out, "now record template v11.0.0 at 11.0.0") {
		t.Fatalf("choose the newest when already at it: %d %s %s", code, out, errOut)
	}
	recordedTemplate(t, projects["newest"], url, "v11.0.0", "11.0.0")
	if out, errOut, code := upgradeIn(t, projects["newest"], cache, "", "upgrade"); code != 0 || !strings.Contains(out, "already at template 11.0.0") {
		t.Errorf("after settling the edit: %d %s %s", code, out, errOut)
	}

	before = upgradeState(t, projects["nothing"])
	out, errOut, code = upgradeIn(t, projects["nothing"], cache, "3\n", "upgrade")
	if code != 0 || !strings.Contains(out, "nothing changed: the project stays at template 9.9.9") {
		t.Fatalf("choose to change nothing: %d %s %s", code, out, errOut)
	}
	if upgradeState(t, projects["nothing"]) != before {
		t.Error("changing nothing changed the project")
	}
}

// A template.ref naming a branch other than the default does not follow
// releases: it is used as given, and a commit on it after the cache cloned
// it is applied.
func TestUpgradeFollowsABranchThatIsNotTheDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	tpl, url, commitVersion := gitTemplate(t)
	gitIn(t, tpl, "branch", "develop")
	dest := newFromGit(t, url, "develop")
	commitVersion("10.0.0", "v10.0.0")
	cache := t.TempDir()

	out, errOut, code := upgradeIn(t, dest, cache, "", "upgrade")
	if code != 0 || !strings.Contains(out, "template develop at 9.9.9: "+reasonAsGiven) || !strings.Contains(out, "already at template 9.9.9") {
		t.Fatalf("upgrade on develop: %d %s %s", code, out, errOut)
	}
	gitIn(t, tpl, "checkout", "-q", "develop")
	commitVersion("9.9.10", "")
	gitIn(t, tpl, "checkout", "-q", "main")
	out, errOut, code = upgradeIn(t, dest, cache, "", "upgrade")
	if code != 0 || !strings.Contains(out, "upgraded to template 9.9.10") {
		t.Fatalf("upgrade after a commit on develop: %d %s %s", code, out, errOut)
	}
	if got := readTeamTxt(t, dest); got != "v9.9.10\n" {
		t.Errorf("team.txt = %q, want develop's new commit", got)
	}
	recordedTemplate(t, dest, url, "develop", "9.9.10")
}

// lsRemote answers git ls-remote with its own text, so a template.Remote is
// read without git.
type lsRemote string

func (l lsRemote) Run(string, string, ...string) (string, error)              { return string(l), nil }
func (l lsRemote) RunInput(string, string, string, ...string) (string, error) { return string(l), nil }
func (lsRemote) LookPath(name string) (string, error)                         { return name, nil }

// remoteWith is a template remote whose default branch is main, with a
// branch develop and the given tags.
func remoteWith(t *testing.T, tags ...string) template.Remote {
	t.Helper()
	out := "ref: refs/heads/main\tHEAD\nabc\tHEAD\nabc\trefs/heads/main\nabc\trefs/heads/develop\n"
	for _, tag := range tags {
		out += "abc\trefs/tags/" + tag + "\n"
	}
	rm, err := template.ListRemote(lsRemote(out), "https://example.org/tpl.git")
	if err != nil {
		t.Fatal(err)
	}
	return rm
}

// The decision of which ref an upgrade with no --ref applies (ADR-0103).
func TestChooseTarget(t *testing.T) {
	rm := remoteWith(t, "v1.0.18", "v1.0.50", "v1.0.60")
	at := func(ref, version string) manifest.Template { return manifest.Template{Ref: ref, Version: version} }
	lockAt := func(ref, version string) *lock.Template { return &lock.Template{Ref: ref, Version: version} }
	newest := upgradeTarget{Ref: "v1.0.60", Version: "1.0.60", Reason: reasonNewest}
	for _, c := range []struct {
		name     string
		rm       template.Remote
		mf       manifest.Template
		lk       *lock.Template
		want     upgradeTarget
		conflict string // the release the edit names, when it asks
		err      string
	}{
		{name: "main follows releases", rm: rm, mf: at("main", "1.0.18"), lk: lockAt("main", "1.0.18"), want: newest},
		{name: "no lock", rm: rm, mf: at("main", "1.0.18"), want: newest},
		{name: "no ref", rm: rm, mf: at("", "1.0.18"), lk: lockAt("", "1.0.18"), want: newest},
		{name: "a recorded tag equal to the lock's is no pin", rm: rm, mf: at("v1.0.50", "1.0.50"), lk: lockAt("v1.0.50", "1.0.50"), want: newest},
		{name: "version edited to the newest", rm: rm, mf: at("main", "1.0.60"), lk: lockAt("main", "1.0.18"),
			want: upgradeTarget{Ref: "v1.0.60", Version: "1.0.60", Reason: reasonManifest + ", the newest release"}},
		{name: "ref edited to the newest, without its v", rm: rm, mf: at("1.0.60", "1.0.18"), lk: lockAt("main", "1.0.18"),
			want: upgradeTarget{Ref: "v1.0.60", Version: "1.0.60", Reason: reasonManifest + ", the newest release"}},
		{name: "version edited to an older tag", rm: rm, mf: at("main", "1.0.50"), lk: lockAt("main", "1.0.18"), conflict: "v1.0.50"},
		{name: "ref edited to an older tag", rm: rm, mf: at("v1.0.50", "1.0.18"), lk: lockAt("main", "1.0.18"), conflict: "v1.0.50"},
		{name: "ref edited back to main", rm: rm, mf: at("main", "1.0.50"), lk: lockAt("v1.0.50", "1.0.50"), want: newest},
		{name: "version edited with no tag", rm: rm, mf: at("main", "1.0.55"), lk: lockAt("main", "1.0.18"), err: "no release tag for it; its releases are v1.0.18, v1.0.50, v1.0.60"},
		{name: "another branch", rm: rm, mf: at("develop", "1.0.18"), lk: lockAt("develop", "1.0.18"), want: upgradeTarget{Ref: "develop", Reason: reasonAsGiven}},
		{name: "ref edited to another branch", rm: rm, mf: at("develop", "1.0.18"), lk: lockAt("main", "1.0.18"), want: upgradeTarget{Ref: "develop", Reason: reasonAsGiven}},
		{name: "a commit", rm: rm, mf: at("0123abc", "1.0.18"), lk: lockAt("0123abc", "1.0.18"), want: upgradeTarget{Ref: "0123abc", Reason: reasonAsGiven}},
		{name: "no version tags", rm: remoteWith(t, "nightly"), mf: at("main", "1.0.18"), lk: lockAt("main", "1.0.18"), want: upgradeTarget{Ref: "main", Reason: reasonNoTags}},
	} {
		got, conflict, err := chooseTarget(c.rm, c.mf, c.lk)
		switch {
		case c.err != "":
			if err == nil || !strings.Contains(err.Error(), c.err) {
				t.Errorf("%s: err %v, want %q", c.name, err, c.err)
			}
		case err != nil:
			t.Errorf("%s: %v", c.name, err)
		case c.conflict != "":
			if conflict == nil || conflict.Edited.Ref != c.conflict || conflict.Newest != newest {
				t.Errorf("%s: conflict %+v, want %s against the newest", c.name, conflict, c.conflict)
			}
		case conflict != nil || got != c.want:
			t.Errorf("%s: %+v %+v, want %+v", c.name, got, conflict, c.want)
		}
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

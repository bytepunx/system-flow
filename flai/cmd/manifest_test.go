package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// manifestProject is a scratch project whose manifest has comments and a
// strategic block, and its manifest's contents.
func manifestProject(t *testing.T) (string, string) {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	file := filepath.Join(root, manifest.File)
	data, _ := os.ReadFile(file)
	data = append(data, "# what the orchestrator may do\norchestration:\n  policy: fifo   # the operator's order\n  permissions:\n    promote_to_ready: true\n"...)
	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return root, string(data)
}

func readManifest(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0229: flai manifest set writes keys of each kind into their blocks,
// --unset removes one, and the rest of the file is as it was.
func TestManifestSetWritesAndUnsets(t *testing.T) {
	root, before := manifestProject(t)
	out, errOut, code := runIn(t, root, "manifest", "set", "orchestration.policy=wsjf", "planning.schedule=0 6 * * 1-5", "orchestration.permissions.publish=true", "--unset", "orchestration.permissions.promote_to_ready")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	want := strings.Replace(before, "  policy: fifo   # the operator's order\n  permissions:\n    promote_to_ready: true\n",
		"  policy: wsjf   # the operator's order\n  permissions:\n    publish: true\n", 1) + "planning:\n  schedule: \"0 6 * * 1-5\"\n"
	if got := readManifest(t, root); got != want {
		t.Errorf("manifest:\n%s\nwant:\n%s", got, want)
	}
	if !strings.Contains(out, "set orchestration.policy = wsjf\n") || !strings.Contains(out, "unset orchestration.permissions.promote_to_ready\n") {
		t.Errorf("output: %q", out)
	}
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	if m.Orchestration.Permissions.Allows(manifest.PermitPromoteToReady) || !m.Orchestration.Permissions.Allows(manifest.PermitPublish) {
		t.Errorf("permissions read back: %+v", m.Orchestration.Permissions)
	}
	out, errOut, code = runIn(t, root, "manifest", "set", "--unset", "planning.schedule", "--json")
	var got struct {
		Set   []manifest.Assignment `json:"set"`
		Unset []string              `json:"unset"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil || len(got.Set) != 0 || len(got.Unset) != 1 || got.Unset[0] != "planning.schedule" {
		t.Errorf("unset --json: %d %q %s", code, out, errOut)
	}
	if strings.Contains(readManifest(t, root), "schedule:") {
		t.Error("planning.schedule is still there")
	}
}

// S-0229: a bad value and a key outside the catalog are refused with the
// field and the reason, as a rule is, and nothing is written.
func TestManifestSetRefuses(t *testing.T) {
	root, before := manifestProject(t)
	out, errOut, code := runIn(t, root, "manifest", "set", "planning.cycle")
	if code != 1 || !strings.Contains(out, "planning.cycle: is not key=value") || !strings.Contains(errOut, "rule: ") {
		t.Errorf("an argument that is not key=value is a rule's refusal: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "manifest", "set", "planning.currency=GBP", "planning.cycle=2h")
	if code != 1 || !strings.Contains(out, "planning.currency: is not a setting flai writes") || !strings.Contains(errOut, "rule: planning.currency") {
		t.Errorf("planning.currency is refused as a rule: %d %q %s", code, out, errOut)
	}
	out, _, code = runIn(t, root, "manifest", "set", "orchestration.policy=fastest", "orchestration.release.value=-1", "planning.colour=blue", "--json")
	var got struct {
		Refused []manifest.Problem `json:"refused"`
	}
	if code != 1 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("--json: %d %q", code, out)
	}
	fields := map[string]string{}
	for _, p := range got.Refused {
		fields[p.Field] = p.Reason
	}
	if len(got.Refused) != 1 || !strings.Contains(fields["planning.colour"], "is not a setting") {
		t.Errorf("the unknown key is refused first, alone: %+v", got.Refused)
	}
	out, _, code = runIn(t, root, "manifest", "set", "orchestration.policy=fastest", "orchestration.release.policy=threshold", "orchestration.release.value=-1", "--json")
	got.Refused = nil
	if code != 1 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("--json: %d %q", code, out)
	}
	fields = map[string]string{}
	for _, p := range got.Refused {
		fields[p.Field] = p.Reason
	}
	if !strings.Contains(fields["orchestration.policy"], "is not an order policy") || !strings.Contains(fields["orchestration.release.value"], "is not an amount of zero or more") {
		t.Errorf("refusals: %+v", got.Refused)
	}
	if after := readManifest(t, root); after != before {
		t.Errorf("the manifest changed:\n%s", after)
	}
}

// S-0229: with --autocommit the change is committed, system-flow.yaml alone,
// with its trailers.
func TestManifestSetAutocommitsTheManifestAlone(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, _ := manifestProject(t)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "start")
	_ = os.WriteFile(filepath.Join(root, "notes.txt"), []byte("not mine\n"), 0o644)
	t.Setenv("GIT_AUTHOR_NAME", "t")
	t.Setenv("GIT_AUTHOR_EMAIL", "t@t")
	t.Setenv("GIT_COMMITTER_NAME", "t")
	t.Setenv("GIT_COMMITTER_EMAIL", "t@t")

	out, errOut, code := runIn(t, root, "manifest", "set", "planning.replan=agent", "--autocommit", "--trailer", "Co-Authored-By: someone <s@s>", "--json")
	if code != 0 || !strings.Contains(out, `"commit"`) {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	if st := gitIn(t, root, "status", "--porcelain", "--", manifest.File); st != "" {
		t.Errorf("the manifest is committed: %q", st)
	}
	if st := gitIn(t, root, "status", "--porcelain", "--", "notes.txt"); st == "" {
		t.Error("another file was committed with it")
	}
	msg := gitIn(t, root, "log", "-1", "--format=%B")
	if !strings.Contains(msg, "planning.replan") || !strings.Contains(msg, "Co-Authored-By: someone <s@s>") {
		t.Errorf("commit message: %q", msg)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); strings.TrimSpace(files) != manifest.File {
		t.Errorf("the commit holds the manifest alone: %q", files)
	}
}

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
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
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
	t.Parallel()
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
	t.Parallel()
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

// S-0692 (I-0292): flai manifest set writes analysis.formula naming a formula
// a file defines, every input with a default, and refuses, writing nothing, a
// name no file defines and a formula with an input that has no default, with
// the reason flai check gives as manifest.analysis; --unset removes it.
func TestManifestSetWritesAnalysisFormula(t *testing.T) {
	t.Parallel()
	root, before := manifestProject(t)
	dir := filepath.Join(root, "design", "workflows", "formulas")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	formula := func(name, input string) string {
		return "---\ntitle: " + name + "\nname: " + name + "\ndescription: d\ninputs:\n  - " + input +
			"\nsteps:\n  - id: list\n    title: List\n    run: {flai: issue list --json}\n---\n\n# " + name + "\n"
	}
	for name, input := range map[string]string{
		"analysis-cadence": `{name: focus, description: f, default: ""}`,
		"needs-a-focus":    `{name: focus, description: f}`,
	} {
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(formula(name, input)), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	refused := func(value, says string) {
		t.Helper()
		out, _, code := runIn(t, root, "manifest", "set", "analysis.formula="+value, "analysis.schedule=0 6 * * 1", "--json")
		var got struct {
			Refused []manifest.Problem `json:"refused"`
		}
		if code != 1 || json.Unmarshal([]byte(out), &got) != nil || len(got.Refused) != 1 ||
			got.Refused[0].Field != "analysis.formula" || !strings.Contains(got.Refused[0].Reason, says) {
			t.Errorf("analysis.formula=%s: %d %s, want refused saying %q", value, code, out, says)
		}
		if after := readManifest(t, root); after != before {
			t.Errorf("analysis.formula=%s changed the manifest:\n%s", value, after)
		}
	}
	refused("nope", "no formula is named nope; the formulas are analysis-cadence, needs-a-focus")
	refused("needs-a-focus", "names needs-a-focus, whose input focus has no default")
	refused("Not Kebab", "is not a formula's name")

	out, errOut, code := runIn(t, root, "manifest", "set", "analysis.formula=analysis-cadence", "analysis.schedule=0 6 * * 1")
	if code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	if want := before + "analysis:\n  formula: analysis-cadence\n  schedule: \"0 6 * * 1\"\n"; readManifest(t, root) != want {
		t.Errorf("manifest:\n%s\nwant:\n%s", readManifest(t, root), want)
	}
	if out, errOut, code := runIn(t, root, "manifest", "set", "--unset", "analysis.formula"); code != 0 || strings.Contains(readManifest(t, root), "formula:") {
		t.Errorf("--unset analysis.formula: %d %s %s\n%s", code, out, errOut, readManifest(t, root))
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

// S-0839 (I-0364): from a manifest whose agent image names a layer and no
// args, as S-0760 found this repository's, flai manifest set writes build
// arguments under agent.image.args, keeping the comments, and the manifest
// then has no manifest.agent problem; --unset removes one argument; BASE is
// refused with the field and the reason, writing nothing; and the help names
// both keys.
func TestManifestSetWritesAgentImageArgs(t *testing.T) {
	t.Parallel()
	root, _ := manifestProject(t)
	file := filepath.Join(root, manifest.File)
	before := readManifest(t, root) + "agent:\n  harness: claude-code\n  image:\n    layer: images/agent/Dockerfile   # the project's tools\n"
	if err := os.WriteFile(file, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	load := func() *manifest.AgentImage {
		t.Helper()
		m, err := manifest.Load(file)
		if err != nil {
			t.Fatal(err)
		}
		if err := m.Agent.Validate(); err != nil {
			t.Errorf("manifest.agent: %v\n%s", err, readManifest(t, root))
		}
		if m.Agent == nil || m.Agent.Image == nil {
			t.Fatalf("no agent image read back:\n%s", readManifest(t, root))
		}
		return m.Agent.Image
	}

	out, errOut, code := runIn(t, root, "manifest", "set", "agent.image.args.GO_VERSION=1.26.0", "agent.image.args.PNPM_VERSION=10.34.5", "--json")
	var got struct {
		Set   []manifest.Assignment `json:"set"`
		Unset []string              `json:"unset"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("exit %d: %s %s", code, out, errOut)
	}
	if len(got.Set) != 2 || got.Set[0].Key != "agent.image.args.GO_VERSION" || got.Set[1].Key != "agent.image.args.PNPM_VERSION" || len(got.Unset) != 0 {
		t.Errorf("--json lists the keys set: %s", out)
	}
	image := load()
	if image.Layer != "images/agent/Dockerfile" || len(image.Args) != 2 || image.Args["GO_VERSION"] != "1.26.0" || image.Args["PNPM_VERSION"] != "10.34.5" {
		t.Errorf("image read back: %+v", image)
	}
	if after := readManifest(t, root); !strings.Contains(after, "# the project's tools") || !strings.Contains(after, "# what the orchestrator may do") {
		t.Errorf("the comments are gone:\n%s", after)
	}

	if out, errOut, code := runIn(t, root, "manifest", "set", "--unset", "agent.image.args.GO_VERSION"); code != 0 || !strings.Contains(out, "unset agent.image.args.GO_VERSION\n") {
		t.Fatalf("--unset: %d %s %s", code, out, errOut)
	}
	if image := load(); len(image.Args) != 1 || image.Args["PNPM_VERSION"] != "10.34.5" || strings.Contains(readManifest(t, root), "GO_VERSION") {
		t.Errorf("--unset removes the one argument: %+v\n%s", image, readManifest(t, root))
	}

	kept := readManifest(t, root)
	out, errOut, code = runIn(t, root, "manifest", "set", "agent.image.args.BASE=x")
	if code != 1 || !strings.Contains(out, "agent: ") || !strings.Contains(out, "BASE") || !strings.Contains(errOut, "rule: ") {
		t.Errorf("BASE is refused with the field and the reason: %d %q %s", code, out, errOut)
	}
	if after := readManifest(t, root); after != kept {
		t.Errorf("the refusal changed the manifest:\n%s", after)
	}

	out, _, _ = runIn(t, root, "manifest", "set", "--help")
	if !strings.Contains(out, "agent.image.layer") || !strings.Contains(out, "agent.image.args.<NAME>") || !strings.Contains(out, "flai manifest set agent.image.layer=images/agent/Dockerfile agent.image.args.GO_VERSION=1.26.0") {
		t.Errorf("the help names the agent image's keys:\n%s", out)
	}
}

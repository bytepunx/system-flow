package kluster

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// The commits of graphGit's history, oldest first: c1 and c2 change code,
// c3 only the design, docs, and wip folders, and c4 code again.
var (
	evC1 = strings.Repeat("1", 40)
	evC2 = strings.Repeat("2", 40)
	evC3 = strings.Repeat("3", 40)
	evC4 = strings.Repeat("4", 40)
)

// graphGit is an execx.Runner that answers git from a linear history:
// rev-parse resolves a story's branch from heads, merge-base --is-ancestor
// walks parent, and diff --name-only -z lists the files each commit after
// the first up to the second changed.
type graphGit struct {
	heads   map[string]string
	parent  map[string]string
	files   map[string][]string
	diffErr error
	calls   [][]string
}

func newGraphGit(head string) *graphGit {
	return &graphGit{
		heads:  map[string]string{"S-0662": head},
		parent: map[string]string{evC2: evC1, evC3: evC2, evC4: evC3},
		files: map[string][]string{
			evC1: {"flai/main.go"},
			evC2: {"flai/internal/kluster/evidence.go"},
			evC3: {"docs/users/flai.md", "wip/agents/S-0662-x.md", "design/issues/I-0001.md"},
			evC4: {"flai/internal/kluster/deploy.go"},
		},
	}
}

func (g *graphGit) Run(_, _ string, args ...string) (string, error) {
	g.calls = append(g.calls, args)
	switch args[0] {
	case "rev-parse":
		story := strings.TrimSuffix(strings.TrimPrefix(args[len(args)-1], "refs/heads/story/"), "^{commit}")
		if head, ok := g.heads[story]; ok {
			return head + "\n", nil
		}
		return "", errors.New("exit status 1")
	case "merge-base":
		for c := args[3]; c != ""; c = g.parent[c] {
			if c == args[2] {
				return "", nil
			}
		}
		return "", errors.New("exit status 1")
	case "diff":
		if g.diffErr != nil {
			return "", g.diffErr
		}
		from, to := args[len(args)-2], args[len(args)-1]
		var out []string
		for c := to; c != "" && c != from; c = g.parent[c] {
			out = append(out, g.files[c]...)
		}
		if len(out) == 0 {
			return "", nil
		}
		return strings.Join(out, "\x00") + "\x00", nil
	}
	return "", errors.New("unexpected git " + strings.Join(args, " "))
}

func (g *graphGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (g *graphGit) LookPath(name string) (string, error) { return name, nil }

// evManifest is a manifest with the default layout and a kluster: block
// that names a deploy.
func evManifest() manifest.Manifest {
	return manifest.Manifest{
		Layout:  map[string]string{"design": "design", "docs": "docs", "wip": "wip"},
		Kluster: &manifest.Kluster{Deploy: []string{"scripts/deploy.sh"}},
	}
}

// evRecord writes s-0662's record: a deploy at commit that ended outcome,
// with a checks run at checks that ended checksOutcome when checks is set.
func evRecord(t *testing.T, root, commit, outcome, checks, checksOutcome string) Record {
	t.Helper()
	r := Record{Story: "S-0662", Namespace: "s-0662", Commit: commit, Outcome: outcome, AskedBy: "agent-S-0662",
		Steps: []Step{
			{Name: StepBuild, Outcome: OutcomePassed, Log: LogPath(root, "s-0662", StepBuild)},
			{Name: StepLoad, Outcome: OutcomePassed, Log: LogPath(root, "s-0662", StepLoad)},
			{Name: StepDeploy, Outcome: outcome, Log: LogPath(root, "s-0662", StepDeploy)},
		}}
	if outcome == OutcomeFailed {
		r.FailedStep, r.ExitCode, r.Error = StepDeploy, 2, "deploy: exit status 2"
	}
	if checks != "" {
		r.ChecksCommit, r.ChecksOutcome, r.ChecksLog = checks, checksOutcome, "/logs/checks.log"
	}
	if err := writeRecord(root, r); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestEvidenceDoesNotApplyWithoutABlockTheActionOrADeploy(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC2, OutcomePassed, evC2, OutcomePassed)
	noDeploy := evManifest()
	noDeploy.Kluster = &manifest.Kluster{}
	for _, c := range []struct {
		name string
		m    manifest.Manifest
		on   bool
		why  string
	}{
		{"no block", manifest.Manifest{}, true, "no kluster: block"},
		{"action off", evManifest(), false, "kluster host action is off"},
		{"no deploy", noDeploy, true, "names no deploy"},
	} {
		e, err := ReadEvidence(newGraphGit(evC2), root, c.m, "s-0662", c.on, true)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if e.Applies || e.Passed || e.Deploy != nil || e.Commands != nil || !strings.Contains(e.Why, c.why) || e.Story != "S-0662" {
			t.Errorf("%s: %+v", c.name, e)
		}
		if s := e.Summary(); !strings.Contains(s, "no kluster evidence applies to S-0662") || !strings.Contains(s, c.why) {
			t.Errorf("%s: summary %q", c.name, s)
		}
		if n := e.Next(); n != "" {
			t.Errorf("%s: next %q", c.name, n)
		}
	}
}

func TestEvidenceWithNoRecordNamesTheDeploy(t *testing.T) {
	e, err := ReadEvidence(newGraphGit(evC2), t.TempDir(), evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !e.Applies || e.Passed || e.Deploy != nil || e.Checks != nil || e.Head != evC2 || e.Namespace != "s-0662" {
		t.Errorf("no record: %+v", e)
	}
	if s := e.Summary(); !strings.Contains(s, "no deploy of it is recorded") {
		t.Errorf("summary %q", s)
	}
	if n := e.Next(); !strings.Contains(n, "flai kluster deploy S-0662 --job") || !strings.Contains(n, "kluster_deploy") {
		t.Errorf("next %q", n)
	}
	want := EvidenceCommands{Deploy: "flai kluster deploy S-0662 --job", DeployTool: "kluster_deploy", Checks: "flai checks run S-0662 --job", ChecksTool: "kluster_checks"}
	if e.Commands == nil || *e.Commands != want {
		t.Errorf("commands %+v", e.Commands)
	}
}

func TestEvidenceOfADeployAtAnOlderCommitIsNotAtTheHead(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC1, OutcomePassed, evC1, OutcomePassed)
	e, err := ReadEvidence(newGraphGit(evC2), root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	d := e.Deploy
	if e.Passed || d == nil || d.PassedAtHead || d.Head.AtHead || d.Head.Commit != evC1 || d.Head.Head != evC2 || !slices.Equal(d.Head.Outside, []string{"flai/internal/kluster/evidence.go"}) {
		t.Fatalf("an older deploy: %+v", d)
	}
	if e.Checks == nil || e.Checks.PassedAtHead {
		t.Errorf("checks at the older commit: %+v", e.Checks)
	}
	if s := e.Summary(); !strings.Contains(s, "not deployed at the head 222222222222") || !strings.Contains(s, "at 111111111111, passed") || !strings.Contains(s, "flai/internal/kluster/evidence.go") {
		t.Errorf("summary %q", s)
	}
	if n := e.Next(); !strings.Contains(n, "kluster deploy") {
		t.Errorf("next %q", n)
	}
}

func TestEvidenceOfAFailedDeploy(t *testing.T) {
	root := t.TempDir()
	r := evRecord(t, root, evC2, OutcomeFailed, "", "")
	e, err := ReadEvidence(newGraphGit(evC2), root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	d := e.Deploy
	if e.Passed || d == nil || d.PassedAtHead || !d.Head.AtHead || !d.Head.Exact || d.FailedStep != StepDeploy || d.ExitCode != 2 || d.AskedBy != "agent-S-0662" || d.Outcome != OutcomeFailed {
		t.Fatalf("a failed deploy: %+v", d)
	}
	if len(d.Steps) != 3 || d.Steps[2].Log != r.Steps[2].Log {
		t.Errorf("steps %+v", d.Steps)
	}
	if s := e.Summary(); !strings.Contains(s, "failed at the deploy step") {
		t.Errorf("summary %q", s)
	}
}

func TestEvidenceOfAPassedDeployWithNoChecksRun(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC2, OutcomePassed, "", "")
	e, err := ReadEvidence(newGraphGit(evC2), root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if e.Passed || e.Deploy == nil || !e.Deploy.PassedAtHead || e.Checks != nil {
		t.Fatalf("checks on: %+v", e)
	}
	if s := e.Summary(); !strings.Contains(s, "no checks run against it has ended there") {
		t.Errorf("summary %q", s)
	}
	if n := e.Next(); !strings.Contains(n, "flai checks run S-0662 --job") || !strings.Contains(n, "kluster_checks") {
		t.Errorf("next %q", n)
	}
	off, err := ReadEvidence(newGraphGit(evC2), root, evManifest(), "S-0662", true, false)
	if err != nil {
		t.Fatal(err)
	}
	if !off.Passed || off.Next() != "" || !strings.Contains(off.Summary(), "checks host action is off") {
		t.Errorf("checks off: %+v, %q", off, off.Summary())
	}
}

func TestEvidenceOfChecksFailedAtTheHead(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC2, OutcomePassed, evC2, OutcomeFailed)
	e, err := ReadEvidence(newGraphGit(evC2), root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	c := e.Checks
	if e.Passed || !e.Deploy.PassedAtHead || c == nil || c.PassedAtHead || !c.Head.AtHead || c.Outcome != OutcomeFailed || c.Log != "/logs/checks.log" {
		t.Fatalf("failed checks: %+v", c)
	}
	if s := e.Summary(); !strings.Contains(s, "its checks run there ended failed") {
		t.Errorf("summary %q", s)
	}
	if n := e.Next(); !strings.Contains(n, "flai checks run") {
		t.Errorf("next %q", n)
	}
}

func TestEvidenceOfBothPassedAtTheHead(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC2, OutcomePassed, evC2, OutcomePassed)
	g := newGraphGit(evC2)
	e, err := ReadEvidence(g, root, evManifest(), "s-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !e.Passed || !e.Deploy.Head.Exact || !e.Checks.PassedAtHead || e.Next() != "" {
		t.Fatalf("both passed: %+v", e)
	}
	if s, want := e.Summary(), "S-0662's namespace s-0662 is deployed and checked at the head 222222222222"; s != want {
		t.Errorf("summary %q, want %q", s, want)
	}
	if len(g.calls) != 1 {
		t.Errorf("git ran %q; the head itself needs only rev-parse", g.calls)
	}
}

func TestEvidenceAtACommitFollowedOnlyByDesignDocsAndWipIsAtTheHead(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC2, OutcomePassed, evC2, OutcomePassed)
	e, err := ReadEvidence(newGraphGit(evC3), root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !e.Passed || !e.Deploy.Head.AtHead || e.Deploy.Head.Exact || e.Deploy.Head.Commit != evC2 || e.Head != evC3 || !e.Checks.PassedAtHead {
		t.Fatalf("docs only since: %+v", e.Deploy.Head)
	}
	if s := e.Summary(); !strings.Contains(s, "deployed and checked at the head 333333333333 (at 222222222222, followed only by") {
		t.Errorf("summary %q", s)
	}

	code, err := ReadEvidence(newGraphGit(evC4), root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if code.Passed || code.Deploy.Head.AtHead || code.Checks.PassedAtHead || !slices.Equal(code.Deploy.Head.Outside, []string{"flai/internal/kluster/deploy.go"}) {
		t.Fatalf("code since: %+v", code.Deploy.Head)
	}
	if s := code.Summary(); !strings.Contains(s, "not deployed at the head 444444444444") {
		t.Errorf("summary %q", s)
	}
}

func TestEvidenceWithNoBranchHasNoHead(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC2, OutcomePassed, evC2, OutcomePassed)
	g := newGraphGit(evC2)
	g.heads = nil
	e, err := ReadEvidence(g, root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	if e.Passed || e.Head != "" || e.Deploy.Head.AtHead || e.Next() != "" || !strings.Contains(e.Summary(), "has no branch story/S-0662") {
		t.Errorf("no branch: %+v, %q", e, e.Summary())
	}
}

func TestEvidenceRefusesANameThatIsNoStory(t *testing.T) {
	if _, err := ReadEvidence(newGraphGit(evC2), t.TempDir(), evManifest(), "main", true, true); err == nil {
		t.Error("a name that is no story ID gave evidence")
	}
}

func TestAtHeadJudgesRebasesLayoutsAndGitFailures(t *testing.T) {
	g := newGraphGit(evC4)
	g.parent[evC2] = "" // c1 is no longer in the history of the head
	j, err := AtHead(g, "/main", evManifest(), evC1, evC4)
	if err != nil || j.AtHead || !strings.Contains(j.Why, "rebased or reset") {
		t.Errorf("a rebased commit: %+v %v", j, err)
	}
	g = newGraphGit(evC3)
	g.files[evC3] = []string{"documentation/guide.md", "wip/agents/x.md", "design/a.md"}
	m := evManifest()
	m.Layout["docs"] = "./documentation/"
	if j, err := AtHead(g, "/main", m, evC2, evC3); err != nil || !j.AtHead || j.Exact {
		t.Errorf("a docs folder of the layout's: %+v %v", j, err)
	}
	if j, err := AtHead(g, "/main", evManifest(), evC2, evC3); err != nil || j.AtHead || !slices.Equal(j.Outside, []string{"documentation/guide.md"}) {
		t.Errorf("the default layout: %+v %v", j, err)
	}
	if j, err := AtHead(g, "/main", evManifest(), "", evC3); err != nil || j.AtHead || j.Why == "" {
		t.Errorf("no commit: %+v %v", j, err)
	}
	g.diffErr = errors.New("fatal: bad object")
	if _, err := AtHead(g, "/main", evManifest(), evC2, evC3); err == nil || !strings.Contains(err.Error(), "git diff") {
		t.Errorf("a diff git cannot take: %v", err)
	}
	g.diffErr = nil
	g.files[evC3] = []string{"a.go", "b.go", "c.go", "d.go", "e.go", "f.go", "g.go"}
	if j, _ := AtHead(g, "/main", evManifest(), evC2, evC3); len(j.Outside) != maxOutside || !strings.HasPrefix(j.Why, "7 paths outside design, docs, and wip changed") {
		t.Errorf("many paths: %+v", j)
	}
}

func TestEvidenceMarshalsInSnakeCase(t *testing.T) {
	root := t.TempDir()
	evRecord(t, root, evC2, OutcomePassed, evC2, OutcomePassed)
	e, err := ReadEvidence(newGraphGit(evC3), root, evManifest(), "S-0662", true, true)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"applies":true`, `"checks_on":true`, `"passed_at_head":true`, `"at_head":true`, `"asked_by":"agent-S-0662"`, `"deploy_tool":"kluster_deploy"`, `"checks":{"commit"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("%s lacks %s", data, key)
		}
	}
}

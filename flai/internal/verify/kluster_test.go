package verify

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// The commits of klusterHistory, oldest first: k1 and k2 change code, k3
// only the docs and wip folders.
var (
	klK1 = strings.Repeat("a", 40)
	klK2 = strings.Repeat("b", 40)
	klK3 = strings.Repeat("c", 40)
)

// klusterGit answers the git a namespace's evidence is judged with, from a
// linear history k1, k2, k3: rev-parse resolves a story's branch to head,
// merge-base --is-ancestor walks the history, and diff --name-only -z lists
// what each commit after the first up to the second changed. It answers
// every other command as scriptGit's out does, when one is given.
type klusterGit struct {
	head string
	rest *scriptGit
}

var klusterParent = map[string]string{klK2: klK1, klK3: klK2}

var klusterFiles = map[string][]string{
	klK1: {"flai/main.go"},
	klK2: {"flai/internal/x/x.go"},
	klK3: {"docs/users/flai.md", "wip/agents/S-0662.md"},
}

func (g klusterGit) Run(dir, name string, args ...string) (string, error) {
	switch {
	case len(args) == 4 && args[0] == "rev-parse" && strings.HasPrefix(args[3], "refs/heads/story/"):
		if g.head == "" {
			return "", errors.New("exit status 1")
		}
		return g.head + "\n", nil
	case args[0] == "merge-base" && len(args) == 4 && args[1] == "--is-ancestor":
		for c := args[3]; c != ""; c = klusterParent[c] {
			if c == args[2] {
				return "", nil
			}
		}
		return "", errors.New("exit status 1")
	case args[0] == "diff" && args[1] == "--name-only" && args[2] == "--no-renames":
		from, to := args[len(args)-2], args[len(args)-1]
		var out []string
		for c := to; c != "" && c != from; c = klusterParent[c] {
			out = append(out, klusterFiles[c]...)
		}
		return strings.Join(out, "\x00"), nil
	case g.rest != nil:
		return g.rest.Run(dir, name, args...)
	}
	return "", errors.New("fatal: not a git repository")
}

func (g klusterGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (klusterGit) LookPath(name string) (string, error) { return name, nil }

// klusterManifest is a manifest with the default layout and a kluster:
// block that names a deploy.
func klusterManifest() manifest.Manifest {
	return manifest.Manifest{
		Layout:  map[string]string{"design": "design", "docs": "docs", "wip": "wip"},
		Kluster: &manifest.Kluster{Config: "integration/kluster.yaml", Deploy: []string{"scripts/deploy.sh"}},
	}
}

// klusterRecord writes story's namespace's record under root: a deploy at
// commit that ended outcome, asked for by the story's agent, with a checks
// run at checks that ended checksOutcome when checks is set.
func klusterRecord(t *testing.T, root, story, commit, outcome, checks, checksOutcome string) {
	t.Helper()
	ns := strings.ToLower(story)
	ended := time.Date(2026, 11, 13, 9, 30, 0, 0, time.UTC)
	r := kluster.Record{Story: story, Namespace: ns, Commit: commit, Outcome: outcome, AskedBy: "agent-" + story,
		StartedAt: ended.Add(-3 * time.Minute), EndedAt: ended,
		Steps: []kluster.Step{
			{Name: kluster.StepBuild, Outcome: kluster.OutcomePassed, Log: kluster.LogPath(root, ns, kluster.StepBuild)},
			{Name: kluster.StepDeploy, Outcome: outcome, Log: kluster.LogPath(root, ns, kluster.StepDeploy)},
		}}
	if outcome == kluster.OutcomeFailed {
		r.FailedStep, r.ExitCode, r.Error = kluster.StepDeploy, 2, "deploy: exit status 2"
	}
	if checks != "" {
		r.ChecksCommit, r.ChecksOutcome, r.ChecksLog = checks, checksOutcome, "/logs/checks.log"
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	path := kluster.RecordPath(root, ns)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// The part is filled from a deploy and checks run that passed at the head:
// it passes, its note saying so.
func TestKlusterEvidenceFillsThePartFromARecordAtTheHead(t *testing.T) {
	root := t.TempDir()
	klusterRecord(t, root, "S-0662", klK2, kluster.OutcomePassed, klK2, kluster.OutcomePassed)
	part, note, found := klusterEvidence(klusterGit{head: klK2}, root, klusterManifest(), "S-0662", true, true)
	if len(found) > 0 || part == nil {
		t.Fatalf("part %+v, found %+v", part, found)
	}
	d, c := part.Deploy, part.Checks
	if !part.Applies || !part.Passed || part.Namespace != "s-0662" || part.Head != klK2 || d == nil || c == nil {
		t.Fatalf("part %+v", part)
	}
	if d.Commit != klK2 || d.Outcome != kluster.OutcomePassed || d.AskedBy != "agent-S-0662" || d.EndedAt.IsZero() || !d.Head.AtHead || !d.PassedAtHead || len(d.Steps) != 2 {
		t.Errorf("deploy %+v", d)
	}
	if c.Commit != klK2 || c.Outcome != kluster.OutcomePassed || c.Log != "/logs/checks.log" || !c.Head.AtHead || !c.PassedAtHead {
		t.Errorf("checks %+v", c)
	}
	if note != "S-0662's namespace s-0662 is deployed and checked at the head "+klK2[:12] {
		t.Errorf("note %q", note)
	}
}

// The part is absent, and the step passes saying why, while the action is
// off, the manifest has no kluster: block or its block no deploy,
// kluster.agent_namespaces is 0, or the record's folder is not there, as in
// a story's container (ADR-0274).
func TestKlusterEvidenceIsAbsentWhileItDoesNotApply(t *testing.T) {
	zero := 0
	noBlock, noDeploy, noAgents := klusterManifest(), klusterManifest(), klusterManifest()
	noBlock.Kluster = nil
	noDeploy.Kluster = &manifest.Kluster{Config: "integration/kluster.yaml"}
	noAgents.Kluster.AgentNamespaces = &zero
	for _, tc := range []struct {
		name     string
		m        manifest.Manifest
		on       bool
		noRecord bool
		want     string
	}{
		{"the action off", klusterManifest(), false, false, "no kluster evidence applies to S-0662: the kluster host action is off for this project"},
		{"no block", noBlock, true, false, "no kluster evidence applies to S-0662: the manifest has no kluster: block"},
		{"no deploy", noDeploy, true, false, "no kluster evidence applies to S-0662: the manifest's kluster: block names no deploy"},
		{"agent_namespaces 0", noAgents, true, false, "no kluster evidence applies to S-0662: kluster.agent_namespaces is 0, so no story's agent deploys before review; the review deploys it"},
		{"no record folder", klusterManifest(), true, true, "the main checkout's .flai-cache/kluster is absent, as in a story's container, which does not mount it (ADR-0274)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if !tc.noRecord {
				// a deploy that failed at the head fails nothing while none applies
				klusterRecord(t, root, "S-0662", klK2, kluster.OutcomeFailed, "", "")
			}
			part, note, found := klusterEvidence(klusterGit{head: klK2}, root, tc.m, "S-0662", tc.on, true)
			if part != nil || len(found) > 0 || !strings.Contains(note, tc.want) {
				t.Errorf("part %+v, found %+v, note %q, want it to say %q", part, found, note, tc.want)
			}
		})
	}
}

// A record at an older commit, whose code the head changed since, passes
// and reports it not at the head, whether it passed or failed there; so
// does a namespace never deployed, naming the deploy command.
func TestKlusterEvidenceAtAnOlderCommitPassesNotAtTheHead(t *testing.T) {
	for _, outcome := range []string{kluster.OutcomePassed, kluster.OutcomeFailed} {
		root := t.TempDir()
		klusterRecord(t, root, "S-0662", klK1, outcome, klK1, outcome)
		part, note, found := klusterEvidence(klusterGit{head: klK2}, root, klusterManifest(), "S-0662", true, true)
		if len(found) > 0 || part == nil || part.Deploy == nil || part.Checks == nil {
			t.Fatalf("%s at k1: part %+v, found %+v", outcome, part, found)
		}
		if part.Passed || part.Deploy.Head.AtHead || part.Deploy.Commit != klK1 || part.Deploy.Outcome != outcome || part.Checks.Head.AtHead {
			t.Errorf("%s at k1: part %+v, deploy %+v", outcome, part, part.Deploy)
		}
		if !strings.Contains(note, "is not deployed at the head "+klK2[:12]) || !strings.Contains(note, "flai kluster deploy S-0662 --job") {
			t.Errorf("%s at k1: note %q", outcome, note)
		}
	}

	root := t.TempDir()
	if err := os.MkdirAll(kluster.Dir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	part, note, found := klusterEvidence(klusterGit{head: klK2}, root, klusterManifest(), "S-0662", true, true)
	if len(found) > 0 || part == nil || !part.Applies || part.Deploy != nil || !strings.Contains(note, "no deploy of it is recorded") {
		t.Errorf("never deployed: part %+v, found %+v, note %q", part, found, note)
	}
}

// A deploy that failed at the head fails the step, at the head as ADR-0292
// §3 widens it too; so does a checks run that failed or timed out there
// while the checks host action is on, and not while it is off.
func TestKlusterEvidenceFailsOnlyOnAFailureAtTheHead(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		commit, outcome       string
		checks, checksOutcome string
		head                  string
		checksOn              bool
		fails                 bool
		want                  []string
	}{
		{"the deploy failed at the head", klK2, kluster.OutcomeFailed, "", "", klK2, true, true,
			[]string{"failed at the deploy step", "(deploy: exit status 2)", "flai kluster log S-0662 --step deploy", "flai kluster deploy S-0662 --job", "kluster_deploy"}},
		{"the deploy failed with only docs and wip after it", klK2, kluster.OutcomeFailed, "", "", klK3, true, true,
			[]string{"failed at the deploy step"}},
		{"the checks failed at the head", klK2, kluster.OutcomePassed, klK2, kluster.OutcomeFailed, klK2, true, true,
			[]string{"its checks run there ended failed", "read its log /logs/checks.log", "flai kluster deploy S-0662 --job"}},
		{"the checks timed out at the head", klK2, kluster.OutcomePassed, klK2, "timed-out", klK2, true, true,
			[]string{"ended timed-out"}},
		{"the checks failed at the head with the checks action off", klK2, kluster.OutcomePassed, klK2, kluster.OutcomeFailed, klK2, false, false, nil},
		{"the checks were cancelled at the head", klK2, kluster.OutcomePassed, klK2, "cancelled", klK2, true, false, nil},
		{"the deploy passed at the head", klK2, kluster.OutcomePassed, klK2, kluster.OutcomePassed, klK2, true, false, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			klusterRecord(t, root, "S-0662", tc.commit, tc.outcome, tc.checks, tc.checksOutcome)
			part, note, found := klusterEvidence(klusterGit{head: tc.head}, root, klusterManifest(), "S-0662", true, tc.checksOn)
			if part == nil || part.Deploy == nil {
				t.Fatalf("part %+v", part)
			}
			if !tc.fails {
				if len(found) > 0 || note == "" {
					t.Errorf("found %+v, note %q, want it to pass", found, note)
				}
				return
			}
			if len(found) != 1 || found[0].Name != StepKluster || note != "" {
				t.Fatalf("found %+v, note %q, want one finding", found, note)
			}
			for _, want := range tc.want {
				if !strings.Contains(found[0].Message, want) {
					t.Errorf("finding %q lacks %q", found[0].Message, want)
				}
			}
		})
	}
}

// klusterVerifyFixture is newVerifyFixture whose manifest has a kluster:
// block naming a deploy, and whose story's branch has its head at the
// worktree's commit, with the kluster and checks host actions on.
func klusterVerifyFixture(t *testing.T) *verifyFixture {
	t.Helper()
	f := newVerifyFixture(t)
	f.opts.Project.Manifest.Kluster = klusterManifest().Kluster
	f.opts.Git = klusterGit{head: "0123abc", rest: f.git}
	f.opts.KlusterOn, f.opts.ChecksOn = true, true
	return f
}

// The kluster step comes after the check and before the tiers; it keeps
// what it read as the report's part, which the stored report holds; a
// failure at the head stops the run there, the tiers not reached.
func TestVerifyKlusterStepKeepsThePartAndStopsAtAFailureAtTheHead(t *testing.T) {
	f := klusterVerifyFixture(t)
	klusterRecord(t, f.root, "S-004", "0123abc", kluster.OutcomePassed, "0123abc", kluster.OutcomePassed)
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "kluster:passed", "go-test:passed", "integration:passed", "smoke:passed"}
	if !rep.Passed || !reflect.DeepEqual(states(rep), want) {
		t.Fatalf("steps %v, want %v", states(rep), want)
	}
	if got := rep.Steps[4].Note; got != "S-004's namespace s-004 is deployed and checked at the head 0123abc" {
		t.Errorf("kluster note %q", got)
	}
	if rep.Kluster == nil || !rep.Kluster.Passed || rep.Kluster.Deploy == nil || rep.Kluster.Deploy.Commit != "0123abc" {
		t.Fatalf("part %+v", rep.Kluster)
	}
	last, ok, err := LastReport(f.opts.Project, "S-004")
	if err != nil || !ok || !reflect.DeepEqual(last.Kluster, rep.Kluster) {
		t.Errorf("stored part %+v (%v, %v), want %+v", last.Kluster, ok, err, rep.Kluster)
	}

	klusterRecord(t, f.root, "S-004", "0123abc", kluster.OutcomeFailed, "", "")
	f.opts.Fresh = true
	rep, err = Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"rebase:passed", "sync:passed", "narrative:passed", "check:passed", "kluster:failed", "go-test:not-reached", "integration:not-reached", "smoke:not-reached"}
	if rep.Passed || rep.StoppedAt != StepKluster || !reflect.DeepEqual(states(rep), want) {
		t.Fatalf("stopped at %q, steps %v, want %v", rep.StoppedAt, states(rep), want)
	}
	if rep.Kluster == nil || rep.Kluster.Deploy.Outcome != kluster.OutcomeFailed || len(rep.Steps[4].Findings) != 1 {
		t.Errorf("part %+v, step %+v", rep.Kluster, rep.Steps[4])
	}
}

// A project without a kluster: block has no kluster step and no part; one
// whose block names no deploy has the step, passing, and no part.
func TestVerifyHasNoKlusterStepWithoutABlockAndNoPartWithoutADeploy(t *testing.T) {
	f := newVerifyFixture(t)
	rep, err := Verify(context.Background(), f.opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range rep.Steps {
		if s.Name == StepKluster {
			t.Errorf("a project without a kluster: block has the step %+v", s)
		}
	}
	if rep.Kluster != nil {
		t.Errorf("part %+v, want none", rep.Kluster)
	}

	f = klusterVerifyFixture(t)
	f.opts.Project.Manifest.Kluster.Deploy = nil
	f.opts.Fresh = true
	if rep, err = Verify(context.Background(), f.opts); err != nil {
		t.Fatal(err)
	}
	if rep.Kluster != nil || rep.Steps[4].Name != StepKluster || rep.Steps[4].State != Passed || !strings.Contains(rep.Steps[4].Note, "names no deploy") {
		t.Errorf("part %+v, step %+v", rep.Kluster, rep.Steps[4])
	}
	data, err := os.ReadFile(ReportPath(f.opts.Project, "S-004"))
	if err != nil || strings.Contains(string(data), `"kluster": {`) {
		t.Errorf("stored report %v:\n%s", err, data)
	}
}

// A report stored before the part was kept, without it, still reads.
func TestLastReportReadsAReportWithoutAKlusterPart(t *testing.T) {
	f := newVerifyFixture(t)
	path := ReportPath(f.opts.Project, "S-004")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	old := `{"story":"S-004","commit":"0123abc","base":"main","ran_at":"2026-09-01T12:00:00Z","duration_ms":10,"passed":true,"steps":[{"name":"check","state":"passed","duration_ms":10}]}`
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, ok, err := LastReport(f.opts.Project, "S-004")
	if err != nil || !ok || rep.Kluster != nil || !rep.Passed || rep.Commit != "0123abc" || len(rep.Steps) != 1 {
		t.Errorf("report %+v (%v, %v)", rep, ok, err)
	}
}

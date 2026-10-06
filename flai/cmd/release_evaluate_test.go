package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/release"
)

// evaluateProject is pendingProject with S-0001 and S-0002 accepted and not
// yet released, S-0001 worth 300 a week and S-0002 with no value, and S-0003
// still in backlog under the same epic.
func evaluateProject(t *testing.T) (root, remote string) {
	t.Helper()
	root, remote = pendingProject(t)
	for _, id := range []string{"S-0001", "S-0002"} {
		if _, errOut, code := runIn(t, root, "accept", id); code != 0 {
			t.Fatalf("accept %s: %s", id, errOut)
		}
	}
	m, _ := filepath.Glob(filepath.Join(root, "wip/archive/kanban/stories/S-0001-*.md"))
	if len(m) != 1 {
		t.Fatalf("S-0001 is archived: %v", m)
	}
	data, _ := os.ReadFile(m[0])
	s := strings.Replace(string(data), "\n---\n", "\ncost_of_delay:\n  value: 300\n  by: t\n  at: 2026-10-01T09:00:00Z\n---\n", 1)
	if s == string(data) {
		t.Fatalf("no front matter to add a value to: %s", data)
	}
	if err := os.WriteFile(m[0], []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: value S-0001")
	return root, remote
}

func setReleasePolicy(t *testing.T, root, yaml string) {
	t.Helper()
	path := filepath.Join(root, "system-flow.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	base, _, _ := strings.Cut(string(data), "orchestration:\n")
	if err := os.WriteFile(path, []byte(base+"orchestration:\n  release:\n"+yaml), 0o644); err != nil {
		t.Fatal(err)
	}
}

// releaseState is everything flai release --evaluate must leave alone: the
// tags, the files, the head, and the remote.
func releaseState(t *testing.T, root, remote string) string {
	t.Helper()
	return strings.Join([]string{
		gitIn(t, root, "tag", "--list"),
		gitIn(t, root, "status", "--porcelain"),
		gitIn(t, root, "rev-parse", "HEAD"),
		gitIn(t, remote, "show-ref"),
	}, "\n--\n")
}

func TestReleaseEvaluateThreshold(t *testing.T) {
	root, remote := evaluateProject(t)
	setReleasePolicy(t, root, "    policy: threshold\n    value: 300\n")
	before := releaseState(t, root, remote)

	out, errOut, code := runIn(t, root, "release", "--evaluate")
	if code != 0 {
		t.Fatalf("evaluate: %d %s", code, errOut)
	}
	lines := strings.Split(out, "\n")
	if lines[0] != "release policy threshold: met" {
		t.Fatalf("the first line is the policy and the verdict: %s", out)
	}
	for _, want := range []string{"at or over the threshold of 300 USD/week", "300 USD/week of 300 USD/week", "count", "unvalued  S-0002", "pending:", "S-0001  300 USD/week", "S-0002  no value"} {
		if !strings.Contains(out, want) {
			t.Errorf("%q is in the output: %s", want, out)
		}
	}
	if after := releaseState(t, root, remote); after != before {
		t.Errorf("evaluate changes no file, tag, or remote:\n%s\nvs\n%s", before, after)
	}

	setReleasePolicy(t, root, "    policy: threshold\n    count: 3\n")
	out, _, code = runIn(t, root, "release", "--evaluate")
	if code != 0 || !strings.HasPrefix(out, "release policy threshold: not met\n") || !strings.Contains(out, "2 of 3") || !strings.Contains(out, "their count, 2, is under the threshold of 3") {
		t.Errorf("two pending stories are under a count of three: %d %s", code, out)
	}
}

func TestReleaseEvaluateThemeAndJudgement(t *testing.T) {
	root, remote := evaluateProject(t)
	setReleasePolicy(t, root, "    policy: theme\n    epic: E-0001\n")
	out, errOut, code := runIn(t, root, "release", "--evaluate")
	if code != 0 || !strings.HasPrefix(out, "release policy theme: not met\n") {
		t.Fatalf("S-0003 is in backlog: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{"not accepted  S-0003", "epic", "theme:", "S-0001  not released", "S-0003  backlog"} {
		if !strings.Contains(out, want) {
			t.Errorf("%q is in the output: %s", want, out)
		}
	}

	setReleasePolicy(t, root, "    policy: judgement\n")
	before := releaseState(t, root, remote)
	out, errOut, code = runIn(t, root, "release", "--evaluate", "--json")
	if code != 0 {
		t.Fatalf("evaluate --json: %d %s", code, errOut)
	}
	var ev release.Evaluation
	if err := json.Unmarshal([]byte(out), &ev); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if ev.Policy != "judgement" || ev.Met || ev.Count != 2 || ev.Value != 300 || len(ev.Pending) != 2 || !strings.Contains(ev.Reason, "the orchestrator's or the operator's call") {
		t.Errorf("judgement is never met, with the pending figures: %+v", ev)
	}
	if after := releaseState(t, root, remote); after != before {
		t.Errorf("evaluate changes no file, tag, or remote:\n%s\nvs\n%s", before, after)
	}
}

func TestReleaseEvaluateRefusesToRelease(t *testing.T) {
	root, remote := evaluateProject(t)
	before := releaseState(t, root, remote)
	for _, args := range [][]string{
		{"release", "--evaluate", "--apply"},
		{"release", "--evaluate", "--pending"},
		{"release", "--evaluate", "--dry-run"},
		{"release", "--evaluate", "--deliver", "cli"},
		{"release", "--evaluate", "S-0001"},
	} {
		_, errOut, code := runIn(t, root, args...)
		if code == 0 || !strings.Contains(errOut, "--evaluate") {
			t.Errorf("flai %v is refused and says why: %d %s", args, code, errOut)
		}
	}
	if after := releaseState(t, root, remote); after != before {
		t.Errorf("a refusal changes nothing:\n%s\nvs\n%s", before, after)
	}
}

// S-0222: with whole_epics set, S-0001 and S-0002 are held back while E-0001,
// in review in evaluateProject, is moved back to in-progress; in review, it
// holds nothing.
func TestReleaseEvaluateWholeEpics(t *testing.T) {
	root, remote := evaluateProject(t)
	setReleasePolicy(t, root, "    policy: threshold\n    value: 300\n    whole_epics: true\n")
	if out, _, _ := runIn(t, root, "release", "--evaluate"); !strings.HasPrefix(out, "release policy threshold: met\n") || strings.Contains(out, "held by epic") {
		t.Errorf("an epic in review holds nothing: %s", out)
	}
	if _, errOut, code := runIn(t, root, "move", "E-0001", "in-progress", "--reason", "S-0003 is still to come"); code != 0 {
		t.Fatalf("move E-0001 back: %s", errOut)
	}
	setReleasePolicy(t, root, "    policy: threshold\n    value: 300\n    whole_epics: true\n")
	before := releaseState(t, root, remote)

	out, errOut, code := runIn(t, root, "release", "--evaluate")
	if code != 0 || !strings.HasPrefix(out, "release policy threshold: not met\n") {
		t.Fatalf("a batch held by its epic is not met: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{"at or over the threshold of 300 USD/week; whole_epics holds the batch back until its epic is in review or done: E-0001 is ", "with S-0001, S-0002", "held by epic:", "S-0001  E-0001  in-progress  What we found", "S-0002  E-0001  in-progress  Another fix", "pending:"} {
		if !strings.Contains(out, want) {
			t.Errorf("%q is in the output: %s", want, out)
		}
	}

	out, errOut, code = runIn(t, root, "release", "--evaluate", "--json")
	if code != 0 {
		t.Fatalf("evaluate --json: %d %s", code, errOut)
	}
	var ev release.Evaluation
	if err := json.Unmarshal([]byte(out), &ev); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if ev.Met || len(ev.HeldByEpic) != 2 || ev.HeldByEpic[0].Epic != "E-0001" || ev.HeldByEpic[0].EpicStatus != "in-progress" || !strings.Contains(out, `"held_by_epic"`) {
		t.Errorf("the JSON holds the stories with their epic and its status: %s", out)
	}
	if after := releaseState(t, root, remote); after != before {
		t.Errorf("evaluate changes no file, tag, or remote:\n%s\nvs\n%s", before, after)
	}

	setReleasePolicy(t, root, "    policy: threshold\n    value: 300\n")
	if out, _, _ = runIn(t, root, "release", "--evaluate"); !strings.HasPrefix(out, "release policy threshold: met\n") || strings.Contains(out, "held by epic") {
		t.Errorf("whole_epics off holds nothing: %s", out)
	}
}

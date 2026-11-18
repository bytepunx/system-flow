package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// systemFlowTiers is this repository's system-flow.yaml tests block as it
// stood when S-0807 began: each tier's name, suite, paths, all_only, covers,
// queue, when, manual, and parallel as the manifest gives them, which are
// what choose a story's tiers. Each command is sh exiting 0 with the tier's
// files, and no tier gives a dir or a format, so that the tiers run in a
// test; none of those changes which tiers a story's paths select. It is a
// copy, not the live manifest, so that a later story's edit of the manifest
// does not change what the tests below expect.
const systemFlowTiers = `tests:
  - {name: gofmt, suite: static, command: [sh, -c, 'exit 0', gofmt, "{files}"], paths: ["flai/**/*.go"]}
  - {name: vet, suite: static, command: [sh, -c, 'exit 0', vet, "{files}"], paths: ["flai/**/*.go", "flai/**/testdata/**"]}
  - {name: golangci-lint, suite: static, command: [sh, -c, 'exit 0', golangci-lint, "{files}"], paths: ["flai/**/*.go", "flai/**/testdata/**"]}
  - {name: go-test, suite: behavior, command: [sh, -c, 'exit 0', go-test, "{files}"], paths: ["flai/**/*.go", "flai/**/testdata/**"]}
  - {name: settings-index, suite: behavior, command: [sh, -c, 'exit 0', settings-index, "{files}"], paths: ["docs/operators/settings.md", "docs/users/flai-reference.md"]}
  - {name: flaiover-lint, suite: static, command: [sh, -c, 'exit 0', flaiover-lint, "{files}"], paths: ["flaiover/**", "!flaiover/node_modules/**", "!flaiover/build/**", "!flaiover/.svelte-kit/**"]}
  - {name: vitest, suite: behavior, command: [sh, -c, 'exit 0', vitest, "{files}"], paths: ["flaiover/src/**", "flaiover/*.ts", "flaiover/*.js"]}
  - {name: flaiover-flai, suite: behavior, command: [sh, -c, 'exit 0', flaiover-flai, "{files}"], paths: ["flai/internal/hostapi/**/*.go", "flai/cmd/hostapi.go"]}
  - {name: design-routes, suite: behavior, command: [sh, -c, 'exit 0', design-routes, "{files}"], paths: ["design/system/flaiover-dashboard.md", "flaiover/src/routes/**"]}
  - {name: svelte-check, suite: static, command: [sh, -c, 'exit 0', svelte-check, "{files}"], paths: ["flaiover/src/**/*.ts", "flaiover/src/**/*.js", "flaiover/src/**/*.svelte", "flaiover/*.ts", "flaiover/*.js", "flaiover/tsconfig.json", "flaiover/package.json", "flaiover/pnpm-workspace.yaml", "scripts/flaiover-check.sh"]}
  - {name: shellcheck, suite: static, command: [sh, -c, 'exit 0', shellcheck, "{files}"], paths: ["scripts/**/*.sh"]}
  - {name: markdown, suite: static, command: [sh, -c, 'exit 0', markdown, "{files}"], paths: ["**/*.md", "!**/node_modules/**", "!**/testdata/**", "!bin/**", "!flaiover/build/**", "!**/.svelte-kit/**", "!.flai-cache/**"]}
  - {name: template, suite: behavior, command: [sh, -c, 'exit 0', template], paths: ["template/**", "scripts/template-test.sh"], all_only: true}
  - {name: flaiover, suite: behavior, command: [sh, -c, 'exit 0', flaiover], paths: ["flaiover/**", "!flaiover/node_modules/**", "!flaiover/build/**", "!flaiover/.svelte-kit/**", "scripts/flaiover-*.sh", "scripts/env.sh"], all_only: true, covers: [flaiover-lint, vitest, flaiover-flai, design-routes, svelte-check], parallel: true}
  - {name: integration, suite: integration, command: [sh, -c, 'exit 0', integration], paths: ["flai/**", "template/**", "scripts/integration.sh"], all_only: true, covers: [go-test, settings-index], queue: true, when: ["flai/**", "template/**"], parallel: true}
  - {name: smoke, suite: smoke, command: [sh, -c, 'exit 0', smoke], paths: ["flai/**", "template/**", "scripts/**/*.sh"], all_only: true, queue: true, when: ["flai/**", "template/**"], parallel: true}
  - {name: gateway, suite: smoke, command: [sh, -c, 'exit 0', gateway], paths: ["scripts/gateway-smoke.sh"], all_only: true, manual: true}
  - {name: race, suite: integration, command: [sh, -c, 'exit 0', race], paths: ["scripts/integration.sh"], all_only: true, manual: true}
`

// baselineStories are the stories the tests below verify on the baseline,
// by the one path each one's branch changes, with the tiers flai verify
// chose for it before S-0807 while the merge queue can batch: those it runs,
// a covered one as "<tier> covered by <tier>" and a queue tier its when
// keeps here as "<tier> ran here", and those it leaves to the queue.
var baselineStories = []struct {
	name, path  string
	run, queued []string
	// steps are the report's tier steps, for the stories flai verify runs
	steps string
}{
	{
		name: "a story that changes flai", path: "flai/cmd/verify.go",
		run:   []string{"gofmt", "vet", "golangci-lint", "go-test covered by integration", "integration ran here", "smoke ran here"},
		steps: "tier gofmt:passed,tier vet:passed,tier golangci-lint:passed,tier go-test:skipped,tier integration:passed,tier smoke:passed",
	},
	{
		name: "a story that changes only markdown", path: "docs/users/flai-reference.md",
		run:   []string{"settings-index", "markdown"},
		steps: "tier settings-index:passed,tier markdown:passed",
	},
	{
		name: "a story that changes the template", path: "template/root/README.md",
		run: []string{"markdown", "template", "integration ran here", "smoke ran here"},
	},
	{
		name: "a story that changes a script", path: "scripts/integration.sh",
		run: []string{"shellcheck"}, queued: []string{"integration", "smoke"},
	},
}

// selectedNames are the selection's tiers as baselineStories gives them.
func selectedNames(selection []verify.Selected) []string {
	var out []string
	for _, s := range selection {
		name := s.Tier.Name
		switch {
		case s.CoveredBy != "":
			name += " covered by " + s.CoveredBy
		case s.RanFor != "":
			name += " ran here"
		}
		out = append(out, name)
	}
	return out
}

// S-0807 criterion 4, ADR-0312 §8: on the baseline workflow a story in
// progress, and one in review, needs every suite, so the tiers chosen for
// this repository's tiers with the suites its close-out needs are those
// chosen with none, the tiers flai verify chose before suites.
func TestBaselineSelectsTheTiersFlaiVerifySelectedBeforeSuites(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"in-progress", "review"} {
		if got := verify.NeededSuites(workflow.Baseline(), state); got != nil {
			t.Errorf("the baseline's close-out of a story in %s needs the suites %q, want every suite (nil)", state, got)
		}
	}
	good, err := os.ReadFile("../internal/metrics/testdata/good/system-flow.yaml")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "system-flow.yaml"), append(good, systemFlowTiers...), 0o644); err != nil {
		t.Fatal(err)
	}
	tiers, err := verify.CheckoutTiers(dir)
	if err != nil {
		t.Fatal(err)
	}
	fsys := os.DirFS(dir)
	needed := verify.NeededSuites(workflow.Baseline(), "in-progress")
	for _, s := range baselineStories {
		paths := []string{s.path}
		run, queued := verify.SelectStoryLeavingQueue(fsys, tiers, paths, nil, needed)
		if got := selectedNames(run); !reflect.DeepEqual(got, s.run) {
			t.Errorf("%s: runs %q, want %q", s.name, got, s.run)
		}
		if got := selectedNames(queued); !reflect.DeepEqual(got, s.queued) {
			t.Errorf("%s: leaves %q to the merge queue, want %q", s.name, got, s.queued)
		}
		allRun, allQueued := verify.SelectStoryLeavingQueue(fsys, tiers, paths, nil, nil)
		if !reflect.DeepEqual(run, allRun) || !reflect.DeepEqual(queued, allQueued) {
			t.Errorf("%s: the needed suites' selection differs from every suite's", s.name)
		}
		if got, want := verify.SelectStorySuites(fsys, tiers, paths, needed), verify.SelectStory(fsys, tiers, paths); !reflect.DeepEqual(got, want) {
			t.Errorf("%s without the merge queue: the needed suites select %q, every suite %q", s.name, selectedNames(got), selectedNames(want))
		}
	}
}

// changedGit is verifyGit for a branch that changed the files changed,
// one per line, in place of docs/README.md.
type changedGit struct {
	verifyGit
	changed string
}

func (g changedGit) Run(dir, name string, args ...string) (string, error) {
	if strings.Join(args, " ") == "diff --name-only --diff-filter=d main...HEAD" {
		return g.changed, nil
	}
	return g.verifyGit.Run(dir, name, args...)
}

func (g changedGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

// S-0807 criterion 4: on the baseline, with no workflow file and with the
// baseline's own, flai verify on a story in progress runs the tiers of every
// suite that it ran before S-0807, the queue tiers whose when the branch
// matches among them, and names no suites, in text or in --json.
func TestVerifyOnTheBaselineRunsTheTiersItRanBeforeSuites(t *testing.T) {
	t.Parallel()
	baseline, err := os.ReadFile("../internal/workflow/baseline/board.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range []struct {
		name  string
		board []byte
	}{{"no workflow file", nil}, {"the baseline's workflow file", baseline}} {
		for _, s := range baselineStories[:2] {
			t.Run(file.name+", "+s.name, func(t *testing.T) {
				t.Parallel()
				root := queueFixture(t, true, func(root string) []string { return []string{root} })
				good, err := os.ReadFile("../internal/metrics/testdata/good/system-flow.yaml")
				if err != nil {
					t.Fatal(err)
				}
				manifest := filepath.Join(root, ".flai-cache", "worktrees", "S-004", "system-flow.yaml")
				if err := os.WriteFile(manifest, append(good, systemFlowTiers...), 0o644); err != nil {
					t.Fatal(err)
				}
				if file.board != nil {
					board := filepath.Join(root, "design", "workflows", "board.md")
					if err := os.WriteFile(board, file.board, 0o644); err != nil {
						t.Fatal(err)
					}
				}
				var git execx.Runner = changedGit{changed: s.path + "\n"}
				out, errOut, code := runVerify(t, root, git, "verify", "S-004", "--json")
				rep := decodeReport(t, out)
				if code != 0 || !rep.Passed {
					t.Fatalf("exit %d, passed %v\n%s%s", code, rep.Passed, out, errOut)
				}
				if want := "rebase:passed,sync:passed,narrative:passed,check:passed," + s.steps; stepStates(rep) != want {
					t.Errorf("steps %s, want %s", stepStates(rep), want)
				}
				if rep.Suites != nil {
					t.Errorf("the report names the suites %q, want none", rep.Suites)
				}
				text, errOut, code := runVerify(t, root, git, "verify", "S-004", "--fresh")
				if code != 0 || strings.Contains(text, "suites ") || strings.Contains(text, "merge queue") {
					t.Errorf("text: exit %d, want every tier run with no suites line\n%s%s", code, text, errOut)
				}
			})
		}
	}
}

// decodeReport is flai verify --json's answer out.
func decodeReport(t *testing.T, out string) verify.Report {
	t.Helper()
	var rep verify.Report
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatalf("flai verify --json: not JSON: %v\n%s", err, out)
	}
	return rep
}

// baselineMoveProject is verifyProject on the baseline workflow: S-0001 in
// progress, with a task, its one criterion ticked, and a worktree whose one
// commit it answers.
func baselineMoveProject(t *testing.T) (root, head string) {
	t.Helper()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	asOperator(t)
	root = tempProject(t)
	for _, args := range [][]string{{"config", "set", "author", "olive"}, {"epic", "new", "Epic"}, {"story", "new", "Verified", "--epic", "E-0001"}, {"task", "new", "--story", "S-0001", "Step"}} {
		mustRun(t, root, args...)
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-verified.md")
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] ok\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{"ready", "in-progress"} {
		mustRun(t, root, "move", "S-0001", to)
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "init", "-q")
	if err := os.WriteFile(filepath.Join(wt, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "add", "a.go")
	gitIn(t, wt, "commit", "-q", "-m", "a")
	return root, gitIn(t, wt, "rev-parse", "HEAD")
}

// S-0807 criterion 4, ADR-0312 §8: on the baseline, flai move takes a story
// to review without a passing verify at its head, as it did before S-0807,
// warning that the move needs one; once the record holds a passing run at
// the head, the move warns of nothing.
func TestMoveToReviewOnTheBaselineWarnsOfVerifyAndRefusesNothing(t *testing.T) {
	for _, tc := range []struct {
		name   string
		record func(t *testing.T, root, head string)
		warns  bool
	}{
		{"no record at the head", func(*testing.T, string, string) {}, true},
		{"a failing record at the head", func(t *testing.T, root, head string) { keepVerdict(t, root, head, false) }, true},
		{"a passing record at the head", func(t *testing.T, root, head string) { keepVerdict(t, root, head, true) }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, head := baselineMoveProject(t)
			tc.record(t, root, head)
			out, errOut, code := runIn(t, root, "move", "S-0001", "review")
			if code != 0 || !strings.Contains(out, "S-0001 → review") {
				t.Fatalf("the baseline must let the move to review through: exit %d\n%s%s", code, out, errOut)
			}
			if got := stateOf(t, root, "S-0001"); got != "review" {
				t.Errorf("status %q, want review", got)
			}
			if !tc.warns {
				if errOut != "" {
					t.Errorf("the move with a passing verify at the head must warn of nothing:\n%s", errOut)
				}
				return
			}
			for _, want := range []string{"workflow policy warning", "needs its last flai verify to have passed every suite at its branch's head", "run flai verify S-0001"} {
				if !strings.Contains(errOut, want) {
					t.Errorf("the warning must say %q:\n%s", want, errOut)
				}
			}
		})
	}
}

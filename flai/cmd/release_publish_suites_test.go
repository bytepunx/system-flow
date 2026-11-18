package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// publishSuitesTiers are tiers for a project's manifest: boom, of suite x,
// prints boom and fails; ok, of suite ok, and later, of suite later, each
// write the marker file named after them beside marker and pass.
func publishSuitesTiers(marker string) string {
	return "tests:\n" +
		"  - name: boom\n    suite: x\n    command: [sh, -c, 'echo boom; exit 1']\n    paths: [\"cli/**\"]\n" +
		"  - name: ok\n    suite: ok\n    command: [sh, -c, 'touch \"$0\"', \"" + marker + ".ok\"]\n    paths: [\"cli/**\"]\n" +
		"  - name: later\n    suite: later\n    command: [sh, -c, 'touch \"$0\"', \"" + marker + ".later\"]\n    paths: [\"cli/**\"]\n"
}

// publishSuitesProject is servedPublishProject, S-0001 accepted and pending
// as cli/v1.1.0, with publishSuitesTiers committed to its manifest on main
// and, unless verifyAt is empty, the baseline workflow file committed with
// its verify_at replaced by verifyAt. It answers the root, the remote, the
// marker prefix the tiers write beside, and main's head.
func publishSuitesProject(t *testing.T, verifyAt string) (root, remote, marker, head string) {
	t.Helper()
	root, remote = servedPublishProject(t)
	marker = filepath.Join(t.TempDir(), "ran")
	m, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	head = commitOnMain(t, root, "system-flow.yaml", string(m)+publishSuitesTiers(marker), "chore: the publish tiers")
	if verifyAt != "" {
		base, err := os.ReadFile(filepath.Join("..", "internal", "workflow", "baseline", "board.md"))
		if err != nil {
			t.Fatal(err)
		}
		const baselineVerifyAt = "verify_at: {batch: all, ci: all}\n"
		if !strings.Contains(string(base), baselineVerifyAt) {
			t.Fatalf("the baseline workflow file no longer gives %q", baselineVerifyAt)
		}
		head = commitOnMain(t, root, "design/workflows/board.md", strings.Replace(string(base), baselineVerifyAt, verifyAt, 1), "chore: the workflow file")
	}
	return root, remote, marker, strings.TrimSpace(head)
}

// assertNothingPublished checks that nothing of the release was written:
// main still at head, a clean tree, no cli/v1.1.0 here or on the remote, and
// the remote's main where it was.
func assertNothingPublished(t *testing.T, root, remote, head, remoteHead string) {
	t.Helper()
	if now := strings.TrimSpace(gitIn(t, root, "rev-parse", "main")); now != head {
		t.Errorf("no commit is made: main is at %s, was %s", now, head)
	}
	if status := gitIn(t, root, "status", "--porcelain"); strings.TrimSpace(status) != "" {
		t.Errorf("no version file or changelog is written: %s", status)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("nothing tagged here: %s", tags)
	}
	if tags := gitIn(t, remote, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("nothing tagged on the remote: %s", tags)
	}
	if now := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")); now != remoteHead {
		t.Errorf("nothing pushed: the remote's main is at %s, was %s", now, remoteHead)
	}
}

// ADR-0312 §6: a tier of a suite verify_at.publish names that fails stops
// the publish before anything is applied: exit 5, the message naming the
// tier, its suite, and the log of its whole output, the suites after it not
// run, and no commit, version file, changelog, tag, or push. With --json the
// answer gives the suites, each tier's outcome, and the log.
func TestReleasePendingRefusedWhenAPublishSuiteFails(t *testing.T) {
	t.Parallel()
	root, remote, marker, head := publishSuitesProject(t, "verify_at: {batch: all, publish: [x, later], ci: all}\n")
	remoteHead := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main"))
	log := filepath.Join(root, ".flai-cache", "publish", "20260915T210000Z", "boom.log")

	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != exitReleaseSuiteFailed {
		t.Fatalf("a failing publish suite exits %d: %d %s %s", exitReleaseSuiteFailed, code, out, errOut)
	}
	for _, want := range []string{"tier boom of suite x failed at " + short(head), "verify_at.publish in design/workflows/board.md", "nothing was applied, committed, tagged, or pushed", "its whole output is in " + log} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal says %q: %s", want, errOut)
		}
	}
	if !strings.Contains(out, "failed boom") {
		t.Errorf("the output shows the tier failed: %s", out)
	}
	if data, err := os.ReadFile(log); err != nil || !strings.Contains(string(data), "boom") {
		t.Errorf("the log keeps the tier's output: %q %v", data, err)
	}
	if _, err := os.Stat(marker + ".later"); err == nil {
		t.Error("the suite after the failing one ran")
	}
	assertNothingPublished(t, root, remote, head, remoteHead)

	out, errOut, code = runIn(t, root, "release", "--pending", "--json")
	if code != exitReleaseSuiteFailed || !strings.Contains(errOut, "tier boom of suite x failed") {
		t.Fatalf("with --json a failing publish suite exits %d naming the tier: %d %s %s", exitReleaseSuiteFailed, code, out, errOut)
	}
	var res struct {
		Pushed        bool `json:"pushed"`
		PublishSuites struct {
			Suites []string `json:"suites"`
			Head   string   `json:"head"`
			Log    string   `json:"log"`
			Passed bool     `json:"passed"`
			Tiers  []struct {
				Suite string `json:"suite"`
				Name  string `json:"name"`
				State string `json:"state"`
				Log   string `json:"log"`
			} `json:"tiers"`
		} `json:"publish_suites"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("the answer is JSON: %v\n%s", err, out)
	}
	ps := res.PublishSuites
	if res.Pushed || ps.Passed || strings.Join(ps.Suites, ",") != "x,later" || ps.Head != head || ps.Log != filepath.Dir(log) {
		t.Errorf("publish_suites gives the suites, the head, the log, and the failure: %+v", res)
	}
	if len(ps.Tiers) != 1 || ps.Tiers[0].Suite != "x" || ps.Tiers[0].Name != "boom" || ps.Tiers[0].State != "failed" || ps.Tiers[0].Log != log {
		t.Errorf("publish_suites gives each tier's outcome: %+v", ps.Tiers)
	}
	assertNothingPublished(t, root, remote, head, remoteHead)
}

// ADR-0312 §6: suites verify_at.publish names that pass publish as before,
// having run, with the folder of their output named.
func TestReleasePendingPublishesWhenThePublishSuitesPass(t *testing.T) {
	t.Parallel()
	root, remote, marker, _ := publishSuitesProject(t, "verify_at: {batch: all, publish: [ok, later], ci: all}\n")
	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.1.0") || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("passing publish suites publish: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(out, "verify_at.publish: the suites ok, later passed") || !strings.Contains(out, filepath.Join(root, ".flai-cache", "publish")) {
		t.Errorf("the answer says the suites passed and where their output is: %s", out)
	}
	for _, suite := range []string{"ok", "later"} {
		if _, err := os.Stat(marker + "." + suite); err != nil {
			t.Errorf("the suite %s ran: %v", suite, err)
		}
	}
	if tags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the tag reaches the remote: %s", tags)
	}
}

// ADR-0312 §8: with no publish in the workflow file, as on the baseline,
// and with no workflow file at all, publishing runs no tier and answers as
// before.
func TestReleasePendingRunsNoTierWithoutPublish(t *testing.T) {
	t.Parallel()
	for name, verifyAt := range map[string]string{"no workflow file": "", "the baseline's verify_at": "verify_at: {batch: all, ci: all}\n"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, remote, marker, _ := publishSuitesProject(t, verifyAt)
			out, errOut, code := runIn(t, root, "release", "--pending", "--json")
			if code != 0 || strings.Contains(out, "publish_suites") || !strings.Contains(out, `"pushed": true`) {
				t.Fatalf("a publish with no publish suites publishes as before: %d %s %s", code, out, errOut)
			}
			for _, suite := range []string{"ok", "later"} {
				if _, err := os.Stat(marker + "." + suite); err == nil {
					t.Errorf("the suite %s ran", suite)
				}
			}
			if tags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(tags, "cli/v1.1.0") {
				t.Errorf("the tag reaches the remote: %s", tags)
			}
		})
	}
}

// ADR-0312 §6: --dry-run names the suites verify_at.publish would run and
// their tiers, in publish_suites with --json, and runs none.
func TestReleasePendingDryRunNamesThePublishSuites(t *testing.T) {
	t.Parallel()
	root, remote, marker, head := publishSuitesProject(t, "verify_at: {batch: all, publish: [ok, x], ci: all}\n")
	remoteHead := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main"))
	out, errOut, code := runIn(t, root, "release", "--pending", "--dry-run")
	if code != 0 || !strings.Contains(out, "verify_at.publish would run the suites ok, x before tagging: the tiers ok (ok), boom (x)") || !strings.Contains(out, "dry run: nothing changed") {
		t.Fatalf("a dry run names the publish suites: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "release", "--pending", "--dry-run", "--json")
	var res struct {
		PublishSuites struct {
			Suites   []string `json:"suites"`
			WouldRun []struct {
				Suite string `json:"suite"`
				Tier  string `json:"tier"`
			} `json:"would_run"`
		} `json:"publish_suites"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil {
		t.Fatalf("a dry run with --json answers JSON: %d %s %s", code, out, errOut)
	}
	if ps := res.PublishSuites; strings.Join(ps.Suites, ",") != "ok,x" || len(ps.WouldRun) != 2 || ps.WouldRun[0].Tier != "ok" || ps.WouldRun[1].Suite != "x" {
		t.Errorf("publish_suites names the suites and the tiers they would run: %+v", ps)
	}
	if _, err := os.Stat(marker + ".ok"); err == nil {
		t.Error("a dry run ran a tier")
	}
	if _, err := os.Stat(filepath.Join(root, ".flai-cache", "publish")); err == nil {
		t.Error("a dry run kept a log")
	}
	assertNothingPublished(t, root, remote, head, remoteHead)
}

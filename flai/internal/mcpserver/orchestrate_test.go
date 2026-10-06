package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// stories reads the ids, in order, of a list of ranked stories in a tool's
// answer.
func stories(v any) (ids []string) {
	list, _ := v.([]any)
	for _, s := range list {
		m, _ := s.(map[string]any)
		id, _ := m["id"].(string)
		ids = append(ids, id)
	}
	return ids
}

// S-0217: order_by_policy computes the ready column's order by a policy, the
// project's when none is given, and writes nothing.
func TestOrderByPolicyOrdersTheReadyColumn(t *testing.T) {
	f := setup(t)
	bare := f.readyStory(t, "Bare", t0)
	valued := f.readyStory(t, "Valued", t0)
	if _, failed := f.call(t, "item_edit", map[string]any{"id": valued.ID, "cost_of_delay": map[string]any{"value": 300}, "forecast": map[string]any{"duration": "3h"}}); failed != "" {
		t.Fatal(failed)
	}
	board := filepath.Join(f.repo.KanbanDir(), workitem.BoardFile)
	before, _ := os.ReadFile(board)

	out, failed := f.call(t, "order_by_policy", map[string]any{"policy": "wsjf"})
	if failed != "" {
		t.Fatal(failed)
	}
	got := stories(out["stories"])
	if out["policy"] != "wsjf" || out["applied"] != false || strings.Join(got, " ") != valued.ID+" "+bare.ID {
		t.Fatalf("by wsjf, the valued story first and the bare one after: %v", out)
	}
	list := out["stories"].([]any)
	if first := list[0].(map[string]any); first["figure"] != 100.0 || first["unit"] != "per hour" {
		t.Errorf("wsjf is 300 over 3 hours: %v", first)
	}
	if second := list[1].(map[string]any); second["missing"] != "cost of delay value and forecast duration" {
		t.Errorf("the bare story names what it lacks: %v", second)
	}
	if out, _ := f.call(t, "order_by_policy", map[string]any{}); out["policy"] != "fifo" || strings.Join(stories(out["stories"]), " ") != bare.ID+" "+valued.ID {
		t.Errorf("no policy is the project's, fifo when it sets none: %v", out)
	}
	if _, failed := f.call(t, "order_by_policy", map[string]any{"policy": "random"}); !strings.Contains(failed, `unknown order policy "random"`) {
		t.Errorf("an unknown policy is refused: %q", failed)
	}
	if after, _ := os.ReadFile(board); string(after) != string(before) {
		t.Error("order_by_policy wrote board.md")
	}
}

// S-0217: promote_candidates lists the backlog stories that could go to
// ready and why each other one cannot, capped by limit.
func TestPromoteCandidatesListsTheBacklog(t *testing.T) {
	f := setup(t)
	backlog := func(title string, ready bool) *workitem.Item {
		t.Helper()
		s, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: f.story.Parent, Owner: "alex", Touches: []string{"docs/" + strings.ToLower(title)}, Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		if ready {
			data, _ := os.ReadFile(s.Path)
			body := strings.Replace(string(data), "## Goal\n", "## Goal\n\nDo it.\n", 1)
			body = strings.Replace(body, "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)
			_ = os.WriteFile(s.Path, []byte(body), 0o644)
			if _, failed := f.call(t, "item_edit", map[string]any{"id": s.ID, "cost_of_delay": map[string]any{"value": 50}, "forecast": map[string]any{"duration": "1h"}}); failed != "" {
				t.Fatal(failed)
			}
		}
		return s
	}
	one, two, bare := backlog("One", true), backlog("Two", true), backlog("Bare", false)

	out, failed := f.call(t, "promote_candidates", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["policy"] != "fifo" || strings.Join(stories(out["candidates"]), " ") != one.ID+" "+two.ID {
		t.Errorf("the candidates by fifo: %v", out)
	}
	others, _ := out["others"].([]any)
	if len(others) != 1 || others[0].(map[string]any)["id"] != bare.ID || !strings.Contains(strings.Join(toStrings(others[0].(map[string]any)["reasons"]), "; "), "no goal") {
		t.Errorf("the bare story with its reasons: %v", out["others"])
	}
	if out, _ := f.call(t, "promote_candidates", map[string]any{"limit": 1}); strings.Join(stories(out["candidates"]), " ") != one.ID {
		t.Errorf("limit 1 lists one candidate: %v", out)
	}
	if _, failed := f.call(t, "promote_candidates", map[string]any{"limit": -1}); !strings.Contains(failed, "limit is a number of candidates") {
		t.Errorf("a negative limit is refused: %q", failed)
	}
}

// S-0217: release_evaluate says whether the release policy is met, reading
// git for what is released, and says so when it cannot read git.
func TestReleaseEvaluateWeighsWhatIsPending(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	if _, failed := setup(t).call(t, "release_evaluate", map[string]any{}); !strings.Contains(failed, "on the host run flai release --evaluate") {
		t.Errorf("without a runner it says how to on the host: %q", failed)
	}

	f := setupWith(t, func(o *Options) { o.Runner = execx.System{} })
	root := f.repo.Root
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"commit", "-q", "-m", "init"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	out, failed := f.call(t, "release_evaluate", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["policy"] != "judgement" || out["met"] != false || out["count"] != 0.0 || out["reason"] == "" {
		t.Errorf("judgement, with nothing pending, is not met: %v", out)
	}
	count := 1
	f.repo.Manifest.Orchestration.Release.Policy, f.repo.Manifest.Orchestration.Release.Count = "threshold", &count
	if out, failed := f.call(t, "release_evaluate", map[string]any{}); failed != "" || out["policy"] != "threshold" || out["met"] != false || out["count_threshold"] != 1.0 ||
		!strings.Contains(out["reason"].(string), "no accepted story is waiting for a release") {
		t.Errorf("a threshold is not met by nothing pending: %v %s", out, failed)
	}
}

// S-0221: item_move refuses the orchestrator a move of a story to done, as it
// refuses every agent, and names the way it accepts one: flai accept with
// --by orchestrator, while accept_reviews is on (ADR-0093). An epic's
// refusal, and the story's agent's, are as before.
func TestItemMoveToDoneNamesTheOrchestratorsWay(t *testing.T) {
	t.Setenv("FLAI_ROLE", guard.RoleOrchestrate)
	f := setupWith(t, func(o *Options) { o.Agent = workitem.ActivityOrchestrator })
	_, failed := f.call(t, "item_move", map[string]any{"id": f.story.ID, "to": workitem.Done})
	want := "flai accept " + f.story.ID + " --by orchestrator --verified <commit> --evidence <file>, while orchestration.permissions.accept_reviews is on (ADR-0093)"
	if !strings.Contains(failed, want) || !strings.Contains(failed, "item_move never makes") {
		t.Errorf("orchestrator: %q, want it to name %q", failed, want)
	}
	if story, _ := f.repo.Get(f.story.ID); story.Status != workitem.InProgress {
		t.Errorf("the refused move must change nothing: %s", story.Status)
	}
	if _, failed := f.call(t, "item_move", map[string]any{"id": f.story.Parent, "to": workitem.Done}); !strings.Contains(failed, "only the operator") {
		t.Errorf("orchestrator, epic: %q", failed)
	}
	t.Setenv("FLAI_ROLE", "")
	agent := setup(t)
	if _, failed := agent.call(t, "item_move", map[string]any{"id": agent.story.ID, "to": workitem.Done}); !strings.Contains(failed, "only the operator") || strings.Contains(failed, "flai accept "+agent.story.ID) {
		t.Errorf("story's agent: %q", failed)
	}
}

// publishing is a project whose cli component was released at 1.0.0 and
// pushed to a bare remote, with S-0001 (tagged batch, worth 300 a week) and
// S-0002 (tagged batch, no value) of E-0001, in progress, accepted since,
// and S-0003 of E-0001 in the backlog. The orchestrator's flai mcp serves
// it with orchestration.permissions.publish on, publishing through a
// hostapi.Publisher whose flai release --pending is flaiPending.
type publishing struct {
	*fixture
	root, remote string
	push         bool            // the push host action
	journal      []hostapi.Entry // the host's journal
	ran          int             // runs of flai release --pending
}

func publishProject(t *testing.T, accepted bool) *publishing {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	pf := &publishing{root: t.TempDir(), remote: filepath.Join(t.TempDir(), "remote.git"), push: true}
	write := func(rel, body string) {
		t.Helper()
		p := filepath.Join(pf.root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range []string{"wip/kanban/tasks", "wip/agents", "wip/archive"} {
		write(filepath.Join(d, ".gitkeep"), "")
	}
	write("system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: cli\n    kind: go\norchestration:\n  permissions:\n    publish: true\n")
	write(".gitignore", ".flai-cache/\n")
	write("cli/main.go", "package main\n")
	write("wip/kanban/epics/E-0001-cli.md", "---\nid: E-0001\ntype: epic\nnature: feature\ntitle: CLI\nstatus: in-progress\n---\n# E-0001 CLI\n")
	write("wip/kanban/stories/S-0003-three.md", "---\nid: S-0003\ntype: story\nnature: feature\ntitle: Three\nstatus: backlog\nparent: E-0001\n---\n# S-0003 Three\n")
	gitIn(t, pf.root, "init", "-q", "-b", "main")
	gitIn(t, pf.root, "add", "-A")
	gitIn(t, pf.root, "commit", "-q", "-m", "init")
	gitIn(t, pf.root, "tag", "-a", "cli/v1.0.0", "-m", "cli 1.0.0")
	gitIn(t, filepath.Dir(pf.remote), "init", "-q", "--bare", "-b", "main", pf.remote)
	gitIn(t, pf.root, "remote", "add", "origin", pf.remote)
	gitIn(t, pf.root, "push", "-q", "-u", "origin", "main", "--tags")
	if accepted {
		// accepted as flai accept leaves a story: archived, done, committed
		for _, s := range []struct{ id, slug, value string }{{"S-0001", "one", "cost_of_delay:\n  value: 300\n  by: t\n  at: 2026-10-01T09:00:00Z\n"}, {"S-0002", "two", ""}} {
			write("wip/archive/kanban/stories/"+s.id+"-"+s.slug+".md", "---\nid: "+s.id+"\ntype: story\nnature: feature\ntitle: "+s.slug+"\nstatus: done\nparent: E-0001\ntags: [batch]\n"+s.value+"---\n# "+s.id+" "+s.slug+"\n")
			write("cli/"+s.slug+".go", "package main\n")
			gitIn(t, pf.root, "add", "-A")
			gitIn(t, pf.root, "commit", "-q", "-m", "chore: ["+s.id+"] accept and archive")
		}
	}
	repo, err := workitem.Open(pf.root)
	if err != nil {
		t.Fatal(err)
	}
	// flai serve sets the orchestrator's role in its session (S-0219)
	t.Setenv("FLAI_ROLE", guard.RoleOrchestrate)
	pb := &hostapi.Publisher{Run: pf.flaiPending(repo), Now: func() time.Time { return t0 }, Host: hostapi.Host{
		Enabled: func(action, root string) bool { return pf.push && action == hostapi.ActionPush && root == pf.root },
		Record:  func(e hostapi.Entry) { pf.journal = append(pf.journal, e) },
	}}
	srv := New(Options{Repo: repo, Agent: "claude", Version: "test", Now: func() time.Time { return t0 }, Runner: execx.System{}, Publish: pb})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "orchestrator", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	pf.fixture = &fixture{repo: repo, cs: cs}
	return pf
}

// flaiPending stands in for the flai release --pending that publish.run
// runs, which is flai's command and tested there (cmd's
// TestReleasePendingFromACloneMissingTheRemotesTags and
// TestReleasePendingFromACloneBehindItsRemoteBranch): it refuses with exit
// 3 and flai's own words, from release.CheckRemote, while the remote is
// ahead of the clone, and otherwise answers with the batch release.Pending
// finds, tagged and pushed, changing nothing.
func (pf *publishing) flaiPending(repo *workitem.Repo) hostapi.Runner {
	return func(_ context.Context, r hostapi.Run) (hostapi.Ran, error) {
		pf.ran++
		if got := strings.Join(r.Args, " "); got != "release --pending --json" || r.Dir != pf.root {
			return hostapi.Ran{Exit: 1, Events: []map[string]any{{"level": "FATAL", "err": "ran " + got + " in " + r.Dir}}}, nil
		}
		if refusal := release.CheckRemoteAfresh(execx.System{}, pf.root, repo.Manifest).Refusal(); refusal != "" {
			return hostapi.Ran{Exit: 3, Events: []map[string]any{{"level": "FATAL", "err": refusal}}}, nil
		}
		b, err := release.PendingBatch(execx.System{}, pf.root, repo.Manifest, repo)
		if err != nil {
			return hostapi.Ran{}, err
		}
		tags := []string{}
		for _, p := range b.Plans {
			tags = append(tags, p.Tag)
		}
		out, _ := json.Marshal(map[string]any{"plans": b.Plans, "tags": tags, "pushed": true, "published": []string{}})
		return hostapi.Ran{Stdout: out}, nil
	}
}

func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// state is what a refused publish leaves alone: the tags, the files, the
// head, and the remote.
func (pf *publishing) state(t *testing.T) string {
	t.Helper()
	return strings.Join([]string{gitIn(t, pf.root, "tag", "--list"), gitIn(t, pf.root, "status", "--porcelain"), gitIn(t, pf.root, "rev-parse", "HEAD"), gitIn(t, pf.remote, "show-ref")}, "\n--\n")
}

func (pf *publishing) policy(rel manifest.Release) { pf.repo.Manifest.Orchestration.Release = rel }

// S-0222: release_publish publishes the batch under threshold met, theme met,
// and judgement with a reason, through publish.run's flai release --pending,
// and returns the policy and its figures, the versions and tags, and the
// items bundled; the host's journal names the orchestrator and its reason.
func TestReleasePublishPublishesWhenThePolicyIsMet(t *testing.T) {
	value := 300.0
	for _, c := range []struct {
		name, reason, want string
		rel                manifest.Release
	}{
		{"threshold", "", "at or over the threshold of 300 USD/week", manifest.Release{Policy: manifest.ReleaseThreshold, Value: &value}},
		{"theme", "every story tagged batch is accepted", "every story tagged batch is accepted", manifest.Release{Policy: manifest.ReleaseTheme, Tag: "batch"}},
		{"judgement", "S-0001 and S-0002 make the batch command whole", "S-0001 and S-0002 make the batch command whole", manifest.Release{Policy: manifest.ReleaseJudgement}},
	} {
		t.Run(c.name, func(t *testing.T) {
			pf := publishProject(t, true)
			pf.policy(c.rel)
			out, failed := pf.call(t, "release_publish", map[string]any{"reason": c.reason})
			if failed != "" {
				t.Fatal(failed)
			}
			if out["policy"] != c.name || !strings.Contains(out["reason"].(string), c.want) || pf.ran != 1 {
				t.Errorf("the policy and the reason it published by: %v, ran %d", out, pf.ran)
			}
			if ev := out["evaluation"].(map[string]any); ev["count"] != 2.0 || ev["value"] != 300.0 || ev["policy"] != c.name {
				t.Errorf("the evaluation's figures: %v", ev)
			}
			versions, _ := out["versions"].([]any)
			if len(versions) != 1 {
				t.Fatalf("one component released: %v", out)
			}
			if v := versions[0].(map[string]any); v["component"] != "cli" || v["from"] != "1.0.0" || v["to"] != "1.1.0" || v["level"] != "minor" || strings.Join(stories(v["items"]), " ") != "S-0001 S-0002" {
				t.Errorf("cli 1.0.0 to 1.1.0 with both stories: %v", v)
			}
			if strings.Join(toStrings(out["tags"]), " ") != "cli/v1.1.0" || strings.Join(toStrings(out["items"]), " ") != "S-0001 S-0002" || out["pushed"] != true {
				t.Errorf("the tags and the items bundled: %v", out)
			}
			if len(pf.journal) != 1 {
				t.Fatalf("one run journalled: %+v", pf.journal)
			}
			if e := pf.journal[0]; e.Action != hostapi.ActionPush || e.Method != "mcp.release_publish" || e.By != "orchestrator (claude)" || e.Outcome != "done" ||
				!strings.HasPrefix(e.Detail, c.name+": ") || !strings.Contains(e.Detail, c.want) || !strings.HasSuffix(e.Detail, "pushed with tags cli/v1.1.0") {
				t.Errorf("the journal tells the orchestrator's release from the operator's: %+v", e)
			}
		})
	}
}

// S-0222: release_publish refuses, saying why and what would allow it, and
// changing no file, tag, or remote, under threshold and theme not met, a
// batch held under whole_epics, judgement with no reason, the push action
// off, and nothing accepted and unreleased; and to anyone but the
// orchestrator, and while publish is off.
func TestReleasePublishRefusesWhenThePolicyIsNotMet(t *testing.T) {
	value, count := 300.0, 3
	for _, c := range []struct {
		name     string
		rel      manifest.Release
		push     bool
		reason   string
		want     []string
		journals bool
	}{
		{"threshold not met", manifest.Release{Policy: manifest.ReleaseThreshold, Count: &count}, true, "", []string{"the release policy threshold is not met", "their count, 2, is under the threshold of 3", "It publishes once it is"}, false},
		{"theme not met", manifest.Release{Policy: manifest.ReleaseTheme, Epic: "E-0001"}, true, "E-0001 is done", []string{"the release policy theme is not met", "not yet accepted: S-0003"}, false},
		{"held by whole_epics", manifest.Release{Policy: manifest.ReleaseThreshold, Value: &value, WholeEpics: true}, true, "", []string{"whole_epics holds the batch back", "S-0001, whose epic E-0001 is in-progress; S-0002, whose epic E-0001 is in-progress", "once each of those epics is in review or done"}, false},
		{"held under judgement", manifest.Release{Policy: manifest.ReleaseJudgement, WholeEpics: true}, true, "coherent", []string{"whole_epics holds the batch back", "S-0001, whose epic E-0001"}, false},
		{"judgement with no reason", manifest.Release{Policy: manifest.ReleaseJudgement}, true, "  ", []string{"the release policy is judgement", "give reason, one sentence"}, false},
		{"push off", manifest.Release{Policy: manifest.ReleaseThreshold, Value: &value}, false, "", []string{"needs the push host action", `the host action "push" is not enabled`, "flai serve enable push"}, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			pf := publishProject(t, true)
			pf.policy(c.rel)
			pf.push = c.push
			before := pf.state(t)
			_, failed := pf.call(t, "release_publish", map[string]any{"reason": c.reason})
			for _, w := range c.want {
				if !strings.Contains(failed, w) {
					t.Errorf("the refusal says %q: %q", w, failed)
				}
			}
			if pf.ran != 0 {
				t.Errorf("a refusal runs no flai release --pending: ran %d", pf.ran)
			}
			if after := pf.state(t); after != before {
				t.Errorf("a refusal changes no file, tag, or remote:\n%s\nvs\n%s", before, after)
			}
			if journalled := len(pf.journal) == 1 && pf.journal[0].Outcome == "disabled" && pf.journal[0].By == "orchestrator (claude)"; journalled != c.journals || (!c.journals && len(pf.journal) != 0) {
				t.Errorf("only the push action's refusal is journalled, as publish.run's is: %+v", pf.journal)
			}
		})
	}

	t.Run("nothing pending", func(t *testing.T) {
		pf := publishProject(t, false)
		pf.policy(manifest.Release{Policy: manifest.ReleaseJudgement})
		before := pf.state(t)
		if _, failed := pf.call(t, "release_publish", map[string]any{"reason": "ready"}); !strings.Contains(failed, "nothing is accepted and not yet released") || pf.ran != 0 || pf.state(t) != before {
			t.Errorf("nothing to publish: %q, ran %d", failed, pf.ran)
		}
	})
	t.Run("publish off and not the orchestrator", func(t *testing.T) {
		pf := publishProject(t, true)
		pf.policy(manifest.Release{Policy: manifest.ReleaseThreshold, Value: &value})
		pf.repo.Manifest.Orchestration.Permissions.Publish = false
		if _, failed := pf.call(t, "release_publish", map[string]any{}); !strings.Contains(failed, "orchestration.permissions.publish, which is off") || pf.ran != 0 {
			t.Errorf("publish off: %q, ran %d", failed, pf.ran)
		}
		t.Setenv("FLAI_ROLE", "")
		if _, failed := setup(t).call(t, "release_publish", map[string]any{}); !strings.Contains(failed, "release_publish is the orchestrator's") || !strings.Contains(failed, "ADR-0067") {
			t.Errorf("a story's agent: %q", failed)
		}
	})
}

// S-0222: when flai release --pending refuses with exit 3, because the remote
// has release tags newer than the clone's (S-0174) or its branch moved,
// release_publish returns flai's message unchanged as a conflict, and
// nothing changes.
func TestReleasePublishReturnsTheRemotesRefusalAsAConflict(t *testing.T) {
	defer func(d time.Duration) { release.RemoteTTL = d }(release.RemoteTTL)
	release.RemoteTTL = 0
	value := 300.0
	for _, c := range []struct {
		name  string
		ahead func(t *testing.T, pf *publishing)
		want  string
	}{
		{"newer release tags", func(t *testing.T, pf *publishing) {
			gitIn(t, pf.remote, "tag", "cli/v1.4.0", "main") // published from another clone
		}, "cli/v1.4.0 (here cli/v1.0.0)"},
		{"moved remote", func(t *testing.T, pf *publishing) {
			other := filepath.Join(t.TempDir(), "other")
			gitIn(t, filepath.Dir(other), "clone", "-q", pf.remote, other)
			gitIn(t, other, "commit", "-q", "--allow-empty", "-m", "docs: from elsewhere")
			gitIn(t, other, "push", "-q", "origin", "main")
		}, "origin/main has commits this clone lacks"},
	} {
		t.Run(c.name, func(t *testing.T) {
			pf := publishProject(t, true)
			pf.policy(manifest.Release{Policy: manifest.ReleaseThreshold, Value: &value})
			c.ahead(t, pf)
			before := pf.state(t)
			refusal := release.CheckRemoteAfresh(execx.System{}, pf.root, pf.repo.Manifest).Refusal()
			_, failed := pf.call(t, "release_publish", map[string]any{})
			if refusal == "" || failed != refusal || !strings.HasPrefix(failed, "conflict: refusing to publish: ") || !strings.Contains(failed, c.want) {
				t.Errorf("flai's refusal, unchanged, as a conflict:\n got %q\nwant %q, naming %q", failed, refusal, c.want)
			}
			if pf.ran != 1 || pf.state(t) != before {
				t.Errorf("the refused publish changes nothing: ran %d", pf.ran)
			}
			if len(pf.journal) != 1 || pf.journal[0].Outcome != "failed" || pf.journal[0].By != "orchestrator (claude)" || !strings.Contains(pf.journal[0].Detail, c.want) {
				t.Errorf("the refusal is journalled: %+v", pf.journal)
			}
		})
	}
}

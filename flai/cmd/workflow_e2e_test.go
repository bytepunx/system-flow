package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// e2eBoard is the template's workflow file with testing, of the active kind
// and with a WIP limit, added between in-progress and review, has-task on
// the move into it, and the timer testing-to-accept from it to the
// acceptance (S-0756).
const e2eBoard = "../internal/workflow/testdata/e2e-board.md"

// e2eProject is a project rendered from the template that runs e2eBoard,
// and the clock its commands run at, an hour on at each, so that the time
// an item spends in a state is never nothing.
type e2eProject struct {
	root string
	now  time.Time
	reqs int
}

// at is the next command's time.
func (p *e2eProject) at() time.Time {
	p.now = p.now.Add(time.Hour)
	return p.now
}

// run runs flai in the project at the next hour.
func (p *e2eProject) run(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	return runInAt(t, p.root, p.at(), args...)
}

// must runs flai in the project and fails the test on a non-zero exit.
func (p *e2eProject) must(t *testing.T, args ...string) string {
	t.Helper()
	out, errOut, code := p.run(t, args...)
	if code != 0 {
		t.Fatalf("flai %s: exit %d\n%s%s", strings.Join(args, " "), code, out, errOut)
	}
	return out
}

// status is the status of the item id, archived or not.
func (p *e2eProject) status(t *testing.T, id string) string {
	t.Helper()
	repo, err := workitem.Open(p.root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	return it.Status
}

// itemMove moves id to to with the MCP tool item_move, as the agent of the
// story story calls it over flai mcp, and says why it failed, or "".
func (p *e2eProject) itemMove(t *testing.T, story, id, to string) string {
	t.Helper()
	repo, err := workitem.Open(p.root)
	if err != nil {
		t.Fatal(err)
	}
	at := p.at()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct, st := mcp.NewInMemoryTransports()
	if _, err := mcpserver.New(mcpserver.Options{Repo: repo, Agent: "agent-" + story, Version: "test", Now: func() time.Time { return at }}).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "agent-" + story, Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "item_move", Arguments: map[string]any{"id": id, "to": to}})
	if err != nil {
		t.Fatalf("item_move %s %s: %v", id, to, err)
	}
	if !res.IsError {
		return ""
	}
	failed := ""
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			failed += tc.Text
		}
	}
	return cmpOr(failed, "failed")
}

// hostMove moves id to to with the host API's item.move, as the dashboard
// asks for it, and says why it was refused, or "".
func (p *e2eProject) hostMove(t *testing.T, id, to string) string {
	t.Helper()
	p.reqs++
	params := fmt.Sprintf(`{"id":%q,"to":%q,"request_id":"3f0c1a52-7d3b-4f0e-9a51-%012d"}`, id, to, p.reqs)
	out, errOut, code := p.run(t, "hostapi", "item.move", params)
	if code == 0 && !strings.Contains(out, `"error"`) {
		return ""
	}
	return cmpOr(strings.TrimSpace(out+errOut), "refused")
}

// cmpOr is s, or else when s is empty.
func cmpOr(s, otherwise string) string {
	if s == "" {
		return otherwise
	}
	return s
}

// commandsAt runs the flai commands behind the host API's writes in this
// process, at the project's clock: under go test the executable is the test
// binary, not flai.
func commandsAt(t *testing.T, p *e2eProject) {
	t.Helper()
	was := hostapi.Commands
	t.Cleanup(func() { hostapi.Commands = was })
	t.Setenv("LOG_FORMAT", "json")
	hostapi.Commands = func(_ context.Context, r hostapi.Run) (hostapi.Ran, error) {
		out, errOut, code := runInAt(t, r.Dir, p.at(), r.Args...)
		ran := hostapi.Ran{Stdout: []byte(out), Exit: code}
		for _, line := range strings.Split(errOut, "\n") {
			var ev map[string]any
			if json.Unmarshal([]byte(line), &ev) == nil {
				ran.Events = append(ran.Events, ev)
			}
		}
		return ran, nil
	}
}

// e2eLane is a lane of flai board --json and board.get.
type e2eLane struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	WIPLimit   int    `json:"wip_limit"`
	TakesLimit bool   `json:"takes_limit"`
}

// e2eStats is the part of flai stats --json and stats.get the test reads.
type e2eStats struct {
	Summary struct {
		Timers       map[string]struct{ Count int } `json:"timers"`
		InStateShare map[string]float64             `json:"time_in_state_share"`
	} `json:"summary"`
	Items []struct {
		ID      string             `json:"id"`
		Timers  map[string]float64 `json:"timers"`
		InState map[string]float64 `json:"time_in_state_seconds"`
	} `json:"items"`
	CFD []struct {
		Counts map[string]int `json:"counts"`
	} `json:"cfd"`
	Timers []struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	} `json:"timers"`
}

// writeCriterion gives the story id a ticked
// acceptance criterion in place of the template's empty one.
func (p *e2eProject) writeCriterion(t *testing.T, id string) {
	t.Helper()
	files, _ := filepath.Glob(filepath.Join(p.root, "wip/kanban/stories", id+"-*.md"))
	if len(files) != 1 {
		t.Fatalf("%s's file: %v", id, files)
	}
	b, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	i := strings.Index(s, "## Acceptance criteria\n")
	if i < 0 {
		t.Fatalf("%s has no acceptance criteria:\n%s", id, s)
	}
	j := strings.Index(s[i:], "- [ ]")
	if j < 0 {
		t.Fatalf("%s has no checkbox under its acceptance criteria:\n%s", id, s)
	}
	j += i
	end := strings.IndexByte(s[j:], '\n') + j
	if err := os.WriteFile(files[0], []byte(s[:j]+"- [x] it works"+s[end:]), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0756 criteria 1 to 7, the Go side: a state added in the workflow file of
// a project rendered from the template, testing between in-progress and
// review, is followed by flai check, flai board and board.get, flai move,
// the MCP tool item_move, and the host API's item.move, flai serve's agent
// restart, flai accept, and flai stats --json and stats.get, each driven by
// the file alone. It needs real git and the template, so it skips under
// -short and runs in the integration tier, which the close-out of a story
// that changes flai/internal/workflow runs (criterion 7).
func TestWorkflowE2EAStateAddedInTheTemplatesWorkflowFileIsFollowed(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: renders the template into a project with real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	tpl := filepath.Join("..", "..", "template")
	if _, err := os.Stat(filepath.Join(tpl, "template.yaml")); err != nil {
		t.Skip("template not present")
	}
	tpl, _ = filepath.Abs(tpl)
	board, err := os.ReadFile(e2eBoard)
	if err != nil {
		t.Fatal(err)
	}
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	asOperator(t)
	t.Setenv("FLAI_AGENT", "olive")
	p := &e2eProject{root: filepath.Join(t.TempDir(), "e2e"), now: time.Date(2026, 9, 15, 8, 0, 0, 0, time.UTC)}
	commandsAt(t, p)

	ok := t.Run("1 the template renders and its workflow file takes testing between in-progress and review", func(t *testing.T) {
		if _, errOut, code := runInAt(t, ".", p.at(), "new", p.root, "--template", tpl, "--defaults", "--no-git", "--var", "description=End to end", "--var", "owner=olive"); code != 0 {
			t.Fatalf("render the template: %s", errOut)
		}
		file := filepath.Join(p.root, "design", "workflows", "board.md")
		rendered, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("the template renders no workflow file: %v", err)
		}
		// the fixture is the template's board with testing added: the same
		// states otherwise, in the same order, with the same kinds
		base, prob := workflow.Parse("design/workflows/board.md", string(rendered))
		if prob != nil {
			t.Fatalf("the template's workflow file: %s", prob)
		}
		added, prob := workflow.Parse("design/workflows/board.md", string(board))
		if prob != nil {
			t.Fatalf("the fixture: %s", prob)
		}
		var rest []string
		for _, s := range added.States {
			if s.Name != "testing" {
				rest = append(rest, s.Name+"/"+s.Kind)
			}
		}
		var want []string
		for _, s := range base.States {
			want = append(want, s.Name+"/"+s.Kind)
		}
		if !slices.Equal(rest, want) {
			t.Errorf("the fixture's states but testing are %v; the template's are %v", rest, want)
		}
		if err := os.WriteFile(file, board, 0o644); err != nil {
			t.Fatal(err)
		}

		repo, err := workitem.Open(p.root)
		if err != nil {
			t.Fatalf("the project does not open on the file: %v", err)
		}
		wf := repo.Flow()
		if got := wf.Names(); !slices.Equal(got, []string{"backlog", "ready", "in-progress", "testing", "review", "done", "cancelled"}) {
			t.Errorf("the project runs the states %v", got)
		}
		if st, ok := wf.State("testing"); !ok || st.Kind != workflow.KindActive || !st.Limit {
			t.Errorf("testing: %+v %v, want the active kind with a limit", st, ok)
		}
		tr, ok := wf.Find(workitem.Story, "in-progress", "testing")
		if !ok || len(tr.Rules) != 1 || tr.Rules[0].Rule != workflow.RuleHasTask {
			t.Errorf("the move into testing: %+v %v, want has-task", tr, ok)
		}
		if _, ok := wf.Find(workitem.Story, "in-progress", "review"); ok {
			t.Error("a story goes to review only through testing")
		}
		if tm, ok := wf.Timer("testing-to-accept"); !ok || tm.From == nil || tm.From.Edge != workflow.EdgeFirst || tm.From.State != "testing" || tm.To == nil || tm.To.At != workflow.PointCompleted {
			t.Errorf("the timer testing-to-accept: %+v %v", tm, ok)
		}

		p.must(t, "config", "set", "author", "olive")
		p.must(t, "agent", "set", "--harness", "command")
		if out := p.must(t, "board", "limit", "testing", "2"); !strings.Contains(out, "testing: WIP limit 2") {
			t.Errorf("flai board limit testing 2: %s", out)
		}
		gitIn(t, p.root, "init", "-q", "-b", "main")
		gitIn(t, p.root, "config", "user.email", "t@t")
		gitIn(t, p.root, "config", "user.name", "t")
		gitIn(t, p.root, "add", "-A")
		gitIn(t, p.root, "commit", "-q", "-m", "chore: render the template with testing added")
	})
	if !ok {
		return
	}

	t.Run("2 flai check --strict passes, and names testing when it is left unreachable", func(t *testing.T) {
		if out, errOut, code := p.run(t, "check", "--strict"); code != 0 {
			t.Fatalf("flai check --strict on the project: exit %d\n%s%s", code, out, errOut)
		}
		file := filepath.Join(p.root, "design", "workflows", "board.md")
		into := "  - {from: in-progress, to: testing, rules: [has-task]}\n"
		if !strings.Contains(string(board), into) {
			t.Fatalf("the fixture lost its move into testing, %q", into)
		}
		if err := os.WriteFile(file, []byte(strings.Replace(string(board), into, "", 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.WriteFile(file, board, 0o644); err != nil {
				t.Fatal(err)
			}
		}()
		out, errOut, code := p.run(t, "check", "--strict")
		if code == 0 || !strings.Contains(out+errOut, workflow.CodeUnreachable) || !strings.Contains(out+errOut, "state testing is entered by no chain of transitions") {
			t.Errorf("flai check --strict must fail naming testing unreachable: exit %d\n%s%s", code, out, errOut)
		}
	})

	t.Run("3 flai board and board.get show testing in its place with its kind and WIP limit", func(t *testing.T) {
		out := p.must(t, "board")
		order := []string{"in-progress", "testing (0/2)", "review"}
		at := -1
		for _, want := range order {
			i := slices.IndexFunc(strings.Split(out, "\n"), func(l string) bool { return strings.HasPrefix(l, want) })
			if i <= at {
				t.Errorf("flai board must show %q after the column before it:\n%s", want, out)
			}
			at = i
		}
		lanes := func(what, js string) {
			t.Helper()
			var v struct {
				Lanes []e2eLane `json:"lanes"`
			}
			if err := json.Unmarshal([]byte(js), &v); err != nil {
				t.Fatalf("%s: %v\n%s", what, err, js)
			}
			i := slices.IndexFunc(v.Lanes, func(l e2eLane) bool { return l.Name == "testing" })
			if i < 1 || i+1 >= len(v.Lanes) || v.Lanes[i-1].Name != "in-progress" || v.Lanes[i+1].Name != "review" {
				t.Errorf("%s must lay testing out between in-progress and review: %+v", what, v.Lanes)
				return
			}
			if l := v.Lanes[i]; l.Kind != workflow.KindActive || l.WIPLimit != 2 || !l.TakesLimit {
				t.Errorf("%s's testing lane is %+v, want the active kind and the limit 2", what, l)
			}
		}
		lanes("flai board --json", p.must(t, "board", "--json"))
		lanes("board.get", p.must(t, "hostapi", "board.get", "{}"))
	})

	t.Run("4 flai move, item_move, and item.move take a story through testing and refuse the same moves", func(t *testing.T) {
		p.must(t, "epic", "new", "Epic")
		p.must(t, "story", "new", "By flai move", "--epic", "E-0001")
		p.must(t, "story", "new", "By item_move")
		p.must(t, "story", "new", "By item.move")
		for _, id := range []string{"S-0001", "S-0002", "S-0003"} {
			p.writeCriterion(t, id)
		}
		// each surface makes a move, saying why it failed, or ""
		surfaces := []struct {
			name, story string
			move        func(id, to string) string
		}{
			{"flai move", "S-0001", func(id, to string) string {
				out, errOut, code := p.run(t, "move", id, to)
				if code == 0 {
					return ""
				}
				return cmpOr(out+errOut, "refused")
			}},
			{"item_move", "S-0002", func(id, to string) string { return p.itemMove(t, "S-0002", id, to) }},
			{"item.move", "S-0003", func(id, to string) string { return p.hostMove(t, id, to) }},
		}
		for _, s := range surfaces {
			moves := func(to string) {
				t.Helper()
				if why := s.move(s.story, to); why != "" {
					t.Errorf("%s must move %s to %s: %s", s.name, s.story, to, why)
				}
				if got := p.status(t, s.story); got != to {
					t.Errorf("after %s to %s, %s is %s", s.name, to, s.story, got)
				}
			}
			refuses := func(to, want string) {
				t.Helper()
				before := p.status(t, s.story)
				if why := s.move(s.story, to); why == "" || !strings.Contains(why, want) {
					t.Errorf("%s must refuse %s's move from %s to %s saying %q: %q", s.name, s.story, before, to, want, why)
				}
				if got := p.status(t, s.story); got != before {
					t.Errorf("a refused move changes nothing: %s is %s, want %s", s.story, got, before)
				}
			}
			moves("ready")
			moves("in-progress")
			if s.story == "S-0001" {
				// its branch and worktree, which the acceptance merges
				p.must(t, "stream", "open", "S-0001")
			}
			// no transition leads from in-progress to review on this board
			refuses("review", "cannot go from in-progress to review")
			// has-task refuses the move into testing until the story has a task
			refuses("testing", "needs at least one task")
			p.must(t, "task", "new", "--story", s.story, "Step")
			moves("testing")
			moves("review")
			if s.story == "S-0001" {
				if got := p.status(t, "E-0001"); got != "review" {
					t.Errorf("E-0001 follows its only story to review, through testing: it is %s", got)
				}
			}
		}
		// a story in testing goes back to in-progress, and no further back
		p.must(t, "story", "new", "Kept")
		p.writeCriterion(t, "S-0004")
		p.must(t, "task", "new", "--story", "S-0004", "Step")
		for _, to := range []string{"ready", "in-progress", "testing", "in-progress", "testing"} {
			p.must(t, "move", "S-0004", to)
		}
		for _, s := range surfaces {
			if why := s.move("S-0004", "ready"); why == "" || !strings.Contains(why, "cannot go from testing to ready") {
				t.Errorf("%s must refuse S-0004's move from testing to ready: %q", s.name, why)
			}
		}
		if got := p.status(t, "S-0004"); got != "testing" {
			t.Errorf("S-0004 is %s, want testing", got)
		}
		// acceptance is no agent's: item_move refuses a story in review the move to done
		if why := p.itemMove(t, "S-0002", "S-0002", "done"); !strings.Contains(why, "moving it to done is acceptance") {
			t.Errorf("item_move must refuse S-0002's move to done as acceptance: %q", why)
		}
		// the tasks are done from in-progress, as the acceptance needs
		for _, id := range []string{"T-0001", "T-0002", "T-0003"} {
			for _, to := range []string{"ready", "in-progress", "done"} {
				p.must(t, "move", id, to)
			}
		}
	})

	t.Run("5 flai serve starts the agent of a story in testing again, and flai accept merges and archives a story that went through it", func(t *testing.T) {
		p.must(t, "serve", "enable", "agent")
		p.must(t, "serve", "agent", "set", "--", "true", "{story}")
		p.must(t, "serve", "agent", "harness", "command", "--deny-protected")
		// review is waiting: its agent's run has ended with its work done
		if _, errOut, code := p.run(t, "serve", "agent", "restart", "S-0002"); code == 0 || !strings.Contains(errOut, "S-0002 is in review; only a story in ready, in-progress, or testing is restarted") {
			t.Errorf("flai serve must restart no agent for a story in review, naming testing among the states it does: exit %d %s", code, errOut)
		}
		// testing is active work: flai serve keeps an agent on it
		if out, errOut, code := p.run(t, "serve", "agent", "restart", "S-0004"); code != 0 || !strings.HasPrefix(out, "started true (command) for S-0004 as agent-S-0004") {
			t.Errorf("flai serve must start the agent of a story in testing: exit %d %s%s", code, out, errOut)
		}

		gitIn(t, p.root, "add", "-A")
		gitIn(t, p.root, "commit", "-q", "-m", "chore: the stories through testing")
		wt := filepath.Join(p.root, ".flai-cache", "worktrees", "S-0001")
		if err := os.WriteFile(filepath.Join(wt, "design", "tested.md"), []byte("the story's work\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		gitIn(t, wt, "add", "-A")
		gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] the story's work")
		out, errOut, code := p.run(t, "accept", "S-0001", "--by", "olive", "--trailer", "Co-Authored-By: t <t@t>")
		if code != 0 {
			t.Fatalf("flai accept S-0001: exit %d\n%s%s", code, out, errOut)
		}
		if !strings.Contains(out, "accepted S-0001: done") || !strings.Contains(out, "E-0001 followed it from review to done") {
			t.Errorf("the acceptance's summary:\n%s", out)
		}
		files, _ := filepath.Glob(filepath.Join(p.root, "wip/archive/kanban/stories/S-0001-*.md"))
		if len(files) != 1 {
			t.Fatalf("S-0001 is not archived: %v", files)
		}
		if b, _ := os.ReadFile(files[0]); !strings.Contains(string(b), "status: done") {
			t.Error("the archived S-0001 is done")
		}
		if log := gitIn(t, p.root, "log", "-1", "--format=%s"); !strings.HasPrefix(log, "chore: [S-0001] accept and archive") {
			t.Errorf("the acceptance's commit: %q", log)
		}
		if got := gitIn(t, p.root, "show", "main:design/tested.md"); got != "the story's work" {
			t.Errorf("main must carry the story's branch after the acceptance: %q", got)
		}
		if gitIn(t, p.root, "branch", "--list", "story/S-0001") != "" {
			t.Error("the accepted story's branch is removed")
		}
		if st := gitIn(t, p.root, "status", "--porcelain"); strings.TrimSpace(st) != "" {
			t.Errorf("the working tree is dirty after the acceptance:\n%s", st)
		}
	})

	t.Run("6 flai stats --json and stats.get report testing-to-accept for the story, and time in state stacks testing", func(t *testing.T) {
		report := func(what, js string) {
			t.Helper()
			var r e2eStats
			if err := json.Unmarshal([]byte(js), &r); err != nil {
				t.Fatalf("%s: %v\n%s", what, err, js)
			}
			if !slices.ContainsFunc(r.Timers, func(d struct {
				Name  string `json:"name"`
				Title string `json:"title"`
			}) bool {
				return d.Name == "testing-to-accept" && d.Title == "Testing to acceptance"
			}) {
				t.Errorf("%s must list the timer testing-to-accept: %+v", what, r.Timers)
			}
			i := slices.IndexFunc(r.Items, func(it struct {
				ID      string             `json:"id"`
				Timers  map[string]float64 `json:"timers"`
				InState map[string]float64 `json:"time_in_state_seconds"`
			}) bool {
				return it.ID == "S-0001"
			})
			if i < 0 {
				t.Fatalf("%s reports no S-0001", what)
			}
			it := r.Items[i]
			// from the first move into testing to the acceptance is the time in
			// testing and in review, for a story that never went back
			want := it.InState["testing"] + it.InState["review"]
			if got, ok := it.Timers["testing-to-accept"]; !ok || got <= 0 || got != want {
				t.Errorf("%s gives S-0001 testing-to-accept %v (%v), want %v, its time in testing and review", what, got, ok, want)
			}
			if it.InState["testing"] <= 0 {
				t.Errorf("%s gives S-0001 no time in testing: %v", what, it.InState)
			}
			if c := r.Summary.Timers["testing-to-accept"].Count; c != 1 {
				t.Errorf("%s's summary counts %d stories under testing-to-accept, want 1", what, c)
			}
			if r.Summary.InStateShare["testing"] <= 0 {
				t.Errorf("%s's time in state share has no testing: %v", what, r.Summary.InStateShare)
			}
			if !slices.ContainsFunc(r.CFD, func(d struct {
				Counts map[string]int `json:"counts"`
			}) bool {
				return d.Counts["testing"] > 0
			}) {
				t.Errorf("%s's cumulative flow never stacks testing", what)
			}
		}
		report("flai stats --json", p.must(t, "stats", "--json"))
		// stats.get answers as a read the host runs: the report under data
		var got struct {
			Data json.RawMessage `json:"data"`
		}
		js := p.must(t, "hostapi", "stats.get", "{}")
		if err := json.Unmarshal([]byte(js), &got); err != nil || len(got.Data) == 0 {
			t.Fatalf("stats.get answers no data: %v\n%s", err, js)
		}
		report("stats.get", string(got.Data))
	})
}

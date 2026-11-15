package workitem

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

const readyBody = "## Goal\n\nDo it.\n\n## Acceptance criteria\n- [ ] it works\n\n## Notes\n"

// promotable is a backlog story under E-0001 that is a candidate unless the
// test changes it.
func promotable(id string, value float64, duration string, touches ...string) *Item {
	return &Item{
		ID: id, Type: Story, Status: Backlog, Parent: "E-0001", Title: "Story " + id,
		Created: "2026-09-01T10:00:00Z", Body: readyBody, Touches: touches,
		CostOfDelay: &CostOfDelay{Value: amount(value)}, Forecast: &Forecast{Duration: duration},
	}
}

func TestPromotionCandidates(t *testing.T) {
	epic := &Item{ID: "E-0001", Type: Epic, Status: InProgress, Title: "Epic"}
	gone := &Item{ID: "E-0002", Type: Epic, Status: Cancelled, Title: "Gone"}
	open := promotable("S-0001", 0, "1h", "flai/cmd")
	open.Status = InProgress

	good := promotable("S-0002", 100, "2h", "docs")
	better := promotable("S-0003", 900, "3h", "template")
	draft := promotable("S-0004", 100, "1h", "a")
	draft.Draft = true
	noGoal := promotable("S-0005", 100, "1h", "b")
	noGoal.Body = "## Goal\n\n## Acceptance criteria\n- [ ]\n"
	orphan := promotable("S-0006", 100, "1h", "c")
	orphan.Parent = "E-0002"
	held := promotable("S-0007", 100, "1h", "flai/cmd/move.go")
	waits := promotable("S-0008", 100, "1h", "d")
	waits.After = []string{"S-0001"}
	noForecast := promotable("S-0009", 100, "", "e")
	noForecast.Forecast = nil
	badForecast := promotable("S-0010", 100, "soon", "f")
	noValue := promotable("S-0011", 0, "1h", "g")
	noValue.CostOfDelay = nil
	ready := promotable("S-0012", 100, "1h", "h")
	ready.Status = Ready

	items := []*Item{epic, gone, open, good, better, draft, noGoal, orphan, held, waits, noForecast, badForecast, noValue, ready}
	order := []string{"S-0012", "S-0011", "S-0010", "S-0009", "S-0008", "S-0007", "S-0006", "S-0005", "S-0004", "S-0003", "S-0002"}
	got, err := PromotionCandidates(items, order, NewHolds(items, nil), PolicyCOD, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got.Policy != PolicyCOD || len(got.Candidates) != 2 || got.Candidates[0].ID != "S-0003" || got.Candidates[1].ID != "S-0002" ||
		got.Candidates[0].Position != 1 || got.Candidates[0].Text != "900" {
		t.Fatalf("candidates by cod: %+v", got.Candidates)
	}
	want := map[string][]string{
		"S-0011": {"no cost of delay value"},
		"S-0010": {"no forecast duration"},
		"S-0009": {"no forecast duration"},
		"S-0008": {"held (after): waits for S-0001 (in progress); starts when S-0001 is done"},
		"S-0007": {"held (overlap): touches flai/cmd/move.go, inside flai/cmd which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		"S-0006": {"its epic E-0002 is cancelled"},
		"S-0005": {"no goal", "no acceptance criteria with a checkbox"},
		"S-0004": {"draft: finalize it first"},
	}
	var ids []string
	for _, o := range got.Others {
		ids = append(ids, o.ID)
		if !reflect.DeepEqual(o.Reasons, want[o.ID]) {
			t.Errorf("%s: reasons %q, want %q", o.ID, o.Reasons, want[o.ID])
		}
	}
	if wantIDs := []string{"S-0011", "S-0010", "S-0009", "S-0008", "S-0007", "S-0006", "S-0005", "S-0004"}; !reflect.DeepEqual(ids, wantIDs) {
		t.Errorf("others in backlog pull order, ready stories left out: got %v want %v", ids, wantIDs)
	}
	for _, o := range got.Others {
		if (o.ID == "S-0007" || o.ID == "S-0008") != (o.Held != nil) {
			t.Errorf("%s: held %+v", o.ID, o.Held)
		}
	}

	// Every reason at once, and a story without touches while one is open.
	bare := &Item{ID: "S-0013", Type: Story, Status: Backlog, Parent: "E-0001", Title: "Bare", Draft: true}
	got, err = PromotionCandidates([]*Item{epic, open, bare}, nil, NewHolds([]*Item{epic, open, bare}, nil), PolicyFIFO, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Candidates) != 0 || len(got.Others) != 1 {
		t.Fatalf("bare: %+v", got)
	}
	r := got.Others[0].Reasons
	if len(r) != 6 || r[0] != "draft: finalize it first" || r[1] != "no goal" || !strings.HasPrefix(r[3], "held (no-touches)") ||
		r[4] != "no forecast duration" || r[5] != "no cost of delay value" {
		t.Errorf("bare's reasons: %q", r)
	}
}

func TestPromotionCandidatesPolicyAndLimit(t *testing.T) {
	epic := &Item{ID: "E-0001", Type: Epic, Status: Ready, Title: "Epic"}
	items := []*Item{epic,
		promotable("S-0001", 100, "1h"),
		promotable("S-0002", 900, "10h"),
		promotable("S-0003", 300, "30m"),
	}
	holds := NewHolds(items, nil)
	for _, c := range []struct {
		policy string
		ids    []string
	}{
		{PolicyCOD, []string{"S-0002", "S-0003", "S-0001"}},
		{PolicyWSJF, []string{"S-0003", "S-0001", "S-0002"}},
		{PolicyThroughput, []string{"S-0003", "S-0001", "S-0002"}},
		{PolicyFIFO, []string{"S-0001", "S-0002", "S-0003"}},
	} {
		got, err := PromotionCandidates(items, nil, holds, c.policy, 0)
		if err != nil {
			t.Fatal(err)
		}
		var ids []string
		for _, r := range got.Candidates {
			ids = append(ids, r.ID)
		}
		if !reflect.DeepEqual(ids, c.ids) {
			t.Errorf("%s: got %v want %v", c.policy, ids, c.ids)
		}
	}

	got, err := PromotionCandidates(items, nil, holds, PolicyCOD, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Candidates) != 2 || got.Candidates[1].ID != "S-0003" || len(got.Others) != 0 {
		t.Errorf("a limit caps the candidates and lists no others: %+v", got)
	}

	if _, err := PromotionCandidates(items, nil, holds, "random", 0); err == nil {
		t.Error("an unknown policy is an error")
	}

	// Empty lists read as [] in JSON, not null.
	got, err = PromotionCandidates([]*Item{epic}, nil, holds, PolicyFIFO, 0)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(got)
	if string(data) != `{"policy":"fifo","candidates":[],"others":[]}` {
		t.Errorf("json: %s", data)
	}
}

// promoteKindRenamedFile is the fixture with every state but review renamed:
// backlog icebox, ready queued, done shipped, and cancelled dropped, besides
// in-progress doing, so that code which reads a state's kind, a transition's
// effect, or an order finds each role, and code which names a state finds
// none (S-0739).
var promoteKindRenamedFile = strings.NewReplacer(
	"backlog", "icebox", "ready", "queued",
	"{name: done,", "{name: shipped,", "to: done", "to: shipped",
	"{name: cancelled,", "{name: dropped,", "to: cancelled", "to: dropped", "from: cancelled", "from: dropped",
).Replace(testFlowFile)

// promoteKindRenamedState is the renamed fixture's name for the baseline
// state named state.
func promoteKindRenamedState(state string) string {
	if name, ok := map[string]string{Backlog: "icebox", Ready: "queued", Done: "shipped", Cancelled: "dropped"}[state]; ok {
		return name
	}
	return inTestFlow(state)
}

// promoteKindFlow is one workflow the cases of the files that read a state's
// kind run under: the workflow, nil for the baseline, its name for each
// baseline state, and the name of an active state after the pull's target,
// in-progress itself for the baseline.
type promoteKindFlow struct {
	name  string
	flow  func(*testing.T) *workflow.Workflow
	state func(string) string
	added string
}

// of is the workflow f runs.
func (f promoteKindFlow) of(t *testing.T) workflow.Workflow {
	if w := f.flow(t); w != nil {
		return *w
	}
	return workflow.Baseline()
}

// open is a new project running f's workflow.
func (f promoteKindFlow) open(t *testing.T) *Repo {
	r := newProject(t)
	r.Workflow = f.flow(t)
	return r
}

// promoteKindFlows are the baseline, the fixture, and the fixture with every
// state but review renamed, for cases that must give the same result under
// each wherever their kinds agree (S-0739, ADR-0301 §2).
func promoteKindFlows() []promoteKindFlow {
	return []promoteKindFlow{
		{"baseline", func(*testing.T) *workflow.Workflow { return nil }, func(s string) string { return s }, InProgress},
		{"fixture", testFlow, inTestFlow, fixtureTesting},
		{"renamed", func(t *testing.T) *workflow.Workflow {
			t.Helper()
			w, p := workflow.Parse("design/workflows/board.md", promoteKindRenamedFile)
			if p != nil {
				t.Fatalf("the renamed workflow does not parse: %s", p)
			}
			return &w
		}, promoteKindRenamedState, fixtureTesting},
	}
}

// S-0739, ADR-0301 §3: a promotion is the commit transition, so the
// candidates are the stories in the state it leaves, an epic is cancelled by
// its state's kind, and a story anywhere else is left out, under every
// workflow alike.
func TestPromoteKindCandidatesFollowTheWorkflow(t *testing.T) {
	for _, f := range promoteKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			w, st := f.of(t), f.state
			if from, to := promoteSource(w), promoteTarget(w); from != st(Backlog) || to != st(Ready) {
				t.Errorf("promotion from %q to %q, want %q to %q", from, to, st(Backlog), st(Ready))
			}
			in := func(it *Item, state string) *Item {
				it.Status = state
				return it
			}
			epic := &Item{ID: "E-0001", Type: Epic, Status: st(InProgress), Title: "Epic"}
			gone := &Item{ID: "E-0002", Type: Epic, Status: st(Cancelled), Title: "Gone"}
			orphan := in(promotable("S-0006", 100, "1h", "c"), st(Backlog))
			orphan.Parent = "E-0002"
			items := []*Item{epic, gone, orphan,
				in(promotable("S-0001", 0, "1h", "flai/cmd"), st(InProgress)),
				in(promotable("S-0002", 100, "2h", "docs"), st(Backlog)),
				in(promotable("S-0003", 900, "3h", "template"), st(Backlog)),
				in(promotable("S-0004", 0, "1h", "flai/serve"), f.added),
				in(promotable("S-0005", 0, "1h", "flai/mcp"), st(Review)),
				in(promotable("S-0007", 100, "1h", "h"), st(Ready)),
				in(promotable("S-0008", 100, "1h", "i"), st(Done)),
			}
			got, err := PromotionCandidatesIn(w, items, []string{"S-0006", "S-0002", "S-0003"}, NewHolds(items, nil), PolicyCOD, 0)
			if err != nil {
				t.Fatal(err)
			}
			var ids []string
			for _, c := range got.Candidates {
				ids = append(ids, c.ID)
			}
			want := []Refused{{ID: "S-0006", Title: "Story S-0006", Reasons: []string{"its epic E-0002 is " + st(Cancelled)}}}
			if !reflect.DeepEqual(ids, []string{"S-0003", "S-0002"}) || !reflect.DeepEqual(got.Others, want) {
				t.Errorf("candidates %v, others %+v; want [S-0003 S-0002], %+v", ids, got.Others, want)
			}
		})
	}
}

// promoteKindStory is a new story under E-0001, S-0001 onwards, that is a
// candidate to be promoted once it is in the backlog: it has a goal, a
// criterion, a forecast, a cost of delay value, and a touch of its own.
func promoteKindStory(t *testing.T, r *Repo) *Item {
	t.Helper()
	s := mustCreate(t, r, Story, "Story", "E-0001")
	s.Body, s.Touches = readyBody, []string{"dir/" + s.ID}
	s.CostOfDelay, s.Forecast = &CostOfDelay{Value: amount(100)}, &Forecast{Duration: "1h"}
	return s
}

// promoteKindSave puts each item in its state and saves it.
func promoteKindSave(t *testing.T, r *Repo, states map[*Item]string) {
	t.Helper()
	for it, state := range states {
		it.Status = state
		if err := r.Save(it); err != nil {
			t.Fatal(err)
		}
	}
}

// S-0739, ADR-0301 §2, §3: the orchestrator promotes a candidate while the
// column the commit enters is under its limit, and sets a backlog story's
// first cost of delay inputs, under every workflow alike, and its refusals
// name the states as the workflow does.
func TestPromoteKindPromotableFollowsTheWorkflow(t *testing.T) {
	for _, f := range promoteKindFlows() {
		t.Run(f.name, func(t *testing.T) {
			r, st := f.open(t), f.state
			r.Manifest.Orchestration.Permissions.PromoteToReady = true
			epic := mustCreate(t, r, Epic, "Epic", "")
			cand, queued, started, closed := promoteKindStory(t, r), promoteKindStory(t, r), promoteKindStory(t, r), promoteKindStory(t, r)
			promoteKindSave(t, r, map[*Item]string{epic: st(InProgress), cand: st(Backlog), queued: st(Ready), started: f.added, closed: st(Done)})
			limit := func(n int) {
				b, err := r.LoadBoard()
				if err != nil {
					t.Fatal(err)
				}
				b.WIPLimits[st(Ready)] = n
				if err := b.Save("2026-11-14"); err != nil {
					t.Fatal(err)
				}
			}
			limit(2)
			if err := r.Promotable(cand); err != nil {
				t.Errorf("a candidate under the limit: %v", err)
			}
			limit(1)
			want := "the " + st(Ready) + " column is at its WIP limit (1 of 1): the orchestrator moves no story to " + st(Ready) + " until one leaves it"
			if err := r.Promotable(cand); err == nil || err.Error() != want {
				t.Errorf("a candidate at the limit: %v, want %q", err, want)
			}
			want = "S-0003 is " + f.added + ", not in the " + st(Backlog) + ": the orchestrator moves a candidate of flai promote --candidates to " + st(Ready) + ", and no other story"
			if err := r.Promotable(started); err == nil || err.Error() != want {
				t.Errorf("a started story: %v, want %q", err, want)
			}

			r.Manifest.Orchestration.Permissions.PlanBacklogStories = true
			if err := r.OrchestratorSetsCostOfDelay(cand, nil); err != nil {
				t.Errorf("a backlog story's first inputs: %v", err)
			}
			want = "S-0003 is " + f.added + ", not in the " + st(Backlog) + ": with plan_backlog_stories"
			if err := r.OrchestratorSetsCostOfDelay(started, nil); err == nil || !strings.HasPrefix(err.Error(), want) {
				t.Errorf("a started story's inputs: %v, want %q...", err, want)
			}
			want = "S-0004 is " + st(Done) + ": the orchestrator sets cost of delay inputs on an open epic or story alone, and none on one " + st(Done) + " or " + st(Cancelled)
			if err := r.OrchestratorSetsCostOfDelay(closed, nil); err == nil || err.Error() != want {
				t.Errorf("a closed story's inputs: %v, want %q", err, want)
			}
		})
	}
}

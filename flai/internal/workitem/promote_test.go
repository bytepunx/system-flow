package workitem

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
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

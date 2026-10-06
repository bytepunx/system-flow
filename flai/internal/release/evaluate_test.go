package release

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func ptr[T any](v T) *T { return &v }

func evalStory(id, status, parent string, value *float64, tags ...string) *workitem.Item {
	it := &workitem.Item{ID: id, Type: workitem.Story, Title: "Story " + id, Status: status, Parent: parent, Tags: tags}
	if value != nil {
		it.CostOfDelay = &workitem.CostOfDelay{Value: value}
	}
	return it
}

func withRelease(rel manifest.Release) manifest.Manifest {
	return manifest.Manifest{Orchestration: manifest.Orchestration{Release: rel}}
}

// thresholdItems are three stories accepted and not yet released, one with
// no value, worth 300 a week together, an epic pending with them, and a
// story released already.
func thresholdItems() ([]*workitem.Item, map[string]bool) {
	items := []*workitem.Item{
		{ID: "E-0001", Type: workitem.Epic, Status: workitem.Done, CostOfDelay: &workitem.CostOfDelay{Value: ptr(1000.0)}},
		evalStory("S-0001", workitem.Done, "E-0001", ptr(200.0)),
		evalStory("S-0002", workitem.Done, "E-0001", ptr(100.0)),
		evalStory("S-0003", workitem.Done, "E-0001", nil),
		evalStory("S-0004", workitem.Done, "E-0001", ptr(5000.0)),
	}
	return items, map[string]bool{"E-0001": true, "S-0001": true, "S-0002": true, "S-0003": true}
}

func TestEvaluateThresholdOnValue(t *testing.T) {
	items, pending := thresholdItems()
	ev := Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Value: ptr(300.0)}), items, pending)
	if !ev.Met || ev.Value != 300 || ev.Count != 3 {
		t.Fatalf("300 a week is at the threshold of 300, from three stories, the epic and the released one left out: %+v", ev)
	}
	if !slices.Equal(ev.Unvalued, []string{"S-0003"}) || len(ev.Pending) != 3 {
		t.Errorf("the story with no value counts and is named: %+v", ev)
	}
	if !strings.Contains(ev.Reason, "at or over the threshold of 300 USD/week") {
		t.Errorf("the reason says at or over: %s", ev.Reason)
	}

	ev = Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Value: ptr(300.01)}), items, pending)
	if ev.Met || !strings.Contains(ev.Reason, "under the threshold of 300.01 USD/week") {
		t.Errorf("300 a week is under 300.01: %+v", ev)
	}
}

func TestEvaluateThresholdOnCount(t *testing.T) {
	items, pending := thresholdItems()
	ev := Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Value: ptr(10000.0), Count: ptr(3)}), items, pending)
	if !ev.Met || !strings.Contains(ev.Reason, "their count, 3, is at or over the threshold of 3") {
		t.Fatalf("three stories reach a count of 3, whatever the value: %+v", ev)
	}
	ev = Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Value: ptr(10000.0), Count: ptr(4)}), items, pending)
	if ev.Met || !strings.Contains(ev.Reason, "their count, 3, is under the threshold of 4") || !strings.Contains(ev.Reason, "under the threshold of 10000 USD/week") {
		t.Errorf("neither figure reached: %+v", ev)
	}
	if *ev.CountThreshold != 4 || *ev.ValueThreshold != 10000 {
		t.Errorf("the thresholds are reported: %+v", ev)
	}
}

func TestEvaluateThresholdWithNothingPending(t *testing.T) {
	items, _ := thresholdItems()
	ev := Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Count: ptr(0)}), items, nil)
	if ev.Met || ev.Count != 0 || !strings.Contains(ev.Reason, "no accepted story is waiting") {
		t.Errorf("nothing pending is not met, even against zero: %+v", ev)
	}
}

func themeItems() []*workitem.Item {
	return []*workitem.Item{
		evalStory("S-0001", workitem.Done, "E-0002", nil, "ui"),
		evalStory("S-0002", workitem.Done, "E-0002", ptr(50.0), "ui"),
		evalStory("S-0003", workitem.Cancelled, "E-0002", nil, "ui"),
		evalStory("S-0004", workitem.InProgress, "E-0003", nil, "ui"),
		evalStory("S-0005", workitem.Done, "E-0003", nil),
	}
}

func TestEvaluateThemeByEpic(t *testing.T) {
	items := themeItems()
	ev := Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Epic: "E-0002"}), items, map[string]bool{"S-0002": true})
	if !ev.Met || len(ev.Theme) != 2 || ev.Epic != "E-0002" {
		t.Fatalf("every story of E-0002 is accepted, the cancelled one aside, and S-0002 is not released: %+v", ev)
	}
	if !ev.Theme[0].Released || ev.Theme[1].Released {
		t.Errorf("S-0001 is released, S-0002 not: %+v", ev.Theme)
	}
	if ev.Count != 1 || ev.Value != 50 {
		t.Errorf("the pending figures are reported too: %+v", ev)
	}

	ev = Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Epic: "E-0002"}), items, nil)
	if ev.Met || !strings.Contains(ev.Reason, "released already") {
		t.Errorf("all accepted and all released is not met: %+v", ev)
	}

	ev = Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Epic: "E-0003"}), items, map[string]bool{"S-0005": true})
	if ev.Met || !slices.Equal(ev.NotAccepted, []string{"S-0004"}) || !strings.Contains(ev.Reason, "1 of the 2 stories of the epic E-0003 is not yet accepted: S-0004") {
		t.Errorf("S-0004 is in progress: %+v", ev)
	}

	ev = Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Epic: "E-0009"}), items, nil)
	if ev.Met || !strings.Contains(ev.Reason, "no story belongs to the epic E-0009") {
		t.Errorf("an epic with no stories is not met: %+v", ev)
	}
}

func TestEvaluateThemeByTag(t *testing.T) {
	items := themeItems()
	ev := Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Tag: "ui"}), items, map[string]bool{"S-0001": true})
	if ev.Met || !slices.Equal(ev.NotAccepted, []string{"S-0004"}) || len(ev.Theme) != 3 {
		t.Fatalf("S-0004 has the tag and is in progress: %+v", ev)
	}
	items[3].Status = workitem.Done
	ev = Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Tag: "ui"}), items, map[string]bool{"S-0001": true})
	if !ev.Met || ev.Tag != "ui" || !strings.Contains(ev.Reason, "every story of the tag ui is accepted, and 1 of the 3 are not yet released") {
		t.Errorf("every ui story accepted, one unreleased: %+v", ev)
	}
	ev = Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Tag: "api"}), items, nil)
	if ev.Met || !strings.Contains(ev.Reason, "no story belongs to the tag api") {
		t.Errorf("a tag on no story is not met: %+v", ev)
	}
}

func TestEvaluateJudgement(t *testing.T) {
	items, pending := thresholdItems()
	for _, m := range []manifest.Manifest{{}, withRelease(manifest.Release{Policy: manifest.ReleaseJudgement})} {
		ev := Evaluate(m, items, pending)
		if ev.Policy != manifest.ReleaseJudgement || ev.Met {
			t.Fatalf("judgement, the default, is never met: %+v", ev)
		}
		if !strings.Contains(ev.Reason, "the orchestrator's or the operator's call") || !strings.Contains(ev.Reason, "3 accepted stories not yet released, worth 300 USD/week") {
			t.Errorf("the reason says whose call it is, with the figures: %s", ev.Reason)
		}
	}
}

func TestEvaluationJSON(t *testing.T) {
	items, pending := thresholdItems()
	m := withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Count: ptr(5)})
	m.Planning.Currency = "EUR"
	data, err := json.Marshal(Evaluate(m, items, pending))
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"policy", "met", "reason", "currency", "value", "count", "count_threshold", "pending", "unvalued"} {
		if _, ok := got[k]; !ok {
			t.Errorf("%s is in the JSON: %s", k, data)
		}
	}
	for _, k := range []string{"value_threshold", "epic", "tag", "theme", "not_accepted", "held_by_epic"} {
		if _, ok := got[k]; ok {
			t.Errorf("%s is left out when unset: %s", k, data)
		}
	}
	if got["currency"] != "EUR" {
		t.Errorf("the project's currency: %s", data)
	}
}

// wholeEpicItems are five accepted stories, all tagged rel: S-0001 and
// S-0005 of E-0001, in progress; S-0002 of E-0002, in review; S-0003 of
// E-0003, done; and S-0004 of no epic.
func wholeEpicItems() []*workitem.Item {
	return []*workitem.Item{
		{ID: "E-0001", Type: workitem.Epic, Status: workitem.InProgress},
		{ID: "E-0002", Type: workitem.Epic, Status: workitem.Review},
		{ID: "E-0003", Type: workitem.Epic, Status: workitem.Done},
		evalStory("S-0001", workitem.Done, "E-0001", ptr(200.0), "rel"),
		evalStory("S-0002", workitem.Done, "E-0002", ptr(100.0), "rel"),
		evalStory("S-0003", workitem.Done, "E-0003", nil, "rel"),
		evalStory("S-0004", workitem.Done, "", nil, "rel"),
		evalStory("S-0005", workitem.Done, "E-0001", nil, "rel"),
	}
}

// S-0222: whole_epics holds a batch back, under every policy, while a story
// in it belongs to an epic in neither review nor done.
func TestEvaluateWholeEpics(t *testing.T) {
	policies := []struct {
		rel manifest.Release
		met bool // what the policy says with nothing held
	}{
		{manifest.Release{Policy: manifest.ReleaseThreshold, Count: ptr(1)}, true},
		{manifest.Release{Policy: manifest.ReleaseTheme, Tag: "rel"}, true},
		{manifest.Release{Policy: manifest.ReleaseJudgement}, false},
	}
	cases := []struct {
		name    string
		whole   bool
		pending []string
		held    []string
	}{
		{"a story whose epic is in progress", true, []string{"S-0001", "S-0002", "S-0004"}, []string{"S-0001"}},
		{"epics all in review or done", true, []string{"S-0002", "S-0003"}, nil},
		{"a story with no epic", true, []string{"S-0004"}, nil},
		{"whole_epics off", false, []string{"S-0001", "S-0002", "S-0004"}, nil},
	}
	for _, p := range policies {
		for _, c := range cases {
			t.Run(p.rel.Policy+"/"+c.name, func(t *testing.T) {
				pending := map[string]bool{}
				for _, id := range c.pending {
					pending[id] = true
				}
				rel := p.rel
				rel.WholeEpics = c.whole
				ev := Evaluate(withRelease(rel), wholeEpicItems(), pending)
				var held []string
				for _, h := range ev.HeldByEpic {
					held = append(held, h.ID)
				}
				if !slices.Equal(held, c.held) {
					t.Fatalf("held %v, want %v: %+v", held, c.held, ev)
				}
				if ev.Count != len(c.pending) {
					t.Errorf("the held stories stay pending: %+v", ev)
				}
				if len(c.held) > 0 {
					if ev.Met {
						t.Errorf("a held batch is not met: %+v", ev)
					}
					if h := ev.HeldByEpic[0]; h.Epic != "E-0001" || h.EpicStatus != workitem.InProgress || h.Title != "Story S-0001" {
						t.Errorf("the held story names its epic and the epic's status: %+v", h)
					}
					if !strings.Contains(ev.Reason, "; whole_epics holds the batch back until its epic is in review or done: E-0001 is in-progress, with S-0001") {
						t.Errorf("the reason names the epic holding the batch back: %s", ev.Reason)
					}
					return
				}
				if ev.Met != p.met || strings.Contains(ev.Reason, "whole_epics") {
					t.Errorf("nothing held leaves the policy's verdict, %v, alone: %+v", p.met, ev)
				}
			})
		}
	}
}

func TestEvaluateWholeEpicsNamesEachEpic(t *testing.T) {
	items := append(wholeEpicItems(),
		&workitem.Item{ID: "E-0004", Type: workitem.Epic, Status: workitem.Backlog},
		evalStory("S-0006", workitem.Done, "E-0004", nil),
		evalStory("S-0007", workitem.Done, "E-0009", nil),
	)
	pending := map[string]bool{"S-0001": true, "S-0005": true, "S-0006": true, "S-0007": true, "S-0003": true}
	ev := Evaluate(withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Count: ptr(1), WholeEpics: true}), items, pending)
	if ev.Met || len(ev.HeldByEpic) != 4 {
		t.Fatalf("S-0001, S-0005, S-0006, and S-0007 are held: %+v", ev)
	}
	if !strings.Contains(ev.Reason, "until their epics are in review or done: E-0001 is in-progress, with S-0001, S-0005; E-0004 is backlog, with S-0006; E-0009 is not found, with S-0007") {
		t.Errorf("each epic is named once, with its stories; one not found holds too: %s", ev.Reason)
	}
	if h := ev.HeldByEpic[3]; h.Epic != "E-0009" || h.EpicStatus != "" {
		t.Errorf("an epic not found has no status: %+v", h)
	}

	data, err := json.Marshal(ev)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Held []map[string]string `json:"held_by_epic"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Held) != 4 || got.Held[0]["id"] != "S-0001" || got.Held[0]["epic"] != "E-0001" || got.Held[0]["epic_status"] != "in-progress" {
		t.Errorf("held_by_epic is in the JSON, each with its epic and status: %s", data)
	}
}

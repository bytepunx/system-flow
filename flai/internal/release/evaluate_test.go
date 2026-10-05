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
	for _, k := range []string{"value_threshold", "epic", "tag", "theme", "not_accepted"} {
		if _, ok := got[k]; ok {
			t.Errorf("%s is left out when unset: %s", k, data)
		}
	}
	if got["currency"] != "EUR" {
		t.Errorf("the project's currency: %s", data)
	}
}

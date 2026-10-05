package workitem

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// rankable is a ready story with a cost of delay value and a forecast duration
// when given, created at a time.
func rankable(id string, value *float64, duration, created string) *Item {
	it := &Item{ID: id, Type: Story, Status: Ready, Title: "Story " + id, Created: created}
	if value != nil {
		it.CostOfDelay = &CostOfDelay{Value: value}
	}
	if duration != "" {
		it.Forecast = &Forecast{Duration: duration}
	}
	return it
}

func amount(v float64) *float64 { return &v }

// In board order: S-0001 lacks a value, S-0002 and S-0004 tie on cod, S-0005
// lacks a forecast, S-0003 lacks both, and S-0006 ties S-0002 on created.
func policyFixture() []*Item {
	return []*Item{
		rankable("S-0001", nil, "2h", "2026-09-03T10:00:00Z"),
		rankable("S-0002", amount(400), "4h", "2026-09-02T10:00:00Z"),
		rankable("S-0003", nil, "", "not a time"),
		rankable("S-0004", amount(400), "1h", "2026-09-04T10:00:00Z"),
		rankable("S-0005", amount(900), "", "2026-09-01T10:00:00Z"),
		rankable("S-0006", amount(100), "30m", "2026-09-02T10:00:00Z"),
	}
}

func TestOrderByPolicy(t *testing.T) {
	cases := []struct {
		policy  string
		ids     []string
		texts   []string
		missing map[string]string
	}{
		{PolicyCOD,
			[]string{"S-0005", "S-0002", "S-0004", "S-0006", "S-0001", "S-0003"},
			[]string{"900", "400", "400", "100", "", ""},
			map[string]string{"S-0001": "cost of delay value", "S-0003": "cost of delay value"}},
		{PolicyWSJF,
			[]string{"S-0004", "S-0006", "S-0002", "S-0001", "S-0003", "S-0005"},
			[]string{"400/h", "200/h", "100/h", "", "", ""},
			map[string]string{"S-0001": "cost of delay value", "S-0003": "cost of delay value and forecast duration", "S-0005": "forecast duration"}},
		{PolicyThroughput,
			[]string{"S-0006", "S-0004", "S-0001", "S-0002", "S-0003", "S-0005"},
			[]string{"30m", "1h", "2h", "4h", "", ""},
			map[string]string{"S-0003": "forecast duration", "S-0005": "forecast duration"}},
		{PolicyFIFO,
			[]string{"S-0005", "S-0002", "S-0006", "S-0001", "S-0004", "S-0003"},
			[]string{"2026-09-01T10:00:00Z", "2026-09-02T10:00:00Z", "2026-09-02T10:00:00Z", "2026-09-03T10:00:00Z", "2026-09-04T10:00:00Z", ""},
			map[string]string{"S-0003": "created"}},
	}
	for _, c := range cases {
		got, err := OrderByPolicy(policyFixture(), c.policy)
		if err != nil {
			t.Fatalf("%s: %v", c.policy, err)
		}
		var ids, texts []string
		for i, r := range got {
			ids = append(ids, r.ID)
			texts = append(texts, r.Text)
			if r.Position != i+1 {
				t.Errorf("%s: %s is at position %d, want %d", c.policy, r.ID, r.Position, i+1)
			}
			if r.Missing != c.missing[r.ID] {
				t.Errorf("%s: %s missing %q, want %q", c.policy, r.ID, r.Missing, c.missing[r.ID])
			}
			if (r.Figure == nil) != (r.Missing != "") {
				t.Errorf("%s: %s has figure %v and missing %q", c.policy, r.ID, r.Figure, r.Missing)
			}
		}
		if !reflect.DeepEqual(ids, c.ids) {
			t.Errorf("%s: got %v want %v", c.policy, ids, c.ids)
		}
		if !reflect.DeepEqual(texts, c.texts) {
			t.Errorf("%s figures: got %q want %q", c.policy, texts, c.texts)
		}
	}
}

func TestOrderByPolicyFiguresAndJSON(t *testing.T) {
	got, err := OrderByPolicy([]*Item{rankable("S-0001", amount(300), "1h30m", "2026-09-01T10:00:00Z")}, PolicyThroughput)
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Figure == nil || *got[0].Figure != 1.5 || got[0].Unit != "hours" {
		t.Errorf("throughput figure is the duration in hours: %+v", got[0])
	}
	data, _ := json.Marshal(got[0])
	if want := `{"position":1,"id":"S-0001","title":"Story S-0001","figure":1.5,"unit":"hours","text":"1h30m"}`; string(data) != want {
		t.Errorf("json:\n got %s\nwant %s", data, want)
	}
	got, _ = OrderByPolicy([]*Item{rankable("S-0002", nil, "", "")}, PolicyCOD)
	data, _ = json.Marshal(got[0])
	if want := `{"position":1,"id":"S-0002","title":"Story S-0002","unit":"per week","missing":"cost of delay value"}`; string(data) != want {
		t.Errorf("json without the figure:\n got %s\nwant %s", data, want)
	}
	got, _ = OrderByPolicy([]*Item{rankable("S-0003", amount(100), "3h", "")}, PolicyWSJF)
	if got[0].Text != "33.33/h" {
		t.Errorf("wsjf text rounds to two places: %q", got[0].Text)
	}
}

func TestOrderByPolicyRefusesAnUnknownPolicy(t *testing.T) {
	_, err := OrderByPolicy(policyFixture(), "random")
	if err == nil || !strings.Contains(err.Error(), "cod, wsjf, throughput, fifo") {
		t.Errorf("an unknown policy names the valid ones: %v", err)
	}
	if got, err := OrderByPolicy(nil, PolicyCOD); err != nil || len(got) != 0 {
		t.Errorf("no stories: %v %v", got, err)
	}
}

func TestApplyReadyOrder(t *testing.T) {
	b, items := orderFixture()
	if err := b.ApplyReadyOrder(items, []string{"S-0001", "S-0002"}); err != nil {
		t.Fatal(err)
	}
	if want := []string{"S-0001", "S-0002", "S-0005"}; !reflect.DeepEqual(b.Order, want) {
		t.Errorf("ready reordered, backlog names kept: got %v want %v", b.Order, want)
	}
	if err := b.ApplyReadyOrder(items, []string{"S-0007"}); err == nil {
		t.Error("an in-progress story is refused")
	}
}

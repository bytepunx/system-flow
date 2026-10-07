package metrics

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// forecastItems are stories started on 2026-08-20 at midnight, with and
// without a forecast, an estimate, a cycle time, and a model, and those the
// aggregates leave out: open, cancelled, and done before the window.
func forecastItems() []*workitem.Item {
	story := func(id, nature, model, done string) *workitem.Item {
		it := doneWith(id, done, nil)
		it.Nature = nature
		if model != "" {
			it.Agent = &manifest.Agent{Harness: "claude-code", Model: model}
		}
		return it
	}
	a := story("S-0001", "feature", "claude-opus-5-5", "2026-08-21T00:00:00Z")
	a.Forecast = &workitem.Forecast{Duration: "20h", Delivery: "2026-08-21T02:00:00Z"}
	a.Estimate = "30h"
	b := story("S-0002", "remediation", "claude-haiku-4-5", "2026-08-20T06:00:00Z")
	b.Forecast = &workitem.Forecast{Duration: "4h", Delivery: "2026-08-20T05:00:00Z"}
	c := story("S-0003", "feature", "", "2026-08-22T00:00:00Z")
	c.Forecast = &workitem.Forecast{Duration: "51h"}
	c.Estimate = "40h"
	d := story("S-0004", "improvement", "", "2026-08-20T01:00:00Z")
	d.Agent = &manifest.Agent{Harness: "claude-code"} // no model
	d.Estimate = "2h"
	e := story("S-0005", "research", "claude-opus-5-5", "2026-08-25T00:00:00Z")
	e.Transitions = e.Transitions[1:] // never started: no cycle time
	e.Forecast = &workitem.Forecast{Duration: "1h", Delivery: "2026-08-24T00:00:00Z"}
	f := story("S-0006", "feature", "claude-opus-5-5", "")
	f.Status = workitem.InProgress
	f.Transitions = f.Transitions[:1]
	f.Forecast = &workitem.Forecast{Duration: "2h", Delivery: "2026-08-31T00:00:00Z"}
	f.Estimate = "3h"
	g := story("S-0007", "feature", "claude-opus-5-5", "2026-08-25T00:00:00Z")
	g.Status = workitem.Cancelled
	g.Transitions[1].To = workitem.Cancelled
	g.Forecast = &workitem.Forecast{Duration: "1h", Delivery: "2026-08-21T00:00:00Z"}
	g.Estimate = "1h"
	h := story("S-0008", "feature", "claude-opus-5-5", "2026-07-01T00:00:00Z")
	h.Transitions[0].At = "2026-06-30T00:00:00Z"
	h.Forecast = &workitem.Forecast{Duration: "1h", Delivery: "2026-06-30T00:00:00Z"}
	h.Estimate = "1h"
	i := story("S-0009", "feature", "claude-opus-5-5", "2026-08-23T00:00:00Z")
	return []*workitem.Item{a, b, c, d, e, f, g, h, i}
}

func val(v float64) *float64 { return &v }

// S-0205: each item's forecast and the errors of its forecast, delivery, and
// estimate in seconds, a positive one later or longer, absent without input.
func TestForecastAndEstimateErrorsPerItem(t *testing.T) {
	rep := Compute(forecastItems(), Options{Now: now})
	want := map[string][4]*float64{ // forecast, forecast error, delivery error, estimate error
		"S-0001": {val(72000), val(14400), val(-7200), val(-21600)},
		"S-0002": {val(14400), val(7200), val(3600), nil},
		"S-0003": {val(183600), val(-10800), nil, val(28800)},
		"S-0004": {nil, nil, nil, val(-3600)},
		"S-0005": {val(3600), nil, val(86400), nil},
		"S-0006": {val(7200), nil, nil, nil},
		"S-0007": {val(3600), val(428400), val(345600), val(428400)},
		"S-0008": {val(3600), val(82800), val(86400), val(82800)},
		"S-0009": {nil, nil, nil, nil},
	}
	names := [4]string{"forecast_seconds", "forecast_error_seconds", "delivery_error_seconds", "estimate_error_seconds"}
	for _, m := range rep.Items {
		w, ok := want[m.ID]
		if !ok {
			t.Fatalf("unexpected item %s", m.ID)
		}
		for k, got := range [4]*float64{m.Forecast, m.ForecastError, m.DeliveryError, m.EstErrorSeconds} {
			if (got == nil) != (w[k] == nil) || (got != nil && *got != *w[k]) {
				t.Errorf("%s %s = %v, want %v", m.ID, names[k], deref(got), deref(w[k]))
			}
		}
	}
	a := rep.Items[0]
	if a.EstError == nil || *a.EstError != -0.2 {
		t.Errorf("S-0001 estimate_error = %v, want -0.2 beside the seconds", deref(a.EstError))
	}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	keys := `"estimate_error_seconds":-21600,"forecast_seconds":72000,"forecast_error_seconds":14400,"delivery_error_seconds":-7200`
	if !strings.Contains(string(data), keys) {
		t.Errorf("S-0001 json lacks %s: %s", keys, data)
	}
	data, err = json.Marshal(rep.Items[8])
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range names {
		if strings.Contains(string(data), k) {
			t.Errorf("S-0009 has no forecast or estimate and carries %s: %s", k, data)
		}
	}
}

// ADR-0111: each item carries the model its errors are grouped under, its
// agent's model, (none) without an agent or a model.
func TestItemModelIsTheOneForecastsGroupBy(t *testing.T) {
	rep := Compute(forecastItems(), Options{Now: now})
	want := map[string]string{"S-0001": "claude-opus-5-5", "S-0002": "claude-haiku-4-5", "S-0003": "(none)", "S-0004": "(none)"}
	for _, m := range rep.Items {
		if w, ok := want[m.ID]; ok && m.Model != w {
			t.Errorf("%s model = %q, want %q", m.ID, m.Model, w)
		}
	}
	data, err := json.Marshal(rep.Items[2])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"model":"(none)"`) {
		t.Errorf("S-0003 json lacks its model: %s", data)
	}
}

// S-0205: the absolute errors of the stories done in the window, in all, by
// nature, and by model, with (none) for an agent without a model.
func TestForecastsSpreadTheAbsoluteErrorsOfItemsDoneInTheWindow(t *testing.T) {
	data, err := json.Marshal(Compute(forecastItems(), Options{Now: now}).Forecasts)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"forecast":{"count":3,"p50_seconds":10800,"p85_seconds":14400,` +
		`"by_nature":{"feature":{"count":2,"p50_seconds":10800,"p85_seconds":14400},"remediation":{"count":1,"p50_seconds":7200,"p85_seconds":7200}},` +
		`"by_model":{"(none)":{"count":1,"p50_seconds":10800,"p85_seconds":10800},"claude-haiku-4-5":{"count":1,"p50_seconds":7200,"p85_seconds":7200},"claude-opus-5-5":{"count":1,"p50_seconds":14400,"p85_seconds":14400}}},` +
		`"delivery":{"count":3,"p50_seconds":7200,"p85_seconds":86400,` +
		`"by_nature":{"feature":{"count":1,"p50_seconds":7200,"p85_seconds":7200},"remediation":{"count":1,"p50_seconds":3600,"p85_seconds":3600},"research":{"count":1,"p50_seconds":86400,"p85_seconds":86400}},` +
		`"by_model":{"claude-haiku-4-5":{"count":1,"p50_seconds":3600,"p85_seconds":3600},"claude-opus-5-5":{"count":2,"p50_seconds":7200,"p85_seconds":86400}}},` +
		`"estimate":{"count":3,"p50_seconds":21600,"p85_seconds":28800,` +
		`"by_nature":{"feature":{"count":2,"p50_seconds":21600,"p85_seconds":28800},"improvement":{"count":1,"p50_seconds":3600,"p85_seconds":3600}},` +
		`"by_model":{"(none)":{"count":2,"p50_seconds":3600,"p85_seconds":28800},"claude-opus-5-5":{"count":1,"p50_seconds":21600,"p85_seconds":21600}}}}`
	if string(data) != want {
		t.Errorf("forecasts =\n%s\nwant\n%s", data, want)
	}
}

func TestForecastsOfNoItemsCountZeroWithoutPercentiles(t *testing.T) {
	data, err := json.Marshal(Compute(nil, Options{Now: now}).Forecasts)
	if err != nil {
		t.Fatal(err)
	}
	empty := `{"count":0,"by_nature":{},"by_model":{}}`
	want := `{"forecast":` + empty + `,"delivery":` + empty + `,"estimate":` + empty + `}`
	if string(data) != want {
		t.Errorf("forecasts = %s, want %s", data, want)
	}
}

func deref(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

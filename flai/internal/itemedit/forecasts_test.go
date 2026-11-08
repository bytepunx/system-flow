package itemedit

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// forecastRepo is planningRepo with two more stories under the epic, each of
// the three with a forecast, as the replanner finds them.
func forecastRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	repo := planningRepo(t)
	created := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	for _, title := range []string{"Story two", "Story three"} {
		if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: "E-0001", Owner: "alex", Now: created}); err != nil {
			t.Fatal(err)
		}
	}
	at := time.Date(2026, 10, 2, 9, 0, 0, 0, time.UTC)
	for _, id := range []string{"S-0001", "S-0002", "S-0003"} {
		if _, err := Apply(repo, nil, id, Change{Forecast: &ForecastEdit{Duration: str("2h"), Delivery: str("2026-10-05T09:00:00Z"), Basis: str("first")}}, Options{By: "planner", NoCommit: true, Now: at}); err != nil {
			t.Fatalf("forecast %s: %v", id, err)
		}
	}
	return repo
}

func moves() []ForecastChange {
	return []ForecastChange{
		{ID: "S-0001", Edit: ForecastEdit{Delivery: str("2026-10-06T10:00:00Z"), Basis: str("replayed one")}},
		{ID: "S-0002", Edit: ForecastEdit{Delivery: str("2026-10-07T11:00:00Z"), Basis: str("replayed two")}},
		{ID: "S-0003", Edit: ForecastEdit{Delivery: str("2026-10-08T12:00:00Z"), Basis: str("replayed three")}},
	}
}

func itemFile(t *testing.T, repo *workitem.Repo, id string) []byte {
	t.Helper()
	it, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(it.Path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// S-0599: the batch writes each item byte for byte as Apply writes the same
// change, stamps and body included.
func TestApplyForecastsWritesWhatApplyWrites(t *testing.T) {
	one, batch := forecastRepo(t), forecastRepo(t)
	at := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	opt := Options{By: "flai", NoCommit: true, Now: at}
	for _, m := range moves() {
		e := m.Edit
		if _, err := Apply(one, nil, m.ID, Change{Forecast: &e}, opt); err != nil {
			t.Fatalf("Apply %s: %v", m.ID, err)
		}
	}
	res, err := ApplyForecasts(batch, nil, moves(), opt)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Refused) != 0 || !slices.Equal(res.Written, []string{"S-0001", "S-0002", "S-0003"}) {
		t.Fatalf("written %v, refused %v", res.Written, res.Refused)
	}
	for _, m := range moves() {
		if got, want := itemFile(t, batch, m.ID), itemFile(t, one, m.ID); string(got) != string(want) {
			t.Errorf("%s: the batch wrote\n%s\nApply wrote\n%s", m.ID, got, want)
		}
	}
	for _, f := range []string{"wip/agents/index.md", "wip/kanban/stories/S-0001-story.md"} {
		if !slices.Contains(res.Files, f) {
			t.Errorf("files %v lack %s", res.Files, f)
		}
	}
	var noticed []string
	for _, n := range Notices(batch) {
		if n.By == "flai" {
			noticed = append(noticed, n.ID)
		}
	}
	if !slices.Equal(noticed, []string{"S-0001", "S-0002", "S-0003"}) {
		t.Errorf("notices by flai %v, want one per item written", noticed)
	}
}

// S-0599: however many forecasts move, the index is written once and the
// check runs once, over the items written; with nothing to change, neither.
func TestApplyForecastsWritesTheIndexAndChecksOnce(t *testing.T) {
	repo := forecastRepo(t)
	var checks, indexes int
	var checked []string
	runItems = func(repo *workitem.Repo, all, items []*workitem.Item, now time.Time) *check.Result {
		checks++
		for _, it := range items {
			checked = append(checked, it.ID)
		}
		return check.Items(repo, all, items, now)
	}
	writeIndex = func(r *workitem.Repo, items []*workitem.Item, now time.Time) error {
		indexes++
		return r.WriteIndex(items, now)
	}
	t.Cleanup(func() { runItems, writeIndex = check.Items, (*workitem.Repo).WriteIndex })
	at := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	if _, err := ApplyForecasts(repo, nil, moves(), Options{By: "flai", Now: at}); err != nil {
		t.Fatal(err)
	}
	if checks != 1 || indexes != 1 || !slices.Equal(checked, []string{"S-0001", "S-0002", "S-0003"}) {
		t.Errorf("checks %d over %v, index writes %d; want one of each over the three", checks, checked, indexes)
	}
	// the same moves again change nothing: no write, no index, no check
	res, err := ApplyForecasts(repo, nil, moves(), Options{By: "flai", Now: at.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Written) != 0 || len(res.Files) != 0 || checks != 1 || indexes != 1 {
		t.Errorf("unchanged forecasts: written %v, files %v, checks %d, index writes %d", res.Written, res.Files, checks, indexes)
	}
}

// S-0599, ADR-0184: a forecast another writer changed between the batch's
// read and its save is refused for that item alone and left as that writer
// saved it; another writer's change to another part of an item is merged
// with the batch's forecast; the rest are written.
func TestApplyForecastsMergesOrRefusesAnotherWritersSave(t *testing.T) {
	repo := forecastRepo(t)
	other := time.Date(2026, 10, 3, 8, 0, 0, 0, time.UTC)
	beforeSave = func() {
		s1, err := repo.Get("S-0001")
		if err != nil {
			t.Fatal(err)
		}
		s1.Forecast.Delivery, s1.Forecast.Basis = "2026-10-09T09:00:00Z", "the planner's"
		s1.Updated = other.Format(workitem.TimeFormat)
		if err := repo.Save(s1); err != nil {
			t.Fatal(err)
		}
		s2, err := repo.Get("S-0002")
		if err != nil {
			t.Fatal(err)
		}
		s2.Tags = []string{"cli"}
		s2.Updated = other.Format(workitem.TimeFormat)
		if err := repo.Save(s2); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { beforeSave = func() {} })
	at := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	res, err := ApplyForecasts(repo, nil, moves(), Options{By: "flai", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	var conflict *workitem.ConflictError
	if !errors.As(res.Refused["S-0001"], &conflict) || !slices.Contains(conflict.Keys, "forecast") {
		t.Errorf("S-0001 refused with %v, want a conflict on forecast", res.Refused["S-0001"])
	}
	if !slices.Equal(res.Written, []string{"S-0002", "S-0003"}) {
		t.Errorf("written %v, want S-0002 and S-0003", res.Written)
	}
	if slices.Contains(res.Files, "wip/kanban/stories/S-0001-story.md") {
		t.Errorf("files %v name the refused S-0001", res.Files)
	}
	s1, err := workitem.ReadItem(filepath.Join(repo.Root, "wip/kanban/stories/S-0001-story.md"))
	if err != nil {
		t.Fatal(err)
	}
	if s1.Forecast.Delivery != "2026-10-09T09:00:00Z" || s1.Forecast.Basis != "the planner's" {
		t.Errorf("S-0001's forecast %+v, want the other writer's", s1.Forecast)
	}
	s2, err := repo.Get("S-0002")
	if err != nil {
		t.Fatal(err)
	}
	if s2.Forecast.Delivery != "2026-10-07T11:00:00Z" || !slices.Equal(s2.Tags, []string{"cli"}) {
		t.Errorf("S-0002: forecast %+v and tags %v, want the batch's forecast and the other writer's tags", s2.Forecast, s2.Tags)
	}
}

// A change to an archived or missing item is refused alone.
func TestApplyForecastsRefusesWhatItCannotEdit(t *testing.T) {
	repo := forecastRepo(t)
	at := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	res, err := ApplyForecasts(repo, nil, append(moves()[:1], ForecastChange{ID: "S-0099", Edit: ForecastEdit{Delivery: str("2026-10-06T10:00:00Z")}}), Options{By: "flai", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	if !errors.Is(res.Refused["S-0099"], workitem.ErrNotFound) || !slices.Equal(res.Written, []string{"S-0001"}) {
		t.Errorf("written %v, refused %v", res.Written, res.Refused)
	}
}

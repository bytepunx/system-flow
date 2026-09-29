package hostapi

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/perf"
)

// phases runs one method with a recorder and returns the phases it marked.
func phases(t *testing.T, m func(context.Context) error) map[string]int {
	t.Helper()
	ctx, rec := perf.Start(context.Background())
	if err := m(ctx); err != nil {
		t.Fatal(err)
	}
	out := map[string]int{}
	for _, p := range rec.Phases() {
		out[p.Name] = p.Count
	}
	return out
}

func TestReadMethodsMarkWhereTheTimeGoes(t *testing.T) {
	p := harbour(t)
	table := Methods("test", func() time.Time { return t0.Add(time.Hour) })
	for method, want := range map[string][]string{
		"board.get":      {"repo.open", "repo.list", "board.load", "release.pending", "board.view"},
		"item.get":       {"repo.open", "repo.get", "repo.list"},
		"items.list":     {"repo.open", "repo.list"},
		"threads.list":   {"repo.open", "threads.read", "threads.view"},
		"inbox.designer": {"repo.open", "threads.read", "repo.list", "narratives.read", "check.run"},
		"activity.get":   {"repo.open", "repo.list", "narratives.read"},
		"docs.tree":      {"docs.walk"},
		"project.info":   {"manifest.load"},
	} {
		params := `{}`
		if method == "item.get" {
			params = `{"id":"S-0001"}`
		}
		got := phases(t, func(ctx context.Context) error {
			if _, rerr := table[method](ctx, p, json.RawMessage(params)); rerr != nil {
				return rerr
			}
			return nil
		})
		for _, name := range want {
			if got[name] != 1 {
				t.Errorf("%s: phase %s counted %d times, want once; phases %v", method, name, got[name], got)
			}
		}
	}
}

func TestAFlaiItRunsIsAPhaseNamedForTheCommand(t *testing.T) {
	p := harbour(t)
	rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
	m := writeMethods(rec.run, time.Now, Host{})["stats.get"]
	got := phases(t, func(ctx context.Context) error {
		if _, rerr := m(ctx, p, json.RawMessage(`{}`)); rerr != nil {
			return rerr
		}
		return nil
	})
	if got["exec.flai.stats"] != 1 {
		t.Errorf("phases %v", got)
	}
}

package execx

import (
	"context"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/perf"
)

type fake struct{ ran []string }

func (f *fake) Run(dir, name string, args ...string) (string, error) {
	return f.RunInput(dir, name, "", args...)
}

func (f *fake) RunInput(_, name, _ string, _ ...string) (string, error) {
	f.ran = append(f.ran, name)
	return "", nil
}

func (f *fake) LookPath(name string) (string, error) { return name, nil }

func TestTimedRecordsEachCommandAsAPhase(t *testing.T) {
	f := &fake{}
	if Timed(context.Background(), f) != Runner(f) {
		t.Error("with nothing timing the context the runner is itself")
	}
	ctx, rec := perf.Start(context.Background())
	r := Timed(ctx, f)
	_, _ = r.Run("/repo", "git", "log", "--oneline")
	_, _ = r.Run("/repo", "git", "log", "-1")
	_, _ = r.RunInput("/repo", "git", "input", "tag", "--list")
	ph := map[string]int{}
	for _, p := range rec.Phases() {
		ph[p.Name] = p.Count
	}
	if len(f.ran) != 3 || ph["exec.git.log"] != 2 || ph["exec.git.tag"] != 1 {
		t.Errorf("ran %v, phases %v", f.ran, ph)
	}
}

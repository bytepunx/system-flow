package hostapi

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0804: resolve.status answers the story's last close-out or rebase and
// resolve job from .flai-cache/verify/jobs, whatever the host actions say,
// in flai serve's process: the story alone while it has had none; the job's
// fields, running, and the step, tier, and steps its log's progress lines
// say while it runs, with where the log is; and how it ended once it ended,
// with no tier. It runs nothing and records nothing, and is asked of a
// story only.
func TestResolveStatusAnswersTheStorysLastJobAndItsProgress(t *testing.T) {
	p := withDocs(t)
	var journal []Entry
	host := Host{Enabled: func(string, string) bool { return false }, Record: func(e Entry) { journal = append(journal, e) }}
	rec := &recorder{}
	call := func(params string) (Written, *channel.Error) {
		t.Helper()
		res, e := writeMethods(rec.run, time.Now, host)["resolve.status"](context.Background(), p, json.RawMessage(params))
		w, _ := res.(Written)
		return w, e
	}
	type answered struct {
		ID  string `json:"id"`
		Job *struct {
			ID        string                  `json:"id"`
			Kind      string                  `json:"kind"`
			By        string                  `json:"by"`
			Started   time.Time               `json:"started"`
			Ended     time.Time               `json:"ended"`
			Verdict   string                  `json:"verdict"`
			StoppedAt string                  `json:"stopped_at"`
			Outcome   string                  `json:"outcome"`
			Thread    string                  `json:"thread"`
			Conflicts []string                `json:"conflicts"`
			Reason    string                  `json:"reason"`
			Running   bool                    `json:"running"`
			Step      string                  `json:"step"`
			Tier      string                  `json:"tier"`
			Steps     []closeout.StepProgress `json:"steps"`
			Log       string                  `json:"log"`
		} `json:"job"`
	}
	decoded := func(w Written) answered {
		t.Helper()
		var a answered
		if err := json.Unmarshal(w.Data, &a); err != nil {
			t.Fatalf("answer %s: %v", w.Data, err)
		}
		return a
	}

	if w, e := call(`{"id":"S-0001"}`); e != nil || string(w.Data) != `{"id":"S-0001"}` {
		t.Errorf("no job: %s %+v", w.Data, e)
	}
	if _, e := call(`{"id":"T-0001"}`); e == nil || !strings.Contains(e.Message, "not a story") {
		t.Errorf("a task must be refused: %+v", e)
	}

	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	cache := repo.CacheDir()
	dir := closeout.Dir(cache)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	started := time.Date(2026, 11, 16, 9, 0, 0, 0, time.UTC)
	at := func(s int) time.Time { return started.Add(time.Duration(s) * time.Second) }
	job := closeout.Job{ID: "S-0001-20261116T090000Z", Story: "S-0001", Kind: closeout.KindResolve, By: "olive", Started: started}
	if err := closeout.WriteJob(dir, job); err != nil {
		t.Fatal(err)
	}
	log := strings.Join([]string{
		closeout.StartedLine(closeout.SourceResolve, "S-0001", closeout.LineStep, closeout.StepRebase, at(0)),
		"Successfully rebased and updated refs/heads/story/S-0001.",
		closeout.EndedLine(closeout.SourceResolve, "S-0001", closeout.LineStep, "passed", closeout.StepRebase, at(2)),
		closeout.StartedLine(closeout.SourceResolve, "S-0001", closeout.LineStep, closeout.StepCloseOut, at(3)),
		"close-out: flai verify S-0001 --json",
		closeout.StartedLine(closeout.SourceVerify, "S-0001", closeout.LineTier, "unit", at(5)),
	}, "\n") + "\n"
	if err := os.WriteFile(closeout.LogPath(cache, "S-0001"), []byte(log), 0o644); err != nil {
		t.Fatal(err)
	}

	w, e := call(`{"id":"S-1"}`)
	if e != nil {
		t.Fatalf("a running job: %+v", e)
	}
	a := decoded(w)
	j := a.Job
	if a.ID != "S-0001" || j == nil || j.ID != job.ID || j.Kind != closeout.KindResolve || j.By != "olive" || !j.Started.Equal(started) || !j.Ended.IsZero() || !j.Running {
		t.Fatalf("the running job: %s", w.Data)
	}
	if j.Step != closeout.StepVerify || j.Tier != "unit" || j.Log != closeout.LogPath(cache, "S-0001") {
		t.Errorf("step %q, tier %q, log %q: %s", j.Step, j.Tier, j.Log, w.Data)
	}
	want := []closeout.StepProgress{
		{Name: closeout.StepRebase, Kind: closeout.LineStep, Source: closeout.SourceResolve, Started: at(0), Ended: at(2), Outcome: "passed"},
		{Name: closeout.StepCloseOut, Kind: closeout.LineStep, Source: closeout.SourceResolve, Started: at(3)},
		{Name: "unit", Kind: closeout.LineTier, Source: closeout.SourceVerify, Started: at(5)},
	}
	if len(j.Steps) != len(want) {
		t.Fatalf("steps %+v, want %+v", j.Steps, want)
	}
	for i := range want {
		if g := j.Steps[i]; g.Name != want[i].Name || g.Kind != want[i].Kind || g.Source != want[i].Source || !g.Started.Equal(want[i].Started) || !g.Ended.Equal(want[i].Ended) || g.Outcome != want[i].Outcome {
			t.Errorf("step %d: %+v, want %+v", i, g, want[i])
		}
	}

	// the job ended: the close-out stopped at unit and the story was sent back
	f, err := os.OpenFile(closeout.LogPath(cache, "S-0001"), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(closeout.EndedLine(closeout.SourceVerify, "S-0001", closeout.LineTier, "failed", "unit", at(9)) + "\n" +
		closeout.EndedLine(closeout.SourceResolve, "S-0001", closeout.LineStep, "failed", closeout.StepCloseOut, at(10)) + "\n")
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		t.Fatal(err)
	}
	job.Ended, job.Verdict, job.StoppedAt = at(10), closeout.VerdictStopped, "unit"
	job.Outcome, job.Thread, job.Reason = closeout.OutcomeSentBack, "TH-0007", "rebase and resolve job S-0001-20261116T090000Z: story/S-0001's close-out stopped at unit"
	if err := closeout.WriteJob(dir, job); err != nil {
		t.Fatal(err)
	}
	w, e = call(`{"id":"S-0001"}`)
	if e != nil {
		t.Fatalf("an ended job: %+v", e)
	}
	j = decoded(w).Job
	if j == nil || j.Running || !j.Ended.Equal(at(10)) || j.Verdict != closeout.VerdictStopped || j.StoppedAt != "unit" || j.Outcome != closeout.OutcomeSentBack || j.Thread != "TH-0007" || j.Reason != job.Reason || j.Tier != "" || len(j.Steps) != 3 || j.Steps[2].Outcome != "failed" {
		t.Errorf("the ended job: %s", w.Data)
	}

	if len(rec.runs) != 0 || len(journal) != 0 {
		t.Errorf("a read ran %d commands and journalled %+v", len(rec.runs), journal)
	}
}

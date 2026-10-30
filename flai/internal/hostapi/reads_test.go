package hostapi

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// readPhase is, per read a command also prints, a valid call and the phase
// its work is timed under (S-0159).
var readPhase = map[string]struct{ params, phase string }{
	"item.show":         {`{"id":"S-0001"}`, "item.show"},
	"item.move.preview": {`{"id":"E-0001"}`, "cancel.preview"},
	"accept.preview":    {`{"id":"S-0001"}`, "accept.preview"},
	"stream.diff":       {`{"id":"S-0001"}`, "stream.diff"},
	"stats.get":         {`{"since":"12w","type":"story","by":"nature","bucket":"week"}`, "stats.compute"},
	"publish.preview":   {`{}`, "release.pending"},
	// S-0217
	"order.by":           {`{"policy":"wsjf"}`, "order.by"},
	"promote.candidates": {`{"limit":3}`, "promote.candidates"},
	"release.evaluate":   {`{}`, "release.evaluate"},
	// S-0394
	"forge.status": {`{"id":"S-0001"}`, "forge.status"},
}

// readRefused is, per read, params refused before anything is read.
var readRefused = map[string][]string{
	"item.show":         {`{"id":"--show"}`, `{"id":"../S-0001"}`},
	"item.move.preview": {`{"id":"../x"}`},
	"accept.preview":    {`{"id":"--no-push"}`},
	"stream.diff":       {`{"id":"../../etc"}`},
	"stats.get":         {`{"since":"30d; ls"}`, `{"type":"folder"}`, `{"by":"owner"}`, `{"bucket":"minute"}`, `{"bucket":"--json"}`},
	"publish.preview":   {`"--force"`},
	// S-0217
	"order.by":           {`{"policy":"random"}`, `{"policy":"--apply"}`, `{"policy":3}`},
	"promote.candidates": {`{"limit":-1}`, `{"limit":"3"}`},
	"release.evaluate":   {`"--apply"`, `[1]`},
	// S-0394
	"forge.status": {`{"id":"T-0001"}`, `{"id":"E-0001"}`, `{"id":"../S-0001"}`, `{"id":"S-0001 --json"}`, `{}`, `[]`},
}

func TestEveryReadIsCovered(t *testing.T) {
	reads := readMethods(time.Now, Host{})
	for name := range reads {
		if _, ok := readPhase[name]; !ok {
			t.Errorf("%s has no call in readPhase", name)
		}
		if len(readRefused[name]) == 0 {
			t.Errorf("%s has no refusals in readRefused", name)
		}
		if _, ok := specs()[name]; ok {
			t.Errorf("%s is still in the write table, which starts flai for it", name)
		}
	}
	if len(readPhase) != len(reads) || len(readRefused) != len(reads) {
		t.Errorf("readPhase names %d, readRefused %d, there are %d reads", len(readPhase), len(readRefused), len(reads))
	}
}

func TestAReadRefusesWhatIsNotDataBeforeReadingAnything(t *testing.T) {
	// a root that does not exist: anything past validation fails otherwise
	p := channel.Project{Key: "nowhere", Root: "/nonexistent/flai-test"}
	reads := readMethods(time.Now, Host{})
	for name, all := range readRefused {
		for _, params := range all {
			_, rerr := reads[name](context.Background(), p, json.RawMessage(params))
			if rerr == nil || rerr.Code != channel.CodeInvalidParams {
				t.Errorf("%s %s: %+v, want invalid params", name, params, rerr)
			}
		}
	}
}

func TestAReadIsAnsweredInThisProcess(t *testing.T) {
	p := harbour(t)
	table := Methods("test", func() time.Time { return t0.Add(time.Hour) })
	for name, call := range readPhase {
		got := phases(t, func(ctx context.Context) error {
			// the harbour is not a git repository: stream.diff and
			// publish.preview fail, as flai itself does there
			_, _ = table[name](ctx, p, json.RawMessage(call.params))
			return nil
		})
		if got[call.phase] != 1 {
			t.Errorf("%s: phase %s counted %d times, want once; phases %v", name, call.phase, got[call.phase], got)
		}
		for phase := range got {
			if strings.HasPrefix(phase, "exec.flai") {
				t.Errorf("%s started flai: phases %v", name, got)
			}
		}
	}
}

func TestAReadAnswersAsFlaiDid(t *testing.T) {
	p := harbour(t)
	table := Methods("test", func() time.Time { return t0.Add(time.Hour) })
	call := func(name, params string) (Written, *channel.Error) {
		res, rerr := table[name](context.Background(), p, json.RawMessage(params))
		w, _ := res.(Written)
		return w, rerr
	}

	w, rerr := call("item.show", `{"id":"S-1"}`)
	var shown struct{ ID, Hash string }
	if rerr != nil || json.Unmarshal(w.Data, &shown) != nil || shown.ID != "S-0001" || len(shown.Hash) != 64 || w.Warnings == nil {
		t.Errorf("item.show: %s %+v %+v", w.Data, w.Warnings, rerr)
	}

	// an ID that names nothing is the caller's mistake, as it was from flai
	if _, rerr := call("item.show", `{"id":"S-0099"}`); rerr == nil || rerr.Code != NotFound || !strings.Contains(rerr.Message, "S-0099 not found") {
		t.Errorf("item.show of nothing: %+v", rerr)
	}
	// a story not in review cannot be accepted: flai's failure, said as it said it
	if _, rerr := call("accept.preview", `{"id":"S-0001"}`); rerr == nil || rerr.Code != channel.CodeInternal || rerr.Message != "S-0001 is ready; a story is accepted from review" {
		t.Errorf("accept.preview of a ready story: %+v", rerr)
	}

	w, rerr = call("item.move.preview", `{"id":"E-0001"}`)
	var cancel struct {
		DryRun    bool `json:"dry_run"`
		Cancelled []struct{ ID string }
	}
	if rerr != nil || json.Unmarshal(w.Data, &cancel) != nil || !cancel.DryRun || len(cancel.Cancelled) == 0 {
		t.Errorf("item.move.preview: %s %+v", w.Data, rerr)
	}

	// S-0217: no policy is the manifest's, fifo when it names none, and the
	// order is computed, never applied
	w, rerr = call("order.by", `{}`)
	var order struct {
		Policy  string
		Applied bool
		Stories []struct{ ID string }
	}
	if rerr != nil || json.Unmarshal(w.Data, &order) != nil || order.Policy != "fifo" || order.Applied || order.Stories == nil {
		t.Errorf("order.by: %s %+v", w.Data, rerr)
	}
	w, rerr = call("promote.candidates", `{}`)
	var promote struct {
		Policy string
		Others []struct{ ID string }
	}
	if rerr != nil || json.Unmarshal(w.Data, &promote) != nil || promote.Policy != "fifo" {
		t.Errorf("promote.candidates: %s %+v", w.Data, rerr)
	}
}

// S-0394: forge.status answers a story's pull request as flai serve kept
// it, without its events, and null for a story it has read none for.
func TestForgeStatusAnswersTheKeptPullRequest(t *testing.T) {
	p := harbour(t)
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	kept := forge.Status{
		Story: "S-0001", PullRequest: 7, URL: "https://forge.test/harbour/pull/7", HeadSHA: "abc123", ReadAt: "2026-10-17T10:00:00Z",
		Checks:  []forge.CheckStatus{{Name: "test", State: forge.CheckFailed, URL: "https://forge.test/runs/1", Since: "2026-10-17T09:59:00Z"}},
		Reviews: []forge.ReviewStatus{{Author: "olive", State: forge.ReviewChangesRequested, URL: "https://forge.test/harbour/pull/7#r1", Body: "Deeper."}},
		Events:  []forge.StatusEvent{{At: "2026-10-17T09:59:00Z", Kind: forge.EventCheck, Name: "test", To: "failed", Commit: "abc123"}},
	}
	if err := forge.WriteStatus(repo.CacheDir(), kept); err != nil {
		t.Fatal(err)
	}
	table := Methods("test", func() time.Time { return t0.Add(time.Hour) })
	call := func(params string) (Written, *channel.Error) {
		res, rerr := table["forge.status"](context.Background(), p, json.RawMessage(params))
		w, _ := res.(Written)
		return w, rerr
	}

	// any zero padding names the story
	w, rerr := call(`{"id":"S-1"}`)
	if rerr != nil {
		t.Fatalf("forge.status of a story with kept state: %+v", rerr)
	}
	var got struct {
		Story       string          `json:"story"`
		PullRequest json.RawMessage `json:"pull_request"`
	}
	if err := json.Unmarshal(w.Data, &got); err != nil {
		t.Fatalf("%v\n%s", err, w.Data)
	}
	var pr forge.Status
	if err := json.Unmarshal(got.PullRequest, &pr); err != nil {
		t.Fatalf("%v\n%s", err, w.Data)
	}
	if got.Story != "S-0001" || pr.PullRequest != 7 || pr.URL != kept.URL || pr.HeadSHA != "abc123" ||
		len(pr.Checks) != 1 || pr.Checks[0].State != forge.CheckFailed || len(pr.Reviews) != 1 || pr.Reviews[0].Body != "Deeper." {
		t.Errorf("forge.status: %s", w.Data)
	}
	if strings.Contains(string(w.Data), `"events"`) {
		t.Errorf("forge.status answered the events: %s", w.Data)
	}

	// a story flai serve has read no pull request for
	w, rerr = call(`{"id":"S-0002"}`)
	if rerr != nil || string(w.Data) != `{"story":"S-0002","pull_request":null}` {
		t.Errorf("forge.status of a story without kept state: %s %+v", w.Data, rerr)
	}

	// an ID that names nothing is the caller's mistake, as for every read
	if _, rerr := call(`{"id":"S-0099"}`); rerr == nil || rerr.Code != NotFound || !strings.Contains(rerr.Message, "S-0099 not found") {
		t.Errorf("forge.status of nothing: %+v", rerr)
	}
	// a task or an epic has no pull request: refused, saying what to give
	for _, id := range []string{"T-0001", "E-1"} {
		if _, rerr := call(`{"id":"` + id + `"}`); rerr == nil || rerr.Code != channel.CodeInvalidParams || !strings.Contains(rerr.Message, "give a story's ID") {
			t.Errorf("forge.status of %s: %+v", id, rerr)
		}
	}

	// a corrupt file is a failure that says what to do, with the path in the project
	if err := os.WriteFile(forge.StatusPath(repo.CacheDir(), "S-0002"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, rerr := call(`{"id":"S-0002"}`); rerr == nil || rerr.Code != channel.CodeInternal ||
		!strings.Contains(rerr.Message, "delete it") || !strings.Contains(rerr.Message, filepath.Join(".flai-cache", "forge", "status", "S-0002.json")) || strings.Contains(rerr.Message, p.Root) {
		t.Errorf("forge.status of a corrupt file: %+v", rerr)
	}
}

// reviewing is harbour with Berths, S-0001, in review, the thread on it
// open and no flai verify result, and a batch in the merge queue's state.
func reviewing(t *testing.T) (channel.Project, *workitem.Repo) {
	t.Helper()
	p := harbour(t)
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	toReview(t, repo, "S-0001", t0.Add(time.Hour))
	state := mergequeue.StatePath(repo.CacheDir())
	_ = os.MkdirAll(filepath.Dir(state), 0o755)
	if err := os.WriteFile(state, []byte(`{"enqueued":[],"batch":{"stories":["S-0001"],"step":"verify"},"outcomes":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	return p, repo
}

// toReview moves story, ready, to in-progress and then to review at at.
func toReview(t *testing.T, repo *workitem.Repo, story string, at time.Time) {
	t.Helper()
	for _, st := range []string{workitem.InProgress, workitem.Review} {
		it, err := repo.Get(story)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := repo.Transition(it, st, "olive", "", at); err != nil {
			t.Fatalf("%s to %s: %v", story, st, err)
		}
	}
}

// boardQueue is board.get's queue at p, and the queue mergequeue.Compute
// gives for repo, both as JSON decodes them.
func boardQueue(t *testing.T, p channel.Project, repo *workitem.Repo) (got map[string]any, want any) {
	t.Helper()
	var v struct {
		Queue map[string]any `json:"queue"`
	}
	if err := call(t, p, "board.get", `{}`, &v); err != nil {
		t.Fatal(err)
	}
	q, err := mergequeue.Compute(execx.System{}, repo)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(q)
	_ = json.Unmarshal(data, &want)
	return v.Queue, want
}

// S-0398, ADR-0173: board.get carries the merge queue as flai queue gives
// it: the reasons of a story in review that is not eligible, and the batch
// and its step. Under pull-request acceptance the queue is off, with the
// reason, and lists no story.
func TestBoardGetCarriesTheMergeQueue(t *testing.T) {
	p, repo := reviewing(t)
	got, want := boardQueue(t, p, repo)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("board.get queue:\n%v\nwant flai queue's:\n%v", got, want)
	}
	in, _ := got["ineligible"].([]any)
	if el, _ := got["eligible"].([]any); got["mode"] != "local" || el == nil || len(el) != 0 || len(in) != 1 {
		t.Fatalf("queue: %v", got)
	}
	var kinds []string
	for _, r := range in[0].(map[string]any)["reasons"].([]any) {
		kinds = append(kinds, r.(map[string]any)["kind"].(string))
	}
	if strings.Join(kinds, " ") != "verify_missing thread_open" {
		t.Errorf("S-0001's reasons: %v", kinds)
	}
	if b, _ := got["batch"].(map[string]any); b["step"] != "verify" || !reflect.DeepEqual(b["stories"], []any{"S-0001"}) {
		t.Errorf("batch: %v", got["batch"])
	}

	m := filepath.Join(p.Root, "system-flow.yaml")
	data, _ := os.ReadFile(m)
	_ = os.WriteFile(m, append(data, "acceptance:\n  mode: pull-request\n"...), 0o644)
	got, _ = boardQueue(t, p, repo)
	if el, _ := got["eligible"].([]any); got["off"] != true || got["off_reason"] != mergequeue.OffUnderPullRequest || len(el) != 0 || got["batch"] != nil {
		t.Errorf("under pull-request: %v", got)
	}
}

// S-0398: board.get lists a story in review whose last flai verify passed
// at its branch's head, with every criterion ticked and no thread open, as
// eligible, ahead of none, and the other with its reasons.
func TestBoardGetListsTheEligibleStories(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	p, repo := reviewing(t)
	r := execx.System{}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "one"}, {"branch", storygit.Branch("S-0002")}} {
		if out, err := r.Run(p.Root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	head, err := r.Run(p.Root, "git", "rev-parse", storygit.Branch("S-0002"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Lift", Parent: "S-0002", Owner: "olive", Now: t0}); err != nil {
		t.Fatal(err)
	}
	cranes, _ := repo.Get("S-0002")
	if _, err := repo.Transition(cranes, workitem.Ready, "olive", "", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	toReview(t, repo, "S-0002", t0.Add(2*time.Hour))
	if err := verify.SaveReport(repo, verify.Report{Story: "S-0002", Commit: strings.TrimSpace(head), Passed: true}); err != nil {
		t.Fatal(err)
	}

	got, want := boardQueue(t, p, repo)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("board.get queue:\n%v\nwant flai queue's:\n%v", got, want)
	}
	el, _ := got["eligible"].([]any)
	in, _ := got["ineligible"].([]any)
	if len(el) != 1 || el[0].(map[string]any)["story"] != "S-0002" || len(in) != 1 || in[0].(map[string]any)["story"] != "S-0001" {
		t.Errorf("queue: %v", got)
	}
}

// ADR-0221: while the merge queue batches, with the merge-queue host action
// on and accept_reviews on under local acceptance, accept.preview lists the
// blocker queue_tiers for a story whose last flai verify left a tier marked
// queue at its branch's head, as flai accept --dry-run does; with the host
// action off it lists none.
func TestAcceptPreviewListsTheQueueTiersWhileTheQueueBatches(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	p, repo := reviewing(t)
	m := filepath.Join(p.Root, "system-flow.yaml")
	data, _ := os.ReadFile(m)
	if err := os.WriteFile(m, append(data, "orchestration:\n  permissions:\n    accept_reviews: true\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	r := execx.System{}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"config", "user.email", "t@t"}, {"config", "user.name", "t"}, {"add", "-A"}, {"commit", "-q", "-m", "one"}, {"branch", storygit.Branch("S-0001")}} {
		if out, err := r.Run(p.Root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	head, err := r.Run(p.Root, "git", "rev-parse", storygit.Branch("S-0001"))
	if err != nil {
		t.Fatal(err)
	}
	if err := verify.SaveReport(repo, verify.Report{Story: "S-0001", Commit: strings.TrimSpace(head), Passed: true, Steps: []verify.Step{
		{Name: "integration", Tier: true, State: verify.Queued},
	}}); err != nil {
		t.Fatal(err)
	}
	preview := func(on bool) map[string]any {
		t.Helper()
		host := Host{Enabled: func(action, root string) bool { return on && action == ActionMergeQueue && root == p.Root }}
		out, rerr := readMethods(time.Now, host)["accept.preview"](context.Background(), p, json.RawMessage(`{"id":"S-0001"}`))
		if rerr != nil {
			t.Fatalf("accept.preview: %+v", rerr)
		}
		data, err := json.Marshal(out)
		if err != nil {
			t.Fatal(err)
		}
		var got struct {
			Data map[string]any `json:"data"`
		}
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		return got.Data
	}
	got := preview(true)
	if b, _ := got["queue_tiers_blocker"].(map[string]any); b["code"] != "queue_tiers" || !reflect.DeepEqual(got["queue_tiers"], []any{"integration"}) {
		t.Errorf("with the queue batching, accept.preview lists queue_tiers naming integration: %v", got)
	}
	if got := preview(false); got["queue_tiers_blocker"] != nil || got["queue_tiers"] != nil {
		t.Errorf("with the merge-queue host action off, accept.preview lists no queue tiers: %v", got)
	}
}

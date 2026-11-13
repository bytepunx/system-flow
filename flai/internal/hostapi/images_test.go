package hostapi

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// fakeImageJobs is Host.Images over a fixed state: the jobs it answers, the
// log of each by id, and what Ensure answers, with each Ensure it was asked.
type fakeImageJobs struct {
	state   ImagesState
	logs    map[string]string // a job's log's path, by id
	started ImageStarted
	err     error
	asked   []string // each Ensure, as key and root
}

func (f *fakeImageJobs) Images(string) ImagesState { return f.state }

func (f *fakeImageJobs) JobLog(_, id string) (ImageJob, string, bool) {
	jobs := append([]ImageJob{}, f.state.Queued...)
	jobs = append(jobs, f.state.Ended...)
	if f.state.Running != nil {
		jobs = append(jobs, *f.state.Running)
	}
	for _, j := range jobs {
		if j.ID == id {
			return j, f.logs[id], true
		}
	}
	return ImageJob{}, "", false
}

func (f *fakeImageJobs) Ensure(key, root string) (ImageStarted, error) {
	f.asked = append(f.asked, key+" "+root)
	return f.started, f.err
}

// imagesHost is a host whose agent action is on when on says so,
// journalling into journal, with jobs as its image jobs.
func imagesHost(on *bool, journal *[]Entry, jobs ImageJobs) Host {
	return Host{
		Enabled: func(action, _ string) bool { return *on && action == ActionAgent },
		Record:  func(e Entry) { *journal = append(*journal, e) },
		Images:  jobs,
	}
}

// harnessStatus is the status of the default claude-code harness image, in state.
func harnessStatus(state agentimage.State) agentimage.Status {
	return agentimage.Status{Image: agentimage.Image{Kind: agentimage.KindHarness, Inner: "claude-code", Ref: "flai-agent-claude-code:1.2.3", Source: agentimage.SourceDefault}, State: state}
}

func TestImagesMethodsAreTwoReadsAndAWriteUnderTheAgentAction(t *testing.T) {
	for _, name := range []string{"images.get", "images.log"} {
		if sp := imageSpecs()[name]; sp.action != "" || !sp.reads || sp.do == nil {
			t.Errorf("%s is a read, answered whatever the action says: %+v", name, sp)
		}
	}
	if sp := imageSpecs()["images.ensure"]; sp.action != ActionAgent || sp.reads || sp.do == nil {
		t.Errorf("images.ensure is a write under the agent action, answered in flai serve: %+v", imageSpecs()["images.ensure"])
	}
}

func TestImagesMethodsRefuseWhereNoImageJobsRun(t *testing.T) {
	p := shipyard(t)
	on := true
	journal := []Entry{}
	host := imagesHost(&on, &journal, nil)
	for method, params := range map[string]string{"images.get": `{}`, "images.log": `{"id":"20261113T093000Z-001","offset":0}`, "images.ensure": reqID(1)} {
		e := ask(t, host, p, method, params, nil)
		if e == nil || e.Code != Rule || !strings.Contains(e.Message, "flai serve runs them") {
			t.Errorf("%s with no image jobs: %v, want a refusal that says flai serve runs them", method, e)
		}
	}
}

func TestImagesGetWithNoJobsAnswersTheImagesAndEmptyLists(t *testing.T) {
	p := shipyard(t)
	jobs := &fakeImageJobs{state: ImagesState{Flai: "1.2.3", Images: []agentimage.Status{harnessStatus(agentimage.StateMissing)}}}
	for _, on := range []bool{true, false} {
		journal := []Entry{}
		var raw json.RawMessage
		if e := ask(t, imagesHost(&on, &journal, jobs), p, "images.get", `{}`, &raw); e != nil {
			t.Fatalf("on %v: %v", on, e)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		for name, want := range map[string]string{"flai": `"1.2.3"`, "running": `null`, "queued": `[]`, "ended": `[]`} {
			if string(fields[name]) != want {
				t.Errorf("on %v: %s is %s, want %s", on, name, fields[name], want)
			}
		}
		var got ImagesState
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatal(err)
		}
		if len(got.Images) != 1 || got.Images[0].Ref != "flai-agent-claude-code:1.2.3" || got.Images[0].State != agentimage.StateMissing {
			t.Errorf("on %v: images %+v, want the missing harness image", on, got.Images)
		}
		if len(journal) != 0 {
			t.Errorf("a read journals nothing: %+v", journal)
		}
	}
}

func TestImagesGetAnswersTheRunningJobAndTheEndedOnesWithoutTheHostsPaths(t *testing.T) {
	p := shipyard(t)
	running := ImageJob{ID: "20261113T093000Z-002", Images: []string{"flai-agent-claude-code:1.2.3"}, State: ImageJobRunning,
		Step: "flai-agent-claude-code:1.2.3 missing; pulling claude-code of flai 1.2.3", Requested: "2026-11-13T09:30:00Z", Started: "2026-11-13T09:30:01Z", LogSize: 120}
	failedJob := ImageJob{ID: "20261113T080000Z-001", Images: []string{"flai-agent-claude-code:1.2.3", "shipyard-agent:latest"}, State: ImageJobEnded,
		Ways: []ImageWay{{Image: "flai-agent-claude-code:1.2.3", Way: "pulled"}}, Step: "build shipyard-agent:latest",
		Requested: "2026-11-13T08:00:00Z", Started: "2026-11-13T08:00:01Z", Ended: "2026-11-13T08:02:01Z", Outcome: "failed",
		Cause: "build shipyard-agent:latest failed in " + filepath.Join(p.Root, "agent.Dockerfile") + " for " + p.Root, DurationMS: 120000, LogSize: 4096}
	jobs := &fakeImageJobs{state: ImagesState{Images: []agentimage.Status{harnessStatus(agentimage.StateMissing)}, Running: &running, Ended: []ImageJob{failedJob}}}
	on := false
	journal := []Entry{}
	var got ImagesState
	if e := ask(t, imagesHost(&on, &journal, jobs), p, "images.get", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	running.Ways = []ImageWay{}
	if got.Running == nil || !reflect.DeepEqual(*got.Running, running) {
		t.Errorf("running:\n got %+v\nwant %+v", got.Running, running)
	}
	if len(got.Ended) != 1 {
		t.Fatalf("ended: %+v, want the failed job", got.Ended)
	}
	ended := got.Ended[0]
	if ended.ID != failedJob.ID || ended.Outcome != "failed" || ended.Step != "build shipyard-agent:latest" || ended.Ended != failedJob.Ended || len(ended.Ways) != 1 || ended.Ways[0].Way != "pulled" {
		t.Errorf("the ended job: %+v", ended)
	}
	if want := "build shipyard-agent:latest failed in agent.Dockerfile for the project's main checkout"; ended.Cause != want {
		t.Errorf("the cause is %q, want %q", ended.Cause, want)
	}
	if len(got.Queued) != 0 {
		t.Errorf("queued: %+v, want none", got.Queued)
	}
}

func TestImagesGetAnswersUnjudgedImagesWithWhy(t *testing.T) {
	p := shipyard(t)
	jobs := &fakeImageJobs{state: ImagesState{Images: []agentimage.Status{{Image: harnessStatus("").Image}}, Unjudged: "docker is not on PATH in " + p.Root}}
	on := true
	journal := []Entry{}
	var got ImagesState
	if e := ask(t, imagesHost(&on, &journal, jobs), p, "images.get", `{}`, &got); e != nil {
		t.Fatal(e)
	}
	if got.Unjudged != "docker is not on PATH in the project's main checkout" || len(got.Images) != 1 || got.Images[0].State != "" {
		t.Errorf("unjudged: %+v", got)
	}
}

func TestImagesLogReadsAJobsLogFromAnOffset(t *testing.T) {
	p := shipyard(t)
	dir := t.TempDir()
	logPath := filepath.Join(dir, "20261113T093000Z-001.log")
	if err := os.WriteFile(logPath, []byte("one\ntwo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	ended := ImageJob{ID: "20261113T093000Z-001", State: ImageJobEnded, Outcome: "failed"}
	queued := ImageJob{ID: "20261113T093000Z-003", State: ImageJobQueued}
	jobs := &fakeImageJobs{state: ImagesState{Ended: []ImageJob{ended}, Queued: []ImageJob{queued}},
		logs: map[string]string{ended.ID: logPath, queued.ID: filepath.Join(dir, queued.ID+".log")}}
	on := false
	journal := []Entry{}
	host := imagesHost(&on, &journal, jobs)
	read := func(params string) ImageLog {
		t.Helper()
		var got ImageLog
		if e := ask(t, host, p, "images.log", params, &got); e != nil {
			t.Fatalf("%s: %v", params, e)
		}
		return got
	}
	if got, want := read(`{"id":"20261113T093000Z-001","offset":0}`), (ImageLog{ID: ended.ID, State: ImageJobEnded, Outcome: "failed", Exists: true, Offset: 0, Next: 8, Size: 8, Text: "one\ntwo\n"}); got != want {
		t.Errorf("from 0\n got %+v\nwant %+v", got, want)
	}
	if got := read(`{"id":"20261113T093000Z-001","offset":4}`); got.Text != "two\n" || got.Offset != 4 || got.Next != 8 || got.Reset {
		t.Errorf("from the middle: %+v", got)
	}
	if got := read(`{"id":"20261113T093000Z-001","offset":8}`); got.Text != "" || got.Next != 8 || got.Reset {
		t.Errorf("at the end: %+v", got)
	}
	if got := read(`{"id":"20261113T093000Z-001","offset":100}`); !got.Reset || got.Offset != 0 || got.Text != "one\ntwo\n" || got.Next != 8 {
		t.Errorf("past the end, a log cleared since: %+v, want it read again from 0 with reset", got)
	}
	if got := read(`{"id":"20261113T093000Z-003","offset":0}`); got.Exists || got.State != ImageJobQueued || got.Text != "" || got.Next != 0 {
		t.Errorf("a queued job's log, not there yet: %+v", got)
	}
	if e := ask(t, host, p, "images.log", `{"id":"20261113T093000Z-009","offset":0}`, nil); e == nil || e.Code != NotFound || !strings.Contains(e.Message, "20261113T093000Z-009") {
		t.Errorf("an unknown job: %v, want not found, naming it", e)
	}
	for _, params := range []string{`{"id":"../agents","offset":0}`, `{"id":"","offset":0}`, `{"id":"20261113T093000Z-001","offset":-1}`} {
		if e := ask(t, host, p, "images.log", params, nil); e == nil || e.Code != channel.CodeInvalidParams {
			t.Errorf("%s: %v, want invalid params", params, e)
		}
	}
	if len(journal) != 0 {
		t.Errorf("a read journals nothing: %+v", journal)
	}
}

func TestReadImageLogEndsOnAWholeCharacter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "job.log")
	text := strings.Repeat("a", imageLogChunk-1) + "é"
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := readImageLog(path, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got.Next != imageLogChunk-1 || strings.Contains(got.Text, "\xc3") {
		t.Errorf("a part cut in a character's middle: next %d, want %d, before the character", got.Next, imageLogChunk-1)
	}
	rest, err := readImageLog(path, got.Next)
	if err != nil {
		t.Fatal(err)
	}
	if rest.Text != "é" || rest.Next != rest.Size {
		t.Errorf("the next part: %+v", rest)
	}
}

func TestImagesEnsureIsRefusedWhileTheAgentActionIsOff(t *testing.T) {
	p := shipyard(t)
	jobs := &fakeImageJobs{}
	on := false
	journal := []Entry{}
	e := ask(t, imagesHost(&on, &journal, jobs), p, "images.ensure", reqID(1), nil)
	if e == nil || e.Code != Disabled || !strings.Contains(e.Message, EnableCommand(ActionAgent)) {
		t.Errorf("with the agent action off: %v, want disabled, naming %s", e, EnableCommand(ActionAgent))
	}
	if len(jobs.asked) != 0 {
		t.Errorf("flai serve was asked for a job with the action off: %v", jobs.asked)
	}
	if len(journal) != 1 || journal[0].Outcome != "disabled" || journal[0].Action != ActionAgent {
		t.Errorf("journal: %+v, want the refusal", journal)
	}
}

func TestImagesEnsureIsRefusedWhileAJobForTheSameImagesRuns(t *testing.T) {
	p := shipyard(t)
	busy := ImageJob{ID: "20261113T093000Z-002", Images: []string{"flai-agent-claude-code:1.2.3"}, State: ImageJobRunning, Started: "2026-11-13T09:30:01Z"}
	jobs := &fakeImageJobs{err: &ImageJobBusy{Job: busy}}
	on := true
	journal := []Entry{}
	e := ask(t, imagesHost(&on, &journal, jobs), p, "images.ensure", reqID(1), nil)
	if e == nil || e.Code != Rule {
		t.Fatalf("while a job for the same images runs: %v, want a rule's refusal", e)
	}
	for _, said := range []string{busy.ID, "flai-agent-claude-code:1.2.3", "runs since 2026-11-13T09:30:01Z", "no other job is started", "images.log"} {
		if !strings.Contains(e.Message, said) {
			t.Errorf("the refusal %q does not say %q", e.Message, said)
		}
	}
	data, _ := e.Data.(map[string]any)
	if job, _ := data["job"].(ImageJob); job.ID != busy.ID {
		t.Errorf("the refusal's data %+v does not carry the running job", e.Data)
	}
	if len(journal) != 1 || journal[0].Outcome != "failed" || !strings.Contains(journal[0].Detail, busy.ID) {
		t.Errorf("journal: %+v, want the refusal, naming the job", journal)
	}
}

func TestImagesEnsureIsRefusedForAProjectWhoseAgentsNeedNoImage(t *testing.T) {
	p := shipyard(t)
	on := true
	journal := []Entry{}
	e := ask(t, imagesHost(&on, &journal, &fakeImageJobs{err: ErrNoImages}), p, "images.ensure", reqID(1), nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "no image job is started") {
		t.Errorf("with no image to make: %v, want a rule's refusal", e)
	}
	e = ask(t, imagesHost(&on, &journal, &fakeImageJobs{err: errors.New("open the project at " + p.Root + ": no manifest")}), p, "images.ensure", reqID(2), nil)
	if e == nil || e.Code != channel.CodeInternal || strings.Contains(e.Message, p.Root) {
		t.Errorf("a project that cannot be read: %v, want a failure without the host's path", e)
	}
}

func TestImagesEnsureStartsAJob(t *testing.T) {
	p := shipyard(t)
	job := ImageJob{ID: "20261113T093000Z-004", Images: []string{"flai-agent-claude-code:1.2.3"}, State: ImageJobQueued, Requested: "2026-11-13T09:30:00Z"}
	jobs := &fakeImageJobs{started: ImageStarted{ImageJob: job, Ahead: 1}}
	on := true
	journal := []Entry{}
	host := imagesHost(&on, &journal, jobs)
	var w Written
	if e := ask(t, host, p, "images.ensure", reqID(1), &w); e != nil {
		t.Fatal(e)
	}
	var got ImageStarted
	if err := json.Unmarshal(w.Data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != job.ID || got.State != ImageJobQueued || got.Ahead != 1 || len(got.Images) != 1 {
		t.Errorf("started: %+v", got)
	}
	if want := []string{"shipyard " + p.Root}; !reflect.DeepEqual(jobs.asked, want) {
		t.Errorf("asked of flai serve: %v, want %v", jobs.asked, want)
	}
	if len(journal) != 1 || journal[0].Action != ActionAgent || journal[0].Method != "images.ensure" || journal[0].Outcome != "done" || !strings.Contains(journal[0].Detail, job.ID+" for flai-agent-claude-code:1.2.3: queued with 1 ahead") {
		t.Errorf("journal: %+v", journal)
	}
}

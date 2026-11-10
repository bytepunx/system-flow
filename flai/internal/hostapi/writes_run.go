package hostapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Run is one flai invocation: what to run, what to give it on standard
// input, and where each log event it writes goes while it runs.
type Run struct {
	Dir     string
	Args    []string
	Stdin   string
	OnEvent func(map[string]any)
}

// Ran is what came back.
type Ran struct {
	Stdout []byte
	Events []map[string]any // the log events from standard error
	Exit   int
}

// Runner runs flai. The real one runs this executable; tests give their own.
type Runner func(ctx context.Context, r Run) (Ran, error)

// ExecRunner runs the flai that is running, as the dashboard would have run
// the one in its container: JSON out, log events on standard error.
func ExecRunner(ctx context.Context, r Run) (Ran, error) {
	exe, err := os.Executable()
	if err != nil {
		return Ran{}, err
	}
	return runFlai(ctx, exe, r)
}

// runFlai runs the flai executable exe as ExecRunner does: a test gives the
// flai it built from this tree.
func runFlai(ctx context.Context, exe string, r Run) (Ran, error) {
	cmd := exec.CommandContext(ctx, exe, r.Args...)
	cmd.Dir = r.Dir
	cmd.Env = append(os.Environ(), "FLAI_AGENT=flaiover", "FLAI_SESSION=", "LOG_FORMAT=json")
	cmd.Stdin = strings.NewReader(r.Stdin)
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 10 * time.Second
	var out bytes.Buffer
	cmd.Stdout = &out
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return Ran{}, err
	}
	if err := cmd.Start(); err != nil {
		return Ran{}, err
	}
	var ran Ran
	sc := bufio.NewScanner(stderr)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for sc.Scan() {
		var ev map[string]any
		if json.Unmarshal(sc.Bytes(), &ev) != nil {
			ev = map[string]any{"level": "INFO", "msg": sc.Text()}
		}
		ran.Events = append(ran.Events, ev)
		if r.OnEvent != nil {
			r.OnEvent(ev)
		}
	}
	err = cmd.Wait()
	ran.Stdout = out.Bytes()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		ran.Exit = exit.ExitCode()
		return ran, nil
	}
	return ran, err
}

// itemsChanged tells the item store that a flai run in another process may
// have written items under root, and under dir when it ran in another folder
// such as a story's worktree, so that a snapshot flai serve keeps of them is
// looked at again on its next list rather than when the watcher reports the
// files 300 to 600 ms later (S-0606).
func itemsChanged(root, dir string) {
	workitem.ItemsChanged(root)
	if dir != "" && dir != root {
		workitem.ItemsChanged(dir)
	}
}

// Written is what a command answered: its JSON, and the warnings it logged.
type Written struct {
	Data     json.RawMessage `json:"data"`
	Warnings []string        `json:"warnings"`
}

// journal remembers the writes that completed, by the request ID the
// dashboard gave, so that a request repeated after a lost connection is
// answered with what happened and not done twice.
type journal struct {
	mu   sync.Mutex
	done map[string]journalled
	now  func() time.Time
}

type journalled struct {
	at  time.Time
	res any
	err *channel.Error
}

const journalKeeps = 10 * time.Minute

// maxWIPLimit bounds a WIP limit set from the dashboard: past it a limit limits nothing.
const maxWIPLimit = 99

func (j *journal) get(key string) (journalled, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	for k, v := range j.done {
		if j.now().Sub(v.at) > journalKeeps {
			delete(j.done, k)
		}
	}
	v, ok := j.done[key]
	return v, ok
}

func (j *journal) put(key string, res any, err *channel.Error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.done[key] = journalled{at: j.now(), res: res, err: err}
}

// Request is the record of one detached write, kept by its request ID from
// before it acts, so that a repeat is answered and not done again even when
// the write ended the flai serve that took it (S-0109).
type Request struct {
	Started time.Time      `json:"started"`
	PID     int            `json:"pid"`   // the flai serve that took it
	Until   time.Time      `json:"until"` // forgotten after this
	Done    bool           `json:"done"`
	Result  *Written       `json:"result,omitempty"`
	Err     *channel.Error `json:"error,omitempty"`
}

// Requests changes the records of detached writes under a lock, and keeps
// what change leaves.
type Requests func(change func(map[string]Request)) error

// inMemory keeps the records for as long as the table lives.
func inMemory() Requests {
	var mu sync.Mutex
	kept := map[string]Request{}
	return func(change func(map[string]Request)) error {
		mu.Lock()
		defer mu.Unlock()
		change(kept)
		return nil
	}
}

// begin records a detached write as started, unless its key is recorded
// already: then it returns that record, and true. Records past their time
// are forgotten first.
func (rs Requests) begin(key string, now time.Time, keep time.Duration) (was Request, found bool, err error) {
	err = rs(func(all map[string]Request) {
		forgetOld(all, now)
		if was, found = all[key]; !found {
			all[key] = Request{Started: now, PID: os.Getpid(), Until: now.Add(keep)}
		}
	})
	return was, found, err
}

// finish records what a detached write answered.
func (rs Requests) finish(key string, now time.Time, res any, rerr *channel.Error) error {
	return rs(func(all map[string]Request) {
		forgetOld(all, now)
		r := all[key]
		r.Done, r.Until, r.Err = true, now.Add(journalKeeps), rerr
		if w, ok := res.(Written); ok {
			r.Result = &w
		}
		all[key] = r
	})
}

// unrecorded adds to an answer that its outcome could not be recorded, so a
// repeat will not be told it.
func unrecorded(res any, rerr *channel.Error, err error) (any, *channel.Error) {
	say := "the outcome could not be recorded, so a repeat of this request will not be told it: " + err.Error()
	if w, ok := res.(Written); ok {
		w.Warnings = append(w.Warnings, say)
		return w, rerr
	}
	if rerr != nil {
		rerr.Message += "; " + say
	}
	return res, rerr
}

func forgetOld(all map[string]Request, now time.Time) {
	for k, r := range all {
		if now.After(r.Until) {
			delete(all, k)
		}
	}
}

// repeated answers a repeat of a detached write from its record: what it
// answered, once it has; until then, that it is under way in this flai serve,
// or that the flai serve that took it ended before it could say, which is
// what a restart of flai serve itself looks like. It is never done again.
func repeated(id string, r Request) (any, *channel.Error) {
	if r.Done {
		if r.Err != nil || r.Result == nil {
			return nil, r.Err
		}
		return *r.Result, nil
	}
	started := r.Started.UTC().Format(time.RFC3339)
	underWay := r.PID == os.Getpid()
	say := fmt.Sprintf("request %s was started at %s and is still under way; it was not started again", id, started)
	if !underWay {
		say = fmt.Sprintf("request %s was started at %s by a flai serve (pid %d) that ended before it could say how it went; it was not done again", id, started, r.PID)
	}
	data, _ := json.Marshal(map[string]any{"repeat": map[string]any{"request_id": id, "started": started, "under_way": underWay}})
	return Written{Data: data, Warnings: []string{say}}, nil
}

var requestID = regexp.MustCompile(`^[A-Za-z0-9._-]{8,64}$`)

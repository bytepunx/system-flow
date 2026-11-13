package agentimage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// errExit stands for a docker client's non-zero exit.
var errExit = errors.New("exit status 1")

// try is what the fake docker does on one try of a step: write out, then
// block until its context ends when block is set, and return err.
type try struct {
	out   string
	err   error
	block bool
	// during runs after out is written, before the try returns.
	during func(out io.Writer)
}

// fakeDocker is an Exec scripted per try, which records how it was called.
type fakeDocker struct {
	t       *testing.T
	mu      sync.Mutex
	tries   []try
	calls   int
	program string
	dir     string
	env     []string
	args    []string
}

// exec is the fake's Exec.
func (f *fakeDocker) exec(ctx context.Context, dir string, env []string, out io.Writer, program string, args ...string) error {
	f.mu.Lock()
	if f.calls >= len(f.tries) {
		f.mu.Unlock()
		f.t.Errorf("docker ran %d times, but the script has %d tries", f.calls+1, len(f.tries))
		return errExit
	}
	tr := f.tries[f.calls]
	f.calls++
	f.program, f.dir, f.env, f.args = program, dir, env, args
	f.mu.Unlock()
	if tr.out != "" {
		if _, err := io.WriteString(out, tr.out); err != nil {
			return err
		}
	}
	if tr.during != nil {
		tr.during(out)
	}
	if tr.block {
		<-ctx.Done()
		return ctx.Err()
	}
	return tr.err
}

// ran is how many times the fake ran.
func (f *fakeDocker) ran() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

// syncBuffer is a bytes.Buffer safe to write and read from two goroutines.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write appends p.
func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// String is what was written.
func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// sleeps is a Sleep that records each wait and returns at once.
type sleeps struct {
	mu    sync.Mutex
	waits []time.Duration
}

// sleep records d.
func (s *sleeps) sleep(_ context.Context, d time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.waits = append(s.waits, d)
	return nil
}

// S-0740: a step's output reaches Out as it is written, not when the step
// ends, and the step runs docker with its folder, environment, and arguments.
func TestRunStreamsOutputAsItIsWritten(t *testing.T) {
	out := &syncBuffer{}
	var seen string
	docker := &fakeDocker{t: t, tries: []try{{
		out:    "#1 [internal] load build definition\n",
		during: func(io.Writer) { seen = out.String() },
	}}}
	r := Runner{Exec: docker.exec, Out: out, Sleep: (&sleeps{}).sleep}
	step := Step{Name: "build flai-agent-bedrock:1.2.3", Dir: "/project", Env: []string{"DOCKER_BUILDKIT=1"}, Args: []string{"build", "."}}
	if err := r.Run(context.Background(), step); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if seen != "#1 [internal] load build definition\n" {
		t.Errorf("Out held %q while the step ran, want its first line already there", seen)
	}
	if docker.program != "docker" || docker.dir != "/project" || !reflect.DeepEqual(docker.env, step.Env) || !reflect.DeepEqual(docker.args, step.Args) {
		t.Errorf("ran %s %v in %s with %v, want docker %v in /project with %v", docker.program, docker.args, docker.dir, docker.env, step.Args, step.Env)
	}
}

// S-0740: a step that writes nothing for the stall limit is stopped and
// reported, naming the step and the limit, and is not tried again.
func TestRunStopsASilentStep(t *testing.T) {
	docker := &fakeDocker{t: t, tries: []try{{out: "#5 [stage 2/4] RUN go mod download\n", block: true}}}
	waits := &sleeps{}
	r := Runner{Exec: docker.exec, Stall: 50 * time.Millisecond, Sleep: waits.sleep}
	err := r.Run(context.Background(), Step{Name: "build flai-agent-bedrock:1.2.3"})
	var stall *StallError
	if !errors.As(err, &stall) {
		t.Fatalf("Run = %v, want a *StallError", err)
	}
	if stall.Step != "build flai-agent-bedrock:1.2.3" || stall.Limit != 50*time.Millisecond {
		t.Errorf("StallError = %+v, want the step's name and a limit of 50ms", stall)
	}
	if msg := err.Error(); !strings.Contains(msg, "build flai-agent-bedrock:1.2.3 wrote nothing for 50ms") || !strings.Contains(msg, "run it again") {
		t.Errorf("message %q does not name the step, the limit, and what to do", msg)
	}
	if docker.ran() != 1 || len(waits.waits) != 0 {
		t.Errorf("docker ran %d times and waited %v, want once and no wait", docker.ran(), waits.waits)
	}
}

// S-0740: a step that keeps writing is not stopped, however long it runs
// in all, because each write starts the stall limit again.
func TestRunKeepsAStepThatWrites(t *testing.T) {
	docker := &fakeDocker{t: t, tries: []try{{during: func(out io.Writer) {
		for range 8 {
			if _, err := io.WriteString(out, ".\n"); err != nil {
				t.Errorf("write: %v", err)
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}}}}
	r := Runner{Exec: docker.exec, Stall: 100 * time.Millisecond}
	if err := r.Run(context.Background(), Step{Name: "pull bedrock"}); err != nil {
		t.Fatalf("Run = %v, want a step that writes every 20ms to pass a 100ms stall limit", err)
	}
}

// networkOutputs are the last output of a step that failed on the network,
// one for each kind ADR-0295 §5 tries again.
var networkOutputs = map[string]string{
	"connection reset": "go: downloading golang.org/x/sys v0.30.0\nread tcp 172.17.0.2:41234->142.250.74.17:443: read: connection reset by peer\nERROR: failed to solve: process \"/bin/sh -c go mod download\" did not complete successfully: exit code: 1\n",
	"i/o timeout":      "Error response from daemon: Get \"https://ghcr.io/v2/\": dial tcp 140.82.121.34:443: i/o timeout\n",
	"tls handshake":    "Error response from daemon: Get \"https://ghcr.io/v2/\": net/http: TLS handshake timeout\n",
	"no such host":     "Error response from daemon: Get \"https://ghcr.io/v2/\": dial tcp: lookup ghcr.io on 127.0.0.53:53: no such host\n",
	"registry 5xx":     "Error response from daemon: received unexpected HTTP status: 503 Service Unavailable\n",
}

// S-0740: a step that fails on the network is tried up to twice more,
// after 30 seconds and then 2 minutes, and reports how many times it ran.
func TestRunTriesANetworkFailureAgain(t *testing.T) {
	for kind, said := range networkOutputs {
		t.Run(kind, func(t *testing.T) {
			fail := try{out: said, err: errExit}
			docker := &fakeDocker{t: t, tries: []try{fail, fail, fail}}
			waits := &sleeps{}
			out := &syncBuffer{}
			r := Runner{Exec: docker.exec, Out: out, Sleep: waits.sleep}
			err := r.Run(context.Background(), Step{Name: "pull flai-agent-bedrock"})
			var failed *StepError
			if !errors.As(err, &failed) {
				t.Fatalf("Run = %v, want a *StepError", err)
			}
			if failed.Tries != 3 || !failed.Network || failed.Said != lastLine(said) || !errors.Is(err, errExit) {
				t.Errorf("StepError = %+v, want 3 tries of a network failure saying %q", failed, lastLine(said))
			}
			if !strings.Contains(err.Error(), "after 3 tries; check the connection to the registry") {
				t.Errorf("message %q does not say it ran out of tries and what to do", err.Error())
			}
			if want := []time.Duration{30 * time.Second, 2 * time.Minute}; docker.ran() != 3 || !reflect.DeepEqual(waits.waits, want) {
				t.Errorf("docker ran %d times, waiting %v, want 3 times waiting %v", docker.ran(), waits.waits, want)
			}
			for _, line := range []string{"trying again in 30s (try 2 of 3)", "trying again in 2m (try 3 of 3)"} {
				if !strings.Contains(out.String(), line) {
					t.Errorf("Out does not say %q:\n%s", line, out.String())
				}
			}
		})
	}
}

// S-0740: a step that fails on the network once and then passes is a pass.
func TestRunPassesOnTheSecondTry(t *testing.T) {
	for kind, said := range networkOutputs {
		t.Run(kind, func(t *testing.T) {
			docker := &fakeDocker{t: t, tries: []try{{out: said, err: errExit}, {out: "done\n"}}}
			waits := &sleeps{}
			r := Runner{Exec: docker.exec, Sleep: waits.sleep}
			if err := r.Run(context.Background(), Step{Name: "pull flai-agent-bedrock"}); err != nil {
				t.Fatalf("Run = %v, want a pass on the second try", err)
			}
			if want := []time.Duration{30 * time.Second}; docker.ran() != 2 || !reflect.DeepEqual(waits.waits, want) {
				t.Errorf("docker ran %d times, waiting %v, want twice waiting %v", docker.ran(), waits.waits, want)
			}
		})
	}
}

// S-0740: a step that fails for any other reason fails at once.
func TestRunFailsAtOnceOnAnythingElse(t *testing.T) {
	said := "ERROR: failed to solve: process \"/bin/sh -c apk add\" did not complete successfully: exit code: 1"
	docker := &fakeDocker{t: t, tries: []try{{out: "#7 [stage 1/3] RUN apk add\n" + said + "\n\n", err: errExit}}}
	waits := &sleeps{}
	r := Runner{Exec: docker.exec, Sleep: waits.sleep}
	err := r.Run(context.Background(), Step{Name: "build flai-agent-bedrock:1.2.3"})
	var failed *StepError
	if !errors.As(err, &failed) {
		t.Fatalf("Run = %v, want a *StepError", err)
	}
	if failed.Tries != 1 || failed.Network || failed.Said != said {
		t.Errorf("StepError = %+v, want one try, not a network failure, saying %q", failed, said)
	}
	if want := "build flai-agent-bedrock:1.2.3 failed (" + said + "); its output is above and in the log"; err.Error() != want {
		t.Errorf("message = %q, want %q", err.Error(), want)
	}
	if docker.ran() != 1 || len(waits.waits) != 0 {
		t.Errorf("docker ran %d times and waited %v, want once and no wait", docker.ran(), waits.waits)
	}
}

// S-0740: a step that writes nothing and fails says its error's first line.
func TestRunSaysTheErrorWhenTheStepWroteNothing(t *testing.T) {
	docker := &fakeDocker{t: t, tries: []try{{err: errors.New("exec: \"docker\": executable file not found in $PATH")}}}
	err := Runner{Exec: docker.exec}.Run(context.Background(), Step{Name: "pull bedrock"})
	var failed *StepError
	if !errors.As(err, &failed) || failed.Said != "exec: \"docker\": executable file not found in $PATH" {
		t.Errorf("Run = %v, want a *StepError saying the exec error", err)
	}
}

// S-0740: Network names a network failure only, never a number in a digest,
// a size, or a step count, nor a refused connection, which ADR-0295 does
// not list.
func TestNetwork(t *testing.T) {
	for _, output := range []string{
		"read: connection reset by peer",
		"dial tcp 1.2.3.4:443: i/o timeout",
		"net/http: TLS handshake timeout",
		"remote error: tls handshake failure",
		"context deadline exceeded",
		"operation timed out",
		"lookup ghcr.io: no such host",
		"Temporary failure in name resolution",
		"lookup proxy.golang.org on 10.0.0.1:53: server misbehaving",
		"Could not resolve host: dl-cdn.alpinelinux.org",
		"500 Internal Server Error",
		"502 Bad Gateway",
		"503 Service Unavailable",
		"504 Gateway Time-out",
		"unexpected status code 502",
		"status: 500",
		"HTTP/1.1 503",
		"HTTP 504",
	} {
		if !Network(output) {
			t.Errorf("Network(%q) = false, want true", output)
		}
	}
	for _, output := range []string{
		"",
		"sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08503",
		"ghcr.io/bytepunx/flai-agent-bedrock@sha256:5034ab12",
		"#5 503B / 2.1kB 0.1s",
		"#12 [stage 503/600] RUN make",
		"Downloading 503.2kB",
		"ERROR: failed to solve: process \"/bin/sh -c apk add\" did not complete successfully: exit code: 1",
		"dial tcp 127.0.0.1:2375: connect: connection refused",
		"http://localhost:5000/v2/",
	} {
		if Network(output) {
			t.Errorf("Network(%q) = true, want false", output)
		}
	}
}

// S-0740: a context that ends while Run waits to try again ends the wait,
// with the context's error.
func TestRunStopsWaitingWhenTheContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	docker := &fakeDocker{t: t, tries: []try{{out: networkOutputs["i/o timeout"], err: errExit, during: func(io.Writer) {
		time.AfterFunc(20*time.Millisecond, cancel)
	}}}}
	r := Runner{Exec: docker.exec, Delays: []time.Duration{time.Hour}}
	start := time.Now()
	err := r.Run(ctx, Step{Name: "pull bedrock"})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Run = %v, want the context's error", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second || docker.ran() != 1 {
		t.Errorf("Run took %s and ran docker %d times, want it to stop waiting at once after one try", elapsed, docker.ran())
	}
}

// S-0740: SystemExec runs a program in its folder, with the environment
// added, writing both its streams to out, and returns its exit's error.
func TestSystemExec(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh is not on PATH")
	}
	dir := t.TempDir()
	var out bytes.Buffer
	err := SystemExec(context.Background(), dir, []string{"FLAI_RUN_TEST=hi"}, &out, "sh", "-c", `echo "$FLAI_RUN_TEST"; pwd; echo oops >&2`)
	if err != nil {
		t.Fatalf("SystemExec: %v", err)
	}
	got := strings.Fields(out.String())
	resolved, _ := filepath.EvalSymlinks(dir)
	if len(got) != 3 || got[0] != "hi" || (got[1] != dir && got[1] != resolved) || got[2] != "oops" {
		t.Errorf("out = %q, want hi, %s, and oops", out.String(), dir)
	}
	var exit *exec.ExitError
	if err := SystemExec(context.Background(), dir, nil, io.Discard, "sh", "-c", "exit 3"); !errors.As(err, &exit) || exit.ExitCode() != 3 {
		t.Errorf("SystemExec of exit 3 = %v, want its exit status", err)
	}
}

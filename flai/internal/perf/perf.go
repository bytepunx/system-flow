// Package perf times a request inside flai, beneath the transport it came
// on (S-0152): the dashboard's channel and the MCP server each start a
// Recorder for a request, the code that answers it marks the phases it
// spends time in, and one event says where the time went. The dashboard's
// own duration for a request, less flai's, is the transport.
package perf

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// SlowEnv names the variable that sets how long a request may take before
// its event is logged at info rather than debug.
const SlowEnv = "FLAI_SLOW_REQUEST"

// DefaultSlow is the threshold when SlowEnv is unset.
const DefaultSlow = 500 * time.Millisecond

// Slow reads SlowEnv as a Go duration (250ms, 2s); unset, empty, or not a
// positive duration gives DefaultSlow.
func Slow() time.Duration {
	if d, err := time.ParseDuration(os.Getenv(SlowEnv)); err == nil && d > 0 {
		return d
	}
	return DefaultSlow
}

// Phase is the time a request spent in one named step, over every time it
// took that step.
type Phase struct {
	Name  string        `json:"name"`
	Count int           `json:"count"`
	Total time.Duration `json:"total_ns"`
}

// Recorder collects the phases of one request. It is safe for the
// goroutines a request starts.
type Recorder struct {
	start  time.Time
	now    func() time.Time
	mu     sync.Mutex
	phases map[string]*Phase
}

type key struct{}

// Start gives ctx a new Recorder, started now.
func Start(ctx context.Context) (context.Context, *Recorder) {
	return StartAt(ctx, time.Now)
}

// StartAt is Start with a clock, for tests.
func StartAt(ctx context.Context, now func() time.Time) (context.Context, *Recorder) {
	r := &Recorder{start: now(), now: now, phases: map[string]*Phase{}}
	return context.WithValue(ctx, key{}, r), r
}

// From is the Recorder in ctx, or nil when nothing is timing it.
func From(ctx context.Context) *Recorder {
	r, _ := ctx.Value(key{}).(*Recorder)
	return r
}

// Track starts a phase and returns what ends it:
//
//	defer perf.Track(ctx, "items.list")()
//
// Without a Recorder in ctx it costs a context lookup and does nothing.
func Track(ctx context.Context, name string) func() {
	r := From(ctx)
	if r == nil {
		return func() {}
	}
	began := r.now()
	return func() { r.Add(name, r.now().Sub(began)) }
}

// Add counts one step of a phase that took d.
func (r *Recorder) Add(name string, d time.Duration) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.phases[name]
	if p == nil {
		p = &Phase{Name: name}
		r.phases[name] = p
	}
	p.Count++
	p.Total += d
}

// Elapsed is the time since the request started.
func (r *Recorder) Elapsed() time.Duration {
	return r.now().Sub(r.start)
}

// Phases are the steps recorded so far, longest first.
func (r *Recorder) Phases() []Phase {
	r.mu.Lock()
	out := make([]Phase, 0, len(r.phases))
	for _, p := range r.phases {
		out = append(out, *p)
	}
	r.mu.Unlock()
	sort.Slice(out, func(i, j int) bool {
		if out[i].Total != out[j].Total {
			return out[i].Total > out[j].Total
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// String is the phases as name=milliseconds, with a count when a step was
// taken more than once: "items.list=182.4 exec.git=40.1x3". Phases overlap
// when a request runs steps together, and nest when one step takes others,
// so they need not add up to the whole.
func (r *Recorder) String() string {
	var b strings.Builder
	for i, p := range r.Phases() {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%s=%.1f", p.Name, ms(p.Total))
		if p.Count > 1 {
			fmt.Fprintf(&b, "x%d", p.Count)
		}
	}
	return b.String()
}

func ms(d time.Duration) float64 { return float64(d.Microseconds()) / 1000 }

// Answered is what is known about a request once it is answered.
type Answered struct {
	Transport string // channel or mcp
	Method    string // the method, or the MCP tool for tools/call
	Project   string // the project's key, when the transport names one
	Bytes     int    // the answer's size, when the transport knows it
	Err       string // the error answered, if any
	Wait      bool   // a request that waits by design: never logged as slow
}

// Log writes one "request answered" event: at info when the request took
// slow or longer and does not wait by design, at debug otherwise.
func (r *Recorder) Log(logger *slog.Logger, a Answered, slow time.Duration) {
	if logger == nil {
		return
	}
	took := r.Elapsed()
	level := slog.LevelDebug
	if !a.Wait && took >= slow {
		level = slog.LevelInfo
	}
	ctx := context.Background()
	if !logger.Enabled(ctx, level) {
		return
	}
	attrs := []slog.Attr{
		slog.String("component", "perf"),
		slog.String("transport", a.Transport),
		slog.String("method", a.Method),
		slog.Float64("duration_ms", ms(took)),
	}
	if a.Project != "" {
		attrs = append(attrs, slog.String("project", a.Project))
	}
	if a.Bytes > 0 {
		attrs = append(attrs, slog.Int("bytes", a.Bytes))
	}
	if phases := r.String(); phases != "" {
		attrs = append(attrs, slog.String("phases", phases))
	}
	if a.Err != "" {
		attrs = append(attrs, slog.String("err", a.Err))
	}
	logger.LogAttrs(ctx, level, "request answered", attrs...)
}

// Exec names the phase of running a program: "exec.git.log" for
// git -C dir log --oneline, "exec.flai.stream.diff" for flai stream diff
// S-0001 --json. Flags and the value of -C and -c are skipped, and at most
// two words are kept, so that the name says what ran and not on what.
func Exec(program string, args []string) string {
	name := "exec." + program
	words := 0
	for i := 0; i < len(args) && words < 2; i++ {
		a := args[i]
		if a == "-C" || a == "-c" {
			i++
			continue
		}
		if strings.HasPrefix(a, "-") || strings.ContainsAny(a, "/=. ") || (words > 0 && strings.ToLower(a) != a) {
			if words > 0 {
				break
			}
			continue
		}
		name += "." + a
		words++
	}
	return name
}

package execx

import (
	"context"

	"github.com/bytepunx/system-flow/flai/internal/perf"
)

// Timed is r with each command it runs recorded as a phase of the request
// ctx is timing (S-0152), named by perf.Exec. With nothing timing ctx it
// is r itself.
func Timed(ctx context.Context, r Runner) Runner {
	if r == nil || perf.From(ctx) == nil {
		return r
	}
	return timed{ctx: ctx, r: r}
}

type timed struct {
	ctx context.Context
	r   Runner
}

// Run runs the command as r does, timed.
func (t timed) Run(dir, name string, args ...string) (string, error) {
	defer perf.Track(t.ctx, perf.Exec(name, args))()
	return t.r.Run(dir, name, args...)
}

// RunInput runs the command with input as r does, timed.
func (t timed) RunInput(dir, name, input string, args ...string) (string, error) {
	defer perf.Track(t.ctx, perf.Exec(name, args))()
	return t.r.RunInput(dir, name, input, args...)
}

// LookPath is r's, untimed.
func (t timed) LookPath(name string) (string, error) { return t.r.LookPath(name) }

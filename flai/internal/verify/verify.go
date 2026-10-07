// Package verify runs a project's test and lint tiers for the paths a change
// touches and answers pass or the first failures as findings (S-0273). A
// tier is one command: the paths it is for, where it runs, and the format of
// its output. The tiers that select the paths run in list order, cheapest
// first, and the run stops at the first that fails; its output is parsed
// into a few findings, never handed back whole. flai test, the MCP tool
// test, and the host method test.run share it. Verify runs a story's
// close-out checks before the tiers its branch selects, as one report
// (S-0270).
package verify

import (
	"context"
	"io/fs"
	"os"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// DefaultMax is how many findings a run reports when Max is not set.
const DefaultMax = 5

// The formats a tier's output can be parsed as.
const (
	FormatGoTestJSON   = manifest.FormatGoTestJSON
	FormatVitestJSON   = manifest.FormatVitestJSON
	FormatGolangciJSON = manifest.FormatGolangciJSON
	FormatGofmtList    = manifest.FormatGofmtList
	FormatPlain        = manifest.FormatPlain
)

// Tier is one command a project runs to test or lint the paths it selects.
type Tier struct {
	// Name names the tier in the answer, such as go-test.
	Name string `json:"name"`
	// Command is the argv run; an argument that is wholly {packages} or
	// {files} stands for the Go packages or the files the paths select.
	Command []string `json:"command"`
	// Dir is where the command runs, relative to the checkout's root, ""
	// for the root itself. The placeholders are filled relative to it.
	Dir string `json:"dir,omitempty"`
	// Paths are globs relative to the root of the files the tier is for: **
	// is zero or more whole segments, * any characters within one, and a
	// pattern with no glob character is the path and everything below it.
	// A pattern beginning with ! excludes what it matches.
	Paths []string `json:"paths,omitempty"`
	// Format says how to read the command's output when it fails: one of
	// the Format constants, "" for plain.
	Format string `json:"format,omitempty"`
	// AllOnly keeps the tier out of every run but one of all the tiers.
	AllOnly bool `json:"all_only,omitempty"`
	// AllCommand is the argv run when every tier runs; without it Command
	// runs with {packages} as ./... and {files} as the directory itself.
	AllCommand []string `json:"all_command,omitempty"`
}

// Finding is one failure a tier's output names.
type Finding struct {
	// Name is the failing test's or check's name, or the linter's.
	Name string `json:"name,omitempty"`
	// Path is the file, relative to the checkout's root where the output
	// says enough to make it so.
	Path string `json:"path,omitempty"`
	// Line is the line in Path, 0 when the output gives none.
	Line int `json:"line,omitempty"`
	// Message is what the output says went wrong.
	Message string `json:"message,omitempty"`
}

// State is where a tier got to in a run.
type State string

// The states of a tier in a run.
const (
	Passed     State = "passed"
	Failed     State = "failed"
	NotReached State = "not-reached"
)

// TierResult is what one tier did in a run.
type TierResult struct {
	Name    string   `json:"name"`
	Command []string `json:"command"`
	Dir     string   `json:"dir,omitempty"`
	State   State    `json:"state"`
	// ExitCode is the command's exit status; nil when it did not run or
	// could not start.
	ExitCode   *int      `json:"exit_code,omitempty"`
	DurationMS int64     `json:"duration_ms"`
	Duration   string    `json:"duration,omitempty"`
	Findings   []Finding `json:"findings,omitempty"`
	// Omitted counts the findings past the run's cap.
	Omitted int `json:"omitted,omitempty"`
}

// Result is the answer to a run: pass, or the tier that failed and why.
type Result struct {
	Passed bool `json:"passed"`
	// Paths are the root-relative files the tiers were selected for; none
	// when every tier ran.
	Paths []string     `json:"paths,omitempty"`
	Tiers []TierResult `json:"tiers"`
}

// RunOptions are what a run needs beyond the tiers it runs.
type RunOptions struct {
	// Max caps the findings across the run; 0 or less is DefaultMax.
	Max int
	// FS is the checkout, read for go.mod; nil is the checkout on disk.
	FS fs.FS
	// Proc runs the commands; nil is OS.
	Proc Proc
	// Now is the clock that times each tier; nil is time.Now.
	Now func() time.Time
	// Env are KEY=value entries added to the environment each tier's
	// command inherits, as Verify adds CLOSE_OUT_STORY.
	Env []string
}

// Options are what Test needs to answer for a checkout.
type Options struct {
	// Root is the checkout's root, a worktree or the main checkout.
	Root string
	// Tiers are the project's tiers, cheapest first.
	Tiers []Tier
	// Args are files or folders relative to Root; a folder stands for the
	// files below it. With none, the run is for what the checkout changed.
	Args []string
	// All runs every tier for the whole checkout, whatever Args says.
	All bool
	// Base is the main branch the checkout's changes are counted against.
	Base string
	// Git runs git for the changed paths; nil is execx.System.
	Git execx.Runner
	RunOptions
}

// Test answers whether the checkout passes the tiers its paths select: the
// paths Args names, or, without them, the files the checkout changed against
// Base; with All, every tier. It returns an error when the paths cannot be
// read, or, with the result so far, when ctx ends during the run.
func Test(ctx context.Context, opts Options) (Result, error) {
	if opts.FS == nil {
		opts.FS = os.DirFS(opts.Root)
	}
	if opts.Git == nil {
		opts.Git = execx.System{}
	}
	var paths []string
	var err error
	switch {
	case opts.All:
	case len(opts.Args) > 0:
		paths, err = Resolve(opts.FS, opts.Args)
	default:
		paths, err = ChangedPaths(opts.Git, opts.Root, opts.Base)
	}
	if err != nil {
		return Result{}, err
	}
	res, err := Run(ctx, opts.Root, Select(opts.FS, opts.Tiers, paths, opts.All), opts.RunOptions)
	res.Paths = paths
	return res, err
}

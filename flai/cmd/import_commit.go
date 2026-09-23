package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/importer"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// exitImportNotCommitted is flai import --commit's exit code when the import
// was applied but not committed because a test failed (S-0098): the JSON
// answer on standard output says which, and the files are left as they are.
const exitImportNotCommitted = 5

// importTestOutputTail is how much of a test command's output is kept.
const importTestOutputTail = 8 << 10

// importTestResult is one test command's outcome.
type importTestResult struct {
	importer.TestCommand
	OK       bool   `json:"ok"`
	ExitCode int    `json:"exit_code"`
	Seconds  int    `json:"seconds"`
	Output   string `json:"output" jsonschema:"the last 8 KiB of what it printed"`
	Error    string `json:"error,omitempty"`
}

// importCommit is what --commit did after the import.
type importCommit struct {
	Tests     []importTestResult `json:"tests"`
	Source    string             `json:"tests_from"` // host, repository, or none
	Committed bool               `json:"committed"`
	Commit    string             `json:"commit,omitempty"`
	Reason    string             `json:"reason,omitempty"` // why it was not committed
}

// importPreflight refuses --commit where the import's commit could not hold
// the import alone: not a git repository, or one with uncommitted changes.
func (a *app) importPreflight(an *importer.Analysis) error {
	if !an.Git {
		return fmt.Errorf("%s is not a git repository; --commit needs one", an.Root)
	}
	out, err := a.runner.Run(an.Root, "git", "status", "--porcelain")
	if err != nil {
		return fmt.Errorf("git status in %s: %w", an.Root, err)
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("%s has uncommitted changes; commit or stash them first, so that the import's commit holds the import and nothing else", an.Root)
	}
	return nil
}

// importTests are the commands --commit runs: the host's checks when it names
// any (flai serve checks set), else what the repository itself has.
func (a *app) importTests(root string, ownMakefile bool, projects []importer.Project) ([]importer.TestCommand, string) {
	cfg, _, err := a.loadConfig()
	if err == nil && len(cfg.Checks.Commands) > 0 {
		var out []importer.TestCommand
		for _, c := range cfg.Checks.Commands {
			out = append(out, importer.TestCommand{Name: c.Name, Dir: ".", Command: serve.Substitute(c.Command, "", root)})
		}
		return out, "host"
	}
	found := importer.DetectTests(root, ownMakefile, projects)
	if len(found) == 0 {
		return nil, "none"
	}
	return found, "repository"
}

// runImportTests runs each command in turn, all of them however the first
// ones went, so the operator sees every failure at once.
func (a *app) runImportTests(root string, cmds []importer.TestCommand) []importTestResult {
	timeout := serve.DefaultChecksTimeout
	if cfg, _, err := a.loadConfig(); err == nil && cfg.Checks.TimeoutMinutes > 0 {
		timeout = time.Duration(cfg.Checks.TimeoutMinutes) * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log := a.logger().With("component", "import")
	out := make([]importTestResult, 0, len(cmds))
	for _, c := range cmds {
		log.Info("running tests", "name", c.Name, "dir", c.Dir)
		started := time.Now()
		cmd := exec.CommandContext(ctx, c.Command[0], c.Command[1:]...)
		cmd.Dir = filepath.Join(root, c.Dir)
		var buf bytes.Buffer
		cmd.Stdout, cmd.Stderr = &buf, &buf
		err := cmd.Run()
		r := importTestResult{TestCommand: c, OK: err == nil, Seconds: int(time.Since(started).Seconds()), Output: tail(buf.String(), importTestOutputTail)}
		var exit *exec.ExitError
		switch {
		case errors.As(err, &exit):
			r.ExitCode = exit.ExitCode()
		case err != nil:
			r.ExitCode, r.Error = -1, err.Error()
		}
		if ctx.Err() != nil {
			r.OK, r.Error = false, fmt.Sprintf("stopped after %s, the time the host allows a run of checks", timeout)
		}
		if r.OK {
			log.Info("tests passed", "name", c.Name, "seconds", r.Seconds)
		} else {
			log.Warn("tests failed", "name", c.Name, "exit", r.ExitCode, "seconds", r.Seconds)
		}
		out = append(out, r)
	}
	return out
}

func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return "…" + s[len(s)-n:]
}

// commitImport runs the tests and, when every one passed or there were none,
// commits exactly the paths the import wrote or moved.
func (a *app) commitImport(root string, ownMakefile bool, projects []importer.Project, written []string, moved []importer.Move, trailers []string) (importCommit, error) {
	cmds, from := a.importTests(root, ownMakefile, projects)
	res := importCommit{Source: from, Tests: a.runImportTests(root, cmds)}
	var failed []string
	for _, t := range res.Tests {
		if !t.OK {
			failed = append(failed, t.Name)
		}
	}
	if len(failed) > 0 {
		res.Reason = fmt.Sprintf("tests failed: %s; the imported files are left uncommitted", strings.Join(failed, ", "))
		return res, nil
	}
	// what the import wrote, and where it moved things to: a move's source is
	// already staged as removed by git mv, and naming a path that is gone
	// would make git add refuse the whole list
	paths := map[string]bool{manifest.File: true}
	for _, p := range written {
		paths[p] = true
	}
	for _, m := range moved {
		paths[m.To] = true
	}
	list := make([]string, 0, len(paths))
	for p := range paths {
		if _, err := os.Stat(filepath.Join(root, p)); err == nil {
			list = append(list, p)
		}
	}
	sort.Strings(list)
	if _, err := a.runner.Run(root, "git", append([]string{"add", "-A", "--"}, list...)...); err != nil {
		return res, fmt.Errorf("git add: %w", err)
	}
	msg := fmt.Sprintf("chore: bring %s under system-flow (flai import)", filepath.Base(root))
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	if _, err := a.runner.Run(root, "git", "commit", "-q", "-m", msg); err != nil {
		return res, fmt.Errorf("git commit: %w", err)
	}
	sha, _ := a.runner.Run(root, "git", "rev-parse", "--short", "HEAD")
	res.Committed, res.Commit = true, strings.TrimSpace(sha)
	return res, nil
}

// ownMakefileAt says whether root has a Makefile of its own, before the
// template brings one.
func ownMakefileAt(root string) bool {
	_, err := os.Stat(filepath.Join(root, "Makefile"))
	return err == nil
}

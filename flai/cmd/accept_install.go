package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// installTailBytes is how much of an install command's latest output a run
// keeps, from which its warning takes the last installTailLines lines.
const (
	installTailBytes = 8 << 10
	installTailLines = 10
)

// installRun is one sub-project's install command that flai accept ran in
// the main checkout after the acceptance commit (ADR-0185).
type installRun struct {
	// Project is the sub-project's name in the manifest.
	Project string `json:"project"`
	// Command is the argument list run, with no shell.
	Command []string `json:"command"`
	// OK says the command started and exited 0.
	OK bool `json:"ok"`
	// Error is why it did not: the exit status, or why it could not start.
	Error string `json:"error,omitempty"`
	// folder is the sub-project's folder in the main checkout, where it ran.
	folder string
	// exit is the command's exit status when it ran and failed, for the
	// warning; -1 when it could not start.
	exit int
	// output is the last lines the command printed, for the warning.
	output string
}

// failedLine is the text output's line for a run that failed: the reason
// and the command to run by hand, and where.
func (r installRun) failedLine() string {
	return fmt.Sprintf("install for %s failed: %s; run %s in %s", r.Project, r.Error, commandLine(r.Command), r.folder)
}

// installAccepted runs, in the main checkout, the install command of each
// sub-project whose install paths match changed, the repository-relative
// paths the acceptance of it brought into the main branch (ADR-0185). It
// reads the manifest the main checkout holds after the commit, so a
// sub-project's install declared by the accepted change itself counts. A run
// that fails is warned of and recorded; none fails the acceptance, which is
// committed by then.
func (a *app) installAccepted(repo *workitem.Repo, it *workitem.Item, changed []string) []installRun {
	if len(changed) == 0 {
		return nil
	}
	m, err := manifest.Load(filepath.Join(repo.MainRoot, manifest.File))
	if err != nil {
		a.logger().Warn("sub-project installs not run after acceptance", "component", "accept", "item", it.ID, "err", err)
		return nil
	}
	runs := runInstalls(context.Background(), verify.OS{}, m.Projects, repo.MainRoot, changed, a.errOut)
	for _, r := range runs {
		if r.OK {
			continue
		}
		attrs := []any{"component", "accept", "item", it.ID, "project", r.Project, "command", commandLine(r.Command), "dir", r.folder}
		if r.exit >= 0 {
			attrs = append(attrs, "exit", r.exit, "output", r.output)
		}
		a.logger().Warn("sub-project install failed after acceptance", append(attrs, "err", r.Error)...)
	}
	return runs
}

// runInstalls runs with proc, in projects' order, the install command of
// each project whose install paths match changed, once, in its folder under
// root, with what it prints written to out. It answers each run, with the
// last lines of the output of one that failed.
func runInstalls(ctx context.Context, proc verify.Proc, projects []manifest.Project, root string, changed []string, out io.Writer) []installRun {
	var runs []installRun
	for _, p := range projects {
		if !p.InstallMatches(changed) {
			continue
		}
		r := installRun{Project: p.Name, Command: slices.Clone(p.Install.Command), folder: filepath.Join(root, filepath.FromSlash(p.Path)), exit: -1}
		if len(r.Command) == 0 {
			r.Error = "the manifest's install command is empty, which flai check reports as manifest.install"
			runs = append(runs, r)
			continue
		}
		tail := &outputTail{}
		w := io.MultiWriter(out, tail)
		code, err := proc.Run(ctx, r.folder, r.Command, nil, w, w)
		switch {
		case err != nil:
			r.Error = "cannot start: " + err.Error()
		case code != 0:
			r.exit, r.Error, r.output = code, fmt.Sprintf("exit status %d", code), tail.lastLines(installTailLines)
		default:
			r.OK, r.exit = true, 0
		}
		runs = append(runs, r)
	}
	return runs
}

// outputTail keeps the last installTailBytes bytes written to it.
type outputTail struct {
	b []byte
}

// Write keeps p's end and never fails.
func (t *outputTail) Write(p []byte) (int, error) {
	t.b = append(t.b, p...)
	if over := len(t.b) - installTailBytes; over > 0 {
		t.b = t.b[over:]
	}
	return len(p), nil
}

// lastLines is the last n lines kept that are not blank, one per line.
func (t *outputTail) lastLines(n int) string {
	var lines []string
	for _, l := range strings.Split(string(bytes.TrimRight(t.b, "\n")), "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, strings.TrimRight(l, "\r"))
		}
	}
	return strings.Join(lines[max(0, len(lines)-n):], "\n")
}

// commandLine is argv as a POSIX shell reads it back: each word as it is
// when it holds only letters, digits, and ./_=:-, else single quoted.
func commandLine(argv []string) string {
	words := make([]string, len(argv))
	for i, w := range argv {
		if w != "" && strings.Trim(w, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789./_=:-") == "" {
			words[i] = w
		} else {
			words[i] = "'" + strings.ReplaceAll(w, "'", `'\''`) + "'"
		}
	}
	return strings.Join(words, " ")
}

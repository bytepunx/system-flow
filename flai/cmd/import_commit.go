package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/importer"
	"github.com/bytepunx/system-flow/flai/internal/lock"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// exitImportNotCommitted is flai import --commit's exit code when the import
// was applied but not committed because a test failed (S-0098): the JSON
// answer on standard output says which, and the files are left as they are.
const exitImportNotCommitted = 5

// exitImportTestsFailed is flai import test's exit code when a test failed
// (S-0424), the code --commit gives for the same.
const exitImportTestsFailed = exitImportNotCommitted

// Where the tests an import runs come from: the host's checks, what the
// repository has, none found, or, for --formula, the import formula's
// instance, which runs them after the commit (S-0424).
const (
	importTestsHost       = "host"
	importTestsRepository = "repository"
	importTestsNone       = "none"
	importTestsFormula    = "formula"
)

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
	Source    string             `json:"tests_from"` // host, repository, none, or formula
	Committed bool               `json:"committed"`
	Commit    string             `json:"commit,omitempty"`
	Reason    string             `json:"reason,omitempty"` // why it was not committed
}

// importPreflight refuses --commit or --formula, flag, where the import's
// commit could not hold the import alone: not a git repository, or one with
// uncommitted changes.
func (a *app) importPreflight(an *importer.Analysis, flag string) error {
	if !an.Git {
		return fmt.Errorf("%s is not a git repository; %s needs one", an.Root, flag)
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
		return out, importTestsHost
	}
	found := importer.DetectTests(root, ownMakefile, projects)
	if len(found) == 0 {
		return nil, importTestsNone
	}
	return found, importTestsRepository
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
	sha, err := a.commitImportPaths(root, written, moved, trailers)
	if err != nil {
		return res, err
	}
	res.Committed, res.Commit = true, sha
	return res, nil
}

// commitImportUntested commits exactly the paths the import wrote or moved,
// running no test first: under --formula the import formula's instance runs
// them after the commit (S-0424).
func (a *app) commitImportUntested(root string, written []string, moved []importer.Move, trailers []string) (importCommit, error) {
	res := importCommit{Source: importTestsFormula, Tests: []importTestResult{}}
	sha, err := a.commitImportPaths(root, written, moved, trailers)
	if err != nil {
		return res, err
	}
	res.Committed, res.Commit = true, sha
	return res, nil
}

// commitImportPaths commits exactly the paths the import wrote or moved, the
// manifest, and the lock, with the trailers, and returns the commit's short
// hash.
func (a *app) commitImportPaths(root string, written []string, moved []importer.Move, trailers []string) (string, error) {
	// what the import wrote, and where it moved things to: a move's source is
	// already staged as removed by git mv, and naming a path that is gone
	// would make git add refuse the whole list
	paths := map[string]bool{manifest.File: true, lock.File: true}
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
		return "", fmt.Errorf("git add: %w", err)
	}
	msg := fmt.Sprintf("chore: bring %s under system-flow %s", filepath.Base(root), importCommitMark)
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	if _, err := a.runner.Run(root, "git", "commit", "-q", "-m", msg); err != nil {
		return "", fmt.Errorf("git commit: %w", err)
	}
	sha, _ := a.runner.Run(root, "git", "rev-parse", "--short", "HEAD")
	return strings.TrimSpace(sha), nil
}

// importCommitMark ends the subject of the import's commit, so that flai
// import test knows a Makefile that commit added for the template's.
const importCommitMark = "(flai import)"

// importFormula is the formula flai import --formula runs in the new
// project (S-0424): its steps check the import and run its tests.
const importFormula = "import"

// importFormulaShipped refuses --formula, before anything changes, when the
// template m at src renders no import formula into the design folder.
func importFormulaShipped(m template.Manifest, src template.Source) error {
	design := m.Layout["design"]
	if design == "" {
		design = "design"
	}
	rel := path.Join(design, "workflows", "formulas", importFormula+".md")
	for _, name := range []string{rel, rel + m.Render.Suffix} {
		if _, err := os.Stat(filepath.Join(src.Dir, template.RootDir, filepath.FromSlash(name))); err == nil {
			return nil
		}
	}
	version := m.Version
	if version == "" {
		version = "with no version"
	}
	return fmt.Errorf("the template %s, version %s, ships no %s formula (%s), which --formula runs in the new project, so nothing was changed: use a template version that ships it (--ref), or import with --commit", src.Repo, version, importFormula, rel)
}

// runImportFormula makes the import formula's instance in the project at
// root, as flai formula run import --input path=<root> --ready --autocommit
// does there, with the trailers on its commit, and returns it.
func (a *app) runImportFormula(root string, trailers []string) (formulaRunResult, error) {
	var buf bytes.Buffer
	there := *a
	there.cwd, there.out, there.jsonOut = root, &buf, true
	var run formulaRunResult
	if err := there.formulaRun(importFormula, map[string]string{"path": root}, "", formulaRunOptions{ready: true, autocommit: true, trailers: trailers}); err != nil {
		return run, err
	}
	if err := json.Unmarshal(buf.Bytes(), &run); err != nil {
		return run, fmt.Errorf("reading what flai formula run made: %w", err)
	}
	return run, nil
}

// importTestRun is what flai import test ran and how each went.
type importTestRun struct {
	Root   string             `json:"root"`
	Source string             `json:"tests_from"` // host, repository, or none
	Tests  []importTestResult `json:"tests"`
	OK     bool               `json:"ok"`
}

// runImportTest runs the tests --commit runs in the project at or above dir,
// from its root, and fails with exitImportTestsFailed when one failed
// (S-0424).
func (a *app) runImportTest(dir string) error {
	there := *a
	if !filepath.IsAbs(dir) {
		wd, err := a.workingDir()
		if err != nil {
			return err
		}
		dir = filepath.Join(wd, dir)
	}
	there.cwd = dir
	repo, err := there.project()
	if err != nil {
		return err
	}
	projects := make([]importer.Project, 0, len(repo.Manifest.Projects))
	for _, p := range repo.Manifest.Projects {
		projects = append(projects, importer.Project{Name: p.Name, Path: p.Path, Kind: p.Kind})
	}
	cmds, from := a.importTests(repo.Root, a.importOwnMakefile(repo), projects)
	run := importTestRun{Root: repo.Root, Source: from, Tests: a.runImportTests(repo.Root, cmds), OK: true}
	var failed []string
	for _, t := range run.Tests {
		if !t.OK {
			failed = append(failed, t.Name)
		}
	}
	run.OK = len(failed) == 0
	if a.jsonOut {
		if err := a.printJSON(run); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(a.out, "Tests of %s\n", repo.Root)
		a.printImportTests(run.Source, run.Tests)
		if run.OK {
			fmt.Fprintf(a.out, "  %s passed\n", plural(len(run.Tests), "test command"))
		}
	}
	if !run.OK {
		return &exitError{code: exitImportTestsFailed, msg: "tests failed: " + strings.Join(failed, ", ")}
	}
	return nil
}

// importOwnMakefile says whether the project's root Makefile is the
// repository's own rather than the one the import brought: the template's
// byte for byte, or one the import's commit added.
func (a *app) importOwnMakefile(repo *workitem.Repo) bool {
	data, err := os.ReadFile(filepath.Join(repo.Root, "Makefile"))
	if err != nil {
		return false
	}
	if repo.TemplateDir != "" {
		if tpl, err := os.ReadFile(filepath.Join(repo.TemplateDir, template.RootDir, "Makefile")); err == nil && bytes.Equal(tpl, data) {
			return false
		}
	}
	// the newest commit that added the Makefile there is
	if out, err := a.runner.Run(repo.Root, "git", "log", "--diff-filter=A", "--format=%s", "--", "Makefile"); err == nil {
		added, _, _ := strings.Cut(strings.TrimSpace(out), "\n")
		if strings.HasSuffix(strings.TrimSpace(added), importCommitMark) {
			return false
		}
	}
	return true
}

// ownMakefileAt says whether root has a Makefile of its own, before the
// template brings one.
func ownMakefileAt(root string) bool {
	_, err := os.Stat(filepath.Join(root, "Makefile"))
	return err == nil
}

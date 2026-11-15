package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// workingDir is the absolute directory flai was run in.
func (a *app) workingDir() (string, error) {
	start := a.cwd
	if start == "" {
		var err error
		if start, err = os.Getwd(); err != nil {
			return "", err
		}
	}
	return filepath.Abs(start)
}

// projectOrNone is the project above the working directory, or nil when there
// is no system-flow.yaml there or in any parent (S-0101): the commands that
// act for the whole host, not one project, work from any folder. Any other
// failure, a manifest that does not parse say, is still an error.
func (a *app) projectOrNone() (*workitem.Repo, error) {
	repo, err := a.project()
	if errors.Is(err, manifest.ErrNotFound) {
		return nil, nil
	}
	return repo, err
}

// testWorkflow, when a test sets it, is the workflow every project the
// commands open runs in place of the baseline (ADR-0301), so that the
// commands are tried against a board with renamed and added states.
var testWorkflow *workflow.Workflow

// mcpProject is projectOrNone for the MCP server: a workflow file flai
// cannot run does not keep it from starting; each tool call is refused with
// the file's finding instead, until the file is fixed (ADR-0307 §3).
func (a *app) mcpProject() (*workitem.Repo, error) {
	repo, err := a.projectUnchecked()
	if errors.Is(err, manifest.ErrNotFound) {
		return nil, nil
	}
	return repo, err
}

// project opens the conforming repo above the working directory and points
// it at the cached template for item bodies when one is available. It never
// fetches. A workflow file flai cannot run refuses it (ADR-0307 §2).
func (a *app) project() (*workitem.Repo, error) {
	return a.openProject(workitem.Open)
}

// projectUnchecked is project for flai check alone: a workflow file flai
// cannot run is kept as the repo's FlowErr, and the repo runs the baseline
// (ADR-0307 §2).
func (a *app) projectUnchecked() (*workitem.Repo, error) {
	return a.openProject(workitem.OpenUnchecked)
}

// openProject is project, opening the repo with open.
func (a *app) openProject(open func(string) (*workitem.Repo, error)) (*workitem.Repo, error) {
	start := a.cwd
	if start == "" {
		var err error
		if start, err = os.Getwd(); err != nil {
			return nil, err
		}
	}
	repo, err := open(start)
	var refused *workflow.LoadError
	switch {
	case errors.As(err, &refused):
		return nil, err
	case err != nil:
		return nil, fmt.Errorf("%w (run flai new or flai import first)", err)
	}
	repo.Git = a.runner
	if testWorkflow != nil {
		repo.Workflow = testWorkflow
	}
	// an overlap holds a ready story only when a trial merge conflicts, as
	// flai story start, flai board, and the promotion candidates judge it
	// (ADR-0230); outside a git checkout NewTrialJudge gives none
	if a.runner != nil {
		repo.Trial = storygit.NewTrialJudge(a.runner, repo, issues.FilesIn(repo, repo.MainRoot, a.runner, time.Now))
	}
	cacheDir := config.Default().CacheDir
	if cfg, _, err := a.loadConfig(); err == nil {
		cacheDir = cfg.CacheDir
	}
	if t := repo.Manifest.Template; t.Repo != "" {
		if src, err := template.Resolve(t.Repo, t.Ref, cacheDir); err == nil && src.Cached() {
			repo.TemplateDir = src.Dir
		}
	}
	a.logger().Debug("project opened", "component", "workitem", "root", repo.Root, "template_dir", repo.TemplateDir)
	return repo, nil
}

// now is the clock, overridable in tests.
func (a *app) now() time.Time {
	if a.clock != nil {
		return a.clock()
	}
	return time.Now().UTC()
}

// author is the --by default: config author, else "agent".
func (a *app) author() string {
	if cfg, _, err := a.loadConfig(); err == nil && cfg.Author != "" {
		return cfg.Author
	}
	return "agent"
}

// agentIdentity reads FLAI_AGENT and FLAI_SESSION for narratives.
func agentIdentity() (agent, session string) {
	agent = os.Getenv("FLAI_AGENT")
	if agent == "" {
		agent = "agent"
	}
	return agent, os.Getenv("FLAI_SESSION")
}

// refreshIndex regenerates wip/agents/index.md after anything that changes
// item status or narratives.
func (a *app) refreshIndex(repo *workitem.Repo) error {
	items, err := repo.List(false)
	if err != nil {
		return err
	}
	return repo.WriteIndex(items, a.now())
}

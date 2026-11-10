package check

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func (c *checker) layout() {
	m := c.repo.Manifest
	// API responses name the project by this key (ADR-0024). A warning now;
	// an error after the template minor that ships the requirement.
	if strings.TrimSpace(m.Key) == "" {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Warning, "manifest.key", mf, keyLine(mf, "name"), "system-flow.yaml has no key; add a short one (for example the project's initials), which the dashboard's API and a hub use to name this project")
	}
	if err := m.Agent.Validate(); err != nil {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.agent", mf, keyLine(mf, "agent"), "%s; flai agent set replaces it", err)
	}
	if _, err := m.Issues.StoryAfterDuration(); err != nil {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.issues", mf, keyLine(mf, "issues"), "%v", err)
	}
	for _, msg := range m.Planning.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.planning", mf, keyLine(mf, "planning"), "%s", msg)
	}
	for _, msg := range m.Orchestration.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.orchestration", mf, keyLine(mf, "orchestration"), "%s", msg)
	}
	for _, msg := range m.RoleErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.roles", mf, keyLine(mf, "roles"), "%s", msg)
	}
	for _, msg := range m.Analysis.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.analysis", mf, keyLine(mf, "analysis"), "%s", msg)
	}
	for _, msg := range m.Digest.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.digest", mf, keyLine(mf, "digest"), "%s", msg)
	}
	for _, msg := range m.Threads.Errors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.threads", mf, keyLine(mf, "threads"), "%s", msg)
	}
	for _, msg := range m.TestErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.tests", mf, keyLine(mf, "tests"), "%s", msg)
	}
	for _, msg := range m.TestWarnings() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Warning, "manifest.tests-timeout", mf, keyLine(mf, "tests"), "%s", msg)
	}
	for _, msg := range m.ProjectErrors(os.DirFS(c.repo.Root)) {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.projects", mf, keyLine(mf, "projects"), "%s", msg)
	}
	for _, msg := range m.ComponentErrors(os.DirFS(c.repo.Root)) {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.components", mf, keyLine(mf, "projects"), "%s", msg)
	}
	for _, msg := range m.InstallErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.install", mf, keyLine(mf, "projects"), "%s", msg)
	}
	for _, msg := range m.ProviderErrors() {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.providers", mf, keyLine(mf, "providers"), "%s", msg)
	}
	// the kluster: block (S-0577, ADR-0254): each finding names its key
	for _, msg := range m.KlusterErrors(os.DirFS(c.repo.Root)) {
		mf := filepath.Join(c.repo.Root, "system-flow.yaml")
		c.add(Error, "manifest.kluster", mf, keyLine(mf, "kluster"), "%s", msg)
	}
	c.sharedPatterns()
	if m.Template.Version != "" {
		if _, err := os.Stat(filepath.Join(c.repo.Root, "system-flow.lock.yaml")); err != nil {
			c.add(Warning, "layout.lock", filepath.Join(c.repo.Root, "system-flow.yaml"), keyLine(filepath.Join(c.repo.Root, "system-flow.yaml"), "template"), "template %s is recorded but there is no system-flow.lock.yaml; run flai upgrade --relock", m.Template.Version)
		}
	}
	for key, subs := range map[string][]string{
		"design": {"adrs", "system", "tech", "conventions"},
		"docs":   nil,
		"wip":    {"kanban", "agents", "archive"},
	} {
		dir := m.Dir(c.repo.Root, key)
		if st, err := os.Stat(dir); err != nil || !st.IsDir() {
			c.add(Error, "layout.folder", c.repo.Root, 1, "layout.%s folder %q does not exist", key, m.Layout[key])
			continue
		}
		for _, sub := range subs {
			if st, err := os.Stat(filepath.Join(dir, sub)); err != nil || !st.IsDir() {
				c.add(Warning, "layout.subfolder", dir, 1, "%s/%s is missing", m.Layout[key], sub)
			}
		}
	}
	for _, sub := range []string{"epics", "stories", "tasks"} {
		if st, err := os.Stat(filepath.Join(c.repo.KanbanDir(), sub)); err != nil || !st.IsDir() {
			c.add(Warning, "layout.subfolder", c.repo.KanbanDir(), 1, "kanban/%s is missing", sub)
		}
	}
	c.cacheIgnored()
}

// sharedPatterns reports each pattern of claims.shared that is not valid, on
// its line of the manifest: it frees nothing (ADR-0096).
func (c *checker) sharedPatterns() {
	errs := c.repo.SharedClaims().Errors() // as the file has them, like the lines
	if len(errs) == 0 {
		return
	}
	mf := filepath.Join(c.repo.Root, manifest.File)
	data, _ := os.ReadFile(mf) // without it, each is reported on line 1
	lines := manifest.SharedLines(data)
	for _, e := range errs {
		line := 1
		if e.Index < len(lines) {
			line = lines[e.Index]
		}
		c.add(Error, "manifest.claims", mf, line, "%s; until it is fixed it frees nothing", e)
	}
}

// cacheIgnored refuses a repository whose .flai-cache (dashboard token,
// worktrees, tool caches) could be committed (ADR-0018).
func (c *checker) cacheIgnored() {
	cache := filepath.Join(c.repo.Root, ".flai-cache")
	if st, err := os.Stat(cache); err != nil || !st.IsDir() {
		return
	}
	gi := filepath.Join(c.repo.Root, ".gitignore")
	data, err := os.ReadFile(gi)
	if err != nil {
		c.add(Error, "layout.gitignore", c.repo.Root, 1, ".flai-cache/ exists but there is no .gitignore; it holds the dashboard token and must be ignored")
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		switch strings.TrimSpace(line) {
		case ".flai-cache", ".flai-cache/", "/.flai-cache", "/.flai-cache/":
			return
		}
	}
	c.add(Error, "layout.gitignore", gi, 1, ".flai-cache/ is not ignored; it holds the dashboard token and must be")
}

var requiredHeadings = map[string][]string{
	workitem.Epic:  {"## Outcome", "## Stories", "## Notes"},
	workitem.Story: {"## Goal", "## Acceptance criteria", "## Tasks", "## Notes"},
	workitem.Task:  {"## Work", "## Done when", "## Notes"},
}

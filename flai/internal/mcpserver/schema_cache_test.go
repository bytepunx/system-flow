package mcpserver

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
)

// toolSchemas lists the tools a server made with opt answers tools/list
// with, as each tool's name and its input and output schemas in JSON.
func toolSchemas(t *testing.T, opt Options) map[string]string {
	t.Helper()
	cs := connectServer(t, opt)
	out := map[string]string{}
	for tool, err := range cs.Tools(context.Background(), nil) {
		if err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(struct{ In, Out any }{tool.InputSchema, tool.OutputSchema})
		if err != nil {
			t.Fatal(err)
		}
		out[tool.Name] = string(data)
	}
	if len(out) == 0 {
		t.Fatal("tools/list answered no tools")
	}
	return out
}

// sameSchemas fails t naming each tool whose schemas differ between got and
// want, or that only one of them lists.
func sameSchemas(t *testing.T, what string, got, want map[string]string) {
	t.Helper()
	names := map[string]bool{}
	for n := range got {
		names[n] = true
	}
	for n := range want {
		names[n] = true
	}
	sorted := make([]string, 0, len(names))
	for n := range names {
		sorted = append(sorted, n)
	}
	slices.Sort(sorted)
	for _, n := range sorted {
		g, gok := got[n]
		w, wok := want[n]
		switch {
		case !gok:
			t.Errorf("%s: %s is not listed, though the uncached server lists it", what, n)
		case !wok:
			t.Errorf("%s: %s is listed, though the uncached server does not list it", what, n)
		case g != w:
			t.Errorf("%s: %s's schemas are\n%s\nwant\n%s", what, n, g, w)
		}
	}
}

// uncached is opt with no schema cache.
func uncached(opt Options) Options {
	opt.noSchemaCache = true
	return opt
}

// Two servers made from the same options share the schema cache and list
// the same tools with the same schemas, which are those a server made
// without the cache lists (S-0721).
func TestServersSharingTheSchemaCacheListTheSameSchemas(t *testing.T) {
	t.Parallel()
	opt := Options{Repo: makeProject(t, t.TempDir(), "p"), Agent: "claude", Version: "test"}
	first, second := toolSchemas(t, opt), toolSchemas(t, opt)
	sameSchemas(t, "second cached server", second, first)
	sameSchemas(t, "cached server", first, toolSchemas(t, uncached(opt)))
}

// Servers made from differing options, a single project and a folder and
// each role, list with the shared cache the schemas each lists without it,
// though the cache was filled by the other option sets first (S-0721).
func TestTheSchemaCacheKeepsEachOptionSetsSchemas(t *testing.T) {
	t.Parallel()
	folder := t.TempDir()
	makeProject(t, filepath.Join(folder, "a"), "a")
	makeProject(t, filepath.Join(folder, "b"), "b")
	repo := makeProject(t, t.TempDir(), "p")
	single := func(with func(*Options)) Options {
		opt := Options{Repo: repo, Agent: "claude", Version: "test"}
		with(&opt)
		return opt
	}
	sets := []struct {
		name string
		opt  Options
	}{
		{"folder", Options{Folder: folder, Agent: "claude", Version: "test"}},
		{"one project", single(func(*Options) {})},
		{"story's agent", single(func(o *Options) { o.Story, o.StartedBy = "S-0001", guard.StartedByServe })},
		{"planner", single(func(o *Options) { o.Role, o.StartedBy = guard.RolePlan, guard.StartedByServe })},
		{"orchestrator", single(func(o *Options) {
			o.Role, o.Reasons, o.StartedBy = guard.RoleOrchestrate, "ready", guard.StartedByServe
		})},
		{"analyzer", single(func(o *Options) { o.Role, o.StartedBy = guard.RoleAnalyze, guard.StartedByServe })},
		{"curator", single(func(o *Options) { o.Role, o.StartedBy = guard.RoleCurate, guard.StartedByServe })},
	}
	cached := make([]map[string]string, len(sets))
	for i, s := range sets {
		cached[i] = toolSchemas(t, s.opt)
	}
	for i, s := range sets {
		sameSchemas(t, s.name, cached[i], toolSchemas(t, uncached(s.opt)))
	}
}

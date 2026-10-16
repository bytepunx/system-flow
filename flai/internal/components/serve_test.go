package components

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// pickProjects are a SvelteKit app, a Go sub-project, and a Svelte package.
var pickProjects = []manifest.Project{
	{Name: "ui", Path: "ui", Kind: manifest.ProjectKindSvelteKit},
	{Name: "cli", Path: "cli", Kind: "go"},
	{Name: "kit", Path: "packages/kit", Kind: manifest.ProjectKindSvelte},
}

func TestPickSubProject(t *testing.T) {
	for _, tc := range []struct {
		name     string
		projects []manifest.Project
		pick     string
		want     string
		err      string
	}{
		{"named", pickProjects, "kit", "kit", ""},
		{"the one supported", pickProjects[:2], "", "ui", ""},
		{"unknown", pickProjects, "web", "", `system-flow.yaml has no sub-project "web"; name one of ui, cli, kit`},
		{"unknown among none", nil, "web", "", `has no sub-project "web", and names none under projects`},
		{"unsupported", pickProjects, "cli", "", `sub-project "cli" is of kind go, and flai serves the components of a sub-project of kind sveltekit or svelte only; name one of ui, kit`},
		{"unsupported alone", pickProjects[1:2], "cli", "", `is of kind go, and flai serves the components of a sub-project of kind sveltekit or svelte only; system-flow.yaml names none`},
		{"none supported", pickProjects[1:2], "", "", "names no sub-project of kind sveltekit or svelte"},
		{"ambiguous", pickProjects, "", "", "has 2 sub-projects of kind sveltekit or svelte (ui, kit), and flai serves one at a time; name the one to serve, such as --project ui"},
	} {
		got, err := PickSubProject(tc.projects, tc.pick)
		switch {
		case tc.err == "" && (err != nil || got.Name != tc.want):
			t.Errorf("%s: %+v, %v; want %s", tc.name, got, err, tc.want)
		case tc.err != "" && (err == nil || !strings.Contains(err.Error(), tc.err)):
			t.Errorf("%s: error %v, want one containing %q", tc.name, err, tc.err)
		}
	}
}

// serveFixture is ServeOptions for the sub-project newFakeVite makes in
// mode, at sub-project path ui in a project named proj whose root is its
// parent, and the folder the fake writes into.
func serveFixture(t *testing.T, mode string) (ServeOptions, string) {
	t.Helper()
	dev, state := newFakeVite(t, mode)
	return ServeOptions{
		Root:    filepath.Dir(dev.Dir),
		Project: "proj",
		Sub: manifest.Project{Name: "ui", Path: filepath.Base(dev.Dir), Kind: manifest.ProjectKindSvelteKit,
			Components: &manifest.Components{Setup: "src/app.css", Exclude: []string{"src/lib/internal/**"}}},
		Timeout: dev.Timeout,
	}, state
}

func TestServeWritesTheHarnessAndStartsTheSubProjectsVite(t *testing.T) {
	opts, state := serveFixture(t, "serve")
	dir := filepath.Join(opts.Root, opts.Sub.Path)
	own := filepath.Join(dir, "vite.config.ts")
	if err := os.WriteFile(own, []byte("export default {};\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out syncBuffer
	opts.Output = &out
	s, err := Serve(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Stop() })
	config := filepath.Join(opts.Root, ".flai-cache", "components", "proj", "ui", HarnessConfig)
	written, err := os.ReadFile(config)
	if err != nil {
		t.Fatalf("Serve wrote no harness config at %s: %v", config, err)
	}
	for _, want := range []string{`"config":"` + own + `"`, `"root":"` + dir + `"`, `"setup":"src/app.css"`, `"include":["src/lib/**/*.svelte"]`, `"exclude":["src/lib/internal/**"]`} {
		if !strings.Contains(string(written), want) {
			t.Errorf("the harness config lacks %s:\n%s", want, written)
		}
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(config), "plugin.mjs")); err != nil {
		t.Errorf("the plugin is not beside the config: %v", err)
	}
	args, _ := os.ReadFile(filepath.Join(state, "args"))
	if !strings.HasPrefix(string(args), "--config\n"+config+"\n") {
		t.Errorf("vite ran with\n%s\nwant --config %s first", args, config)
	}
	if !strings.HasSuffix(s.URL(), HarnessPath) || s.PID() != readPID(t, state, "pid") {
		t.Errorf("URL %s, PID %d", s.URL(), s.PID())
	}
	resp, err := http.Get(s.URL())
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET %s: %d", s.URL(), resp.StatusCode)
	}
	if !strings.Contains(out.String(), "VITE ready") {
		t.Errorf("Output got %q, want vite's output", out.String())
	}
	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	waitGone(t, s.PID())
}

func TestServeRefusesASubProjectWithoutVite(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "ui"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := Serve(context.Background(), ServeOptions{Root: root, Project: "proj", Sub: pickProjects[0]})
	if !errors.Is(err, ErrViteMissing) || !strings.Contains(err.Error(), "sub-project ui's components") || !strings.Contains(err.Error(), "npm install") {
		t.Errorf("error %v, want ErrViteMissing naming ui, with the fix", err)
	}
}

func TestServeRefuses(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct {
		name string
		opts ServeOptions
		err  string
	}{
		{"unsupported kind", ServeOptions{Root: root, Project: "proj", Sub: pickProjects[1]}, "sub-project cli's components: its kind go is not sveltekit or svelte"},
		{"missing folder", ServeOptions{Root: root, Project: "proj", Sub: pickProjects[2]}, "its folder " + filepath.Join(root, "packages", "kit") + " is not there; check its path"},
		{"relative root", ServeOptions{Root: "proj", Project: "proj", Sub: pickProjects[0]}, "the project's root as an absolute path"},
		{"no project name", ServeOptions{Root: root, Sub: pickProjects[0]}, "and its name"},
	} {
		if _, err := Serve(context.Background(), tc.opts); err == nil || !strings.Contains(err.Error(), tc.err) {
			t.Errorf("%s: error %v, want one containing %q", tc.name, err, tc.err)
		}
	}
}

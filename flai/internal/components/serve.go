package components

// Serving a sub-project's components (S-0438, ADR-0163): PickSubProject
// chooses the sub-project from the manifest's projects, as the caller loaded
// them, and Serve writes its harness and starts its own Vite dev server on
// it. flai components serve and flai serve's supervisor both call them.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// PickSubProject is the sub-project whose components to serve among
// projects, the manifest's: the one name names, or, when name is empty, the
// one whose kind manifest.SupportsComponents takes. It refuses a name the
// manifest does not have, a named sub-project of another kind, a manifest
// with no such sub-project, and, with no name, one with more than one.
func PickSubProject(projects []manifest.Project, name string) (manifest.Project, error) {
	var supported []string
	for _, p := range projects {
		if manifest.SupportsComponents(p.Kind) {
			supported = append(supported, p.Name)
		}
	}
	if name == "" {
		switch len(supported) {
		case 0:
			return manifest.Project{}, fmt.Errorf("system-flow.yaml names no sub-project of kind %s, whose components flai serves; add one under projects, with its kind", componentKinds())
		case 1:
			name = supported[0]
		default:
			return manifest.Project{}, fmt.Errorf("system-flow.yaml has %d sub-projects of kind %s (%s), and flai serves one at a time; name the one to serve, such as --project %s", len(supported), componentKinds(), strings.Join(supported, ", "), supported[0])
		}
	}
	i := slices.IndexFunc(projects, func(p manifest.Project) bool { return p.Name == name })
	switch {
	case i < 0:
		return manifest.Project{}, UnknownSubProject(name, projects)
	case manifest.SupportsComponents(projects[i].Kind):
		return projects[i], nil
	case len(supported) == 0:
		return manifest.Project{}, fmt.Errorf("sub-project %q is of kind %s, and flai serves the components of a sub-project of kind %s only; system-flow.yaml names none, so add one under projects", name, projects[i].Kind, componentKinds())
	default:
		return manifest.Project{}, fmt.Errorf("sub-project %q is of kind %s, and flai serves the components of a sub-project of kind %s only; name one of %s", name, projects[i].Kind, componentKinds(), strings.Join(supported, ", "))
	}
}

// componentKinds is manifest.ComponentKinds as prose: sveltekit or svelte.
func componentKinds() string {
	return strings.Join(manifest.ComponentKinds, " or ")
}

// UnknownSubProject is the error for a sub-project name the manifest's
// projects do not have, listing the names they do.
func UnknownSubProject(name string, projects []manifest.Project) error {
	if len(projects) == 0 {
		return fmt.Errorf("system-flow.yaml has no sub-project %q, and names none under projects; add the sub-project there first", name)
	}
	names := make([]string, len(projects))
	for i, p := range projects {
		names[i] = p.Name
	}
	return fmt.Errorf("system-flow.yaml has no sub-project %q; name one of %s", name, strings.Join(names, ", "))
}

// ServeOptions says whose components Serve serves, and how.
type ServeOptions struct {
	// Root is the project's root folder, absolute.
	Root string
	// Project is the project's name, the manifest's name, which names the
	// folder of its harnesses under .flai-cache/components.
	Project string
	// Sub is the sub-project to serve, as PickSubProject picks it.
	Sub manifest.Project
	// Timeout is how long to wait for the dev server to answer; zero is
	// DevServerTimeout.
	Timeout time.Duration
	// Output, when set, also gets Vite's standard output and standard error,
	// as they come.
	Output io.Writer
}

// Serve writes the harness of opts.Sub into HarnessDir, over the
// sub-project's own Vite config when it has one, and starts the
// sub-project's own Vite dev server on it, returning once the harness page
// at HarnessPath answers, as StartDevServer does. It refuses a sub-project
// whose kind manifest.SupportsComponents does not take, one whose folder is
// missing, and one without vite installed, with an error that wraps
// ErrViteMissing; each error names the sub-project.
func Serve(ctx context.Context, opts ServeOptions) (*DevServer, error) {
	name := opts.Sub.Name
	if !filepath.IsAbs(opts.Root) || opts.Project == "" {
		return nil, fmt.Errorf("cannot serve sub-project %s's components: flai needs the project's root as an absolute path (got %q) and its name (got %q)", name, opts.Root, opts.Project)
	}
	settings, ok := opts.Sub.ComponentSettings()
	if !ok {
		return nil, fmt.Errorf("cannot serve sub-project %s's components: its kind %s is not %s, whose components flai serves; name a sub-project of one of those kinds", name, opts.Sub.Kind, componentKinds())
	}
	dir := filepath.Join(opts.Root, filepath.FromSlash(opts.Sub.Path))
	if info, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) || err == nil && !info.IsDir() {
		return nil, fmt.Errorf("cannot serve sub-project %s's components: its folder %s is not there; check its path in system-flow.yaml", name, dir)
	} else if err != nil {
		return nil, fmt.Errorf("cannot serve sub-project %s's components: cannot read its folder: %w", name, err)
	}
	own, err := ViteConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot serve sub-project %s's components: %w", name, err)
	}
	config, err := WriteHarness(HarnessDir(opts.Root, opts.Project, name), Harness{
		Dir:     dir,
		Config:  own,
		Setup:   settings.Setup,
		Include: settings.Include,
		Exclude: settings.Exclude,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot serve sub-project %s's components: %w", name, err)
	}
	s, err := StartDevServer(ctx, DevServerOptions{
		Dir:     dir,
		Config:  config,
		Path:    HarnessPath,
		Timeout: opts.Timeout,
		Output:  opts.Output,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot serve sub-project %s's components: %w", name, err)
	}
	return s, nil
}

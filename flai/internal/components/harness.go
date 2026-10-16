package components

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
)

// harnessFiles are the harness's plugin, page, and script, which
// WriteHarness writes beside the config it writes (S-0438, ADR-0163).
//
//go:embed harness
var harnessFiles embed.FS

// HarnessPath is the path, on the dev server, of the harness page the
// plugin serves; ?c=<path> names the component it renders and ?origin= the
// origin of the page that frames it.
const HarnessPath = "/__flai/components/"

// HarnessConfig is the name of the Vite config WriteHarness writes, which
// the dev server is run with by --config.
const HarnessConfig = "vite.config.mjs"

// viteConfigs are the names of a Vite config, in the order Vite looks for
// them.
var viteConfigs = []string{"vite.config.js", "vite.config.mjs", "vite.config.ts", "vite.config.cjs", "vite.config.mts", "vite.config.cts"}

// Harness is what WriteHarness needs to write one sub-project's harness.
type Harness struct {
	// Dir is the sub-project's folder, absolute. The dev server runs in it,
	// with it as Vite's root.
	Dir string
	// Config is the sub-project's own Vite config, absolute, as ViteConfig
	// finds it; empty when it has none.
	Config string
	// Setup is the module the harness imports before any component, relative
	// to Dir in slash form, as the components key's setup; empty for none.
	Setup string
	// Include are the components key's include globs; none means
	// DefaultInclude.
	Include []string
	// Exclude are the components key's exclude globs.
	Exclude []string
}

// HarnessDir is the folder of the harness of the sub-project named
// subProject in the project named project, whose root is root:
// .flai-cache/components/<project>/<sub-project>/.
func HarnessDir(root, project, subProject string) string {
	return filepath.Join(root, ".flai-cache", "components", project, subProject)
}

// ViteConfig is the Vite config in dir, absolute, that Vite would load there
// without --config, or empty when there is none.
func ViteConfig(dir string) (string, error) {
	for _, name := range viteConfigs {
		p := filepath.Join(dir, name)
		info, err := os.Stat(p)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			continue
		case err != nil:
			return "", fmt.Errorf("cannot look for a Vite config in %s: %w", dir, err)
		case info.Mode().IsRegular():
			return filepath.Abs(p)
		}
	}
	return "", nil
}

// harnessOptions are flaiComponents' options, which the written config
// carries as JSON.
type harnessOptions struct {
	Root    string   `json:"root"`
	Config  string   `json:"config"`
	Files   string   `json:"files"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
	Setup   string   `json:"setup"`
}

// WriteHarness writes h's harness into dir, which HarnessDir names: the
// embedded plugin.mjs, harness.mjs, and harness.html, and a HarnessConfig
// that loads h.Config, when there is one, with the sub-project's own Vite,
// merges the plugin into it, and sets Vite's root to h.Dir. It returns the
// config's absolute path. It refuses a Dir that is not absolute and a glob
// CheckGlob refuses, and rewrites what an earlier call wrote.
func WriteHarness(dir string, h Harness) (string, error) {
	if !filepath.IsAbs(h.Dir) {
		return "", fmt.Errorf("cannot write the component harness: the sub-project's folder %q is not absolute; pass it joined to the project root", h.Dir)
	}
	include := h.Include
	if len(include) == 0 {
		include = DefaultInclude
	}
	exclude := h.Exclude
	if exclude == nil {
		exclude = []string{}
	}
	for _, g := range slices.Concat(include, exclude) {
		if err := CheckGlob(g); err != nil {
			return "", fmt.Errorf("cannot write the component harness for %s: %w", h.Dir, err)
		}
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("cannot write the component harness: %w", err)
	}
	opts, err := json.Marshal(harnessOptions{Root: h.Dir, Config: h.Config, Files: dir, Include: include, Exclude: exclude, Setup: h.Setup})
	if err != nil {
		return "", fmt.Errorf("cannot write the component harness for %s: %w", h.Dir, err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot make the component harness's folder: %w; check that .flai-cache is writable", err)
	}
	entries, err := harnessFiles.ReadDir("harness")
	if err != nil {
		return "", fmt.Errorf("cannot read the embedded component harness: %w", err)
	}
	for _, e := range entries {
		data, err := harnessFiles.ReadFile(path.Join("harness", e.Name()))
		if err != nil {
			return "", fmt.Errorf("cannot read the embedded component harness: %w", err)
		}
		if err := writeHarnessFile(dir, e.Name(), data); err != nil {
			return "", err
		}
	}
	if err := writeHarnessFile(dir, HarnessConfig, []byte(harnessConfig(opts))); err != nil {
		return "", err
	}
	return filepath.Join(dir, HarnessConfig), nil
}

func writeHarnessFile(dir, name string, data []byte) error {
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		return fmt.Errorf("cannot write the component harness's %s: %w; check that .flai-cache is writable", name, err)
	}
	return nil
}

// harnessConfig is the source of the config, given flaiComponents' options
// as JSON, which is a JavaScript expression.
func harnessConfig(opts []byte) string {
	return `// Written by flai for its component harness (S-0438, ADR-0163), and written
// again each time flai starts the dev server: edit the sub-project's own
// config, not this. It loads that config with the sub-project's own Vite and
// merges flai's plugin into it.
import flaiComponents, { importVite } from './plugin.mjs';

const flai = ` + string(opts) + `;

const vite = await importVite(flai.root);
const own = flai.config
	? await vite.loadConfigFromFile({ command: 'serve', mode: 'development' }, flai.config, flai.root)
	: null;

// The cache is the harness's own, so that this config and the app's own dev
// server do not each make the other optimize the dependencies again; it is
// under root, as Vite's own is, so that the dev server may serve it.
export default vite.mergeConfig(own?.config ?? {}, {
	root: flai.root,
	cacheDir: 'node_modules/.vite-flai-components',
	plugins: [flaiComponents(flai)]
});
`
}

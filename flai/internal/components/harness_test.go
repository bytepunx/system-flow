package components

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// writeFile writes content to the slash path rel below dir.
func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// writtenOptions are the options the config at p passes the plugin, read
// back from its JSON.
func writtenOptions(t *testing.T, p string) (harnessOptions, string) {
	t.Helper()
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	src := string(data)
	_, rest, ok := strings.Cut(src, "\nconst flai = ")
	if !ok {
		t.Fatalf("config has no options:\n%s", src)
	}
	line, _, _ := strings.Cut(rest, "\n")
	var opts harnessOptions
	if err := json.Unmarshal([]byte(strings.TrimSuffix(line, ";")), &opts); err != nil {
		t.Fatalf("config's options %q are not JSON: %v", line, err)
	}
	return opts, src
}

func TestWriteHarnessForASvelteKitSubProject(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "apps", "web")
	writeFile(t, sub, "package.json", `{"devDependencies":{"@sveltejs/kit":"^2","vite":"^8"}}`)
	writeFile(t, sub, "svelte.config.js", "export default {};\n")
	writeFile(t, sub, "vite.config.ts", "export default {};\n")
	own, err := ViteConfig(sub)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(sub, "vite.config.ts"); own != want {
		t.Fatalf("ViteConfig = %q, want %q", own, want)
	}
	dir := HarnessDir(root, "acme", "web")
	if want := filepath.Join(root, ".flai-cache", "components", "acme", "web"); dir != want {
		t.Fatalf("HarnessDir = %q, want %q", dir, want)
	}
	h := Harness{Dir: sub, Config: own, Setup: "src/routes/layout.css", Include: []string{"src/lib/**/*.svelte"}, Exclude: []string{"src/lib/harness"}}
	got, err := WriteHarness(dir, h)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, HarnessConfig); got != want {
		t.Fatalf("WriteHarness = %q, want %q", got, want)
	}
	for _, name := range []string{"plugin.mjs", "harness.mjs", "harness.html"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("%s not written: %v", name, err)
		}
		embedded, err := harnessFiles.ReadFile(path.Join("harness", name))
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != string(embedded) || len(data) == 0 {
			t.Errorf("%s is not the embedded file", name)
		}
	}
	opts, src := writtenOptions(t, got)
	want := harnessOptions{Root: sub, Config: own, Files: dir, Include: h.Include, Exclude: h.Exclude, Setup: "src/routes/layout.css"}
	if opts.Root != want.Root || opts.Config != want.Config || opts.Files != want.Files || opts.Setup != want.Setup ||
		!slices.Equal(opts.Include, want.Include) || !slices.Equal(opts.Exclude, want.Exclude) {
		t.Errorf("options %+v, want %+v", opts, want)
	}
	for _, s := range []string{
		"import flaiComponents, { importVite } from './plugin.mjs';",
		"await importVite(flai.root)",
		"vite.loadConfigFromFile({ command: 'serve', mode: 'development' }, flai.config, flai.root)",
		"vite.mergeConfig(own?.config ?? {}, {",
		"root: flai.root,",
		"cacheDir: 'node_modules/.vite-flai-components',",
		"plugins: [flaiComponents(flai)]",
	} {
		if !strings.Contains(src, s) {
			t.Errorf("config lacks %q:\n%s", s, src)
		}
	}
	if strings.Contains(src, "from 'vite'") {
		t.Errorf("config imports vite by name, which resolves from .flai-cache, not the sub-project:\n%s", src)
	}
}

func TestWriteHarnessForAPlainSvelteSubProjectWithTheDefaults(t *testing.T) {
	sub := filepath.Join(t.TempDir(), "widgets")
	writeFile(t, sub, "vite.config.js", "export default {};\n")
	writeFile(t, sub, "vite.config.ts", "export default {};\n")
	own, err := ViteConfig(sub)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(sub, "vite.config.js"); own != want {
		t.Fatalf("ViteConfig = %q, want %q, the first Vite itself loads", own, want)
	}
	got, err := WriteHarness(filepath.Join(t.TempDir(), "harness"), Harness{Dir: sub, Config: own})
	if err != nil {
		t.Fatal(err)
	}
	opts, _ := writtenOptions(t, got)
	if opts.Config != own || opts.Setup != "" || !slices.Equal(opts.Include, DefaultInclude) || opts.Exclude == nil || len(opts.Exclude) != 0 {
		t.Errorf("options %+v, want its own config, no setup, the default include, and an empty exclude", opts)
	}
}

func TestWriteHarnessWithoutAnOwnConfigAndAgainOverItself(t *testing.T) {
	sub := t.TempDir()
	own, err := ViteConfig(sub)
	if err != nil || own != "" {
		t.Fatalf("ViteConfig = %q, %v, want none", own, err)
	}
	dir := filepath.Join(t.TempDir(), "harness")
	if _, err := WriteHarness(dir, Harness{Dir: sub, Setup: "src/app.css"}); err != nil {
		t.Fatal(err)
	}
	writeFile(t, dir, "plugin.mjs", "stale")
	got, err := WriteHarness(dir, Harness{Dir: sub})
	if err != nil {
		t.Fatal(err)
	}
	opts, src := writtenOptions(t, got)
	if opts.Config != "" || opts.Setup != "" {
		t.Errorf("options %+v, want no own config and no setup, as the second call asked", opts)
	}
	if !strings.Contains(src, "const own = flai.config\n\t?") {
		t.Errorf("config does not load its own config only when there is one:\n%s", src)
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "plugin.mjs")); string(data) == "stale" {
		t.Error("plugin.mjs was not written again")
	}
}

func TestWriteHarnessEscapesPathsAsJavaScriptStrings(t *testing.T) {
	sub := filepath.Join(t.TempDir(), `it's "web" \ `+"\u2028")
	got, err := WriteHarness(filepath.Join(t.TempDir(), "h"), Harness{Dir: sub, Setup: "src/a'b.css"})
	if err != nil {
		t.Fatal(err)
	}
	opts, src := writtenOptions(t, got)
	if opts.Root != sub || opts.Setup != "src/a'b.css" {
		t.Errorf("options %+v, want root %q and setup src/a'b.css", opts, sub)
	}
	if strings.Contains(src, "\u2028") {
		t.Error("config carries a raw line separator, which ends a JavaScript string")
	}
}

func TestWriteHarnessRefusals(t *testing.T) {
	sub := t.TempDir()
	if _, err := WriteHarness(t.TempDir(), Harness{Dir: "apps/web"}); err == nil || !strings.Contains(err.Error(), "not absolute") {
		t.Errorf("relative folder: err %v, want one saying it is not absolute", err)
	}
	if _, err := WriteHarness(t.TempDir(), Harness{Dir: sub, Exclude: []string{"src/[lib"}}); err == nil || !strings.Contains(err.Error(), "does not parse") {
		t.Errorf("bad glob: err %v, want one saying it does not parse", err)
	}
	blocked := filepath.Join(t.TempDir(), "file")
	writeFile(t, filepath.Dir(blocked), "file", "")
	if _, err := WriteHarness(filepath.Join(blocked, "harness"), Harness{Dir: sub}); err == nil || !strings.Contains(err.Error(), "writable") {
		t.Errorf("unwritable folder: err %v, want one saying to check it is writable", err)
	}
}

func TestThePluginServesTheHarnessPath(t *testing.T) {
	plugin, err := harnessFiles.ReadFile("harness/plugin.mjs")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(plugin), "'"+strings.TrimSuffix(HarnessPath, "/")+"'") {
		t.Errorf("plugin.mjs does not serve %s", HarnessPath)
	}
	page, err := harnessFiles.ReadFile("harness/harness.html")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(page), `src="/@id/virtual:flai-components/harness.mjs"`) || !strings.Contains(string(page), `id="flai-component"`) {
		t.Errorf("harness.html does not load the harness's module into its container:\n%s", page)
	}
}

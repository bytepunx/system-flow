package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// settingsPage is the operators' settings index, relative to this package.
const settingsPage = "../../docs/operators/settings.md"

// The settings index names every key of the configuration and the manifest,
// each in its own section, and every environment variable flai, flaiover, and
// install.sh read or set, so that a setting added without a row fails here.
func TestSettingsIndexIsComplete(t *testing.T) {
	if testing.Short() {
		t.Skip("reads the monorepo's docs and sources")
	}
	data, err := os.ReadFile(settingsPage)
	if err != nil {
		t.Fatalf("read the settings index: %v", err)
	}
	page := string(data)
	missing := func(where, text string, names []string) {
		t.Helper()
		for _, n := range names {
			if !strings.Contains(text, "`"+n+"`") {
				t.Errorf("docs/operators/settings.md: %s does not name `%s`", where, n)
			}
		}
	}
	missing("## flai configuration", pageSection(page, "## flai configuration"), keyPaths(reflect.TypeFor[config.Config](), "json", ""))
	missing("## Project manifest", pageSection(page, "## Project manifest"), keyPaths(reflect.TypeFor[manifest.Manifest](), "yaml", ""))
	// flai also hands the container its secrets' paths and an image its build
	// arguments, which the flaiover section lists.
	missing("## Environment variables or ## flaiover container", pageSection(page, "## Environment variables")+pageSection(page, "## flaiover container"), append(
		namesIn(t, "..", `\.go$`, `_test\.go$`, regexp.MustCompile(`"((?:FLAI|FLAIOVER)_[A-Z0-9_]+)|os\.(?:Getenv|LookupEnv)\("([A-Z][A-Z0-9_]*)"`)),
		namesIn(t, "../../install.sh", `install\.sh$`, `^$`, regexp.MustCompile(`\$\{([A-Z][A-Z0-9_]*):-`))...))
	missing("## flaiover container", pageSection(page, "## flaiover container"),
		namesIn(t, "../../flaiover", `/(src/.*\.ts|server\.js|scheme\.js)$`, `(\.test\.ts|\.spec\.ts|/testing\.ts)$|/node_modules/|/\.svelte-kit/|/build/`, regexp.MustCompile(`env\.([A-Z][A-Z0-9_]+)`)))
}

// The walk that feeds the test above finds what it should: a key of every
// shape the schemas use, and a variable where it is only a prefix of a string.
func TestSettingsIndexWalks(t *testing.T) {
	type leaf struct {
		On *bool `json:"on"`
	}
	type schema struct {
		Plain  string              `json:"plain"`
		Nested struct{ A int }     `json:"nested"`
		Map    map[string][]string `json:"map,omitempty"`
		Rich   map[string]leaf     `json:"rich,omitzero"`
		List   []leaf              `json:"list"`
		Ptr    *leaf               `json:"ptr"`
		Words  []string            `json:"words"`
	}
	got := strings.Join(keyPaths(reflect.TypeFor[schema](), "json", ""), " ")
	if want := "list[].on map.<name> nested.A plain ptr.on rich.<name>.on words"; got != want {
		t.Errorf("key paths:\n got %s\nwant %s", got, want)
	}
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "a.go"), []byte(`x := "FLAI_STORY="+id; os.Getenv("HOME"); "FLAI_HOST_TEST_ENV"`), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "a_test.go"), []byte(`"FLAI_ONLY_IN_TESTS"`), 0o644)
	re := regexp.MustCompile(`"((?:FLAI|FLAIOVER)_[A-Z0-9_]+)|os\.(?:Getenv|LookupEnv)\("([A-Z][A-Z0-9_]*)"`)
	if got := strings.Join(namesIn(t, dir, `\.go$`, `_test\.go$`, re), " "); got != "FLAI_STORY HOME" {
		t.Errorf("names: got %q", got)
	}
}

// pageSection is the text from a level-two heading to the next one.
func pageSection(page, heading string) string {
	_, rest, ok := strings.Cut(page, "\n"+heading+"\n")
	if !ok {
		return ""
	}
	if i := strings.Index(rest, "\n## "); i >= 0 {
		return rest[:i]
	}
	return rest
}

// keyPaths are the dotted paths of every leaf of a schema, by the given
// struct tag: a map's key is <name> and a list's element is [].
func keyPaths(t reflect.Type, tag, prefix string) []string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var out []string
	for i := range t.NumField() {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get(tag), ",")
		if name == "" {
			name = f.Name
		}
		out = append(out, fieldPaths(f.Type, tag, prefix+name)...)
	}
	sort.Strings(out)
	return out
}

func fieldPaths(t reflect.Type, tag, path string) []string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch {
	case t.Kind() == reflect.Struct:
		return keyPaths(t, tag, path+".")
	case t.Kind() == reflect.Map && isStruct(t.Elem()):
		return keyPaths(t.Elem(), tag, path+".<name>.")
	case t.Kind() == reflect.Map:
		return []string{path + ".<name>"}
	case t.Kind() == reflect.Slice && isStruct(t.Elem()):
		return keyPaths(t.Elem(), tag, path+"[].")
	}
	return []string{path}
}

func isStruct(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

// namesIn are the variable names re's first non-empty group matches in the
// files under root whose slash path matches include and not exclude, leaving
// out the hooks only tests set (a name with _TEST in it), in order of first
// appearance.
func namesIn(t *testing.T, root, include, exclude string, re *regexp.Regexp) []string {
	t.Helper()
	in, ex := regexp.MustCompile(include), regexp.MustCompile(exclude)
	seen := map[string]bool{}
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		slash := filepath.ToSlash(path)
		if !in.MatchString(slash) || ex.MatchString(slash) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range re.FindAllStringSubmatch(string(data), -1) {
			for _, n := range m[1:] {
				if n != "" && !strings.Contains(n, "_TEST") && !seen[n] {
					seen[n] = true
					out = append(out, n)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("read %s: %v", root, err)
	}
	return out
}

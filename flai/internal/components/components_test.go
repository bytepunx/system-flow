package components

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// tree writes each file, by slash path relative to dir, with empty content.
func tree(t *testing.T, dir string, files ...string) {
	t.Helper()
	for _, f := range files {
		p := filepath.Join(dir, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func paths(comps []Component) []string {
	out := make([]string, len(comps))
	for i, c := range comps {
		out[i] = c.Path
	}
	return out
}

func TestFindDefaultsToTheLibraryAndSkipsNodeModulesAndDotFolders(t *testing.T) {
	root := t.TempDir()
	tree(t, root,
		"web/src/lib/Zeta.svelte",
		"web/src/lib/Alpha.svelte",
		"web/src/lib/forms/Field.svelte",
		"web/src/lib/forms/field.ts",
		"web/src/routes/+page.svelte",
		"web/src/lib/node_modules/pkg/Dep.svelte",
		"web/src/lib/.cache/Hidden.svelte",
		"web/node_modules/svelte/src/lib/Inner.svelte",
	)
	got, err := Find(root, "web", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"src/lib/Alpha.svelte", "src/lib/Zeta.svelte", "src/lib/forms/Field.svelte"}
	if !slices.Equal(paths(got), want) {
		t.Fatalf("paths %v, want %v", paths(got), want)
	}
	c := got[2]
	if c.Project != "web" || c.Name != "Field" || c.Props == nil || c.Error != "" {
		t.Errorf("component %+v, want project web, name Field, empty props, no error", c)
	}
}

func TestFindByIncludeAndExcludeGlobs(t *testing.T) {
	root := t.TempDir()
	tree(t, root,
		"apps/web/src/lib/A.svelte",
		"apps/web/src/lib/harness/H.svelte",
		"apps/web/src/lib/deep/er/B.svelte",
		"apps/web/src/lib/B.test.svelte",
		"apps/web/src/components/C.svelte",
		"apps/web/src/routes/+page.svelte",
	)
	tests := []struct {
		name             string
		include, exclude []string
		want             []string
	}{
		{"double star at the front", []string{"**/*.svelte"}, []string{"src/routes"},
			[]string{"src/components/C.svelte", "src/lib/A.svelte", "src/lib/B.test.svelte", "src/lib/deep/er/B.svelte", "src/lib/harness/H.svelte"}},
		{"plain folder excludes below it", []string{"src/lib/**/*.svelte"}, []string{"src/lib/harness", "**/*.test.svelte"},
			[]string{"src/lib/A.svelte", "src/lib/deep/er/B.svelte"}},
		{"plain folder includes below it", []string{"src/components"}, nil,
			[]string{"src/components/C.svelte"}},
		{"several includes", []string{"src/lib/*.svelte", "src/components/*.svelte"}, []string{"src/lib/B.*"},
			[]string{"src/components/C.svelte", "src/lib/A.svelte"}},
		{"one file", []string{"src/lib/deep/er/B.svelte"}, nil,
			[]string{"src/lib/deep/er/B.svelte"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Find(root, "apps/web", tt.include, tt.exclude)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(paths(got), tt.want) {
				t.Errorf("paths %v, want %v", paths(got), tt.want)
			}
			for _, c := range got {
				if c.Project != "apps/web" {
					t.Errorf("project %q, want apps/web", c.Project)
				}
			}
		})
	}
}

func TestFindRefusesABadGlobAndAMissingSubProject(t *testing.T) {
	root := t.TempDir()
	tree(t, root, "web/src/lib/A.svelte")
	if _, err := Find(root, "web", []string{"src/[lib"}, nil); err == nil || !strings.Contains(err.Error(), "does not parse") {
		t.Errorf("bad include: err %v, want one saying it does not parse", err)
	}
	if _, err := Find(root, "web", nil, []string{"/abs"}); err == nil || !strings.Contains(err.Error(), "absolute") {
		t.Errorf("bad exclude: err %v, want one saying it is absolute", err)
	}
	if _, err := Find(root, "gone", nil, nil); err == nil || !strings.Contains(err.Error(), "gone") || !strings.Contains(err.Error(), "manifest") {
		t.Errorf("missing sub-project: err %v, want one naming it and the manifest", err)
	}
}

func TestCheckGlob(t *testing.T) {
	for _, g := range []string{"src/lib/**/*.svelte", "**", "src/lib", "src/*/X?.svelte", "src/[a-z]*.svelte"} {
		if err := CheckGlob(g); err != nil {
			t.Errorf("CheckGlob(%q) = %v, want nil", g, err)
		}
	}
	bad := map[string]string{
		"":              "empty",
		"/src":          "absolute",
		"C:/src":        "absolute",
		"src//lib":      "empty, . or ..",
		"../x":          "empty, . or ..",
		"src/lib/":      "empty, . or ..",
		"src/a**/x":     "inside a segment",
		"src/[lib/*.sv": "does not parse",
	}
	for g, want := range bad {
		if err := CheckGlob(g); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("CheckGlob(%q) = %v, want an error saying %q", g, err, want)
		}
	}
}

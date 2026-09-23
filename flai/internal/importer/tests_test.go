package importer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func names(cmds []TestCommand) []string {
	out := []string{}
	for _, c := range cmds {
		out = append(out, c.Name)
	}
	return out
}

// S-0098: what tests an imported repository already has.
func TestDetectTests(t *testing.T) {
	t.Run("its own Makefile's test target stands for everything", func(t *testing.T) {
		root := t.TempDir()
		write(t, root, "Makefile", "build:\n\tgo build\n\ntest: build\n\tgo test ./...\n")
		write(t, root, "go.mod", "module x\n")
		if got := names(DetectTests(root, true, nil)); !reflect.DeepEqual(got, []string{"make test"}) {
			t.Errorf("got %v", got)
		}
		// the same file, brought by the template, is not the repository's
		if got := names(DetectTests(root, false, nil)); !reflect.DeepEqual(got, []string{"go test ./..."}) {
			t.Errorf("template Makefile: %v", got)
		}
	})
	t.Run("a Makefile without a test target is not a way to test", func(t *testing.T) {
		root := t.TempDir()
		write(t, root, "Makefile", "build:\n\tgo build\n")
		write(t, root, "Cargo.toml", "[package]\n")
		if got := names(DetectTests(root, true, nil)); !reflect.DeepEqual(got, []string{"cargo test"}) {
			t.Errorf("got %v", got)
		}
	})
	t.Run("each sub-project the way its ecosystem tests", func(t *testing.T) {
		root := t.TempDir()
		write(t, root, "go.mod", "module x\n")
		write(t, root, "web/package.json", `{"scripts":{"test":"vitest run"}}`)
		write(t, root, "web/pnpm-lock.yaml", "")
		write(t, root, "tools/package.json", `{"scripts":{"test":"echo \"Error: no test specified\" && exit 1"}}`)
		write(t, root, "ml/pyproject.toml", "[project]\n")
		got := DetectTests(root, false, []Project{{Path: "web"}, {Path: "tools"}, {Path: "ml"}, {Path: "."}})
		if n := names(got); !reflect.DeepEqual(n, []string{"go test ./...", "pnpm test (web)", "python3 -m pytest (ml)"}) {
			t.Errorf("got %v", n)
		}
		if got[1].Dir != "web" || !reflect.DeepEqual(got[1].Command, []string{"pnpm", "test"}) {
			t.Errorf("web: %+v", got[1])
		}
	})
	t.Run("none found is an empty list", func(t *testing.T) {
		root := t.TempDir()
		write(t, root, "README.md", "# x\n")
		if got := DetectTests(root, false, nil); len(got) != 0 {
			t.Errorf("got %v", got)
		}
	})
}

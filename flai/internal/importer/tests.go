package importer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TestCommand is one command that tests an imported repository (S-0098): an
// argument list run as it stands in Dir, relative to the root, never through
// a shell.
type TestCommand struct {
	Name    string   `json:"name"`
	Dir     string   `json:"dir"`
	Command []string `json:"command"`
}

var makeTestTarget = regexp.MustCompile(`(?m)^test\s*:`)

// npm writes this test script into every package.json it creates.
const npmPlaceholderTest = `echo "Error: no test specified" && exit 1`

// DetectTests finds what tests the repository already has. A Makefile of its
// own with a test target is taken as its way to run everything, as the
// people who wrote it do; ownMakefile says whether the root had one before
// the import, since the template brings one that is not the repository's.
// Otherwise each place with a build file (the root and every sub-project)
// is tested the way that ecosystem is: go test, the package manager's test
// script, cargo test, pytest. None found is an empty list, not an error.
func DetectTests(root string, ownMakefile bool, projects []Project) []TestCommand {
	if ownMakefile {
		if data, err := os.ReadFile(filepath.Join(root, "Makefile")); err == nil && makeTestTarget.Match(data) {
			return []TestCommand{{Name: "make test", Dir: ".", Command: []string{"make", "test"}}}
		}
	}
	dirs := []string{"."}
	for _, p := range projects {
		if d := filepath.ToSlash(filepath.Clean(p.Path)); d != "." {
			dirs = append(dirs, d)
		}
	}
	var out []TestCommand
	seen := map[string]bool{}
	for _, d := range dirs {
		if seen[d] {
			continue
		}
		seen[d] = true
		abs := filepath.Join(root, d)
		add := func(args ...string) {
			name := strings.Join(args, " ")
			if d != "." {
				name += " (" + d + ")"
			}
			out = append(out, TestCommand{Name: name, Dir: d, Command: args})
		}
		if exists(filepath.Join(abs, "go.mod")) {
			add("go", "test", "./...")
		}
		if script := npmTestScript(filepath.Join(abs, "package.json")); script != "" {
			switch {
			case exists(filepath.Join(abs, "pnpm-lock.yaml")):
				add("pnpm", "test")
			case exists(filepath.Join(abs, "yarn.lock")):
				add("yarn", "test")
			default:
				add("npm", "test")
			}
		}
		if exists(filepath.Join(abs, "Cargo.toml")) {
			add("cargo", "test")
		}
		if exists(filepath.Join(abs, "pyproject.toml")) || exists(filepath.Join(abs, "pytest.ini")) {
			add("python3", "-m", "pytest")
		}
	}
	return out
}

func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// npmTestScript is package.json's test script, empty when there is none or
// it is the placeholder npm init writes.
func npmTestScript(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if json.Unmarshal(data, &pkg) != nil {
		return ""
	}
	s := strings.TrimSpace(pkg.Scripts["test"])
	if s == npmPlaceholderTest {
		return ""
	}
	return s
}

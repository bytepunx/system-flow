package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// lintDoc runs flai doc lint in root with stdin and answers stdout, stderr,
// and the exit code.
func lintDoc(t *testing.T, root, stdin string, args ...string) (string, string, int) {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root}
	rootCmd := newRootCmdWith(a)
	rootCmd.SetArgs(append([]string{"doc", "lint"}, args...))
	rootCmd.SetIn(strings.NewReader(stdin))
	code := 0
	if err := rootCmd.Execute(); err != nil {
		var ee *exitError
		if !errors.As(err, &ee) {
			t.Fatalf("doc lint failed without an exit code: %v", err)
		}
		code = ee.code
		if ee.msg != "" {
			a.fail(err)
		}
	}
	return out.String(), errOut.String(), code
}

func docLintProject(t *testing.T) string {
	t.Helper()
	root := tempProject(t)
	if err := os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD013: false\nMD041: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

type lintAnswer struct {
	Path     string
	Findings []struct {
		Line                 int
		Rule, Alias, Message string
	}
}

func TestDocLintAnswersTheFindingsOfTheContentGiven(t *testing.T) {
	root := docLintProject(t)
	content := "# Title\n\n\n## Same\n\ntext \n\n## Same\n"
	out, errOut, code := lintDoc(t, root, content, "docs/x.md", "--stdin", "--json")
	if code != exitDocLintFindings {
		t.Fatalf("exit %d, want %d: %s", code, exitDocLintFindings, errOut)
	}
	var got lintAnswer
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got.Path != "docs/x.md" {
		t.Errorf("path %q", got.Path)
	}
	want := map[string]int{"MD012": 3, "MD009": 6, "MD024": 8}
	if len(got.Findings) != len(want) {
		t.Fatalf("findings %+v, want %v", got.Findings, want)
	}
	for _, f := range got.Findings {
		if want[f.Rule] != f.Line || f.Alias == "" || f.Message == "" {
			t.Errorf("finding %+v, want %s at line %d with an alias and a message", f, f.Rule, want[f.Rule])
		}
	}
}

func TestDocLintAnswersAnEmptyListForCleanContent(t *testing.T) {
	root := docLintProject(t)
	out, errOut, code := lintDoc(t, root, "# Title\n\ntext\n", "design/system/new.md", "--stdin", "--json")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if !strings.Contains(out, `"findings": []`) {
		t.Errorf("clean content answers no empty list: %s", out)
	}
}

func TestDocLintPrintsPathLineRuleAndMessage(t *testing.T) {
	root := docLintProject(t)
	out, _, code := lintDoc(t, root, "# Title\n\ntext \n", "wip/agents/S-0001.md", "--stdin")
	if code != exitDocLintFindings {
		t.Fatalf("exit %d", code)
	}
	if want := "wip/agents/S-0001.md:3 MD009/no-trailing-spaces Trailing spaces"; !strings.HasPrefix(out, want) {
		t.Errorf("printed %q, want it to start %q", out, want)
	}
}

func TestDocLintRefusesWhatItCannotLint(t *testing.T) {
	root := docLintProject(t)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"no --stdin", []string{"docs/x.md"}},
		{"outside the folders", []string{"src/x.md", "--stdin"}},
		{"above the repository", []string{"../x.md", "--stdin"}},
		{"not markdown", []string{"docs/x.txt", "--stdin"}},
	} {
		if _, errOut, code := lintDoc(t, root, "# Title\n", tc.args...); code != exitDocLintUnusable {
			t.Errorf("%s: exit %d, want %d: %s", tc.name, code, exitDocLintUnusable, errOut)
		}
	}
}

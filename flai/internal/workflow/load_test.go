package workflow_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// boardRel is the workflow file's path in a project whose design folder is
// design.
const boardRel = "design/workflows/board.md"

// baselineEdited is the baseline file's text with each pair of old and new
// strings replaced once.
func baselineEdited(t *testing.T, edits ...string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("baseline", "board.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for i := 0; i+1 < len(edits); i += 2 {
		if !strings.Contains(text, edits[i]) {
			t.Fatalf("the baseline lacks %q", edits[i])
		}
		text = strings.Replace(text, edits[i], edits[i+1], 1)
	}
	return text
}

// projectWith is a project root whose design folder holds content as its
// workflow file.
func projectWith(t *testing.T, content string) string {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(boardRel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// The edits that give the baseline a state of its own, parked, that items
// reach from backlog and leave back to it.
var parkedEdits = []string{
	"  - {name: cancelled, kind: cancelled}\n", "  - {name: cancelled, kind: cancelled}\n  - {name: parked, kind: queue}\n",
	"  - {from: ready, to: backlog}\n", "  - {from: ready, to: backlog}\n  - {from: backlog, to: parked}\n  - {from: parked, to: backlog}\n",
}

// S-0755: a project's valid workflow file is the workflow Load gives, on the
// path relative to the root; a project with none gives nil, to run the
// baseline.
func TestLoadReadsTheProjectsFile(t *testing.T) {
	w, err := workflow.Load(projectWith(t, baselineEdited(t, parkedEdits...)), "design")
	if err != nil {
		t.Fatal(err)
	}
	if w == nil || w.Path != boardRel || !slices.Contains(w.Names(), "parked") {
		t.Fatalf("Load gave %+v, want the file's workflow with parked at %s", w, boardRel)
	}
	if w, err := workflow.Load(t.TempDir(), "design"); w != nil || err != nil {
		t.Errorf("a project with no file: %+v, %v; want nil, nil", w, err)
	}
}

// S-0755: a file flai cannot run is a *LoadError naming its first error
// finding's code, path, and line, as flai check reports it; a file with a
// warning alone runs.
func TestLoadRefusesAFileItCannotRun(t *testing.T) {
	cases := []struct {
		name, content, code string
		line                int
		says                string
	}{
		{"a file that does not parse", "no front matter\n", workflow.CodeFile, 1, "the file has no front matter"},
		{"a field outside the allowed values", baselineEdited(t, "{name: done, kind: accepted}", "{name: done, kind: finished}"), workflow.CodeFile, 10, `states[4].kind: state done's kind "finished"`},
		{"no pull", baselineEdited(t, "effects: [pull], ", ""), workflow.CodeNoPull, 12, "no transition gives stories the pull effect"},
		{"a dead end", baselineEdited(t,
			"  - {name: cancelled, kind: cancelled}\n", "  - {name: cancelled, kind: cancelled}\n  - {name: stuck, kind: waiting}\n",
			"  - {from: ready, to: backlog}\n", "  - {from: ready, to: backlog}\n  - {from: review, to: stuck}\n"),
			workflow.CodeDeadEnd, 12, "state stuck leads by no chain of transitions"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, err := workflow.Load(projectWith(t, tc.content), "design")
			var le *workflow.LoadError
			if w != nil || !errors.As(err, &le) {
				t.Fatalf("Load gave %+v, %v; want a *LoadError", w, err)
			}
			if le.Code != tc.code || le.Path != boardRel || le.Line != tc.line || !strings.Contains(le.Message, tc.says) {
				t.Errorf("Load refused with %+v; want %s at %s:%d saying %q", le, tc.code, boardRel, tc.line, tc.says)
			}
			if want := fmt.Sprintf("%s:%d: %s: ", boardRel, tc.line, tc.code); !strings.HasPrefix(err.Error(), want) || !strings.HasSuffix(err.Error(), "remove it to run the baseline workflow flai carries") {
				t.Errorf("the error reads %q; want path:line: code: message, then what to do", err)
			}
		})
	}
	w, err := workflow.Load(projectWith(t, baselineEdited(t,
		"  - {name: cancelled, kind: cancelled}\n", "  - {name: cancelled, kind: cancelled}\n  - {name: parked, kind: queue}\n",
		"  - {from: ready, to: backlog}\n", "  - {from: ready, to: backlog}\n  - {from: parked, to: backlog}\n")), "design")
	if err != nil || w == nil || !slices.Contains(w.Names(), "parked") {
		t.Fatalf("a file whose only finding is %s: %+v, %v; want it to run", workflow.CodeUnreachable, w, err)
	}
	if _, ok := findingOf(w.Findings(), workflow.CodeUnreachable); !ok {
		t.Errorf("the file has no %s finding: %+v", workflow.CodeUnreachable, w.Findings())
	}
}

// findingOf is the first finding of code.
func findingOf(fs []workflow.Finding, code string) (workflow.Finding, bool) {
	for _, f := range fs {
		if f.Code == code {
			return f, true
		}
	}
	return workflow.Finding{}, false
}

package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0220, ADR-0090: flai thread reply posts a recommendation citing a
// source, which leaves the thread awaiting the operator; a source that names
// no file, or a heading not in it, is refused; flai thread show marks the
// recommendation; and flai thread confirm makes it the answer.
func TestAThreadRecommendationIsPostedWithItsSourceAndConfirmed(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	if err := os.MkdirAll(filepath.Join(root, "design", "system"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "design", "system", "plan.md"), []byte("---\ntitle: Plan\n---\n\n# Plan\n\n## Shape\n\nThe CLI first.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "thread", "new", "--on", "design/system/plan.md", "--by", "planner", "Which first?", "The CLI or the dashboard?"); code != 0 {
		t.Fatalf("thread new: %d %s", code, errOut)
	}

	// errOut is the JSON log line, its quotes escaped
	for source, want := range map[string]string{
		"design/system/nope.md":       "design/system/nope.md does not exist",
		"design/system/plan.md#Nope":  `Nope\" is not in design/system/plan.md`,
		"#Shape":                      `#Shape\" names no file: give <path> or <path>#<heading>`,
		"design/system/plan.md #Nope": `Nope\" is not in design/system/plan.md`,
	} {
		if _, errOut, code := runIn(t, root, "thread", "reply", "TH-0001", "--by", "orchestrator", "--recommend", "--source", source, "The CLI."); code == 0 || !strings.Contains(errOut, want) {
			t.Errorf("--source %q: exit %d, %q, want it refused with %q", source, code, errOut, want)
		}
	}

	out, errOut, code := runIn(t, root, "thread", "reply", "TH-0001", "--by", "orchestrator", "--recommend", "--source", "design/system/plan.md#Shape", "The CLI, as the plan says.")
	if code != 0 || !strings.HasPrefix(out, "TH-0001 is open (2 entries)") || !strings.Contains(out, "flai thread confirm TH-0001") {
		t.Fatalf("recommend: %d %q %s", code, out, errOut)
	}
	out, _, _ = runIn(t, root, "thread", "show", "TH-0001")
	for _, want := range []string{"the recommendation of 2026-09-15T21:00:00Z orchestrator awaits the operator", "2026-09-15T21:00:00Z orchestrator (recommendation)\nThe CLI, as the plan says.\n\nSource: design/system/plan.md § Shape"} {
		if !strings.Contains(out, want) {
			t.Errorf("show lacks %q:\n%s", want, out)
		}
	}

	if _, errOut, code := runIn(t, root, "thread", "confirm", "TH-0001", "--by", "orchestrator"); code == 0 || !strings.Contains(errOut, "cannot confirm it") {
		t.Errorf("the recommender confirming: exit %d, %q", code, errOut)
	}
	out, errOut, code = runIn(t, root, "thread", "confirm", "TH-0001", "--by", "alex", "--json")
	if code != 0 {
		t.Fatalf("confirm: %d %s", code, errOut)
	}
	var view struct {
		Status  string `json:"status"`
		Entries []struct {
			Author string `json:"author"`
			Text   string `json:"text"`
			Source *struct {
				Path    string `json:"path"`
				Heading string `json:"heading"`
			} `json:"source"`
		} `json:"entries"`
		Pending any `json:"pending_recommendation"`
	}
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatalf("confirm --json: %v\n%s", err, out)
	}
	last := view.Entries[len(view.Entries)-1]
	if view.Status != "answered" || view.Pending != nil || last.Author != "alex" || !strings.HasPrefix(last.Text, "Confirmed the recommendation of 2026-09-15T21:00:00Z orchestrator.") || last.Source == nil || last.Source.Path != "design/system/plan.md" || last.Source.Heading != "Shape" {
		t.Errorf("confirmed: %s", out)
	}
	if _, errOut, code := runIn(t, root, "thread", "confirm", "TH-0001", "--by", "alex"); code == 0 || !strings.Contains(errOut, "no recommendation awaiting confirmation") {
		t.Errorf("a second confirm: exit %d, %q", code, errOut)
	}
}

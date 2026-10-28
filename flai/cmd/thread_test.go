package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// ADR-0206: flai thread reply --recommend --operator-only marks the
// recommendation for the operator alone, written above its source line and
// shown as operator_only in --json; --operator-only without --recommend is
// refused, writing nothing.
func TestAThreadRecommendationIsMarkedForTheOperatorAlone(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	if err := os.MkdirAll(filepath.Join(root, "design", "system"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "design", "system", "plan.md"), []byte("---\ntitle: Plan\n---\n\n# Plan\n\n## Hosts\n\nOne host.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "thread", "new", "--on", "design/system/plan.md", "--by", "planner-S-0001", "Which host?", "Buy a second host?"); code != 0 {
		t.Fatalf("thread new: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "thread", "reply", "TH-0001", "--by", "orchestrator", "--operator-only", "Buy it."); code == 0 || !strings.Contains(errOut, "give --recommend with it") {
		t.Errorf("--operator-only without --recommend: exit %d, %q", code, errOut)
	}
	out, errOut, code := runIn(t, root, "thread", "reply", "TH-0001", "--by", "orchestrator", "--recommend", "--operator-only", "--source", "design/system/plan.md#Hosts", "Buy it.", "--json")
	if code != 0 {
		t.Fatalf("reply: %d %s", code, errOut)
	}
	var view struct {
		Entries []struct {
			Text           string `json:"text"`
			Recommendation bool   `json:"recommendation"`
			OperatorOnly   bool   `json:"operator_only"`
			Source         *struct {
				Path string `json:"path"`
			} `json:"source"`
		} `json:"entries"`
		Pending *struct {
			OperatorOnly bool `json:"operator_only"`
		} `json:"pending_recommendation"`
	}
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatalf("reply --json: %v\n%s", err, out)
	}
	if len(view.Entries) != 2 {
		t.Fatalf("entries: %s", out)
	}
	last := view.Entries[1]
	if !last.Recommendation || !last.OperatorOnly || last.Source == nil || last.Source.Path != "design/system/plan.md" || view.Pending == nil || !view.Pending.OperatorOnly {
		t.Errorf("marked: %s", out)
	}
	if want := "Buy it.\n\nFor the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.\n\nSource: design/system/plan.md § Hosts"; last.Text != want {
		t.Errorf("text %q, want %q", last.Text, want)
	}
}

// S-0369: a thread entry without --by is FLAI_AGENT's when it is set, else
// the config's author, so the operator's own shell answers as the operator.
func TestAThreadEntryIsTheConfigAuthorsWithoutFLAIAGENT(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	for _, args := range [][]string{{"config", "set", "author", "olive"}, {"thread", "new", "--on", "system-flow.yaml", "--by", "planner", "Which first?", "The CLI or the dashboard?"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
	}
	lastAuthor := func(args ...string) string {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"thread", "reply", "TH-0001", "--json"}, args...)...)
		if code != 0 {
			t.Fatalf("reply %v: %d %s", args, code, errOut)
		}
		var view struct {
			Entries []struct {
				Author string `json:"author"`
			} `json:"entries"`
		}
		if err := json.Unmarshal([]byte(out), &view); err != nil {
			t.Fatalf("reply --json: %v\n%s", err, out)
		}
		return view.Entries[len(view.Entries)-1].Author
	}
	if got := lastAuthor("The CLI."); got != "olive" {
		t.Errorf("without FLAI_AGENT the config author replied, got %q", got)
	}
	t.Setenv("FLAI_AGENT", "claude")
	if got := lastAuthor("The dashboard."); got != "claude" {
		t.Errorf("with FLAI_AGENT the agent replied, got %q", got)
	}
	if got := lastAuthor("--by", "dana", "Both."); got != "dana" {
		t.Errorf("--by wins, got %q", got)
	}
}

// answerProject is a project owned by alex with TH-0001, "Which first?",
// which claude opened on the manifest.
func answerProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nowner: alex\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "thread", "new", "--on", "system-flow.yaml", "--by", "claude", "Which first?", "The CLI or the dashboard?"); code != 0 {
		t.Fatalf("thread new: %d %s", code, errOut)
	}
	return root
}

// answered is a thread reply or confirm's --json: the thread's status and
// entries, and the digest entry its answer became, if any.
type answered struct {
	Status  string `json:"status"`
	Entries []struct {
		Author string `json:"author"`
		Text   string `json:"text"`
	} `json:"entries"`
	Digest *struct {
		ID      string `json:"id"`
		Path    string `json:"path"`
		Action  string `json:"action"`
		Warning string `json:"warning"`
	} `json:"digest"`
}

// answerAt is the clock at minute m past the thread's opening, so each
// answer has a time of its own, as the entry citing its thread records it.
func answerAt(m int) time.Time { return time.Date(2026, 9, 15, 21, m, 0, 0, time.UTC) }

// answer runs a thread reply or confirm with --json in root at minute m, and
// fails the test unless it exits 0 with the thread's view; raw is the output.
func answer(t *testing.T, root string, m int, args ...string) (view answered, raw map[string]any) {
	t.Helper()
	out, errOut, code := runInAt(t, root, answerAt(m), append([]string{"--json"}, args...)...)
	if code != 0 {
		t.Fatalf("%v: %d %s", args, code, errOut)
	}
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatalf("%v --json: %v\n%s", args, err, out)
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		t.Fatal(err)
	}
	return view, raw
}

// entryText is the digest entry at path, from root.
func entryText(t *testing.T, root, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0433 criteria 1 to 3: the operator's reply on a thread an agent opened
// is written as a digest entry, named under digest in --json and on a line of
// its own in the text; a second answer updates that entry.
func TestThreadReplyByTheOperatorWritesADigestEntry(t *testing.T) {
	root := answerProject(t)
	view, _ := answer(t, root, 1, "thread", "reply", "TH-0001", "--by", "alex", "--", "The CLI. It is what agents use.")
	d := view.Digest
	if view.Status != "answered" || d == nil || d.ID != "D-0001" || d.Action != "written" || d.Warning != "" || !strings.HasPrefix(d.Path, "design/digest/D-0001-") {
		t.Fatalf("reply --json: %+v digest %+v", view, d)
	}
	text := entryText(t, root, d.Path)
	for _, want := range []string{"kind: answer", "TH-0001", "The CLI. It is what agents use."} {
		if !strings.Contains(text, want) {
			t.Errorf("entry lacks %q:\n%s", want, text)
		}
	}

	out, errOut, code := runInAt(t, root, answerAt(2), "thread", "reply", "TH-0001", "--by", "alex", "Then the dashboard.")
	want := "TH-0001 is answered (3 entries)\n  digest: updated D-0001 (" + d.Path + ")\n"
	if code != 0 || out != want {
		t.Fatalf("second reply: %d %q %s, want %q", code, out, errOut, want)
	}
	if files := digestFiles(t, root); len(files) != 1 {
		t.Errorf("one entry, updated: %v", files)
	}
	if text := entryText(t, root, d.Path); !strings.Contains(text, "Then the dashboard.") || !strings.Contains(text, "The CLI. It is what agents use.") {
		t.Errorf("the entry does not hold both answers:\n%s", text)
	}
	view, _ = answer(t, root, 3, "thread", "reply", "TH-0001", "--by", "alex", "And the docs.")
	if view.Digest == nil || view.Digest.ID != "D-0001" || view.Digest.Action != "updated" || view.Digest.Path != d.Path {
		t.Errorf("a third answer updates D-0001: %+v", view.Digest)
	}
}

// S-0433 criterion 2: the operator's confirmation of an agent's pending
// recommendation writes an entry with the recommendation's text.
func TestThreadConfirmByTheOperatorWritesADigestEntry(t *testing.T) {
	root := answerProject(t)
	view, raw := answer(t, root, 1, "thread", "reply", "TH-0001", "--by", "orchestrator", "--recommend", "The CLI, as the plan says.")
	if _, ok := raw["digest"]; ok || view.Digest != nil {
		t.Fatalf("a recommendation writes no entry: %v", raw["digest"])
	}
	view, _ = answer(t, root, 2, "thread", "confirm", "TH-0001", "--by", "alex")
	d := view.Digest
	if view.Status != "answered" || d == nil || d.ID != "D-0001" || d.Action != "written" {
		t.Fatalf("confirm --json: %+v digest %+v", view, d)
	}
	if text := entryText(t, root, d.Path); !strings.Contains(text, "The CLI, as the plan says.") {
		t.Errorf("the entry lacks the recommendation's text:\n%s", text)
	}
}

// S-0433 criterion 1: an agent's reply, and the operator's on a thread of
// their own, write no entry and print no digest.
func TestThreadReplyByAnAgentWritesNoDigestEntry(t *testing.T) {
	root := answerProject(t)
	view, raw := answer(t, root, 1, "thread", "reply", "TH-0001", "--by", "planner", "The CLI.")
	if _, ok := raw["digest"]; ok || view.Status != "answered" {
		t.Errorf("an agent's reply: %v", raw)
	}
	if _, errOut, code := runInAt(t, root, answerAt(2), "thread", "new", "--on", "system-flow.yaml", "--by", "alex", "Mine?", "My own question."); code != 0 {
		t.Fatalf("thread new: %d %s", code, errOut)
	}
	out, errOut, code := runInAt(t, root, answerAt(3), "thread", "reply", "TH-0002", "--by", "alex", "My own answer.")
	if code != 0 || strings.Contains(out, "digest") {
		t.Errorf("the operator on their own thread: %d %q %s", code, out, errOut)
	}
	if files := digestFiles(t, root); len(files) != 0 {
		t.Errorf("no entry is written: %v", files)
	}
}

// S-0433 criterion 4: an entry that cannot be written leaves the reply saved
// and the command's exit 0, with the warning logged and under digest.warning.
func TestThreadReplyKeepsTheAnswerWhenItsDigestEntryFails(t *testing.T) {
	root := answerProject(t)
	if err := os.MkdirAll(filepath.Join(root, "design"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "design", "digest"), []byte("not a folder\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runInAt(t, root, answerAt(1), "thread", "reply", "TH-0001", "--by", "alex", "--json", "The CLI.")
	if code != 0 {
		t.Fatalf("reply: %d %s", code, errOut)
	}
	var view answered
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatalf("reply --json: %v\n%s", err, out)
	}
	if d := view.Digest; d == nil || d.ID != "" || d.Action != "" || !strings.Contains(d.Warning, "the answer is saved on the thread") {
		t.Errorf("digest: %+v", d)
	}
	if last := view.Entries[len(view.Entries)-1]; view.Status != "answered" || last.Author != "alex" || last.Text != "The CLI." {
		t.Errorf("the reply is not saved: %s", out)
	}
	if !strings.Contains(errOut, `"level":"WARN"`) || !strings.Contains(errOut, `"thread":"TH-0001"`) || !strings.Contains(errOut, `"detail":"no digest entry was written`) {
		t.Errorf("no warning logged: %s", errOut)
	}

	out, errOut, code = runInAt(t, root, answerAt(2), "thread", "reply", "TH-0001", "--by", "alex", "The dashboard.")
	if code != 0 || out != "TH-0001 is answered (3 entries)\n" || !strings.Contains(errOut, `"level":"WARN"`) {
		t.Errorf("text reply: %d %q %s", code, out, errOut)
	}
}

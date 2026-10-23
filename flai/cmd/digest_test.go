package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/digest"
)

// digestProject is a project whose sources are a story, a thread, and an
// issue, with no digest folder yet.
func digestProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for path, body := range map[string]string{
		"wip/kanban/stories/S-0001-a-story.md": "---\nid: S-0001\n---\n",
		"wip/threads/TH-0002.md":               "---\nid: TH-0002\n---\n",
		"wip/threads/TH-0003.md":               "---\nid: TH-0003\n---\n",
		"design/issues/I-0004-an-issue.md":     "---\nid: I-0004\n---\n",
	} {
		file := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// digestFiles are the names in the project's digest folder.
func digestFiles(t *testing.T, root string) []string {
	t.Helper()
	files, err := os.ReadDir(filepath.Join(root, "design", "digest"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, f := range files {
		out = append(out, f.Name())
	}
	return out
}

// ADR-0191, S-0428 criterion 3: new allocates the ID, writes the entry in the
// checkout it runs in, and prints it; the next one takes the next number.
func TestDigestNew(t *testing.T) {
	root := digestProject(t)
	out, errOut, code := runIn(t, root, "digest", "new", "The digest lives under design/digest", "--kind", "answer", "--source", "TH-0002,S-0001", "--topic", "digest", "--path", "design/digest", "--body", "The operator chose one file per entry.")
	if code != 0 {
		t.Fatalf("digest new: %d %s %s", code, out, errOut)
	}
	if !strings.HasPrefix(out, "D-0001 The digest lives under design/digest (answer, active)\n  design/digest/D-0001-the-digest-lives-under-design-digest.md\n") {
		t.Errorf("text output:\n%s", out)
	}
	data, err := os.ReadFile(filepath.Join(root, "design", "digest", "D-0001-the-digest-lives-under-design-digest.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"sources: [TH-0002, S-0001]\n", "topics: [digest]\n", "paths: [design/digest]\n", "date: 2026-09-15\n", "status: active\n", "# D-0001 The digest lives under design/digest\n\nThe operator chose one file per entry.\n"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("missing %q in:\n%s", want, data)
		}
	}

	out, errOut, code = runIn(t, root, "digest", "new", "Issues name their cause", "--kind", "cause", "--source", "I-0004", "--json")
	if code != 0 {
		t.Fatalf("digest new --json: %d %s %s", code, out, errOut)
	}
	var res digest.Result
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if res.ID != "D-0002" || res.Path != "design/digest/D-0002-issues-name-their-cause.md" || res.Entry.Kind != "cause" || len(res.Entry.Sources) != 1 || res.Entry.Topics == nil {
		t.Errorf("the second entry takes the next number: %+v", res)
	}
}

// --body-file reads the body from a file, or standard input with -.
func TestDigestNewBodyFile(t *testing.T) {
	root := digestProject(t)
	if err := os.WriteFile(filepath.Join(root, "body.md"), []byte("From a file.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, errOut, code := runIn(t, root, "digest", "new", "From a file", "--kind", "decision", "--source", "S-0001", "--body-file", "body.md"); code != 0 {
		t.Fatalf("--body-file: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runStdin(t, root, "From standard input.\n", "digest", "new", "From stdin", "--kind", "decision", "--source", "S-0001", "--body-file", "-"); code != 0 {
		t.Fatalf("--body-file -: %d %s %s", code, out, errOut)
	}
	for id, want := range map[string]string{"D-0001": "From a file.", "D-0002": "From standard input."} {
		out, _, _ := runIn(t, root, "digest", "show", id, "--json")
		var e digest.Entry
		if err := json.Unmarshal([]byte(out), &e); err != nil || e.Body != want {
			t.Errorf("%s body %q, want %q (%v)", id, e.Body, want, err)
		}
	}
}

// new refuses a missing --kind, an unknown kind, and no --source before it
// writes anything, saying what to give.
func TestDigestNewNeedsKindAndSource(t *testing.T) {
	root := digestProject(t)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--source", "S-0001"}, "--kind is required"},
		{[]string{"--kind", "guess", "--source", "S-0001"}, "none of decision, answer, recommendation, cause"},
		{[]string{"--kind", "decision"}, "--source is required"},
	} {
		args := append([]string{"digest", "new", "A fact"}, c.args...)
		_, errOut, code := runIn(t, root, args...)
		if code != 1 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, want 1 and %q in:\n%s", c.args, code, c.want, errOut)
		}
	}
	if files := digestFiles(t, root); files != nil {
		t.Errorf("nothing is written: %v", files)
	}
}

// S-0428 criterion 3: a write the check refuses prints each problem, exits 4,
// and leaves no file, nor the folder it made.
func TestDigestNewRefused(t *testing.T) {
	root := digestProject(t)
	out, errOut, code := runIn(t, root, "digest", "new", "A fact with a key", "--kind", "decision", "--source", "S-0001", "--body", "The key is AKIAIOSFODNN7EXAMPLE.")
	if code != exitDocRefused {
		t.Fatalf("a secret: code %d, want %d\n%s%s", code, exitDocRefused, out, errOut)
	}
	if !strings.Contains(out, "design/digest/D-0001-a-fact-with-a-key.md:") || !strings.Contains(out, ": error: digest.secret: ") || strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("each problem is printed with its line and rule, never the value:\n%s", out)
	}
	if files := digestFiles(t, root); files != nil {
		t.Errorf("nothing is written: %v", files)
	}

	out, _, code = runIn(t, root, "digest", "new", "A fact from nowhere", "--kind", "decision", "--source", "S-0099", "--json")
	if code != exitDocRefused {
		t.Fatalf("an unknown source: code %d, want %d\n%s", code, exitDocRefused, out)
	}
	var res struct {
		Refused struct {
			Path     string
			Reason   string
			Findings []struct {
				Level, Rule, Path, Message string
				Line                       int
			}
		}
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if len(res.Refused.Findings) != 1 || res.Refused.Findings[0].Rule != digest.RuleEntry || !strings.Contains(res.Refused.Findings[0].Message, "S-0099") || res.Refused.Findings[0].Line < 2 {
		t.Errorf("the refusal names the unknown source: %+v", res)
	}
	if files := digestFiles(t, root); files != nil {
		t.Errorf("nothing is written: %v", files)
	}
}

// S-0428 criterion 3: edit replaces the lists it is given, empties one given
// empty, and leaves the rest of the file as it was.
func TestDigestEdit(t *testing.T) {
	root := digestProject(t)
	if out, errOut, code := runIn(t, root, "digest", "new", "A fact", "--kind", "decision", "--source", "S-0001", "--topic", "cli", "--path", "flai/cmd", "--body", "Why."); code != 0 {
		t.Fatalf("digest new: %d %s %s", code, out, errOut)
	}
	out, errOut, code := runIn(t, root, "digest", "edit", "D-1", "--topic", "digest", "--topic", "go", "--asked-again", "TH-0002,TH-0003", "--path=", "--json")
	if code != 0 {
		t.Fatalf("digest edit: %d %s %s", code, out, errOut)
	}
	var res digest.Result
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	e := res.Entry
	if res.ID != "D-0001" || strings.Join(e.Topics, ",") != "digest,go" || strings.Join(e.AskedAgain, ",") != "TH-0002,TH-0003" || len(e.Paths) != 0 || strings.Join(e.Sources, ",") != "S-0001" || e.Body != "Why." {
		t.Errorf("topics and asked_again replaced, paths emptied, the rest kept: %+v", res)
	}

	out, errOut, code = runIn(t, root, "digest", "edit", "D-0001", "--body", "Because.")
	if code != 0 || !strings.HasPrefix(out, "D-0001 A fact (decision, active)\n") || !strings.Contains(out, "edited") {
		t.Fatalf("digest edit text: %d %s %s", code, out, errOut)
	}
	data, _ := os.ReadFile(filepath.Join(root, "design", "digest", "D-0001-a-fact.md"))
	if !strings.Contains(string(data), "asked_again: [TH-0002, TH-0003]\n") || !strings.HasSuffix(string(data), "# D-0001 A fact\n\nBecause.\n") {
		t.Errorf("the body is replaced and asked_again kept:\n%s", data)
	}
}

// edit with nothing to change is refused naming the flags; one the check
// refuses puts the file back and exits 4.
func TestDigestEditRefused(t *testing.T) {
	root := digestProject(t)
	if out, errOut, code := runIn(t, root, "digest", "new", "A fact", "--kind", "decision", "--source", "S-0001"); code != 0 {
		t.Fatalf("digest new: %d %s %s", code, out, errOut)
	}
	file := filepath.Join(root, "design", "digest", "D-0001-a-fact.md")
	before, _ := os.ReadFile(file)

	_, errOut, code := runIn(t, root, "digest", "edit", "D-0001")
	if code != 1 || !strings.Contains(errOut, "nothing to change: give --body or --body-file, --topic, --path, --source, or --asked-again") {
		t.Errorf("no change flag: code %d\n%s", code, errOut)
	}
	out, _, code := runIn(t, root, "digest", "edit", "D-0001", "--asked-again", "S-0001")
	if code != exitDocRefused || !strings.Contains(out, "digest.entry") || !strings.Contains(out, "not a thread ID") {
		t.Errorf("asked_again holding a story: code %d\n%s", code, out)
	}
	after, _ := os.ReadFile(file)
	if string(after) != string(before) {
		t.Errorf("the file is put back:\n%s", after)
	}
}

// S-0428 criterion 2: list prints each entry with its status, kind, date,
// sources, topics, and paths, filtered by --status and --topic; --json is
// an array, [] when empty.
func TestDigestList(t *testing.T) {
	root := digestProject(t)
	out, _, code := runIn(t, root, "digest", "list", "--json")
	if code != 0 || strings.TrimSpace(out) != "[]" {
		t.Errorf("no folder is no entries, []: %d %q", code, out)
	}
	if out, _, _ := runIn(t, root, "digest", "list"); out != "no digest entries\n" {
		t.Errorf("text when empty: %q", out)
	}
	for _, args := range [][]string{
		{"First fact", "--kind", "decision", "--source", "S-0001", "--topic", "cli", "--path", "flai/cmd"},
		{"Second fact", "--kind", "answer", "--source", "TH-0002", "--topic", "digest"},
	} {
		if out, errOut, code := runIn(t, root, append([]string{"digest", "new"}, args...)...); code != 0 {
			t.Fatalf("digest new: %d %s %s", code, out, errOut)
		}
	}
	// An absorbed entry, as the curator will leave one.
	file := filepath.Join(root, "design", "digest", "D-0002-second-fact.md")
	data, _ := os.ReadFile(file)
	absorbed := strings.NewReplacer("status: active\n", "status: absorbed\n", "absorbed_by: \"\"\n", "absorbed_by: wip/threads/TH-0002.md\n").Replace(string(data))
	if !strings.Contains(absorbed, "absorbed_by: wip/") {
		t.Fatalf("the fixture did not set absorbed_by:\n%s", absorbed)
	}
	if err := os.WriteFile(file, []byte(absorbed), 0o644); err != nil {
		t.Fatal(err)
	}

	out, errOut, code := runIn(t, root, "digest", "list")
	if code != 0 {
		t.Fatalf("digest list: %d %s", code, errOut)
	}
	want := "D-0001 active    decision       2026-09-15 First fact\n       sources: S-0001  topics: cli  paths: flai/cmd\n" +
		"D-0002 absorbed  answer         2026-09-15 Second fact\n       sources: TH-0002  topics: digest  paths: -\n"
	if out != want {
		t.Errorf("list text:\n%s\nwant:\n%s", out, want)
	}

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--status", "absorbed"}, "D-0002"},
		{[]string{"--status", "active"}, "D-0001"},
		{[]string{"--topic", "cli"}, "D-0001"},
		{[]string{"--topic", "digest", "--status", "active"}, ""},
	} {
		out, errOut, code := runIn(t, root, append([]string{"digest", "list", "--json"}, c.args...)...)
		var got []digest.Entry
		if err := json.Unmarshal([]byte(out), &got); err != nil || code != 0 {
			t.Fatalf("%v: %d %v %s %s", c.args, code, err, out, errOut)
		}
		var ids []string
		for _, e := range got {
			ids = append(ids, e.ID)
		}
		if strings.Join(ids, ",") != c.want {
			t.Errorf("%v: got %v, want %q", c.args, ids, c.want)
		}
	}

	_, errOut, code = runIn(t, root, "digest", "list", "--status", "stale")
	if code != 1 || !strings.Contains(errOut, "none of active, absorbed, compacted") {
		t.Errorf("an unknown status is refused naming the statuses: %d %s", code, errOut)
	}
}

// S-0428 criterion 2: show prints the entry's fields and body; --json the
// entry; an unknown ID is an error naming it and flai digest list.
func TestDigestShow(t *testing.T) {
	root := digestProject(t)
	if out, errOut, code := runIn(t, root, "digest", "new", "A fact", "--kind", "recommendation", "--source", "TH-0002", "--source", "S-0001", "--topic", "cli", "--body", "Line one.\n\nLine two."); code != 0 {
		t.Fatalf("digest new: %d %s %s", code, out, errOut)
	}
	out, errOut, code := runIn(t, root, "digest", "show", "d-1")
	if code != 0 {
		t.Fatalf("digest show: %d %s", code, errOut)
	}
	want := "D-0001 A fact\n" +
		"  kind:        recommendation\n" +
		"  status:      active\n" +
		"  date:        2026-09-15\n" +
		"  sources:     TH-0002, S-0001\n" +
		"  topics:      cli\n" +
		"  paths:       -\n" +
		"  asked_again: -\n" +
		"  cited:       -\n" +
		"  path:        design/digest/D-0001-a-fact.md\n" +
		"\nLine one.\n\nLine two.\n"
	if out != want {
		t.Errorf("show text:\n%s\nwant:\n%s", out, want)
	}

	out, _, code = runIn(t, root, "digest", "show", "D-0001", "--json")
	var e digest.Entry
	if err := json.Unmarshal([]byte(out), &e); err != nil || code != 0 {
		t.Fatalf("show --json: %d %v %s", code, err, out)
	}
	if e.ID != "D-0001" || e.Kind != "recommendation" || e.Path != "design/digest/D-0001-a-fact.md" || e.Body != "Line one.\n\nLine two." || e.Cited == nil {
		t.Errorf("show --json: %+v", e)
	}

	for _, args := range [][]string{{"digest", "show", "D-0042"}, {"digest", "edit", "D-0042", "--topic", "cli"}} {
		_, errOut, code := runIn(t, root, args...)
		if code != 1 || !strings.Contains(errOut, "D-0042") || !strings.Contains(errOut, "flai digest list") {
			t.Errorf("%v: an unknown ID is named with flai digest list: %d %s", args, code, errOut)
		}
	}
}

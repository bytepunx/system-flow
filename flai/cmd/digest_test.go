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

// seedProject is a project owned by alex with one candidate source of each
// kind, learned on different days: TH-0001 with two answers by alex on
// 2026-10-19, the archived S-0001 done on 2026-10-20 with one decision,
// I-0001 closed on 2026-10-21 with a cause and a fix, and TH-0002 with alex
// overruling a recommendation on 2026-10-22.
func seedProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	thread := func(id, title, status string, entries ...string) string {
		var b strings.Builder
		b.WriteString("---\nid: " + id + "\ntitle: " + title + "\nanchor:\n  path: wip/archive/kanban/stories/S-0001-s.md\n  item: S-0001\nstatus: " + status + "\nparticipants: [claude, alex]\ncreated: 2026-10-19T10:00:00Z\nupdated: 2026-10-23T10:00:00Z\n---\n\n# " + id + " " + title + "\n\n## Entries\n")
		for i := 0; i+1 < len(entries); i += 2 {
			b.WriteString("\n### " + entries[i] + "\n" + entries[i+1] + "\n")
		}
		return b.String()
	}
	for path, body := range map[string]string{
		"system-flow.yaml":                       "version: 1\nname: t\nkey: t\nowner: alex\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n",
		"wip/archive/kanban/stories/S-0001-s.md": "---\nid: S-0001\ntype: story\nnature: feature\ntitle: S-0001 title\nstatus: done\nparent: \nowner: alex\ncreated: 2026-08-01T00:00:00Z\nupdated: 2026-08-01T00:00:00Z\ntransitions:\n  - to: in-progress\n    at: 2026-08-01T00:00:00Z\n    by: alex\n  - to: done\n    at: 2026-10-20T12:00:00Z\n    by: alex\n---\n\n# S-0001 S-0001 title\n",
		"wip/archive/agents/S-0001.md":           "---\nstream: S-0001\ntitle: t\nupdated: 2026-08-01T00:00:00Z\nagent: a\nsession: s\n---\n\n# S-0001 t\n\n## Context\n\nWhy.\n\n## Decisions\n\n- Keep one file per entry. It merges cleanly.\n\n## Log\n",
		"wip/agents/index.md":                    "# Narratives\n",
		"design/issues/I-0001-a.md":              "---\nid: I-0001\ntitle: I-0001 went wrong\nclass: defect\nstatus: closed\ncount: 1\nfirst_reported: 2026-08-01T00:00:00Z\nlast_reported: 2026-08-01T00:00:00Z\nupdated: 2026-08-02T00:00:00Z\n---\n\n# I-0001 went wrong\n\n## Description\nTwo writers took one number.\n\n## Remediation\nClosed 2026-10-21T09:00:00Z: S-0001: numbers come from every branch.\n",
		"wip/threads/TH-0001-shape.md": thread("TH-0001", "Which shape?", "answered",
			"2026-10-19T10:00:00Z claude", "Left or right?",
			"2026-10-19T10:05:00Z alex", "Left. It is cheaper.",
			"2026-10-19T10:10:00Z alex", "Narrow, too."),
		"wip/threads/TH-0002-store.md": thread("TH-0002", "Which store?", "resolved",
			"2026-10-22T08:00:00Z claude", "Files or a database?",
			"2026-10-22T08:01:00Z orchestrator (recommendation)", "A database.",
			"2026-10-22T08:02:00Z alex", "Files: git keeps them."),
		"wip/threads/TH-0003-key.md": thread("TH-0003", "Which key?", "answered",
			"2026-10-20T09:00:00Z claude", "Which key does the deploy use?",
			"2026-10-20T09:01:00Z alex", "The key is AKIAIOSFODNN7EXAMPLE."),
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

// harvested is each candidate's kind, source, and time, one string each.
func harvested(list []digest.Candidate) []string {
	out := []string{}
	for _, c := range list {
		out = append(out, c.Kind+" "+c.Source+" "+c.At)
	}
	return out
}

// S-0429 criterion 1: harvest lists every candidate with its kind, source,
// where in it, date, title, and the entries citing it, then a count by kind;
// --json is the candidates with every field.
func TestDigestHarvest(t *testing.T) {
	root := seedProject(t)
	if err := os.Remove(filepath.Join(root, "wip", "threads", "TH-0003-key.md")); err != nil {
		t.Fatal(err)
	}
	if out, errOut, code := runIn(t, root, "digest", "new", "Left is cheaper", "--kind", "answer", "--source", "I-0001"); code != 0 {
		t.Fatalf("digest new: %d %s %s", code, out, errOut)
	}
	out, errOut, code := runIn(t, root, "digest", "harvest")
	if code != 0 {
		t.Fatalf("digest harvest: %d %s", code, errOut)
	}
	want := "answer         2026-10-19 TH-0001 2026-10-19T10:05:00Z: Which shape: Left.\n" +
		"answer         2026-10-19 TH-0001 2026-10-19T10:10:00Z: Which shape: Narrow, too.\n" +
		"decision       2026-10-20 S-0001 wip/archive/agents/S-0001.md: Keep one file per entry.\n" +
		"cause          2026-10-21 I-0001: I-0001 went wrong  (cited by D-0001)\n" +
		"recommendation 2026-10-22 TH-0002 2026-10-22T08:02:00Z: Which store: Files: git keeps them.\n" +
		"5 candidates, 4 cited by no entry: 1 decision, 2 answer, 1 recommendation, 1 cause\n"
	if out != want {
		t.Errorf("harvest text:\n%s\nwant:\n%s", out, want)
	}
	if files := digestFiles(t, root); len(files) != 1 {
		t.Errorf("harvest writes nothing: %v", files)
	}

	out, _, code = runIn(t, root, "digest", "harvest", "--json")
	var got []digest.Candidate
	if err := json.Unmarshal([]byte(out), &got); err != nil || code != 0 {
		t.Fatalf("harvest --json: %d %v %s", code, err, out)
	}
	if len(got) != 5 {
		t.Fatalf("harvest --json: %+v", got)
	}
	if c := got[4]; c.Kind != "recommendation" || c.Source != "TH-0002" || c.Ref != "2026-10-22T08:02:00Z" || c.Date != "2026-10-22" || c.Text != "Files: git keeps them." || c.Overruled != "A database." || c.Question != "Files or a database?" || c.Topics == nil || c.Paths == nil {
		t.Errorf("the recommendation, every field: %+v", c)
	}
	if c := got[3]; strings.Join(c.CitedBy, ",") != "D-0001" || c.Fix != "S-0001: numbers come from every branch." {
		t.Errorf("the cause, cited by D-0001: %+v", c)
	}
}

// S-0429 criterion 1: --since keeps what was learned on or after it, as a
// day or a time; --kind keeps its kinds, repeated or comma separated; a
// --since or --kind with nothing to list is [] and "no candidates".
func TestDigestHarvestSinceAndKind(t *testing.T) {
	root := seedProject(t)
	for _, c := range []struct {
		args []string
		want []string
	}{
		{[]string{"--since", "2026-10-21"}, []string{"cause I-0001 2026-10-21T09:00:00Z", "recommendation TH-0002 2026-10-22T08:02:00Z"}},
		{[]string{"--since", "2026-10-20T12:00:00Z"}, []string{"decision S-0001 2026-10-20T12:00:00Z", "cause I-0001 2026-10-21T09:00:00Z", "recommendation TH-0002 2026-10-22T08:02:00Z"}},
		{[]string{"--kind", "cause", "--kind", "recommendation"}, []string{"cause I-0001 2026-10-21T09:00:00Z", "recommendation TH-0002 2026-10-22T08:02:00Z"}},
		{[]string{"--kind", "answer,decision", "--since", "2026-10-20"}, []string{"answer TH-0003 2026-10-20T09:01:00Z", "decision S-0001 2026-10-20T12:00:00Z"}},
		{[]string{"--since", "2026-11-01"}, []string{}},
	} {
		out, errOut, code := runIn(t, root, append([]string{"digest", "harvest", "--json"}, c.args...)...)
		var got []digest.Candidate
		if err := json.Unmarshal([]byte(out), &got); err != nil || code != 0 {
			t.Fatalf("%v: %d %v %s %s", c.args, code, err, out, errOut)
		}
		if strings.Join(harvested(got), "|") != strings.Join(c.want, "|") {
			t.Errorf("%v: %v\nwant %v", c.args, harvested(got), c.want)
		}
	}
	if out, _, _ := runIn(t, root, "digest", "harvest", "--since", "2026-11-01", "--json"); strings.TrimSpace(out) != "[]" {
		t.Errorf("nothing to list is [], not null: %q", out)
	}
	if out, _, _ := runIn(t, root, "digest", "harvest", "--since", "2026-11-01"); out != "no candidates\n" {
		t.Errorf("text when empty: %q", out)
	}
}

// harvest and seed refuse a --since that is not a date, and harvest a
// --kind that is none of the kinds, each saying what to give.
func TestDigestHarvestRefusesABadSinceOrKind(t *testing.T) {
	root := seedProject(t)
	for _, c := range []struct {
		args []string
		want []string
	}{
		{[]string{"digest", "harvest", "--kind", "answer,rumour"}, []string{"rumour", "is none of decision, answer, recommendation, cause; give one of them"}},
		{[]string{"digest", "harvest", "--since", "last week"}, []string{"last week", "give a day, YYYY-MM-DD in UTC"}},
		{[]string{"digest", "seed", "--since", "2026-13-01"}, []string{"2026-13-01", "give a day, YYYY-MM-DD in UTC"}},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		for _, want := range c.want {
			if code != 1 || !strings.Contains(errOut, want) {
				t.Errorf("%v: code %d, want 1 and %q in:\n%s", c.args, code, want, errOut)
			}
		}
	}
	if files := digestFiles(t, root); files != nil {
		t.Errorf("nothing is written: %v", files)
	}
}

// S-0429 criterion 3: seed writes an entry for each answer, overruled
// recommendation, and cause, never a decision; the secret-shaped answer the
// check refuses is skipped and named by its rule, never its text, while the
// rest are written; a second run writes nothing.
func TestDigestSeed(t *testing.T) {
	root := seedProject(t)
	out, errOut, code := runIn(t, root, "digest", "seed")
	if code != 0 {
		t.Fatalf("digest seed: %d %s %s", code, out, errOut)
	}
	want := "D-0001 answer TH-0001 Which shape: Left\n" +
		"D-0002 answer TH-0001 Which shape: Narrow, too\n" +
		"D-0003 cause I-0001 I-0001 went wrong\n" +
		"D-0004 recommendation TH-0002 Which store: Files: git keeps them\n" +
		"skipped answer TH-0003 2026-10-20T09:01:00Z: refused: digest.secret\n" +
		"answer:         listed 3, written 2, skipped 1\n" +
		"recommendation: listed 1, written 1, skipped 0\n" +
		"cause:          listed 1, written 1, skipped 0\n"
	if out != want {
		t.Errorf("seed text:\n%s\nwant:\n%s", out, want)
	}
	files := digestFiles(t, root)
	if len(files) != 4 {
		t.Fatalf("four entries are written: %v", files)
	}
	for _, f := range files {
		data, _ := os.ReadFile(filepath.Join(root, "design", "digest", f))
		if strings.Contains(string(data), "AKIA") {
			t.Errorf("%s holds the key", f)
		}
	}

	out, errOut, code = runIn(t, root, "digest", "seed", "--json")
	if code != 0 {
		t.Fatalf("second digest seed: %d %s %s", code, out, errOut)
	}
	var res digest.SeedResult
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	var reasons []string
	for _, s := range res.Skipped {
		reasons = append(reasons, s.Source+" "+s.Reason)
	}
	wantReasons := "TH-0001 cited by D-0001|TH-0001 cited by D-0002|TH-0003 refused: digest.secret|I-0001 cited by D-0003|TH-0002 cited by D-0004"
	if res.Written == nil || len(res.Written) != 0 || res.DryRun || strings.Join(reasons, "|") != wantReasons || res.Listed["answer"] != 3 || res.Listed["cause"] != 1 {
		t.Errorf("a second run writes nothing: %+v\nreasons %q", res, reasons)
	}
	if files := digestFiles(t, root); len(files) != 4 {
		t.Errorf("still four entries: %v", files)
	}
}

// S-0429 criterion 3: a dry run prints what it would write and writes
// nothing; --json gives the result with no IDs, and an answer with a key's
// shape is skipped as the real run skips it.
func TestDigestSeedDryRun(t *testing.T) {
	root := seedProject(t)
	out, errOut, code := runIn(t, root, "digest", "seed", "--dry-run", "--since", "2026-10-21")
	if code != 0 {
		t.Fatalf("digest seed --dry-run: %d %s %s", code, out, errOut)
	}
	want := "would write cause I-0001 I-0001 went wrong\n" +
		"would write recommendation TH-0002 Which store: Files: git keeps them\n" +
		"answer:         listed 0, would write 0, skipped 0\n" +
		"recommendation: listed 1, would write 1, skipped 0\n" +
		"cause:          listed 1, would write 1, skipped 0\n"
	if out != want {
		t.Errorf("dry run text:\n%s\nwant:\n%s", out, want)
	}

	out, _, code = runIn(t, root, "digest", "seed", "--dry-run", "--json")
	var res digest.SeedResult
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("dry run --json: %d %v %s", code, err, out)
	}
	refused := len(res.Skipped) == 1 && res.Skipped[0].Source == "TH-0003" && res.Skipped[0].Reason == "refused: digest.secret"
	if !res.DryRun || len(res.Written) != 4 || res.Written[0].ID != "" || res.Written[0].Path != "" || res.Written[0].Kind != "answer" || !refused {
		t.Errorf("dry run --json: %+v", res)
	}
	if strings.Contains(out, "AKIA") {
		t.Errorf("a dry run printed the key: %s", out)
	}
	if files := digestFiles(t, root); files != nil {
		t.Errorf("a dry run writes nothing: %v", files)
	}
}

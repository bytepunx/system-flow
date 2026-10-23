package digest

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/threads"
)

// thread parses a thread titled title, anchored by anchor (the lines under
// "anchor:"), whose entries are given as "<ts> <author>" headings each
// followed by its text.
func thread(t *testing.T, title, anchor string, entries ...string) *threads.Thread {
	t.Helper()
	var b strings.Builder
	b.WriteString("---\nid: TH-0042\ntitle: " + title + "\nanchor:\n" + anchor + "status: answered\nparticipants: [claude, alex]\ncreated: 2026-09-17T21:00:00Z\nupdated: 2026-09-17T21:05:00Z\n---\n\n# TH-0042 " + title + "\n\n## Entries\n")
	for i := 0; i+1 < len(entries); i += 2 {
		b.WriteString("\n### " + entries[i] + "\n" + entries[i+1] + "\n")
	}
	th, err := threads.Parse(b.String())
	if err != nil {
		t.Fatal(err)
	}
	return th
}

const pathAnchor = "  path: design/system/plan.md\n"

func TestAnswersGivesTheOperatorsAnswers(t *testing.T) {
	cases := []struct {
		name string
		th   *threads.Thread
		want []Candidate
	}{
		{
			name: "a plain owner reply is an answer",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z alex", "Left. The plan says so.\n\nAnd it is cheaper."),
			want: []Candidate{{Kind: "answer", Ref: "2026-09-17T21:01:00Z", Title: "Which shape: Left.", Text: "Left. The plan says so.\n\nAnd it is cheaper."}},
		},
		{
			name: "a confirmation answers with the recommendation's text",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z orchestrator (recommendation)", "**Right**, as the plan says.\n\nSource: design/system/plan.md § Shape",
				"2026-09-17T21:02:00Z alex", "Confirmed the recommendation of 2026-09-17T21:01:00Z orchestrator.\n\nSource: design/system/plan.md § Shape"),
			want: []Candidate{{Kind: "answer", Ref: "2026-09-17T21:02:00Z", Title: "Which shape: Right, as the plan says.", Text: "**Right**, as the plan says.\n\nSource: design/system/plan.md § Shape"}},
		},
		{
			name: "a confirmation naming no recommendation takes the latest before it",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z orchestrator (recommendation)", "Right.",
				"2026-09-17T21:02:00Z alex", "Confirmed the recommendation of 2026-09-17T20:00:00Z nobody."),
			want: []Candidate{{Kind: "answer", Ref: "2026-09-17T21:02:00Z", Title: "Which shape: Right.", Text: "Right."}},
		},
		{
			name: "an owner reply to a pending recommendation overrules it",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z orchestrator (recommendation)", "Right, as the plan says.",
				"2026-09-17T21:02:00Z alex", "- Left: right is too dear.",
				"2026-09-17T21:03:00Z alex", "Resolved."),
			want: []Candidate{{Kind: "recommendation", Ref: "2026-09-17T21:02:00Z", Title: "Which shape: Left: right is too dear.", Text: "- Left: right is too dear.", Overruled: "Right, as the plan says."}},
		},
		{
			name: "a recommendation the owner already answered is not overruled by a later reply",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z orchestrator (recommendation)", "Right.",
				"2026-09-17T21:02:00Z alex", "Left.",
				"2026-09-17T21:03:00Z claude", "Left it is.",
				"2026-09-17T21:04:00Z alex", "Narrow, too."),
			want: []Candidate{
				{Kind: "recommendation", Ref: "2026-09-17T21:02:00Z", Title: "Which shape: Left.", Text: "Left.", Overruled: "Right."},
				{Kind: "answer", Ref: "2026-09-17T21:04:00Z", Title: "Which shape: Narrow, too.", Text: "Narrow, too."},
			},
		},
		{
			name: "a reply and its resolution in the same second keep the reply only",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z alex", "Left.\n\nResolved: S-0042 was accepted"),
			want: []Candidate{{Kind: "answer", Ref: "2026-09-17T21:01:00Z", Title: "Which shape: Left.", Text: "Left."}},
		},
		{
			name: "a thread resolved with no answer gives nothing",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z alex", "Resolved.",
				"2026-09-17T21:02:00Z alex", "Resolved: S-0042 was accepted"),
		},
		{
			name: "a thread whose only replies are agents' gives nothing",
			th: thread(t, "Which shape?", pathAnchor,
				"2026-09-17T21:00:00Z claude", "Left or right?",
				"2026-09-17T21:01:00Z orchestrator (recommendation)", "Right.",
				"2026-09-17T21:02:00Z reviewer", "Confirmed the recommendation of 2026-09-17T21:01:00Z orchestrator.",
				"2026-09-17T21:03:00Z claude", "Resolved."),
		},
		{
			name: "the owner's own thread gives nothing",
			th: thread(t, "Is the shape right?", pathAnchor,
				"2026-09-17T21:00:00Z alex", "I wonder.",
				"2026-09-17T21:01:00Z claude", "Yes.",
				"2026-09-17T21:02:00Z alex", "Then go on."),
		},
		{
			name: "a permission-prompt thread gives nothing",
			th: thread(t, "Allow Write .claude/settings.json?", "  path: wip/kanban/stories/S-0042-x.md\n  item: S-0042\n",
				"2026-09-17T21:00:00Z claude", "May I write .claude/settings.json?",
				"2026-09-17T21:01:00Z alex", "yes"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Answers(tc.th, "alex", nil)
			if len(got) != len(tc.want) {
				t.Fatalf("got %d candidates, want %d: %+v", len(got), len(tc.want), got)
			}
			for i, w := range tc.want {
				g := got[i]
				if g.Kind != w.Kind || g.Ref != w.Ref || g.At != w.Ref || g.Title != w.Title || g.Text != w.Text || g.Overruled != w.Overruled {
					t.Errorf("candidate %d:\n got %+v\nwant %+v", i, g, w)
				}
				if g.Source != "TH-0042" || g.Date != "2026-09-17" || g.Question != "Left or right?" {
					t.Errorf("candidate %d: source %q, date %q, question %q", i, g.Source, g.Date, g.Question)
				}
			}
		})
	}
}

func TestAnswersTakesTopicsAndPathsFromTheAnchoredItem(t *testing.T) {
	asked := ""
	lookup := func(id string) ([]string, []string) {
		asked = id
		return []string{"digest"}, []string{"flai/internal/digest/"}
	}
	onItem := thread(t, "Which shape?", "  path: wip/kanban/stories/S-0042-x.md\n  item: S-0042\n",
		"2026-09-17T21:00:00Z claude", "Left or right?",
		"2026-09-17T21:01:00Z alex", "Left.")
	got := Answers(onItem, "alex", lookup)
	if len(got) != 1 || asked != "S-0042" || strings.Join(got[0].Topics, ",") != "digest" || strings.Join(got[0].Paths, ",") != "flai/internal/digest/" {
		t.Fatalf("item anchor: asked %q, got %+v", asked, got)
	}

	asked = ""
	onPath := thread(t, "Which shape?", pathAnchor,
		"2026-09-17T21:00:00Z claude", "Left or right?",
		"2026-09-17T21:01:00Z alex", "Left.")
	got = Answers(onPath, "alex", lookup)
	if len(got) != 1 || asked != "" {
		t.Fatalf("path anchor: asked %q, got %+v", asked, got)
	}
	data, err := json.Marshal(got[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"topics":[]`, `"paths":[]`, `"cited_by":[]`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("path anchor JSON lacks %s: %s", want, data)
		}
	}
}

func TestAnswerTitleIsCutOnAWord(t *testing.T) {
	cases := []struct{ title, answer, want string }{
		{"How should X be done?", "yes", "How should X be done: yes"},
		{"Which shape.", "1. **Left**  first.\nThen right.", "Which shape: Left first."},
		{"Which shape?", "Left? Or right.", "Which shape: Left?"},
	}
	for _, tc := range cases {
		if got := answerTitle(tc.title, tc.answer); got != tc.want {
			t.Errorf("answerTitle(%q, %q) = %q, want %q", tc.title, tc.answer, got, tc.want)
		}
	}
	long := answerTitle("Which shape?", strings.Repeat("word ", 60))
	if r := []rune(long); len(r) > titleLimit || !strings.HasSuffix(long, "word…") {
		t.Errorf("long title is %d runes: %q", len(r), long)
	}
}

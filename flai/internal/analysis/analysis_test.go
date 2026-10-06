package analysis

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

var m = manifest.Manifest{Layout: map[string]string{"design": "notes", "docs": "docs", "wip": "wip"}}

// S-0223: a report is <date>-<focus>.md in the design folder's analysis
// folder, and a run asked for no focus writes <date>-all.md.
func TestReportPathIsTheDateAndTheFocus(t *testing.T) {
	if got := Dir(m); got != "notes/analysis" {
		t.Errorf("dir: %q", got)
	}
	for focus, want := range map[string]string{
		"bottlenecks": "2026-10-06-bottlenecks.md",
		"intent":      "2026-10-06-intent.md",
		"risk":        "2026-10-06-risk.md",
		"":            "2026-10-06-all.md",
		AllFocus:      "2026-10-06-all.md",
	} {
		if got := ReportPath("2026-10-06", focus); got != want {
			t.Errorf("ReportPath(%q) = %q, want %q", focus, got, want)
		}
	}
}

const good = "---\ntitle: Bottlenecks in September\nupdated: 2026-10-06T08:00:00Z\nstatus: active\nfocus: bottlenecks\nfrom: 2026-09-01\nto: 2026-09-30\n---\n\n# Bottlenecks in September\n"

// S-0223: a report with every field there and well formed has no problems,
// whether updated is a timestamp or a date.
func TestValidateAcceptsAWellFormedReport(t *testing.T) {
	f, problems := Validate("2026-10-06-bottlenecks.md", good)
	if len(problems) != 0 {
		t.Errorf("problems: %+v", problems)
	}
	if f.Focus != "bottlenecks" || f.From != "2026-09-01" || f.To != "2026-09-30" {
		t.Errorf("front matter: %+v", f)
	}
	dated := strings.Replace(good, "2026-10-06T08:00:00Z", "2026-10-06", 1)
	if _, problems := Validate("2026-10-06-bottlenecks.md", dated); len(problems) != 0 {
		t.Errorf("updated as a date: %+v", problems)
	}
	all := strings.Replace(good, "focus: bottlenecks", "focus: all", 1)
	if _, problems := Validate("2026-10-06-all.md", all); len(problems) != 0 {
		t.Errorf("a report of all three: %+v", problems)
	}
}

// S-0223: Validate names each field that is missing or bad, and a file not
// named <date>-<focus>.md.
func TestValidateNamesEachMissingOrBadField(t *testing.T) {
	for _, field := range []string{"title", "updated", "status", "focus", "from", "to"} {
		var kept []string
		for _, l := range strings.Split(good, "\n") {
			if !strings.HasPrefix(l, field+":") {
				kept = append(kept, l)
			}
		}
		_, problems := Validate("2026-10-06-bottlenecks.md", strings.Join(kept, "\n"))
		if len(problems) != 1 || problems[0].Field != field || problems[0].Message != "front matter needs "+field {
			t.Errorf("without %s: %+v", field, problems)
		}
	}
	for _, c := range []struct{ name, from, to, field, says string }{
		{"2026-10-06-bottlenecks.md", "updated: 2026-10-06T08:00:00Z", "updated: yesterday", "updated", `updated is "yesterday"`},
		{"2026-10-06-bottlenecks.md", "status: active", "status: final", "status", `status is "final", which is none of active, draft, deprecated`},
		{"2026-10-06-bottlenecks.md", "focus: bottlenecks", "focus: speed", "focus", `focus is "speed", which is none of bottlenecks, intent, risk, all`},
		{"2026-10-06-bottlenecks.md", "focus: bottlenecks", "focus: risk", "focus", "focus is risk but the file is named for bottlenecks"},
		{"2026-10-06-bottlenecks.md", "from: 2026-09-01", "from: September", "from", `from is "September", not a date`},
		{"2026-10-06-bottlenecks.md", "to: 2026-09-30", "to: 2026-08-31", "to", "the window ends on 2026-08-31, before it starts on 2026-09-01"},
		{"report.md", "", "", "", "a report is named <date>-<focus>.md"},
		{"2026-13-06-bottlenecks.md", "", "", "", "the file is named for 2026-13-06, which is not a date"},
		{"2026-10-06-speed.md", "", "", "", "the file is named for focus speed, which is none of"},
	} {
		content := good
		if c.from != "" {
			content = strings.Replace(good, c.from, c.to, 1)
		}
		_, problems := Validate(c.name, content)
		if len(problems) != 1 || problems[0].Field != c.field || !strings.Contains(problems[0].Message, c.says) {
			t.Errorf("%s with %q: %+v, want one problem on %q saying %q", c.name, c.to, problems, c.field, c.says)
		}
	}
	if _, problems := Validate("2026-10-06-risk.md", "# No front matter\n"); len(problems) != 1 || !strings.Contains(problems[0].Message, "no front matter") {
		t.Errorf("no front matter: %+v", problems)
	}
}

// S-0223: the index lists a report by a link to it.
func TestListsFindsALinkToTheReport(t *testing.T) {
	index := "| Report | Focus |\n|---|---|\n| [2026-10-06-risk.md](2026-10-06-risk.md) | risk |\n- [Intent](./2026-10-05-intent.md)\n"
	for name, want := range map[string]bool{"2026-10-06-risk.md": true, "2026-10-05-intent.md": true, "2026-10-06-all.md": false} {
		if got := Lists(index, name); got != want {
			t.Errorf("Lists(%q) = %v, want %v", name, got, want)
		}
	}
	if !Skipped("README.md") || Skipped("2026-10-06-risk.md") {
		t.Error("only the index is skipped")
	}
}

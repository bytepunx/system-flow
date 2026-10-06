// Package analysis knows the analyzer's reports (S-0223): one file per run
// under design/analysis, named for the day it ran and its focus, with front
// matter that says what it looked for and the window its metrics cover, and
// listed in the folder's README.md.
package analysis

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folder is the reports folder's name inside the design folder.
const Folder = "analysis"

// Index is the folder's index, which lists every report.
const Index = "README.md"

// Dir is the reports folder, relative to the project root, slash separated.
func Dir(m manifest.Manifest) string {
	return path.Join(m.Layout["design"], Folder)
}

// Focuses are what the analyzer can be asked to look for, harness.Focuses,
// which its prompt names; a report's focus is one of them or AllFocus.
var Focuses = harness.Focuses

// AllFocus is the focus of a report from a run asked for no focus.
const AllFocus = harness.AllFocus

// Statuses are what a report's status is, the documentation standard's:
// draft while the analyzer writes it, active once the run has ended, and
// deprecated when a later report replaces it.
var Statuses = []string{"active", "draft", "deprecated"}

// ReportPath is a report's file name in the reports folder, <date>-<focus>.md,
// where date is YYYY-MM-DD and focus is AllFocus for a run asked for none,
// focus empty. Join it to Dir for its path in the project.
func ReportPath(date, focus string) string {
	if focus == "" {
		focus = AllFocus
	}
	return date + "-" + focus + ".md"
}

// Front is a report's front matter: From and To are the window its metrics
// cover.
type Front struct {
	Title   string `yaml:"title"`
	Updated string `yaml:"updated"`
	Status  string `yaml:"status"`
	Focus   string `yaml:"focus"`
	From    string `yaml:"from"`
	To      string `yaml:"to"`
}

// Problem is one way a report falls short: Field names the front matter key
// it is about, or is empty for the file.
type Problem struct {
	Field   string
	Message string
}

var fileName = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-([a-z]+)\.md$`)

// Skipped reports whether a file in the reports folder is not a report: the
// folder's index.
func Skipped(name string) bool { return name == Index }

// Validate checks a report named name with content: that the file is named
// <date>-<focus>.md, and its front matter, each field there and well formed,
// the focus the one the file is named for, and the window in order.
func Validate(name, content string) (Front, []Problem) {
	var f Front
	var out []Problem
	focuses := append(slices.Clone(Focuses), AllFocus)
	m := fileName.FindStringSubmatch(name)
	switch {
	case m == nil:
		out = append(out, Problem{Message: "a report is named <date>-<focus>.md, the date as YYYY-MM-DD and the focus one of " + strings.Join(focuses, ", ")})
	case !isDate(m[1]):
		out = append(out, Problem{Message: fmt.Sprintf("the file is named for %s, which is not a date", m[1])})
		m = nil
	case !slices.Contains(focuses, m[2]):
		out = append(out, Problem{Message: fmt.Sprintf("the file is named for focus %s, which is none of %s", m[2], strings.Join(focuses, ", "))})
		m = nil
	}
	fm, _, err := workitem.SplitFrontMatter(content)
	if err != nil {
		return f, append(out, Problem{Message: "no front matter; a report needs title, updated, status, focus, from, and to"})
	}
	if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
		return f, append(out, Problem{Message: fmt.Sprintf("front matter does not parse: %v", err)})
	}
	for _, k := range []struct{ name, value string }{{"title", f.Title}, {"updated", f.Updated}, {"status", f.Status}, {"focus", f.Focus}, {"from", f.From}, {"to", f.To}} {
		if k.value == "" {
			out = append(out, Problem{Field: k.name, Message: "front matter needs " + k.name})
		}
	}
	if f.Updated != "" && !isDate(f.Updated) && !isTimestamp(f.Updated) {
		out = append(out, Problem{Field: "updated", Message: fmt.Sprintf("updated is %q, neither a date, YYYY-MM-DD, nor a timestamp, YYYY-MM-DDTHH:MM:SSZ", f.Updated)})
	}
	if f.Status != "" && !slices.Contains(Statuses, f.Status) {
		out = append(out, Problem{Field: "status", Message: fmt.Sprintf("status is %q, which is none of %s", f.Status, strings.Join(Statuses, ", "))})
	}
	switch {
	case f.Focus == "":
	case !slices.Contains(focuses, f.Focus):
		out = append(out, Problem{Field: "focus", Message: fmt.Sprintf("focus is %q, which is none of %s", f.Focus, strings.Join(focuses, ", "))})
	case m != nil && f.Focus != m[2]:
		out = append(out, Problem{Field: "focus", Message: fmt.Sprintf("focus is %s but the file is named for %s", f.Focus, m[2])})
	}
	for _, k := range []struct{ name, value string }{{"from", f.From}, {"to", f.To}} {
		if k.value != "" && !isDate(k.value) {
			out = append(out, Problem{Field: k.name, Message: fmt.Sprintf("%s is %q, not a date, YYYY-MM-DD", k.name, k.value)})
		}
	}
	if isDate(f.From) && isDate(f.To) && f.From > f.To {
		out = append(out, Problem{Field: "to", Message: fmt.Sprintf("the window ends on %s, before it starts on %s", f.To, f.From)})
	}
	return f, out
}

func isDate(s string) bool {
	_, err := time.Parse(time.DateOnly, s)
	return err == nil
}

func isTimestamp(s string) bool {
	_, err := time.Parse(workitem.TimeFormat, s)
	return err == nil
}

// Lists reports whether index, the folder's README.md, lists the report
// named name: a link to it.
func Lists(index, name string) bool {
	return strings.Contains(index, "]("+name+")") || strings.Contains(index, "](./"+name+")")
}

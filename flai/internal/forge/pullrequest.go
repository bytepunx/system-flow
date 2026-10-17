package forge

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// StoryPullRequest is what a story's pull request is built from: the story,
// its last flai verify result, and the project's pull request template.
type StoryPullRequest struct {
	// ID is the story's ID: "S-0393".
	ID string
	// Title is the story's title, without its ID.
	Title string
	// Criteria are the story's acceptance criteria, in the order it lists them.
	Criteria []Criterion
	// Verify is the story's last flai verify result; nil when it has none.
	Verify *VerifyResult
	// Template is the project's .github/pull_request_template.md; empty, or
	// only blank, when the project has none.
	Template string
}

// Criterion is one acceptance criterion of a story.
type Criterion struct {
	// Text is the criterion's words, after its box.
	Text string
	// Ticked is whether the story has it ticked.
	Ticked bool
}

// VerifyResult is a story's last flai verify result, as a pull request reports it.
type VerifyResult struct {
	// Passed is whether every step passed.
	Passed bool
	// StoppedAt names the step that failed; empty when every step passed.
	StoppedAt string
	// Commit is the story's worktree HEAD when it was verified.
	Commit string
	// RanAt is when the run began.
	RanAt time.Time
}

// templatePart is one place the pull request template marks for flai to
// fill, between <!-- flai:name --> and <!-- /flai:name -->.
type templatePart struct {
	name    string
	heading string
	// keep is whether what the template has between the marks stays as
	// the project wrote it rather than being replaced.
	keep bool
}

// templateParts are the four parts of a story's pull request, in the order
// they are appended when a template lacks their marks.
var templateParts = []templatePart{
	{name: "story", heading: "Story"},
	{name: "criteria", heading: "Acceptance criteria"},
	{name: "definition-of-done", heading: "Definition of done", keep: true},
	{name: "verify", heading: "Last flai verify"},
}

// defaultDefinitionOfDone is the definition of done a pull request has when
// the project's template does not mark its own.
const defaultDefinitionOfDone = "- [ ] Acceptance criteria checked, tasks done or cancelled\n" +
	"- [ ] Decisions captured (ADR, design/system, design/tech)\n" +
	"- [ ] docs/ updated for user-facing changes\n" +
	"- [ ] `flai check` passes\n" +
	"- [ ] Narrative has a closing log entry"

// PullRequestTitle is the title of a story's pull request: "[S-nnnn] <title>".
func (s StoryPullRequest) PullRequestTitle() string {
	return fmt.Sprintf("[%s] %s", s.ID, strings.TrimSpace(s.Title))
}

// PullRequestBody is the body of a story's pull request: the project's
// template with the places it marks filled, and each part it does not mark
// appended under a heading of its own; without a template, the four parts
// under headings of their own. It is refused when the template has a begin
// mark without its end mark, or an end mark without its begin mark.
func (s StoryPullRequest) PullRequestBody() (string, error) {
	body := s.Template
	if strings.TrimSpace(body) == "" {
		body = ""
	}
	for _, p := range templateParts {
		content := s.partContent(p)
		filled, found, err := fillPart(body, p, content)
		if err != nil {
			return "", fmt.Errorf("build the pull request body of %s from .github/pull_request_template.md: %w", s.ID, err)
		}
		if !found {
			filled = appendPart(filled, p, content)
		}
		body = filled
	}
	return strings.TrimLeft(body, "\n"), nil
}

// partContent is what flai writes between the marks of p; for the
// definition of done, which a template keeps as the project wrote it, the
// default a template without its marks gets.
func (s StoryPullRequest) partContent(p templatePart) string {
	switch p.name {
	case "story":
		return s.ID
	case "criteria":
		return criteriaText(s.Criteria)
	case "definition-of-done":
		return defaultDefinitionOfDone
	case "verify":
		return verifyText(s.ID, s.Verify)
	}
	return ""
}

// criteriaText is the criteria as checkboxes, ticked as the story has them.
func criteriaText(criteria []Criterion) string {
	if len(criteria) == 0 {
		return "The story lists no acceptance criteria."
	}
	lines := make([]string, 0, len(criteria))
	for _, c := range criteria {
		box := "[ ]"
		if c.Ticked {
			box = "[x]"
		}
		lines = append(lines, fmt.Sprintf("- %s %s", box, c.Text))
	}
	return strings.Join(lines, "\n")
}

// verifyText is the outcome of the story's last flai verify and the commit
// it ran at, or that it has none.
func verifyText(id string, v *VerifyResult) string {
	if v == nil {
		return fmt.Sprintf("No `flai verify` result is stored for %s; run `flai verify %s`.", id, id)
	}
	outcome := "Passed"
	if !v.Passed {
		outcome = fmt.Sprintf("Stopped at step `%s`", v.StoppedAt)
	}
	text := fmt.Sprintf("%s at commit %s", outcome, v.Commit)
	if !v.RanAt.IsZero() {
		text += fmt.Sprintf(", run %s", v.RanAt.UTC().Format(time.RFC3339))
	}
	return text + "."
}

// partMarks are the begin and end marks of p, with any spaces inside the
// comment.
func partMarks(p templatePart) (begin, end *regexp.Regexp) {
	name := regexp.QuoteMeta(p.name)
	return regexp.MustCompile(`<!--[ \t]*flai:` + name + `[ \t]*-->`),
		regexp.MustCompile(`<!--[ \t]*/flai:` + name + `[ \t]*-->`)
}

// fillPart writes content between each pair of p's marks in body, or, for a
// part the template keeps, leaves what is there; found says whether body
// has p's marks.
func fillPart(body string, p templatePart, content string) (string, bool, error) {
	beginMark, endMark := partMarks(p)
	var b strings.Builder
	found := false
	for rest := body; ; {
		begin := beginMark.FindStringIndex(rest)
		end := endMark.FindStringIndex(rest)
		if end != nil && (begin == nil || end[0] < begin[0]) {
			return "", false, fmt.Errorf("it has <!-- /flai:%s --> without <!-- flai:%s --> before it; add the begin mark before the part, or remove the end mark", p.name, p.name)
		}
		if begin == nil {
			b.WriteString(rest)
			return b.String(), found, nil
		}
		found = true
		after := rest[begin[1]:]
		end = endMark.FindStringIndex(after)
		if end == nil {
			return "", false, fmt.Errorf("it has <!-- flai:%s --> without <!-- /flai:%s --> after it; add the end mark after the part, or remove the begin mark", p.name, p.name)
		}
		b.WriteString(rest[:begin[1]])
		if p.keep {
			b.WriteString(after[:end[0]])
		} else {
			b.WriteString("\n" + content + "\n")
		}
		b.WriteString(after[end[0]:end[1]])
		rest = after[end[1]:]
	}
}

// appendPart adds p to the end of body under its heading, with its marks.
func appendPart(body string, p templatePart, content string) string {
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return body + fmt.Sprintf("\n## %s\n\n<!-- flai:%s -->\n%s\n<!-- /flai:%s -->\n", p.heading, p.name, content, p.name)
}

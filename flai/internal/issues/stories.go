package issues

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var (
	storyPattern = regexp.MustCompile(`^S-\d+$`)
	storyLineRe  = regexp.MustCompile(`(?m)^Story: (S-\d+)\.$`)
)

// storyID normalises the story an instance belongs to: empty for none, or a
// padded story ID.
func storyID(story string) (string, error) {
	story = strings.TrimSpace(story)
	if story == "" {
		return "", nil
	}
	id := workitem.CanonicalID(story)
	if !storyPattern.MatchString(id) {
		return "", fmt.Errorf("story %q is not a story ID like S-0001; give the story the issue was recorded for, or none", story)
	}
	return id, nil
}

// storyLine is the line naming an instance's story, with its newline, or ""
// for none.
func storyLine(story string) string {
	if story == "" {
		return ""
	}
	return "Story: " + story + ".\n"
}

// section finds the text under a "## " heading: from the line after it to the
// next "## " heading or the end of the body.
func section(body, heading string) (start, end int, ok bool) {
	i := strings.Index("\n"+body, "\n"+heading+"\n")
	if i < 0 {
		return 0, 0, false
	}
	start = i + len(heading) + 1
	end = len(body)
	if j := strings.Index(body[start-1:], "\n## "); j >= 0 {
		end = start + j
	}
	return start, end, true
}

// Stories returns the stories the issue's instances name, in order, once each.
func Stories(is *Issue) []string {
	start, end, ok := section(is.Body, "## Instances")
	if !ok {
		return nil
	}
	var out []string
	for _, m := range storyLineRe.FindAllStringSubmatch(is.Body[start:end], -1) {
		if id := workitem.CanonicalID(m[1]); !contains(out, id) {
			out = append(out, id)
		}
	}
	return out
}

// RecordedBy returns the issues with an instance that names the story.
func RecordedBy(list []*Issue, story string) []*Issue {
	story = workitem.CanonicalID(strings.TrimSpace(story))
	var out []*Issue
	for _, is := range list {
		if contains(Stories(is), story) {
			out = append(out, is)
		}
	}
	return out
}

// StoryDraft is a story made from an issue: what flai story new needs, and
// the body that goes below the story's "# S-nnnn Title" heading.
type StoryDraft struct {
	Title, Nature, Body string
}

// ForStory drafts the story that remediates the issue.
func ForStory(is *Issue) StoryDraft {
	nature := "improvement"
	if is.Class == "defect" || is.Class == "blocker" {
		nature = "remediation"
	}
	link := fmt.Sprintf("[%s](../../../design/issues/%s)", is.ID, filepath.Base(is.Path))
	var b strings.Builder
	fmt.Fprintf(&b, "## Goal\n\nThis story remediates %s, \"%s\".", link, is.Title)
	if solution := remediation(is); solution != "" {
		fmt.Fprintf(&b, " The issue recommends this solution:\n\n%s\n", solution)
	} else {
		b.WriteString(" The issue recommends no solution yet: propose one from its instances before building it.\n")
	}
	fmt.Fprintf(&b, "\n## Acceptance criteria\n- [ ] The cause %s describes no longer occurs, with a test that reproduces it where one fits\n", is.ID)
	fmt.Fprintf(&b, "- [ ] %s is closed with `flai issue close %s --reason` saying what fixed it\n", is.ID, is.ID)
	b.WriteString("\n## Tasks\n\n## Notes\n")
	return StoryDraft{Title: is.Title, Nature: nature, Body: b.String()}
}

// remediation is the text of the issue's Remediation section.
func remediation(is *Issue) string {
	start, end, ok := section(is.Body, "## Remediation")
	if !ok {
		return ""
	}
	return strings.TrimSpace(is.Body[start:end])
}

// Links reports whether a work item body names the issue by its ID, as a
// whole word, which a link to its document does too.
func Links(body string, is *Issue) bool {
	return namesWord(body, is.ID) || namesWord(body, workitem.CanonicalID(is.ID))
}

// namesWord reports whether word appears in s with no letter, digit, or
// underscore on either side.
func namesWord(s, word string) bool {
	for from := 0; ; {
		i := strings.Index(s[from:], word)
		if i < 0 {
			return false
		}
		i += from
		after := i + len(word)
		if (i == 0 || !wordByte(s[i-1])) && (after == len(s) || !wordByte(s[after])) {
			return true
		}
		from = i + 1
	}
}

func wordByte(c byte) bool {
	return c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// LinkedBy is the ID of the first open story, one not archived, done, or
// cancelled, whose body links the issue, or "" when none does.
func LinkedBy(is *Issue, items []*workitem.Item) string {
	for _, it := range items {
		if it.Type == workitem.Story && !it.Archived && !it.Closed() && Links(it.Body, is) {
			return it.ID
		}
	}
	return ""
}

// NoStory returns the open issues that no open story links.
func NoStory(list []*Issue, items []*workitem.Item) []*Issue {
	var out []*Issue
	for _, is := range list {
		if is.Status == "open" && LinkedBy(is, items) == "" {
			out = append(out, is)
		}
	}
	return out
}

// RepoFor is the project to read issues from for a story: a copy whose Root
// is the story's worktree when it has one, since the issues a story records
// are committed on its branch until it is accepted, else r. Work items are
// read from the main checkout either way.
func RepoFor(r *workitem.Repo, story string) *workitem.Repo {
	if strings.TrimSpace(story) == "" {
		return r
	}
	path := r.WorktreePath(workitem.CanonicalID(strings.TrimSpace(story)))
	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		return r
	}
	wt := *r
	wt.Root = path
	return &wt
}

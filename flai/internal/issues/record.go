package issues

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// FindOpenByTitle returns the open issue titled exactly title, the lowest
// numbered when several are, or nil when none is.
func FindOpenByTitle(list []*Issue, title string) *Issue {
	var found *Issue
	for _, is := range list {
		if is.Status == "open" && is.Title == title && (found == nil || num(is.ID) < num(found.ID)) {
			found = is
		}
	}
	return found
}

// HasInstance reports whether an instance of the issue names story and has
// note as one of its paragraphs, whatever its timestamp, so that the same
// occurrence is not recorded twice.
func HasInstance(is *Issue, story, note string) bool {
	story, err := storyID(story)
	if err != nil {
		return false
	}
	note = strings.TrimSpace(note)
	start, end, ok := section(is.Body, "## Instances")
	if !ok || note == "" {
		return false
	}
	line := strings.TrimSuffix(storyLine(story), "\n")
	for _, block := range strings.Split("\n"+is.Body[start:end], "\n### ") {
		_, text, _ := strings.Cut(block, "\n")
		named := line == ""
		var kept []string
		for _, l := range strings.Split(text, "\n") {
			if storyLineRe.MatchString(l) {
				if line != "" && workitem.CanonicalID(storyLineRe.FindStringSubmatch(l)[1]) == story {
					named = true
				}
				continue
			}
			kept = append(kept, l)
		}
		if !named {
			continue
		}
		for _, para := range strings.Split(strings.Join(kept, "\n"), "\n\n") {
			if strings.TrimSpace(para) == note {
				return true
			}
		}
	}
	return false
}

// Outcome says what RecordOnce did.
type Outcome string

// What RecordOnce did with an occurrence.
const (
	Opened  Outcome = "opened"
	Bumped  Outcome = "bumped"
	Already Outcome = "already recorded"
)

// RecordOnce records an occurrence in the open issue titled opt.Title: it
// opens the issue, with this occurrence as its first instance, when none is
// open; bumps it when no instance names opt.Story with opt.Note; and leaves
// it as it is when one does. It does not regenerate summary.md.
func RecordOnce(r *workitem.Repo, opt NewOptions) (*Issue, Outcome, error) {
	list, err := List(r)
	if err != nil {
		return nil, "", fmt.Errorf("read the issues in %s: %w", Dir(r), err)
	}
	is := FindOpenByTitle(list, opt.Title)
	if is == nil {
		is, err := New(r, opt)
		if err != nil {
			return nil, "", err
		}
		return is, Opened, nil
	}
	if HasInstance(is, opt.Story, opt.Note) {
		return is, Already, nil
	}
	if err := Bump(is, opt.Story, opt.Cost, opt.Note, opt.Now); err != nil {
		return nil, "", err
	}
	return is, Bumped, nil
}

// NewOrBump records a finding. Without opt.Report it is New, Opened. With
// it, the finding of an analysis report (S-0224), it bumps the open issue
// titled opt.Title, when one is, with the report, impact, story, cost, and
// note opt gives, rather than open a second, and leaves it as it is when an
// instance of it names that report already, so that filing a report's
// findings again does not count them twice. It does not regenerate
// summary.md. Anything it refuses leaves nothing written.
func NewOrBump(r *workitem.Repo, opt NewOptions) (*Issue, Outcome, error) {
	if strings.TrimSpace(opt.Report) == "" {
		is, err := New(r, opt)
		if err != nil {
			return nil, "", err
		}
		return is, Opened, nil
	}
	_, report, err := opt.check(r)
	if err != nil {
		return nil, "", err
	}
	list, err := List(r)
	if err != nil {
		return nil, "", fmt.Errorf("read the issues in %s: %w", Dir(r), err)
	}
	is := FindOpenByTitle(list, opt.Title)
	if is == nil {
		is, err := New(r, opt)
		if err != nil {
			return nil, "", err
		}
		return is, Opened, nil
	}
	if hasReport(is, report) {
		return is, Already, nil
	}
	if err := BumpWith(r, is, BumpOptions{Story: opt.Story, Cost: opt.Cost, Note: opt.Note, Impact: opt.Impact, Report: report, Now: opt.Now}); err != nil {
		return nil, "", err
	}
	return is, Bumped, nil
}

package digest

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// titleLimit is the most characters a candidate's title keeps before it is
// cut on a word and ends in an ellipsis.
const titleLimit = 200

// The lines that fence the threads mirror flai writes into a narrative.
const (
	mirrorStart = "<!-- threads:start -->"
	mirrorEnd   = "<!-- threads:end -->"
)

var (
	// narrativeFile is a story's narrative, named for the story.
	narrativeFile = regexp.MustCompile(`^S-(\d+)\.md$`)
	// bulletLine is a top-level list item: its marker and its text.
	bulletLine = regexp.MustCompile(`^(?:[-*+]|\d+[.)])\s+(.*)$`)
	// boldLead is a bullet opening with a bold label, "**Label.** text".
	boldLead = regexp.MustCompile(`^\*\*(.+?)\*\*[.:]?\s*(.*)$`)
	// closedLine is the line issues.Close appends: its time and its reason.
	closedLine = regexp.MustCompile(`(?m)^Closed (\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ):[ \t]*(.*)$`)
	// remediatesLine is the line issues.LinkStory writes, which is no fix.
	remediatesLine = regexp.MustCompile(`(?m)^Story S-\d+ remediates this issue, created from it at \S+\.$`)
	// storyMention is a story ID anywhere in a text.
	storyMention = regexp.MustCompile(`\bS-\d{3,}\b`)
)

// Lookup is the ItemLookup of a project, archived items included. A story's
// topics are those topics.OfStory gives it, with its epic and its open tasks,
// and its paths are its touches; a task has its story's; an epic is worked
// out as a story with no epic or tasks would be, from its own topics, tags,
// and touches. It reads each item once and remembers the answer.
func Lookup(repo *workitem.Repo) ItemLookup {
	type answer struct{ topics, paths []string }
	seen := map[string]answer{}
	var open []*workitem.Item
	loaded := false
	tasks := func() []*workitem.Item {
		if !loaded {
			loaded = true
			// A story's archived tasks are closed and claim nothing; when the
			// board cannot be read, the story's topics come without its tasks'.
			open, _ = repo.List(false)
		}
		return open
	}
	var find func(id string) answer
	find = func(id string) answer {
		id = workitem.CanonicalID(strings.TrimSpace(id))
		if a, ok := seen[id]; ok {
			return a
		}
		a := answer{topics: []string{}, paths: []string{}}
		if it, err := repo.Get(id); err == nil {
			switch it.Type {
			case workitem.Story:
				var epic *workitem.Item
				if it.Parent != "" {
					epic, _ = repo.Get(it.Parent)
				}
				a = answer{topics.Names(topics.OfStory(it, epic, tasks(), repo.Manifest.Projects)), clean(it.Touches, filepath.ToSlash)}
			case workitem.Task:
				if it.Parent != "" {
					a = find(it.Parent)
				}
			case workitem.Epic:
				a = answer{topics.Names(topics.OfStory(it, nil, nil, repo.Manifest.Projects)), clean(it.Touches, filepath.ToSlash)}
			}
		}
		seen[id] = a
		return a
	}
	return func(id string) ([]string, []string) {
		a := find(id)
		return slices.Clone(a.topics), slices.Clone(a.paths)
	}
}

// Decisions lists each top-level bullet under "## Decisions" in the
// narrative of every done story completed at or after since, a zero since
// being every one: the narratives in the archive, and those in the agents
// folder whose story is done and not yet archived. They come in story order,
// each narrative's bullets in the order written. A nil lookup is Lookup(repo).
func Decisions(repo *workitem.Repo, since time.Time, lookup ItemLookup) ([]Candidate, error) {
	if lookup == nil {
		lookup = Lookup(repo)
	}
	files := map[int]string{}
	for _, dir := range []string{repo.AgentsDir(), filepath.Join(repo.ArchiveDir(), "agents")} {
		entries, err := os.ReadDir(dir)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("listing the narratives in %s: %w", repoPath(repo, dir), err)
		}
		for _, e := range entries {
			if m := narrativeFile.FindStringSubmatch(e.Name()); m != nil && !e.IsDir() {
				n, _ := strconv.Atoi(m[1])
				files[n] = filepath.Join(dir, e.Name()) // the archive's copy wins
			}
		}
	}
	numbers := make([]int, 0, len(files))
	for n := range files {
		numbers = append(numbers, n)
	}
	slices.Sort(numbers)

	out := []Candidate{}
	for _, n := range numbers {
		file := files[n]
		id := fmt.Sprintf("S-%0*d", workitem.IDWidth, n)
		story, err := repo.Get(id)
		if err != nil || story.Status != workitem.Done {
			continue // a narrative whose story is open, cancelled, or gone records nothing learned
		}
		at := story.CompletedAt()
		if at.IsZero() || at.Before(since) {
			continue
		}
		narr, err := workitem.ReadNarrative(file)
		if err != nil {
			return nil, fmt.Errorf("reading the narrative %s for its decisions: %w; fix its front matter so it parses", repoPath(repo, file), err)
		}
		section, ok := workitem.NarrativeSection(narr.Body, "Decisions")
		if !ok {
			continue
		}
		topicNames, paths := lookup(story.ID)
		for _, text := range bullets(section.Text) {
			out = append(out, Candidate{
				Kind:    "decision",
				Source:  story.ID,
				Ref:     repoPath(repo, file),
				Title:   leadSentence(text),
				Topics:  slices.Clone(topicNames),
				Paths:   slices.Clone(paths),
				At:      at.UTC().Format(workitem.TimeFormat),
				Date:    at.UTC().Format(time.DateOnly),
				Text:    text,
				CitedBy: []string{},
			})
		}
	}
	return out, nil
}

// Causes lists each issue closed at or after since, a zero since being every
// one, that has a cause and a fix: the cause is its Description, the fix the
// reason its last "Closed <time>: <reason>" line gives, or its Remediation
// when that reason is empty. The close time is that line's, or the issue's
// updated time when it has none. Its topics and paths are those of the first
// story the fix names, else of the first story made from it. Issues come in
// ID order. A nil lookup is Lookup(repo).
func Causes(repo *workitem.Repo, since time.Time, lookup ItemLookup) ([]Candidate, error) {
	if lookup == nil {
		lookup = Lookup(repo)
	}
	list, err := issues.List(repo)
	if err != nil {
		return nil, fmt.Errorf("reading the issues for their causes: %w; fix the issue file it names so its front matter parses", err)
	}
	out := []Candidate{}
	for _, is := range list {
		if is.Status != "closed" {
			continue
		}
		stamp, reason := is.Updated, ""
		if all := closedLine.FindAllStringSubmatch(is.Body, -1); len(all) > 0 {
			last := all[len(all)-1]
			stamp, reason = last[1], strings.TrimSpace(last[2])
		}
		at, err := time.Parse(workitem.TimeFormat, stamp)
		if err != nil || at.Before(since) {
			continue
		}
		description, _ := workitem.NarrativeSection(is.Body, "Description")
		cause := strings.TrimSpace(description.Text)
		fix := reason
		if fix == "" {
			fix = remediationText(is.Body)
		}
		if cause == "" || fix == "" {
			continue
		}
		topicNames, paths := []string{}, []string{}
		if story := closingStory(is, fix); story != "" {
			topicNames, paths = lookup(story)
		}
		out = append(out, Candidate{
			Kind:    "cause",
			Source:  workitem.CanonicalID(is.ID),
			Title:   clipOnWord(strings.Join(strings.Fields(is.Title), " "), titleLimit),
			Topics:  topicNames,
			Paths:   paths,
			At:      at.UTC().Format(workitem.TimeFormat),
			Date:    at.UTC().Format(time.DateOnly),
			Text:    cause,
			Fix:     fix,
			CitedBy: []string{},
		})
	}
	return out, nil
}

// bullets are the top-level list items of a section's text, each with its
// marker dropped and its indented lines joined to it by single spaces. Text
// outside a list item and the threads mirror block are left out.
func bullets(text string) []string {
	var out, item []string
	open, mirror := false, false
	flush := func() {
		if open {
			out = append(out, strings.Join(strings.Fields(strings.Join(item, " ")), " "))
		}
		item, open = nil, false
	}
	for _, line := range strings.Split(text, "\n") {
		switch trimmed := strings.TrimSpace(line); {
		case mirror:
			mirror = trimmed != mirrorEnd
		case trimmed == mirrorStart:
			flush()
			mirror = true
		case trimmed == "":
			// A blank line ends nothing by itself: an indented line after it
			// still continues the item.
		case line[0] == ' ' || line[0] == '\t':
			if open {
				item = append(item, trimmed)
			}
		default:
			flush()
			if m := bulletLine.FindStringSubmatch(line); m != nil {
				item, open = []string{m[1]}, true
			}
		}
	}
	flush()
	return slices.DeleteFunc(out, func(s string) bool { return s == "" })
}

// leadSentence is a decision's title: its first sentence, a bold label it
// opens with joined to that sentence by a colon, cut to titleLimit.
func leadSentence(text string) string {
	if m := boldLead.FindStringSubmatch(text); m != nil {
		label := strings.TrimSpace(m[1])
		switch rest := firstSentence(m[2]); {
		case rest == "":
			text = strings.TrimRight(label, ".:")
		case strings.HasSuffix(label, "?") || strings.HasSuffix(label, "!"):
			text = label + " " + rest
		default:
			text = strings.TrimRight(label, ".:") + ": " + rest
		}
	} else {
		text = firstSentence(text)
	}
	return clipOnWord(strings.ReplaceAll(text, "**", ""), titleLimit)
}

// firstSentence is text up to and including the first full stop, question
// mark, or exclamation mark outside a code span that a space and a capital,
// a backtick, a parenthesis, or a quote follow, and that does not close a
// list number such as "1." or the abbreviations e.g. and i.e.
func firstSentence(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	runes := []rune(text)
	code := false
	for i, r := range runes {
		if r == '`' {
			code = !code
		}
		if code || !strings.ContainsRune(".!?", r) || i+2 >= len(runes) || runes[i+1] != ' ' {
			continue
		}
		next := runes[i+2]
		if !unicode.IsUpper(next) && !strings.ContainsRune("`(\"'", next) {
			continue
		}
		before := string(runes[:i])
		word := strings.ToLower(before[strings.LastIndex(before, " ")+1:])
		if word == "e.g" || word == "i.e" || strings.Trim(word, "0123456789") == "" {
			continue
		}
		return string(runes[:i+1])
	}
	return text
}

// clipOnWord cuts a text longer than limit characters at the last space
// before the limit and ends it with an ellipsis.
func clipOnWord(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	cut := string(runes[:limit])
	if i := strings.LastIndex(cut, " "); i > limit/2 {
		cut = cut[:i]
	}
	return strings.TrimRight(cut, " \n,;:-(") + "…"
}

// remediationText is an issue's Remediation section without its Closed
// lines and the lines naming the stories made from it.
func remediationText(body string) string {
	section, ok := workitem.NarrativeSection(body, "Remediation")
	if !ok {
		return ""
	}
	text := closedLine.ReplaceAllString(section.Text, "")
	text = remediatesLine.ReplaceAllString(text, "")
	return strings.TrimSpace(blankRuns.ReplaceAllString(text, "\n\n"))
}

// blankRuns is more than one blank line, which a removed line can leave.
var blankRuns = regexp.MustCompile(`\n{3,}`)

// closingStory is the story that closed an issue: the first the fix names,
// else the first story made from the issue, else none.
func closingStory(is *issues.Issue, fix string) string {
	if id := storyMention.FindString(fix); id != "" {
		return workitem.CanonicalID(id)
	}
	if made := issues.StoriesMade(is); len(made) > 0 {
		return made[0]
	}
	return ""
}

// repoPath is the path of a file in the work-in-process folder from the
// project root, slash separated, as the layout names it: the folder is the
// main checkout's even when repo is a worktree's.
func repoPath(repo *workitem.Repo, file string) string {
	rel, err := filepath.Rel(repo.WipDir(), file)
	if err != nil || !filepath.IsLocal(rel) {
		return filepath.ToSlash(file)
	}
	return path.Join(repo.Manifest.Layout["wip"], filepath.ToSlash(rel))
}

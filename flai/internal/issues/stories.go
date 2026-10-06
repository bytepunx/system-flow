package issues

import (
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
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

// StoryDraft is a story made from an issue: what flai story new needs, the
// body that goes below the story's "# S-nnnn Title" heading, and its planning
// data. Draft is always true: such a story is an agent's draft until the
// operator finalizes it (S-0203). CostOfDelay is nil when the issue gives no
// inputs. Strategic is a copy of what strategic agents spent on the issue,
// for CarryStrategic to charge to the story once it is made (S-0227); nil
// when the issue has no usage.
type StoryDraft struct {
	Title, Nature, Body string
	Draft               bool
	CostOfDelay         *workitem.CostOfDelay
	Strategic           []usage.Strategic
}

// ForStory drafts the story that remediates the issue, made at now in a
// project whose planning cycle, longer than zero, is cycle.
func ForStory(is *Issue, now time.Time, cycle time.Duration) StoryDraft {
	nature := "improvement"
	if is.Class == "defect" || is.Class == "blocker" {
		nature = "remediation"
	}
	link := fmt.Sprintf("[%s](%s%s/%s)", is.ID, storyToRoot, issuesFromRoot, filepath.Base(is.Path))
	reports := Reports(is)
	var b strings.Builder
	fmt.Fprintf(&b, "## Goal\n\nThis story remediates %s, \"%s\"", link, is.Title)
	if len(reports) > 0 {
		var links []string
		for _, rep := range reports {
			links = append(links, fmt.Sprintf("[%s](%s%s)", path.Base(rep), storyToRoot, rep))
		}
		analyses := "analysis"
		if len(reports) > 1 {
			analyses = "analyses"
		}
		fmt.Fprintf(&b, ", which the %s in %s found", analyses, joinAnd(links))
	}
	b.WriteString(".")
	if solution := remediation(is); solution != "" {
		fmt.Fprintf(&b, " The issue recommends this solution:\n\n%s\n", solution)
	} else {
		b.WriteString(" The issue recommends no solution yet: propose one from its instances before building it.\n")
	}
	fmt.Fprintf(&b, "\n## Acceptance criteria\n- [ ] The cause %s describes no longer occurs, with a test that reproduces it where one fits\n", is.ID)
	fmt.Fprintf(&b, "- [ ] %s is closed with `flai issue close %s --reason` saying what fixed it\n", is.ID, is.ID)
	b.WriteString("\n## Tasks\n\n## Notes\n")
	if len(reports) == 1 {
		fmt.Fprintf(&b, "\nThe analyzer found %s: the analysis the goal links gives its evidence.\n", is.ID)
	} else if len(reports) > 1 {
		fmt.Fprintf(&b, "\nThe analyzer found %s: the analyses the goal links give its evidence.\n", is.ID)
	}
	cod, notes := costOfDelay(is, now, cycle)
	if notes != "" {
		fmt.Fprintf(&b, "\n%s\n", notes)
	}
	var strategic []usage.Strategic
	if is.Usage != nil && len(is.Usage.Strategic) > 0 {
		strategic = is.Usage.Clone().Strategic
		var spent []string
		for _, s := range strategic {
			spent = append(spent, fmt.Sprintf("%s %s, %s USD", s.Kind, normalise((time.Duration(s.Seconds)*time.Second).String()), strconv.FormatFloat(s.Cost(), 'f', 4, 64)))
		}
		fmt.Fprintf(&b, "\nStrategic usage carried over from %s: %s.\n", is.ID, strings.Join(spent, "; "))
	}
	return StoryDraft{Title: is.Title, Nature: nature, Body: b.String(), Draft: true, CostOfDelay: cod, Strategic: strategic}
}

// CarryStrategic charges each strategic entry carried from an issue to the
// story with id, and so to its epic, as ADR-0083 charges a planner's
// (S-0227). It returns the IDs of the items it changed; on an error, it says
// which kinds were not carried.
func CarryStrategic(r *workitem.Repo, id string, entries []usage.Strategic) ([]string, error) {
	var changed []string
	for i, s := range entries {
		ids, err := r.ChargeStrategic(id, s.Kind, &usage.Usage{Seconds: s.Seconds, Models: s.Models})
		for _, c := range ids {
			if !contains(changed, c) {
				changed = append(changed, c)
			}
		}
		if err != nil {
			var left []string
			for _, l := range entries[i:] {
				left = append(left, l.Kind)
			}
			return changed, fmt.Errorf("charging the %s's usage to %s and each item above it: %w; what the %s spent is not carried in full", s.Kind, id, err, joinAnd(left))
		}
	}
	return changed, nil
}

// costOfDelay is the cost of delay inputs the issue gives, the inputs set by
// flai at now (ADR-0080), or nil when it gives none, and the sentences that say how each was set and
// which values were left out.
func costOfDelay(is *Issue, now time.Time, cycle time.Duration) (*workitem.CostOfDelay, string) {
	in, carried, skipped := impact(is)
	lost, how, derived := derive(is, now, cycle)
	fromImpact := in.TimeLostPerCycle != ""
	var set []string
	if derived && !fromImpact {
		in.TimeLostPerCycle = lost
		set = append(set, how)
	}
	if len(carried) > 0 {
		set = append(set, fmt.Sprintf("%s carried over from %s's Impact section.", joinAnd(carried), is.ID))
	}
	switch {
	case derived && fromImpact:
		set = append(set, fmt.Sprintf("Its Impact time_lost_per_cycle was taken rather than the %s derived from its cost and count.", lost))
	case !derived && how != "":
		skipped = append(skipped, how)
	}
	if in.IsZero() {
		return nil, strings.Join(skipped, " ")
	}
	set = append([]string{fmt.Sprintf("Cost of delay inputs set by flai from %s.", is.ID)}, set...)
	in.By, in.At = "flai", now.UTC().Format(workitem.TimeFormat)
	cod := &workitem.CostOfDelay{Inputs: &in}
	return cod, strings.Join(append(set, skipped...), " ")
}

var impactLineRe = regexp.MustCompile(`^(?:[-*]\s*)?(revenue_per_week|penalty_per_week|time_lost_per_cycle):(.*)$`)

// impact reads the inputs the issue's Impact section gives, one per line such
// as "- revenue_per_week: 1200"; other lines are evidence. The first value of
// a key that parses is used. It returns the inputs, each one carried over as
// "key value", and a sentence for each value left out.
func impact(is *Issue) (in workitem.CostInputs, carried, skipped []string) {
	start, end, ok := section(is.Body, "## Impact")
	if !ok {
		return in, nil, nil
	}
	for _, line := range strings.Split(is.Body[start:end], "\n") {
		m := impactLineRe.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		key, value := m[1], strings.TrimSpace(m[2])
		if key == "time_lost_per_cycle" {
			if in.TimeLostPerCycle != "" {
				continue
			}
			d, err := time.ParseDuration(value)
			if err != nil || d <= 0 {
				skipped = append(skipped, fmt.Sprintf("%s's Impact gives %s %q, which is not a duration longer than zero, so it was left out.", is.ID, key, value))
				continue
			}
			in.TimeLostPerCycle = normalise(d.String())
			carried = append(carried, key+" "+in.TimeLostPerCycle)
			continue
		}
		amount := &in.RevenuePerWeek
		if key == "penalty_per_week" {
			amount = &in.PenaltyPerWeek
		}
		if *amount != nil {
			continue
		}
		v, err := strconv.ParseFloat(value, 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
			skipped = append(skipped, fmt.Sprintf("%s's Impact gives %s %q, which is not an amount of zero or more, so it was left out.", is.ID, key, value))
			continue
		}
		*amount = &v
		carried = append(carried, key+" "+strconv.FormatFloat(v, 'f', -1, 64))
	}
	return in, carried, skipped
}

// derive works out the time lost per cycle from the issue's cost and count:
// cost × count ÷ the cycles since it was first reported, at least one. It
// returns the duration, the sentence saying how it was worked out, and true;
// or false with no duration, and a sentence when the cost is not a duration
// longer than zero.
func derive(is *Issue, now time.Time, cycle time.Duration) (lost, how string, ok bool) {
	if is.Cost == "" || is.Count < 1 {
		return "", "", false
	}
	cost, err := time.ParseDuration(is.Cost)
	if err != nil || cost <= 0 {
		return "", fmt.Sprintf("%s's cost %q is not a duration longer than zero, so no time_lost_per_cycle was derived from it.", is.ID, is.Cost), false
	}
	cycles, since := 1.0, ""
	if first, err := time.Parse(workitem.TimeFormat, is.FirstReported); err != nil {
		since = fmt.Sprintf("first reported %q is not a timestamp, so it counts as one cycle", is.FirstReported)
	} else {
		elapsed := now.Sub(first)
		days := tenths(elapsed.Hours() / 24)
		since = fmt.Sprintf("first reported %s, %s %s before this story", is.FirstReported, decimal(days), plural(days, "day"))
		if elapsed < cycle {
			since += "; under one cycle counts as one"
		} else {
			cycles = tenths(float64(elapsed) / float64(cycle))
		}
	}
	per := time.Duration(float64(cost) * float64(is.Count) / cycles)
	if per >= time.Minute {
		per = per.Round(time.Minute)
	} else {
		per = per.Round(time.Second)
	}
	per = max(per, time.Second)
	lost = normalise(per.String())
	how = fmt.Sprintf("time_lost_per_cycle %s: %s per occurrence × %d %s ÷ %s %s of %s (%s).",
		lost, normalise(cost.String()), is.Count, plural(float64(is.Count), "occurrence"), decimal(cycles), plural(cycles, "cycle"), normalise(cycle.String()), since)
	return lost, how, true
}

// tenths rounds x to one decimal place.
func tenths(x float64) float64 {
	return math.Round(x*10) / 10
}

// decimal writes x with as few digits as it needs: 3, 2.5.
func decimal(x float64) string {
	return strconv.FormatFloat(x, 'f', -1, 64)
}

// plural is word for one of it, else word with an s.
func plural(n float64, word string) string {
	if n == 1 {
		return word
	}
	return word + "s"
}

// joinAnd lists the parts as prose: "a", "a and b", "a, b and c".
func joinAnd(parts []string) string {
	if len(parts) < 2 {
		return strings.Join(parts, "")
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

// LinkStory names the story that remediates the issue as the last paragraph
// of its Remediation section, adding the section when it has none, and saves
// the issue. The story is named by ID only: its file moves on acceptance.
func LinkStory(is *Issue, story string, now time.Time) error {
	id, err := storyID(story)
	if err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("no story given to link from %s: give the ID of the story that remediates it, like S-0001", is.ID)
	}
	ts := now.UTC().Format(workitem.TimeFormat)
	line := fmt.Sprintf("Story %s remediates this issue, created from it at %s.\n", id, ts)
	body := strings.TrimRight(is.Body, "\n") + "\n"
	if start, end, ok := section(body, "## Remediation"); ok {
		text, rest := strings.TrimSpace(body[start:end]), body[end:]
		if text != "" {
			text += "\n\n"
		}
		if rest != "" {
			rest = "\n" + rest
		}
		body = body[:start] + "\n" + text + line + rest
	} else {
		body += "\n## Remediation\n\n" + line
	}
	is.Body, is.Updated = body, ts
	return is.Save()
}

// linkedLineRe is a line LinkStory writes.
var linkedLineRe = regexp.MustCompile(`(?m)^Story S-\d+ remediates this issue, created from it at \S+\.\n?`)

// reportLinkLineRe is a line linkReport writes, with its newline.
var reportLinkLineRe = regexp.MustCompile(reportLinkRe.String() + `\n?`)

// remediation is the text of the issue's Remediation section, without the
// lines naming stories made from it before, which are no solution, and
// without the lines linking the reports that found it, whose links are from
// the issue's folder: the story's goal links those reports instead (S-0224).
func remediation(is *Issue) string {
	start, end, ok := section(is.Body, "## Remediation")
	if !ok {
		return ""
	}
	text := linkedLineRe.ReplaceAllString(is.Body[start:end], "")
	if stripped := reportLinkLineRe.ReplaceAllString(text, ""); stripped != text {
		text = blankRunRe.ReplaceAllString(stripped, "\n\n")
	}
	return strings.TrimSpace(text)
}

// blankRunRe is more than one blank line, which a stripped line can leave.
var blankRunRe = regexp.MustCompile(`\n{3,}`)

// The story made from an issue links the issue and its reports from the
// story's file, in wip/kanban/stories, through the project root.
const (
	storyToRoot    = "../../../"
	issuesFromRoot = "design/" + Folder
)

// reportLineRe is an instance's line naming its report, capturing its path
// from the project root.
var reportLineRe = regexp.MustCompile(`(?m)^Report: (.+)\.$`)

// Reports are the analysis reports the issue names, as paths from the project
// root, once each in order of first appearance: those its instances' Report
// lines name, then those its Remediation section links, whose links are from
// the issue's folder. A path that leaves the project is not a report.
func Reports(is *Issue) []string {
	var out []string
	add := func(p string) {
		p = path.Clean(p)
		if p == "." || path.IsAbs(p) || p == ".." || strings.HasPrefix(p, "../") || contains(out, p) {
			return
		}
		out = append(out, p)
	}
	if start, end, ok := section(is.Body, "## Instances"); ok {
		for _, m := range reportLineRe.FindAllStringSubmatch(is.Body[start:end], -1) {
			add(strings.TrimSpace(m[1]))
		}
	}
	for _, target := range ReportLinks(is) {
		add(path.Join(issuesFromRoot, target))
	}
	return out
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

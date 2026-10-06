package preview

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/experiment"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// An acceptance by the orchestrator (ADR-0093, S-0221) is refused unless the
// verifier passed at the story branch's head, every acceptance criterion is
// ticked, every file the branch changes is under the story's touches, no
// thread on the story or its tasks is open, and the orchestrator's evidence,
// when given, names for each criterion a file the branch changes. Each
// failing condition is a blocker of the preview, beside those of any
// acceptance.

// Codes of the blockers an acceptance by the orchestrator adds.
const (
	BlockNotStory             = "not_story"             // the item is not a story
	BlockUnverified           = "unverified"            // the verified commit is missing or not the branch head
	BlockCriterionUnticked    = "criterion_unticked"    // an acceptance criterion is unticked
	BlockOutsideTouches       = "outside_touches"       // the branch changes a file under none of the touches
	BlockThreadOpen           = "thread_open"           // a thread on the story or a task is not resolved
	BlockCriterionUnevidenced = "criterion_unevidenced" // the evidence names no changed file for a criterion
)

// Blocker is one failing condition of an acceptance by the orchestrator: a
// code naming the condition and a message naming what fails.
type Blocker struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AcceptOptions are what an acceptance adds to the preview (ADR-0093).
type AcceptOptions struct {
	Orchestrator bool      // an acceptance by the orchestrator: adds ADR-0093's blockers
	Verified     string    // the commit the orchestrator's verifier passed, as given
	Evidence     *Evidence // the orchestrator's evidence, nil when not given
}

// Evidence is the markdown an orchestrator writes from its verifier's report:
// a Verdict line and, per acceptance criterion, the changed files that meet
// it (ADR-0093).
type Evidence struct {
	Verdict  string              `json:"verdict"`
	Criteria []CriterionEvidence `json:"criteria"` // in the order written
	Text     string              `json:"text"`     // the evidence as given
}

// CriterionEvidence is one list item of the evidence: the files that meet
// acceptance criterion N.
type CriterionEvidence struct {
	N     int      `json:"n"`
	Files []string `json:"files"`
}

var (
	verdictLine  = regexp.MustCompile(`^\s*Verdict:(.*)$`)
	criterionRow = regexp.MustCompile(`^\s*[-*]\s+(\d+):(.*)$`)
	quoted       = regexp.MustCompile("`([^`]+)`")
)

// ParseEvidence reads an orchestrator's evidence: one "Verdict: <verdict>"
// line, and one list item "- <n>: <files>" per criterion. Files are the
// backtick-quoted spans of the item when it has any, so a note may follow
// them; otherwise they are the item's text split on commas. It is an error
// when there is no Verdict line or it is empty, when there are two, when no
// criterion is listed, or when a criterion is listed twice.
func ParseEvidence(text string) (*Evidence, error) {
	ev := &Evidence{Text: text, Criteria: []CriterionEvidence{}}
	seenVerdict := false
	listed := map[int]bool{}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimRight(line, "\r")
		if m := verdictLine.FindStringSubmatch(line); m != nil {
			if seenVerdict {
				return nil, fmt.Errorf("the evidence has two Verdict lines: keep the one line \"Verdict: <verdict>\" from the verifier's report")
			}
			seenVerdict = true
			ev.Verdict = strings.TrimSpace(m[1])
			continue
		}
		m := criterionRow.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, fmt.Errorf("the evidence lists criterion %q, which is not a criterion's number: %w", m[1], err)
		}
		if listed[n] {
			return nil, fmt.Errorf("the evidence lists criterion %d twice: name all the files that meet it in one item \"- %d: <files>\"", n, n)
		}
		listed[n] = true
		ev.Criteria = append(ev.Criteria, CriterionEvidence{N: n, Files: evidenceFiles(m[2])})
	}
	switch {
	case !seenVerdict:
		return nil, fmt.Errorf("the evidence has no Verdict line: start it with \"Verdict: <verdict>\" from the verifier's report")
	case ev.Verdict == "":
		return nil, fmt.Errorf("the evidence's Verdict line is empty: give the verifier's verdict after \"Verdict:\"")
	case len(ev.Criteria) == 0:
		return nil, fmt.Errorf("the evidence lists no criterion: add one list item \"- <n>: <files>\" per acceptance criterion, naming the changed files that meet it")
	}
	return ev, nil
}

// evidenceFiles is the files a criterion's item names: its backtick-quoted
// spans when it has any, otherwise its text split on commas.
func evidenceFiles(s string) []string {
	var parts []string
	if spans := quoted.FindAllStringSubmatch(s, -1); spans != nil {
		for _, sp := range spans {
			parts = append(parts, sp[1])
		}
	} else {
		parts = strings.Split(s, ",")
	}
	files := []string{}
	for _, p := range parts {
		if p = cleanPath(p); p != "" {
			files = append(files, p)
		}
	}
	return files
}

// cleanPath is a path as the evidence and the diff are compared: trimmed of
// spaces, backticks, a leading ./ and a trailing /.
func cleanPath(p string) string {
	p = strings.Trim(strings.TrimSpace(p), "`")
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(p), "./"), "/")
}

// orchestratorBlockers is what stops the orchestrator accepting it under
// ADR-0093, beside what stops any acceptance; verified is the commit given,
// resolved to its full name, or "" when it names none.
func orchestratorBlockers(r execx.Runner, repo *workitem.Repo, it *workitem.Item, opts AcceptOptions) (blockers []Blocker, verified string) {
	add := func(code, format string, args ...any) {
		blockers = append(blockers, Blocker{Code: code, Message: fmt.Sprintf(format, args...)})
	}
	if it.Type != workitem.Story {
		add(BlockNotStory, "the orchestrator accepts only a story in review; %s is a %s, which the operator accepts", it.ID, it.Type)
		return blockers, ""
	}
	branch := storygit.Branch(it.ID)

	// 1. the verifier passed at the branch's head
	head, headErr := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", branch+"^{commit}")
	head = strings.TrimSpace(head)
	switch given := strings.TrimSpace(opts.Verified); {
	case headErr != nil || head == "":
		add(BlockUnverified, "%s has no branch %s to verify here, so no commit its verifier passed can be its head; the operator accepts it", it.ID, branch)
	case given == "":
		add(BlockUnverified, "the orchestrator's acceptance needs the commit its verifier passed: run the verifier on %s at its head %s and give --verified %s", branch, short(head), short(head))
	case strings.HasPrefix(given, "-"):
		add(BlockUnverified, "--verified %s is not a commit: give the commit the verifier passed, the head of %s, %s", given, branch, short(head))
	default:
		full, err := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", given+"^{commit}")
		full = strings.TrimSpace(full)
		switch {
		case err != nil || full == "":
			add(BlockUnverified, "--verified %s names no commit in this repository: give the commit the verifier passed, the head of %s, %s", given, branch, short(head))
		case full != head:
			verified = full
			add(BlockUnverified, "--verified %s is not the head of %s, which is %s: the branch changed after the verifier's run, so verify again at %s", short(full), branch, short(head), short(head))
		default:
			verified = full
		}
	}

	// 2. every criterion is ticked
	criteria := itemedit.Criteria(it.Body)
	for _, c := range criteria {
		if !c.Ticked {
			add(BlockCriterionUnticked, "acceptance criterion %d of %s is unticked: %s", c.N, it.ID, c.Text)
		}
	}

	// 3. every changed file is under the story's touches, read as a claim
	// reads them; the workflow's own records (the wip folder, the issues an
	// agent records, an experiment's results) are the story's whatever it
	// touches
	var changed []string
	diff, err := storygit.StoryDiff(r, repo, it.ID)
	if err != nil {
		add(BlockOutsideTouches, "cannot read what %s changes to check it against its touches: %v", branch, err)
	} else {
		for _, f := range diff.Files {
			changed = append(changed, f.Path)
			if f.OldPath != "" {
				changed = append(changed, f.OldPath)
			}
		}
		claim := workitem.NewHolds(nil, repo.Manifest.Projects).Claim(it)
		records := []string{repo.Manifest.Layout["wip"], path.Join(repo.Manifest.Layout["design"], issues.Folder)}
		if experiment.Needs(it) {
			records = append(records, experiment.Dir(repo.Manifest))
		}
		for _, p := range changed {
			if !underAny(p, claim) && !underAny(p, records) {
				add(BlockOutsideTouches, "%s changes %s, which is under none of the touches of %s", branch, p, it.ID)
			}
		}
	}

	// 4. no thread on the story or one of its tasks is open
	ids := []string{it.ID}
	if items, err := repo.List(false); err != nil {
		add(BlockThreadOpen, "cannot list the tasks of %s to check their threads: %v", it.ID, err)
	} else {
		for _, x := range items {
			if x.Type == workitem.Task && workitem.CanonicalID(x.Parent) == it.ID {
				ids = append(ids, x.ID)
			}
		}
	}
	seen := map[string]bool{}
	for _, id := range ids {
		ths, err := threads.For(repo, id)
		if err != nil {
			add(BlockThreadOpen, "cannot read the threads on %s: %v", id, err)
			continue
		}
		for _, th := range ths {
			if th.Open() && !seen[th.ID] {
				seen[th.ID] = true
				add(BlockThreadOpen, "thread %s on %s is %s, not resolved: %s", th.ID, id, th.Status, th.Title)
			}
		}
	}

	// 5. the evidence, when given, names a changed file for every criterion
	if opts.Evidence != nil && diff != nil {
		blockers = append(blockers, unevidenced(criteria, opts.Evidence, changed)...)
	}
	return blockers, verified
}

// unevidenced is a blocker for each criterion the evidence has no item for,
// or whose item names no file in changed.
func unevidenced(criteria []itemedit.Criterion, ev *Evidence, changed []string) []Blocker {
	isChanged := map[string]bool{}
	for _, p := range changed {
		isChanged[p] = true
	}
	byN := map[int][]string{}
	listed := map[int]bool{}
	for _, c := range ev.Criteria {
		byN[c.N], listed[c.N] = c.Files, true
	}
	var out []Blocker
	for _, c := range criteria {
		if !listed[c.N] {
			out = append(out, Blocker{Code: BlockCriterionUnevidenced, Message: fmt.Sprintf("the evidence has no item for acceptance criterion %d (%s): add \"- %d: <files>\" naming the changed files that meet it", c.N, c.Text, c.N)})
			continue
		}
		met := false
		for _, f := range byN[c.N] {
			if isChanged[f] {
				met = true
				break
			}
		}
		if !met {
			named := "no file"
			if len(byN[c.N]) > 0 {
				named = strings.Join(byN[c.N], ", ")
			}
			out = append(out, Blocker{Code: BlockCriterionUnevidenced, Message: fmt.Sprintf("the evidence for acceptance criterion %d (%s) names no file the branch changes (it names %s), so it cannot be checked against the diff", c.N, c.Text, named)})
		}
	}
	return out
}

// underAny says whether path p is one of dirs or under one of them.
func underAny(p string, dirs []string) bool {
	for _, d := range dirs {
		if d != "" && workitem.PathsOverlap(p, d) {
			return true
		}
	}
	return false
}

// short is a commit's name abbreviated to twelve characters.
func short(commit string) string {
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}

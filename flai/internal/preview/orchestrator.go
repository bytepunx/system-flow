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
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// An acceptance by the orchestrator (ADR-0093, S-0221) is refused unless the
// verifier passed at the story branch's head, or, as the merge queue's
// landing, the batch it lands from verified the story's rebased head
// (ADR-0176), every acceptance criterion is ticked, every file the branch
// changes is under the story's touches, no
// thread on the story or its tasks is open, and the orchestrator's evidence,
// when given, names for each criterion a file the branch changes, or, for a
// formula's instance, whose steps may change nothing, the done task of the
// criterion's step (S-0422), or the story's kluster record, when its
// namespace's deploy and checks passed at the verified head (S-0662); and the
// branch changes no path Claude Code protects, which the operator alone accepts (ADR-0106). Each failing
// condition is a blocker of the preview, beside those of any acceptance.

// Codes of the blockers an acceptance by the orchestrator adds.
const (
	BlockNotStory             = "not_story"             // the item is not a story
	BlockUnverified           = "unverified"            // the verified commit is missing or not the branch head
	BlockCriterionUnticked    = "criterion_unticked"    // an acceptance criterion is unticked
	BlockOutsideTouches       = "outside_touches"       // the branch changes a file under none of the touches
	BlockThreadOpen           = "thread_open"           // a thread on the story or a task is not resolved
	BlockOwnThreadOpen        = "own_thread_open"       // a thread the orchestrator opened and alone wrote on is not resolved (ADR-0189)
	BlockCriterionUnevidenced = "criterion_unevidenced" // the evidence names no changed file for a criterion, no kluster record that passed at the verified head, or, for a formula's instance, no done task of its step
	BlockProtectedPaths       = "protected_paths"       // the branch changes a path Claude Code protects, which the operator accepts (ADR-0106)
)

// Blocker is one failing condition of an acceptance by the orchestrator: a
// code naming the condition and a message naming what fails.
type Blocker struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// Name is the check's name or the reviewer's login, when the blocker is
	// one check's or one reviewer's.
	Name string `json:"name,omitempty"`
	// URL is the check's log, the review's page, or the pull request's page,
	// when the blocker has one.
	URL string `json:"url,omitempty"`
}

// AcceptOptions are what an acceptance adds to the preview (ADR-0093).
type AcceptOptions struct {
	Orchestrator bool      // an acceptance by the orchestrator: adds ADR-0093's blockers
	Verified     string    // the commit the orchestrator's verifier passed, as given
	Evidence     *Evidence // the orchestrator's evidence, nil when not given
	// Batch is the merge queue batch whose verification the acceptance
	// takes, as flai serve lands the story from it (ADR-0176); "" for any
	// other acceptance.
	Batch string
	// Keep are the paths whose removal of main's lines the operator keeps,
	// each with its reason (ADR-0155).
	Keep []KeptRemoval
	// QueueBatches says the merge queue can batch for the project, by
	// ADR-0174's three conditions, which the caller reads from the
	// configuration; WithoutQueueTiers is the operator's override of the
	// blocker it brings (ADR-0221).
	QueueBatches      bool
	WithoutQueueTiers bool
	// KlusterOn says the kluster host action is on for the project, and
	// ChecksOn the checks host action, which the caller reads from the
	// configuration: while KlusterOn and the manifest has a kluster: block,
	// every acceptance of a story waits for its namespace's deploy at its
	// branch's head, and with ChecksOn for the checks run against it there;
	// WithoutKluster is the operator's override of the blocker it brings,
	// which no acceptance by the orchestrator or the queue's landing takes
	// (ADR-0254 §5a).
	KlusterOn      bool
	ChecksOn       bool
	WithoutKluster bool
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
// acceptance criterion N, or, in place of files, the commit of the story's
// kluster record it names (S-0662).
type CriterionEvidence struct {
	N     int      `json:"n"`
	Files []string `json:"files"`
	// Kluster is the commit an item "- <n>: kluster deploy and checks passed
	// at <commit>" names; "" for an item naming files.
	Kluster string `json:"kluster,omitempty"`
}

var (
	verdictLine  = regexp.MustCompile(`^\s*Verdict:(.*)$`)
	criterionRow = regexp.MustCompile(`^\s*[-*]\s+(\d+):(.*)$`)
	quoted       = regexp.MustCompile("`([^`]+)`")
	// stepDone is a formula instance's evidence for a criterion: the task of
	// its step and its done transition, "T-nnnn done" (S-0422)
	stepDone = regexp.MustCompile("^`?([Tt]-\\d+)`?\\s+done\\b")
	// stepSkipped is a formula instance's evidence for a criterion whose step
	// a stop skipped: its cancelled task, "T-nnnn skipped" (ADR-0271 §4)
	stepSkipped = regexp.MustCompile("^`?([Tt]-\\d+)`?\\s+skipped\\b")
	// stopCancel is the line flai formula stop leaves in the Notes of each
	// step task it cancels (ADR-0271 §2), with the stopping step and the exit
	stopCancel = regexp.MustCompile(`moved to cancelled: step (\S+) stopped the run with exit (\d+)`)
	// klusterLine is evidence for a criterion that the story's namespace
	// shows it, in place of files: "kluster deploy and checks passed at
	// <commit>", or "kluster deploy passed at <commit>" while the checks host
	// action is off, a note may follow (ADR-0292 §6, S-0662)
	klusterLine = regexp.MustCompile("(?i)^kluster\\s+deploy(?:\\s+and\\s+checks)?\\s+passed\\s+at\\s+`?([0-9a-f]{7,40})\\b`?")
)

// ParseEvidence reads an orchestrator's evidence: one "Verdict: <verdict>"
// line, and one list item "- <n>: <files>" per criterion. Files are the
// backtick-quoted spans of the item when it has any, so a note may follow
// them; otherwise they are the item's text split on commas. An item "- <n>:
// kluster deploy and checks passed at <commit>" names no files but the
// commit of the story's kluster record (S-0662). It is an error
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
		if k := klusterLine.FindStringSubmatch(strings.TrimSpace(m[2])); k != nil {
			ev.Criteria = append(ev.Criteria, CriterionEvidence{N: n, Files: []string{}, Kluster: k[1]})
			continue
		}
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
// ADR-0093, beside what stops any acceptance, with prot the protected files
// its branch changes; verified is the commit given, resolved to its full
// name, or "" when it names none. As the merge queue's landing, opts.Batch,
// the batch's verification stands in for the verifier's (ADR-0176): rebased
// is the story's rebased head in the batch, which verified is, and landing
// every reason the queue does not let the story land, each a blocker.
func orchestratorBlockers(r execx.Runner, repo *workitem.Repo, it *workitem.Item, opts AcceptOptions, prot []string, rebased string, landing []string) (blockers []Blocker, verified string) {
	add := func(code, format string, args ...any) {
		blockers = append(blockers, Blocker{Code: code, Message: fmt.Sprintf(format, args...)})
	}
	if it.Type != workitem.Story {
		add(BlockNotStory, "the orchestrator accepts only a story in review; %s is a %s, which the operator accepts", it.ID, it.Type)
		return blockers, ""
	}
	branch := storygit.Branch(it.ID)

	// 1. the verifier passed at the branch's head, or the batch the story
	// lands from passed at its rebased head
	if opts.Batch != "" {
		for _, reason := range landing {
			add(BlockUnverified, "%s", reason)
		}
		verified = rebased
	} else {
		verified = verifiedHead(r, repo, it.ID, it.Formula != "", opts.Verified, add)
	}
	// an instance whose branch is gone changes nothing beyond main: its
	// close-out's commit is on main, which verifiedHead checked (ADR-0271 §4)
	branchless := it.Formula != "" && opts.Batch == "" && !storygit.BranchExists(r, repo.MainRoot, branch)

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
	var diff *storygit.Diff
	var err error
	if branchless {
		diff = &storygit.Diff{Story: it.ID, Branch: branch, Files: []storygit.DiffFile{}}
	} else {
		diff, err = storygit.StoryDiff(r, repo, it.ID)
	}
	if err != nil {
		add(BlockOutsideTouches, "cannot read what %s changes to check it against its touches: %v", branch, err)
	} else {
		changed = changedPaths(diff)
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

	// 4. no thread on the story or one of its tasks is open; one the
	// orchestrator opened and nobody else wrote on is its own blocker, the
	// orchestrator's to resolve (ADR-0189)
	ids := []string{it.ID}
	tasks := map[string]*workitem.Item{}
	items, listErr := repo.List(false)
	if listErr != nil {
		add(BlockThreadOpen, "cannot list the tasks of %s to check their threads: %v", it.ID, listErr)
	} else {
		for _, x := range items {
			if x.Type == workitem.Task && workitem.CanonicalID(x.Parent) == it.ID {
				ids = append(ids, x.ID)
				tasks[x.ID] = x
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
			if !th.Open() || seen[th.ID] {
				continue
			}
			seen[th.ID] = true
			if th.OrchestratorsOwn() {
				add(BlockOwnThreadOpen, "thread %s on %s is %s, not resolved: %s; the orchestrator opened it and nobody else wrote on it, so the orchestrator resolves it with a reason (thread_resolve, or flai thread resolve %s --reason) once its concern is gone, or leaves %s in review (ADR-0189)", th.ID, id, th.Status, th.Title, th.ID, it.ID)
				continue
			}
			add(BlockThreadOpen, "thread %s on %s is %s, not resolved: %s", th.ID, id, th.Status, th.Title)
		}
	}

	// 5. the evidence, when given, names a changed file for every criterion,
	// or the story's kluster record, judged at the verified head: the story
	// branch's, which the queue's landing rebased a copy of (S-0662); for a
	// formula's instance, whose steps may change no file, the task of the
	// criterion's step, done (S-0422)
	switch {
	case opts.Evidence == nil:
	case it.Formula != "" && listErr != nil:
		add(BlockCriterionUnevidenced, "cannot list the tasks of %s to check the evidence against its steps: %v", it.ID, listErr)
	case it.Formula != "":
		blockers = append(blockers, stepsUndone(it, criteria, opts.Evidence, tasks)...)
	case diff != nil:
		at := verified
		if opts.Batch != "" {
			at = ""
		}
		kl := klusterJudgeOf(r, repo.MainRoot, repo.Manifest, it.ID, at, opts.KlusterOn, opts.ChecksOn)
		blockers = append(blockers, unevidenced(criteria, opts.Evidence, changed, kl)...)
	}

	// 6. the branch changes no path Claude Code protects
	if len(prot) > 0 {
		add(BlockProtectedPaths, "%s changes paths Claude Code protects, %s, so only the operator (%s) accepts %s, whatever the orchestrator judges (ADR-0106)", branch, strings.Join(prot, ", "), operatorsOf(repo, it), it.ID)
	}
	return blockers, verified
}

// verifiedHead is the commit given, the one the orchestrator's verifier
// passed, resolved to its full name, or "" when it names none; add takes a
// blocker when it is not the head of story id's branch. For a formula's
// instance, instance, whose branch is gone, the commit its close-out passed
// at stands for the head when it is on the main branch (ADR-0271 §4).
func verifiedHead(r execx.Runner, repo *workitem.Repo, id string, instance bool, given string, add func(code, format string, args ...any)) (verified string) {
	branch := storygit.Branch(id)
	head, headErr := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", branch+"^{commit}")
	head = strings.TrimSpace(head)
	switch given = strings.TrimSpace(given); {
	case (headErr != nil || head == "") && instance:
		return verifiedOnMain(r, repo, id, given, add)
	case headErr != nil || head == "":
		add(BlockUnverified, "%s has no branch %s to verify here, so no commit its verifier passed can be its head; the operator accepts it", id, branch)
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
	return verified
}

// verifiedOnMain is the commit given resolved to its full name, for id, a
// formula's instance with no branch (ADR-0271 §4): add takes a blocker when
// it names no commit, or one not on the main branch, where the instance's
// close-out must have passed for a branch that carried nothing beyond it.
func verifiedOnMain(r execx.Runner, repo *workitem.Repo, id, given string, add func(code, format string, args ...any)) string {
	branch := storygit.Branch(id)
	if given == "" || strings.HasPrefix(given, "-") {
		add(BlockUnverified, "%s, a formula's instance, has no branch %s, so the orchestrator's acceptance needs the commit its close-out passed at, on the main branch: give it with --verified, from flai verify %s --last (ADR-0271)", id, branch, id)
		return ""
	}
	full, err := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", given+"^{commit}")
	if full = strings.TrimSpace(full); err != nil || full == "" {
		add(BlockUnverified, "--verified %s names no commit in this repository: give the commit %s's close-out passed at, from flai verify %s --last", given, id, id)
		return ""
	}
	main, err := storygit.MainBranch(r, repo.MainRoot)
	if err != nil {
		add(BlockUnverified, "cannot name the main branch to check that %s, %s's verified commit, is on it: %v", short(full), id, err)
		return full
	}
	if _, err := r.Run(repo.MainRoot, "git", "merge-base", "--is-ancestor", full, main); err != nil {
		add(BlockUnverified, "--verified %s is not on %s, and %s, a formula's instance, has no branch %s that could hold it: give the commit its close-out passed at, on %s, or the operator accepts it (ADR-0271)", short(full), main, id, branch, main)
	}
	return full
}

// unevidenced is a blocker for each criterion the evidence has no item for,
// whose item names no file in changed, or whose item names the story's
// kluster record at a commit kl says is no evidence, saying why.
func unevidenced(criteria []itemedit.Criterion, ev *Evidence, changed []string, kl klusterJudge) []Blocker {
	isChanged := map[string]bool{}
	for _, p := range changed {
		isChanged[p] = true
	}
	byN := map[int][]string{}
	klusterAt := map[int]string{}
	listed := map[int]bool{}
	for _, c := range ev.Criteria {
		byN[c.N], klusterAt[c.N], listed[c.N] = c.Files, c.Kluster, true
	}
	var out []Blocker
	for _, c := range criteria {
		if !listed[c.N] {
			out = append(out, Blocker{Code: BlockCriterionUnevidenced, Message: fmt.Sprintf("the evidence has no item for acceptance criterion %d (%s): add \"- %d: <files>\" naming the changed files that meet it", c.N, c.Text, c.N)})
			continue
		}
		if at := klusterAt[c.N]; at != "" {
			if why := kl(at); why != "" {
				out = append(out, Blocker{Code: BlockCriterionUnevidenced, Message: fmt.Sprintf("the evidence for acceptance criterion %d (%s) names the kluster record at %s, but %s", c.N, c.Text, at, why)})
			}
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

// klusterJudge is why the story's kluster record is no evidence for a
// criterion whose item names it at a commit, or "" when it is (S-0662).
type klusterJudge func(named string) string

// klusterJudgeOf judges the kluster record of story in the main checkout
// root, with m its manifest and git r, as kluster.ReadEvidence reads it
// with klusterOn and checksOn, the host actions, reading it once, on the
// first item that names it. verified is the commit the acceptance verified,
// "" to judge at the story branch's head alone.
func klusterJudgeOf(r execx.Runner, root string, m manifest.Manifest, story, verified string, klusterOn, checksOn bool) klusterJudge {
	var (
		ev   kluster.Evidence
		err  error
		read bool
	)
	return func(named string) string {
		if !read {
			ev, err = kluster.ReadEvidence(r, root, m, story, klusterOn, checksOn)
			read = true
		}
		if err != nil {
			return fmt.Sprintf("its record cannot be read: %v", err)
		}
		return klusterRecordWhy(r, root, m, ev, verified, named)
	}
}

// klusterRecordWhy is why ev, a story's kluster evidence, is no evidence for
// a criterion whose item names it at named, or "" when it is (ADR-0292 §6,
// S-0662): it is when ev applies, its deploy passed at the story branch's
// head, as kluster.AtHead widens it, and, while the checks host action is
// on, its checks run did too, the head is verified, the commit the
// acceptance verified when that is not "", and named is the deploy's or the
// checks run's commit, or a commit at the head the same way.
func klusterRecordWhy(r execx.Runner, root string, m manifest.Manifest, ev kluster.Evidence, verified, named string) string {
	switch {
	case !ev.Applies:
		return ev.Summary() + ", so name the changed files that meet it"
	case !ev.Passed:
		return "the record says " + ev.Summary()
	case verified != "" && verified != ev.Head:
		return fmt.Sprintf("the record is judged at the head of story/%s, %s, which is not the commit the acceptance verified, %s", ev.Story, short(ev.Head), short(verified))
	}
	full, err := r.Run(root, "git", "rev-parse", "--verify", "--quiet", named+"^{commit}")
	if full = strings.TrimSpace(full); err != nil || full == "" {
		return fmt.Sprintf("%s names no commit in this repository: name the commit of the record's deploy, %s", named, short(ev.Deploy.Commit))
	}
	if full == ev.Deploy.Commit || ev.Checks != nil && full == ev.Checks.Commit {
		return ""
	}
	j, err := kluster.AtHead(r, root, m, full, ev.Head)
	switch {
	case err != nil:
		return fmt.Sprintf("cannot tell whether %s stands at the head %s: %v", short(full), short(ev.Head), err)
	case !j.AtHead:
		return fmt.Sprintf("%s is not the commit of the record's deploy, %s, nor at the head %s: %s", short(full), short(ev.Deploy.Commit), short(ev.Head), j.Why)
	}
	return ""
}

// stepsUndone is a blocker for each criterion of it, a formula's instance,
// that the evidence has no item for, whose item names no task's done
// transition, "- <n>: T-nnnn done", or names one that is not a task of it in
// tasks or is not done: a step's evidence is its task's done transition, not
// a changed file, since a step may change nothing (S-0422). A step a stop
// skipped is evidenced by its task cancelled by that stop, "- <n>: T-nnnn
// skipped", as stopSkipped checks it (ADR-0271 §4).
func stepsUndone(it *workitem.Item, criteria []itemedit.Criterion, ev *Evidence, tasks map[string]*workitem.Item) []Blocker {
	said := evidenceItems(ev.Text)
	var out []Blocker
	block := func(format string, args ...any) {
		out = append(out, Blocker{Code: BlockCriterionUnevidenced, Message: fmt.Sprintf(format, args...)})
	}
	for _, c := range criteria {
		text, listed := said[c.N]
		if !listed {
			block("the evidence has no item for acceptance criterion %d (%s): %s is an instance of formula %s, so add \"- %d: T-nnnn done\" naming the task of its step", c.N, c.Text, it.ID, it.Formula, c.N)
			continue
		}
		if m := stepSkipped.FindStringSubmatch(text); m != nil {
			if why := stopSkipped(it, workitem.CanonicalID(m[1]), tasks); why != "" {
				block("the evidence for acceptance criterion %d (%s) names %s skipped, but %s", c.N, c.Text, workitem.CanonicalID(m[1]), why)
			}
			continue
		}
		m := stepDone.FindStringSubmatch(text)
		if m == nil {
			block("the evidence for acceptance criterion %d (%s) names no task's done transition (it says %q): %s is an instance of formula %s, so write \"- %d: T-nnnn done\" naming the task of its step", c.N, c.Text, text, it.ID, it.Formula, c.N)
			continue
		}
		id := workitem.CanonicalID(m[1])
		switch task := tasks[id]; {
		case task == nil:
			block("the evidence for acceptance criterion %d (%s) names %s, which is not a task of %s", c.N, c.Text, id, it.ID)
		case task.Status != workitem.Done:
			block("the evidence for acceptance criterion %d (%s) names %s, which is %s, not done: its step has not passed", c.N, c.Text, id, task.Status)
		}
	}
	return out
}

// stopSkipped is why task id is no evidence that a stop skipped its step, or
// "" when it is (ADR-0271 §4): it is a task of the instance it in tasks, it
// is cancelled, its Notes hold the cancel flai formula stop writes, "moved to
// cancelled: step <id> stopped the run with exit <n>", and a task of it is
// done with that stop's record, "Stopped the run with exit <n>:", in its
// Notes: the stopping step's task, which flai formula stop moved to done.
// The record names the stopping step's task, so the formula's file, which
// may have changed since the run, is not read.
func stopSkipped(it *workitem.Item, id string, tasks map[string]*workitem.Item) string {
	task := tasks[id]
	switch {
	case task == nil:
		return fmt.Sprintf("%s is not a task of %s", id, it.ID)
	case task.Status != workitem.Cancelled:
		return fmt.Sprintf("%s is %s, not cancelled: a skipped step's task is cancelled by flai formula stop", id, task.Status)
	}
	m := stopCancel.FindStringSubmatch(task.Body)
	if m == nil {
		return fmt.Sprintf("%s was not cancelled by a step's stop: its Notes hold no \"moved to cancelled: step <id> stopped the run with exit <n>\", which flai formula stop writes", id)
	}
	record := "Stopped the run with exit " + m[2] + ":"
	for _, t := range tasks {
		if t.Status == workitem.Done && strings.Contains(t.Body, record) {
			return ""
		}
	}
	return fmt.Sprintf("no task of %s is done with the record of step %s's stop, %q, in its Notes: the step that stopped the run has not passed", it.ID, m[1], record)
}

// evidenceItems is the text of each criterion's list item in an evidence,
// "- <n>: <text>", trimmed, by criterion number.
func evidenceItems(text string) map[int]string {
	out := map[int]string{}
	for _, line := range strings.Split(text, "\n") {
		m := criterionRow.FindStringSubmatch(strings.TrimRight(line, "\r"))
		if m == nil {
			continue
		}
		if n, err := strconv.Atoi(m[1]); err == nil {
			out[n] = strings.TrimSpace(m[2])
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

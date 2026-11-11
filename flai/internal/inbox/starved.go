package inbox

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A starved board (ADR-0285, I-0288): no ready story can start, the
// in-progress limit has room, and no backlog story is a promotion candidate.
// flai serve tells the operator, on one thread, each backlog story that cannot
// be promoted and the decision it waits on. This file builds that list; when
// the board is starved, and the thread, are flai serve's.

// What a starved story waits on, the first that applies (ADR-0285 §3).
const (
	WaitsDraft   = "draft"   // it is a draft: finalize it
	WaitsAfter   = "after"   // held by after: the stories at the root of its chain
	WaitsBlocked = "blocked" // blocked: its block reason
	WaitsOther   = "other"   // anything else: its reasons as the candidates give them
)

// Ref names a story or a thread by its ID and title, and a story's state.
type Ref struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status,omitempty"`
}

// StarvedStory is one backlog story that cannot be promoted, with the
// decision it waits on (ADR-0285 §3).
type StarvedStory struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Reasons are why it is not a candidate, as PromotionCandidates gives them.
	Reasons []string `json:"reasons"`
	// Waits is what it waits on: WaitsDraft, WaitsAfter, WaitsBlocked, or
	// WaitsOther.
	Waits string `json:"waits"`
	// Threads are, for a draft, the open threads on it or its tasks that
	// await the operator, in the order of the summaries given, and for a
	// blocked story, the threads its block reason names, in its order.
	Threads []Ref `json:"threads,omitempty"`
	// Roots are, for a story held by after, the stories at the root of its
	// chain: each named in after and not done, followed until one names
	// none that is not done.
	Roots []Ref `json:"roots,omitempty"`
	// Loops says that its after chain comes back to a story already on it.
	Loops bool `json:"loops,omitempty"`
	// Held is the hold's reason, for a story held by after.
	Held string `json:"held,omitempty"`
	// Block is the reason of its open block, for a blocked story.
	Block string `json:"block,omitempty"`
}

// Starved is the backlog as a starved board shows it to the operator: the
// promotion candidates, and every other backlog story with what it waits on,
// in the pull order.
type Starved struct {
	Candidates []workitem.Ranked `json:"candidates"`
	Stories    []StarvedStory    `json:"stories"`
}

// ReadStarved judges the backlog among items, the active items as
// repo.List(false) lists them, in the board's pull order by the project's
// policy, and says what each story that cannot be promoted waits on. Whom an
// open thread awaits is as r reads it (ADR-0181).
func ReadStarved(repo *workitem.Repo, items []*workitem.Item, r Reader) (Starved, error) {
	board, err := repo.LoadBoard()
	if err != nil {
		return Starved{}, fmt.Errorf("read the board's pull order to list the backlog's decisions: %w", err)
	}
	c, err := workitem.PromotionCandidates(items, board.Order, repo.Holds(items), repo.Manifest.Orchestration.PolicyOrDefault(), 0)
	if err != nil {
		return Starved{}, fmt.Errorf("judge the backlog's promotion candidates: %w; check orchestration.policy in system-flow.yaml", err)
	}
	all, err := threads.List(repo)
	if err != nil {
		return Starved{}, fmt.Errorf("read the threads to list the backlog's decisions: %w", err)
	}
	sums := make([]ThreadSummary, 0, len(all))
	for _, th := range all {
		if th.Open() {
			sums = append(sums, Summarize(repo, r, th))
			continue
		}
		// a resolved thread awaits no one; a block reason may still name it
		sums = append(sums, ThreadSummary{ID: th.ID, Title: th.Title, Status: th.Status, Anchor: th.Anchor})
	}
	return StarvedFrom(items, c, sums), nil
}

// StarvedFrom says what each story c refuses waits on, in c's order, given
// the items c judged and the summaries of the threads, each Story and
// Awaiting set as Summarize sets them. It reads nothing.
func StarvedFrom(items []*workitem.Item, c workitem.Candidates, sums []ThreadSummary) Starved {
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	out := Starved{Candidates: c.Candidates, Stories: []StarvedStory{}}
	if out.Candidates == nil {
		out.Candidates = []workitem.Ranked{}
	}
	for _, o := range c.Others {
		out.Stories = append(out.Stories, starvedStory(o, byID[o.ID], byID, sums))
	}
	return out
}

// starvedStory is what refused, the item story, waits on: a draft, held by
// after, blocked, or anything else, the first that applies.
func starvedStory(refused workitem.Refused, story *workitem.Item, byID map[string]*workitem.Item, sums []ThreadSummary) StarvedStory {
	out := StarvedStory{ID: refused.ID, Title: refused.Title, Reasons: refused.Reasons, Waits: WaitsOther}
	switch {
	case story == nil:
		// not among the items: its reasons are all there is to say
	case story.Draft:
		out.Waits, out.Threads = WaitsDraft, awaitingOperator(story.ID, sums)
	case refused.Held != nil && refused.Held.Code == workitem.HoldAfter:
		out.Waits, out.Held = WaitsAfter, refused.Held.Reason
		roots, loops := afterRoots(story, byID)
		out.Loops = loops
		for _, it := range roots {
			out.Roots = append(out.Roots, Ref{ID: it.ID, Title: it.Title, Status: it.Status})
		}
	case story.IsBlocked():
		out.Waits, out.Block = WaitsBlocked, openBlock(story)
		out.Threads = namedThreads(out.Block, sums)
	}
	return out
}

// awaitingOperator are the open threads on story or its tasks that await the
// operator; a resolved one is left out.
func awaitingOperator(story string, sums []ThreadSummary) []Ref {
	var out []Ref
	for _, s := range sums {
		if s.Status != "resolved" && s.Awaiting == "operator" && s.Story != "" && workitem.CanonicalID(s.Story) == story {
			out = append(out, Ref{ID: s.ID, Title: s.Title})
		}
	}
	return out
}

// afterRoots are the stories at the root of story's after chain, in the order
// the chain reaches them, and whether the chain loops. A story named in after
// that is done, cancelled, or absent, as an archived one is, counts as done.
func afterRoots(story *workitem.Item, byID map[string]*workitem.Item) ([]*workitem.Item, bool) {
	const (
		onChain = 1
		walked  = 2
	)
	state := map[string]int{}
	var roots []*workitem.Item
	loops := false
	var walk func(it *workitem.Item)
	walk = func(it *workitem.Item) {
		state[it.ID] = onChain
		waits := pendingAfter(it, byID)
		if len(waits) == 0 && it != story {
			roots = append(roots, it)
		}
		for _, w := range waits {
			switch state[w.ID] {
			case onChain:
				loops = true
			case 0:
				walk(w)
			}
		}
		state[it.ID] = walked
	}
	walk(story)
	return roots, loops
}

// pendingAfter are the stories it names in after that are among byID and
// neither done nor cancelled, each once, itself left out.
func pendingAfter(it *workitem.Item, byID map[string]*workitem.Item) []*workitem.Item {
	var out []*workitem.Item
	seen := map[string]bool{}
	for _, e := range it.After {
		id := workitem.CanonicalID(e)
		dep, ok := byID[id]
		if !ok || id == it.ID || seen[id] || dep.Closed() {
			continue
		}
		seen[id] = true
		out = append(out, dep)
	}
	return out
}

// openBlock is the reason of its newest open block.
func openBlock(it *workitem.Item) string {
	for i := len(it.Blocked) - 1; i >= 0; i-- {
		if it.Blocked[i].Until == "" {
			return it.Blocked[i].Reason
		}
	}
	return ""
}

// threadRef finds a thread's ID in a block reason, in any case and padding.
var threadRef = regexp.MustCompile(`(?i)\bTH-\d+\b`)

// namedThreads are the threads reason names, each once, with the title sums
// give it, or none when sums do not hold it.
func namedThreads(reason string, sums []ThreadSummary) []Ref {
	titles := map[string]string{}
	for _, s := range sums {
		titles[threads.CanonicalID(s.ID)] = s.Title
	}
	var out []Ref
	seen := map[string]bool{}
	for _, m := range threadRef.FindAllString(reason, -1) {
		id := threads.CanonicalID(m)
		if !seen[id] {
			seen[id] = true
			out = append(out, Ref{ID: id, Title: titles[id]})
		}
	}
	return out
}

// Text is the list as markdown, one bullet line per story in the pull order,
// each ended by a newline, and nothing when there are none: the same list
// gives the same text, so a caller compares it to tell whether the list
// changed.
func (s Starved) Text() string {
	var b strings.Builder
	for _, st := range s.Stories {
		b.WriteString("- " + oneLine(st.ID+" "+st.Title) + ": " + st.decision() + "\n")
	}
	return b.String()
}

// decision says in words what st waits on.
func (st StarvedStory) decision() string {
	switch st.Waits {
	case WaitsDraft:
		text := "a draft; finalize it."
		if len(st.Threads) > 0 {
			text += " " + sentence("It waits on "+listRefs(st.Threads))
		}
		return text
	case WaitsAfter:
		var text string
		switch len(st.Roots) {
		case 0:
			if st.Loops {
				return "held by `after` in a chain that loops back on itself; fix `after`."
			}
			return sentence(oneLine(st.Held))
		case 1:
			text = sentence("held by `after`; the root of its chain is " + listRefs(st.Roots))
		default:
			text = sentence("held by `after`; the roots of its chain are " + listRefs(st.Roots))
		}
		if st.Loops {
			text += " Its `after` chain loops back on itself; fix `after`."
		}
		return text
	case WaitsBlocked:
		text := sentence("blocked: " + oneLine(st.Block))
		if len(st.Threads) > 0 {
			text += " " + sentence("It waits on "+listRefs(st.Threads))
		}
		return text
	}
	return sentence(oneLine(strings.Join(st.Reasons, "; ")))
}

// listRefs names refs, a story's with its state in brackets, as "a", "a and
// b", or "a, b and c".
func listRefs(refs []Ref) string {
	names := make([]string, len(refs))
	for i, r := range refs {
		names[i] = oneLine(r.ID + " " + r.Title)
		if r.Status != "" {
			names[i] += " (" + r.Status + ")"
		}
	}
	if n := len(names); n > 1 {
		return strings.Join(names[:n-1], ", ") + " and " + names[n-1]
	}
	return strings.Join(names, "")
}

// oneLine is s on one line, each run of white space, tabs and newlines among
// it, one space.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// sentence is s ended by a full stop, unless it ends in one, a question
// mark, or an exclamation mark already.
func sentence(s string) string {
	if s == "" || strings.HasSuffix(s, ".") || strings.HasSuffix(s, "?") || strings.HasSuffix(s, "!") {
		return s
	}
	return s + "."
}

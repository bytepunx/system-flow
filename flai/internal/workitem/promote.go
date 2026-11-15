package workitem

import (
	"regexp"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// Promotion candidates (S-0217). A backlog story is a candidate to go to
// ready when it is not a draft, meets the definition of ready, would not be
// held if it were ready, and carries a forecast duration and a cost of delay
// value. The candidates are ordered by the project's policy; every other
// backlog story is listed with each reason it is not one. The backlog and
// ready are roles, not names: a promotion is the story transition with the
// commit effect (ADR-0301 §3), and the backlog is the state it leaves.

// promoteSource is the state a promotion takes a story from: the source of
// the story transition with the commit effect (ADR-0301 §3), the first it
// names when it names several; empty when the workflow has no commit.
func promoteSource(w workflow.Workflow) string {
	t, ok := w.WithEffect(workflow.TypeStory, workflow.EffectCommit)
	if !ok || len(t.From) == 0 {
		return ""
	}
	return t.From[0]
}

// promoteTarget is the state a promotion puts a story in, the state that
// commits it to be worked (ADR-0301 §3).
func promoteTarget(w workflow.Workflow) string {
	t, _ := w.WithEffect(workflow.TypeStory, workflow.EffectCommit)
	return t.To
}

// Candidates are the backlog stories that could go to ready, in the order a
// policy gives them, and the others with why not.
type Candidates struct {
	Policy     string    `json:"policy"`
	Candidates []Ranked  `json:"candidates"`
	Others     []Refused `json:"others"`
}

// Refused is a backlog story that is not a candidate, with every reason.
type Refused struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Reasons []string `json:"reasons"`
	// Held is the hold it would have in ready, when it would have one.
	Held *Hold `json:"held,omitempty"`
}

// PromotionCandidates judges the backlog stories among items, taken in their
// pull order by order, and orders the candidates by policy. A limit above
// zero caps the candidates listed; those beyond it are not listed at all.
// It runs the baseline workflow, as PromotionCandidatesIn does with it.
func PromotionCandidates(items []*Item, order []string, holds *Holds, policy string, limit int) (Candidates, error) {
	return PromotionCandidatesIn(baselineFlow, items, order, holds, policy, limit)
}

// PromotionCandidatesIn is PromotionCandidates under workflow w: the stories
// judged are those in the state its promotion leaves, and an epic is
// cancelled by its state's kind (ADR-0301 §2, §3).
func PromotionCandidatesIn(w workflow.Workflow, items []*Item, order []string, holds *Holds, policy string, limit int) (Candidates, error) {
	out := Candidates{Policy: policy, Candidates: []Ranked{}, Others: []Refused{}}
	byID := map[string]*Item{}
	for _, it := range items {
		if !it.Archived {
			byID[it.ID] = it
		}
	}
	var ok []*Item
	for _, id := range PullSequence(order, items, promoteSource(w)) {
		it := byID[id]
		reasons, held := promotionRefusals(w, it, byID, holds)
		if len(reasons) == 0 {
			ok = append(ok, it)
			continue
		}
		out.Others = append(out.Others, Refused{ID: it.ID, Title: it.Title, Reasons: reasons, Held: held})
	}
	ranked, err := OrderByPolicy(ok, policy)
	if err != nil {
		return Candidates{}, err
	}
	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}
	out.Candidates = append(out.Candidates, ranked...)
	return out, nil
}

// PromotionCandidates reads the active items, the board's pull order, and
// the holds, and judges the backlog by orchestration.policy, under the
// workflow the repository runs.
func (r *Repo) PromotionCandidates(limit int) (Candidates, error) {
	items, err := r.List(false)
	if err != nil {
		return Candidates{}, err
	}
	board, err := r.LoadBoard()
	if err != nil {
		return Candidates{}, err
	}
	return PromotionCandidatesIn(r.Flow(), items, board.Order, r.Holds(items), r.Manifest.Orchestration.PolicyOrDefault(), limit)
}

// promotionRefusals is every reason story is not a candidate, and the hold
// it would have in ready, under workflow w: an epic in a state of the
// cancelled kind refuses it, named by its own status (ADR-0301 §2).
func promotionRefusals(w workflow.Workflow, story *Item, byID map[string]*Item, holds *Holds) ([]string, *Hold) {
	var reasons []string
	if story.Draft {
		reasons = append(reasons, "draft: finalize it first")
	}
	if !hasGoal(story.Body) {
		reasons = append(reasons, "no goal")
	}
	if !hasCriteria(story.Body) {
		reasons = append(reasons, "no acceptance criteria with a checkbox")
	}
	if p, ok := byID[story.Parent]; ok && w.Kind(p.Status) == workflow.KindCancelled {
		reasons = append(reasons, "its epic "+p.ID+" is "+p.Status)
	}
	held := holds.Of(story)
	if held != nil {
		reasons = append(reasons, held.Reason)
	}
	if _, ok := forecastHours(story); !ok {
		reasons = append(reasons, "no forecast duration")
	}
	if codValue(story) == nil {
		reasons = append(reasons, "no cost of delay value")
	}
	return reasons, held
}

var goalHeading = regexp.MustCompile(`(?m)^## Goal\s*$`)

// hasGoal says whether a body's "## Goal" section has any words.
func hasGoal(body string) bool {
	loc := goalHeading.FindStringIndex(body)
	if loc == nil {
		return false
	}
	rest := body[loc[1]:]
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	return strings.TrimSpace(rest) != ""
}

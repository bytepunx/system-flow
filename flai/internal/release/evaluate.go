package release

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Evaluation is whether orchestration.release is met, and the figures it
// rests on (S-0217). It tags, bumps, and pushes nothing.
type Evaluation struct {
	// Policy is orchestration.release.policy, or judgement when unset.
	Policy string `json:"policy"`
	Met    bool   `json:"met"`
	// Reason says why the policy is met or not, in one sentence.
	Reason string `json:"reason"`
	// Currency is the currency of every amount here.
	Currency string `json:"currency"`
	// Value is the cost of delay per week summed over the pending stories;
	// Count is how many there are.
	Value float64 `json:"value"`
	Count int     `json:"count"`
	// ValueThreshold and CountThreshold are a threshold policy's figures.
	ValueThreshold *float64 `json:"value_threshold,omitempty"`
	CountThreshold *int     `json:"count_threshold,omitempty"`
	// Epic or Tag is a theme policy's.
	Epic string `json:"epic,omitempty"`
	Tag  string `json:"tag,omitempty"`
	// Pending are the stories accepted and not yet released.
	Pending []EvaluatedStory `json:"pending"`
	// Unvalued are the pending stories with no cost of delay value.
	Unvalued []string `json:"unvalued,omitempty"`
	// Theme are a theme policy's stories, cancelled ones left out;
	// NotAccepted are those of them not yet done.
	Theme       []EvaluatedStory `json:"theme,omitempty"`
	NotAccepted []string         `json:"not_accepted,omitempty"`
	// HeldByEpic are, with orchestration.release.whole_epics set, the
	// pending stories whose epic is in neither review nor done (S-0222).
	// While any is held the policy is not met, whichever it is; under
	// judgement they are what the orchestrator must not publish.
	HeldByEpic []HeldStory `json:"held_by_epic,omitempty"`
}

// HeldStory is a pending story whole_epics holds back, with its epic and
// the epic's status; the status is empty when the epic is not found.
type HeldStory struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Epic       string `json:"epic"`
	EpicStatus string `json:"epic_status"`
}

// EvaluatedStory is one story an Evaluation weighs.
type EvaluatedStory struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	// Value is its cost of delay per week; nil when it has none.
	Value    *float64 `json:"value,omitempty"`
	Accepted bool     `json:"accepted"`
	Released bool     `json:"released"`
}

// EvaluateRepo is Evaluate on the repository's items and what its history
// says is pending: accepted and not yet published (PendingBatch).
func EvaluateRepo(r execx.Runner, root string, m manifest.Manifest, repo *workitem.Repo) (Evaluation, error) {
	b, err := PendingBatch(r, root, m, repo)
	if err != nil {
		return Evaluation{}, err
	}
	items, err := repo.List(true)
	if err != nil {
		return Evaluation{}, err
	}
	return Evaluate(m, items, b.IDs()), nil
}

// Evaluate says whether m's orchestration.release is met by items, of which
// pending are the IDs accepted and not yet released. A threshold is met when
// the pending stories' summed value or their count is at or over a figure
// the manifest sets; a theme when every story of its epic or tag, cancelled
// ones aside, is accepted and one at least is not yet released; judgement
// never, the call being the orchestrator's or the operator's. With
// whole_epics set, none is met while a pending story's epic is in neither
// review nor done.
func Evaluate(m manifest.Manifest, items []*workitem.Item, pending map[string]bool) Evaluation {
	rel := m.Orchestration.Release
	ev := Evaluation{Policy: rel.PolicyOrDefault(), Currency: m.Planning.CurrencyCode(), Pending: []EvaluatedStory{}}
	for _, it := range items {
		if it.Type != workitem.Story || !pending[it.ID] {
			continue
		}
		s := story(it, pending)
		ev.Pending = append(ev.Pending, s)
		ev.Count++
		if s.Value == nil {
			ev.Unvalued = append(ev.Unvalued, it.ID)
			continue
		}
		ev.Value += *s.Value
	}
	switch ev.Policy {
	case manifest.ReleaseThreshold:
		ev.ValueThreshold, ev.CountThreshold = rel.Value, rel.Count
		ev.threshold()
	case manifest.ReleaseTheme:
		ev.Epic, ev.Tag = strings.TrimSpace(rel.Epic), strings.TrimSpace(rel.Tag)
		ev.theme(items, pending)
	case manifest.ReleaseJudgement:
		ev.Reason = fmt.Sprintf("the release policy is judgement, which is never met by itself: whether to release is the orchestrator's or the operator's call, with %s", ev.pendingPhrase())
	default:
		ev.Reason = fmt.Sprintf("%q is not a release policy, so it is never met; write judgement, threshold, or theme", rel.Policy)
	}
	if rel.WholeEpics {
		ev.holdByEpic(items, pending)
	}
	return ev
}

// holdByEpic sets HeldByEpic to the pending stories whose epic is in neither
// review nor done, and, when there is one, unmeets the policy and says which
// epics hold the batch back. A story with no epic is never held; one whose
// epic is not found is, as nothing says the epic is in review or done.
func (ev *Evaluation) holdByEpic(items []*workitem.Item, pending map[string]bool) {
	epics := map[string]*workitem.Item{}
	for _, it := range items {
		if it.Type == workitem.Epic {
			epics[workitem.CanonicalID(it.ID)] = it
		}
	}
	// order is the holding epics as first met; stories, each one's stories.
	var order []string
	stories := map[string][]string{}
	for _, it := range items {
		if it.Type != workitem.Story || !pending[it.ID] || strings.TrimSpace(it.Parent) == "" {
			continue
		}
		h := HeldStory{ID: it.ID, Title: it.Title, Epic: workitem.CanonicalID(it.Parent)}
		if e := epics[h.Epic]; e != nil {
			h.Epic, h.EpicStatus = e.ID, e.Status
		}
		if h.EpicStatus == workitem.Review || h.EpicStatus == workitem.Done {
			continue
		}
		ev.HeldByEpic = append(ev.HeldByEpic, h)
		if _, ok := stories[h.Epic]; !ok {
			order = append(order, h.Epic)
		}
		stories[h.Epic] = append(stories[h.Epic], it.ID)
	}
	if len(ev.HeldByEpic) == 0 {
		return
	}
	ev.Met = false
	status := map[string]string{}
	for _, h := range ev.HeldByEpic {
		status[h.Epic] = h.EpicStatus
	}
	held := make([]string, 0, len(order))
	for _, id := range order {
		state := "is not found"
		if s := status[id]; s != "" {
			state = "is " + s
		}
		held = append(held, fmt.Sprintf("%s %s, with %s", id, state, strings.Join(stories[id], ", ")))
	}
	whose := "its epic is"
	if len(order) > 1 {
		whose = "their epics are"
	}
	ev.Reason += fmt.Sprintf("; whole_epics holds the batch back until %s in review or done: %s", whose, strings.Join(held, "; "))
}

func story(it *workitem.Item, pending map[string]bool) EvaluatedStory {
	s := EvaluatedStory{ID: it.ID, Title: it.Title, Status: it.Status, Accepted: it.Status == workitem.Done}
	s.Released = s.Accepted && !pending[it.ID]
	if c := it.CostOfDelay; c != nil && c.Value != nil {
		v := *c.Value
		s.Value = &v
	}
	return s
}

// pendingPhrase is the pending stories' count and value, as a phrase.
func (ev *Evaluation) pendingPhrase() string {
	switch ev.Count {
	case 0:
		return "no accepted story waiting for a release"
	case 1:
		return fmt.Sprintf("1 accepted story not yet released, worth %s", ev.Amount(ev.Value))
	}
	return fmt.Sprintf("%d accepted stories not yet released, worth %s", ev.Count, ev.Amount(ev.Value))
}

// Amount writes v as a cost of delay per week in the evaluation's currency.
func (ev *Evaluation) Amount(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64) + " " + ev.Currency + "/week"
}

func (ev *Evaluation) threshold() {
	if ev.ValueThreshold == nil && ev.CountThreshold == nil {
		ev.Reason = "the threshold sets neither value nor count, so it is never met; write one or both under orchestration.release"
		return
	}
	if ev.Count == 0 {
		ev.Reason = "no accepted story is waiting for a release, so the threshold is not met"
		return
	}
	var reached, short []string
	if v := ev.ValueThreshold; v != nil {
		if ev.Value >= *v {
			reached = append(reached, fmt.Sprintf("their value, %s, is at or over the threshold of %s", ev.Amount(ev.Value), ev.Amount(*v)))
		} else {
			short = append(short, fmt.Sprintf("their value, %s, is under the threshold of %s", ev.Amount(ev.Value), ev.Amount(*v)))
		}
	}
	if c := ev.CountThreshold; c != nil {
		if ev.Count >= *c {
			reached = append(reached, fmt.Sprintf("their count, %d, is at or over the threshold of %d", ev.Count, *c))
		} else {
			short = append(short, fmt.Sprintf("their count, %d, is under the threshold of %d", ev.Count, *c))
		}
	}
	ev.Met = len(reached) > 0
	head := fmt.Sprintf("%d accepted stories are not yet released", ev.Count)
	if ev.Count == 1 {
		head = "1 accepted story is not yet released"
	}
	if ev.Met {
		ev.Reason = head + ", and " + strings.Join(reached, ", and ")
		return
	}
	ev.Reason = head + ", and " + strings.Join(short, ", and ")
}

func (ev *Evaluation) theme(items []*workitem.Item, pending map[string]bool) {
	if ev.Epic == "" && ev.Tag == "" {
		ev.Reason = "the theme names neither epic nor tag, so it is never met; write one under orchestration.release"
		return
	}
	what := "the tag " + ev.Tag
	match := func(it *workitem.Item) bool { return slices.Contains(it.Tags, ev.Tag) }
	if ev.Epic != "" {
		what = "the epic " + ev.Epic
		epic := workitem.CanonicalID(ev.Epic)
		match = func(it *workitem.Item) bool { return workitem.CanonicalID(it.Parent) == epic }
	}
	unreleased := 0
	for _, it := range items {
		if it.Type != workitem.Story || it.Status == workitem.Cancelled || !match(it) {
			continue
		}
		s := story(it, pending)
		ev.Theme = append(ev.Theme, s)
		switch {
		case !s.Accepted:
			ev.NotAccepted = append(ev.NotAccepted, it.ID)
		case !s.Released:
			unreleased++
		}
	}
	switch {
	case len(ev.Theme) == 0:
		ev.Reason = fmt.Sprintf("no story belongs to %s, so the theme is not met", what)
	case len(ev.NotAccepted) > 0:
		verb := "are"
		if len(ev.NotAccepted) == 1 {
			verb = "is"
		}
		ev.Reason = fmt.Sprintf("%d of the %d stories of %s %s not yet accepted: %s", len(ev.NotAccepted), len(ev.Theme), what, verb, strings.Join(ev.NotAccepted, ", "))
	case unreleased == 0:
		ev.Reason = fmt.Sprintf("every story of %s is accepted, and every one is released already", what)
	default:
		ev.Met = true
		ev.Reason = fmt.Sprintf("every story of %s is accepted, and %d of the %d are not yet released", what, unreleased, len(ev.Theme))
	}
}

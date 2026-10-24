package context

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// FromRole is the kind of the source of a strategic pack's topic that comes
// from its role rather than from an item.
const FromRole = "role"

// roleTopics are the topics each strategic role brings to its pack, in the
// order the pack lists them.
var roleTopics = map[string][]string{
	conventions.RolePlan:        {"planning"},
	conventions.RoleOrchestrate: {"orchestration"},
	conventions.RoleAnalyze:     {"analysis"},
	conventions.RoleCurate:      {"digest", "planning"},
}

// roleAgent names the agent in each strategic role.
var roleAgent = map[string]string{
	conventions.RolePlan:        "the planner",
	conventions.RoleOrchestrate: "the orchestrator",
	conventions.RoleAnalyze:     "the analyzer",
	conventions.RoleCurate:      "the curator",
}

// ForStrategic builds the pack for a strategic agent, one that works above
// a story (E-0016): the conventions whose roles are empty or list its role,
// with the sections its topics leave out taken out; the open issues; and
// briefs of the design, tech, and ADRs its topics select, with the ADRs one
// step reaches, and a catalog of the rest. Its topics are its role's
// (planning, orchestration, analysis, or, for the curator, digest and
// planning, ADR-0196). The planner, role plan, is primed for an epic or a
// story, given as epic or story: its topics are the item's as well, and the
// pack loads what the item names, as a story's pack does, and ranks sections
// against it to fill the budget. The orchestrator, the analyzer, and the
// curator are primed for the whole project and take no item. The
// orchestrator's pack adds the answers and recommendations of the digest
// that rank highest against the threads awaiting the operator, after the
// briefs (ADR-0194); the planner's, the analyzer's, and the curator's take no
// digest step.
// budget is a size as ParseSize reads it; empty means the project's
// prime.budget, else DefaultBudget.
func ForStrategic(repo *workitem.Repo, role, epic, story, budget string) (*Pack, error) {
	if !slices.Contains(conventions.StrategicRoles, role) {
		return nil, fmt.Errorf("flai prime --role %s: no such role; a strategic agent's role is %s", role, strings.Join(conventions.StrategicRoles, ", "))
	}
	cmd := "flai prime --role " + role
	it, err := strategicItem(repo, cmd, role, epic, story)
	if err != nil {
		return nil, err
	}
	size, err := Budget(budget, repo.Manifest.Prime.Budget)
	if err != nil {
		return nil, err
	}
	pr, err := loadProject(repo, cmd)
	if err != nil {
		return nil, err
	}
	var (
		itemTopics   []topics.StoryTopic
		sources      []Source
		title, query string
	)
	if it != nil {
		if itemTopics, sources, err = planned(repo, it); err != nil {
			return nil, err
		}
		title, query = it.Title, Query(it.Title, it.Body)
	}
	packTopics := withRole(itemTopics, role)
	list, _ := issues.List(repo)
	pack, err := Build(repo.Root, "", title, packTopics, pr.readBy(role, false, topics.Names(packTopics)), issues.SummaryTable(list), size)
	if err != nil {
		return nil, err
	}
	pack.Role = role
	if it != nil {
		pack.Item = it.ID
	}
	var answers []Item
	if role == conventions.RoleOrchestrate {
		answers = orchestratorDigest(repo)
	}
	pack.addDesign(Design(pr.docs, topics.Names(packTopics), sources, BriefOver(size), nil), query, answers)
	return pack, nil
}

// answerKinds are the kinds of digest entry the orchestrator's pack selects
// (ADR-0194): what the operator answered, and the recommendations they
// confirmed.
var answerKinds = []string{"answer", "recommendation"}

// orchestratorDigest is the digest items the orchestrator's pack tries, as
// answersFor orders the digest's entries against the threads awaiting the
// operator. A digest entry or a thread that does not parse is flai check's
// to report; the pack does without the digest rather than fail, as a
// story's does.
func orchestratorDigest(repo *workitem.Repo) []Item {
	entries, err := digest.List(repo)
	if err != nil {
		return nil
	}
	all, err := threads.List(repo)
	if err != nil {
		return nil
	}
	return answersFor(entries, awaitingOperator(repo, all))
}

// answersFor is the digest items the orchestrator's pack tries, best first
// (ADR-0194): the active and compacted answers and recommendations of
// entries, ranked for each thread of awaiting by BM25 against its title and
// first entry, as rankDigest ranks them, and taken round by round, each
// thread's best entry not yet taken, threads in the order given. An entry
// with no word in common with a thread is not taken for it. An active entry
// with a body loads it; a compacted one is a brief.
func answersFor(entries []digest.Entry, awaiting []*threads.Thread) []Item {
	var picks []pick
	for _, e := range entries {
		if (e.Status == digestActive || e.Status == digestCompacted) && slices.Contains(answerKinds, e.Kind) {
			picks = append(picks, pick{entry: e})
		}
	}
	if len(picks) == 0 {
		return nil
	}
	ranked := make([][]pick, len(awaiting))
	for i, th := range awaiting {
		own := slices.Clone(picks)
		scoreDigest(own, threadQuery(th))
		ranked[i] = rankDigest(own)
	}
	var out []Item
	taken := map[string]bool{}
	next := make([]int, len(awaiting))
	for more := true; more; {
		more = false
		for i, th := range awaiting {
			for next[i] < len(ranked[i]) && taken[ranked[i][next[i]].entry.ID] {
				next[i]++
			}
			if next[i] == len(ranked[i]) {
				continue
			}
			p := ranked[i][next[i]]
			next[i]++
			taken[p.entry.ID], more = true, true
			out = append(out, digestItem(p.entry, threadReason(th.ID, next[i], p.entry), true))
		}
	}
	return out
}

// threadQuery is what the digest is ranked against for a thread: its title
// and its first entry.
func threadQuery(th *threads.Thread) string {
	q := th.Title
	if entries := th.Entries(); len(entries) > 0 {
		q += "\n" + entries[0].Text
	}
	return q
}

// threadReason is the reason of an entry the orchestrator's pack takes for
// thread id at rank n, "thread TH-nnnn; rank n"; a compacted entry's ends
// "; compacted".
func threadReason(id string, n int, e digest.Entry) string {
	r := "thread " + id + "; rank " + strconv.Itoa(n)
	if e.Status == digestCompacted {
		r += "; " + digestCompacted
	}
	return r
}

// awaitingOperator is the threads of all that await the operator
// (ADR-0181), in the order given: those open, with entries, and with a
// recommendation pending the operator's confirmation or a last entry by
// none of the operator's names. It is the rule inbox.Summarize gives a
// thread its reader takes no part in; this package cannot import inbox,
// which imports it.
func awaitingOperator(repo *workitem.Repo, all []*threads.Thread) []*threads.Thread {
	var out []*threads.Thread
	for _, th := range all {
		entries := th.Entries()
		if !th.Open() || len(entries) == 0 {
			continue
		}
		if th.PendingRecommendation() != nil || !threads.OperatorName(repo, threads.StoryOf(repo, th), entries[len(entries)-1].Author) {
			out = append(out, th)
		}
	}
	return out
}

// strategicItem is the item a strategic role is primed for: the epic or the
// story the planner is given, checked to be what it was given as, and none
// for the others.
func strategicItem(repo *workitem.Repo, cmd, role, epic, story string) (*workitem.Item, error) {
	const give = "give --epic E-nnnn or --story S-nnnn"
	id, want := epic, workitem.Epic
	if story != "" {
		id, want = story, workitem.Story
	}
	switch {
	case role != conventions.RolePlan && id != "":
		return nil, fmt.Errorf("%s %s: %s primes for the whole project; give no --story or --epic", cmd, id, roleAgent[role])
	case role != conventions.RolePlan:
		return nil, nil
	case epic != "" && story != "":
		return nil, fmt.Errorf("%s: the planner plans one item; give --epic or --story, not both", cmd)
	case id == "":
		return nil, fmt.Errorf("%s: the planner plans an epic or a story; %s", cmd, give)
	}
	it, err := repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w; %s", cmd, id, err, give)
	}
	if it.Type != want {
		return nil, fmt.Errorf("%s --%s %s: %s is %s; %s", cmd, want, id, it.ID, it.Type, give)
	}
	return it, nil
}

// planned is the topics of an item the planner plans and the sources of
// what it names: a story's topics and its story, epic, and tasks, as its
// own pack has them; an epic's own topics and those its tags and touches
// reach, and the epic alone, whose stories the planner reads as items.
func planned(repo *workitem.Repo, it *workitem.Item) ([]topics.StoryTopic, []Source, error) {
	if it.Type == workitem.Story {
		t, err := topics.ForStory(repo, it.ID)
		if err != nil {
			return nil, nil, err
		}
		s, err := storySources(repo, it)
		return t, s, err
	}
	src := Source{ID: it.ID, Path: rel(repo.Root, it.Path), Body: it.Body}
	return topics.OfStory(it, nil, nil, repo.Manifest.Projects), []Source{src}, nil
}

// withRole puts the role's topics first in list, in their order, save one
// list has already, which keeps its place and gains the role as a source.
func withRole(list []topics.StoryTopic, role string) []topics.StoryTopic {
	from := topics.Source{Kind: FromRole, Via: role}
	out := slices.Clone(list)
	var front []topics.StoryTopic
	for _, name := range roleTopics[role] {
		if i := slices.IndexFunc(out, func(t topics.StoryTopic) bool { return t.Topic == name }); i >= 0 {
			out[i].Sources = append(slices.Clone(out[i].Sources), from)
			continue
		}
		front = append(front, topics.StoryTopic{Topic: name, Sources: []topics.Source{from}})
	}
	return append(front, out...)
}

// strategicLines are the header's lines naming a strategic pack's item and
// role.
func (p *Pack) strategicLines() string {
	if p.Item == "" {
		answers := ""
		if p.Role == conventions.RoleOrchestrate {
			answers = " the answers and recommendations of the digest that rank highest against the threads awaiting the operator,"
		}
		return fmt.Sprintf("role: %s, %s, for the whole project: the conventions the role reads, the open issues, briefs of the design its topics select,%s and a catalog of the rest; doc_search and doc_get find what you need.\n", p.Role, roleAgent[p.Role], answers)
	}
	return fmt.Sprintf("item: %s %s\nrole: %s, %s, for %s: the conventions the role reads, the open issues, what the item names, briefs of the design its topics select, the sections that rank highest against it, and a catalog of the rest; doc_search and doc_get find what you need.\n", p.Item, p.Title, p.Role, roleAgent[p.Role], p.Item)
}

// strategic reports whether the pack is a strategic agent's.
func (p *Pack) strategic() bool { return slices.Contains(conventions.StrategicRoles, p.Role) }

// subAgent reports whether the pack is a sub-agent's (ADR-0059).
func (p *Pack) subAgent() bool { return p.Role != "" && !p.strategic() }

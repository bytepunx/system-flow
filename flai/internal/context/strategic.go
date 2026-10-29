package context

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// FromRole is the kind of the source of a strategic pack's topic that comes
// from its role rather than from an item.
const FromRole = "role"

// roleTopics are the topics each built-in strategic role brings to its pack,
// in the order the pack lists them, for a project with no role file for it
// (ADR-0207): one made before role files, whose packs stay as they were.
var roleTopics = map[string][]string{
	conventions.RolePlan:        {"planning"},
	conventions.RoleOrchestrate: {"orchestration"},
	conventions.RoleAnalyze:     {"analysis"},
	conventions.RoleCurate:      {"digest", "planning"},
}

// roleAgent names the agent in each built-in strategic role, for a project
// with no role file for it, as roleTopics gives its topics.
var roleAgent = map[string]string{
	conventions.RolePlan:        "the planner",
	conventions.RoleOrchestrate: "the orchestrator",
	conventions.RoleAnalyze:     "the analyzer",
	conventions.RoleCurate:      "the curator",
}

// strategicRole is a strategic role as flai prime --role resolves it: its
// verb, the name of its agent, the topics it brings to its pack, and the
// kinds of item it is primed for, none meaning the whole project.
type strategicRole struct {
	verb   string
	agent  string
	topics []string
	items  roles.Kinds
}

// takes reports whether the role is primed for an item of kind.
func (r *strategicRole) takes(kind string) bool { return slices.Contains(r.items, kind) }

// projectRoles is the project's role files (ADR-0207) that give a strategic
// role, those whose verb is empty or the story's agent's or a sub-agent's
// left out, as flai check refuses them. A role file that does not parse is
// refused with where it is and what is wrong, under cmd, rather than leave
// its role out of the pack it defines.
func projectRoles(repo *workitem.Repo, cmd string) ([]roles.Role, error) {
	all, problems, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", cmd, err)
	}
	if len(problems) > 0 {
		var each []string
		for _, p := range problems {
			each = append(each, fmt.Sprintf("role file %s:%d: %s", p.Path, p.Line, p.Message))
		}
		fix := "fix it"
		if len(problems) > 1 {
			fix = "fix them"
		}
		return nil, fmt.Errorf("%s: %s; %s, as flai check reports it", cmd, strings.Join(each, "; "), fix)
	}
	var out []roles.Role
	for _, r := range all {
		if r.Role != "" && !slices.Contains(roles.ReservedRoles, r.Role) {
			out = append(out, r)
		}
	}
	return out, nil
}

// StrategicRoles is the verbs flai prime --role takes for a strategic agent
// in the project: the built-in roles' in their order (plan, orchestrate,
// analyze, curate), then, sorted, those of the role files that are not
// built-in (ADR-0207). A role file that does not parse is refused, as
// ForStrategic refuses it.
func StrategicRoles(repo *workitem.Repo) ([]string, error) {
	files, err := projectRoles(repo, "flai prime --role")
	if err != nil {
		return nil, err
	}
	return strategicVerbs(files), nil
}

// strategicVerbs is StrategicRoles for the role files given.
func strategicVerbs(files []roles.Role) []string {
	var defined []string
	for _, r := range files {
		if !slices.Contains(conventions.StrategicRoles, r.Role) && !slices.Contains(defined, r.Role) {
			defined = append(defined, r.Role)
		}
	}
	slices.Sort(defined)
	return append(slices.Clone(conventions.StrategicRoles), defined...)
}

// resolveRole is the strategic role verb names in the project: the role
// file whose role it is, the first by path when more than one claims it, as
// flai check reports; else, for a built-in role, its topics and agent from
// roleTopics and roleAgent, the planner primed for an epic or a story and
// the others for the whole project. Any other verb is refused with the
// project's strategic roles.
func resolveRole(repo *workitem.Repo, cmd, verb string) (*strategicRole, error) {
	files, err := projectRoles(repo, cmd)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Role == verb {
			agent := roleAgent[verb]
			if f.Name != "" {
				agent = "the " + f.Name
			} else if agent == "" {
				agent = "the " + verb + " agent"
			}
			return &strategicRole{verb: verb, agent: agent, topics: f.Topics, items: f.Item}, nil
		}
	}
	if !slices.Contains(conventions.StrategicRoles, verb) {
		return nil, fmt.Errorf("%s: no such role; a strategic agent's role is %s", cmd, strings.Join(strategicVerbs(files), ", "))
	}
	r := &strategicRole{verb: verb, agent: roleAgent[verb], topics: roleTopics[verb]}
	if verb == conventions.RolePlan {
		r.items = roles.Kinds{workitem.Epic, workitem.Story}
	}
	return r, nil
}

// ForStrategic builds the pack for a strategic agent, one that works above
// a story (E-0016): the conventions whose roles are empty or list its role,
// with the sections its topics leave out taken out; the open issues; and
// briefs of the design, tech, and ADRs its topics select, with the ADRs one
// step reaches, and a catalog of the rest. The role is any the project
// defines in a role file (ADR-0207), whose topics, item kinds, and agent's
// name the pack takes from it; a built-in role with no file takes its
// topics from roleTopics (planning, orchestration, analysis, or, for the
// curator, digest and planning, ADR-0196). A role whose file names item
// kinds, as the planner's names epic and story, is primed for one item of
// those kinds, given as epic or story: its topics are the item's as well,
// and the pack loads what the item names, as a story's pack does, and ranks
// sections against it to fill the budget. A role with no item kinds, as the
// orchestrator, the analyzer, and the curator are, is primed for the whole
// project and takes no item. The orchestrator's pack adds the answers and
// recommendations of the digest that rank highest against the threads
// awaiting the operator, after the briefs (ADR-0194); no other role's takes
// a digest step.
// budget is a size as ParseSize reads it; empty means the project's
// prime.budget, else DefaultBudget.
func ForStrategic(repo *workitem.Repo, role, epic, story, budget string) (*Pack, error) {
	cmd := "flai prime --role " + role
	r, err := resolveRole(repo, cmd, role)
	if err != nil {
		return nil, err
	}
	it, err := strategicItem(repo, cmd, r, epic, story)
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
	packTopics := withTopics(itemTopics, role, r.topics)
	list, _ := issues.List(repo)
	pack, err := Build(repo.Root, "", title, packTopics, pr.readBy(role, false, topics.Names(packTopics)), issues.SummaryTable(list), size)
	if err != nil {
		return nil, err
	}
	pack.Role, pack.agent = role, r.agent
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

// strategicItem is the item a strategic role is primed for: for a role that
// takes items, the epic or the story it is given, of a kind the role takes
// and checked to be what it was given as; none for a role primed for the
// whole project.
func strategicItem(repo *workitem.Repo, cmd string, r *strategicRole, epic, story string) (*workitem.Item, error) {
	takesEpic, takesStory := r.takes(workitem.Epic), r.takes(workitem.Story)
	var flags, items []string
	if takesEpic {
		flags, items = append(flags, "--epic E-nnnn"), append(items, "an epic")
	}
	if takesStory {
		flags, items = append(flags, "--story S-nnnn"), append(items, "a story")
	}
	give := "give " + strings.Join(flags, " or ")
	does := "primes for"
	if r.verb == conventions.RolePlan {
		does = "plans"
	}
	id, want := epic, workitem.Epic
	if story != "" {
		id, want = story, workitem.Story
	}
	switch {
	case !takesEpic && !takesStory && id != "":
		return nil, fmt.Errorf("%s %s: %s primes for the whole project; give no --story or --epic", cmd, id, r.agent)
	case !takesEpic && !takesStory:
		return nil, nil
	case epic != "" && !takesEpic:
		return nil, fmt.Errorf("%s --epic %s: %s takes %s, not an epic; %s", cmd, epic, r.agent, strings.Join(items, " or "), give)
	case story != "" && !takesStory:
		return nil, fmt.Errorf("%s --story %s: %s takes %s, not a story; %s", cmd, story, r.agent, strings.Join(items, " or "), give)
	case epic != "" && story != "":
		return nil, fmt.Errorf("%s: %s %s one item; give --epic or --story, not both", cmd, r.agent, does)
	case id == "":
		return nil, fmt.Errorf("%s: %s %s %s; %s", cmd, r.agent, does, strings.Join(items, " or "), give)
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

// planned is the topics of an item a strategic role is primed for, as the
// planner plans one, and the sources of what it names: a story's topics and
// its story, epic, and tasks, as its own pack has them; an epic's own topics
// and those its tags and touches reach, and the epic alone, whose stories
// the role reads as items.
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

// withRole is withTopics with a built-in role's topics, as roleTopics gives
// them.
func withRole(list []topics.StoryTopic, role string) []topics.StoryTopic {
	return withTopics(list, role, roleTopics[role])
}

// withTopics puts names, the topics role brings, first in list, in their
// order, save one list has already, which keeps its place and gains the
// role as a source.
func withTopics(list []topics.StoryTopic, role string, names []string) []topics.StoryTopic {
	from := topics.Source{Kind: FromRole, Via: role}
	out := slices.Clone(list)
	var front []topics.StoryTopic
	for _, name := range names {
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
	if len(p.Reasons) > 0 {
		return p.reasonLines()
	}
	if p.Item == "" {
		answers := ""
		if p.Role == conventions.RoleOrchestrate {
			answers = " the answers and recommendations of the digest that rank highest against the threads awaiting the operator,"
		}
		return fmt.Sprintf("role: %s, %s, for the whole project: the conventions the role reads, the open issues, briefs of the design its topics select,%s and a catalog of the rest; doc_search and doc_get find what you need.\n", p.Role, p.agent, answers)
	}
	return fmt.Sprintf("item: %s %s\nrole: %s, %s, for %s: the conventions the role reads, the open issues, what the item names, briefs of the design its topics select, the sections that rank highest against it, and a catalog of the rest; doc_search and doc_get find what you need.\n", p.Item, p.Title, p.Role, p.agent, p.Item)
}

// strategic reports whether the pack is a strategic agent's: one with a
// role and no story, as ForStrategic builds it, whatever role the project
// defines (ADR-0207).
func (p *Pack) strategic() bool { return p.Role != "" && p.Story == "" }

// subAgent reports whether the pack is a sub-agent's (ADR-0059): one with a
// role and the story it works for, as ForRole builds it.
func (p *Pack) subAgent() bool { return p.Role != "" && p.Story != "" }

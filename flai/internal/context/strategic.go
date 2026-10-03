package context

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// FromRole is the kind of the source of a strategic pack's topic that comes
// from its role rather than from an item.
const FromRole = "role"

// roleTopic is the topic each strategic role brings to its pack.
var roleTopic = map[string]string{
	conventions.RolePlan:        "planning",
	conventions.RoleOrchestrate: "orchestration",
	conventions.RoleAnalyze:     "analysis",
}

// roleAgent names the agent in each strategic role.
var roleAgent = map[string]string{
	conventions.RolePlan:        "the planner",
	conventions.RoleOrchestrate: "the orchestrator",
	conventions.RoleAnalyze:     "the analyzer",
}

// ForStrategic builds the pack for a strategic agent, one that works above
// a story (E-0016): the conventions whose roles are empty or list its role,
// with the sections its topics leave out taken out; the open issues; and
// briefs of the design, tech, and ADRs its topics select, with the ADRs one
// step reaches, and a catalog of the rest. Its topics are its role's
// (planning, orchestration, or analysis). The planner, role plan, is primed
// for an epic or a story, given as epic or story: its topics are the item's
// as well, and the pack loads what the item names, as a story's pack does,
// and ranks sections against it to fill the budget. The orchestrator and the
// analyzer are primed for the whole project and take no item. budget is a size as ParseSize reads it; empty means the
// project's prime.budget, else DefaultBudget.
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
	pack, err := Build(repo.Root, "", title, packTopics, pr.readBy(role, false), issues.SummaryTable(list), size)
	if err != nil {
		return nil, err
	}
	pack.Role = role
	if it != nil {
		pack.Item = it.ID
	}
	pack.AddDesign(Design(pr.docs, topics.Names(packTopics), sources, BriefOver(size)), query)
	return pack, nil
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

// withRole puts the role's topic first in list, or adds the role to its
// sources when list has it already.
func withRole(list []topics.StoryTopic, role string) []topics.StoryTopic {
	from := topics.Source{Kind: FromRole, Via: role}
	for i, t := range list {
		if t.Topic == roleTopic[role] {
			out := slices.Clone(list)
			out[i].Sources = append(slices.Clone(t.Sources), from)
			return out
		}
	}
	return append([]topics.StoryTopic{{Topic: roleTopic[role], Sources: []topics.Source{from}}}, list...)
}

// strategicLines are the header's lines naming a strategic pack's item and
// role.
func (p *Pack) strategicLines() string {
	if p.Item == "" {
		return fmt.Sprintf("role: %s, %s, for the whole project: the conventions the role reads, the open issues, briefs of the design its topics select, and a catalog of the rest; doc_search and doc_get find what you need.\n", p.Role, roleAgent[p.Role])
	}
	return fmt.Sprintf("item: %s %s\nrole: %s, %s, for %s: the conventions the role reads, the open issues, what the item names, briefs of the design its topics select, the sections that rank highest against it, and a catalog of the rest; doc_search and doc_get find what you need.\n", p.Item, p.Title, p.Role, roleAgent[p.Role], p.Item)
}

// strategic reports whether the pack is a strategic agent's.
func (p *Pack) strategic() bool { return slices.Contains(conventions.StrategicRoles, p.Role) }

// subAgent reports whether the pack is a sub-agent's (ADR-0059).
func (p *Pack) subAgent() bool { return p.Role != "" && !p.strategic() }

package context

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ForStory builds the context pack for a story (ADR-0047), fitted to a
// budget (ADR-0049): the conventions the story's agent reads (ADR-0068)
// with the sections its topics leave out taken out (S-0136), the open
// issues, then what the story, its epic, and its tasks name, briefs of the
// design, tech, and ADRs its topics and one link step select, the sections
// that rank highest to fill the budget, and a catalog of the rest (S-0137,
// S-0146). budget is a size as ParseSize reads it; empty means the
// project's prime.budget, else DefaultBudget. flai prime --story and the MCP
// prime tool both call it. An archived story gets the pack it would get
// today, which is how a pack is replayed.
func ForStory(repo *workitem.Repo, id, budget string) (*Pack, error) {
	size, err := Budget(budget, repo.Manifest.Prime.Budget)
	if err != nil {
		return nil, err
	}
	st, err := loadStory(repo, id)
	if err != nil {
		return nil, err
	}
	list, _ := issues.List(repo)
	pack, err := Build(repo.Root, st.item.ID, st.item.Title, st.topics, st.readBy(conventions.RoleStory, true), issues.SummaryTable(list), size)
	if err != nil {
		return nil, err
	}
	pack.AddDesign(Design(st.docs, topics.Names(st.topics), st.sources, BriefOver(size)), Query(st.item.Title, st.item.Body))
	return pack, nil
}

// ForRole builds the pack for a sub-agent of a story's agent (ADR-0059):
// the conventions it reads (ADR-0068), with the sections the story's
// topics leave out taken out; the story's goal and acceptance criteria; and
// briefs, never bodies, of what the story, its epic, and its tasks name,
// what their topics select, and the ADRs one step reaches, in the order
// Items gives them, each while the budget has room. budget is a size as
// ParseSize reads it; empty means half the story's agent's.
func ForRole(repo *workitem.Repo, id, role, budget string) (*Pack, error) {
	if !slices.Contains(conventions.SubAgentRoles, role) {
		return nil, fmt.Errorf("flai prime --role %s: no such role; a sub-agent's role is %s, and a strategic agent's %s", role, strings.Join(conventions.SubAgentRoles, " or "), strings.Join(conventions.StrategicRoles, ", "))
	}
	size, err := Budget(budget, repo.Manifest.Prime.Budget)
	if err != nil {
		return nil, err
	}
	if budget == "" {
		size /= 2
	}
	st, err := loadStory(repo, id)
	if err != nil {
		return nil, err
	}
	pack, err := Build(repo.Root, st.item.ID, st.item.Title, st.topics, st.readBy(role, false), "", size)
	if err != nil {
		return nil, err
	}
	pack.Role, pack.Goal = role, Goal(st.item.Body)
	s := NewSelection(st.docs)
	s.BriefNamed = true
	s.Linked(st.sources)
	s.ByTopics(topics.Names(st.topics))
	s.Step()
	pack.AddBriefs(s)
	return pack, nil
}

// project is what every pack is built from: the conventions and the design
// documents.
type project struct {
	set  *conventions.Set
	docs []*Doc
}

// story is what a story's pack is built from: the project, the story, its
// topics, and what links them.
type story struct {
	*project
	item    *workitem.Item
	topics  []topics.StoryTopic
	sources []Source
}

// readBy is the conventions the agent in role reads (ADR-0068): those whose
// roles are empty or list it, with the README when readme is set.
func (p *project) readBy(role string, readme bool) *conventions.Set {
	sub := &conventions.Set{Dir: p.set.Dir}
	if readme {
		sub.README = p.set.README
	}
	for _, f := range p.set.Files {
		if f.ReadBy(role) {
			sub.Files = append(sub.Files, f)
		}
	}
	return sub
}

// loadProject reads the conventions and the design documents; cmd is the
// command the errors name.
func loadProject(repo *workitem.Repo, cmd string) (*project, error) {
	set, _, err := conventions.Load(repo)
	if err != nil {
		return nil, err
	}
	if set.Missing {
		return nil, fmt.Errorf("no conventions folder at %s; render it from the template or run flai upgrade", rel(repo.Root, set.Dir))
	}
	docs, err := LoadDocs(repo.Root, repo.Manifest.Dir(repo.Root, "design"))
	if err != nil {
		return nil, fmt.Errorf("%s: reading the design documents: %w", cmd, err)
	}
	return &project{set: set, docs: docs}, nil
}

func loadStory(repo *workitem.Repo, id string) (*story, error) {
	it, err := repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("flai prime --story %s: %w; give the ID of a story", id, err)
	}
	if it.Type != workitem.Story {
		return nil, fmt.Errorf("flai prime --story %s: %s is %s, not a story", id, it.ID, it.Type)
	}
	pr, err := loadProject(repo, "flai prime --story "+id)
	if err != nil {
		return nil, err
	}
	storyTopics, err := topics.ForStory(repo, it.ID)
	if err != nil {
		return nil, err
	}
	sources, err := storySources(repo, it)
	if err != nil {
		return nil, err
	}
	return &story{project: pr, item: it, topics: storyTopics, sources: sources}, nil
}

// storySources is the story, its epic, and its tasks, archived or not, for
// the documents they link.
func storySources(repo *workitem.Repo, story *workitem.Item) ([]Source, error) {
	src := func(it *workitem.Item) Source {
		return Source{ID: it.ID, Path: rel(repo.Root, it.Path), Body: it.Body}
	}
	out := []Source{src(story)}
	if story.Parent != "" {
		if epic, err := repo.Get(story.Parent); err == nil {
			out = append(out, src(epic))
		}
	}
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		if it.Type == workitem.Task && it.Parent == story.ID {
			out = append(out, src(it))
		}
	}
	return out, nil
}

// rel is p relative to root, or p when it is not under root.
func rel(root, p string) string {
	if r, err := filepath.Rel(root, p); err == nil {
		return r
	}
	return p
}

package context

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ForStory builds the context pack for a story (ADR-0047), fitted to a
// budget (ADR-0049): the conventions the story's agent reads (ADR-0068)
// with the sections its topics leave out taken out (S-0136), the open
// issues, then what the story, its epic, and its tasks name, and the design
// system when its guidance's topics reach the story (ADR-0152), briefs of the
// design, tech, and ADRs its topics and one link step select, the digest
// entries its topics or its claim select, ranked against it within their
// share (ADR-0194), the sections that rank highest to fill the budget, and a
// catalog of the rest (S-0137, S-0146). budget is a size as ParseSize reads
// it; empty means the project's prime.budget, else DefaultBudget. flai prime --story and the MCP
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
	names := topics.Names(st.topics)
	pack, err := Build(repo.Root, st.item.ID, st.item.Title, st.topics, st.readBy(conventions.RoleStory, true, names), issues.SummaryTable(list), size)
	if err != nil {
		return nil, err
	}
	query := Query(st.item.Title, st.item.Body)
	pack.addDesign(Design(st.docs, names, st.sources, BriefOver(size), st.design), query, digestFor(st.digest, names, st.claim, query, true))
	return pack, nil
}

// ForRole builds the pack for a sub-agent of a story's agent (ADR-0059):
// the conventions it reads (ADR-0068), with the sections the story's
// topics leave out taken out; the story's goal and acceptance criteria; and
// briefs, never bodies, of what the story, its epic, and its tasks name,
// the design system its guidance's topics reach (ADR-0152), what their
// topics select, and the ADRs one step reaches, in the order Items gives
// them, then the digest entries the story's would load, ranked as its are
// (ADR-0194), each while the budget has room. budget is a size as
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
	names := topics.Names(st.topics)
	pack, err := Build(repo.Root, st.item.ID, st.item.Title, st.topics, st.readBy(role, false, names), "", size)
	if err != nil {
		return nil, err
	}
	pack.Role, pack.Goal = role, Goal(st.item.Body)
	s := NewSelection(st.docs)
	s.BriefNamed = true
	s.Linked(st.sources)
	s.DesignSystem(st.design, names)
	s.ByTopics(names)
	s.Step()
	pack.AddBriefs(s, digestFor(st.digest, names, st.claim, Query(st.item.Title, st.item.Body), false))
	return pack, nil
}

// project is what every pack is built from: the conventions, the design
// documents, and the design system.
type project struct {
	set    *conventions.Set
	docs   []*Doc
	design *DesignSystem
}

// story is what a story's pack is built from: the project, the story, its
// topics, what links them, its claim, and the digest.
type story struct {
	*project
	item    *workitem.Item
	topics  []topics.StoryTopic
	sources []Source
	claim   []string       // the story's touches and its tasks', as paths
	digest  []digest.Entry // every entry, whatever its status
}

// readBy is the conventions the agent in role reads (ADR-0068): those whose
// roles are empty or list it, with the README when readme is set, and
// without the design-system convention unless the design system reaches a
// pack with these topics (ADR-0152).
func (p *project) readBy(role string, readme bool, packTopics []string) *conventions.Set {
	sub := &conventions.Set{Dir: p.set.Dir}
	if readme {
		sub.README = p.set.README
	}
	for _, f := range p.set.Files {
		if f.ReadBy(role) {
			sub.Files = append(sub.Files, f)
		}
	}
	if len(p.design.Reaches(packTopics)) == 0 {
		return sub.WithoutDesignSystem()
	}
	return sub
}

// loadProject reads the conventions, the design documents, and the design
// system; cmd is the command the errors name.
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
	ds, err := loadDesignSystem(repo.Root, repo.Manifest, docs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", cmd, err)
	}
	return &project{set: set, docs: docs, design: ds}, nil
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
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	// A digest entry that does not parse is flai check's to report; the pack
	// does without the digest rather than fail, as it does without the issues.
	entries, _ := digest.List(repo)
	return &story{project: pr, item: it, topics: storyTopics, sources: sourcesIn(repo, it, items), claim: storyClaim(repo, it, items), digest: entries}, nil
}

// storySources is the story, its epic, and its tasks, archived or not, for
// the documents they link.
func storySources(repo *workitem.Repo, story *workitem.Item) ([]Source, error) {
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	return sourcesIn(repo, story, items), nil
}

// sourcesIn is storySources with the story's tasks found in items.
func sourcesIn(repo *workitem.Repo, story *workitem.Item, items []*workitem.Item) []Source {
	src := func(it *workitem.Item) Source {
		return Source{ID: it.ID, Path: rel(repo.Root, it.Path), Body: it.Body}
	}
	out := []Source{src(story)}
	if story.Parent != "" {
		if epic, err := repo.Get(story.Parent); err == nil {
			out = append(out, src(epic))
		}
	}
	for _, it := range items {
		if it.Type == workitem.Task && it.Parent == story.ID {
			out = append(out, src(it))
		}
	}
	return out
}

// storyClaim is what a digest entry's paths are met against (ADR-0194): the
// story's touches and those of its tasks in items that are not cancelled,
// archived or not, each read as a path as a claim reads it, a component's
// name or tag as its path, without duplicates.
func storyClaim(repo *workitem.Repo, story *workitem.Item, items []*workitem.Item) []string {
	touches := slices.Clone(story.Touches)
	for _, it := range items {
		if it.Type == workitem.Task && it.Parent == story.ID && it.Status != workitem.Cancelled {
			touches = append(touches, it.Touches...)
		}
	}
	return workitem.NewHolds(nil, repo.Manifest.Projects).Claim(&workitem.Item{Touches: touches})
}

// rel is p relative to root, or p when it is not under root.
func rel(root, p string) string {
	if r, err := filepath.Rel(root, p); err == nil {
		return r
	}
	return p
}

package context

import (
	"fmt"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ForStory builds the context pack for a story (ADR-0047), fitted to a
// budget (ADR-0049): the conventions its topics select (S-0136), the open
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
	it, err := repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("flai prime --story %s: %w; give the ID of a story", id, err)
	}
	if it.Type != workitem.Story {
		return nil, fmt.Errorf("flai prime --story %s: %s is %s, not a story", id, it.ID, it.Type)
	}
	set, _, err := conventions.Load(repo)
	if err != nil {
		return nil, err
	}
	if set.Missing {
		return nil, fmt.Errorf("no conventions folder at %s; render it from the template or run flai upgrade", rel(repo.Root, set.Dir))
	}
	list, _ := issues.List(repo)
	storyTopics, err := topics.ForStory(repo, it.ID)
	if err != nil {
		return nil, err
	}
	pack, err := Build(repo.Root, it.ID, it.Title, storyTopics, set, issues.SummaryTable(list), size)
	if err != nil {
		return nil, err
	}
	docs, err := LoadDocs(repo.Root, repo.Manifest.Dir(repo.Root, "design"))
	if err != nil {
		return nil, fmt.Errorf("flai prime --story %s: reading the design documents: %w", id, err)
	}
	sources, err := storySources(repo, it)
	if err != nil {
		return nil, err
	}
	pack.AddDesign(Design(docs, topics.Names(storyTopics), sources, BriefOver(size)), Query(it.Title, it.Body))
	return pack, nil
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

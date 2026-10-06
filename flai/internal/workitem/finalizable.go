package workitem

import (
	"fmt"
	"strings"
	"time"
)

// Finalizable drafts (S-0219). A draft story is complete when it has every
// section of the project's story template, a goal, at least one acceptance
// criterion as a checkbox, at least one touch, a forecast duration and
// delivery, a cost of delay value, and an open parent epic or none. That is
// arithmetic, and flai says it; whether the criteria, touches, and forecast
// describe the same work is the orchestrator's judgement.

// Drafts are the draft stories in the backlog, in the order a policy gives
// them, each complete or with what it lacks.
type Drafts struct {
	Policy string       `json:"policy"`
	Drafts []DraftCheck `json:"drafts"`
}

// DraftCheck is one draft story and what it lacks to be finalized.
type DraftCheck struct {
	// Position counts from 1, in the policy's order.
	Position int    `json:"position"`
	ID       string `json:"id"`
	Title    string `json:"title"`
	// Text is the figure the draft was ordered by, as a person reads it,
	// when it has one.
	Text     string   `json:"text,omitempty"`
	Complete bool     `json:"complete"`
	Lacks    []string `json:"lacks"`
}

// FinalizableLacks is each thing story lacks to be finalized. parent is its
// epic, nil when it names none or the epic is not found, and sections are the
// level-two headings of the project's story template. It does not look at
// the draft flag, so a story can be checked as it would be finalized.
func FinalizableLacks(story, parent *Item, sections []string) []string {
	lacks := []string{}
	have := map[string]bool{}
	for _, h := range levelTwoHeadings(story.Body) {
		have[h] = true
	}
	for _, h := range sections {
		if !have[h] {
			lacks = append(lacks, fmt.Sprintf("no %q section", "## "+h))
		}
	}
	if !hasGoal(story.Body) {
		lacks = append(lacks, "no goal")
	}
	if !hasCriteria(story.Body) {
		lacks = append(lacks, "no acceptance criteria with a checkbox")
	}
	if len(story.Touches) == 0 {
		lacks = append(lacks, "no touches")
	}
	if _, ok := forecastHours(story); !ok {
		lacks = append(lacks, "no forecast duration")
	}
	if !hasDelivery(story) {
		lacks = append(lacks, "no forecast delivery")
	}
	if codValue(story) == nil {
		lacks = append(lacks, "no cost of delay value")
	}
	switch {
	case story.Parent == "":
	case parent == nil:
		lacks = append(lacks, "its epic "+story.Parent+" is not found")
	case parent.Closed():
		lacks = append(lacks, "its epic "+parent.ID+" is "+parent.Status)
	}
	return lacks
}

// Finalizable is each thing story lacks to be finalized, judged against the
// project's story template and its parent epic as they are on disk. None
// means it is complete.
func (r *Repo) Finalizable(story *Item) ([]string, error) {
	if story.Type != Story {
		return nil, fmt.Errorf("%s is %s: only a story is finalized", story.ID, articled(story.Type))
	}
	sections, err := r.storySections()
	if err != nil {
		return nil, err
	}
	parent, err := r.parentOf(story, nil)
	if err != nil {
		return nil, err
	}
	return FinalizableLacks(story, parent, sections), nil
}

// Drafts lists each draft story in the backlog, complete or with what it
// lacks, ordered by orchestration.policy from the board's pull order.
func (r *Repo) Drafts() (Drafts, error) {
	items, err := r.List(false)
	if err != nil {
		return Drafts{}, err
	}
	board, err := r.LoadBoard()
	if err != nil {
		return Drafts{}, err
	}
	sections, err := r.storySections()
	if err != nil {
		return Drafts{}, err
	}
	policy := r.Manifest.Orchestration.PolicyOrDefault()
	byID := map[string]*Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	var drafts []*Item
	for _, id := range PullSequence(board.Order, items, Backlog) {
		if it := byID[id]; it.Draft {
			drafts = append(drafts, it)
		}
	}
	ranked, err := OrderByPolicy(drafts, policy)
	if err != nil {
		return Drafts{}, err
	}
	out := Drafts{Policy: policy, Drafts: []DraftCheck{}}
	for _, rk := range ranked {
		story := byID[rk.ID]
		parent, err := r.parentOf(story, byID)
		if err != nil {
			return Drafts{}, err
		}
		lacks := FinalizableLacks(story, parent, sections)
		out.Drafts = append(out.Drafts, DraftCheck{
			Position: rk.Position, ID: rk.ID, Title: rk.Title, Text: rk.Text,
			Complete: len(lacks) == 0, Lacks: lacks,
		})
	}
	return out, nil
}

// parentOf is story's epic, from active when it is there, else from disk,
// archive included; nil when it names none or the epic is not found.
func (r *Repo) parentOf(story *Item, active map[string]*Item) (*Item, error) {
	if story.Parent == "" {
		return nil, nil
	}
	if p, ok := active[story.Parent]; ok {
		return p, nil
	}
	if TypeOfID(CanonicalID(story.Parent)) == "" {
		return nil, nil
	}
	p, err := r.Get(story.Parent)
	if err != nil {
		if strings.HasSuffix(err.Error(), " not found") {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

// storySections are the level-two headings of the project's story template.
func (r *Repo) storySections() ([]string, error) {
	body, err := r.TemplateBody(Story)
	if err != nil {
		return nil, err
	}
	return levelTwoHeadings(body), nil
}

// levelTwoHeadings are the "## " headings of a markdown body, without the
// marks, in order.
func levelTwoHeadings(body string) []string {
	var out []string
	for _, l := range strings.Split(body, "\n") {
		if h, ok := strings.CutPrefix(l, "## "); ok {
			if h = strings.TrimSpace(h); h != "" {
				out = append(out, h)
			}
		}
	}
	return out
}

// hasDelivery says whether the story's forecast has a delivery that is a UTC
// timestamp.
func hasDelivery(it *Item) bool {
	if it.Forecast == nil || it.Forecast.Delivery == "" {
		return false
	}
	_, err := time.Parse(TimeFormat, it.Forecast.Delivery)
	return err == nil
}

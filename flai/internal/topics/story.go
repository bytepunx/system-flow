package topics

import (
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Code is the topic of every story that reaches a sub-project whose kind is
// not template (ADR-0047).
const Code = "code"

// Where a story's topic came from (ADR-0047, S-0135).
const (
	FromOwn   = "own"   // the story's own topics:
	FromEpic  = "epic"  // its epic's topics:
	FromTag   = "tag"   // a sub-project one of its tags names
	FromClaim = "claim" // a sub-project its claim reaches: its touches or an open task's
	FromCode  = "code"  // it reaches a sub-project that is not the template
	FromAll   = "all"   // every story has all
)

// Source is one reason a story has a topic.
type Source struct {
	Kind    string `json:"kind"`              // own, epic, tag, claim, code, all
	Item    string `json:"item,omitempty"`    // the story, epic, or task that says it
	Project string `json:"project,omitempty"` // the sub-project reached, for tag, claim, and code
	Via     string `json:"via,omitempty"`     // the tag or the touches entry that reached it
}

// String is the source in a few words: "tag cli reaches flai".
func (s Source) String() string {
	switch s.Kind {
	case FromOwn:
		return "own"
	case FromEpic:
		return "epic " + s.Item
	case FromTag:
		return fmt.Sprintf("tag %s reaches %s", s.Via, s.Project)
	case FromClaim:
		return fmt.Sprintf("%s touches %s, in %s", s.Item, s.Via, s.Project)
	case FromCode:
		return s.Project + " is code"
	case FromAll:
		return "every story"
	}
	return s.Kind
}

// StoryTopic is one of a story's topics and every source of it.
type StoryTopic struct {
	Topic   string   `json:"topic"`
	Sources []Source `json:"sources"`
}

// OfStory works out a story's topics (ADR-0047): its own and its epic's
// topics; the name, tags, and kind of every sub-project its tags name or its
// claim reaches; code when one of those is not the template; and all. The
// claim is the story's touches and those of its tasks that are still open
// (ADR-0046); an entry reaches a sub-project when it is the sub-project's
// name or one of its tags, or when it and the sub-project's path overlap.
// epic may be nil, and tasks may hold any items: only the story's open tasks
// count. Topics come in that order, each once, with every source.
func OfStory(story, epic *workitem.Item, tasks []*workitem.Item, projects []manifest.Project) []StoryTopic {
	var out []StoryTopic
	at := map[string]int{}
	add := func(topic string, s Source) {
		if topic == "" {
			return
		}
		i, ok := at[topic]
		if !ok {
			i = len(out)
			at[topic] = i
			out = append(out, StoryTopic{Topic: topic})
		}
		for _, have := range out[i].Sources {
			if have == s {
				return // a sub-project's name, tag, and kind may be one word
			}
		}
		out[i].Sources = append(out[i].Sources, s)
	}
	for _, t := range story.Topics {
		add(t, Source{Kind: FromOwn, Item: story.ID})
	}
	if epic != nil {
		for _, t := range epic.Topics {
			add(t, Source{Kind: FromEpic, Item: epic.ID})
		}
	}

	claim := touchesOf(story)
	for _, t := range tasks {
		if t.Type == workitem.Task && t.Parent == story.ID && !t.Archived && !t.Closed() {
			claim = append(claim, touchesOf(t)...)
		}
	}
	var code []Source
	for _, p := range projects {
		var reached []Source
		for _, tag := range story.Tags {
			if names(p, tag) {
				reached = append(reached, Source{Kind: FromTag, Item: story.ID, Project: p.Name, Via: tag})
			}
		}
		for _, c := range claim {
			if names(p, c.entry) || workitem.PathsOverlap(c.entry, cleanPath(p.Path)) {
				reached = append(reached, Source{Kind: FromClaim, Item: c.item, Project: p.Name, Via: c.entry})
			}
		}
		if len(reached) == 0 {
			continue
		}
		for _, w := range append(append([]string{p.Name}, p.Tags...), p.Kind) {
			for _, s := range reached {
				add(w, s)
			}
		}
		if p.Kind != "template" {
			code = append(code, Source{Kind: FromCode, Project: p.Name})
		}
	}
	for _, s := range code {
		add(Code, s)
	}
	add(All, Source{Kind: FromAll})
	return out
}

// touch is one touches entry of a claim and the item that declares it.
type touch struct{ item, entry string }

// touchesOf is an item's touches entries as a claim compares them.
func touchesOf(it *workitem.Item) []touch {
	var out []touch
	for _, e := range it.Touches {
		if e = cleanPath(e); e != "" {
			out = append(out, touch{it.ID, e})
		}
	}
	return out
}

// names says whether a tag or touches entry is the sub-project's name or one
// of its tags.
func names(p manifest.Project, w string) bool {
	if w == p.Name {
		return true
	}
	for _, t := range p.Tags {
		if w == t {
			return true
		}
	}
	return false
}

func cleanPath(p string) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(p), "./"), "/")
}

// ForStory reads a story, its epic, and its tasks from repo and works out
// its topics. The story and its epic may be archived.
func ForStory(repo *workitem.Repo, id string) ([]StoryTopic, error) {
	story, err := repo.Get(id)
	if err != nil {
		return nil, err
	}
	if story.Type != workitem.Story {
		return nil, fmt.Errorf("%s is %s; topics are worked out for stories", story.ID, story.Type)
	}
	var epic *workitem.Item
	if story.Parent != "" {
		if e, err := repo.Get(story.Parent); err == nil {
			epic = e
		}
	}
	items, err := repo.List(false) // an archived task is closed and claims nothing
	if err != nil {
		return nil, err
	}
	return OfStory(story, epic, items, repo.Manifest.Projects), nil
}

// Summary is a topic's sources in one line, each sub-project once: "own;
// flai by tag cli and 3 touches (S-0135, T-0496)". --json keeps every source.
func (t StoryTopic) Summary() string {
	type reach struct {
		tags      []string
		touches   []string
		items     []string
		seenTouch map[string]bool
		seenItem  map[string]bool
	}
	var parts, order []string
	by := map[string]*reach{}
	get := func(p string) *reach {
		if by[p] == nil {
			by[p] = &reach{seenTouch: map[string]bool{}, seenItem: map[string]bool{}}
			order = append(order, p)
		}
		return by[p]
	}
	var code []string
	for _, s := range t.Sources {
		switch s.Kind {
		case FromTag:
			r := get(s.Project)
			r.tags = append(r.tags, "tag "+s.Via)
		case FromClaim:
			r := get(s.Project)
			if !r.seenTouch[s.Via] {
				r.seenTouch[s.Via] = true
				r.touches = append(r.touches, s.Via)
			}
			if !r.seenItem[s.Item] {
				r.seenItem[s.Item] = true
				r.items = append(r.items, s.Item)
			}
		case FromCode:
			code = append(code, s.Project)
		default:
			parts = append(parts, s.String())
		}
	}
	for _, p := range order {
		r := by[p]
		how := r.tags
		switch len(r.touches) {
		case 0:
		case 1:
			how = append(how, "touches "+r.touches[0]+" ("+strings.Join(r.items, ", ")+")")
		default:
			how = append(how, fmt.Sprintf("%d touches (%s)", len(r.touches), strings.Join(r.items, ", ")))
		}
		parts = append(parts, p+" by "+strings.Join(how, " and "))
	}
	switch len(code) {
	case 0:
	case 1:
		parts = append(parts, code[0]+" is code")
	default:
		parts = append(parts, strings.Join(code, ", ")+" are code")
	}
	return strings.Join(parts, "; ")
}

// Names are the topics alone, in order.
func Names(list []StoryTopic) []string {
	out := make([]string, len(list))
	for i, t := range list {
		out[i] = t.Topic
	}
	return out
}

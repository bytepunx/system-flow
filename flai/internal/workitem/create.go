package workitem

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

//go:embed defaults/*.tmpl
var defaults embed.FS

// NewOptions describe an item to create.
type NewOptions struct {
	Type    string
	Title   string
	Nature  string
	Parent  string
	Owner   string
	Tags    []string
	Touches []string
	// Topics are what a story or epic is about (S-0135, ADR-0047).
	Topics []string
	// After names what must be done before the item starts: stories for a
	// story (ADR-0046), tasks of the same story for a task (S-0176). That
	// each exists and no cycle forms is flai check's.
	After []string
	// Agent is who works a story, over the project's default (S-0103): what
	// it sets wins, and the project's default fills in the rest.
	Agent *manifest.Agent
	// Draft makes a story a draft (S-0199): an agent wrote it, and it is
	// finalized before it is ready.
	Draft bool
	// CostOfDelay is what waiting for a story or an epic costs (S-0203): a
	// story made from an issue carries the issue's cost of delay inputs.
	CostOfDelay *CostOfDelay
	Now         time.Time
	// Body replaces what the template puts below the item's heading: the
	// author's goal, criteria, and notes, written before the item exists
	// (S-0059). The heading stays the one flai renders, so the ID and the
	// title in it cannot disagree with the front matter.
	Body string
}

// Create allocates an ID, renders the body template, links the parent, and
// writes both files. It returns the new item.
func (r *Repo) Create(opt NewOptions) (*Item, error) {
	if !contains(Types, opt.Type) {
		return nil, fmt.Errorf("unknown type %q", opt.Type)
	}
	if opt.Title = CleanTitle(opt.Title); opt.Title == "" {
		return nil, fmt.Errorf("a title is required")
	}
	if opt.Nature == "" {
		opt.Nature = "feature"
	}
	if !contains(Natures, opt.Nature) {
		return nil, fmt.Errorf("nature %q must be one of %s", opt.Nature, strings.Join(Natures, ", "))
	}
	var parent *Item
	switch {
	case opt.Type == Epic:
		if opt.Parent != "" {
			return nil, fmt.Errorf("epics have no parent")
		}
	case opt.Type == Story && opt.Parent == "":
		// a story need not belong to an epic (S-0092): not every story fits
		// an active one, and an epic created only to hold one is not wanted.
	default:
		wantParent := Epic
		if opt.Type == Task {
			wantParent = Story
		}
		if opt.Parent == "" {
			return nil, fmt.Errorf("a %s needs a parent %s (--%s)", opt.Type, wantParent, wantParent)
		}
		p, err := r.Get(opt.Parent)
		if err != nil {
			return nil, err
		}
		if p.Type != wantParent {
			return nil, fmt.Errorf("%s is a %s, a %s's parent must be a %s", p.ID, p.Type, opt.Type, wantParent)
		}
		if p.Archived || p.Closed() {
			return nil, fmt.Errorf("%s is %s; reopen or pick another parent", p.ID, p.Status)
		}
		parent = p
	}
	id, err := r.NextID(opt.Type)
	if err != nil {
		return nil, err
	}
	now := opt.Now.UTC().Format(TimeFormat)
	data := map[string]any{
		"id": id, "title": opt.Title, "nature": opt.Nature, "parent": opt.Parent,
		"owner": orDefault(opt.Owner, "agent"), "now": now, "today": now[:10],
	}
	doc, err := r.renderItemTemplate(opt.Type, data)
	if err != nil {
		return nil, err
	}
	it, err := ParseItem(doc)
	if err != nil {
		return nil, fmt.Errorf("item template for %s produced an invalid item: %w", opt.Type, err)
	}
	if body := strings.TrimSpace(opt.Body); body != "" {
		it.Body = itemHeading(it.Body) + "\n\n" + body + "\n"
	}
	it.Tags = append(it.Tags, opt.Tags...)
	it.Touches = append(it.Touches, opt.Touches...)
	topics, err := CleanTopics(opt.Topics)
	if err != nil {
		return nil, err
	}
	if len(topics) > 0 && opt.Type == Task {
		return nil, fmt.Errorf("only stories and epics carry topics, not a task")
	}
	it.Topics = append(it.Topics, topics...)
	if it.After, err = CleanAfter(opt.Type, id, opt.After); err != nil {
		return nil, err
	}
	if opt.Type == Story {
		it.Agent = r.Manifest.Agent.With(opt.Agent)
	} else if !opt.Agent.IsZero() {
		return nil, fmt.Errorf("only a story carries an agent, not a %s", opt.Type)
	}
	if opt.Draft && opt.Type != Story {
		return nil, fmt.Errorf("only a story is a draft, not %s", articled(opt.Type))
	}
	it.Draft = opt.Draft
	if c := opt.CostOfDelay; !c.IsZero() {
		if !Carries(opt.Type, "cost_of_delay") {
			return nil, fmt.Errorf("only a story or an epic carries a cost of delay, not %s", articled(opt.Type))
		}
		// checked here, so that a wrong one is not blamed on the template below
		if errs := planningErrors(&Item{Type: opt.Type, CostOfDelay: c}); len(errs) > 0 {
			return nil, fmt.Errorf("the new %s's cost of delay: %s", opt.Type, strings.Join(errs, "; "))
		}
		it.CostOfDelay = c.clone()
	}
	if it.Tags == nil {
		it.Tags = []string{}
	}
	if it.Transitions == nil {
		it.Transitions = []Transition{}
	}
	if err := it.Validate(); err != nil {
		return nil, fmt.Errorf("item template for %s produced an invalid item: %w", opt.Type, err)
	}
	it.Path = filepath.Join(r.ItemDir(it.Type, it.Archived), FileName(it.ID, it.Title))
	if err := r.LintGuard(it.Path, "", it.Marshal()); err != nil {
		return nil, err
	}
	if err := r.Save(it); err != nil {
		return nil, err
	}
	if parent != nil {
		AppendChild(parent, it)
		parent.Updated = now
		if err := r.Save(parent); err != nil {
			return nil, err
		}
	}
	return it, nil
}

// headingPunctuation is what markdownlint rejects at the end of a heading
// by default (MD026).
const headingPunctuation = ".,;:!。，；：！"

// CleanTitle is a title as flai writes it: one line, single spaces, and no
// trailing punctuation, which the lint rejects in the item's heading (MD026,
// S-0179). A question mark stays.
func CleanTitle(s string) string {
	return strings.TrimRight(strings.Join(strings.Fields(s), " "), headingPunctuation+" ")
}

// LintGuard refuses a write to path, a file in the project, whose markdown
// the project's lint rejects where before did not (S-0179): what flai writes
// in the main checkout reaches main without the lint a story runs. The error
// is an *mdlint.Error naming each finding's rule and line.
func (r *Repo) LintGuard(path, before, after string) error {
	root := r.MainRoot
	if root == "" {
		root = r.Root
	}
	rel := path
	if x, err := filepath.Rel(root, path); err == nil && !strings.HasPrefix(x, "..") {
		rel = filepath.ToSlash(x)
	}
	return mdlint.Guard(root, rel, before, after)
}

// CleanTopics trims a topics list, drops empty entries and repeats, and
// refuses an entry that is not one word.
func CleanTopics(in []string) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, t := range in {
		t = strings.TrimSpace(t)
		if t == "" || seen[t] {
			continue
		}
		if !ValidTopic(t) {
			return nil, fmt.Errorf("topic %q is not one word: use letters, digits, dot, dash, or underscore, such as logging or release", t)
		}
		seen[t] = true
		out = append(out, t)
	}
	return out, nil
}

// CleanAfter is an after list as flai writes it, for the item self of type
// typ: canonical IDs, without empty entries or repeats. It refuses an epic's,
// an entry that is not an ID of the item's own type, and the item itself.
// That each exists, a task's is of the same story, and no cycle forms is
// flai check's. Nothing given is nil.
func CleanAfter(typ, self string, in []string) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if !Carries(typ, "after") {
			return nil, fmt.Errorf("only a story or a task waits for others in after, not %s", articled(typ))
		}
		id := CanonicalID(v)
		if !idPattern.MatchString(id) || !strings.HasPrefix(id, strings.ToUpper(typ[:1])+"-") {
			what := "stories, such as S-0001"
			if typ == Task {
				what = "tasks of the same story, such as T-0001"
			}
			return nil, fmt.Errorf("after %q: name %s", v, what)
		}
		if id == self {
			return nil, fmt.Errorf("after %s: %s cannot wait for itself", id, articled(typ))
		}
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out, nil
}

// itemHeading is the first heading of a rendered item body, "# ID Title".
func itemHeading(body string) string {
	for _, l := range strings.Split(body, "\n") {
		if strings.HasPrefix(l, "# ") {
			return strings.TrimRight(l, " ")
		}
	}
	return ""
}

// TemplateBody is what the project's template puts below the heading of a
// new item of this type: the sections an author fills in. A form starts
// from it, so a project that changed its item template gets its own
// sections and not flai's defaults.
func (r *Repo) TemplateBody(typ string) (string, error) {
	if !contains(Types, typ) {
		return "", fmt.Errorf("unknown type %q", typ)
	}
	doc, err := r.renderItemTemplate(typ, map[string]any{
		"id": "X-0000", "title": "title", "nature": "feature", "parent": "",
		"owner": "owner", "now": "2000-01-01T00:00:00Z", "today": "2000-01-01",
	})
	if err != nil {
		return "", err
	}
	it, err := ParseItem(doc)
	if err != nil {
		return "", err
	}
	body := it.Body
	if h := itemHeading(body); h != "" {
		body = body[strings.Index(body, h)+len(h):]
	}
	return strings.TrimLeft(body, "\n"), nil
}

// renderItemTemplate uses the project's template items/ when available,
// else the embedded defaults.
func (r *Repo) renderItemTemplate(name string, data map[string]any) (string, error) {
	var text string
	if r.TemplateDir != "" {
		if m, err := template.LoadManifest(r.TemplateDir); err == nil {
			if rel, ok := m.Items[name]; ok {
				if b, err := os.ReadFile(filepath.Join(r.TemplateDir, rel)); err == nil {
					text = string(b)
				}
			}
		}
	}
	if text == "" {
		b, err := defaults.ReadFile("defaults/" + name + ".md.tmpl")
		if err != nil {
			return "", fmt.Errorf("no template for %s", name)
		}
		text = string(b)
	}
	return template.RenderText("item:"+name, text, data)
}

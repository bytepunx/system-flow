package workitem

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
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
	// Agent is who works a story, over the project's default (S-0103): what
	// it sets wins, and the project's default fills in the rest.
	Agent *manifest.Agent
	Now   time.Time
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
	if strings.TrimSpace(opt.Title) == "" {
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
	if opt.Type == Story {
		it.Agent = r.Manifest.Agent.With(opt.Agent)
	} else if !opt.Agent.IsZero() {
		return nil, fmt.Errorf("only a story carries an agent, not a %s", opt.Type)
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

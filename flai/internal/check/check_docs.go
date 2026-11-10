package check

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/adrindex"
	"github.com/bytepunx/system-flow/flai/internal/analysis"
	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/experiment"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

type docFront struct {
	Title   string `yaml:"title"`
	Updated string `yaml:"updated"`
	ID      string `yaml:"id"`
	Status  string `yaml:"status"`
	Date    string `yaml:"date"`
}

// ADRStatuses are the statuses an ADR may carry (ADR-0223).
var ADRStatuses = []string{"proposed", "accepted", "superseded", "deprecated", "rejected"}

var adrName = regexp.MustCompile(`^(\d{4})-.*\.md$`)

// adrIndex keeps design/adrs/README.md what its ADR files render
// (ADR-0228): an index that differs is one warning, on its first line that
// differs, naming the rows that differ, are missing, or link no file. A
// project without the folder or the index has nothing to keep, and an ADR
// whose front matter does not read is doc.front-matter's to report.
func (c *checker) adrIndex() {
	dir := filepath.Join(c.repo.Manifest.Dir(c.repo.Root, "design"), "adrs")
	index := filepath.Join(dir, "README.md")
	data, err := os.ReadFile(index)
	if err != nil {
		return
	}
	rendered, err := adrindex.Render(os.DirFS(dir), string(data))
	if err != nil {
		return // the unreadable ADR is doc.front-matter's to report
	}
	d := adrindex.Diff(string(data), rendered)
	if d.Line == 0 {
		return
	}
	var what []string
	for _, p := range []struct {
		nums        []string
		one, plural string
	}{
		{d.Changed, "row %s differs", "rows %s differ"},
		{d.Missing, "row %s is missing", "rows %s are missing"},
		{d.Extra, "row %s links no ADR file", "rows %s link no ADR file"},
	} {
		switch len(p.nums) {
		case 0:
		case 1:
			what = append(what, fmt.Sprintf(p.one, p.nums[0]))
		default:
			what = append(what, fmt.Sprintf(p.plural, andList(p.nums)))
		}
	}
	if len(what) == 0 {
		what = []string{"the order or layout of its table differs"}
	}
	c.add(Warning, "adr.index", index, d.Line, "the ADR index differs from what its ADR files render: %s; run flai check --fix or flai adr index, which rewrite it", andList(what))
}

func (c *checker) documentation() {
	for _, key := range []string{"design", "docs"} {
		root := c.repo.Manifest.Dir(c.repo.Root, key)
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil //nolint:nilerr // an unreadable entry is skipped, the walk continues
			}
			if d.IsDir() {
				if (d.Name() == conventions.Folder || d.Name() == issues.Folder || (key == "design" && (d.Name() == experiment.Folder || d.Name() == analysis.Folder || d.Name() == digest.Folder))) && filepath.Dir(path) == root {
					return filepath.SkipDir // validated by their own rules
				}
				return nil
			}
			if !strings.HasSuffix(path, ".md") || d.Name() == "README.md" {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil //nolint:nilerr // unreadable files are skipped, not fatal
			}
			fm, _, err := workitem.SplitFrontMatter(string(data))
			if err != nil {
				c.add(Warning, "doc.front-matter", path, 1, "no front matter; every design and docs file needs title and updated")
				return nil //nolint:nilerr // recorded as a finding
			}
			var f docFront
			if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
				c.add(Warning, "doc.front-matter", path, 1, "front matter does not parse: %v", err)
				return nil //nolint:nilerr // recorded as a finding
			}
			isADR := filepath.Base(filepath.Dir(path)) == "adrs" && adrName.MatchString(d.Name())
			if key == "design" {
				c.docTopics(path, string(data), isADR)
			}
			if isADR {
				num := adrName.FindStringSubmatch(d.Name())[1]
				if f.ID != "ADR-"+num {
					c.add(Warning, "adr.id", path, keyLine(path, "id"), "id should be ADR-%s to match the file name", num)
				}
				if f.Status == "" || f.Date == "" || f.Title == "" {
					c.add(Warning, "adr.front-matter", path, 1, "ADRs need id, title, status, and date")
				}
				// The template, 0000, carries the status template.
				if f.Status != "" && num != "0000" && !slices.Contains(ADRStatuses, f.Status) {
					c.add(Warning, "adr.status", path, keyLine(path, "status"), "status %q is not one of %s (ADR-0223)", f.Status, strings.Join(ADRStatuses, ", "))
				}
				if _, ok := ctxpack.DecisionSentence(string(data)); !ok {
					c.add(Warning, "adr.decision", path, headingLine(path, "## Decision"), "## Decision does not open with a sentence stating the decision; a context pack briefs the ADR by that sentence (ADR-0049)")
				}
				return nil
			}
			if f.Title == "" {
				c.add(Warning, "doc.title", path, 1, "front matter needs a title")
			}
			if f.Updated == "" {
				c.add(Warning, "doc.updated", path, 1, "front matter needs an updated date")
			}
			return nil
		})
	}
}

// Code is the topic of every sub-project that is not the template (ADR-0047).
const Code = topics.Code

// Vocabulary is every topic a story can have: all, code, each sub-project's
// name, tags, and kind (ADR-0047), and the topics that the stories and epics
// among items declare (S-0135); and the words a pack primed for the
// orchestrator's reasons reads on a heading beside all, the kinds of reason
// and the built-in strategic roles' topics (ADR-0217).
func Vocabulary(repo *workitem.Repo, items []*workitem.Item) map[string]bool {
	v := map[string]bool{topics.All: true, Code: true}
	for _, w := range ctxpack.SectionWords() {
		v[w] = true
	}
	for _, it := range items {
		for _, t := range it.Topics {
			v[t] = true
		}
	}
	for _, p := range repo.Manifest.Projects {
		for _, w := range append([]string{p.Name, p.Kind}, p.Tags...) {
			if w != "" {
				v[w] = true
			}
		}
	}
	return v
}

// docTopics checks the topics of a design file: design/system and
// design/tech files declare them, and every topic on the file or a heading
// is one a story can have. ADRs may leave them out.
func (c *checker) docTopics(path, content string, isADR bool) {
	rel, _ := filepath.Rel(c.repo.Manifest.Dir(c.repo.Root, "design"), path)
	folder := strings.SplitN(filepath.ToSlash(rel), "/", 2)[0]
	if !isADR && folder != "system" && folder != "tech" {
		return
	}
	doc, err := topics.Parse(content, nil)
	if err != nil {
		c.add(Warning, "doc.topics", path, keyLine(path, "topics"), "%v", err)
		return
	}
	if !doc.Declared && !isADR {
		c.add(Warning, "doc.topics", path, 1, "no topics; say which stories it is for, such as topics: [all] or a sub-project's name, tag, or kind (ADR-0047)")
	}
	c.topicWords(path, doc)
}

// topicWords warns on each topic of a document that no story can have.
func (c *checker) topicWords(path string, doc *topics.Doc) {
	if c.vocab == nil {
		c.vocab = Vocabulary(c.repo, c.items)
	}
	unknown := func(line int, where string, list []string) {
		for _, t := range list {
			if !c.vocab[t] {
				c.add(Warning, "doc.topic", path, line, "topic %q %s is used by nothing: not all, code, a sub-project's name, tag, or kind in system-flow.yaml, nor a story's or epic's topics (%s)", t, where, strings.Join(c.words(), ", "))
			}
		}
	}
	if doc.Declared {
		unknown(keyLine(path, "topics"), "in the front matter", doc.Topics)
	}
	for _, s := range doc.Sections {
		unknown(s.Line, "on heading \""+s.Heading+"\"", s.Own)
	}
}

func (c *checker) words() []string {
	var out []string
	for w := range c.vocab {
		out = append(out, w)
	}
	sort.Strings(out)
	return out
}

var (
	criteriaHeading = regexp.MustCompile(`(?m)^## Acceptance criteria\s*$`)
	checkbox        = regexp.MustCompile(`(?m)^\s*- \[[ xX]\] \S`)
	unchecked       = regexp.MustCompile(`(?m)^\s*- \[ \] \S`)
)

func criteriaSection(body string) string {
	loc := criteriaHeading.FindStringIndex(body)
	if loc == nil {
		return ""
	}
	rest := body[loc[1]:]
	if i := strings.Index(rest, "\n## "); i >= 0 {
		rest = rest[:i]
	}
	return rest
}

func hasCriteria(body string) bool  { return checkbox.MatchString(criteriaSection(body)) }
func hasUnchecked(body string) bool { return unchecked.MatchString(criteriaSection(body)) }

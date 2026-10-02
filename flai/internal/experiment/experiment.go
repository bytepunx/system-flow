// Package experiment knows an experiment story's results document: one file
// per experiment under design/experiments, named for the story, without which
// the story is not accepted (ADR-0066).
package experiment

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Nature is the work item nature whose stories record results here.
const Nature = "experiment"

// Folder is the results folder's name inside the design folder.
const Folder = "experiments"

// Dir is the results folder, relative to the project root, slash separated.
func Dir(m manifest.Manifest) string {
	return path.Join(m.Layout["design"], Folder)
}

// Needs reports whether an item is a story that needs a results document.
func Needs(it *workitem.Item) bool {
	return it.Type == workitem.Story && it.Nature == Nature
}

// Expected is the results document an experiment story is told to write:
// <design>/experiments/<S-nnnn>-<slug>.md, named as the story's file is.
func Expected(m manifest.Manifest, it *workitem.Item) string {
	return path.Join(Dir(m), workitem.FileName(it.ID, it.Title))
}

// Names reports whether a file name in the results folder is a results
// document for the story id: <id>.md or <id>-<slug>.md.
func Names(name, id string) bool {
	return name == id+".md" || (strings.HasPrefix(name, id+"-") && strings.HasSuffix(name, ".md"))
}

// Find returns the first of names, file names in the results folder, that is
// a results document for the story id.
func Find(names []string, id string) (string, bool) {
	for _, n := range names {
		if Names(n, id) {
			return n, true
		}
	}
	return "", false
}

// Missing is the reason an experiment story without its results document is
// not accepted, naming the document expected.
func Missing(m manifest.Manifest, it *workitem.Item, where string) string {
	return it.ID + " is an experiment and has no results document on " + where + ": write " + Expected(m, it) + " with its hypothesis, success measure, what was done, results, and recommendation (adopt, adapt, or drop), commit it, and accept again (ADR-0066)"
}

// Sections are the level-two headings every results document has.
var Sections = []string{"Hypothesis", "Success measure", "What was done", "Results", "Recommendation"}

// Recommendations are the words a recommendation is one of.
var Recommendations = []string{"adopt", "adapt", "drop"}

// Front is a results document's front matter.
type Front struct {
	Title   string `yaml:"title"`
	Updated string `yaml:"updated"`
	Status  string `yaml:"status"`
	Story   string `yaml:"story"`
}

// Problem is one way a results document falls short: Field names the front
// matter key it is about, Heading the section, or neither for the file.
type Problem struct {
	Field   string
	Heading string
	Message string
}

var fileID = regexp.MustCompile(`^(S-\d+)(-.*)?\.md$`)

// Skipped reports whether a file in the results folder is not a results
// document: the folder's index and the template a new one is copied from.
func Skipped(name string) bool { return name == "README.md" || name == "template.md" }

// Validate checks a results document named name with content: its front
// matter, that the file is named for the story it names, and its sections.
func Validate(name, content string) (Front, []Problem) {
	var f Front
	var out []Problem
	m := fileID.FindStringSubmatch(name)
	if m == nil {
		out = append(out, Problem{Message: "a results document is named <S-nnnn>-<slug>.md for its experiment story"})
	}
	fm, body, err := workitem.SplitFrontMatter(content)
	if err != nil {
		return f, append(out, Problem{Message: "no front matter; a results document needs title, updated, status, and story"})
	}
	if err := yaml.Unmarshal([]byte(fm), &f); err != nil {
		return f, append(out, Problem{Message: fmt.Sprintf("front matter does not parse: %v", err)})
	}
	for _, k := range []struct{ name, value string }{{"title", f.Title}, {"updated", f.Updated}, {"status", f.Status}, {"story", f.Story}} {
		if k.value == "" {
			out = append(out, Problem{Field: k.name, Message: "front matter needs " + k.name})
		}
	}
	if m != nil && f.Story != "" && f.Story != m[1] {
		out = append(out, Problem{Field: "story", Message: fmt.Sprintf("story is %s but the file is named for %s", f.Story, m[1])})
	}
	sections := sectionBodies(body)
	for _, h := range Sections {
		text, ok := sections[h]
		switch {
		case !ok:
			out = append(out, Problem{Message: "no ## " + h + " section"})
		case strings.TrimSpace(text) == "":
			out = append(out, Problem{Heading: "## " + h, Message: "## " + h + " is empty"})
		}
	}
	if text, ok := sections["Recommendation"]; ok && strings.TrimSpace(text) != "" && !recommends(text) {
		out = append(out, Problem{Heading: "## Recommendation", Message: "## Recommendation says none of " + strings.Join(Recommendations, ", ")})
	}
	return f, out
}

// sectionBodies maps each level-two heading of body to the text under it.
func sectionBodies(body string) map[string]string {
	out := map[string]string{}
	current, inFence := "", false
	for _, l := range strings.Split(body, "\n") {
		if strings.HasPrefix(l, "```") {
			inFence = !inFence
		}
		if !inFence && strings.HasPrefix(l, "## ") {
			current = strings.TrimSpace(strings.TrimPrefix(l, "## "))
			out[current] = ""
			continue
		}
		if current != "" {
			out[current] += l + "\n"
		}
	}
	return out
}

var word = regexp.MustCompile(`[A-Za-z]+`)

// recommends reports whether text says adopt, adapt, or drop.
func recommends(text string) bool {
	for _, w := range word.FindAllString(text, -1) {
		if slices.Contains(Recommendations, strings.ToLower(w)) {
			return true
		}
	}
	return false
}

// Package experiment knows an experiment story's results document: one file
// per experiment under design/experiments, named for the story, without which
// the story is not accepted (ADR-0066).
package experiment

import (
	"path"
	"strings"

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

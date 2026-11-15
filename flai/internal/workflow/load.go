package workflow

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// Rel is the path, relative to a project's root, of its workflow file under
// its design folder designDir, whether or not the project has it.
func Rel(designDir string) string {
	return path.Join(designDir, File)
}

// LoadError is why flai cannot run a project's workflow file (ADR-0307 §2):
// the first error-level finding on it, by its code, its path relative to the
// project's root, its line, and its message, and how many more error-level
// findings the file has.
type LoadError struct {
	Path    string
	Line    int
	Code    string
	Message string
	More    int
}

// Error is the finding as flai check reports it, and what to do.
func (e *LoadError) Error() string {
	more := ""
	if e.More > 0 {
		more = fmt.Sprintf(" (and %d more; flai check lists them all)", e.More)
	}
	return fmt.Sprintf("%s:%d: %s: %s%s; fix the file, or remove it to run the baseline workflow flai carries", e.Path, e.Line, e.Code, e.Message, more)
}

// Load reads the workflow file of the project at root whose design folder is
// designDir (ADR-0307 §1): nil and no error when the project has none, so
// that it runs the baseline; a *LoadError when the file cannot be read, does
// not parse, or has an error-level Finding; else the workflow, whose Path is
// the file's path relative to root, warnings and all.
func Load(root, designDir string) (*Workflow, error) {
	rel := Rel(designDir)
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil, nil
	case err != nil:
		return nil, &LoadError{Path: rel, Line: 1, Code: CodeFile, Message: fmt.Sprintf("%s cannot be read: %v; make it readable", rel, err)}
	}
	return judgeFile(rel, string(data))
}

// judgeFile is the workflow file at rel with content, or the *LoadError that
// keeps flai from running it.
func judgeFile(rel, content string) (*Workflow, error) {
	w, prob := Parse(rel, content)
	if prob != nil {
		msg := prob.Message
		if prob.Field != "" {
			msg = prob.Field + ": " + msg
		}
		return nil, &LoadError{Path: prob.Path, Line: prob.Line, Code: prob.Code, Message: msg}
	}
	if errs := errorsOf(w.Findings()); len(errs) > 0 {
		f := errs[0]
		return nil, &LoadError{Path: rel, Line: f.LineIn(content), Code: f.Code, Message: f.Message, More: len(errs) - 1}
	}
	return &w, nil
}

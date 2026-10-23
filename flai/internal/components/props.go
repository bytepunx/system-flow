package components

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// script reads the props of the components named on its standard input; see
// its own header for what it answers.
//
//go:embed props.mjs
var script string

// scriptHash is script's SHA-256, which keys its cache, so that a flai with
// another script resolves the skeletons again.
var scriptHash = func() string {
	sum := sha256.Sum256([]byte(script))
	return hex.EncodeToString(sum[:])
}()

// PropsCacheDir is the folder where the script keeps the skeletons it
// resolved for the sub-project at project, relative to root:
// .flai-cache/component-props/<project>/ (S-0478).
func PropsCacheDir(root, project string) string {
	return filepath.Join(root, ".flai-cache", "component-props", filepath.FromSlash(project))
}

// Timeout is how long ReadProps lets node read one sub-project's props.
const Timeout = time.Minute

// Kind is the kind of value a prop takes, read from its type.
type Kind string

// The kinds of prop. A type flai cannot see whole, such as one imported from
// another module, is KindUnknown, with its type as written.
const (
	KindString   Kind = "string"
	KindNumber   Kind = "number"
	KindBoolean  Kind = "boolean"
	KindUnion    Kind = "union"
	KindArray    Kind = "array"
	KindObject   Kind = "object"
	KindFunction Kind = "function"
	KindSnippet  Kind = "snippet"
	KindUnknown  Kind = "unknown"
)

// Prop is one prop a component takes.
type Prop struct {
	// Name is the prop's name.
	Name string `json:"name"`
	// Kind is the kind of value it takes. A union of string, number, or
	// boolean literals is KindUnion; undefined and null are left out of a
	// union before it is classified, so 'a' | undefined is a union and
	// number | null a number.
	Kind Kind `json:"kind"`
	// Type is its type as the source writes it, empty when it has none.
	Type string `json:"type,omitempty"`
	// Default is its default as the source writes it, empty when it has
	// none.
	Default string `json:"default,omitempty"`
	// Required is set when the type does not mark it optional and it has no
	// default.
	Required bool `json:"required"`
	// Literals are a union's members, each as JSON: "a", 1, true.
	Literals []string `json:"literals,omitempty"`
	// Skeleton is an empty value of its type, as JSON, for a prop of kind
	// array, object, or unknown only: [], {}, or null when its type cannot be
	// resolved.
	Skeleton json.RawMessage `json:"skeleton,omitempty"`
}

// ErrNodeMissing is wrapped by the error ReadProps returns when node is not
// on PATH.
var ErrNodeMissing = errors.New("node is not on PATH")

// ErrSvelteMissing is wrapped by the error ReadProps returns when a
// sub-project cannot resolve svelte/compiler.
var ErrSvelteMissing = errors.New("svelte is not installed")

// Runner runs node. It is execx.Runner's pattern with a context, for the
// timeout, and standard output apart from standard error, so that a warning
// node prints cannot spoil the answer.
type Runner interface {
	// LookPath reports where name is on PATH.
	LookPath(name string) (string, error)
	// Run runs name with args in dir, feeding input on standard input, and
	// returns its standard output and standard error.
	Run(ctx context.Context, dir, name, input string, args ...string) (stdout, stderr []byte, err error)
}

// System is the real Runner.
type System struct{}

// LookPath reports where name is on PATH.
func (System) LookPath(name string) (string, error) { return exec.LookPath(name) }

// Run runs name with args in dir, killing it when ctx ends.
func (System) Run(ctx context.Context, dir, name, input string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)
	cmd.WaitDelay = time.Second
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	return out.Bytes(), errOut.Bytes(), err
}

// List is Find and then ReadProps: the components of the sub-project at
// project, relative to root, with their props.
func List(ctx context.Context, r Runner, root, project string, include, exclude []string) ([]Component, error) {
	comps, err := Find(root, project, include, exclude)
	if err != nil {
		return nil, err
	}
	return ReadProps(ctx, r, root, comps)
}

// ReadProps is comps with their props, read by running the embedded script
// with node once for each sub-project, in its folder under root, within
// Timeout, keeping the skeletons it resolves under PropsCacheDir for the
// next run. A component whose props cannot be read, such as one that does not
// parse, keeps no props and carries the reason in Error. The error is for
// what stops every component of a sub-project: node not on PATH
// (ErrNodeMissing), svelte not installed there (ErrSvelteMissing), a svelte
// older than 5, or node failing or taking too long.
func ReadProps(ctx context.Context, r Runner, root string, comps []Component) ([]Component, error) {
	if len(comps) == 0 {
		return comps, nil
	}
	if _, err := r.LookPath("node"); err != nil {
		return nil, fmt.Errorf("%w, and flai reads components' props with it: install Node.js (https://nodejs.org), then run this again", ErrNodeMissing)
	}
	out := make([]Component, len(comps))
	copy(out, comps)
	var order []string
	byProject := map[string][]int{}
	for i, c := range out {
		if _, ok := byProject[c.Project]; !ok {
			order = append(order, c.Project)
		}
		byProject[c.Project] = append(byProject[c.Project], i)
	}
	for _, project := range order {
		if err := readProject(ctx, r, root, project, out, byProject[project]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// answer is what the script prints.
type answer struct {
	Components []struct {
		Path  string `json:"path"`
		Props []Prop `json:"props"`
		Error string `json:"error"`
	} `json:"components"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// readProject reads the props of out's components at idx, all of project.
func readProject(ctx context.Context, r Runner, root, project string, out []Component, idx []int) error {
	paths := make([]string, len(idx))
	for i, n := range idx {
		paths[i] = out[n].Path
	}
	input, err := json.Marshal(paths)
	if err != nil {
		return fmt.Errorf("cannot list the components of %s for node: %w", project, err)
	}
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	dir := filepath.Join(root, filepath.FromSlash(project))
	stdout, stderr, err := r.Run(ctx, dir, "node", string(input), "--input-type=module", "-e", script, PropsCacheDir(root, project), scriptHash)
	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		return fmt.Errorf("node did not finish reading the props of %d components in %s in time; run this again, or narrow the sub-project's components include", len(idx), project)
	case err != nil:
		return fmt.Errorf("node failed reading the props of the components in %s: %w%s", project, err, detail(stderr))
	}
	var a answer
	if err := json.Unmarshal(stdout, &a); err != nil {
		return fmt.Errorf("node's answer for the components in %s is not the JSON flai expects: %w%s", project, err, detail(stderr))
	}
	if a.Error != nil {
		switch a.Error.Code {
		case "svelte-missing":
			return fmt.Errorf("%w in %s, and flai reads components' props with its own svelte/compiler: run npm install (or its package manager's install) in %s, then run this again", ErrSvelteMissing, project, project)
		case "svelte-old":
			return fmt.Errorf("%s has svelte %s, and flai reads props with the parser of Svelte 5 or later: upgrade svelte there, then run this again", project, a.Error.Message)
		default:
			return fmt.Errorf("node could not read the props of the components in %s: %s (%s)", project, a.Error.Message, a.Error.Code)
		}
	}
	got := make(map[string]int, len(a.Components))
	for i, c := range a.Components {
		got[c.Path] = i
	}
	for _, n := range idx {
		i, ok := got[out[n].Path]
		if !ok {
			out[n].Error = "node gave no answer for it"
			continue
		}
		c := a.Components[i]
		out[n].Props, out[n].Error = c.Props, c.Error
		if out[n].Props == nil {
			out[n].Props = []Prop{}
		}
	}
	return nil
}

// detail is node's standard error as the end of an error message, or
// nothing when it printed none.
func detail(stderr []byte) string {
	s := strings.TrimSpace(string(stderr))
	if s == "" {
		return ""
	}
	return "\n" + s
}

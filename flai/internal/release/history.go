package release

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// commit is what Pending uses of one commit. A hash's commit never changes,
// so what is read of it is kept for as long as the process runs.
type commit struct {
	parents []string
	subject string
	ids     map[string]bool // the IDs git log --fixed-strings --grep=[<id>] finds it by
	files   []string
	listed  bool // files is read; a merge lists none
}

// history is HEAD, the tags, and the commits reachable from HEAD, as Pending
// needs them (S-0157). It is kept per repository root between calls, and
// read again only in the part that changed: while HEAD and the tags are
// unchanged, Pending costs one git process, not one per accepted item and
// commit.
type history struct {
	key     string // git show-ref --head --tags -d: HEAD and every tag
	head    string
	tags    []string          // every tag, as git tag --list names them
	tagAt   map[string]string // tag → the commit it points at
	commits map[string]*commit
	order   []string // the commits reachable from head, newest first, as git log lists them
	reach   map[string]bool
	// Worked out from the above, and so kept with it.
	ranges    map[string][]string // boundary commit → IDs accepted after it up to head, oldest first
	templates map[string]string   // template.yaml path and version → the commit git log -S finds
}

// kept holds each repository's history. Pending holds its lock throughout,
// which also guards the commits a history shares with the one before it.
var kept = struct {
	sync.Mutex
	roots map[string]*history
}{roots: map[string]*history{}}

// logFormat separates commits with RS and fields with US; the files
// --name-only lists follow the last US.
const logFormat = "--format=%x1e%H%x1f%P%x1f%s%x1f%B%x1f"

var bracketedID = regexp.MustCompile(`^[EST]-\d+$`)

// readHistory returns root's history, reading from git only what changed
// since the last call: one process when nothing did, two when commits were
// added, and the whole history when there is none kept or it was rewritten.
// The caller holds kept.
func readHistory(r execx.Runner, root string) (*history, error) {
	key, err := r.Run(root, "git", "show-ref", "--head", "--tags", "-d")
	if err != nil {
		return nil, err
	}
	prev := kept.roots[root]
	if prev != nil && prev.key == key {
		return prev, nil
	}
	h := &history{key: key, tagAt: map[string]string{}, commits: map[string]*commit{}, ranges: map[string][]string{}, templates: map[string]string{}}
	for _, l := range strings.Split(key, "\n") {
		hash, ref, _ := strings.Cut(l, " ")
		switch {
		case ref == "HEAD":
			h.head = hash
		case strings.HasPrefix(ref, "refs/tags/"):
			name := strings.TrimPrefix(ref, "refs/tags/")
			if tag, peeled := strings.CutSuffix(name, "^{}"); peeled {
				h.tagAt[tag] = hash
			} else if _, seen := h.tagAt[name]; !seen {
				h.tagAt[name] = hash
				h.tags = append(h.tags, name)
			}
		}
	}
	if h.head == "" {
		return nil, fmt.Errorf("git show-ref --head names no HEAD in %s", root)
	}
	if prev != nil {
		h.commits = prev.commits
		if h.follow(r, root, prev) {
			kept.roots[root] = h
			return h, nil
		}
	}
	out, err := r.Run(root, "git", "log", logFormat, h.head)
	if err != nil {
		return nil, err
	}
	h.walk(h.parse(out, false))
	kept.roots[root] = h
	return h, nil
}

// follow reads the commits new since prev and carries over what prev worked
// out that they cannot change. It reports false when prev's commits do not
// reach far enough for that: the history was rewritten, or HEAD went where
// prev never was.
func (h *history) follow(r execx.Runner, root string, prev *history) bool {
	var fresh []string
	touched := map[string]bool{}
	merged := false
	if _, known := h.commits[h.head]; !known {
		out, err := r.Run(root, "git", "log", logFormat, "--name-only", "--no-renames", "--root", h.head, "--not", prev.head)
		if err != nil {
			return false
		}
		fresh = h.parse(out, true)
		for _, hash := range fresh {
			c := h.commits[hash]
			for _, p := range c.parents {
				if h.commits[p] == nil {
					return false
				}
			}
			merged = merged || len(c.parents) > 1
			for _, f := range c.files {
				touched[f] = true
			}
		}
	}
	if !h.walk(append(fresh, prev.order...)) {
		return false
	}
	switch {
	case h.head == prev.head:
		h.ranges, h.templates = prev.ranges, prev.templates
	case h.reach[prev.head] && !merged:
		// git log -S finds the newest commit that changed the version line;
		// commits that did not touch the file cannot be it.
		for k, v := range prev.templates {
			if path, _, _ := strings.Cut(k, "\x00"); !touched[path] {
				h.templates[k] = v
			}
		}
	}
	return true
}

// parse adds the commits of a git log in logFormat, keeping what is already
// known of any, and returns their hashes in the log's order. listed says the
// log carried --name-only.
func (h *history) parse(out string, listed bool) []string {
	var hashes []string
	for _, rec := range strings.Split(out, "\x1e") {
		f := strings.SplitN(rec, "\x1f", 5)
		if len(f) < 4 {
			continue
		}
		hash := strings.TrimSpace(f[0])
		hashes = append(hashes, hash)
		if h.commits[hash] != nil {
			continue
		}
		c := &commit{parents: strings.Fields(f[1]), subject: f[2], ids: bracketed(f[3]), listed: listed}
		if listed && len(f) == 5 {
			for _, l := range strings.Split(f[4], "\n") {
				if l = strings.TrimSpace(l); l != "" {
					c.files = append(c.files, l)
				}
			}
		}
		h.commits[hash] = c
	}
	return hashes
}

// bracketed is every [E-n], [S-n], or [T-n] in a message: an ID is among
// them exactly when git log --fixed-strings --grep=[<id>] matches it, since
// an ID holds no bracket.
func bracketed(msg string) map[string]bool {
	ids := map[string]bool{}
	for i := 0; i < len(msg); i++ {
		if msg[i] != '[' {
			continue
		}
		end := strings.IndexByte(msg[i+1:], ']')
		if end < 0 {
			break
		}
		if tok := msg[i+1 : i+1+end]; bracketedID.MatchString(tok) {
			ids[tok] = true
		}
	}
	return ids
}

// walk sets reach to the commits reachable from head and order to those of
// candidates, and reports whether candidates held all of them. A parent not
// known (a shallow clone's edge) ends the walk there, as it ends git log.
func (h *history) walk(candidates []string) bool {
	h.reach = map[string]bool{}
	stack := []string{h.head}
	for len(stack) > 0 {
		hash := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		c := h.commits[hash]
		if c == nil || h.reach[hash] {
			continue
		}
		h.reach[hash] = true
		stack = append(stack, c.parents...)
	}
	h.order = h.order[:0]
	for _, hash := range candidates {
		if h.reach[hash] {
			h.order = append(h.order, hash)
		}
	}
	return len(h.order) == len(h.reach)
}

// currentVersion is CurrentVersion from the tags already read.
func (h *history) currentVersion(root string, p manifest.Project) (Version, error) {
	if p.Kind == "template" {
		return templateVersion(root, p)
	}
	return highestTag(h.tags, p.Name), nil
}

// boundary is the commit, or the ref, marking a component's last publish,
// or "" when it has never been released: everything accepted since is
// pending. A code component's tag is that ref. A template has no local tag
// (Tag creates none for it; its version lives in template.yaml, published to
// its own remote); the commit that set its current version stands in, and
// when none is found the whole history is walked rather than the batch
// failed over it.
func (h *history) boundary(r execx.Runner, root string, p manifest.Project, cur Version) (string, error) {
	if cur == (Version{}) {
		return "", nil
	}
	if p.Kind != "template" {
		tag := p.Name + "/v" + cur.String()
		if c, ok := h.tagAt[tag]; ok {
			return c, nil
		}
		return tag, nil
	}
	file := filepath.Join(p.Path, "template.yaml")
	key := filepath.ToSlash(file) + "\x00" + cur.String()
	if c, ok := h.templates[key]; ok {
		return c, nil
	}
	out, err := r.Run(root, "git", "log", "-n1", "--format=%H", "-S", "version: "+cur.String(), h.head, "--", file)
	if err != nil {
		return "", err
	}
	h.templates[key] = strings.TrimSpace(out)
	return h.templates[key], nil
}

// accepted lists the items a "chore: [ID] accept and archive" commit names
// after boundary up to head, each once, oldest first.
func (h *history) accepted(r execx.Runner, root, boundary string) ([]string, error) {
	if ids, ok := h.ranges[boundary]; ok {
		return ids, nil
	}
	inRange := func(string) bool { return true }
	if boundary != "" {
		if h.commits[boundary] != nil {
			before := map[string]bool{}
			stack := []string{boundary}
			for len(stack) > 0 {
				hash := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if c := h.commits[hash]; c != nil && !before[hash] {
					before[hash] = true
					stack = append(stack, c.parents...)
				}
			}
			inRange = func(hash string) bool { return !before[hash] }
		} else {
			// A boundary outside what was read (a tag off HEAD's history):
			// git knows its ancestors, so git says what lies after it.
			out, err := r.Run(root, "git", "log", "--format=%H", boundary+".."+h.head)
			if err != nil {
				return nil, err
			}
			after := map[string]bool{}
			for _, l := range strings.Fields(out) {
				after[l] = true
			}
			inRange = func(hash string) bool { return after[hash] }
		}
	}
	ids := []string{}
	seen := map[string]bool{}
	for i := len(h.order) - 1; i >= 0; i-- {
		hash := h.order[i]
		if !inRange(hash) {
			continue
		}
		if m := acceptedSubject.FindStringSubmatch(h.commits[hash].subject); m != nil && !seen[m[1]] {
			seen[m[1]] = true
			ids = append(ids, m[1])
		}
	}
	h.ranges[boundary] = ids
	return ids, nil
}

// itemCommits is Commits from the history: the commits reachable from head
// whose message names the item, newest first.
func (h *history) itemCommits(id string) []string {
	var hashes []string
	for _, hash := range h.order {
		if h.commits[hash].ids[id] {
			hashes = append(hashes, hash)
		}
	}
	return hashes
}

// list reads the files of the commits not yet listed, with one git diff-tree
// for all of them.
func (h *history) list(r execx.Runner, root string, hashes []string) error {
	want := map[string]bool{}
	var input strings.Builder
	for _, hash := range hashes {
		if c := h.commits[hash]; !c.listed && !want[hash] {
			want[hash] = true
			input.WriteString(hash + "\n")
		}
	}
	if len(want) == 0 {
		return nil
	}
	out, err := r.RunInput(root, "git", input.String(), "diff-tree", "--stdin", "--name-only", "-r", "--root")
	if err != nil {
		return err
	}
	var c *commit
	for _, l := range strings.Split(out, "\n") {
		if want[l] {
			c = h.commits[l]
			continue
		}
		if l != "" && c != nil {
			c.files = append(c.files, l)
		}
	}
	for hash := range want {
		h.commits[hash].listed = true
	}
	return nil
}

// files is TouchedFiles from the history, for commits already listed.
func (h *history) files(hashes []string) []string {
	seen := map[string]bool{}
	var files []string
	for _, hash := range hashes {
		for _, f := range h.commits[hash].files {
			if !seen[f] {
				seen[f] = true
				files = append(files, f)
			}
		}
	}
	sort.Strings(files)
	return files
}

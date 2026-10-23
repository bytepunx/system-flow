package manifest

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

// The sub-projects under projects, and the edit of them that flai subproject
// makes (S-0491): ProjectErrors says what is wrong with the list, and
// EditProjects adds, replaces, or removes one entry, keeping every other line
// and comment, and refuses a manifest that ProjectErrors, ComponentErrors, or
// InstallErrors would find wrong after it.

// exampleProjectPath is the folder a refusal of a sub-project's path gives as
// one to write when the sub-project's name is not one.
const exampleProjectPath = "web"

// projectKinds are the kinds of sub-project flai knows today, as a refusal of
// an empty kind lists them; kinds are open, so another is not an error.
const projectKinds = "go, node, python, rust, sveltekit, svelte, or template"

// ProjectErrors are what is wrong with the projects list, one sentence each,
// each beginning with the sub-project's index and field, such as
// projects[1].path: a name that is empty or repeats another's, a path that is
// empty, the repository root, absolute, not clean, leaves the repository, or
// repeats another's, and a kind that is empty. When fsys, the repository
// root, is not nil, it also reports a path that names no folder there. A kind
// flai does not know is not an error. None when every entry is valid.
func (m Manifest) ProjectErrors(fsys fs.FS) []string {
	var errs []string
	names := map[string]int{}
	paths := map[string]int{}
	for i, p := range m.Projects {
		at := func(field string) string { return projectField(i, p, field, "") }
		name, dir := strings.TrimSpace(p.Name), strings.TrimSpace(p.Path)
		switch prev, seen := names[name]; {
		case name == "":
			errs = append(errs, at("name")+" is empty; write the sub-project's name, such as "+exampleProjectPath)
		case seen:
			errs = append(errs, fmt.Sprintf("%s repeats the name of projects[%d]; give each sub-project a name of its own", at("name"), prev))
		default:
			names[name] = i
		}
		example := exampleProjectPath
		if plainValue.MatchString(name) && repoPathProblem(name, "", "", "") == "" {
			example = name
		}
		switch r := repoPathProblem(dir, "a folder", example, path.Clean(dir)); {
		case dir == "":
			errs = append(errs, at("path")+" is empty; write the sub-project's folder, relative to the repository root, such as "+example)
		case r != "":
			errs = append(errs, fmt.Sprintf("%s %q %s", at("path"), p.Path, r))
		case dir == ".":
			errs = append(errs, fmt.Sprintf("%s %q is the repository root; write the sub-project's own folder under it, such as %s", at("path"), p.Path, example))
		default:
			if prev, seen := paths[dir]; seen {
				errs = append(errs, fmt.Sprintf("%s %q is also the path of projects[%d]; give each sub-project a folder of its own", at("path"), p.Path, prev))
			} else {
				paths[dir] = i
				if fsys != nil && !isDir(fsys, dir) {
					errs = append(errs, fmt.Sprintf("%s %q names no folder in the repository; create the folder, or write the path of the sub-project's folder, relative to the repository root", at("path"), p.Path))
				}
			}
		}
		if strings.TrimSpace(p.Kind) == "" {
			errs = append(errs, at("kind")+" is empty; write what the sub-project is, such as "+projectKinds)
		}
	}
	return errs
}

// isDir reports whether name is a folder in fsys; anything it cannot stat is
// not.
func isDir(fsys fs.FS, name string) bool {
	st, err := fs.Stat(fsys, name)
	return err == nil && st.IsDir()
}

// ProjectAction is what an edit of projects does to an entry.
type ProjectAction string

// The actions of an edit of projects.
const (
	// ProjectAdd adds the entry given after the last, refusing a name the
	// list has already.
	ProjectAdd ProjectAction = "add"
	// ProjectSet replaces the entry named with the one given: its name, path,
	// kind, and tags, and its components and install when the entry given
	// sets them, keeping them when it leaves them out.
	ProjectSet ProjectAction = "set"
	// ProjectRemove removes the entry named.
	ProjectRemove ProjectAction = "remove"
)

// ProjectOp is an edit of projects: its action, and the entry it is on.
type ProjectOp struct {
	Action ProjectAction
	// Name is the name of the entry to set or remove as the manifest has it;
	// empty means the name of the entry given. A set whose entry has another
	// name renames it. Add does not use it.
	Name string
}

// SetProjects writes the edit of projects to the manifest at file, as
// EditProjects does, checking each path against fsys, the repository root,
// unless it is nil; on a refusal it writes nothing.
func SetProjects(file string, fsys fs.FS, op ProjectOp, p Project) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out, err := EditProjects(data, fsys, op, p)
	if err != nil {
		return err
	}
	return os.WriteFile(file, out, 0o644)
}

// EditProjects is the manifest data with one entry of projects added,
// replaced by name, or removed by name, as op says; every other line, the
// order of the keys, and the comments are kept. An add to projects: [] turns
// it into a block holding the entry, and removing the last entry leaves
// projects: []. It refuses an add of a name the list has, a set or remove of
// one it has not, and an entry laid out in a way it cannot edit or that does
// not read back as given; and, with a *RefusedError naming each field and
// why, any problem ProjectErrors, ComponentErrors, or InstallErrors finds in
// the manifest as it would be, checking paths against fsys, the repository
// root, unless fsys is nil.
func EditProjects(data []byte, fsys fs.FS, op ProjectOp, p Project) ([]byte, error) {
	p.Name, p.Path, p.Kind = strings.TrimSpace(p.Name), strings.TrimSpace(p.Path), strings.TrimSpace(p.Kind)
	name := strings.TrimSpace(op.Name)
	if name == "" {
		name = p.Name
	}
	if name == "" {
		return nil, fmt.Errorf("give the name of the sub-project to %s", op.Action)
	}
	var before Manifest
	if err := yaml.Unmarshal(data, &before); err != nil {
		return nil, fmt.Errorf("cannot read system-flow.yaml to %s sub-project %q: %w; fix its YAML by hand, then run this again", op.Action, name, err)
	}
	idx := slices.IndexFunc(before.Projects, func(q Project) bool { return strings.TrimSpace(q.Name) == name })
	lines := strings.Split(string(data), "\n")
	var want []Project
	var err error
	switch op.Action {
	case ProjectAdd:
		if idx >= 0 {
			return nil, fmt.Errorf("system-flow.yaml has a sub-project %q already; give the new one another name, or set that one", name)
		}
		want = append(slices.Clone(before.Projects), p)
		lines, err = addProjectLines(lines, p)
	case ProjectRemove:
		if idx < 0 {
			return nil, noProject(before, name)
		}
		want = slices.Delete(slices.Clone(before.Projects), idx, idx+1)
		lines, err = removeProjectLines(lines, idx, name)
	case ProjectSet:
		if idx < 0 {
			return nil, noProject(before, name)
		}
		if p.Components == nil {
			p.Components = before.Projects[idx].Components
		}
		if p.Install == nil {
			p.Install = before.Projects[idx].Install
		}
		want = slices.Clone(before.Projects)
		want[idx] = p
		lines, err = setProjectLines(lines, idx, name, p)
	default:
		return nil, fmt.Errorf("cannot edit projects with %q; add, set, or remove a sub-project", op.Action)
	}
	if err != nil {
		return nil, err
	}
	out := []byte(strings.Join(lines, "\n"))
	var after Manifest
	if err := yaml.Unmarshal(out, &after); err != nil || !slices.EqualFunc(after.Projects, want, sameProject) {
		return nil, fmt.Errorf("cannot %s sub-project %q so that system-flow.yaml reads it back as given; the file is left as it was: edit projects by hand", op.Action, name)
	}
	msgs := after.ProjectErrors(fsys)
	msgs = append(msgs, after.ComponentErrors(fsys)...)
	msgs = append(msgs, after.InstallErrors()...)
	if len(msgs) > 0 {
		problems := make([]Problem, len(msgs))
		for i, msg := range msgs {
			problems[i] = splitProblem(msg)
		}
		return nil, &RefusedError{Problems: problems}
	}
	return out, nil
}

// noProject is the refusal of a name the manifest's projects do not have,
// naming those it has.
func noProject(m Manifest, name string) error {
	names := make([]string, len(m.Projects))
	for i, p := range m.Projects {
		names[i] = p.Name
	}
	if len(names) == 0 {
		return fmt.Errorf("system-flow.yaml has no sub-project %q, and names none under projects; add the sub-project there first", name)
	}
	return fmt.Errorf("system-flow.yaml has no sub-project %q; name one of %s", name, strings.Join(names, ", "))
}

// sameProject reports whether two entries of projects say the same, a list
// absent and one empty alike.
func sameProject(a, b Project) bool {
	return a.Name == b.Name && a.Path == b.Path && a.Kind == b.Kind && slices.Equal(a.Tags, b.Tags) &&
		sameComponents(a.Components, b.Components) && sameInstall(a.Install, b.Install)
}

// sameInstall reports whether two install keys say the same, a list absent
// and one empty alike.
func sameInstall(a, b *Install) bool {
	if a == nil || b == nil {
		return a == b
	}
	return slices.Equal(a.Command, b.Command) && slices.Equal(a.Paths, b.Paths)
}

// projectsList is where projects is among a manifest's lines.
type projectsList struct {
	// key is the projects: line and the line after its value's last, its at
	// -1 when there is no projects key.
	key span
	// value and comment are what follows projects: on its line.
	value, comment string
	// items are the first line of each of its items, at the indent of the
	// first, and itemIn is that indent.
	items  []int
	itemIn string
}

// locateProjects finds the top-level projects key among a manifest's lines
// and the items of its list.
func locateProjects(lines []string) projectsList {
	key, _ := findKey(lines, 0, len(lines), "projects")
	l := projectsList{key: key}
	if key.at < 0 {
		return l
	}
	rest, _ := keyRest(lines[key.at], "projects")
	l.value, l.comment = splitComment(rest)
	for i := key.at + 1; i < key.end; i++ {
		if !isContent(lines[i]) {
			continue
		}
		t := strings.TrimSpace(lines[i])
		if len(l.items) == 0 {
			l.itemIn = indent(lines[i])
		}
		if indent(lines[i]) == l.itemIn && (t == "-" || strings.HasPrefix(t, "- ")) {
			l.items = append(l.items, i)
		}
	}
	return l
}

// projectEntry is where an entry of projects is among a manifest's lines.
type projectEntry struct {
	// from is its first line, the dash's, and to the next item's or the
	// line after the list's last.
	from, to int
	// keyIn is the indent of its keys, and head what follows the dash on its
	// first line, but a comment.
	keyIn, head string
}

// entry is where the idx-th entry of the list is, or why it cannot be edited:
// a list written on one line, an index past its items, an entry that is
// empty, and one written on one line.
func (l projectsList) entry(lines []string, idx int) (projectEntry, string) {
	switch {
	case l.key.at < 0:
		return projectEntry{}, "system-flow.yaml has no projects block"
	case l.value != "":
		return projectEntry{}, "projects is written on one line; write it as a block, a sub-project an item under it"
	case idx >= len(l.items):
		return projectEntry{}, "projects is not a list of items, each a sub-project"
	}
	e := projectEntry{from: l.items[idx], to: l.key.end}
	if idx+1 < len(l.items) {
		e.to = l.items[idx+1]
	}
	first := strings.TrimPrefix(strings.TrimSpace(lines[e.from]), "-")
	e.head, _ = splitComment(first)
	switch {
	case e.head == "":
		next := e.from + 1
		for next < e.to && !isContent(lines[next]) {
			next++
		}
		if next == e.to {
			return e, "its entry under projects is empty"
		}
		e.keyIn = indent(lines[next])
	case strings.HasPrefix(e.head, "{"):
		return e, "its entry under projects is written on one line; write it as a block, a key a line"
	default:
		e.keyIn = l.itemIn + " " + strings.Repeat(" ", len(first)-len(strings.TrimLeft(first, " \t")))
	}
	return e, ""
}

// find is where key is among the entry's keys, its at -1 when it is not
// there, and whether it is on the entry's first line, after the dash.
func (e projectEntry) find(lines []string, key string) (span, bool) {
	if _, ok := keyRest(e.head, key); ok {
		return span{at: e.from, end: e.from + 1}, true
	}
	for i := e.from + 1; i < e.to; i++ {
		if isContent(lines[i]) && indent(lines[i]) == e.keyIn {
			if _, ok := keyRest(lines[i], key); ok {
				return span{at: i, end: valueEnd(lines, i, e.to)}, false
			}
		}
	}
	return span{at: -1}, false
}

// projectKeyLines are the lines of an entry's keys that flai subproject
// writes, in the order it writes them, at the indent in: name, path, kind,
// tags when there are any, and components and install when set.
func projectKeyLines(p Project, in string) []string {
	out := []string{
		in + "name: " + yamlScalar(p.Name),
		in + "path: " + yamlScalar(p.Path),
		in + "kind: " + yamlScalar(p.Kind),
	}
	if len(p.Tags) > 0 {
		out = append(out, in+"tags: "+flowList(p.Tags))
	}
	if p.Components != nil {
		out = append(out, componentsBlock(in, "", p.Components)...)
	}
	if p.Install != nil {
		out = append(out, installBlock(in, "", p.Install)...)
	}
	return out
}

// installBlock is the lines of an install key at the indent in, keeping the
// comment of the key's line old, which it replaces, when there is one.
func installBlock(in, old string, i *Install) []string {
	return []string{
		blockHead(in, "install", old),
		in + "  command: " + flowList(i.Command),
		in + "  paths: " + flowList(i.Paths),
	}
}

// blockHead is the key line of a block at the indent in: old, the line it
// replaces, when it holds the key and nothing but a comment after it, so its
// comment is kept, and key: with old's comment otherwise.
func blockHead(in, key, old string) string {
	head := in + key + ":"
	if rest, ok := keyRest(old, key); ok {
		value, comment := splitComment(rest)
		if value == "" {
			return old
		}
		return withComment(head, comment)
	}
	return head
}

// addProjectLines are the manifest's lines with p added as the last entry of
// projects, its keys laid out as the first entry's are: after the last item,
// in a projects: [] turned into a block, or in a projects block added at the
// end of the file when there is none.
func addProjectLines(lines []string, p Project) ([]string, error) {
	l := locateProjects(lines)
	if l.key.at < 0 {
		block := append([]string{"projects:"}, entryLines(p, "  ", "    ")...)
		at := len(lines)
		if at > 0 && lines[at-1] == "" {
			at--
		} else {
			block = append(block, "")
		}
		return slices.Concat(lines[:at], block, lines[at:]), nil
	}
	switch l.value {
	case "", "[]", "null", "~":
	default:
		return nil, fmt.Errorf("cannot add sub-project %q: projects is written on one line (projects: %s); write it as a block, a sub-project an item under it, then run this again", p.Name, l.value)
	}
	itemIn, keyIn := "  ", "    "
	if len(l.items) > 0 {
		itemIn, keyIn = l.itemIn, l.itemIn+"  "
		if e, why := l.entry(lines, 0); why == "" && len(e.keyIn) >= len(itemIn)+2 {
			keyIn = e.keyIn
		}
	}
	out := slices.Clone(lines)
	if l.value != "" {
		out[l.key.at] = withComment("projects:", l.comment)
	}
	return slices.Concat(out[:l.key.end], entryLines(p, itemIn, keyIn), out[l.key.end:]), nil
}

// entryLines are p's lines as an item of projects, its dash at itemIn and its
// keys at keyIn, the first on the dash's line.
func entryLines(p Project, itemIn, keyIn string) []string {
	out := projectKeyLines(p, keyIn)
	out[0] = itemIn + "-" + strings.Repeat(" ", len(keyIn)-len(itemIn)-1) + strings.TrimPrefix(out[0], keyIn)
	return out
}

// removeProjectLines are the manifest's lines without the idx-th entry of
// projects, named name: its lines up to its last content line, so that
// comments before the next entry stay with it. Removing the last entry
// leaves projects: [].
func removeProjectLines(lines []string, idx int, name string) ([]string, error) {
	l := locateProjects(lines)
	e, why := l.entry(lines, idx)
	if why != "" {
		return nil, fmt.Errorf("cannot remove sub-project %q: %s; edit projects by hand", name, why)
	}
	out := slices.Delete(slices.Clone(lines), e.from, contentEnd(lines, e.from, e.to))
	if len(l.items) == 1 {
		out[l.key.at] = withComment("projects: []", l.comment)
	}
	return out, nil
}

// setProjectLines are the manifest's lines with the idx-th entry of projects,
// named name, given p's keys: name, path, and kind replaced where they are,
// tags replaced or removed when p has none, components and install replaced
// when p sets them, each kept on its line with its comment, and a key the
// entry lacks added after the key before it; every other key of the entry is
// kept.
func setProjectLines(lines []string, idx int, name string, p Project) ([]string, error) {
	type write struct {
		key    string
		render func(in, old string) []string // nil removes the key
	}
	scalar := func(key, value string) write {
		return write{key, func(in, old string) []string { return []string{keyLine(in, key, value, old)} }}
	}
	writes := []write{scalar("name", yamlScalar(p.Name)), scalar("path", yamlScalar(p.Path)), scalar("kind", yamlScalar(p.Kind)), {key: "tags"}}
	if len(p.Tags) > 0 {
		writes[3] = scalar("tags", flowList(p.Tags))
	}
	if c := p.Components; c != nil {
		writes = append(writes, write{"components", func(in, old string) []string { return componentsBlock(in, old, c) }})
	}
	if i := p.Install; i != nil {
		writes = append(writes, write{"install", func(in, old string) []string { return installBlock(in, old, i) }})
	}
	prev := ""
	for _, w := range writes {
		var why string
		lines, why = entryKeyLines(lines, idx, w.key, prev, w.render)
		if why != "" {
			return nil, fmt.Errorf("cannot set sub-project %q: %s; edit the entry by hand, or write it that way and run this again", name, why)
		}
		prev = w.key
	}
	return lines, nil
}

// entryKeyLines are the manifest's lines with key of the idx-th entry of
// projects written by render, or removed when render is nil: replaced where
// it is, or added after the key named after, or after the entry's last line
// when it has no such key. A key on the entry's first line, after the dash,
// may be given another one-line value only.
func entryKeyLines(lines []string, idx int, key, after string, render func(in, old string) []string) ([]string, string) {
	e, why := locateProjects(lines).entry(lines, idx)
	if why != "" {
		return nil, why
	}
	at, onHead := e.find(lines, key)
	if at.at < 0 {
		if render == nil {
			return lines, ""
		}
		end := contentEnd(lines, e.from, e.to)
		if prev, _ := e.find(lines, after); after != "" && prev.at >= 0 {
			end = prev.end
		}
		return slices.Concat(lines[:end], render(e.keyIn, ""), lines[end:]), ""
	}
	if onHead {
		rest, _ := keyRest(e.head, key)
		value, _ := splitComment(rest)
		var block []string
		if render != nil {
			line := lines[at.at]
			pos := strings.Index(line, key+":")
			block = render(line[:pos], strings.Repeat(" ", pos)+line[pos:])
		}
		if value == "" || len(block) != 1 {
			return nil, fmt.Sprintf("its entry under projects begins with %s; move name to its first line", key)
		}
		return slices.Concat(lines[:at.at], block, lines[at.end:]), ""
	}
	if render == nil {
		return slices.Delete(slices.Clone(lines), at.at, at.end), ""
	}
	return slices.Concat(lines[:at.at], render(e.keyIn, lines[at.at]), lines[at.end:]), ""
}

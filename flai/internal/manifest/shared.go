package manifest

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

// Claims is the project's say about how stories' touches claim paths
// (S-0295, ADR-0096).
type Claims struct {
	// Shared are glob patterns of the paths that many stories change in
	// separate sections or new files: an overlap wholly inside one holds no
	// story. A pattern that is not valid frees nothing.
	Shared []string `yaml:"shared,omitempty" json:"shared,omitempty"`
}

// PatternError is a pattern of claims.shared that is not valid.
type PatternError struct {
	// Index is the pattern's place in claims.shared, from 0.
	Index int
	// Pattern is the pattern as written.
	Pattern string
	// Reason says what is wrong and what to write instead, beginning with a
	// verb: "is empty; write ...".
	Reason string
}

func (e PatternError) Error() string {
	return fmt.Sprintf("claims.shared pattern %q %s", e.Pattern, e.Reason)
}

// Errors are the patterns of claims.shared that are not valid, in list
// order; none when every one is.
func (c Claims) Errors() []PatternError {
	var errs []PatternError
	for i, p := range c.Shared {
		if r := patternProblem(p); r != "" {
			errs = append(errs, PatternError{Index: i, Pattern: p, Reason: r})
		}
	}
	return errs
}

// patternProblem says what is wrong with a pattern, or "" when it is valid:
// a path relative to the repository root, segments separated by /, each a
// path.Match pattern or ** alone.
func patternProblem(p string) string {
	switch {
	case strings.TrimSpace(p) == "":
		return "is empty; write a path relative to the repository root, such as design/adrs, or remove it"
	case strings.HasPrefix(p, "/") || strings.HasPrefix(p, `\`) || len(p) > 1 && p[1] == ':':
		return "is an absolute path; write it relative to the repository root, such as design/adrs"
	}
	for _, seg := range strings.Split(p, "/") {
		switch {
		case seg == "":
			return "has an empty segment, a doubled or trailing /; write it with single slashes and none at the end"
		case seg == "..":
			return "has a .. segment, which leaves the repository; write a path inside it, relative to its root"
		case seg == ".":
			return "has a . segment; write the path without it, such as design/adrs for ./design/adrs"
		case seg != "**" && strings.Contains(seg, "**"):
			return "has ** inside a segment; ** stands for whole segments only, so make it a segment of its own, such as design/**/x.md, or use * for characters within one"
		}
		if _, err := path.Match(seg, ""); err != nil {
			return "is not a valid glob, such as an unclosed [ or a trailing \\; fix it or remove it"
		}
	}
	return ""
}

// Match is the first valid pattern of claims.shared that matches p, a file's
// path relative to the repository root, and whether there is one. A pattern
// with no glob character (* ? [ \) matches the path itself and everything
// below it, so design/adrs and design/adrs/** are the same; * is any
// characters within one segment, ** zero or more whole segments, and ? one
// character within a segment.
func (c Claims) Match(p string) (string, bool) {
	return c.first(p, false)
}

// Covers is the first valid pattern of claims.shared that the touches entry
// lies wholly inside, and whether there is one. The entry is a path relative
// to the repository root, a component's name already read as its path. A
// touches entry may name a file or a folder, and flai cannot always tell
// which, so an entry is read as a folder unless its last segment has an
// extension (a dot after its first character and before its last, such as
// flai.md): a folder lies inside a pattern only when the pattern matches it
// and every path below it, which only a plain pattern or one ending in **
// does, and a file lies inside a pattern that matches it, as Match says. So
// docs/users is inside docs and docs/users/**, and not inside
// docs/users/*.md, because the folder may hold other files. The rule errs
// towards holding: a file with no extension, such as Makefile, is read as a
// folder, inside a plain pattern or a ** but not a *.
func (c Claims) Covers(entry string) (string, bool) {
	e := cleanEntry(entry)
	return c.first(e, !hasExtension(path.Base(e)))
}

// first is the first valid pattern that matches entry, or, when folder is
// set, matches entry and everything below it.
func (c Claims) first(entry string, folder bool) (string, bool) {
	e := cleanEntry(entry)
	if e == "" || patternProblem(e) != "" {
		return "", false
	}
	segs := strings.Split(e, "/")
	for _, p := range c.Shared {
		if patternProblem(p) != "" {
			continue
		}
		ps := strings.Split(p, "/")
		if !strings.ContainsAny(p, `*?[\`) {
			ps = append(ps, "**")
		}
		if matchSegments(ps, segs, folder) {
			return p, true
		}
	}
	return "", false
}

// cleanEntry is a path or touches entry as it is matched, without spaces
// around it, a leading ./, or a trailing /.
func cleanEntry(entry string) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(entry), "./"), "/")
}

// hasExtension reports whether a path's last segment looks like a file's
// name: a dot after its first character and before its last.
func hasExtension(base string) bool {
	i := strings.LastIndex(base, ".")
	return i > 0 && i < len(base)-1
}

// matchSegments reports whether the pattern's segments match the path's.
// With below set, the pattern must also match every path below it: what is
// left of the pattern once the path is used up must be one or more ** and
// nothing else, since anything else needs a segment more or matches no
// deeper path.
func matchSegments(p, s []string, below bool) bool {
	if len(s) == 0 {
		if below && len(p) == 0 {
			return false
		}
		for _, x := range p {
			if x != "**" {
				return false
			}
		}
		return true
	}
	if len(p) == 0 {
		return false
	}
	if p[0] == "**" {
		return matchSegments(p, s[1:], below) || matchSegments(p[1:], s, below)
	}
	ok, _ := path.Match(p[0], s[0])
	return ok && matchSegments(p[1:], s[1:], below)
}

// Change is what an edit of claims.shared did.
type Change struct {
	// Added are the patterns added, in the order given.
	Added []string `json:"added,omitempty"`
	// Removed are the patterns removed, in the order given.
	Removed []string `json:"removed,omitempty"`
	// Shared is the list after the edit.
	Shared []string `json:"shared"`
}

// AddShared adds patterns to the end of claims.shared in the manifest at
// file, as AddSharedBytes does, and writes it back.
func AddShared(file string, patterns ...string) (Change, error) {
	return editFile(file, func(data []byte) ([]byte, Change, error) { return AddSharedBytes(data, patterns...) })
}

// RemoveShared removes patterns from claims.shared in the manifest at file,
// as RemoveSharedBytes does, and writes it back.
func RemoveShared(file string, patterns ...string) (Change, error) {
	return editFile(file, func(data []byte) ([]byte, Change, error) { return RemoveSharedBytes(data, patterns...) })
}

func editFile(file string, edit func([]byte) ([]byte, Change, error)) (Change, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return Change{}, err
	}
	out, ch, err := edit(data)
	if err != nil {
		return Change{}, fmt.Errorf("%s: %w", file, err)
	}
	if err := os.WriteFile(file, out, 0o644); err != nil {
		return Change{}, err
	}
	return ch, nil
}

// AddSharedBytes is the manifest data with patterns added to the end of
// claims.shared, creating claims and shared when they are absent. Every
// other line is left as it was, and the list's own items and comments too.
// It refuses, changing nothing, a pattern that is not valid or is in the
// list already.
func AddSharedBytes(data []byte, patterns ...string) ([]byte, Change, error) {
	if len(patterns) == 0 {
		return nil, Change{}, fmt.Errorf("name at least one pattern to add to claims.shared, such as design/adrs")
	}
	old, err := sharedList(data)
	if err != nil {
		return nil, Change{}, err
	}
	next := slices.Clone(old)
	var added []string
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if r := patternProblem(p); r != "" {
			return nil, Change{}, fmt.Errorf("cannot add %q to claims.shared: it %s", p, r)
		}
		if slices.Contains(next, p) {
			return nil, Change{}, fmt.Errorf("cannot add %q to claims.shared: it is in the list already; nothing to add", p)
		}
		next = append(next, p)
		added = append(added, p)
	}
	return rewriteShared(data, old, Change{Added: added, Shared: next})
}

// RemoveSharedBytes is the manifest data with patterns removed from
// claims.shared, leaving shared: [] when none is left. Every other line is
// left as it was, and the list's other items and comments too. It refuses,
// changing nothing, a pattern that is not in the list; one that is not valid
// may be removed.
func RemoveSharedBytes(data []byte, patterns ...string) ([]byte, Change, error) {
	if len(patterns) == 0 {
		return nil, Change{}, fmt.Errorf("name at least one pattern to remove from claims.shared")
	}
	old, err := sharedList(data)
	if err != nil {
		return nil, Change{}, err
	}
	next := slices.Clone(old)
	var removed []string
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		i := slices.Index(next, p)
		if i < 0 {
			has := "is empty"
			if len(old) > 0 {
				has = "holds " + strings.Join(old, ", ")
			}
			return nil, Change{}, fmt.Errorf("cannot remove %q from claims.shared: it is not in the list, which %s; name a pattern as the list writes it", p, has)
		}
		next = slices.Delete(next, i, i+1)
		removed = append(removed, p)
	}
	return rewriteShared(data, old, Change{Removed: removed, Shared: next})
}

// SharedLines are the lines, from 1, of claims.shared's patterns in the
// manifest data, in list order, for a report that points at each. A list
// written on the key's line, as [a, b], gives that line for each; nil when
// there is no list.
func SharedLines(data []byte) []int {
	old, err := sharedList(data)
	if err != nil || len(old) == 0 {
		return nil
	}
	b := locate(strings.Split(string(data), "\n"))
	if b.key < 0 {
		return nil
	}
	out := make([]int, len(old))
	for i := range out {
		out[i] = b.key + 1
		if len(b.items) == len(old) {
			out[i] = b.items[i] + 1
		}
	}
	return out
}

// sharedList is claims.shared as the manifest data has it.
func sharedList(data []byte) ([]string, error) {
	var doc struct {
		Claims Claims `yaml:"claims"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("cannot read claims.shared: %w; fix the manifest's YAML, then run this again", err)
	}
	return doc.Claims.Shared, nil
}

// sharedBlock is where claims and claims.shared are in a manifest's lines,
// by index from 0.
type sharedBlock struct {
	claims        int    // the claims: line, or -1
	claimsValue   string // what follows claims: on its line, but a comment
	claimsComment string
	claimsEnd     int    // the line after the claims block's last key or value
	childIndent   string // the indent of claims's keys
	key           int    // the shared: line, or -1
	keyValue      string // what follows shared: on its line, but a comment
	keyComment    string
	items         []int // the list's item lines, in order
	end           int   // the line after the list's last line
}

// locate finds claims and claims.shared in the lines of a manifest, a
// top-level claims: block with shared: among its keys.
func locate(lines []string) sharedBlock {
	b := sharedBlock{claims: -1, key: -1, childIndent: "  "}
	for i, l := range lines {
		if rest, ok := keyRest(l, "claims"); ok && indent(l) == "" {
			b.claims = i
			b.claimsValue, b.claimsComment = splitComment(rest)
			break
		}
	}
	if b.claims < 0 {
		return b
	}
	b.claimsEnd = b.claims + 1
	child := ""
	for i := b.claims + 1; i < len(lines); i++ {
		l := lines[i]
		if !isContent(l) {
			continue
		}
		in := indent(l)
		if in == "" {
			break
		}
		b.claimsEnd = i + 1
		if child == "" {
			child = in
			b.childIndent = in
		}
		if rest, ok := keyRest(l, "shared"); ok && in == child && b.key < 0 {
			b.key = i
			b.keyValue, b.keyComment = splitComment(rest)
		}
	}
	if b.key < 0 {
		return b
	}
	keyIn := indent(lines[b.key])
	b.end = b.key + 1
	for i := b.key + 1; i < b.claimsEnd; i++ {
		l := lines[i]
		if !isContent(l) {
			continue
		}
		in, t := indent(l), strings.TrimSpace(l)
		item := t == "-" || strings.HasPrefix(t, "- ")
		if belongs := len(in) > len(keyIn) || in == keyIn && item; !belongs {
			break
		}
		b.end = i + 1
		if item && (len(b.items) == 0 || in == indent(lines[b.items[0]])) {
			b.items = append(b.items, i)
		}
	}
	return b
}

// rewriteShared is data with claims.shared rewritten from old to ch.Shared.
// When the list is written one item a line, items not removed keep their
// lines, comments between them stay, and added ones follow the last; a list
// written any other way is written again one item a line, and comments
// inside it are lost.
func rewriteShared(data []byte, old []string, ch Change) ([]byte, Change, error) {
	if ch.Shared == nil {
		ch.Shared = []string{}
	}
	lines := strings.Split(string(data), "\n")
	b := locate(lines)
	var out []string
	switch {
	case b.key >= 0:
		out = slices.Concat(lines[:b.key], sharedSpan(lines, b, old, ch), lines[b.end:])
	case b.claims >= 0:
		switch b.claimsValue {
		case "", "{}", "null", "~":
		default:
			return nil, Change{}, fmt.Errorf("cannot add to claims.shared: claims is written on one line (claims: %s); write it as a block with shared: under it, then run this again", b.claimsValue)
		}
		head := slices.Clone(lines[:b.claimsEnd])
		if b.claimsValue != "" {
			head[b.claims] = withComment("claims:", b.claimsComment)
		}
		block := append([]string{b.childIndent + "shared:"}, items(b.childIndent+"  - ", ch.Added)...)
		out = slices.Concat(head, block, lines[b.claimsEnd:])
	default:
		block := append([]string{"claims:", "  shared:"}, items("    - ", ch.Added)...)
		at := len(lines)
		if lines[at-1] == "" {
			at--
		} else {
			block = append(block, "")
		}
		out = slices.Concat(lines[:at], block, lines[at:])
	}
	result := []byte(strings.Join(out, "\n"))
	if got, err := sharedList(result); err != nil || !slices.Equal(got, ch.Shared) {
		return nil, Change{}, fmt.Errorf("cannot rewrite claims.shared so that it reads back as %v; edit it by hand in the manifest", ch.Shared)
	}
	return result, ch, nil
}

// sharedSpan is the shared: key's lines rewritten, from its key line to its
// list's last line.
func sharedSpan(lines []string, b sharedBlock, old []string, ch Change) []string {
	keyIn := indent(lines[b.key])
	prefix := keyIn + "  - "
	byLine := b.keyValue == "" && len(b.items) == len(old)
	for i, at := range b.items {
		if !byLine {
			break
		}
		v, ok := itemValue(lines[at])
		byLine = ok && v == old[i]
	}
	if byLine && len(b.items) > 0 {
		prefix = indent(lines[b.items[0]]) + "- "
	}
	key := withComment(keyIn+"shared:", b.keyComment)
	switch {
	case len(ch.Shared) == 0:
		key = withComment(keyIn+"shared: []", b.keyComment)
	case byLine:
		key = lines[b.key]
	}
	if !byLine {
		return append([]string{key}, items(prefix, ch.Shared)...)
	}
	span := []string{key}
	at := 1
	for i := b.key + 1; i < b.end; i++ {
		n := slices.Index(b.items, i)
		if n >= 0 && slices.Contains(ch.Removed, old[n]) {
			continue
		}
		span = append(span, lines[i])
		if n >= 0 {
			at = len(span)
		}
	}
	return slices.Insert(span, at, items(prefix, ch.Added)...)
}

// items are the patterns as list items, each after prefix.
func items(prefix string, patterns []string) []string {
	out := make([]string, len(patterns))
	for i, p := range patterns {
		out[i] = prefix + yamlScalar(p)
	}
	return out
}

// itemValue is the string a one-line list item holds, and whether it holds
// one.
func itemValue(line string) (string, bool) {
	v, _ := splitComment(strings.TrimPrefix(strings.TrimSpace(line), "-"))
	if v == "" {
		return "", false
	}
	var doc struct {
		V string `yaml:"v"`
	}
	if err := yaml.Unmarshal([]byte("v: "+v), &doc); err != nil {
		return "", false
	}
	return doc.V, true
}

// keyRest is what follows key: on a line whose key it is.
func keyRest(line, key string) (string, bool) {
	t := strings.TrimLeft(line, " \t")
	rest, ok := strings.CutPrefix(t, key+":")
	if !ok || rest != "" && rest[0] != ' ' && rest[0] != '\t' {
		return "", false
	}
	return rest, true
}

// splitComment splits a YAML value from a comment after it: a # outside
// quotes, at the start or after a space.
func splitComment(s string) (value, comment string) {
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case quote == '"' && c == '\\':
			i++
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '#' && (i == 0 || s[i-1] == ' ' || s[i-1] == '\t'):
			return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i:])
		}
	}
	return strings.TrimSpace(s), ""
}

func withComment(line, comment string) string {
	if comment == "" {
		return line
	}
	return line + " " + comment
}

func indent(line string) string {
	return line[:len(line)-len(strings.TrimLeft(line, " \t"))]
}

// isContent reports whether a line holds a key or a value, not only space
// or a comment.
func isContent(line string) bool {
	t := strings.TrimSpace(line)
	return t != "" && !strings.HasPrefix(t, "#")
}

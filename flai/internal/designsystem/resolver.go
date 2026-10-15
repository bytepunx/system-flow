package designsystem

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strconv"
	"strings"
)

// The themes a token file gives.
const (
	ThemeLight = "light"
	ThemeDark  = "dark"
)

// Themes are the theme names, in the order they are resolved.
var Themes = []string{ThemeLight, ThemeDark}

// ResolverVersion is the DTCG resolver format flai reads.
const ResolverVersion = "2025.10"

// Theme is one theme's tokens: the file's sources merged for it, aliases
// resolved.
type Theme struct {
	Name   string    `json:"name"`   // ThemeLight or ThemeDark
	Tokens *TokenSet `json:"tokens"` // an alias token's Value is its target's, its Ref still the path it names
}

// System is a token file resolved into its themes, with every problem.
type System struct {
	Resolver bool      `json:"resolver"` // a resolver document, not a plain token file
	Themes   []Theme   `json:"themes"`   // light then dark; nil when the file does not parse or has no resolutionOrder to apply
	Problems []Problem `json:"problems"` // by line, then path; never nil
}

// Load reads the token file name in fsys and resolves it. A file that is not
// there is an error that wraps fs.ErrNotExist; a file that does not parse or
// is wrong is no error but the System's problems.
func Load(fsys fs.FS, name string) (*System, error) {
	f, err := Read(fsys, name)
	var pe *ParseError
	switch {
	case errors.As(err, &pe):
		return &System{Problems: []Problem{pe.Problem()}}, nil
	case err != nil:
		return nil, err
	}
	return Resolve(f), nil
}

// Resolve resolves a parsed token file into its light and dark themes. A
// plain token file gives both themes the same tokens. A resolver document
// gives each theme its sets' sources merged in resolutionOrder, then its
// theme context's, a later token replacing the properties it gives of an
// earlier one; a token dark does not override keeps its light value. Each
// theme's aliases resolve within it. The problems are reading's, each
// theme's Validate, the aliases', and the resolver document's, a problem
// both themes report the same given once.
func Resolve(f *File) *System {
	root := f.Root
	if !f.Resolver && root.Get("sets") == nil && root.Get("modifiers") == nil {
		s := &TokenSet{Tokens: slices.Clone(f.Tokens.Tokens), Groups: f.Tokens.Groups, Problems: f.Tokens.Problems}
		probs := resolveAliases(s)
		probs = append(probs, Validate(s)...)
		return &System{Themes: []Theme{{ThemeLight, s}, {ThemeDark, s}}, Problems: tidy(probs)}
	}
	r := &resolver{sets: map[string][]*Node{}, contexts: map[string][]*Node{}, others: map[string]bool{}}
	r.document(root)
	r.readSets(root)
	r.readModifiers(root)
	sys := &System{Resolver: true}
	order, ok := r.order(root)
	if !ok {
		sys.Problems = tidy(r.probs)
		return sys
	}
	r.overrides(order)
	probs := r.probs
	for _, name := range Themes {
		s := ReadTokens(r.merged(order, name))
		probs = append(probs, resolveAliases(s)...)
		probs = append(probs, Validate(s)...)
		sys.Themes = append(sys.Themes, Theme{Name: name, Tokens: s})
	}
	sys.Problems = tidy(probs)
	return sys
}

// tidy drops a problem given twice and sorts the rest by line, then path.
func tidy(probs []Problem) []Problem {
	out := make([]Problem, 0, len(probs))
	seen := map[Problem]bool{}
	for _, p := range probs {
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	slices.SortStableFunc(out, func(a, b Problem) int {
		return cmp.Or(cmp.Compare(a.Line, b.Line), cmp.Compare(a.Path, b.Path))
	})
	return out
}

// resolveAliases points each alias token of s at its target's value and,
// when it has no type of its own, its target's type, following chains. An
// alias whose chain names no token or goes round a circle keeps its alias
// value and is a problem, as is one whose type is not its target's.
func resolveAliases(s *TokenSet) []Problem {
	index := make(map[string]int, len(s.Tokens))
	for i, t := range s.Tokens {
		index[t.Path] = i
	}
	var probs []Problem
	resolved := slices.Clone(s.Tokens)
	for i, t := range s.Tokens {
		if t.Ref == "" {
			continue
		}
		target, msg := follow(s.Tokens, index, t)
		switch {
		case msg != "": // the alias keeps its value
		case t.Type == "":
			resolved[i].Value, resolved[i].Type = target.Value, target.Type
		case target.Type != "" && target.Type != t.Type:
			msg = fmt.Sprintf("it is a %s token, and the token it names, %s, is a %s token; name a %s token or change its $type", t.Type, target.Path, target.Type, t.Type)
		default:
			resolved[i].Value = target.Value
		}
		if msg != "" {
			probs = append(probs, Problem{Path: t.Path, Line: t.Line, Message: msg})
		}
	}
	s.Tokens = resolved
	return probs
}

// follow answers the token at the end of t's chain of aliases, or why there
// is none.
func follow(tokens []Token, index map[string]int, t Token) (Token, string) {
	ref := "{" + t.Ref + "}"
	chain := []string{t.Path}
	for t.Ref != "" {
		if slices.Contains(chain, t.Ref) {
			return Token{}, fmt.Sprintf("the reference %s does not resolve: it goes round in a circle, %s; break the circle", ref, strings.Join(append(chain, t.Ref), " → "))
		}
		j, ok := index[t.Ref]
		if !ok {
			return Token{}, fmt.Sprintf("the reference %s does not resolve: no token has the path %s; name a token the file defines", ref, t.Ref)
		}
		chain = append(chain, t.Ref)
		t = tokens[j]
	}
	return t, ""
}

// resolver reads a DTCG resolver document in the subset flai reads
// (ADR-0146): sets with inline sources, one modifier, theme, with the
// contexts light and dark, and a resolutionOrder of references to them.
type resolver struct {
	probs    []Problem
	sets     map[string][]*Node // each set, by name, and its inline sources
	hasTheme bool               // modifiers.theme is given
	contexts map[string][]*Node // the theme modifier's contexts and their inline sources
	others   map[string]bool    // the modifiers other than theme, reported where they are given
}

// step is one entry of resolutionOrder flai applies: a set, or the theme
// modifier.
type step struct {
	set   string
	theme bool
}

func (r *resolver) problem(path string, line int, msg string) {
	r.probs = append(r.probs, Problem{Path: path, Line: line, Message: msg})
}

// The properties flai reads of a resolver document, a set, and a modifier.
var (
	documentKeys = []string{"$schema", "$extensions", "name", "description", "version", "sets", "modifiers", "resolutionOrder"}
	setKeys      = []string{"sources", "description", "$extensions"}
	modifierKeys = []string{"contexts", "default", "description", "$extensions"}
)

// known reports each member of n, at path, that is not one of keys; what
// names n in the message.
func (r *resolver) known(n *Node, path, what string, keys []string) {
	for _, m := range n.Members {
		if !slices.Contains(keys, m.Key) {
			r.problem(join(path, m.Key), m.Line, fmt.Sprintf("flai does not read the property %s of %s; remove it", m.Key, what))
		}
	}
}

// document checks the document's own properties and its version.
func (r *resolver) document(root *Node) {
	r.known(root, "", "a resolver document", documentKeys)
	v := root.Get("version")
	switch {
	case v == nil:
		r.problem("version", root.Line, fmt.Sprintf("a resolver document needs its version, %q; add it", ResolverVersion))
	case v.Kind != KindString || v.Text != ResolverVersion:
		r.problem("version", v.Line, fmt.Sprintf("version is %s, and flai reads the resolver format %s; write %q", describe(v), ResolverVersion, ResolverVersion))
	}
}

// readSets reads each set's inline sources.
func (r *resolver) readSets(root *Node) {
	sets := root.Get("sets")
	if sets == nil {
		return
	}
	if sets.Kind != KindObject {
		r.problem("sets", sets.Line, fmt.Sprintf("sets is %s; write an object of named sets, such as {\"tokens\": {\"sources\": [...]}}", sets.Kind))
		return
	}
	for _, m := range sets.Members {
		path := "sets." + m.Key
		r.sets[m.Key] = nil
		if m.Value.Kind != KindObject {
			r.problem(path, m.Line, fmt.Sprintf("a set is an object with its sources, not %s", m.Value.Kind))
			continue
		}
		r.known(m.Value, path, "a set", setKeys)
		src := m.Value.Get("sources")
		if src == nil {
			r.problem(path+".sources", m.Line, "a set needs its sources, an array of token trees such as [{\"color\": {...}}]; add it")
			continue
		}
		r.sets[m.Key] = r.sources(src, path+".sources")
	}
}

// readModifiers reads the theme modifier's contexts, reporting any other
// modifier.
func (r *resolver) readModifiers(root *Node) {
	const want = "a resolver document gives the dark values through one modifier, theme, with the contexts light and dark"
	mods := root.Get("modifiers")
	switch {
	case mods == nil:
		r.problem("modifiers.theme", root.Line, want+"; add it, or make this a plain token file")
		return
	case mods.Kind != KindObject:
		r.problem("modifiers", mods.Line, fmt.Sprintf("modifiers is %s; write an object holding the modifier theme", mods.Kind))
		return
	}
	for _, m := range mods.Members {
		if m.Key == "theme" {
			r.theme(m)
			continue
		}
		r.others[m.Key] = true
		r.problem("modifiers."+m.Key, m.Line, fmt.Sprintf("flai reads one modifier, theme, with the contexts light and dark; remove %s, or give its values in theme's contexts", m.Key))
	}
	if !r.hasTheme {
		r.problem("modifiers.theme", mods.Line, want+"; add it")
	}
}

// theme reads the theme modifier: its contexts light and dark, and its
// default.
func (r *resolver) theme(m Member) {
	const path = "modifiers.theme"
	r.hasTheme = true
	n := m.Value
	if n.Kind != KindObject {
		r.problem(path, m.Line, fmt.Sprintf("the modifier theme is %s; write an object with its contexts, light and dark", n.Kind))
		return
	}
	r.known(n, path, "a modifier", modifierKeys)
	ctx := n.Get("contexts")
	switch {
	case ctx == nil:
		r.problem(path+".contexts", m.Line, "the modifier theme needs its contexts, light and dark; add them")
	case ctx.Kind != KindObject:
		r.problem(path+".contexts", ctx.Line, fmt.Sprintf("contexts is %s; write an object with the contexts light and dark, each an array of sources", ctx.Kind))
	default:
		for _, c := range ctx.Members {
			cp := path + ".contexts." + c.Key
			if !slices.Contains(Themes, c.Key) {
				r.problem(cp, c.Line, fmt.Sprintf("flai reads the theme contexts light and dark only; remove %s", c.Key))
				continue
			}
			r.contexts[c.Key] = r.sources(c.Value, cp)
		}
		for _, name := range Themes {
			if ctx.Get(name) == nil {
				r.problem(path+".contexts."+name, ctx.Line, fmt.Sprintf("the modifier theme needs the context %s; add it, as [] where %s has the sets' values", name, name))
			}
		}
	}
	if d := n.Get("default"); d != nil && (d.Kind != KindString || !slices.Contains(Themes, d.Text)) {
		r.problem(path+".default", d.Line, fmt.Sprintf("default is %s; name a context, \"light\" or \"dark\"", describe(d)))
	}
}

// sources reads an array of sources at path, answering the inline ones.
func (r *resolver) sources(n *Node, path string) []*Node {
	if n.Kind != KindArray {
		r.problem(path, n.Line, fmt.Sprintf("this is %s; write an array of sources, each a token tree such as {\"color\": {...}}", n.Kind))
		return nil
	}
	var out []*Node
	for i, src := range n.Items {
		sp := path + "." + strconv.Itoa(i)
		ref := src.Get("$ref")
		switch {
		case src.Kind != KindObject:
			r.problem(sp, src.Line, fmt.Sprintf("a source is a token tree, an object such as {\"color\": {...}}, not %s", src.Kind))
		case ref != nil && ref.Kind == KindString && strings.HasPrefix(ref.Text, "#"):
			r.problem(sp, ref.Line, fmt.Sprintf("flai reads inline sources only; put the tokens here in place of the reference %s", ref.Text))
		case ref != nil:
			r.problem(sp, ref.Line, "flai reads inline sources only; move the tokens into this file, in place of this reference")
		default:
			out = append(out, src)
		}
	}
	return out
}

// order reads resolutionOrder into the steps flai applies; false when there
// is none to apply.
func (r *resolver) order(root *Node) ([]step, bool) {
	ro := root.Get("resolutionOrder")
	switch {
	case ro == nil:
		r.problem("resolutionOrder", root.Line, "a resolver document needs its resolutionOrder, its sets and then its modifier, such as [{\"$ref\": \"#/sets/tokens\"}, {\"$ref\": \"#/modifiers/theme\"}]; add it")
		return nil, false
	case ro.Kind != KindArray:
		r.problem("resolutionOrder", ro.Line, fmt.Sprintf("resolutionOrder is %s; write an array of references, such as [{\"$ref\": \"#/sets/tokens\"}, {\"$ref\": \"#/modifiers/theme\"}]", ro.Kind))
		return nil, false
	}
	var steps []step
	themed := false
	for i, item := range ro.Items {
		if st, ok := r.step(item, "resolutionOrder."+strconv.Itoa(i)); ok {
			steps = append(steps, st)
			themed = themed || st.theme
		}
	}
	if r.hasTheme && !themed {
		r.problem("resolutionOrder", ro.Line, "resolutionOrder does not apply the modifier theme, so dark has the sets' values; add {\"$ref\": \"#/modifiers/theme\"} after the sets")
	}
	return steps, true
}

// step reads one entry of resolutionOrder; false when flai cannot apply it.
func (r *resolver) step(item *Node, path string) (step, bool) {
	ref := item.Get("$ref")
	switch {
	case ref == nil:
		r.problem(path, item.Line, "flai reads resolutionOrder as references to this document's sets and modifier; give this a name under sets and write {\"$ref\": \"#/sets/<name>\"} here")
		return step{}, false
	case ref.Kind != KindString:
		r.problem(path, ref.Line, fmt.Sprintf("$ref is %s; write a reference such as \"#/sets/tokens\"", ref.Kind))
		return step{}, false
	case !strings.HasPrefix(ref.Text, "#"):
		r.problem(path, ref.Line, fmt.Sprintf("flai reads inline sources only; move the tokens of %s into this file's sets and reference their set as #/sets/<name>", ref.Text))
		return step{}, false
	}
	if name, ok := strings.CutPrefix(ref.Text, "#/sets/"); ok {
		name = unescape(name)
		if _, ok := r.sets[name]; ok {
			return step{set: name}, true
		}
		r.problem(path, ref.Line, fmt.Sprintf("the reference %s does not resolve: there is no set %s; name a set under sets", ref.Text, name))
		return step{}, false
	}
	if name, ok := strings.CutPrefix(ref.Text, "#/modifiers/"); ok {
		name = unescape(name)
		switch {
		case name == "theme" && r.hasTheme:
			return step{theme: true}, true
		case r.others[name]:
			return step{}, false
		}
		r.problem(path, ref.Line, fmt.Sprintf("the reference %s does not resolve: there is no modifier %s; reference the modifier theme, #/modifiers/theme", ref.Text, name))
		return step{}, false
	}
	r.problem(path, ref.Line, fmt.Sprintf("the reference %s does not resolve to a set or the modifier; write #/sets/<name> or #/modifiers/theme", ref.Text))
	return step{}, false
}

// unescape decodes one JSON Pointer reference token.
func unescape(s string) string { return strings.NewReplacer("~1", "/", "~0", "~").Replace(s) }

// overrides reports each token a theme context gives that no set in the
// resolution order defines.
func (r *resolver) overrides(order []step) {
	base := ReadTokens(r.merged(order, ""))
	for _, name := range Themes {
		for _, src := range r.contexts[name] {
			for _, t := range ReadTokens(src).Tokens {
				if _, ok := base.Get(t.Path); !ok {
					r.problem(t.Path, t.Line, fmt.Sprintf("the %s context overrides %s, which no set defines; define it in a set, or correct its path", name, t.Path))
				}
			}
		}
	}
}

// merged is the token tree for theme: each step's sources overlaid in order,
// the theme modifier's being theme's context. No theme, "", gives the sets'
// alone.
func (r *resolver) merged(order []step, theme string) *Node {
	tree := &Node{Kind: KindObject, Line: 1}
	for _, st := range order {
		srcs := r.sets[st.set]
		if st.theme {
			srcs = r.contexts[theme]
		}
		for _, src := range srcs {
			tree = overlay(tree, src)
		}
	}
	return tree
}

// overlay answers the group with over laid on it, as new nodes, leaving both
// as they are. A token in over gives the same token the properties over
// gives, its $value among them, and its line; a group in over merges into
// the same group, which keeps its line; anything else in over replaces what
// the group holds under its key.
func overlay(group, over *Node) *Node {
	out := &Node{Kind: KindObject, Line: group.Line, Members: slices.Clone(group.Members)}
	for _, m := range over.Members {
		i := slices.IndexFunc(out.Members, func(b Member) bool { return b.Key == m.Key })
		if i < 0 {
			out.Members = append(out.Members, m)
			continue
		}
		b := out.Members[i]
		switch {
		case isToken(b.Value) && isToken(m.Value):
			out.Members[i] = Member{Key: m.Key, Line: m.Line, Value: withProperties(b.Value, m.Value)}
		case isGroup(b.Value) && isGroup(m.Value):
			out.Members[i] = Member{Key: m.Key, Line: b.Line, Value: overlay(b.Value, m.Value)}
		default:
			out.Members[i] = m
		}
	}
	return out
}

// withProperties answers token with each property over gives in place of
// its own.
func withProperties(token, over *Node) *Node {
	out := &Node{Kind: KindObject, Line: over.Line, Members: slices.Clone(token.Members)}
	for _, m := range over.Members {
		if i := slices.IndexFunc(out.Members, func(b Member) bool { return b.Key == m.Key }); i >= 0 {
			out.Members[i] = m
		} else {
			out.Members = append(out.Members, m)
		}
	}
	return out
}

func isToken(n *Node) bool { return n.Kind == KindObject && n.Get("$value") != nil }

func isGroup(n *Node) bool { return n.Kind == KindObject && n.Get("$value") == nil }

// describe names a value as a message quotes it: a string in quotes, else its
// kind.
func describe(n *Node) string {
	if n.Kind == KindString {
		return strconv.Quote(n.Text)
	}
	return n.Kind.String()
}

func join(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

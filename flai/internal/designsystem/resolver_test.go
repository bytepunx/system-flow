package designsystem

import (
	"cmp"
	"errors"
	"io/fs"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

// resolverDoc wraps the full file as a resolver document's one set, so that
// its line n is line n+3, and gives dark the sources in dark, starting at
// line 65.
func resolverDoc(t *testing.T, dark string) []byte {
	t.Helper()
	return []byte(`{
  "version": "2025.10",
  "sets": { "tokens": { "sources": [
` + string(fullFile(t)) + `  ] } },
  "modifiers": { "theme": { "contexts": { "light": [], "dark": [
` + dark + `
  ] }, "default": "light" } },
  "resolutionOrder": [ { "$ref": "#/sets/tokens" }, { "$ref": "#/modifiers/theme" } ]
}`)
}

func load(t *testing.T, src []byte) *System {
	t.Helper()
	sys, err := Load(fstest.MapFS{"tokens.json": {Data: src}}, "tokens.json")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !slices.IsSortedFunc(sys.Problems, func(a, b Problem) int {
		return cmp.Or(cmp.Compare(a.Line, b.Line), cmp.Compare(a.Path, b.Path))
	}) {
		t.Errorf("problems are not sorted by line, then path: %+v", sys.Problems)
	}
	return sys
}

// tokens is the theme name's tokens.
func tokens(t *testing.T, sys *System, name string) *TokenSet {
	t.Helper()
	for _, th := range sys.Themes {
		if th.Name == name {
			return th.Tokens
		}
	}
	t.Fatalf("no theme %s in %+v", name, sys.Themes)
	return nil
}

func hex(t *testing.T, s *TokenSet, path string) string {
	t.Helper()
	tok, ok := s.Get(path)
	if !ok {
		t.Fatalf("no %s", path)
	}
	c, err := DecodeColor(tok.Value)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return c.Hex
}

func TestPlainFileGivesBothThemesItsValues(t *testing.T) {
	sys := load(t, fullFile(t))
	if sys.Resolver || len(sys.Problems) != 0 {
		t.Fatalf("Load(plain) = resolver %v, problems %+v; want a plain file with none", sys.Resolver, sys.Problems)
	}
	if len(sys.Themes) != 2 || sys.Themes[0].Name != ThemeLight || sys.Themes[1].Name != ThemeDark {
		t.Fatalf("Themes = %+v, want light then dark", sys.Themes)
	}
	light, dark := sys.Themes[0].Tokens, sys.Themes[1].Tokens
	if light != dark {
		t.Error("a plain file's light and dark tokens differ")
	}
	if got := hex(t, dark, "color.background"); got != "#ffffff" {
		t.Errorf("dark color.background = %s, want the one value #ffffff", got)
	}
}

func TestResolverGivesDarkValues(t *testing.T) {
	dark := `      { "color": {
        "background": { "$value": { "colorSpace": "srgb", "components": [0.07, 0.07, 0.08], "hex": "#121214" } },
        "primary": { "$value": { "colorSpace": "srgb", "components": [0.4, 0.6, 1], "hex": "#6699ff" } }
      } }`
	sys := load(t, resolverDoc(t, dark))
	if !sys.Resolver || len(sys.Problems) != 0 {
		t.Fatalf("Load(resolver) = resolver %v, problems %+v; want a resolver document with none", sys.Resolver, sys.Problems)
	}
	light, darkSet := tokens(t, sys, ThemeLight), tokens(t, sys, ThemeDark)
	bg, _ := darkSet.Get("color.background")
	if bg.Type != TypeColor || bg.Line != 66 || hex(t, darkSet, "color.background") != "#121214" {
		t.Errorf("dark color.background = %+v, want the override at line 66, a color by its group", bg)
	}
	if lbg, _ := light.Get("color.background"); lbg.Line != 8 || hex(t, light, "color.background") != "#ffffff" {
		t.Errorf("light color.background = %+v, want the set's at line 8", lbg)
	}
	if s, _ := darkSet.Get("color.surface"); s.Line != 9 || hex(t, darkSet, "color.surface") != "#f5f5f7" {
		t.Errorf("dark color.surface = %+v, want light's, at line 9", s)
	}
	if got := hex(t, darkSet, "color.link"); got != "#6699ff" {
		t.Errorf("dark color.link = %s, want dark's primary, #6699ff", got)
	}
	if link, _ := light.Get("color.link"); link.Ref != "color.primary" || link.Value.Kind != KindObject {
		t.Errorf("light color.link = %+v, want light's primary, an object, naming color.primary", link)
	}
	if len(darkSet.Tokens) != len(light.Tokens) {
		t.Errorf("dark has %d tokens, light %d", len(darkSet.Tokens), len(light.Tokens))
	}
}

// resolverBase is a resolver document with one colour; each case of
// TestResolverProblems replaces one of its lines.
var resolverBase = []string{
	`{`,
	`  "version": "2025.10",`,
	`  "sets": { "tokens": { "sources": [ { "color": { "$type": "color", "background": { "$value": { "colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff" } } } } ] } },`,
	`  "modifiers": { "theme": { "contexts": { "light": [], "dark": [] }, "default": "light" } },`,
	`  "resolutionOrder": [ { "$ref": "#/sets/tokens" }, { "$ref": "#/modifiers/theme" } ]`,
	`}`,
}

func TestResolverBaseHasOnlyMissingTokens(t *testing.T) {
	sys := load(t, []byte(strings.Join(resolverBase, "\n")))
	for _, p := range sys.Problems {
		if !strings.Contains(p.Message, "is required") {
			t.Errorf("problem %+v, want only required tokens missing", p)
		}
	}
	if got := at(sys.Problems, "color.text"); len(got) != 1 || got[0].Line != 3 {
		t.Errorf("problems at color.text = %+v, want one, at line 3, though both themes lack it", got)
	}
	if len(sys.Themes) != 2 {
		t.Errorf("Themes = %+v, want light and dark", sys.Themes)
	}
}

func TestResolverProblems(t *testing.T) {
	cases := []struct {
		name     string
		line     int // of resolverBase, 1-based, to replace
		with     string
		path     string
		at       int
		want     string
		noThemes bool
	}{
		{"a modifier other than theme", 4, `  "modifiers": {
    "theme": { "contexts": { "light": [], "dark": [] } },
    "density": { "contexts": { "compact": [] } }
  },`, "modifiers.density", 6, "one modifier, theme", false},
		{"a context other than light and dark", 4, `  "modifiers": { "theme": { "contexts": {
    "light": [],
    "dark": [],
    "sepia": []
  } } },`, "modifiers.theme.contexts.sepia", 7, "light and dark only; remove sepia", false},
		{"a context missing", 4, `  "modifiers": { "theme": {
    "contexts": { "light": [] } } },`, "modifiers.theme.contexts.dark", 5, "needs the context dark", false},
		{"a bad default", 4, `  "modifiers": { "theme": { "contexts": { "light": [], "dark": [] },
    "default": "sepia" } },`, "modifiers.theme.default", 5, `default is "sepia"`, false},
		{"no theme modifier", 4, `  "name": "Example",`, "modifiers.theme", 1, "one modifier, theme", false},
		{"theme not in the order", 5, `  "resolutionOrder": [ { "$ref": "#/sets/tokens" } ]`, "resolutionOrder", 5, "does not apply the modifier theme", false},
		{"no version", 2, `  "name": "Example",`, "version", 1, "needs its version", false},
		{"a wrong version", 2, `  "version": "2025.09",`, "version", 2, `version is "2025.09"`, false},
		{"no resolutionOrder", 5, `  "name": "Example"`, "resolutionOrder", 1, "needs its resolutionOrder", true},
		{"resolutionOrder not an array", 5, `  "resolutionOrder": {}`, "resolutionOrder", 5, "is an object", true},
		{"a set that does not resolve", 5, `  "resolutionOrder": [
    { "$ref": "#/sets/tokens" },
    { "$ref": "#/sets/x" },
    { "$ref": "#/modifiers/theme" }
  ]`, "resolutionOrder.1", 7, "#/sets/x does not resolve: there is no set x", false},
		{"a modifier that does not resolve", 5, `  "resolutionOrder": [ { "$ref": "#/sets/tokens" }, { "$ref": "#/modifiers/theme" },
    { "$ref": "#/modifiers/mode" } ]`, "resolutionOrder.2", 6, "there is no modifier mode", false},
		{"an inline entry in the order", 5, `  "resolutionOrder": [ { "type": "set", "name": "more", "sources": [] }, { "$ref": "#/sets/tokens" }, { "$ref": "#/modifiers/theme" } ]`, "resolutionOrder.0", 5, "references to this document's sets", false},
		{"a file in the order", 5, `  "resolutionOrder": [ { "$ref": "#/sets/tokens" },
    { "$ref": "dark.tokens.json" }, { "$ref": "#/modifiers/theme" } ]`, "resolutionOrder.1", 6, "move the tokens of dark.tokens.json into this file", false},
		{"a file source", 3, `  "sets": { "tokens": { "sources": [
    { "$ref": "base.tokens.json" }
  ] } },`, "sets.tokens.sources.0", 4, "flai reads inline sources only; move the tokens into this file", false},
		{"a set reference as a source", 4, `  "modifiers": { "theme": { "contexts": { "light": [], "dark": [
    { "$ref": "#/sets/tokens" } ] } } },`, "modifiers.theme.contexts.dark.0", 5, "in place of the reference #/sets/tokens", false},
		{"a set without sources", 3, `  "sets": { "tokens": { "sources": [] },
    "more": { "description": "none" } },`, "sets.more.sources", 4, "needs its sources", false},
		{"an override of a token no set has", 4, `  "modifiers": { "theme": { "contexts": { "light": [], "dark": [
    { "color": { "backgrund": { "$value": { "colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000" } } } }
  ] } } },`, "color.backgrund", 5, "the dark context overrides color.backgrund, which no set defines", false},
		{"a property flai does not read", 2, `  "version": "2025.10", "$defs": {},`, "$defs", 2, "does not read the property $defs of a resolver document", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lines := slices.Clone(resolverBase)
			lines[c.line-1] = c.with
			sys := load(t, []byte(strings.Join(lines, "\n")))
			if !sys.Resolver {
				t.Error("Load did not read a resolver document")
			}
			got := at(sys.Problems, c.path)
			if !slices.ContainsFunc(got, func(p Problem) bool { return p.Line == c.at && strings.Contains(p.Message, c.want) }) {
				t.Errorf("problems at %s = %+v, want one at line %d containing %q", c.path, got, c.at, c.want)
			}
			if (sys.Themes == nil) != c.noThemes {
				t.Errorf("Themes = %+v, want none: %v", sys.Themes, c.noThemes)
			}
		})
	}
}

func TestAliasesResolve(t *testing.T) {
	src := `{
  "color": {
    "$type": "color",
    "primary": { "$value": { "colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000" } },
    "link": { "$value": "{color.accent}" },
    "accent": { "$value": "{color.primary}" },
    "a": { "$value": "{color.b}" },
    "b": { "$value": "{color.a}" },
    "lost": { "$value": "{color.nowhere}" }
  },
  "radius": { "$type": "dimension", "pill": { "$value": "{color.primary}" } },
  "font": { "sans": { "$type": "fontFamily", "$value": "Inter" }, "mono": { "$value": "{font.sans}" } }
}`
	f := parsePlain(t, []byte(src))
	sys := Resolve(f)
	s := tokens(t, sys, ThemeLight)
	primary, _ := s.Get("color.primary")
	if link, _ := s.Get("color.link"); link.Value != primary.Value || link.Ref != "color.accent" || link.Type != TypeColor {
		t.Errorf("color.link = %+v, want primary's value through color.accent", link)
	}
	if mono, _ := s.Get("font.mono"); mono.Type != TypeFontFamily || mono.Value.Text != "Inter" {
		t.Errorf("font.mono = %+v, want font.sans's type and value", mono)
	}
	if orig, _ := f.Tokens.Get("color.link"); orig.Value.Kind != KindString {
		t.Error("Resolve changed the file's own tokens")
	}
	for _, c := range []struct {
		path string
		line int
		want string
	}{
		{"color.a", 7, "{color.b} does not resolve: it goes round in a circle, color.a → color.b → color.a"},
		{"color.b", 8, "{color.a} does not resolve: it goes round in a circle"},
		{"color.lost", 9, "{color.nowhere} does not resolve: no token has the path color.nowhere"},
		{"radius.pill", 11, "it is a dimension token, and the token it names, color.primary, is a color token"},
	} {
		got := at(sys.Problems, c.path)
		if len(got) != 1 || got[0].Line != c.line || !strings.Contains(got[0].Message, c.want) {
			t.Errorf("problems at %s = %+v, want one at line %d containing %q", c.path, got, c.line, c.want)
		}
	}
	if lost, _ := s.Get("color.lost"); lost.Value.Text != "{color.nowhere}" {
		t.Errorf("color.lost = %+v, want its alias value kept", lost)
	}
}

func TestLoadMissingOrBrokenFile(t *testing.T) {
	fsys := fstest.MapFS{"broken.json": {Data: []byte("{\n  \"color\": \n}")}}
	sys, err := Load(fsys, "absent.json")
	if !errors.Is(err, fs.ErrNotExist) || sys != nil {
		t.Errorf("Load(absent) = %+v, %v; want an error wrapping fs.ErrNotExist", sys, err)
	}
	sys, err = Load(fsys, "broken.json")
	if err != nil || sys.Themes != nil || len(sys.Problems) != 1 || sys.Problems[0].Line != 3 || !strings.Contains(sys.Problems[0].Message, "does not parse") {
		t.Errorf("Load(broken) = %+v, %v; want its one parse problem, at line 3", sys, err)
	}
}

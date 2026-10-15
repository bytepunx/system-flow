package designsystem

import (
	"errors"
	"io/fs"
	"os"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

// fullFile is a valid token file with every required token, extra colours,
// an alias, and both forms of each value.
func fullFile(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/full.tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func parsePlain(t *testing.T, src []byte) *File {
	t.Helper()
	f, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if f.Resolver || f.Tokens == nil {
		t.Fatalf("Parse read a plain file as a resolver document")
	}
	return f
}

// at is the problems at path.
func at(probs []Problem, path string) []Problem {
	var out []Problem
	for _, p := range probs {
		if p.Path == path {
			out = append(out, p)
		}
	}
	return out
}

func TestFullFileHasNoProblems(t *testing.T) {
	f := parsePlain(t, fullFile(t))
	if probs := Validate(f.Tokens); len(probs) != 0 {
		t.Fatalf("Validate on a valid file = %+v, want none", probs)
	}
	for _, path := range RequiredPaths() {
		if _, ok := f.Tokens.Get(path); !ok {
			t.Errorf("the full file has no %s", path)
		}
	}
	for _, c := range []struct {
		path, typ string
		line      int
	}{
		{"color.background", TypeColor, 5},
		{"text.xs", TypeTypography, 34},
		{"text.lg", TypeDimension, 37},
		{"spacing", TypeDimension, 43},
		{"shadow.lg", TypeShadow, 57},
	} {
		tok, _ := f.Tokens.Get(c.path)
		if tok.Type != c.typ || tok.Line != c.line {
			t.Errorf("%s: type %q line %d, want %q line %d", c.path, tok.Type, tok.Line, c.typ, c.line)
		}
	}
	if link, _ := f.Tokens.Get("color.link"); link.Ref != "color.primary" {
		t.Errorf("color.link Ref = %q, want color.primary", link.Ref)
	}
	if got := f.Tokens.Groups; got["color"] != 2 || got["radius"] != 44 || got["shadow"] != 50 {
		t.Errorf("Groups = %v, want color at 2, radius at 44, shadow at 50", got)
	}
}

func TestDecodedValues(t *testing.T) {
	f := parsePlain(t, fullFile(t))
	get := func(path string) *Node {
		tok, ok := f.Tokens.Get(path)
		if !ok {
			t.Fatalf("no %s", path)
		}
		return tok.Value
	}
	c, err := DecodeColor(get("color.surface"))
	if err != nil || c.Space != ColorSpaceSRGB || c.Hex != "#f5f5f7" || c.Alpha != 1 {
		t.Errorf("color.surface = %+v, %v", c, err)
	}
	c, err = DecodeColor(get("color.overlay"))
	if err != nil || !c.NoHue || c.Alpha != 0.5 || c.Hex != "" {
		t.Errorf("color.overlay = %+v, %v", c, err)
	}
	fam, err := DecodeFontFamily(get("font.mono"))
	if err != nil || !slices.Equal(fam, []string{"ui-monospace"}) {
		t.Errorf("font.mono = %v, %v", fam, err)
	}
	ty, err := DecodeTypography(get("text.sm"))
	if err != nil || ty.LineHeight.Length == nil || *ty.LineHeight.Length != (Dimension{1.25, UnitRem}) {
		t.Errorf("text.sm = %+v, %v", ty, err)
	}
	ty, err = DecodeTypography(get("text.xl"))
	if err != nil || ty.FontWeight != 600 || ty.LineHeight.Number != 1.4 {
		t.Errorf("text.xl = %+v, %v", ty, err)
	}
	sh, err := DecodeShadow(get("shadow.md"))
	if err != nil || len(sh) != 2 || sh[1].Spread != (Dimension{-2, UnitPx}) {
		t.Errorf("shadow.md = %+v, %v", sh, err)
	}
}

func TestValidateProblems(t *testing.T) {
	cases := []struct {
		name, src string
		path      string
		line      int
		want      string
	}{
		{"type flai does not read", `{
  "color": {
    "$type": "number",
    "x": { "$value": 1 }
  }
}`, "color.x", 4, `does not read the $type "number"`},
		{"no type", `{
  "color": {
    "x": { "$value": { "colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff" } }
  }
}`, "color.x", 3, "no $type"},
		{"outside the groups", `{
  "size": {
    "$type": "dimension",
    "big": { "$value": { "value": 2, "unit": "rem" } }
  }
}`, "size.big", 4, "only in the groups"},
		{"type does not fit the group", `{
  "radius": {
    "pill": { "$type": "color", "$value": { "colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff" } }
  }
}`, "radius.pill", 3, "the radius group holds dimension tokens"},
		{"text step is neither typography nor dimension", `{
  "text": {
    "huge": { "$type": "fontFamily", "$value": "Inter" }
  }
}`, "text.huge", 3, "typography or dimension"},
		{"spacing as a group", `{
  "spacing": {
    "$type": "dimension",
    "sm": { "$value": { "value": 4, "unit": "px" } }
  }
}`, "spacing.sm", 4, "spacing is one dimension token"},
		{"a group's own value", `{
  "radius": { "$type": "dimension", "$value": { "value": 4, "unit": "px" } }
}`, "radius", 2, "group of named tokens"},
		{"nested too deep", `{
  "color": {
    "$type": "color",
    "brand": {
      "sun": { "$value": { "colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff" } }
    }
  }
}`, "color.brand.sun", 5, "name it color.brand-sun"},
		{"colour space", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "display-p3", "components": [1, 0, 0] } }
  }
}`, "color.x", 4, `colour space "display-p3"`},
		{"colour as a string", `{
  "color": {
    "$type": "color",
    "x": { "$value": "#ffffff" }
  }
}`, "color.x", 4, "a colour is an object"},
		{"srgb without its hex", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "srgb", "components": [1, 0, 0] } }
  }
}`, "color.x", 4, "needs its hex, the six-digit #ff0000"},
		{"three-digit hex", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "srgb", "components": [1, 1, 1], "hex": "#fff" } }
  }
}`, "color.x", 4, "six-digit"},
		{"hex and components disagree", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "srgb", "components": [1, 0.5, 0], "hex": "#ff0000" } }
  }
}`, "color.x", 4, "give #ff8000"},
		{"srgb component out of range", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "srgb", "components": [255, 0, 0], "hex": "#ff0000" } }
  }
}`, "color.x", 4, "0 to 1"},
		{"two components", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "oklch", "components": [0.5, 0.1] } }
  }
}`, "color.x", 4, "three numbers"},
		{"oklch hue out of range", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "oklch", "components": [0.5, 0.1, 400] } }
  }
}`, "color.x", 4, "hue is 400"},
		{"oklch none outside the hue", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "oklch", "components": ["none", 0.1, 120] } }
  }
}`, "color.x", 4, "component 1 is a string"},
		{"oklch negative chroma", `{
  "color": {
    "$type": "color",
    "x": { "$value": { "colorSpace": "oklch", "components": [0.5, -0.1, 120] } }
  }
}`, "color.x", 4, "chroma"},
		{"alpha on a role", `{
  "color": {
    "$type": "color",
    "primary": { "$value": { "colorSpace": "oklch", "components": [0.5, 0.1, 120], "alpha": 0.8 } }
  }
}`, "color.primary", 4, "alpha is only for shadows and extra tokens"},
		{"dimension unit", `{
  "radius": {
    "$type": "dimension",
    "sm": { "$value": { "value": 0.25, "unit": "em" } }
  }
}`, "radius.sm", 4, "rem or px"},
		{"dimension as a string", `{
  "spacing": { "$type": "dimension", "$value": "0.25rem" }
}`, "spacing", 2, "a dimension is an object"},
		{"font family empty", `{
  "font": {
    "$type": "fontFamily",
    "sans": { "$value": [] }
  }
}`, "font.sans", 4, "non-empty array of names"},
		{"typography without lineHeight", `{
  "text": {
    "$type": "typography",
    "xs": { "$value": { "fontSize": { "value": 0.75, "unit": "rem" } } }
  }
}`, "text.xs", 4, "needs its lineHeight"},
		{"typography without fontSize", `{
  "text": {
    "$type": "typography",
    "xs": { "$value": { "lineHeight": 1.5 } }
  }
}`, "text.xs", 4, "needs its fontSize"},
		{"typography weight", `{
  "text": {
    "$type": "typography",
    "xs": { "$value": { "fontSize": { "value": 0.75, "unit": "rem" }, "lineHeight": 1.5, "fontWeight": "chunky" } }
  }
}`, "text.xs", 4, "fontWeight"},
		{"shadow without spread", `{
  "shadow": {
    "$type": "shadow",
    "sm": { "$value": { "color": { "colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000", "alpha": 0.1 }, "offsetX": { "value": 0, "unit": "px" }, "offsetY": { "value": 1, "unit": "px" }, "blur": { "value": 2, "unit": "px" } } }
  }
}`, "shadow.sm", 4, "needs its spread"},
		{"shadow layer colour space", `{
  "shadow": {
    "$type": "shadow",
    "sm": { "$value": [
      { "color": { "colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000" }, "offsetX": { "value": 0, "unit": "px" }, "offsetY": { "value": 1, "unit": "px" }, "blur": { "value": 2, "unit": "px" }, "spread": { "value": 0, "unit": "px" } },
      { "color": { "colorSpace": "hsl", "components": [0, 0, 0] }, "offsetX": { "value": 0, "unit": "px" }, "offsetY": { "value": 1, "unit": "px" }, "blur": { "value": 2, "unit": "px" }, "spread": { "value": 0, "unit": "px" } }
    ] }
  }
}`, "shadow.sm", 4, `layer 2: color: flai does not read the colour space "hsl"`},
		{"malformed alias", `{
  "color": {
    "$type": "color",
    "link": { "$value": "{color.primary" }
  }
}`, "color.link", 4, "malformed"},
		{"neither token nor group", `{
  "color": {
    "$type": "color",
    "x": "#ffffff"
  }
}`, "color.x", 4, "neither a token nor a group"},
		{"group property flai does not read", `{
  "color": {
    "$type": "color",
    "$extends": "{palette}"
  }
}`, "color", 4, "group property $extends"},
		{"token property flai does not read", `{
  "radius": {
    "$type": "dimension",
    "sm": {
      "$value": { "value": 4, "unit": "px" },
      "$root": true
    }
  }
}`, "radius.sm", 6, "token property $root"},
		{"token holding a token", `{
  "radius": {
    "$type": "dimension",
    "sm": {
      "$value": { "value": 4, "unit": "px" },
      "half": { "$value": { "value": 2, "unit": "px" } }
    }
  }
}`, "radius.sm.half", 6, "a token holds no tokens"},
		{"name with a dot", `{
  "radius": {
    "$type": "dimension",
    "s.m": { "$value": { "value": 4, "unit": "px" } }
  }
}`, "radius.s.m", 4, "rename it"},
		{"name not safe in a CSS custom property", `{
  "color": {
    "$type": "color",
    "brandSun": { "$value": { "colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff" } }
  }
}`, "color.brandSun", 4, "rename it brand-sun"},
		{"name with nothing safe in it", `{
  "radius": {
    "$type": "dimension",
    "__": { "$value": { "value": 4, "unit": "px" } }
  }
}`, "radius.__", 4, "rename it in lower-case letters"},
		{"type not a string", `{
  "radius": {
    "$type": 3,
    "sm": { "$value": { "value": 4, "unit": "px" } }
  }
}`, "radius", 3, "$type is a number"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := at(Validate(parsePlain(t, []byte(c.src)).Tokens), c.path)
			if len(got) != 1 {
				t.Fatalf("problems at %s = %+v, want one", c.path, got)
			}
			if got[0].Line != c.line || !strings.Contains(got[0].Message, c.want) {
				t.Errorf("problem = line %d %q, want line %d containing %q", got[0].Line, got[0].Message, c.line, c.want)
			}
		})
	}
}

func TestSafeName(t *testing.T) {
	for name, want := range map[string]string{
		"brandSun": "brand-sun", "Brand_Sun": "brand-sun", "2XL": "2xl", "-x": "x",
		"on primary!": "on-primary", "HTMLColor": "htmlcolor", "__": "",
	} {
		if got := safeName(name); got != want || (got != "" && !CSSSafe(got)) {
			t.Errorf("safeName(%q) = %q, want %q", name, got, want)
		}
	}
	for _, name := range []string{"2xl", "text-muted", "a", "0"} {
		if !CSSSafe(name) {
			t.Errorf("CSSSafe(%q) = false, want true", name)
		}
	}
}

func TestExtrasMayHaveAlphaAndAliasesAreLeft(t *testing.T) {
	src := `{
  "color": {
    "$type": "color",
    "scrim": { "$value": { "colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000", "alpha": 0.4 } },
    "link": { "$value": "{color.primary}" },
    "loose": { "$type": "color", "$value": "{nowhere}" }
  },
  "radius": { "pill": { "$value": "{radius.lg}" } }
}`
	s := parsePlain(t, []byte(src)).Tokens
	probs := Validate(s)
	for _, path := range []string{"color.scrim", "color.link", "color.loose", "radius.pill"} {
		if got := at(probs, path); len(got) != 0 {
			t.Errorf("problems at %s = %+v, want none", path, got)
		}
	}
	if pill, _ := s.Get("radius.pill"); pill.Ref != "radius.lg" || pill.Type != "" {
		t.Errorf("radius.pill = %+v, want an untyped alias of radius.lg", pill)
	}
}

// without removes the member at path from the tree under n.
func without(n *Node, path ...string) {
	for _, key := range path[:len(path)-1] {
		n = n.Get(key)
	}
	n.Members = slices.DeleteFunc(n.Members, func(m Member) bool { return m.Key == path[len(path)-1] })
}

func TestMissingRequiredTokens(t *testing.T) {
	cases := []struct {
		name   string
		remove []string
		want   []Problem
	}{
		{"a colour role", []string{"color", "on-primary"}, []Problem{{Path: "color.on-primary", Line: 2}}},
		{"a type step", []string{"text", "2xl"}, []Problem{{Path: "text.2xl", Line: 32}}},
		{"spacing", []string{"spacing"}, []Problem{{Path: "spacing", Line: 1}}},
		{"the shadow group", []string{"shadow"}, []Problem{{Path: "shadow.sm", Line: 1}, {Path: "shadow.md", Line: 1}, {Path: "shadow.lg", Line: 1}}},
		{"a font", []string{"font", "mono"}, []Problem{{Path: "font.mono", Line: 27}}},
		{"a radius", []string{"radius", "md"}, []Problem{{Path: "radius.md", Line: 44}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := parsePlain(t, fullFile(t))
			without(f.Root, c.remove...)
			probs := Validate(ReadTokens(f.Root))
			if len(probs) != len(c.want) {
				t.Fatalf("Validate = %+v, want %d problems", probs, len(c.want))
			}
			for i, p := range probs {
				if p.Path != c.want[i].Path || p.Line != c.want[i].Line || !strings.Contains(p.Message, "is required") {
					t.Errorf("problem %d = %+v, want %s at line %d, required", i, p, c.want[i].Path, c.want[i].Line)
				}
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		name, src string
		line      int
		want      string
	}{
		{"trailing comma", "{\n  \"color\": {\n    \"$type\": \"color\",\n  }\n}", 4, "invalid character '}'"},
		{"missing comma", "{\n  \"a\": 1\n\"b\": 2\n}", 3, "invalid character '\"' after object key:value pair"},
		{"bad literal", "{\n  \"a\":\nnope\n}", 3, "invalid character 'o' in literal null"},
		{"unclosed", "{\n  \"color\": {\n    \"$type\": \"color\"\n", 3, "ends before its value does"},
		{"key twice", "{\n  \"color\": {},\n  \"font\": {},\n  \"color\": {}\n}", 4, `"color" is given twice in one object, first at line 2`},
		{"data after the document", "{\n  \"color\": {}\n}\n{}\n", 4, "more after"},
		{"not an object", "\n[1, 2]\n", 2, "is an array, not an object"},
		{"not JSON", "color:\n  primary: red\n", 1, "invalid character"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse([]byte(c.src))
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse error = %v, want a *ParseError", err)
			}
			if pe.Line != c.line || !strings.Contains(pe.Msg, c.want) {
				t.Errorf("ParseError = line %d %q, want line %d containing %q", pe.Line, pe.Msg, c.line, c.want)
			}
			if p := pe.Problem(); p.Path != "" || p.Line != c.line || !strings.Contains(p.Message, "does not parse") {
				t.Errorf("Problem = %+v", p)
			}
		})
	}
}

func TestParseJSONKeepsLinesAndOrder(t *testing.T) {
	n, err := ParseJSON([]byte("{\n  \"b\": [1,\n    \"two\"],\n  \"a\": {\"x\": null, \"y\": true}\n}"))
	if err != nil {
		t.Fatal(err)
	}
	if keys := []string{n.Members[0].Key, n.Members[1].Key}; !slices.Equal(keys, []string{"b", "a"}) {
		t.Errorf("keys = %v, want b then a", keys)
	}
	b := n.Get("b")
	if b.Line != 2 || b.Items[0].Number != 1 || b.Items[1].Line != 3 || b.Items[1].Text != "two" {
		t.Errorf("b = %+v", b)
	}
	if a := n.Get("a"); n.Members[1].Line != 4 || a.Get("x").Kind != KindNull || !a.Get("y").Bool {
		t.Errorf("a = %+v", a)
	}
	if n.Get("missing") != nil || b.Get("x") != nil {
		t.Error("Get answered a member that is not there")
	}
}

func TestRead(t *testing.T) {
	resolver := `{
  "$schema": "https://www.designtokens.org/schemas/2025.10/resolver.json",
  "version": "2025.10",
  "sets": { "tokens": { "sources": [ { "color": { "$type": "color" } } ] } },
  "modifiers": { "theme": { "contexts": { "light": [], "dark": [] }, "default": "light" } },
  "resolutionOrder": [ { "$ref": "#/sets/tokens" }, { "$ref": "#/modifiers/theme" } ]
}`
	fsys := fstest.MapFS{
		"resolver.json": {Data: []byte(resolver)},
		"plain.json":    {Data: fullFile(t)},
		"broken.json":   {Data: []byte("{\n  \"color\": \n}")},
	}
	f, err := Read(fsys, "resolver.json")
	if err != nil || !f.Resolver || f.Tokens != nil || f.Root.Get("sets").Line != 4 {
		t.Errorf("Read(resolver) = %+v, %v; want a resolver document with its root", f, err)
	}
	if f, err := Read(fsys, "plain.json"); err != nil || f.Resolver || len(f.Tokens.Tokens) == 0 {
		t.Errorf("Read(plain) = %+v, %v; want its tokens", f, err)
	}
	_, err = Read(fsys, "absent.json")
	if !errors.Is(err, fs.ErrNotExist) || !strings.Contains(err.Error(), "absent.json") {
		t.Errorf("Read(absent) = %v, want fs.ErrNotExist naming the file", err)
	}
	_, err = Read(fsys, "broken.json")
	var pe *ParseError
	if !errors.As(err, &pe) || pe.Line != 3 || !strings.Contains(err.Error(), "broken.json") {
		t.Errorf("Read(broken) = %v, want a ParseError at line 3 naming the file", err)
	}
}

func TestHeadings(t *testing.T) {
	md := "---\ntitle: Design system\n# a yaml comment, not a heading\ntopics: [flaiover]\n---\n\n" +
		"# Design system\n\nProse.\n\n## Colour ##\n\n```md\n# not a heading\n```\n\n### Type <!-- topics: ui -->\n####### seven is not a heading\n#also not\n"
	got, err := Headings([]byte(md))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"Design system", "Colour", "Type"}; !slices.Equal(got, want) {
		t.Errorf("Headings = %q, want %q", got, want)
	}
	if _, err := Headings([]byte("---\ntopics: [unclosed\n---\n# A\n")); err == nil {
		t.Error("Headings on front matter that does not parse: no error")
	}
}

func TestRequiredPaths(t *testing.T) {
	paths := RequiredPaths()
	if len(paths) != len(ColorRoles)+len(FontNames)+len(TextSteps)+1+len(RadiusNames)+len(ShadowNames) {
		t.Errorf("RequiredPaths has %d paths", len(paths))
	}
	for _, p := range []string{"color.on-secondary", "font.mono", "text.4xl", "spacing", "radius.lg", "shadow.sm"} {
		if !slices.Contains(paths, p) {
			t.Errorf("RequiredPaths lacks %s", p)
		}
	}
	if !IsRole("info-soft") || IsRole("brand-sun") {
		t.Error("IsRole is wrong about info-soft or brand-sun")
	}
}

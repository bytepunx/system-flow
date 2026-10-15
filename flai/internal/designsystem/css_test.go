package designsystem

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata")

// golden compares got with the golden file testdata/name, rewriting it with
// -update. A golden file that is missing is written, and the test fails so
// that it is reviewed.
func golden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	want, err := os.ReadFile(path)
	if *update || errors.Is(err, fs.ErrNotExist) {
		if werr := os.WriteFile(path, got, 0o644); werr != nil {
			t.Fatal(werr)
		}
		if !*update {
			t.Errorf("wrote %s, which was missing; review it and run the test again", path)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s differs from what was rendered; if the change is intended, run go test -update and review the diff\ngot:\n%s", path, got)
	}
}

func loadFixture(t *testing.T, name string) *System {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return load(t, data)
}

func TestRenderingsGolden(t *testing.T) {
	for _, fixture := range []string{"full", "themes"} {
		t.Run(fixture, func(t *testing.T) {
			sys := loadFixture(t, fixture+".tokens.json")
			if len(sys.Problems) != 0 {
				t.Fatalf("problems = %+v, want none", sys.Problems)
			}
			golden(t, fixture+".css.golden", []byte(CSS(sys)))
			golden(t, fixture+".tailwind.css.golden", []byte(Tailwind(sys)))
		})
	}
}

func TestTailwindDarkVariantFollowsTheDarkSelector(t *testing.T) {
	const want = "@custom-variant dark (&:where([data-theme='dark'], [data-theme='dark'] *));\n\n@theme {\n"
	for _, fixture := range []string{"full", "themes"} {
		if got := Tailwind(loadFixture(t, fixture+".tokens.json")); !strings.HasPrefix(got, want) {
			t.Errorf("%s: Tailwind does not open with the dark variant, then @theme:\n%s", fixture, got)
		}
	}
}

func TestTailwindDarkOverridesOnlyWhatDiffers(t *testing.T) {
	plain := Tailwind(loadFixture(t, "full.tokens.json"))
	if strings.Contains(plain, DarkSelector+" {") {
		t.Errorf("a plain file's Tailwind rendering has a dark rule:\n%s", plain)
	}
	themed := Tailwind(loadFixture(t, "themes.tokens.json"))
	_, dark, ok := strings.Cut(themed, DarkSelector+" {")
	if !ok {
		t.Fatalf("no dark rule in:\n%s", themed)
	}
	for _, want := range []string{"--color-background:", "--color-link:", "--shadow-md:"} {
		if !strings.Contains(dark, want) {
			t.Errorf("the dark rule has no %s:\n%s", want, dark)
		}
	}
	for _, unwanted := range []string{"--color-secondary:", "--font-sans:", "--spacing:", "--shadow-sm:"} {
		if strings.Contains(dark, unwanted) {
			t.Errorf("the dark rule has %s, which dark does not change:\n%s", unwanted, dark)
		}
	}
}

func TestCSSEveryBlockStandsAlone(t *testing.T) {
	css := CSS(loadFixture(t, "themes.tokens.json"))
	light, dark, ok := strings.Cut(css, DarkSelector)
	if !ok || !strings.HasPrefix(light, ":root {\n") || !strings.HasSuffix(css, "}\n") {
		t.Fatalf("CSS is not a :root rule then a dark rule:\n%s", css)
	}
	if a, b := strings.Count(light, ";\n"), strings.Count(dark, ";\n"); a != b {
		t.Errorf("the light rule has %d properties and the dark %d, want the same", a, b)
	}
}

// token is a one-token file's token, read and resolved as a plain file.
func token(t *testing.T, group, member string) Token {
	t.Helper()
	sys := load(t, []byte(`{"`+group+`": {`+member+`}}`))
	toks := sys.Themes[0].Tokens.Tokens
	if len(toks) != 1 {
		t.Fatalf("tokens = %+v, want one", toks)
	}
	return toks[0]
}

func TestDeclarations(t *testing.T) {
	cases := []struct {
		name, group, member string
		want                []Declaration
	}{
		{"srgb from its components", "color", `"x": {"$type": "color", "$value": {"colorSpace": "srgb", "components": [1, 0.5, 0], "hex": "#FF8000"}}`,
			[]Declaration{{"--color-x", "#ff8000"}}},
		{"srgb with alpha", "color", `"x": {"$type": "color", "$value": {"colorSpace": "srgb", "components": [0, 0, 0], "hex": "#000000", "alpha": 0.5}}`,
			[]Declaration{{"--color-x", "#00000080"}}},
		{"oklch as Tailwind writes it", "color", `"x": {"$type": "color", "$value": {"colorSpace": "oklch", "components": [0.637, 0.237, 25.331]}}`,
			[]Declaration{{"--color-x", "oklch(63.7% 0.237 25.331)"}}},
		{"oklch with no hue and alpha", "color", `"x": {"$type": "color", "$value": {"colorSpace": "oklch", "components": [0, 0, "none"], "alpha": 0.25}}`,
			[]Declaration{{"--color-x", "oklch(0% 0 none / 0.25)"}}},
		{"font families quoted where CSS needs it", "font", `"x": {"$type": "fontFamily", "$value": ["Inter", "Inter Variable", "1Password", "inherit", "My \"Font\"", "ui-rounded", "emoji"]}`,
			[]Declaration{{"--font-x", `Inter, "Inter Variable", "1Password", "inherit", "My \"Font\"", ui-rounded, emoji`}}},
		{"one font family", "font", `"x": {"$type": "fontFamily", "$value": "Georgia"}`,
			[]Declaration{{"--font-x", "Georgia"}}},
		{"typography with every property", "text", `"x": {"$type": "typography", "$value": {"fontSize": {"value": 1.250, "unit": "rem"}, "lineHeight": {"value": 28, "unit": "px"}, "fontWeight": "bold", "letterSpacing": {"value": -0.5, "unit": "px"}, "fontFamily": "Inter"}}`,
			[]Declaration{{"--text-x", "1.25rem"}, {"--text-x--line-height", "28px"}, {"--text-x--font-weight", "700"}, {"--text-x--letter-spacing", "-0.5px"}}},
		{"typography with a unitless line height", "text", `"x": {"$type": "typography", "$value": {"fontSize": {"value": 12, "unit": "px"}, "lineHeight": 1.5}}`,
			[]Declaration{{"--text-x", "12px"}, {"--text-x--line-height", "1.5"}}},
		{"a dimension text step", "text", `"x": {"$type": "dimension", "$value": {"value": 1.125, "unit": "rem"}}`,
			[]Declaration{{"--text-x", "1.125rem"}}},
		{"radius", "radius", `"full": {"$type": "dimension", "$value": {"value": 9999, "unit": "px"}}`,
			[]Declaration{{"--radius-full", "9999px"}}},
		{"inset shadow", "shadow", `"x": {"$type": "shadow", "$value": {"color": {"colorSpace": "oklch", "components": [0.5, 0.1, 120], "alpha": 0.2}, "offsetX": {"value": 0, "unit": "px"}, "offsetY": {"value": 2, "unit": "px"}, "blur": {"value": 4, "unit": "px"}, "spread": {"value": 0, "unit": "px"}, "inset": true}}`,
			[]Declaration{{"--shadow-x", "inset 0px 2px 4px 0px oklch(50% 0.1 120 / 0.2)"}}},
		{"a name CSS cannot hold", "color", `"brandSun": {"$type": "color", "$value": {"colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff"}}`, nil},
		{"a value that does not decode", "radius", `"x": {"$type": "dimension", "$value": {"value": 4, "unit": "em"}}`, nil},
		{"a type its group does not hold", "radius", `"x": {"$type": "color", "$value": {"colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff"}}`, nil},
		{"outside the groups", "size", `"x": {"$type": "dimension", "$value": {"value": 4, "unit": "px"}}`, nil},
		{"an alias that does not resolve", "color", `"x": {"$type": "color", "$value": "{color.nowhere}"}`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tok := token(t, c.group, c.member)
			if got := Declarations(tok); !slices.Equal(got, c.want) {
				t.Errorf("Declarations = %q, want %q", got, c.want)
			}
		})
	}
	spacing := load(t, []byte(`{"spacing": {"$type": "dimension", "$value": {"value": 0.250, "unit": "rem"}}}`)).Themes[0].Tokens.Tokens[0]
	if got := Declarations(spacing); !slices.Equal(got, []Declaration{{"--spacing", "0.25rem"}}) {
		t.Errorf("Declarations(spacing) = %q, want --spacing: 0.25rem", got)
	}
}

func TestRenderingsOfAFileThatDoesNotParse(t *testing.T) {
	sys := load(t, []byte(`{"color": `))
	if CSS(sys) != "" || Tailwind(sys) != "" {
		t.Errorf("CSS = %q, Tailwind = %q, want both empty", CSS(sys), Tailwind(sys))
	}
}

const guidance = `---
title: Design system
topics: [ui]
---

# Design system

## Colour

## Type
`

// show runs Show on the token file src, "" for none, beside the guidance.
func show(t *testing.T, src []byte, withGuidance bool) *Shown {
	t.Helper()
	fsys := fstest.MapFS{}
	if src != nil {
		fsys["design/system/design-system.tokens.json"] = &fstest.MapFile{Data: src}
	}
	if withGuidance {
		fsys["design/system/design-system.md"] = &fstest.MapFile{Data: []byte(guidance)}
	}
	s, err := Show(fsys, "design/system/design-system.tokens.json", "design/system/design-system.md")
	if err != nil {
		t.Fatalf("Show: %v", err)
	}
	return s
}

func goldenJSON(t *testing.T, name string, v any) {
	t.Helper()
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	golden(t, name, append(data, '\n'))
}

func TestShowGolden(t *testing.T) {
	themes, err := os.ReadFile("testdata/themes.tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	problems := []byte(`{
  "color": {
    "$type": "color",
    "brandSun": { "$value": { "colorSpace": "srgb", "components": [0.961, 0.718, 0], "hex": "#f5b700" } },
    "text": { "$value": { "colorSpace": "oklch", "components": [0.21, 0.02, 265] } },
    "background": { "$value": { "colorSpace": "srgb", "components": [1, 1, 1], "hex": "#ffffff" } },
    "link": { "$value": "{color.nowhere}" }
  },
  "spacing": { "$type": "dimension", "$value": "0.25rem" }
}
`)
	for _, c := range []struct {
		name string
		src  []byte
	}{{"show.json.golden", themes}, {"show-problems.json.golden", problems}, {"show-missing.json.golden", nil}} {
		t.Run(c.name, func(t *testing.T) { goldenJSON(t, c.name, show(t, c.src, true)) })
	}
}

func TestShowSaysWhenThereIsNoTokenFile(t *testing.T) {
	s := show(t, nil, true)
	if s.Exists || !strings.Contains(s.NoTokenFile, "flai design init") || s.Guidance != "" || len(s.Headings) != 0 || s.CSS != "" {
		t.Errorf("Show with no token file = %+v, want exists false, a message naming flai design init, and nothing else", s)
	}
}

func TestShowWithoutGuidance(t *testing.T) {
	s := show(t, fullFile(t), false)
	if !s.Exists || s.Headings == nil || len(s.Headings) != 0 {
		t.Errorf("headings = %#v, want empty", s.Headings)
	}
	if len(s.Contrast) != 2 || len(s.Contrast[0].Pairs) == 0 {
		t.Errorf("contrast = %+v, want light and dark pairs", s.Contrast)
	}
	for _, p := range s.Problems {
		if strings.Contains(p.Message, "contrast") {
			t.Errorf("contrast is a problem: %+v", p)
		}
	}
}

func TestShowOfAFileThatDoesNotParse(t *testing.T) {
	s := show(t, []byte(`{"color": `), true)
	if !s.Exists || len(s.Problems) != 1 || len(s.Themes) != 0 || s.CSS != "" || s.Tailwind != "" || len(s.Headings) != 3 {
		t.Errorf("Show of a file that does not parse = %+v", s)
	}
}

// deniedFS refuses every read.
type deniedFS struct{}

func (deniedFS) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
}

func TestShowReadErrors(t *testing.T) {
	if _, err := Show(deniedFS{}, "tokens.json", "guidance.md"); !errors.Is(err, fs.ErrPermission) {
		t.Errorf("Show on an unreadable token file: err = %v, want ErrPermission", err)
	}
	fsys := fstest.MapFS{"tokens.json": {Data: fullFile(t)}, "guidance.md": {Data: []byte("---\ntitle: [\n---\n# x\n")}}
	if _, err := Show(fsys, "tokens.json", "guidance.md"); err == nil || !strings.Contains(err.Error(), "guidance.md") {
		t.Errorf("Show with guidance whose front matter does not parse: err = %v, want one naming it", err)
	}
}

func TestNodeJSON(t *testing.T) {
	n, err := ParseJSON([]byte(`{"a": [1, "x", true, null], "b": {"c": 0.5}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(n.JSON())
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"a":[1,"x",true,null],"b":{"c":0.5}}` {
		t.Errorf("JSON = %s", got)
	}
}

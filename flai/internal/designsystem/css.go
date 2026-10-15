package designsystem

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Declaration is one CSS custom property and its value.
type Declaration struct {
	Name  string `json:"name"`  // such as --color-primary
	Value string `json:"value"` // such as oklch(52% 0.18 255)
}

// Declarations answers the CSS custom properties a token renders as, named
// as Tailwind v4's theme namespaces name them, with its resolved value: one
// for most tokens, and for a typography step its size, line height, and the
// weight and letter spacing it gives. A token renders nothing, nil, when it
// sits outside the groups, its name is not CSSSafe, or its value does not
// decode as a type its group holds; Validate reports each of those.
func Declarations(t Token) []Declaration {
	if placement(t) != "" || !CSSSafe(t.Name) || !slices.Contains(groupTypes[t.Group], t.Type) || t.Value == nil {
		return nil
	}
	name := "--" + t.Group
	if t.Group != GroupSpacing {
		name += "-" + t.Name
	}
	if t.Type == TypeTypography {
		decls, err := typographyCSS(name, t.Value)
		if err != nil {
			return nil
		}
		return decls
	}
	v, err := valueCSS[t.Type](t.Value)
	if err != nil {
		return nil
	}
	return []Declaration{{name, v}}
}

// valueCSS renders a value of each type but typography as one CSS value.
var valueCSS = map[string]func(*Node) (string, error){
	TypeColor:      colorCSS,
	TypeFontFamily: fontFamilyCSS,
	TypeDimension:  dimensionCSS,
	TypeShadow:     shadowCSS,
}

// CSSValue is the value of the first property a token renders as: a
// typography step's is its size; "" when it renders nothing.
func CSSValue(t Token) string {
	if d := Declarations(t); len(d) > 0 {
		return d[0].Value
	}
	return ""
}

// colorCSS is a colour as CSS: srgb as its #rrggbb, or #rrggbbaa with alpha
// below 1; oklch as Tailwind v4's palette writes it, oklch(63.7% 0.237
// 25.331), with "/ alpha" below 1.
func colorCSS(n *Node) (string, error) {
	c, err := DecodeColor(n)
	if err != nil {
		return "", err
	}
	return colorString(c), nil
}

func colorString(c Color) string {
	if c.Space == ColorSpaceSRGB {
		rgb, _ := c.SRGB()
		hex := rgb.Hex()
		if c.Alpha < 1 {
			hex += fmt.Sprintf("%02x", byte8(c.Alpha))
		}
		return hex
	}
	hue := "none"
	if !c.NoHue {
		hue = number(c.Components[2])
	}
	s := fmt.Sprintf("oklch(%s%% %s %s", number(c.Components[0]*100), number(c.Components[1]), hue)
	if c.Alpha < 1 {
		s += " / " + number(c.Alpha)
	}
	return s + ")"
}

// number writes a number as CSS does, to six decimals at most, with no
// trailing zeros.
func number(v float64) string {
	s := strconv.FormatFloat(math.Round(v*1e6)/1e6, 'f', -1, 64)
	if s == "-0" {
		return "0"
	}
	return s
}

func dimensionCSS(n *Node) (string, error) {
	d, err := DecodeDimension(n)
	if err != nil {
		return "", err
	}
	return dimensionString(d), nil
}

func dimensionString(d Dimension) string { return number(d.Value) + d.Unit }

// genericFamilies are CSS's generic font families, written bare.
var genericFamilies = []string{
	"serif", "sans-serif", "monospace", "cursive", "fantasy", "system-ui",
	"ui-serif", "ui-sans-serif", "ui-monospace", "ui-rounded", "math", "emoji", "fangsong",
}

// cssWideKeywords may not be a family name unquoted.
var cssWideKeywords = []string{"inherit", "initial", "unset", "revert", "revert-layer", "default"}

var identRe = regexp.MustCompile(`^-?[A-Za-z_][A-Za-z0-9_-]*$`)

// fontFamilyCSS is a font stack as CSS: comma-joined, each family that is
// not a generic family and not a plain identifier in double quotes.
func fontFamilyCSS(n *Node) (string, error) {
	fams, err := DecodeFontFamily(n)
	if err != nil {
		return "", err
	}
	out := make([]string, len(fams))
	for i, f := range fams {
		switch {
		case slices.Contains(genericFamilies, f):
			out[i] = f
		case identRe.MatchString(f) && !slices.Contains(cssWideKeywords, strings.ToLower(f)):
			out[i] = f
		default:
			out[i] = `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(f) + `"`
		}
	}
	return strings.Join(out, ", "), nil
}

// typographyCSS is a type step's properties, name being its size's: a plain
// dimension step gives only its size. Tailwind v4 has no per-step font
// family, so a step's fontFamily renders nothing.
func typographyCSS(name string, n *Node) ([]Declaration, error) {
	t, err := DecodeTypography(n)
	if err != nil {
		return nil, err
	}
	lh := number(t.LineHeight.Number)
	if t.LineHeight.Length != nil {
		lh = dimensionString(*t.LineHeight.Length)
	}
	out := []Declaration{{name, dimensionString(t.FontSize)}, {name + "--line-height", lh}}
	if t.FontWeight != 0 {
		out = append(out, Declaration{name + "--font-weight", number(t.FontWeight)})
	}
	if t.LetterSpacing != nil {
		out = append(out, Declaration{name + "--letter-spacing", dimensionString(*t.LetterSpacing)})
	}
	return out, nil
}

// shadowCSS is a shadow as CSS: each layer [inset ]x y blur spread colour,
// comma-joined.
func shadowCSS(n *Node) (string, error) {
	layers, err := DecodeShadow(n)
	if err != nil {
		return "", err
	}
	out := make([]string, len(layers))
	for i, l := range layers {
		s := strings.Join([]string{
			dimensionString(l.OffsetX), dimensionString(l.OffsetY),
			dimensionString(l.Blur), dimensionString(l.Spread), colorString(l.Color),
		}, " ")
		if l.Inset {
			s = "inset " + s
		}
		out[i] = s
	}
	return strings.Join(out, ", "), nil
}

// DarkSelector is the selector the dark theme's properties sit under.
const DarkSelector = "[data-theme='dark']"

// themeDeclarations is a theme's properties: by group in the order Groups
// gives, then in file order within a group.
func themeDeclarations(s *TokenSet) []Declaration {
	var out []Declaration
	for _, g := range Groups {
		for _, t := range s.Tokens {
			if t.Group == g {
				out = append(out, Declarations(t)...)
			}
		}
	}
	return out
}

// block writes one rule: its selector, then each property on its own line.
func block(b *strings.Builder, selector string, decls []Declaration) {
	b.WriteString(selector + " {\n")
	for _, d := range decls {
		fmt.Fprintf(b, "  %s: %s;\n", d.Name, d.Value)
	}
	b.WriteString("}\n")
}

// theme answers the system's theme name; nil when it has none.
func (sys *System) theme(name string) *TokenSet {
	for _, th := range sys.Themes {
		if th.Name == name {
			return th.Tokens
		}
	}
	return nil
}

// CSS renders the tokens as CSS custom properties with Tailwind v4's names:
// a :root rule with the light theme's, then a [data-theme='dark'] rule with
// the dark theme's, each with every token so that either stands alone.
// Properties go by group in the order Groups gives, then in file order. ""
// when the system has no themes, because its file does not parse.
func CSS(sys *System) string {
	light, dark := sys.theme(ThemeLight), sys.theme(ThemeDark)
	if light == nil || dark == nil {
		return ""
	}
	var b strings.Builder
	block(&b, ":root", themeDeclarations(light))
	b.WriteString("\n")
	block(&b, DarkSelector, themeDeclarations(dark))
	return b.String()
}

// darkVariant redefines Tailwind's dark variant to follow DarkSelector, on
// the element that carries it and every element below it.
const darkVariant = "@custom-variant dark (&:where(" + DarkSelector + ", " + DarkSelector + " *));\n"

// Tailwind renders the tokens as a Tailwind v4 stylesheet, as ADR-0159
// decides: a @custom-variant rule that makes the dark: variant follow
// [data-theme='dark'], then a @theme rule with the light values, then a
// [data-theme='dark'] rule overriding the properties whose dark value
// differs, left out when none does. "" when the system has no themes.
func Tailwind(sys *System) string {
	light, dark := sys.theme(ThemeLight), sys.theme(ThemeDark)
	if light == nil || dark == nil {
		return ""
	}
	lightDecls := themeDeclarations(light)
	var b strings.Builder
	b.WriteString(darkVariant + "\n")
	block(&b, "@theme", lightDecls)
	values := make(map[string]string, len(lightDecls))
	for _, d := range lightDecls {
		values[d.Name] = d.Value
	}
	var differ []Declaration
	for _, d := range themeDeclarations(dark) {
		if v, ok := values[d.Name]; !ok || v != d.Value {
			differ = append(differ, d)
		}
	}
	if len(differ) > 0 {
		b.WriteString("\n")
		block(&b, DarkSelector, differ)
	}
	return b.String()
}

// Shown is a project's design system as flai design show --json and the host
// method design.get answer it.
type Shown struct {
	Path        string         `json:"path"`                    // the token file
	Exists      bool           `json:"exists"`                  // false when there is no token file, and every field below is empty
	NoTokenFile string         `json:"no_token_file,omitempty"` // when there is no token file: so, and what to run
	Guidance    string         `json:"guidance"`                // the guidance document's path
	Headings    []string       `json:"headings"`                // the guidance's headings; empty when it is missing
	Resolver    bool           `json:"resolver"`                // the file is a DTCG resolver document
	Themes      []ShownTheme   `json:"themes"`                  // light then dark; empty when the file does not parse
	Problems    []Problem      `json:"problems"`                // the file's problems, by line, then path
	CSS         string         `json:"css"`                     // CSS's rendering; "" when the file does not parse
	Tailwind    string         `json:"tailwind"`                // Tailwind's rendering; "" when the file does not parse
	Contrast    []ShownPairing `json:"contrast"`                // each theme's role pairs; a report, never a problem
}

// ShownTheme is one theme's tokens.
type ShownTheme struct {
	Name   string       `json:"name"`
	Tokens []ShownToken `json:"tokens"` // in file order
}

// ShownToken is one token of a theme, its value resolved.
type ShownToken struct {
	Path      string        `json:"path"`
	Type      string        `json:"type"`
	Value     any           `json:"value"`         // its $value as JSON; an alias's is its target's
	CSS       string        `json:"css"`           // CSSValue: its first property's value; "" when it renders nothing
	Variables []Declaration `json:"variables"`     // the properties it renders as; empty when none
	Ref       string        `json:"ref,omitempty"` // the token path an alias names
	Line      int           `json:"line"`
}

// ShownPairing is one theme's role pair contrast.
type ShownPairing struct {
	Theme string `json:"theme"`
	Pairs []Pair `json:"pairs"`
}

// Show reads the design system for flai design show --json and design.get:
// the token file tokens and the guidance document guidance, in fsys. No
// token file is no error but a Shown that says so; a missing guidance has no
// headings. Any other read error is an error.
func Show(fsys fs.FS, tokens, guidance string) (*Shown, error) {
	sys, err := Load(fsys, tokens)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return &Shown{
			Path:        tokens,
			NoTokenFile: fmt.Sprintf("there is no token file at %s; run flai design init to write one and its guidance", tokens),
			Headings:    []string{}, Themes: []ShownTheme{}, Problems: []Problem{}, Contrast: []ShownPairing{},
		}, nil
	case err != nil:
		return nil, err
	}
	headings, err := guidanceHeadings(fsys, guidance)
	if err != nil {
		return nil, err
	}
	out := &Shown{
		Path: tokens, Exists: true, Guidance: guidance, Headings: headings,
		Resolver: sys.Resolver, Themes: []ShownTheme{}, Problems: sys.Problems,
		CSS: CSS(sys), Tailwind: Tailwind(sys), Contrast: []ShownPairing{},
	}
	if out.Problems == nil {
		out.Problems = []Problem{}
	}
	for _, th := range sys.Themes {
		st := ShownTheme{Name: th.Name, Tokens: make([]ShownToken, 0, len(th.Tokens.Tokens))}
		for _, t := range th.Tokens.Tokens {
			vars := Declarations(t)
			if vars == nil {
				vars = []Declaration{}
			}
			st.Tokens = append(st.Tokens, ShownToken{
				Path: t.Path, Type: t.Type, Value: t.Value.JSON(), CSS: CSSValue(t),
				Variables: vars, Ref: t.Ref, Line: t.Line,
			})
		}
		out.Themes = append(out.Themes, st)
		pairs := Contrast(th.Name, th.Tokens)
		if pairs == nil {
			pairs = []Pair{}
		}
		out.Contrast = append(out.Contrast, ShownPairing{Theme: th.Name, Pairs: pairs})
	}
	return out, nil
}

// guidanceHeadings is the guidance's headings, empty when it is missing.
func guidanceHeadings(fsys fs.FS, name string) ([]string, error) {
	data, err := fs.ReadFile(fsys, name)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return []string{}, nil
	case err != nil:
		return nil, fmt.Errorf("read the design system guidance %s: %w", name, err)
	}
	h, err := Headings(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	if h == nil {
		h = []string{}
	}
	return h, nil
}

// JSON is the value as plain JSON values: nil, bool, float64, string,
// []any, and map[string]any.
func (n *Node) JSON() any {
	if n == nil {
		return nil
	}
	switch n.Kind {
	case KindBool:
		return n.Bool
	case KindNumber:
		return n.Number
	case KindString:
		return n.Text
	case KindArray:
		out := make([]any, len(n.Items))
		for i, item := range n.Items {
			out[i] = item.JSON()
		}
		return out
	case KindObject:
		out := make(map[string]any, len(n.Members))
		for _, m := range n.Members {
			out[m.Key] = m.Value.JSON()
		}
		return out
	}
	return nil
}

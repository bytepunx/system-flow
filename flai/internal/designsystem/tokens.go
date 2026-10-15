// Package designsystem reads a project's design system (ADR-0146, ADR-0148):
// its token file, a W3C Design Tokens Community Group (DTCG) file in format
// 2025.10, in the subset of it flai reads, and the headings of its guidance,
// design/system/design-system.md. It reports each token flai does not read
// at its token path and line.
//
// The token file is a plain DTCG token file, whose values serve both themes,
// or a DTCG resolver document with its sources inline, recognised by a
// top-level resolutionOrder. This file reads JSON with lines and a token
// tree into tokens, and validates them.
package designsystem

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/topics"
)

// Kind is the kind of a JSON value.
type Kind int

// The kinds of JSON value.
const (
	KindNull Kind = iota
	KindBool
	KindNumber
	KindString
	KindArray
	KindObject
)

// String names the kind as a sentence would: a number, an object.
func (k Kind) String() string {
	switch k {
	case KindNull:
		return "null"
	case KindBool:
		return "a boolean"
	case KindNumber:
		return "a number"
	case KindString:
		return "a string"
	case KindArray:
		return "an array"
	case KindObject:
		return "an object"
	}
	return fmt.Sprintf("kind %d", int(k))
}

// Node is a JSON value and the line it starts on.
type Node struct {
	Kind    Kind
	Line    int      // 1-based line of the value's first character
	Bool    bool     // KindBool
	Number  float64  // KindNumber
	Text    string   // KindString
	Items   []*Node  // KindArray, in order
	Members []Member // KindObject, in file order
}

// Member is one key of a JSON object and its value.
type Member struct {
	Key   string
	Line  int // 1-based line of the key
	Value *Node
}

// Get answers the value of the object member key, or nil when n is nil, not
// an object, or has no such key.
func (n *Node) Get(key string) *Node {
	if n == nil || n.Kind != KindObject {
		return nil
	}
	for _, m := range n.Members {
		if m.Key == key {
			return m.Value
		}
	}
	return nil
}

// ParseError is a token file flai cannot read as a document, at the line
// where reading stopped.
type ParseError struct {
	Line int
	Msg  string
}

// Error says what is wrong and at which line.
func (e *ParseError) Error() string { return fmt.Sprintf("line %d: %s", e.Line, e.Msg) }

// Problem is the error as a token file problem, at no token path.
func (e *ParseError) Problem() Problem {
	return Problem{Line: e.Line, Message: "the token file does not parse: " + e.Msg + "; fix the JSON at this line"}
}

// ParseJSON reads one JSON value, keeping each value's line and each
// object's key order. JSON that does not parse, a key given twice in one
// object, or data after the value is a *ParseError.
func ParseJSON(data []byte) (*Node, error) {
	p := &parser{data: data, dec: json.NewDecoder(bytes.NewReader(data))}
	p.dec.UseNumber()
	for i, b := range data {
		if b == '\n' {
			p.newlines = append(p.newlines, i)
		}
	}
	root, err := p.value()
	if err != nil {
		return nil, err
	}
	at := p.next()
	if _, err := p.dec.Token(); !errors.Is(err, io.EOF) {
		return nil, &ParseError{Line: p.line(at), Msg: "there is more after the document's closing brace; remove it"}
	}
	return root, nil
}

// parser walks a json.Decoder's tokens, mapping their offsets to lines.
type parser struct {
	data     []byte
	dec      *json.Decoder
	newlines []int // offsets of each '\n'
}

// next is the offset where the next token starts: past the white space and
// the separators the decoder has yet to consume.
func (p *parser) next() int {
	i := int(p.dec.InputOffset())
	for i < len(p.data) && strings.IndexByte(" \t\r\n,:", p.data[i]) >= 0 {
		i++
	}
	return i
}

// line is the 1-based line of the byte at offset.
func (p *parser) line(offset int) int {
	return sort.SearchInts(p.newlines, offset) + 1
}

// token reads the next token and the line it starts on. A token that does
// not parse is reported at the line it starts on: the decoder's offsets
// count from the start of the value it was decoding, or stop before the
// character it refused, so they cannot place the error.
func (p *parser) token() (json.Token, int, error) {
	line := p.line(p.next())
	tok, err := p.dec.Token()
	switch {
	case err == nil:
		return tok, line, nil
	case errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF):
		end := len(bytes.TrimRight(p.data, " \t\r\n")) - 1
		return nil, 0, &ParseError{Line: p.line(max(end, 0)), Msg: "the document ends before its value does; close every brace and bracket"}
	default:
		return nil, 0, &ParseError{Line: line, Msg: err.Error()}
	}
}

func (p *parser) value() (*Node, error) {
	tok, line, err := p.token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		if t == '{' {
			return p.object(line)
		}
		if t == '[' {
			return p.array(line)
		}
		return nil, &ParseError{Line: line, Msg: fmt.Sprintf("unexpected %q", rune(t))}
	case string:
		return &Node{Kind: KindString, Line: line, Text: t}, nil
	case json.Number:
		f, err := t.Float64()
		if err != nil {
			return nil, &ParseError{Line: line, Msg: fmt.Sprintf("the number %s is out of range; write a smaller one", t)}
		}
		return &Node{Kind: KindNumber, Line: line, Number: f}, nil
	case bool:
		return &Node{Kind: KindBool, Line: line, Bool: t}, nil
	default:
		return &Node{Kind: KindNull, Line: line}, nil
	}
}

func (p *parser) object(line int) (*Node, error) {
	n := &Node{Kind: KindObject, Line: line}
	first := map[string]int{}
	for p.dec.More() {
		tok, keyLine, err := p.token()
		if err != nil {
			return nil, err
		}
		key, _ := tok.(string) // the decoder refuses an object key that is not a string
		if at, ok := first[key]; ok {
			return nil, &ParseError{Line: keyLine, Msg: fmt.Sprintf("the key %q is given twice in one object, first at line %d; keep one", key, at)}
		}
		first[key] = keyLine
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		n.Members = append(n.Members, Member{Key: key, Line: keyLine, Value: v})
	}
	if _, _, err := p.token(); err != nil {
		return nil, err
	}
	return n, nil
}

func (p *parser) array(line int) (*Node, error) {
	n := &Node{Kind: KindArray, Line: line}
	for p.dec.More() {
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		n.Items = append(n.Items, v)
	}
	if _, _, err := p.token(); err != nil {
		return nil, err
	}
	return n, nil
}

// File is a parsed token file.
type File struct {
	Root     *Node     // the whole document, an object
	Resolver bool      // a DTCG resolver document: it has a top-level resolutionOrder
	Tokens   *TokenSet // a plain token file's tokens; nil for a resolver document
}

// Parse reads a token file: a plain DTCG token file, whose tokens it reads,
// or a resolver document, which it leaves to be resolved. A document that
// does not parse, or is not an object, is a *ParseError.
func Parse(data []byte) (*File, error) {
	root, err := ParseJSON(data)
	if err != nil {
		return nil, err
	}
	if root.Kind != KindObject {
		return nil, &ParseError{Line: root.Line, Msg: fmt.Sprintf("the document is %s, not an object of token groups such as {\"color\": {...}}", root.Kind)}
	}
	f := &File{Root: root, Resolver: root.Get("resolutionOrder") != nil}
	if !f.Resolver {
		f.Tokens = ReadTokens(root)
	}
	return f, nil
}

// Read parses the token file name in fsys. A file that is not there is an
// error that wraps fs.ErrNotExist; one that does not parse wraps a
// *ParseError.
func Read(fsys fs.FS, name string) (*File, error) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, fmt.Errorf("read the token file %s: %w", name, err)
	}
	f, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("the token file %s does not parse, fix it at %w", name, err)
	}
	return f, nil
}

// Problem is something wrong with a token file: what, at which token path
// and line, and what to do.
type Problem struct {
	Path    string `json:"path,omitempty"` // dotted token path, such as color.primary; "" for the whole file
	Line    int    `json:"line"`           // 1-based line
	Message string `json:"message"`
}

// Token is one design token: an object with a $value in a token tree.
type Token struct {
	Path  string // dotted, from the top-level group down: color.primary, spacing
	Group string // the first part of Path
	Name  string // the last part of Path
	Type  string // its $type, else its nearest group's; "" when none has one
	Line  int    // 1-based line of the token's key
	Value *Node  // its $value
	Ref   string // the token path an alias $value such as "{color.primary}" names; "" when the value is no alias
}

// TokenSet is a token tree read into its tokens.
type TokenSet struct {
	Tokens   []Token        // in file order
	Groups   map[string]int // each top-level group and the line of its key
	Problems []Problem      // members that are neither token nor group, and properties flai does not read
}

// Get answers the token at path.
func (s *TokenSet) Get(path string) (Token, bool) {
	for _, t := range s.Tokens {
		if t.Path == path {
			return t, true
		}
	}
	return Token{}, false
}

// The $type values flai reads.
const (
	TypeColor      = "color"
	TypeFontFamily = "fontFamily"
	TypeDimension  = "dimension"
	TypeShadow     = "shadow"
	TypeTypography = "typography"
)

// Types are the $type values flai reads.
var Types = []string{TypeColor, TypeFontFamily, TypeDimension, TypeShadow, TypeTypography}

// The token groups flai reads, named after Tailwind v4's theme namespaces.
const (
	GroupColor   = "color"
	GroupFont    = "font"
	GroupText    = "text"
	GroupSpacing = "spacing" // one token, the group key itself
	GroupRadius  = "radius"
	GroupShadow  = "shadow"
)

// Groups are the token groups flai reads, in the order they are rendered.
var Groups = []string{GroupColor, GroupFont, GroupText, GroupSpacing, GroupRadius, GroupShadow}

// groupTypes are the types each group's tokens may have.
var groupTypes = map[string][]string{
	GroupColor:   {TypeColor},
	GroupFont:    {TypeFontFamily},
	GroupText:    {TypeTypography, TypeDimension},
	GroupSpacing: {TypeDimension},
	GroupRadius:  {TypeDimension},
	GroupShadow:  {TypeShadow},
}

// ColorRoles are the required colour roles (ADR-0148); any other colour is an
// extra token.
var ColorRoles = []string{
	"background", "surface", "text", "text-muted", "border",
	"primary", "on-primary", "secondary", "on-secondary", "focus",
	"success", "success-soft", "warning", "warning-soft",
	"danger", "danger-soft", "info", "info-soft",
}

// FontNames are the required font stacks.
var FontNames = []string{"sans", "mono"}

// TextSteps are the required steps of the type scale.
var TextSteps = []string{"xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl"}

// RadiusNames are the required radii.
var RadiusNames = []string{"sm", "md", "lg"}

// ShadowNames are the required shadows.
var ShadowNames = []string{"sm", "md", "lg"}

// IsRole reports whether name is a required colour role.
func IsRole(name string) bool { return slices.Contains(ColorRoles, name) }

// RequiredPaths lists the path of every required token: the colour roles,
// fonts, type steps, spacing, radii, and shadows.
func RequiredPaths() []string {
	var out []string
	for _, g := range []struct {
		group string
		names []string
	}{{GroupColor, ColorRoles}, {GroupFont, FontNames}, {GroupText, TextSteps}} {
		for _, n := range g.names {
			out = append(out, g.group+"."+n)
		}
	}
	out = append(out, GroupSpacing)
	for _, n := range RadiusNames {
		out = append(out, GroupRadius+"."+n)
	}
	for _, n := range ShadowNames {
		out = append(out, GroupShadow+"."+n)
	}
	return out
}

var aliasRe = regexp.MustCompile(`^\{([^{}.]+(?:\.[^{}.]+)*)\}$`)

// Alias answers the token path a DTCG alias such as "{color.primary}" names.
func Alias(s string) (string, bool) {
	m := aliasRe.FindStringSubmatch(s)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// ReadTokens reads a token tree, an object of groups and tokens, into its
// tokens. A token is an object with $value; a group is an object without
// one, and its $type is its tokens' unless they give their own.
func ReadTokens(root *Node) *TokenSet {
	s := &TokenSet{Groups: map[string]int{}}
	s.group(root, nil, "")
	return s
}

// group reads one group's members; path is the group's, nil for the root.
func (s *TokenSet) group(n *Node, path []string, typ string) {
	if t := n.Get("$type"); t != nil && t.Kind == KindString {
		typ = t.Text
	}
	where := strings.Join(path, ".")
	for _, m := range n.Members {
		if strings.HasPrefix(m.Key, "$") {
			switch {
			case m.Key == "$type" && m.Value.Kind != KindString:
				s.problem(where, m.Line, fmt.Sprintf("$type is %s; write the type as a string such as \"color\"", m.Value.Kind))
			case m.Key == "$type", m.Key == "$description", m.Key == "$extensions", m.Key == "$deprecated":
			case m.Key == "$schema" && path == nil:
			default:
				s.problem(where, m.Line, fmt.Sprintf("flai does not read the group property %s; remove it", m.Key))
			}
			continue
		}
		p := append(slices.Clip(path), m.Key)
		full := strings.Join(p, ".")
		if strings.ContainsAny(m.Key, ".{}") {
			s.problem(full, m.Line, fmt.Sprintf("the name %q holds '.', '{', or '}', which DTCG names may not; rename it", m.Key))
			continue
		}
		if m.Value.Kind != KindObject {
			s.problem(full, m.Line, fmt.Sprintf("this is %s, neither a token nor a group: a token is an object with $value, a group an object of tokens", m.Value.Kind))
			continue
		}
		if m.Value.Get("$value") != nil {
			s.token(m, p, typ)
			continue
		}
		if path == nil {
			s.Groups[m.Key] = m.Line
		}
		s.group(m.Value, p, typ)
	}
}

// token reads one token; typ is its group's type.
func (s *TokenSet) token(m Member, path []string, typ string) {
	t := Token{
		Path: strings.Join(path, "."), Group: path[0], Name: path[len(path)-1],
		Type: typ, Line: m.Line, Value: m.Value.Get("$value"),
	}
	for _, pm := range m.Value.Members {
		switch {
		case pm.Key == "$type" && pm.Value.Kind == KindString:
			t.Type = pm.Value.Text
		case pm.Key == "$type":
			s.problem(t.Path, pm.Line, fmt.Sprintf("$type is %s; write the type as a string such as \"color\"", pm.Value.Kind))
		case pm.Key == "$value", pm.Key == "$description", pm.Key == "$extensions", pm.Key == "$deprecated":
		case strings.HasPrefix(pm.Key, "$"):
			s.problem(t.Path, pm.Line, fmt.Sprintf("flai does not read the token property %s; remove it", pm.Key))
		default:
			s.problem(t.Path+"."+pm.Key, pm.Line, "a token holds no tokens; move this out of the token's object")
		}
	}
	if t.Value.Kind == KindString {
		if ref, ok := Alias(t.Value.Text); ok {
			t.Ref = ref
		}
	}
	s.Tokens = append(s.Tokens, t)
}

func (s *TokenSet) problem(path string, line int, msg string) {
	s.Problems = append(s.Problems, Problem{Path: path, Line: line, Message: msg})
}

// Validate reports what is wrong with a token set against flai's subset of
// DTCG (ADR-0148): what reading found, then each token's placement, type,
// and value in file order, then each required token that is missing, at
// the line of its group, or line 1 when the group is missing. An alias is
// left for the resolver: its target is not looked up here.
func Validate(s *TokenSet) []Problem {
	out := slices.Clone(s.Problems)
	for _, t := range s.Tokens {
		if msg := check(t); msg != "" {
			out = append(out, Problem{Path: t.Path, Line: t.Line, Message: msg})
		}
	}
	for _, path := range RequiredPaths() {
		if _, ok := s.Get(path); ok {
			continue
		}
		group, _, _ := strings.Cut(path, ".")
		line, ok := s.Groups[group]
		if !ok {
			line = 1
		}
		out = append(out, Problem{Path: path, Line: line, Message: fmt.Sprintf("%s is required; add it", path)})
	}
	return out
}

// check says what is wrong with one token, or "".
func check(t Token) string {
	if msg := placement(t); msg != "" {
		return msg
	}
	switch {
	case t.Type == "" && t.Ref != "":
		return "" // its type is its target's, which the resolver finds
	case t.Type == "":
		return fmt.Sprintf("it has no $type; give it, or its group, one of %s", strings.Join(Types, ", "))
	case !slices.Contains(Types, t.Type):
		return fmt.Sprintf("flai does not read the $type %q; use one of %s", t.Type, strings.Join(Types, ", "))
	case !slices.Contains(groupTypes[t.Group], t.Type):
		return fmt.Sprintf("it is a %s token, and the %s group holds %s tokens; change its $type or move it", t.Type, t.Group, strings.Join(groupTypes[t.Group], " or "))
	case t.Ref != "":
		return ""
	case t.Value.Kind == KindString && strings.ContainsAny(t.Value.Text, "{}"):
		return fmt.Sprintf("the alias %q is malformed; write the path of the token it names in braces, such as {color.primary}", t.Value.Text)
	}
	if err := checkValue(t); err != nil {
		return err.Error()
	}
	return ""
}

// placement says what is wrong with where a token sits, or "".
func placement(t Token) string {
	depth := strings.Count(t.Path, ".") + 1
	switch {
	case !slices.Contains(Groups, t.Group):
		return fmt.Sprintf("flai reads tokens only in the groups %s; move it into one or remove it", strings.Join(Groups, ", "))
	case t.Group == GroupSpacing && depth != 1:
		return "spacing is one dimension token, the group key itself: give spacing its own $value and remove this"
	case t.Group != GroupSpacing && depth == 1:
		return fmt.Sprintf("%s is a group of named tokens; put this value under a name, such as %s.md", t.Group, t.Group)
	case depth > 2:
		return fmt.Sprintf("flai reads one level of names under %s; name it %s.%s", t.Group, t.Group, strings.Join(strings.Split(t.Path, ".")[1:], "-"))
	}
	return ""
}

// checkValue decodes a token's value by its type.
func checkValue(t Token) error {
	switch t.Type {
	case TypeColor:
		c, err := DecodeColor(t.Value)
		if err != nil {
			return err
		}
		if c.Alpha != 1 && IsRole(t.Name) {
			return fmt.Errorf("the colour role %s has alpha %g; a role is opaque, alpha is only for shadows and extra tokens: remove alpha", t.Name, c.Alpha)
		}
		return nil
	case TypeFontFamily:
		_, err := DecodeFontFamily(t.Value)
		return err
	case TypeDimension:
		_, err := DecodeDimension(t.Value)
		return err
	case TypeTypography:
		_, err := DecodeTypography(t.Value)
		return err
	default:
		_, err := DecodeShadow(t.Value)
		return err
	}
}

// The colour spaces flai reads.
const (
	ColorSpaceSRGB  = "srgb"
	ColorSpaceOKLCH = "oklch"
)

// Color is a colour in a colour space flai reads.
type Color struct {
	Space      string     // ColorSpaceSRGB or ColorSpaceOKLCH
	Components [3]float64 // srgb: red, green, blue, 0 to 1; oklch: lightness 0 to 1, chroma, hue in degrees
	NoHue      bool       // oklch: the hue is "none", an achromatic colour, and Components[2] is 0
	Alpha      float64    // 0 to 1; 1 when the file gives none
	Hex        string     // the six-digit #rrggbb the file gives, lower case; "" when it gives none
}

var hexRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// DecodeColor reads a DTCG colour value: srgb with its hex, or oklch.
func DecodeColor(n *Node) (Color, error) {
	if n.Kind != KindObject {
		return Color{}, fmt.Errorf("a colour is an object such as {\"colorSpace\": \"srgb\", \"components\": [1, 1, 1], \"hex\": \"#ffffff\"}, not %s", n.Kind)
	}
	c := Color{Alpha: 1}
	space := n.Get("colorSpace")
	if space == nil || space.Kind != KindString {
		return Color{}, errors.New("a colour needs its colorSpace, srgb or oklch")
	}
	c.Space = space.Text
	if c.Space != ColorSpaceSRGB && c.Space != ColorSpaceOKLCH {
		return Color{}, fmt.Errorf("flai does not read the colour space %q; write the colour in srgb, with its hex, or in oklch", c.Space)
	}
	comps := n.Get("components")
	if comps == nil || comps.Kind != KindArray || len(comps.Items) != 3 {
		return Color{}, fmt.Errorf("an %s colour needs components, an array of three numbers", c.Space)
	}
	if c.Space == ColorSpaceOKLCH && comps.Items[2].Kind == KindString && comps.Items[2].Text == "none" {
		c.NoHue = true
	}
	for i, item := range comps.Items {
		if i == 2 && c.NoHue {
			continue
		}
		if item.Kind != KindNumber {
			return Color{}, fmt.Errorf("component %d is %s; write a number", i+1, item.Kind)
		}
		c.Components[i] = item.Number
	}
	if err := c.inRange(); err != nil {
		return Color{}, err
	}
	if a := n.Get("alpha"); a != nil {
		if a.Kind != KindNumber || a.Number < 0 || a.Number > 1 {
			return Color{}, errors.New("alpha is a number from 0 to 1; correct it or remove it")
		}
		c.Alpha = a.Number
	}
	hex := n.Get("hex")
	switch {
	case hex == nil && c.Space == ColorSpaceSRGB:
		return Color{}, fmt.Errorf("an srgb colour needs its hex, the six-digit %s", c.hexOf())
	case hex == nil:
		return c, nil
	case hex.Kind != KindString || !hexRe.MatchString(hex.Text):
		return Color{}, errors.New("hex is a six-digit #rrggbb such as #1a2b3c; correct it")
	}
	c.Hex = strings.ToLower(hex.Text)
	if c.Space == ColorSpaceSRGB && !c.hexAgrees() {
		return Color{}, fmt.Errorf("its hex %s and its components, which give %s, are not the same colour; make them agree", c.Hex, c.hexOf())
	}
	return c, nil
}

// inRange says which component is outside its colour space's range.
func (c Color) inRange() error {
	if c.Space == ColorSpaceSRGB {
		for i, v := range c.Components {
			if v < 0 || v > 1 {
				return fmt.Errorf("srgb component %d is %g; red, green, and blue are 0 to 1", i+1, v)
			}
		}
		return nil
	}
	l, ch, h := c.Components[0], c.Components[1], c.Components[2]
	switch {
	case l < 0 || l > 1:
		return fmt.Errorf("oklch lightness is %g; it is 0 to 1", l)
	case ch < 0:
		return fmt.Errorf("oklch chroma is %g; it is 0 or more", ch)
	case h < 0 || h > 360:
		return fmt.Errorf("oklch hue is %g; it is 0 to 360 degrees, or \"none\"", h)
	}
	return nil
}

// hexOf is the #rrggbb an srgb colour's components round to.
func (c Color) hexOf() string {
	r, g, b := c.Components[0], c.Components[1], c.Components[2]
	return fmt.Sprintf("#%02x%02x%02x", int(math.Round(r*255)), int(math.Round(g*255)), int(math.Round(b*255)))
}

// hexAgrees reports whether hex and components are one colour, to within
// rounding: each channel within 1/255.
func (c Color) hexAgrees() bool {
	for i := range 3 {
		b, _ := strconv.ParseUint(c.Hex[1+2*i:3+2*i], 16, 8) // hexRe has matched
		if math.Abs(c.Components[i]-float64(b)/255) > 1.0/255+1e-9 {
			return false
		}
	}
	return true
}

// The units of a dimension flai reads.
const (
	UnitRem = "rem"
	UnitPx  = "px"
)

// Dimension is a length.
type Dimension struct {
	Value float64
	Unit  string // UnitRem or UnitPx
}

// DecodeDimension reads a DTCG dimension, {"value": n, "unit": "rem"} or px.
func DecodeDimension(n *Node) (Dimension, error) {
	if n.Kind != KindObject {
		return Dimension{}, fmt.Errorf("a dimension is an object such as {\"value\": 0.25, \"unit\": \"rem\"}, not %s", n.Kind)
	}
	v, u := n.Get("value"), n.Get("unit")
	if v == nil || v.Kind != KindNumber {
		return Dimension{}, errors.New("a dimension needs its value, a number")
	}
	if u == nil || u.Kind != KindString || (u.Text != UnitRem && u.Text != UnitPx) {
		return Dimension{}, errors.New("a dimension's unit is rem or px; flai reads no other")
	}
	return Dimension{Value: v.Number, Unit: u.Text}, nil
}

// DecodeFontFamily reads a DTCG font family: a name, or a stack of names.
func DecodeFontFamily(n *Node) ([]string, error) {
	const want = "a font family is a name, or a non-empty array of names such as [\"Inter\", \"sans-serif\"]"
	if n.Kind == KindString && n.Text != "" {
		return []string{n.Text}, nil
	}
	if n.Kind != KindArray || len(n.Items) == 0 {
		return nil, errors.New(want)
	}
	out := make([]string, 0, len(n.Items))
	for _, item := range n.Items {
		if item.Kind != KindString || item.Text == "" {
			return nil, errors.New(want)
		}
		out = append(out, item.Text)
	}
	return out, nil
}

// Typography is a step of the type scale.
type Typography struct {
	FontSize      Dimension
	LineHeight    LineHeight
	FontWeight    float64    // 1 to 1000; 0 when not given
	FontFamily    []string   // nil when not given
	LetterSpacing *Dimension // nil when not given
}

// LineHeight is a line height: a multiple of the font size, or a length.
type LineHeight struct {
	Number float64    // the multiple, when Length is nil
	Length *Dimension // the length, when given as one
}

// fontWeights are DTCG's named font weights.
var fontWeights = map[string]float64{
	"thin": 100, "hairline": 100, "extra-light": 200, "ultra-light": 200, "light": 300,
	"normal": 400, "regular": 400, "book": 400, "medium": 500,
	"semi-bold": 600, "demi-bold": 600, "bold": 700, "extra-bold": 800, "ultra-bold": 800,
	"black": 900, "heavy": 900, "extra-black": 950, "ultra-black": 950,
}

// DecodeTypography reads a DTCG typography value: fontSize and lineHeight,
// and optionally fontWeight, fontFamily, and letterSpacing.
func DecodeTypography(n *Node) (Typography, error) {
	if n.Kind != KindObject {
		return Typography{}, fmt.Errorf("a typography value is an object with fontSize and lineHeight, not %s", n.Kind)
	}
	var t Typography
	size := n.Get("fontSize")
	if size == nil {
		return Typography{}, errors.New("a typography value needs its fontSize, a dimension")
	}
	d, err := DecodeDimension(size)
	if err != nil {
		return Typography{}, fmt.Errorf("fontSize: %w", err)
	}
	t.FontSize = d
	lh := n.Get("lineHeight")
	switch {
	case lh == nil:
		return Typography{}, errors.New("a typography value needs its lineHeight, a number or a dimension")
	case lh.Kind == KindNumber:
		t.LineHeight.Number = lh.Number
	default:
		d, err := DecodeDimension(lh)
		if err != nil {
			return Typography{}, fmt.Errorf("lineHeight is a number or a dimension: %w", err)
		}
		t.LineHeight.Length = &d
	}
	if w := n.Get("fontWeight"); w != nil {
		named, ok := fontWeights[w.Text]
		switch {
		case w.Kind == KindNumber && w.Number >= 1 && w.Number <= 1000:
			t.FontWeight = w.Number
		case w.Kind == KindString && ok:
			t.FontWeight = named
		default:
			return Typography{}, errors.New("fontWeight is a number from 1 to 1000 or a DTCG weight name such as bold")
		}
	}
	if f := n.Get("fontFamily"); f != nil {
		if t.FontFamily, err = DecodeFontFamily(f); err != nil {
			return Typography{}, fmt.Errorf("fontFamily: %w", err)
		}
	}
	if ls := n.Get("letterSpacing"); ls != nil {
		d, err := DecodeDimension(ls)
		if err != nil {
			return Typography{}, fmt.Errorf("letterSpacing: %w", err)
		}
		t.LetterSpacing = &d
	}
	return t, nil
}

// Shadow is one layer of a shadow.
type Shadow struct {
	Color                          Color // alpha allowed
	OffsetX, OffsetY, Blur, Spread Dimension
	Inset                          bool
}

// DecodeShadow reads a DTCG shadow value: one layer, or an array of layers.
func DecodeShadow(n *Node) ([]Shadow, error) {
	layers := []*Node{n}
	if n.Kind == KindArray {
		if len(n.Items) == 0 {
			return nil, errors.New("a shadow is a layer or a non-empty array of layers")
		}
		layers = n.Items
	}
	out := make([]Shadow, 0, len(layers))
	for i, l := range layers {
		s, err := decodeLayer(l)
		if err != nil {
			if len(layers) > 1 {
				return nil, fmt.Errorf("layer %d: %w", i+1, err)
			}
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func decodeLayer(n *Node) (Shadow, error) {
	if n.Kind != KindObject {
		return Shadow{}, fmt.Errorf("a shadow layer is an object with color, offsetX, offsetY, blur, and spread, not %s", n.Kind)
	}
	var s Shadow
	c := n.Get("color")
	if c == nil {
		return Shadow{}, errors.New("a shadow layer needs its color")
	}
	var err error
	if s.Color, err = DecodeColor(c); err != nil {
		return Shadow{}, fmt.Errorf("color: %w", err)
	}
	for _, f := range []struct {
		key string
		to  *Dimension
	}{{"offsetX", &s.OffsetX}, {"offsetY", &s.OffsetY}, {"blur", &s.Blur}, {"spread", &s.Spread}} {
		v := n.Get(f.key)
		if v == nil {
			return Shadow{}, fmt.Errorf("a shadow layer needs its %s, a dimension", f.key)
		}
		if *f.to, err = DecodeDimension(v); err != nil {
			return Shadow{}, fmt.Errorf("%s: %w", f.key, err)
		}
	}
	if in := n.Get("inset"); in != nil {
		if in.Kind != KindBool {
			return Shadow{}, errors.New("inset is true or false")
		}
		s.Inset = in.Bool
	}
	return s, nil
}

// Headings lists the headings of a markdown document, such as the guidance
// design/system/design-system.md, in order: its ATX headings outside its
// front matter and fenced code, without their #s.
func Headings(markdown []byte) ([]string, error) {
	doc, err := topics.Parse(string(markdown), nil)
	if err != nil {
		return nil, fmt.Errorf("read the guidance's headings: %w; correct its front matter", err)
	}
	var out []string
	for _, s := range doc.Sections {
		if s.Level > 0 {
			out = append(out, s.Heading)
		}
	}
	return out, nil
}

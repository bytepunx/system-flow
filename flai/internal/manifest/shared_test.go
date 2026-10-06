package manifest

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// S-0295 (ADR-0096): a plain pattern is itself and everything below it, *
// and ? stay within a segment, ** is whole segments; a folder entry lies
// inside a pattern only when the whole folder does.
func TestSharedMatchAndCovers(t *testing.T) {
	tests := []struct {
		pattern, entry string
		match, covers  bool
	}{
		{"design/adrs", "design/adrs/0099-x.md", true, true},
		{"design/adrs/**", "design/adrs/0099-x.md", true, true},
		{"design/adrs", "design/adrs", true, true},
		{"design/adrs/**", "design/adrs", true, true},
		{"design/adrs", "design/adrs-old/x.md", false, false},
		{"design/adrs", "design", false, false},
		{"docs/users/*.md", "docs/users/flai.md", true, true},
		{"docs/users/*.md", "docs/users/a/b.md", false, false},
		{"docs/users/*.md", "docs/users", false, false},
		{"docs/users/*.md", "docs/users/notes", false, false},
		{"docs", "docs/users", true, true},
		{"docs/users/**", "docs/users", true, true},
		{"docs/users/**", "docs", false, false},
		{"docs/users/flai.md", "docs/users/flai.md", true, true},
		{"docs/users/flai.md", "docs/users", false, false},
		{"docs/*", "docs/users", true, false},
		{"design/**/x.md", "design/x.md", true, true},
		{"design/**/x.md", "design/a/b/x.md", true, true},
		{"**/CHANGELOG.md", "template/CHANGELOG.md", true, true},
		{"flai/internal/*/testdata/**", "flai/internal/check/testdata", true, true},
		{"flai/internal/*/testdata/**", "flai/internal/check", false, false},
		{"design/issue?", "design/issues/I-0001.md", false, false},
		{"design/issue?/**", "design/issues/I-0001.md", true, true},
		{"Makefile", "Makefile", true, true},
		{"Make*", "Makefile", true, false},
		{"design/adrs", "./design/adrs/", true, true},
		{"design/adrs", "/design/adrs", false, false},
		{"design/adrs", "design/../design/adrs", false, false},
		{"design/adrs", "", false, false},
	}
	for _, tt := range tests {
		c := Claims{Shared: []string{"unrelated", tt.pattern}}
		p, ok := c.Match(tt.entry)
		if ok != tt.match || ok && p != tt.pattern {
			t.Errorf("%q Match(%q) = %q, %v; want %v", tt.pattern, tt.entry, p, ok, tt.match)
		}
		p, ok = c.Covers(tt.entry)
		if ok != tt.covers || ok && p != tt.pattern {
			t.Errorf("%q Covers(%q) = %q, %v; want %v", tt.pattern, tt.entry, p, ok, tt.covers)
		}
	}
}

// The first pattern that matches is named, and a pattern that is not valid
// frees nothing.
func TestSharedMatchNamesFirstAndSkipsInvalid(t *testing.T) {
	c := Claims{Shared: []string{"/design", "design/[", "design/**", "design/adrs"}}
	if p, ok := c.Covers("design/adrs/x.md"); !ok || p != "design/**" {
		t.Errorf("Covers = %q, %v; want design/**", p, ok)
	}
	bad := Claims{Shared: []string{"/design", "design/../x", ""}}
	if p, ok := bad.Match("design/x"); ok {
		t.Errorf("an invalid pattern matched: %q", p)
	}
}

func TestSharedErrors(t *testing.T) {
	tests := []struct {
		pattern string
		reason  string // a word the reason holds; "" when valid
	}{
		{"design/adrs", ""},
		{"docs/users/*.md", ""},
		{"**/x.md", ""},
		{"design/[a-c]*", ""},
		{"", "empty"},
		{"  ", "empty"},
		{"/design/adrs", "absolute"},
		{`C:/design`, "absolute"},
		{"design/../x", ".."},
		{"..", ".."},
		{"./design", ". segment"},
		{"design//adrs", "empty segment"},
		{"design/adrs/", "empty segment"},
		{"design/a**", "inside a segment"},
		{"design/[a-", "not a valid glob"},
		{`design/a\`, "not a valid glob"},
	}
	for _, tt := range tests {
		errs := Claims{Shared: []string{"ok", tt.pattern}}.Errors()
		if tt.reason == "" {
			if len(errs) != 0 {
				t.Errorf("%q: unexpected %v", tt.pattern, errs)
			}
			continue
		}
		if len(errs) != 1 || errs[0].Index != 1 || errs[0].Pattern != tt.pattern || !strings.Contains(errs[0].Reason, tt.reason) {
			t.Errorf("%q: errors %v; want one at index 1 saying %q", tt.pattern, errs, tt.reason)
			continue
		}
		if msg := errs[0].Error(); !strings.HasPrefix(msg, "claims.shared pattern ") || !strings.Contains(msg, "; ") {
			t.Errorf("%q: message %q names neither the key nor what to do", tt.pattern, msg)
		}
	}
}

// A bad pattern does not stop the manifest loading.
func TestLoadKeepsInvalidSharedPatterns(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	body := "version: 1\nname: demo\nlayout:\n  design: d\n  docs: docs\n  wip: wip\nclaims:\n  shared:\n    - d/adrs\n    - /abs\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(m.Claims.Shared, []string{"d/adrs", "/abs"}) || len(m.Claims.Errors()) != 1 {
		t.Fatalf("claims %+v, errors %v", m.Claims, m.Claims.Errors())
	}
	if lines := SharedLines([]byte(body)); !reflect.DeepEqual(lines, []int{9, 10}) {
		t.Errorf("SharedLines = %v; want [9 10]", lines)
	}
}

const sharedFixture = `# the manifest
version: 1
name: demo   # the project's name
layout:
  design: design
  docs: docs
  wip: wip
planning:
  hour_rate: 95 # per hour
claims:            # how touches claim paths
  # shared paths hold no story
  shared:          # many stories change these
    - design/adrs  # a new file each
    # the guide
    - "docs/users/*.md"
  other: kept
dashboard:
  port: 4242
`

func TestAddAndRemoveSharedKeepTheRest(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		add     []string
		remove  []string
		want    string
		shared  []string
		comment string
	}{
		{
			name:   "add after the last item",
			in:     sharedFixture,
			add:    []string{"design/issues", "**/CHANGELOG.md"},
			want:   strings.Replace(sharedFixture, "    - \"docs/users/*.md\"\n", "    - \"docs/users/*.md\"\n    - design/issues\n    - \"**/CHANGELOG.md\"\n", 1),
			shared: []string{"design/adrs", "docs/users/*.md", "design/issues", "**/CHANGELOG.md"},
		},
		{
			name:   "remove one item and its line only",
			in:     sharedFixture,
			remove: []string{"design/adrs"},
			want:   strings.Replace(sharedFixture, "    - design/adrs  # a new file each\n", "", 1),
			shared: []string{"docs/users/*.md"},
		},
		{
			name:   "remove every item leaves shared: []",
			in:     sharedFixture,
			remove: []string{"docs/users/*.md", "design/adrs"},
			want: strings.Replace(strings.Replace(strings.Replace(sharedFixture,
				"    - design/adrs  # a new file each\n", "", 1),
				"    - \"docs/users/*.md\"\n", "", 1),
				"  shared:          # many stories change these\n", "  shared: [] # many stories change these\n", 1),
			shared: []string{},
		},
		{
			name:   "add to shared: []",
			in:     "version: 1\nclaims:\n  shared: [] # none yet\nflai:\n  minimum: 1.0.0\n",
			add:    []string{"design/adrs"},
			want:   "version: 1\nclaims:\n  shared: # none yet\n    - design/adrs\nflai:\n  minimum: 1.0.0\n",
			shared: []string{"design/adrs"},
		},
		{
			name:   "a flow list is written again a line an item",
			in:     "version: 1\nclaims:\n  shared: [a, \"b # c\"]\n",
			add:    []string{"d"},
			want:   "version: 1\nclaims:\n  shared:\n    - a\n    - \"b # c\"\n    - d\n",
			shared: []string{"a", "b # c", "d"},
		},
		{
			name:   "items at the key's indent",
			in:     "claims:\n    shared:\n    - a\n    other: 1\n",
			add:    []string{"b"},
			want:   "claims:\n    shared:\n    - a\n    - b\n    other: 1\n",
			shared: []string{"a", "b"},
		},
		{
			name:   "claims without shared",
			in:     "version: 1\nclaims:   # mine\n    other: 1\n# next\nflai:\n  minimum: 1.0.0\n",
			add:    []string{"design/adrs"},
			want:   "version: 1\nclaims:   # mine\n    other: 1\n    shared:\n      - design/adrs\n# next\nflai:\n  minimum: 1.0.0\n",
			shared: []string{"design/adrs"},
		},
		{
			name:   "no claims at all",
			in:     "version: 1\nname: x # n\n",
			add:    []string{"design/adrs", "design/issues"},
			want:   "version: 1\nname: x # n\nclaims:\n  shared:\n    - design/adrs\n    - design/issues\n",
			shared: []string{"design/adrs", "design/issues"},
		},
		{
			name:   "no claims and no final newline",
			in:     "version: 1",
			add:    []string{"x"},
			want:   "version: 1\nclaims:\n  shared:\n    - x\n",
			shared: []string{"x"},
		},
		{
			name:   "claims: {}",
			in:     "version: 1\nclaims: {} # empty\n",
			add:    []string{"x"},
			want:   "version: 1\nclaims: # empty\n  shared:\n    - x\n",
			shared: []string{"x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out []byte
			var ch Change
			var err error
			if tt.add != nil {
				out, ch, err = AddSharedBytes([]byte(tt.in), tt.add...)
			} else {
				out, ch, err = RemoveSharedBytes([]byte(tt.in), tt.remove...)
			}
			if err != nil {
				t.Fatal(err)
			}
			if string(out) != tt.want {
				t.Errorf("got\n%s\nwant\n%s", out, tt.want)
			}
			want := Change{Added: tt.add, Removed: tt.remove, Shared: tt.shared}
			if !reflect.DeepEqual(ch, want) {
				t.Errorf("change %+v; want %+v", ch, want)
			}
		})
	}
}

func TestSharedEditRefusals(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		add    []string
		remove []string
		says   string
	}{
		{"duplicate", sharedFixture, []string{"design/adrs"}, nil, "in the list already"},
		{"duplicate within the call", sharedFixture, []string{"x", "x"}, nil, "in the list already"},
		{"invalid", sharedFixture, []string{"design/../x"}, nil, ".. segment"},
		{"empty", sharedFixture, []string{" "}, nil, "is empty"},
		{"nothing to add", sharedFixture, []string{}, nil, "at least one pattern"},
		{"not in the list", sharedFixture, nil, []string{"design/issues"}, "holds design/adrs, docs/users/*.md"},
		{"not in an empty list", "version: 1\n", nil, []string{"x"}, "is empty"},
		{"claims inline", "claims: {shared: [a]}\n", nil, []string{"b"}, "not in the list"},
		{"claims inline add", "claims: {other: 1}\n", []string{"b"}, nil, "written on one line"},
		{"not yaml", "claims: [\n", []string{"b"}, nil, "fix the manifest's YAML"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out []byte
			var err error
			if tt.remove == nil {
				out, _, err = AddSharedBytes([]byte(tt.in), tt.add...)
			} else {
				out, _, err = RemoveSharedBytes([]byte(tt.in), tt.remove...)
			}
			if err == nil || !strings.Contains(err.Error(), tt.says) || out != nil {
				t.Errorf("err %v, out %q; want an error saying %q", err, out, tt.says)
			}
		})
	}
}

// The file forms write the manifest back, refuse without writing, and name
// the file.
func TestAddAndRemoveSharedFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(p, []byte(sharedFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	ch, err := AddShared(p, "design/issues")
	if err != nil || !reflect.DeepEqual(ch.Added, []string{"design/issues"}) {
		t.Fatalf("add: %+v, %v", ch, err)
	}
	m, err := Load(p)
	if err != nil || !reflect.DeepEqual(m.Claims.Shared, []string{"design/adrs", "docs/users/*.md", "design/issues"}) {
		t.Fatalf("after add: %+v, %v", m.Claims, err)
	}
	if _, err := AddShared(p, "design/issues"); err == nil || !strings.Contains(err.Error(), p) {
		t.Errorf("duplicate add: %v; want an error naming %s", err, p)
	}
	if _, err := RemoveShared(p, "design/issues"); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(p)
	if err != nil || string(got) != sharedFixture {
		t.Errorf("add then remove did not give the file back:\n%s", got)
	}
}

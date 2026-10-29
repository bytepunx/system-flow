package adr

import "testing"

// I-0222: a quoted title keeps <, >, and & as themselves, and titleOf reads
// it back, escapes and all, as YAML reads it.
func TestYamlStringWritesNoHTMLEscapes(t *testing.T) {
	for _, c := range []struct{ title, want string }{
		{"Plain words", "Plain words"},
		{"Tokens: one per project", `"Tokens: one per project"`},
		{"Escaped as < and >", `"Escaped as < and >"`},
		{"Rock & roll", `"Rock & roll"`},
		{`Says "hi"`, `"Says \"hi\""`},
	} {
		got := yamlString(c.title)
		if got != c.want {
			t.Errorf("yamlString(%q) = %s, want %s", c.title, got, c.want)
		}
		if back := titleOf("id: ADR-0001\ntitle: " + got + "\nstatus: proposed"); back != c.title {
			t.Errorf("titleOf reads %s back as %q, want %q", got, back, c.title)
		}
	}
}

func TestTitleOfReadsAnEscapeAsTheCharacter(t *testing.T) {
	if got := titleOf("title: \"a \\u003c b \\u003e c \\u0026 d\"\n"); got != "a < b > c & d" {
		t.Errorf("titleOf = %q", got)
	}
	if got := titleOf("title: \"unclosed\n"); got != "unclosed" {
		t.Errorf("a quote YAML would refuse is trimmed as before: %q", got)
	}
}

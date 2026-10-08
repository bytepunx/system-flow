package mdlint

import "testing"

// I-0118: a bare www. literal reached wip/agents/orchestrator.md, which
// MD034 finds. QuoteBareURLs wraps each bare URL MD034 finds in a code
// span, so that a line flai takes from elsewhere is written lint clean, and
// leaves the rest of the line as it was.
func TestQuoteBareURLs(t *testing.T) {
	c, err := Parse([]byte("default: true\n"), false, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, in, want string }{
		{"https", "Read https://example.com/a?b=c. Then go on", "Read `https://example.com/a?b=c`. Then go on"},
		{"www", "Published at www.example.com, see TH-0361", "Published at `www.example.com`, see TH-0361"},
		{"www first", "www.example.com is up", "`www.example.com` is up"},
		{"email", "Mail alex@example.com now", "Mail `alex@example.com` now"},
		{"several", "http://a.example.org and www.b.example.org", "`http://a.example.org` and `www.b.example.org`"},
		{"backtick", "See https://example.com/a`b here", "See ``https://example.com/a`b`` here"},
		{"code span", "Read `https://example.com` first", "Read `https://example.com` first"},
		{"link", "Read [www.example.com](https://www.example.com) and <https://example.org>", "Read [www.example.com](https://www.example.com) and <https://example.org>"},
		{"none", "Ordered S-0346 and S-0347 by WSJF", "Ordered S-0346 and S-0347 by WSJF"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := QuoteBareURLs(tc.in)
			if got != tc.want {
				t.Errorf("QuoteBareURLs(%q) = %q, want %q", tc.in, got, tc.want)
			}
			doc := "# Orchestrator\n\n## Entries\n\n### 2026-10-08T05:31:00Z\n\n- Summary: " + got + "\n"
			if f := c.Lint(doc); len(f) > 0 {
				t.Errorf("%q lints with %v", got, f)
			}
			if f := c.Lint("# Orchestrator\n\n- Summary: " + tc.in + "\n"); (len(f) > 0) != (tc.in != tc.want) {
				t.Errorf("MD034 on %q: %v", tc.in, f)
			}
		})
	}
}

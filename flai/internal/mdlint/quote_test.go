package mdlint

import (
	"path/filepath"
	"strings"
	"testing"
)

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

// I-0242: the merge queue's send-back thread quoted go test's summary line
// with its tabs, the lint refused the thread (MD010), and the story was never
// sent back. Fenced quotes tool output so that the project's own markdownlint
// configuration passes it, whatever the output carries.
func TestFencedLintsCleanWithTheProjectsConfiguration(t *testing.T) {
	// the repository's own .markdownlint.yaml, a file read: the package sits
	// three folders below the monorepo's root
	c, err := Load(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("no markdownlint configuration at the monorepo's root")
	}
	summary := "FAIL\tgithub.com/bytepunx/system-flow/flai/internal/serve\t131.960s"
	out := strings.Join([]string{
		"--- FAIL: TestServe (0.01s)",
		"    serve_test.go:12: got\twant  ",
		"\x1b[31mFAIL\x1b[0m \x1b[1mbuild\x1b]0;title\x07 red\x1b[m",
		"step one\r\nstep two\r\nstep three\r",
		"trailing spaces   ",
		"trailing tab\t",
		"",
		"",
		"",
		"after three blank lines",
		"a run of four backticks: ````",
		"\x00nul\x7fdel\bbackspace",
		summary,
	}, "\n")
	got := Fenced(out)
	doc := "# Merge queue batch 20261031-110330 fails with S-0546 at smoke\n\n## Entries\n\nThe tier failed:\n\n" + got + "\n\nThe whole output is in the log folder.\n"
	if f := c.Lint(doc); len(f) > 0 {
		t.Errorf("Fenced output lints with %v:\n%s", f, doc)
	}
	if want := "FAIL    github.com/bytepunx/system-flow/flai/internal/serve    131.960s"; !strings.Contains(got, "\n"+want+"\n") {
		t.Errorf("Fenced(...) does not quote the summary line as %q:\n%s", want, got)
	}
	if !strings.HasPrefix(got, "`````text\n") || !strings.HasSuffix(got, "\n`````") {
		t.Errorf("Fenced(...) fence is not five backticks:\n%s", got)
	}
	if f := c.Lint("# Thread\n\n```text\n" + summary + "\n```\n"); len(f) == 0 {
		t.Error("the summary line fenced as it is lints clean; the test no longer shows I-0242's refusal")
	}
}

func TestFenced(t *testing.T) {
	c, err := Parse([]byte("default: true\n"), false, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, in, want string }{
		{"empty", "", "```text\n```"},
		{"blank", " \n\t\n\r\n", "```text\n```"},
		{"plain", "ok", "```text\nok\n```"},
		{"tabs", "a\tb", "```text\na    b\n```"},
		{"crlf", "a\r\nb\rc", "```text\na\nb\nc\n```"},
		{"trailing", "a  \nb\t", "```text\na\nb\n```"},
		{"blank runs", "\n\na\n\n\n\nb\n\n", "```text\na\n\nb\n```"},
		{"csi", "\x1b[1;31mred\x1b[0m done", "```text\nred done\n```"},
		{"osc bel", "\x1b]0;title\x07x", "```text\nx\n```"},
		{"osc st", "\x1b]8;;https://example.com\x1b\\link\x1b]8;;\x1b\\", "```text\nlink\n```"},
		{"two byte escape", "\x1b(Bx", "```text\nx\n```"},
		{"controls", "a\x00b\x07c\x7fd\u0085e", "```text\nabcde\n```"},
		{"unicode kept", "✓ passed — 5 tests", "```text\n✓ passed — 5 tests\n```"},
		{"three backticks", "see ``` here", "````text\nsee ``` here\n````"},
		{"one backtick", "a `b` c", "```text\na `b` c\n```"},
		{"fence line", "```\ncode\n```", "````text\n```\ncode\n```\n````"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Fenced(tc.in)
			if got != tc.want {
				t.Errorf("Fenced(%q) = %q, want %q", tc.in, got, tc.want)
			}
			if f := c.Lint("# Thread\n\n" + got + "\n"); len(f) > 0 {
				t.Errorf("Fenced(%q) lints with %v", tc.in, f)
			}
		})
	}
}

package search

import (
	"strings"
	"testing"
)

func corpus() []Doc {
	return []Doc{
		{Path: "wip/kanban/stories/S-0042-tunnel.md", Kind: "item", ItemID: "S-0042", Title: "Research tunnel options", Tags: "flairport cli", Headings: "Goal Notes", Body: "Which tunnels can reach the dashboard from outside the network.", Type: "story", Status: "done", Scope: "wip"},
		{Path: "design/system/channel.md", Kind: "doc", Title: "A channel to flai on the host", Headings: "Security Tunnel clients", Body: "The dashboard never connects to the host. A tunnel client dials out. " + strings.Repeat("Filler words that say little. ", 40), Scope: "design"},
		{Path: "design/system/board.md", Kind: "doc", Title: "The board", Headings: "Columns", Body: "Cards move between columns. Nothing here about the word in question, only tunnelling once.", Scope: "design"},
		{Path: "docs/users/index.md", Kind: "doc", Title: "Tunnel guide for users", Headings: "Tunnel", Body: "How to use a tunnel.", Scope: "docs"},
		{Path: "wip/kanban/stories/S-0007-other.md", Kind: "item", ItemID: "S-0007", Title: "Something else", Body: "Unrelated.", Type: "story", Status: "backlog", Scope: "wip"},
	}
}

func paths(hits []Hit) string {
	var out []string
	for _, h := range hits {
		out = append(out, h.Path)
	}
	return strings.Join(out, " ")
}

func TestExactPrefixAndFuzzyMatchesRankInThatOrderOfTrust(t *testing.T) {
	ix := Build(corpus())
	got := ix.Search("tunnel", false, 0)
	// title beats headings beats body; the docs scope is left out; "tunnels" and "tunnelling" match by prefix
	if paths(got) != "wip/kanban/stories/S-0042-tunnel.md design/system/channel.md design/system/board.md" {
		t.Errorf("tunnel: %s", paths(got))
	}
	if got[0].ItemID != "S-0042" || got[0].Status != "done" || got[0].Type != "story" || got[0].Score <= got[1].Score {
		t.Errorf("first hit: %+v", got[0])
	}
	with := ix.Search("tunnel", true, 0)
	if !strings.Contains(paths(with), "docs/users/index.md") || len(with) != 4 {
		t.Errorf("with docs: %s", paths(with))
	}
	if got := ix.Search("tunel", false, 0); len(got) == 0 || got[0].Path != "wip/kanban/stories/S-0042-tunnel.md" {
		t.Errorf("a typo within a fifth of the word: %s", paths(got))
	}
	if got := ix.Search("tun", false, 0); len(got) != 3 {
		t.Errorf("a prefix: %s", paths(got))
	}
	if got := ix.Search("xyzzy", true, 0); len(got) != 0 {
		t.Errorf("nothing matches: %s", paths(got))
	}
	if got := ix.Search("   ", true, 0); len(got) != 0 {
		t.Errorf("an empty query: %s", paths(got))
	}
}

func TestAnItemIDFindsItsItemFirst(t *testing.T) {
	ix := Build(corpus())
	for _, q := range []string{"S-0042", "s-0042", "0042"} {
		if got := ix.Search(q, false, 0); len(got) == 0 || got[0].ItemID != "S-0042" {
			t.Errorf("%s: %s", q, paths(got))
		}
	}
}

func TestMoreMatchingTermsRankHigherAndTheLimitHolds(t *testing.T) {
	ix := Build(corpus())
	got := ix.Search("tunnel dashboard host", false, 0)
	if got[0].Path != "design/system/channel.md" && got[0].Path != "wip/kanban/stories/S-0042-tunnel.md" {
		t.Errorf("several terms: %s", paths(got))
	}
	if got := ix.Search("tunnel", true, 2); len(got) != 2 {
		t.Errorf("limit: %d", len(got))
	}
}

func TestSnippet(t *testing.T) {
	body := strings.Repeat("lead in words ", 20) + "the Tunnel is here " + strings.Repeat("and trailing words ", 20)
	s := Snippet(body, []string{"tunnel"}, 160)
	if !strings.HasPrefix(s, "…") || !strings.HasSuffix(s, "…") || !strings.Contains(s, "the Tunnel is here") || len([]rune(s)) != 162 {
		t.Errorf("snippet: %q (%d)", s, len([]rune(s)))
	}
	if s := Snippet("short  body\nover lines", []string{"absent"}, 160); s != "short body over lines" {
		t.Errorf("no term found: %q", s)
	}
	if s := Snippet("ünïcödé 🚢 "+strings.Repeat("x ", 200), nil, 10); s != "ünïcödé 🚢 …" {
		t.Errorf("runes, not bytes: %q", s)
	}
}

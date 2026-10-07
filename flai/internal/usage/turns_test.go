package usage

import (
	"slices"
	"testing"
	"time"
)

// bash is a Bash call of command.
func bash(command string) toolCall { return toolCall{name: "Bash", command: command} }

// S-0293: a turn is put in the first class it matches: a test run by hand,
// then a hand edit of a narrative, task, or story, then an empty wake, then
// ceremony alone, and otherwise work.
func TestATurnIsClassifiedByTheFirstClassItMatches(t *testing.T) {
	inbox := toolCall{name: "mcp__flai__inbox"}
	move := toolCall{name: "mcp__flai__item_move"}
	emptyWait := toolCall{name: waitTool, emptyWake: true}
	wait := toolCall{name: waitTool}
	cases := []struct {
		name  string
		calls []toolCall
		want  string
	}{
		// test runs, by hand, whatever else the turn calls
		{"go test", []toolCall{bash("cd flai && go test ./internal/usage/ 2>&1 | tail -20")}, TurnTestRuns},
		{"go vet", []toolCall{bash("go vet ./...")}, TurnTestRuns},
		{"golangci-lint", []toolCall{bash("cd flai && golangci-lint run ./...")}, TurnTestRuns},
		{"gofmt", []toolCall{bash("gofmt -l .")}, TurnTestRuns},
		{"vitest", []toolCall{bash("npx vitest run")}, TurnTestRuns},
		{"svelte-check", []toolCall{bash("pnpm exec svelte-check")}, TurnTestRuns},
		{"eslint", []toolCall{bash("npx eslint src")}, TurnTestRuns},
		{"prettier", []toolCall{bash("npx prettier --check .")}, TurnTestRuns},
		{"markdownlint", []toolCall{bash("markdownlint-cli2 design/**/*.md")}, TurnTestRuns},
		{"make test", []toolCall{bash("make test")}, TurnTestRuns},
		{"make lint-md", []toolCall{bash("make lint-md")}, TurnTestRuns},
		{"flai-test.sh", []toolCall{bash("scripts/flai-test.sh")}, TurnTestRuns},
		{"flaiover-test.sh", []toolCall{bash("scripts/flaiover-test.sh")}, TurnTestRuns},
		{"lint-md.sh", []toolCall{bash("scripts/lint-md.sh")}, TurnTestRuns},
		{"template-test.sh", []toolCall{bash("scripts/template-test.sh")}, TurnTestRuns},
		{"pnpm run check", []toolCall{bash("cd flaiover && pnpm run check")}, TurnTestRuns},
		{"pnpm test", []toolCall{bash("pnpm test")}, TurnTestRuns},
		{"with ceremony", []toolCall{bash("git status"), bash("go test ./...")}, TurnTestRuns},
		{"before a hand edit", []toolCall{{name: "Edit", path: "wip/agents/S-0001.md"}, bash("go test ./...")}, TurnTestRuns},
		// what replaces a run by hand, and a test named in text
		{"flai test", []toolCall{bash("scripts/flai.sh test flai/internal/usage")}, TurnWork},
		{"flai verify", []toolCall{bash("flai verify --story S-0001")}, TurnWork},
		{"close-out", []toolCall{bash("scripts/close-out.sh S-0001")}, TurnWork},
		{"in a heredoc", []toolCall{bash("git add -A && git commit -F - <<'EOF'\nfeat: run go test and make test\nEOF")}, TurnCeremony},
		{"quoted", []toolCall{bash(`git commit -m "go test passes"`)}, TurnCeremony},
		// hand edits
		{"Edit a narrative", []toolCall{{name: "Edit", path: "/w/S-0001/wip/agents/S-0001.md"}}, TurnHandEdits},
		{"Write a task", []toolCall{{name: "Write", path: "wip/kanban/tasks/T-0001-x.md"}}, TurnHandEdits},
		{"MultiEdit a story", []toolCall{{name: "MultiEdit", path: "wip/kanban/stories/S-0001-x.md"}}, TurnHandEdits},
		{"with an empty wake", []toolCall{emptyWait, {name: "Edit", path: "wip/agents/S-0001.md"}}, TurnHandEdits},
		{"sed -i", []toolCall{bash("sed -i 's/a/b/' wip/agents/S-0001.md")}, TurnHandEdits},
		{"sed -E -i", []toolCall{bash("sed -E -i.bak 's/a/b/' wip/kanban/tasks/T-1.md")}, TurnHandEdits},
		{"perl -pi", []toolCall{bash("perl -pi -e 's/a/b/' wip/kanban/stories/S-1.md")}, TurnHandEdits},
		{"tee", []toolCall{bash("printf x | tee -a wip/agents/S-0001.md")}, TurnHandEdits},
		{"redirection", []toolCall{bash("cat >> wip/agents/S-0001.md <<'EOF'\n- a line\nEOF")}, TurnHandEdits},
		{"python", []toolCall{bash("python3 - <<'EOF'\nfrom pathlib import Path\nPath('wip/agents/S-0001.md').write_text('x')\nEOF")}, TurnHandEdits},
		{"Edit an archived story", []toolCall{{name: "Edit", path: "wip/archive/kanban/stories/S-1.md"}}, TurnWork},
		{"Edit code", []toolCall{{name: "Edit", path: "flai/internal/usage/turns.go"}}, TurnWork},
		{"read a narrative", []toolCall{bash("cat wip/agents/S-0001.md 2>&1 >/dev/null")}, TurnWork},
		{"sed -n a narrative", []toolCall{bash("sed -n '1,20p' wip/agents/S-0001.md")}, TurnWork},
		{"write elsewhere", []toolCall{bash("sed -i 's/a/b/' flai/go.mod")}, TurnWork},
		// empty wakes
		{"an empty wake", []toolCall{emptyWait}, TurnEmptyWakes},
		{"two empty wakes", []toolCall{emptyWait, emptyWait}, TurnEmptyWakes},
		{"a wake with news", []toolCall{wait}, TurnWork},
		{"one of two empty", []toolCall{emptyWait, wait}, TurnWork},
		{"an empty wake and the inbox", []toolCall{emptyWait, inbox}, TurnWork},
		// ceremony
		{"inbox", []toolCall{inbox}, TurnCeremony},
		{"item_move and inbox", []toolCall{move, inbox}, TurnCeremony},
		{"stream sync", []toolCall{bash("cd /w/S-0001 && ../../../scripts/flai.sh stream sync 2>&1 | tail -3")}, TurnCeremony},
		{"stream log", []toolCall{bash(`FLAI_AGENT=x /abs/bin/flai stream log "did a thing"`)}, TurnCeremony},
		{"stream open", []toolCall{bash("flai stream open S-0001")}, TurnCeremony},
		{"flai move in a variable", []toolCall{bash("F=../../scripts/flai.sh; \"$F\" move T-0001 done && $FLAI touches T-0002")}, TurnCeremony},
		{"flai check", []toolCall{bash("set -euo pipefail\nscripts/flai.sh check --strict | head -20")}, TurnCeremony},
		{"git", []toolCall{bash("git -C /w/S-0001 status --short; git log --oneline -5 | head -3\ngit diff --stat && git diff --cached --stat\ngit show --stat HEAD")}, TurnCeremony},
		{"git in a loop", []toolCall{bash("for f in a b; do git add $f; done\nif true; then git commit -m x; fi\necho done")}, TurnCeremony},
		{"rebase --continue", []toolCall{bash("GIT_EDITOR=true git rebase --continue")}, TurnCeremony},
		{"commit with a substituted heredoc", []toolCall{bash("git commit -m \"$(cat <<'EOF'\nfix: it\n\nCo-Authored-By: x\nEOF\n)\"")}, TurnCeremony},
		// work
		{"ceremony and cat", []toolCall{bash("git status && cat wip/kanban/board.md")}, TurnWork},
		{"a filter not piped", []toolCall{bash("git status; grep -n x file")}, TurnWork},
		{"sed -n piped", []toolCall{bash("git log | sed -n 1,3p")}, TurnCeremony},
		{"git diff", []toolCall{bash("git diff")}, TurnWork},
		{"task done", []toolCall{bash("scripts/flai.sh task done T-0001")}, TurnWork},
		{"story start", []toolCall{bash("flai story start S-0001")}, TurnWork},
		{"stream state", []toolCall{bash("flai stream state S-0001")}, TurnWork},
		{"nothing but cd", []toolCall{bash("cd /w && echo hi")}, TurnWork},
		{"ceremony and Read", []toolCall{inbox, {name: "Read", path: "wip/agents/S-0001.md"}}, TurnWork},
	}
	for _, c := range cases {
		if got := classify(c.calls); got != c.want {
			t.Errorf("%s: class = %s, want %s (commands %q)", c.name, got, c.want, commands(c.calls[len(c.calls)-1].command))
		}
	}
}

// S-0293: a Bash command's heredoc bodies and quoted strings are not its
// text; the line that opens a heredoc is, and a quoted variable is kept.
func TestShellTextLeavesOutHeredocsAndQuotes(t *testing.T) {
	cases := []struct{ in, want string }{
		{"git commit -F - <<'EOF'\ngo test\nEOF\ngit status", "git commit -F - <<Q\ngit status"},
		{"cat <<-END\n\tbody\n\tEND", "cat <<-END"},
		{"cat <<\"A\" <<B\na\nA\nb\nB\nls", "cat <<Q <<B\nls"},
		{"grep <<<'here' x", "grep <<<Q x"},
		{`echo 'a "b"' "c \" d" "$F" "${FLAI}"`, "echo Q Q $F ${FLAI}"},
		{`echo "open`, "echo Q"},
	}
	for _, c := range cases {
		if got := shellText(c.in); got != c.want {
			t.Errorf("shellText(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// S-0293: a command's leading keywords and assignments are dropped, and
// commands that only change directory, echo, or filter a pipe are left out.
func TestCommandsAreSplitAndBared(t *testing.T) {
	got := commands("set -e\ncd /x && A=1 B=2 git status 2>&1 | grep x | wc -l || (echo no)\nwhile true; do flai move T-1 done; done & X=1")
	want := []string{"git status", "flai move T-1 done"}
	if !slices.Equal(got, want) {
		t.Errorf("commands = %q, want %q", got, want)
	}
}

// S-0293: turns are counted per UTC day, in order of day; a day's counts add
// up with another's, and a merge adds day by day without sharing either.
func TestTurnDaysAreCountedAndMergedByDay(t *testing.T) {
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }
	days := turnDays([]turn{
		{at: at("2026-10-04T01:00:00+02:00"), calls: []toolCall{{id: "w", name: waitTool}}},
		{at: at("2026-10-03T08:00:00Z"), calls: []toolCall{bash("go test ./...")}},
		{at: at("2026-10-03T09:00:00Z")},
		{calls: []toolCall{bash("go test ./...")}},
	}, map[string]bool{"w": true})
	want := []TurnDay{{Day: "2026-10-03", TestRuns: 1, EmptyWakes: 1}}
	if !slices.Equal(days, want) {
		t.Errorf("days = %+v, want %+v", days, want)
	}
	a := []TurnDay{{Day: "2026-10-02", Work: 1}, {Day: "2026-10-04", Ceremony: 2}}
	b := []TurnDay{{Day: "2026-10-04", Ceremony: 1, Work: 3}, {Day: "2026-10-03", HandEdits: 1}, {Day: "2026-10-05"}}
	got := mergeTurns(a, b)
	want = []TurnDay{{Day: "2026-10-02", Work: 1}, {Day: "2026-10-03", HandEdits: 1}, {Day: "2026-10-04", Ceremony: 3, Work: 3}}
	if !slices.Equal(got, want) {
		t.Errorf("merged = %+v, want %+v", got, want)
	}
	if a[1].Ceremony != 2 || b[0].Ceremony != 1 {
		t.Errorf("a merge changed what it merged: %+v %+v", a, b)
	}
	if mergeTurns(nil, nil) != nil {
		t.Error("merging no turns gives some")
	}
	d := TurnDay{Day: "x", Ceremony: 1, TestRuns: 2, EmptyWakes: 3, HandEdits: 4, Work: 5}
	for i, c := range TurnClasses {
		if d.Count(c) != i+1 {
			t.Errorf("count of %s = %d, want %d", c, d.Count(c), i+1)
		}
	}
	if d.Total() != 15 || d.Count("other") != 0 {
		t.Errorf("total = %d, other = %d", d.Total(), d.Count("other"))
	}
}

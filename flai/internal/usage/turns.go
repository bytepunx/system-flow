package usage

import (
	"regexp"
	"slices"
	"strings"
	"time"
)

// A story's turns (S-0293). A turn is one call to the model by the session's
// own agent, an assistant message without parent_tool_use_id, that uses at
// least one tool; Claude Code repeats a message once per content block, so
// a turn's tool calls are gathered across the repeats of its message ID, and
// its time is its first event's. Each turn is put in the first class it
// matches:
//
//  1. test_runs: a Bash call runs a test, lint, or format tool by hand, in its
//     command with heredoc bodies and quoted strings removed.
//  2. hand_edits: a call edits a narrative, a task, or a story by hand: an
//     Edit, Write, or MultiEdit of its file, or a Bash call that names its
//     path and writes.
//  3. empty_wakes: every call is a wait_for_events that was an empty wake.
//  4. ceremony: every call is inbox, item_move, or a Bash call whose every
//     command is git's add, commit, status, log, or stat, or flai's stream
//     sync, log, or open, move, touches, or check.
//  5. work: every other turn.
//
// The turns are counted per UTC day and class into Usage.Turns.

// The classes of a turn, as front matter and JSON name them.
const (
	TurnCeremony   = "ceremony"
	TurnTestRuns   = "test_runs"
	TurnEmptyWakes = "empty_wakes"
	TurnHandEdits  = "hand_edits"
	TurnWork       = "work"
)

// TurnClasses are the classes of a turn in the order they are written.
var TurnClasses = []string{TurnCeremony, TurnTestRuns, TurnEmptyWakes, TurnHandEdits, TurnWork}

// TurnDay counts the turns of a story's own agent on one day, by class (S-0293).
type TurnDay struct {
	// Day is the UTC date, as 2006-01-02.
	Day        string `yaml:"day" json:"day"`
	Ceremony   int    `yaml:"ceremony,omitempty" json:"ceremony"`
	TestRuns   int    `yaml:"test_runs,omitempty" json:"test_runs"`
	EmptyWakes int    `yaml:"empty_wakes,omitempty" json:"empty_wakes"`
	HandEdits  int    `yaml:"hand_edits,omitempty" json:"hand_edits"`
	Work       int    `yaml:"work,omitempty" json:"work"`
}

// DayFormat is the layout of a TurnDay's day.
const DayFormat = "2006-01-02"

// field is the count of class, one of TurnClasses; nil for another.
func (d *TurnDay) field(class string) *int {
	switch class {
	case TurnCeremony:
		return &d.Ceremony
	case TurnTestRuns:
		return &d.TestRuns
	case TurnEmptyWakes:
		return &d.EmptyWakes
	case TurnHandEdits:
		return &d.HandEdits
	case TurnWork:
		return &d.Work
	}
	return nil
}

// Count is the day's turns of class, one of TurnClasses; 0 for another.
func (d TurnDay) Count(class string) int {
	if n := d.field(class); n != nil {
		return *n
	}
	return 0
}

// Total is the day's turns of every class.
func (d TurnDay) Total() int {
	return d.Ceremony + d.TestRuns + d.EmptyWakes + d.HandEdits + d.Work
}

// Add adds o's counts to d's, whatever o's day.
func (d *TurnDay) Add(o TurnDay) {
	for _, c := range TurnClasses {
		*d.field(c) += o.Count(c)
	}
}

// mergeTurns is a and b added day by day, in order of day, sharing neither's
// array; nil when both have none.
func mergeTurns(a, b []TurnDay) []TurnDay {
	if len(a)+len(b) == 0 {
		return nil
	}
	out := slices.Clone(a)
	for _, d := range b {
		if i := slices.IndexFunc(out, func(x TurnDay) bool { return x.Day == d.Day }); i >= 0 {
			out[i].Add(d)
			continue
		}
		out = append(out, d)
	}
	return tidyTurns(out)
}

// tidyTurns sorts days by date and drops those with no turns; nil when none
// is left.
func tidyTurns(days []TurnDay) []TurnDay {
	days = slices.DeleteFunc(days, func(d TurnDay) bool { return d.Total() == 0 })
	if len(days) == 0 {
		return nil
	}
	slices.SortFunc(days, func(a, b TurnDay) int { return strings.Compare(a.Day, b.Day) })
	return days
}

// toolCall is one tool_use of a turn.
type toolCall struct {
	// id is the tool_use's ID, name the tool's.
	id, name string
	// command is a Bash call's command; path the file an Edit, Write, or
	// MultiEdit call changes.
	command, path string
	// emptyWake says a wait_for_events call's result was an empty wake.
	emptyWake bool
}

// turn is one message of the session's own agent and the tools it called.
type turn struct {
	at    time.Time
	calls []toolCall
}

// turnDays counts turns per UTC day and class, in order of day; a message
// that called no tool is not a turn. emptied are the IDs of the
// wait_for_events calls that were empty wakes.
func turnDays(turns []turn, emptied map[string]bool) []TurnDay {
	byDay := map[string]*TurnDay{}
	for _, t := range turns {
		if len(t.calls) == 0 || t.at.IsZero() {
			continue
		}
		calls := slices.Clone(t.calls)
		for i := range calls {
			calls[i].emptyWake = calls[i].name == waitTool && emptied[calls[i].id]
		}
		day := t.at.UTC().Format(DayFormat)
		if byDay[day] == nil {
			byDay[day] = &TurnDay{Day: day}
		}
		*byDay[day].field(classify(calls))++
	}
	var out []TurnDay
	for _, d := range byDay {
		out = append(out, *d)
	}
	return tidyTurns(out)
}

// classify is the class of a turn that made calls, the first that matches.
func classify(calls []toolCall) string {
	switch {
	case slices.ContainsFunc(calls, runsTests):
		return TurnTestRuns
	case slices.ContainsFunc(calls, editsByHand):
		return TurnHandEdits
	case !slices.ContainsFunc(calls, func(c toolCall) bool { return !c.emptyWake }):
		return TurnEmptyWakes
	case !slices.ContainsFunc(calls, func(c toolCall) bool { return !isCeremony(c) }):
		return TurnCeremony
	}
	return TurnWork
}

// testTool is a test, lint, or format tool run by hand rather than through
// flai test, flai verify, or the close-out.
var testTool = regexp.MustCompile(`\bgo (test|vet)\b|golangci-lint|\bgofmt\b|\bvitest\b|svelte-check|\beslint\b|prettier --check|markdownlint|` +
	`\bmake (test|integration|smoke|flai-test|flaiover-test|lint|lint-md)\b|flai-test\.sh|flaiover-test\.sh|lint-md\.sh|template-test\.sh|` +
	`\bpnpm (run )?(test|lint|check)\b`)

// runsTests says c is a Bash call that runs a test, lint, or format tool;
// what heredocs and quoted strings say is not run.
func runsTests(c toolCall) bool {
	return c.name == "Bash" && testTool.MatchString(shellText(c.command))
}

var (
	// workItemPath is the path of a story's narrative, or of a task or story.
	workItemPath = regexp.MustCompile(`wip/agents/S-\d+\.md|wip/kanban/(tasks|stories)/`)
	// inPlace is sed, perl, or tee writing a file, in shell text.
	inPlace = regexp.MustCompile(`\bsed\s+(-\S+\s+)*(-[a-zA-Z]*i|--in-place)|\bperl\s+(-\S+\s+)*-[a-zA-Z]*i|\btee\s`)
	// redirect is output redirected to a file, in shell text; its target is
	// the first group.
	redirect = regexp.MustCompile(`(?:^|[^0-9&<>=-])\d?>>?\s*([^\s;&|<>()]+)`)
	// pythonWrite is Python writing a file.
	pythonWrite = regexp.MustCompile(`\bpython3?\b[\s\S]*(\.write\(|write_text\()`)
)

// editsByHand says c edits a narrative, a task, or a story: an Edit, Write,
// or MultiEdit of its file, or a Bash command that names its path and writes
// with sed -i, perl -i, tee, a redirection to a file, or Python.
func editsByHand(c toolCall) bool {
	switch c.name {
	case "Edit", "Write", "MultiEdit":
		return workItemPath.MatchString(c.path)
	case "Bash":
		if !workItemPath.MatchString(c.command) {
			return false
		}
		text := shellText(c.command)
		return inPlace.MatchString(text) || redirects(text) || pythonWrite.MatchString(c.command)
	}
	return false
}

// redirects says shell text redirects output to a file other than
// /dev/null.
func redirects(text string) bool {
	for _, m := range redirect.FindAllStringSubmatch(text, -1) {
		if m[1] != "/dev/null" {
			return true
		}
	}
	return false
}

// isCeremony says c is a ceremony call: inbox, item_move, or a Bash call of
// ceremony commands alone.
func isCeremony(c toolCall) bool {
	switch c.name {
	case "mcp__flai__inbox", "mcp__flai__item_move":
		return true
	case "Bash":
		cmds := commands(c.command)
		return len(cmds) > 0 && !slices.ContainsFunc(cmds, func(cmd string) bool { return !ceremonyCommand(cmd) })
	}
	return false
}

var (
	// gitCeremony and flaiCeremony are what follows git, or flai, in a
	// ceremony command.
	gitCeremony  = []string{"add", "commit", "status", "log", "diff --stat", "diff --cached --stat", "show --stat", "rebase --continue"}
	flaiCeremony = []string{"stream sync", "stream log", "stream open", "move", "touches", "check"}
	// flaiVars are the shell variables an agent keeps flai's path in.
	flaiVars = []string{"$F", "$FLAI", "${F}", "${FLAI}"}
)

// ceremonyCommand says a command is git's add, commit, status, log, or
// stat, or flai's stream sync, log, or open, move, touches, or check.
func ceremonyCommand(cmd string) bool {
	words := strings.Fields(cmd)
	if len(words) < 2 {
		return false
	}
	first, rest := words[0], words[1:]
	var follows []string
	switch {
	case first == "git":
		if rest[0] == "-C" && len(rest) > 2 {
			rest = rest[2:]
		}
		follows = gitCeremony
	case strings.HasSuffix(first, "flai") || strings.HasSuffix(first, "flai.sh") || slices.Contains(flaiVars, first):
		follows = flaiCeremony
	default:
		return false
	}
	args := strings.Join(rest, " ")
	return slices.ContainsFunc(follows, func(f string) bool { return args == f || strings.HasPrefix(args, f+" ") })
}

var (
	// separator separates the commands of shell text.
	separator = regexp.MustCompile(`\|\||&&|[;|&\n]`)
	// redirection is a redirection or a heredoc's opening, with its target.
	redirection = regexp.MustCompile(`\d*(>>?|<<-?|<)&?\s*[^\s;&|()]*`)
	// assignment is a variable assigned before a command, or alone.
	assignment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=\S*(\s+|$)`)
	// keyword is a shell keyword a command follows.
	keyword = regexp.MustCompile(`^(do|then|else)(\s+|$)`)
)

// skipped are the first words of commands that are neither ceremony nor
// work; filters are those of commands that only filter what a pipe gives
// them.
var skipped, filters = []string{"cd", "export", "echo", "true", "for", "while", "if", "done", "fi", "exit"},
	[]string{"head", "tail", "grep", "wc", "jq", "sort", "uniq", "cut", "tr"}

// commands are the commands of a Bash call's command that count: heredoc
// bodies, quoted strings, and redirections removed, split on &&, ||, ;, |,
// &, and newlines, each without its leading keyword or assignments; those
// that only change directory, export, echo, or open or close a loop or
// condition, and the filters a pipe feeds, are left out.
func commands(command string) []string {
	text := redirection.ReplaceAllString(shellText(command), " ")
	var out []string
	piped := false
	for {
		loc := separator.FindStringIndex(text)
		cmd := text
		if loc != nil {
			cmd = text[:loc[0]]
		}
		if cmd = bare(cmd); cmd != "" && !skip(cmd, piped) {
			out = append(out, cmd)
		}
		if loc == nil {
			return out
		}
		piped = text[loc[0]:loc[1]] == "|"
		text = text[loc[1]:]
	}
}

// bare is a command without the spaces, parentheses, and braces around it,
// its leading keywords, and the variables assigned before it.
func bare(cmd string) string {
	for {
		was := cmd
		cmd = strings.Trim(cmd, " \t(){}")
		cmd = keyword.ReplaceAllString(cmd, "")
		if a := assignment.FindString(cmd); a != "" && strings.TrimSpace(a) != cmd {
			cmd = cmd[len(a):]
		}
		if cmd == was {
			return cmd
		}
	}
}

// skip says a command does not count: a bare assignment, one whose first
// word is in skipped, set with options, or, after a pipe, a filter.
func skip(cmd string, piped bool) bool {
	words := strings.Fields(cmd)
	switch {
	case assignment.MatchString(cmd) && len(words) == 1:
		return true
	case slices.Contains(skipped, words[0]):
		return true
	case words[0] == "set" && len(words) > 1 && strings.HasPrefix(words[1], "-"):
		return true
	case piped && slices.Contains(filters, words[0]):
		return true
	case piped && words[0] == "sed" && len(words) > 1 && words[1] == "-n":
		return true
	}
	return false
}

// heredoc is a heredoc's opening, its delimiter in one of the groups.
var heredoc = regexp.MustCompile(`<<-?\s*(?:'([A-Za-z_]\w*)'|"([A-Za-z_]\w*)"|\\?([A-Za-z_]\w*))`)

// shellText is a shell command without its heredocs' bodies, the lines that
// open them kept, and with each quoted string replaced by Q, save a double-
// quoted variable, which is kept without its quotes.
func shellText(command string) string {
	var kept []string
	var ends []string // the delimiters of the heredocs still open
	for _, l := range strings.Split(command, "\n") {
		if len(ends) > 0 {
			if strings.TrimSpace(l) == ends[0] {
				ends = ends[1:]
			}
			continue
		}
		kept = append(kept, l)
		for _, m := range heredoc.FindAllStringSubmatchIndex(l, -1) {
			if m[0] > 0 && l[m[0]-1] == '<' {
				continue // a here-string, <<<
			}
			for g := 2; g < len(m); g += 2 {
				if m[g] >= 0 {
					ends = append(ends, l[m[g]:m[g+1]])
				}
			}
		}
	}
	return unquoted(strings.Join(kept, "\n"))
}

// quotedVar is a variable alone, as quoted strings hold one.
var quotedVar = regexp.MustCompile(`^\$\{?[A-Za-z_]\w*\}?$`)

// unquoted is text with each single- or double-quoted string replaced by Q,
// save a double-quoted variable, kept without its quotes. A quote left open
// runs to the end.
func unquoted(text string) string {
	var b strings.Builder
	for i := 0; i < len(text); i++ {
		q := text[i]
		if q != '\'' && q != '"' {
			b.WriteByte(q)
			continue
		}
		j := i + 1
		for j < len(text) && text[j] != q {
			if q == '"' && text[j] == '\\' {
				j++
			}
			j++
		}
		inside := text[i+1 : min(j, len(text))]
		if q == '"' && quotedVar.MatchString(inside) {
			b.WriteString(inside)
		} else {
			b.WriteByte('Q')
		}
		i = j
	}
	return b.String()
}

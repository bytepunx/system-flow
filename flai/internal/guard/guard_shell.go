package guard

import (
	"path"
	"slices"
	"strings"
)

// shells are the programs that run a script they read, so that a heredoc fed
// to one is read as commands.
var shells = []string{"sh", "bash", "zsh", "dash", "ksh"}

// heredoc is a here-document a line opens: the delimiter that ends its body,
// with quotes removed, and whether <<- strips the body's leading tabs.
type heredoc struct {
	delim string
	tabs  bool
}

// commands splits a shell command line into simple commands, each as its
// words with quotes removed: at ;, &, |, newlines, and the openings of
// subshells and command substitutions. It is not a shell; it is enough to
// find each program a line runs and the words after it. A heredoc's body is
// input, not commands, so it is left out, unless a shell is named on the
// line that opens it, as in bash <<EOF or cat <<EOF | sh, when its body is
// split as a script (S-0246). A redirection's operator stays in its word, so
// that 2>&1, &>f, and >|f split no command, and so does find's and xargs's
// {} (S-0368).
func commands(line string) [][]string {
	var out [][]string
	for _, c := range parse(line) {
		out = append(out, c.words)
	}
	return out
}

// arguments splits a shell command line into simple commands as commands
// does, each as its arguments: its words without its redirections and their
// targets, so that a redirection such as 2>&1 or 2>/dev/null, wherever it
// stands, neither allows nor refuses a command a check judges (I-0230).
func arguments(line string) [][]string {
	var out [][]string
	for _, c := range parse(line) {
		if len(c.args) > 0 {
			out = append(out, c.args)
		}
	}
	return out
}

// simple is one simple command of a line, as parse splits it: its words, its
// arguments, the words without its redirections and their targets, and the
// files its output redirections write, as written (S-0368). A descriptor
// duplication, such as 2>&1 or >&2, writes no file.
type simple struct {
	words, args, targets []string
}

// The kinds of word a redirection operator written apart from its target
// waits for.
const (
	noTarget = iota
	fileTarget
	dupTarget
	inputTarget
)

// parse splits a shell command line into simple commands, as commands
// describes: a > or < outside quotes and arithmetic begins a redirection,
// with the descriptor number or & before it, and its target is the rest of
// its word or, when nothing follows the operator, the next word.
func parse(line string) []simple {
	var out []simple
	var cur simple
	var word strings.Builder
	inWord := false
	var quote rune
	var pending []heredoc
	start, arith, skip := 0, 0, 0
	redir, awaits := -1, noTarget
	end := func() {
		if !inWord {
			return
		}
		w := word.String()
		cur.words = append(cur.words, w)
		switch {
		case redir >= 0:
			if redir > 0 {
				cur.args = append(cur.args, w[:redir])
			}
			target, kind := operator(w[redir:])
			if target == "" {
				awaits = kind
			} else if kind == fileTarget || kind == dupTarget && !descriptor(target) {
				cur.targets = append(cur.targets, target)
			}
		case awaits != noTarget:
			if awaits == fileTarget || awaits == dupTarget && !descriptor(w) {
				cur.targets = append(cur.targets, w)
			}
			awaits = noTarget
		default:
			cur.args = append(cur.args, w)
		}
		word.Reset()
		inWord, redir = false, -1
	}
	split := func() {
		end()
		if len(cur.words) > 0 {
			out = append(out, cur)
		}
		cur, awaits = simple{}, noTarget
	}
	for i, r := range line {
		if i < skip {
			continue
		}
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			word.WriteRune(r)
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == '<' && arith == 0 && strings.HasPrefix(line[i:], "<<") && !strings.HasPrefix(line[i:], "<<<") && !strings.HasSuffix(line[:i], "<"):
			end()
			h, next, ok := opens(line, i)
			if !ok {
				word.WriteString("<<")
				inWord, skip = true, i+2
				continue
			}
			pending, skip = append(pending, h), next
		case r == '\n' && len(pending) > 0:
			split()
			body, next := bodies(line, i+1, pending)
			if slices.ContainsFunc(out[start:], func(c simple) bool { return runsShell(c.words) }) {
				out = append(out, parse(body)...)
			}
			pending, start, skip = nil, len(out), next
		case r == '(' && strings.HasPrefix(line[i:], "(("):
			arith++
			split()
		case r == ')' && arith > 0 && strings.HasPrefix(line[i:], "))"):
			arith--
			split()
		case (r == '&' || r == '|') && redir >= 0 && line[i-1] == '>', r == '&' && redir >= 0 && line[i-1] == '<':
			word.WriteRune(r)
		case r == '{' && strings.HasPrefix(line[i:], "{}"):
			word.WriteString("{}")
			inWord, skip = true, i+2
		case r == '&' && strings.HasPrefix(line[i+1:], ">"):
			end()
			word.WriteRune(r)
			inWord = true
		case r == ';' || r == '&' || r == '|' || r == '\n' || r == '(' || r == ')' || r == '`' || r == '{' || r == '}':
			split()
			if r == '\n' {
				start = len(out)
			}
		case r == '$':
			end()
		case r == ' ' || r == '\t':
			end()
		default:
			if (r == '>' || r == '<') && arith == 0 && redir < 0 {
				redir = word.Len()
				if strings.Trim(word.String(), "0123456789&") == "" {
					redir = 0
				}
			}
			word.WriteRune(r)
			inWord = true
		}
	}
	split()
	return out
}

// operator reads a redirection written from its operator on, as 2>, >>f, or
// >&2: its target, "" when the next word is, and what kind the target is.
func operator(w string) (target string, kind int) {
	w = strings.TrimLeft(w, "0123456789&")
	if rest, ok := strings.CutPrefix(w, "<"); ok {
		return strings.TrimLeft(rest, "<>&"), inputTarget
	}
	w = strings.TrimPrefix(w, ">")
	switch {
	case strings.HasPrefix(w, "&"):
		return w[1:], dupTarget
	case strings.HasPrefix(w, ">"), strings.HasPrefix(w, "|"):
		return w[1:], fileTarget
	}
	return w, fileTarget
}

// descriptor says whether a duplication's target is a descriptor, as in
// >&2, 2>&1-, or >&-, rather than a file, as in >&out.txt.
func descriptor(target string) bool {
	return strings.Trim(strings.TrimSuffix(target, "-"), "0123456789") == ""
}

// opens reads the heredoc operator at line[i:], << or <<-, and the delimiter
// word after it: the heredoc, and where the line goes on after the word. It
// is not one when no word follows.
func opens(line string, i int) (h heredoc, next int, ok bool) {
	j := i + 2
	if j < len(line) && line[j] == '-' {
		h.tabs = true
		j++
	}
	for j < len(line) && (line[j] == ' ' || line[j] == '\t') {
		j++
	}
	var delim strings.Builder
	var quote byte
	for ; j < len(line); j++ {
		c := line[j]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
				continue
			}
			delim.WriteByte(c)
		case c == '\'' || c == '"':
			quote = c
		case c == '\\' && j+1 < len(line):
			j++
			delim.WriteByte(line[j])
		case strings.IndexByte(" \t\n;&|<>()`", c) >= 0:
			h.delim = delim.String()
			return h, j, h.delim != ""
		default:
			delim.WriteByte(c)
		}
	}
	h.delim = delim.String()
	return h, j, h.delim != ""
}

// bodies reads the bodies of the heredocs pending, in order, from line[from:]
// on: their text, and where the line goes on after the last one's delimiter.
// A body with no delimiter runs to the end, as in bash.
func bodies(line string, from int, pending []heredoc) (string, int) {
	var text strings.Builder
	at := from
	for _, h := range pending {
		for at < len(line) {
			l, _, _ := strings.Cut(line[at:], "\n")
			at = min(at+len(l)+1, len(line))
			mark := strings.TrimSuffix(l, "\r")
			if h.tabs {
				mark = strings.TrimLeft(mark, "\t")
			}
			if mark == h.delim {
				break
			}
			text.WriteString(l + "\n")
		}
	}
	return text.String(), at
}

// runsShell says whether a simple command names a shell, which reads a
// heredoc fed to it, or piped to it, as a script.
func runsShell(words []string) bool {
	return slices.ContainsFunc(words, func(w string) bool {
		return slices.Contains(shells, path.Base(w))
	})
}

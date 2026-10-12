package guard

import (
	"cmp"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// builtInReaders are Claude Code's built-in sub-agents that have neither Edit
// nor Write, whatever a project's definitions say (ADR-0141).
var builtInReaders = []string{"Explore", "Plan"}

// ReadOnlyAgents says whether a sub-agent of a type is read-only, from the
// project at root (S-0368, ADR-0141): Explore and Plan are, and so is one
// whose definition, .claude/agents/<type>.md, has a tools line that names
// none of Edit, MultiEdit, NotebookEdit, Write, or *. A definition it cannot
// read, one without a tools line, which grants every tool, or a type that is
// not a file name is not read-only, and with root "" only Explore and Plan
// are. It reads the definition on each call, so that a definition given
// Edit or Write makes its sub-agent a writer at once.
func ReadOnlyAgents(root string) func(agentType string) bool {
	return func(agentType string) bool {
		if slices.Contains(builtInReaders, agentType) {
			return true
		}
		if root == "" || agentType == "" || strings.ContainsAny(agentType, `/\`) || strings.Contains(agentType, "..") {
			return false
		}
		data, err := os.ReadFile(filepath.Join(root, ".claude", "agents", agentType+".md"))
		if err != nil {
			return false
		}
		fm, _, err := workitem.SplitFrontMatter(string(data))
		if err != nil {
			return false
		}
		for line := range strings.SplitSeq(fm, "\n") {
			tools, ok := strings.CutPrefix(line, "tools:")
			if !ok {
				continue
			}
			tools = strings.Trim(strings.TrimSpace(tools), "[]")
			if tools == "" {
				return false
			}
			for tool := range strings.SplitSeq(tools, ",") {
				tool = strings.Trim(strings.TrimSpace(tool), `"'`)
				if tool == "*" || slices.Contains(fileWrites, tool) {
					return false
				}
			}
			return true
		}
		return false
	}
}

// wrappers are the programs that run a command given after their own flags,
// each with its flags that take a value; timeout takes its duration too.
var wrappers = map[string]map[string]bool{
	"command": nil,
	"env":     {"-u": true, "-C": true, "--unset": true, "--chdir": true},
	"exec":    nil,
	"nice":    {"-n": true},
	"nohup":   nil,
	"stdbuf":  {"-i": true, "-o": true, "-e": true},
	"sudo":    {"-u": true, "-g": true, "-C": true, "-D": true, "-p": true, "-U": true},
	"time":    nil,
	"timeout": {"-s": true, "-k": true, "--signal": true, "--kill-after": true},
	"xargs":   {"-I": true, "-n": true, "-P": true, "-L": true, "-d": true, "-E": true, "-s": true, "-a": true},
}

// operandValues are the flags that take a value of the programs whose every
// operand is a file they write.
var operandValues = map[string]map[string]bool{
	"mkdir":    {"-m": true, "--mode": true},
	"mv":       {"-S": true, "--suffix": true},
	"rm":       nil,
	"rmdir":    nil,
	"tee":      nil,
	"touch":    {"-d": true, "-r": true, "-t": true, "--date": true, "--reference": true},
	"truncate": {"-s": true, "-r": true, "--size": true, "--reference": true},
}

// copyValues are the flags that take a value of cp, install, and ln, whose
// last operand is the file they write.
var copyValues = map[string]bool{"-S": true, "-m": true, "-o": true, "-g": true, "-t": true, "--suffix": true, "--mode": true, "--owner": true, "--group": true}

// patchValues are the flags of patch that take a value.
var patchValues = map[string]bool{"-B": true, "-D": true, "-F": true, "-V": true, "-Y": true, "-d": true, "-i": true, "-o": true, "-p": true, "-r": true, "-z": true, "--input": true, "--output": true}

// patchNames stands for what patch writes when no operand names the file.
const patchNames = "the files the patch names"

// written is a file the simple command c writes outside the temporary folder
// (S-0368, ADR-0141), "" when it writes none: a target of its redirections,
// or an operand of tee, sed -i, perl -i, cp, mv, rm, touch, mkdir, rmdir,
// truncate, ln, install, patch, or dd of=, or of a formatter given its write
// flag. Its program is found past assignments and wrappers such as env,
// sudo, timeout, and xargs; a shell's -c script, eval's words, and find
// -exec's command are judged too.
func written(c simple) string {
	if file := outside(c.targets); file != "" {
		return file
	}
	i := program(c.args)
	if i == len(c.args) {
		return ""
	}
	name, rest := path.Base(c.args[i]), c.args[i+1:]
	switch {
	case slices.Contains(shells, name):
		for j, f := range rest {
			if shellFlags.MatchString(f) && j+1 < len(rest) {
				return writtenBy(rest[j+1])
			}
		}
		return ""
	case name == "eval":
		return writtenBy(strings.Join(rest, " "))
	case name == "find":
		for j, f := range rest {
			if f == "-exec" || f == "-execdir" || f == "-ok" || f == "-okdir" {
				return written(simple{args: rest[j+1:]})
			}
		}
		return ""
	}
	return outside(operands(name, rest))
}

// writtenBy is a file a script writes outside the temporary folder, as
// written finds one in each of its simple commands.
func writtenBy(script string) string {
	for _, c := range parse(script) {
		if file := written(c); file != "" {
			return file
		}
	}
	return ""
}

// outside is the first of files that is not under the temporary folder.
func outside(files []string) string {
	for _, f := range files {
		if !temporary(f) {
			return f
		}
	}
	return ""
}

// temporary says whether a file is one a read-only sub-agent may write: an
// absolute path in or under /tmp, os.TempDir, or /dev, so that a folder
// there may be a cp's or mv's destination. A relative path resolves inside
// the checkout, and a variable is not expanded, so neither is (ADR-0141). An
// os.TempDir of / is not taken as the temporary folder.
func temporary(file string) bool {
	if !path.IsAbs(file) {
		return false
	}
	clean := path.Clean(file)
	return slices.ContainsFunc([]string{"/tmp", "/dev", path.Clean(filepath.ToSlash(os.TempDir()))}, func(dir string) bool {
		return dir != "/" && (clean == dir || strings.HasPrefix(clean, dir+"/"))
	})
}

// program is the index in args of the program a simple command runs, past
// assignments and wrappers and their flags; len(args) when there is none.
func program(args []string) int {
	i := 0
	for i < len(args) {
		w := args[i]
		if assignment.MatchString(w) {
			i++
			continue
		}
		values, ok := wrappers[path.Base(w)]
		if !ok {
			return i
		}
		for i++; i < len(args) && strings.HasPrefix(args[i], "-"); i++ {
			if values[args[i]] {
				i++
			}
		}
		if path.Base(w) == "timeout" {
			i++
		}
	}
	return len(args)
}

// operands are the files a program given rest writes, whichever it reaches.
func operands(name string, rest []string) []string {
	if values, ok := operandValues[name]; ok {
		return positionals(rest, values)
	}
	switch name {
	case "sed":
		return sedFiles(rest)
	case "perl":
		return perlFiles(rest)
	case "cp", "install", "ln":
		if t := cmp.Or(value(rest, "-t"), value(rest, "--target-directory")); t != "" {
			return []string{t}
		}
		files := positionals(rest, copyValues)
		switch {
		case name == "install" && slices.Contains(rest, "-d"):
			return files
		case name == "ln" && len(files) == 1:
			return []string{"."}
		case len(files) > 0:
			return files[len(files)-1:]
		}
	case "patch":
		if slices.Contains(rest, "--dry-run") {
			return nil
		}
		if o := cmp.Or(value(rest, "-o"), value(rest, "--output")); o != "" {
			return []string{o}
		}
		if files := positionals(rest, patchValues); len(files) > 0 {
			return files[:1]
		}
		return []string{patchNames}
	case "dd":
		var files []string
		for _, w := range rest {
			if f, ok := strings.CutPrefix(w, "of="); ok {
				files = append(files, f)
			}
		}
		return files
	case "gofmt", "goimports":
		if slices.Contains(rest, "-w") || slices.Contains(rest, "-w=true") {
			return positionals(rest, map[string]bool{"-r": true, "-local": true, "-srcdir": true})
		}
	case "go":
		if files := positionals(rest, nil); len(files) > 0 && files[0] == "fmt" {
			return orHere(files[1:])
		}
	case "prettier":
		if slices.Contains(rest, "--write") || slices.Contains(rest, "-w") {
			return orHere(positionals(rest, nil))
		}
	case "eslint":
		if slices.Contains(rest, "--fix") {
			return orHere(positionals(rest, nil))
		}
	case "golangci-lint":
		if slices.Contains(rest, "--fix") || slices.Contains(rest, "--fix=true") {
			files := positionals(rest, nil)
			if len(files) > 0 && files[0] == "run" {
				files = files[1:]
			}
			return orHere(files)
		}
	}
	return nil
}

// orHere is files, or the working directory, which a formatter given no
// operand writes in, when there are none.
func orHere(files []string) []string {
	if len(files) == 0 {
		return []string{"."}
	}
	return files
}

// sedFiles are the files sed edits in place, none without -i or --in-place:
// its operands but the script, unless -e or -f gives the script. A short
// flag's letters are read as sed reads them: i takes the rest of the word as
// a suffix, and e, f, and l a value, the rest of the word or the next.
func sedFiles(rest []string) []string {
	inPlace, script := false, false
	var files []string
	for i := 0; i < len(rest); i++ {
		w := rest[i]
		switch {
		case w == "--in-place" || strings.HasPrefix(w, "--in-place="):
			inPlace = true
		case w == "--expression" || w == "--file":
			script = true
			i++
		case strings.HasPrefix(w, "--expression=") || strings.HasPrefix(w, "--file="):
			script = true
		case w == "--line-length":
			i++
		case strings.HasPrefix(w, "--"):
			// another long flag, which takes no value
		case strings.HasPrefix(w, "-") && len(w) > 1:
			for j, l := range w[1:] {
				if l == 'i' {
					inPlace = true
					break
				}
				if l == 'e' || l == 'f' || l == 'l' {
					script = script || l != 'l'
					if j == len(w)-2 {
						i++
					}
					break
				}
			}
		default:
			files = append(files, w)
		}
	}
	if !inPlace {
		return nil
	}
	if !script && len(files) > 0 {
		files = files[1:]
	}
	return files
}

// perlFiles are the files perl edits in place, none without -i: its
// operands but the script's file, unless -e or -E gives the script. A short
// flag's letters are read as perl reads them: i takes the rest of the word as
// a suffix, e and E the next word as code, and the others that take a value
// the rest of the word or, alone, the next.
func perlFiles(rest []string) []string {
	inPlace, script := false, false
	var files []string
	for i := 0; i < len(rest); i++ {
		w := rest[i]
		if !strings.HasPrefix(w, "-") || w == "-" {
			files = append(files, w)
			continue
		}
		if strings.HasPrefix(w, "--") {
			continue
		}
		for j, l := range w[1:] {
			last := j == len(w)-2
			if l == 'i' {
				inPlace = true
				break
			}
			if l == 'e' || l == 'E' {
				script = true
				if last {
					i++
				}
				break
			}
			if strings.ContainsRune("CdDIMmVx", l) {
				if last && strings.ContainsRune("IMm", l) {
					i++
				}
				break
			}
		}
	}
	if !inPlace {
		return nil
	}
	if !script && len(files) > 0 {
		files = files[1:]
	}
	return files
}

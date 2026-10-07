package verify

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// plainLines is how many of a command's last lines a plain finding keeps.
const plainLines = 20

// messageLines is how many lines a parsed finding's message keeps.
const messageLines = 10

// env is what parsing needs to make a tool's paths relative to the root.
type env struct {
	// root is the checkout's root on disk.
	root string
	// dir is the tier's Dir, cleaned: "" for the root.
	dir string
	// fsys is the checkout.
	fsys fs.FS
}

// rel is a path a tool printed, relative to the checkout's root: an absolute
// one made relative when it lies below the root, and one relative to the
// tier's Dir joined to it.
func (e env) rel(p string) string {
	if filepath.IsAbs(p) || path.IsAbs(p) {
		if r, err := filepath.Rel(e.root, p); err == nil && r != ".." && !strings.HasPrefix(r, ".."+string(filepath.Separator)) {
			return filepath.ToSlash(r)
		}
		return p
	}
	return path.Join(e.dir, filepath.ToSlash(p))
}

// parse is the findings a failing tier's output gives, as its format reads
// them, or the plain finding when the format finds none.
func parse(t Tier, out output, e env) []Finding {
	var got []Finding
	if out.err == nil {
		switch t.Format {
		case FormatGoTestJSON:
			got = parseGoTest(out.stdout, out.stderr, e)
		case FormatVitestJSON:
			got = parseVitest(out.stdout, e)
		case FormatGolangciJSON:
			got = parseGolangci(out.stdout, e)
		case FormatGofmtList:
			got = parseGofmt(out.stdout, out.stderr, e)
		}
	}
	if len(got) == 0 {
		got = []Finding{plain(t.Name, out)}
	}
	return got
}

// plain is a tier's failure told by its last lines and its exit status, or
// by why it could not run.
func plain(name string, out output) Finding {
	tail := lastLines(out.combined, plainLines)
	status := fmt.Sprintf("exit status %d", out.exit)
	if out.err != nil {
		status = out.err.Error()
	}
	if tail != "" {
		status = tail + "\n" + status
	}
	return Finding{Name: name, Message: status}
}

// parseGofmt is a finding for each file gofmt -l lists, after one for each
// file it could not parse.
func parseGofmt(stdout, stderr string, e env) []Finding {
	var got []Finding
	for _, line := range strings.Split(stderr, "\n") {
		if m := compileLine.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			n, _ := strconv.Atoi(m[2])
			got = append(got, Finding{Name: "gofmt", Path: e.rel(m[1]), Line: n, Message: m[3]})
		}
	}
	for _, line := range strings.Split(stdout, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			got = append(got, Finding{Name: "gofmt", Path: e.rel(line), Message: "not gofmt-formatted"})
		}
	}
	return got
}

// golangciReport is the part of golangci-lint's JSON output read here.
type golangciReport struct {
	Issues []struct {
		FromLinter string
		Text       string
		Pos        struct {
			Filename string
			Line     int
		}
	}
}

// parseGolangci is a finding for each issue golangci-lint's JSON output
// reports.
func parseGolangci(stdout string, e env) []Finding {
	var r golangciReport
	if !decodeFirst(stdout, "Issues", &r) {
		return nil
	}
	got := make([]Finding, 0, len(r.Issues))
	for _, is := range r.Issues {
		got = append(got, Finding{
			Name: is.FromLinter, Path: e.rel(is.Pos.Filename), Line: is.Pos.Line,
			Message: clip(is.Text),
		})
	}
	return got
}

// decodeFirst decodes into v the first JSON object in s, at the start of a
// line, that has the key: tools print other lines before and after it.
func decodeFirst(s, key string, v any) bool {
	for i := 0; i < len(s); {
		if s[i] == '{' {
			var obj map[string]json.RawMessage
			dec := json.NewDecoder(strings.NewReader(s[i:]))
			if dec.Decode(&obj) == nil && obj[key] != nil {
				end := i + int(dec.InputOffset())
				return json.Unmarshal([]byte(s[i:end]), v) == nil
			}
		}
		next := strings.IndexByte(s[i:], '\n')
		if next < 0 {
			return false
		}
		i += next + 1
	}
	return false
}

// lastLines are the last n lines of s that are not blank at its end.
func lastLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, " \t\r\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

// clip is s trimmed to its first messageLines lines, with a note of how
// many more there were.
func clip(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) <= messageLines {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[:messageLines], "\n") + fmt.Sprintf("\n… %d more lines", len(lines)-messageLines)
}

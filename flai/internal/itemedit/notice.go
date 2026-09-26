package itemedit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A Notice says that someone edited an item, and what of it. It is how an
// agent hears of an edit the designer made (S-0085). The notices are a log
// under .flai-cache, outside git, like an agent's cursor: tooling never owns
// state, and losing them loses only the telling, never the edit, which is in
// the files and in git.
//
// They are not a front matter key, on purpose. Items are parsed strictly, so
// a flai that does not know a key refuses the whole item, and two flai of
// different versions do serve one repository: the installed one runs the MCP
// server and may run flai serve while a newer one is used from the tree.
type Notice struct {
	At      string   `json:"at"`
	By      string   `json:"by"`
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Title   string   `json:"title"`
	Changed []string `json:"changed"`
}

// keep is how many notices are kept; older ones are dropped when it doubles.
const keep = 500

// NoticesPath is where the notices of a repository are.
func NoticesPath(repo *workitem.Repo) string { return filepath.Join(repo.CacheDir(), "edits.jsonl") }

// note records an edit that was just made.
func note(repo *workitem.Repo, n Notice) {
	// The notice is stamped when it is written, not when the edit began: an
	// agent woken by the item's file may have looked, and moved its cursor,
	// in the second between the two.
	if at := time.Now().UTC().Format(workitem.TimeFormat); at > n.At {
		n.At = at
	}
	Record(repo, n)
}

// Record appends a notice as it is. A failure is not the edit's: it is dropped.
func Record(repo *workitem.Repo, n Notice) { appendKept(NoticesPath(repo), n, Notices(repo)) }

// appendKept appends v as one JSON line to the log at path, first cutting the
// log to its newest keep entries when it holds twice that. Failures are
// dropped: a log of notices never fails what it tells of.
func appendKept[T any](path string, v T, all []T) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return
	}
	if len(all) >= 2*keep {
		var b strings.Builder
		for _, old := range all[len(all)-keep:] {
			if line, err := json.Marshal(old); err == nil {
				b.Write(line)
				b.WriteByte('\n')
			}
		}
		_ = os.WriteFile(path, []byte(b.String()), 0o600)
	}
	line, err := json.Marshal(v)
	if err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write(append(line, '\n'))
}

// Notices returns every notice kept, oldest first.
func Notices(repo *workitem.Repo) []Notice {
	return readKept(NoticesPath(repo), func(n Notice) bool { return n.ID != "" && n.At != "" })
}

// readKept returns the entries of the log at path that parse and that ok
// accepts, oldest first.
func readKept[T any](path string, ok func(T) bool) []T {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()
	var out []T
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		var v T
		if json.Unmarshal(sc.Bytes(), &v) == nil && ok(v) {
			out = append(out, v)
		}
	}
	return out
}

// An Overlap tells an open story that an accepted story changed paths its
// claim covers (S-0132, ADR-0046), so that its agent syncs and tests against
// them before its own acceptance finds them. Overlaps are a log of their own
// beside the edit notices, not in it: a flai that predates them reads every
// line of edits.jsonl as an edit, and would tell the agent to reread its
// story instead of to sync.
type Overlap struct {
	At       string   `json:"at"`
	By       string   `json:"by"`       // who accepted
	ID       string   `json:"id"`       // the open story told
	Title    string   `json:"title"`    // its title
	Accepted string   `json:"accepted"` // the story accepted
	Paths    []string `json:"paths"`    // what it changed that the open story's claim covers
}

// OverlapsPath is where the overlap notices of a repository are.
func OverlapsPath(repo *workitem.Repo) string {
	return filepath.Join(repo.CacheDir(), "overlaps.jsonl")
}

// RecordOverlap appends an overlap notice. A failure is dropped, as for edits:
// the acceptance happened, and git has what it changed.
func RecordOverlap(repo *workitem.Repo, o Overlap) {
	appendKept(OverlapsPath(repo), o, Overlaps(repo))
}

// Overlaps returns every overlap notice kept, oldest first.
func Overlaps(repo *workitem.Repo) []Overlap {
	return readKept(OverlapsPath(repo), func(o Overlap) bool { return o.ID != "" && o.At != "" && o.Accepted != "" })
}

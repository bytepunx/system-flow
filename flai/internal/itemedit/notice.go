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
func Record(repo *workitem.Repo, n Notice) {
	path := NoticesPath(repo)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return
	}
	if all := Notices(repo); len(all) >= 2*keep {
		var b strings.Builder
		for _, old := range all[len(all)-keep:] {
			if line, err := json.Marshal(old); err == nil {
				b.Write(line)
				b.WriteByte('\n')
			}
		}
		_ = os.WriteFile(path, []byte(b.String()), 0o600)
	}
	line, err := json.Marshal(n)
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
	f, err := os.Open(NoticesPath(repo))
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()
	var out []Notice
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		var n Notice
		if json.Unmarshal(sc.Bytes(), &n) == nil && n.ID != "" && n.At != "" {
			out = append(out, n)
		}
	}
	return out
}

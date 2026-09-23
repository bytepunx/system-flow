package hostapi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// StreamActivity is one narrative: who is working on what, read from the
// file. Nothing is written for presence; an agent that stops logging ages out.
type StreamActivity struct {
	Stream     string    `json:"stream"`
	Title      string    `json:"title"`
	Agent      string    `json:"agent"`
	Session    string    `json:"session"`
	Updated    string    `json:"updated"`
	AgeSeconds int64     `json:"age_seconds"`
	Status     string    `json:"status"` // the story's state; unknown when the story is not found
	Blocked    bool      `json:"blocked"`
	Task       *TaskRef  `json:"task,omitempty"`
	LastLog    *LogEntry `json:"last_log,omitempty"`
	Path       string    `json:"path"`
}

// TaskRef names the task in progress.
type TaskRef struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// LogEntry is the last entry of a narrative's log.
type LogEntry struct {
	At   string `json:"at"`
	Text string `json:"text"`
}

// InboxEntry is one thing that needs the designer. Item and Path say where
// it lives; the dashboard makes its own link from them.
type InboxEntry struct {
	Key    string `json:"key"` // stable: the same thing keeps the same key, so new means something
	Kind   string `json:"kind"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
	Item   string `json:"item,omitempty"`
	Path   string `json:"path,omitempty"`
	At     string `json:"at,omitempty"`
}

// DesignerInbox is what needs a human (S-0042). It is not the agent's MCP
// inbox: that one is what needs an agent.
type DesignerInbox struct {
	Total   int            `json:"total"`
	Counts  map[string]int `json:"counts"`
	Entries []InboxEntry   `json:"entries"`
	Notes   []string       `json:"notes"`
}

var logHeading = regexp.MustCompile(`^### (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)`)
var sectionHeading = regexp.MustCompile(`^##\s`)

// lastLogEntry is the last entry under ## Log: its timestamp and the lines under it.
func lastLogEntry(body string) *LogEntry {
	lines := strings.Split(body, "\n")
	start := -1
	for i, l := range lines {
		if strings.TrimSpace(l) == "## Log" {
			start = i
			break
		}
	}
	var last *LogEntry
	var text []string
	for _, line := range lines[start+1:] {
		if sectionHeading.MatchString(line) {
			break
		}
		if m := logHeading.FindStringSubmatch(line); m != nil {
			if last != nil {
				last.Text = strings.TrimSpace(strings.Join(text, "\n"))
			}
			last, text = &LogEntry{At: m[1]}, nil
		} else if last != nil {
			text = append(text, line)
		}
	}
	if last != nil {
		last.Text = strings.TrimSpace(strings.Join(text, "\n"))
	}
	return last
}

// keyHash is the hash the dashboard has always keyed inbox entries with
// (djb2 over UTF-16 units, in base 36), kept so that nothing the designer
// has already seen becomes new because flai composes the inbox now.
func keyHash(s string) string {
	h := uint32(5381)
	for _, u := range utf16.Encode([]rune(s)) {
		h = (h << 5) + h + uint32(u)
	}
	return strconv.FormatUint(uint64(h), 36)
}

type narrative struct {
	name string
	fm   map[string]any
	body string
}

// narratives reads wip/agents, leaving out the index and the README.
func narratives(repo *workitem.Repo) (dir string, out []narrative) {
	abs := repo.AgentsDir()
	dir = abs
	if rel, err := filepath.Rel(repo.MainRoot, abs); err == nil {
		dir = filepath.ToSlash(rel)
	}
	entries, _ := os.ReadDir(abs)
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".md") || name == "index.md" || name == "README.md" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(abs, name))
		if err != nil {
			continue
		}
		fm, body := frontMatter(string(data))
		out = append(out, narrative{name: name, fm: fm, body: body})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return dir, out
}

func openBlock(it *workitem.Item) *workitem.Block {
	if it == nil {
		return nil
	}
	for i := range it.Blocked {
		if it.Blocked[i].Until == "" {
			return &it.Blocked[i]
		}
	}
	return nil
}

var leadingID = regexp.MustCompile(`^([EST]-\d+)`)

func peopleMethods(now func() time.Time) map[string]channel.Method {
	return map[string]channel.Method{
		// activity.get: who is working on what, from the narratives.
		"activity.get": func(_ context.Context, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
			repo, err := workitem.Open(p.Root)
			if err != nil {
				return nil, failed(err)
			}
			items, err := repo.List(false)
			if err != nil {
				return nil, failed(err)
			}
			byID := map[string]*workitem.Item{}
			for _, it := range items {
				byID[it.ID] = it
			}
			dir, list := narratives(repo)
			streams := []StreamActivity{}
			for _, n := range list {
				if n.fm == nil {
					continue
				}
				stream := str(n.fm["stream"], strings.TrimSuffix(n.name, ".md"))
				story := byID[stream]
				a := StreamActivity{Stream: stream, Title: str(n.fm["title"], ""), Agent: str(n.fm["agent"], ""), Session: str(n.fm["session"], ""),
					Updated: str(n.fm["updated"], ""), Status: "unknown", LastLog: lastLogEntry(n.body), Path: dir + "/" + n.name}
				if story != nil {
					a.Status = story.Status
					if a.Title == "" {
						a.Title = story.Title
					}
					a.Blocked = openBlock(story) != nil
				}
				for _, it := range items {
					if it.Type != workitem.Task || it.Parent != stream {
						continue
					}
					a.Blocked = a.Blocked || openBlock(it) != nil
					if a.Task == nil && it.Status == workitem.InProgress {
						a.Task = &TaskRef{ID: it.ID, Title: it.Title}
					}
				}
				if t, err := time.Parse(time.RFC3339, a.Updated); err == nil {
					if age := int64(now().Sub(t).Round(time.Second).Seconds()); age > 0 {
						a.AgeSeconds = age
					}
				}
				streams = append(streams, a)
			}
			sort.SliceStable(streams, func(i, j int) bool {
				if streams[i].Updated != streams[j].Updated {
					return streams[i].Updated > streams[j].Updated
				}
				return streams[i].Stream < streams[j].Stream
			})
			return map[string]any{"streams": streams}, nil
		},

		// inbox.designer: what needs a human. Not the agent's MCP inbox.
		"inbox.designer": func(_ context.Context, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
			repo, err := workitem.Open(p.Root)
			if err != nil {
				return nil, failed(err)
			}
			who := repo.Manifest.Owner
			if who == "" {
				who = "designer"
			}
			out := DesignerInbox{Counts: map[string]int{"thread": 0, "question": 0, "review": 0, "blocked": 0, "overlap": 0}, Entries: []InboxEntry{}, Notes: []string{}}
			add := func(e InboxEntry) {
				out.Entries = append(out.Entries, e)
				out.Counts[e.Kind]++
			}

			all, err := threads.List(repo)
			if err != nil {
				return nil, failed(err)
			}
			for _, th := range all {
				entries := th.Entries()
				if !th.Open() || len(entries) == 0 || entries[len(entries)-1].Author == who {
					continue
				}
				last := entries[len(entries)-1]
				on := th.Anchor.Item
				if on == "" {
					on = th.Anchor.Path
				}
				detail := last.Author + " wrote last, on " + on
				if th.Anchor.Heading != "" {
					detail += " (" + th.Anchor.Heading + ")"
				}
				add(InboxEntry{Key: "thread:" + th.ID, Kind: "thread", Title: th.Title, Detail: detail, Item: th.Anchor.Item, Path: th.Anchor.Path, At: last.At})
			}

			items, err := repo.List(false)
			if err != nil {
				return nil, failed(err)
			}
			byID := map[string]*workitem.Item{}
			for _, it := range items {
				byID[it.ID] = it
			}

			// A question lingers in a narrative until whoever wrote it removes
			// it by hand (it is not a thread, so nothing marks it answered);
			// once the story it belongs to is in review, done, or cancelled, or
			// gone from this list (archived), it no longer needs the designer,
			// so it is left out here whatever the narrative still says.
			dir, list := narratives(repo)
			for _, n := range list {
				stream := str(n.fm["stream"], strings.TrimSuffix(n.name, ".md"))
				story := byID[stream]
				if story == nil || story.Status == workitem.Review || story.Status == workitem.Done || story.Status == workitem.Cancelled {
					continue
				}
				for _, q := range workitem.OpenQuestions(n.body) {
					add(InboxEntry{Key: "question:" + stream + ":" + keyHash(q), Kind: "question", Title: q,
						Detail: "asked in the narrative of " + stream, Item: stream, Path: dir + "/" + n.name, At: str(n.fm["updated"], "")})
				}
			}

			for _, it := range items {
				if it.Type == workitem.Story && it.Status == workitem.Review {
					at := ""
					if n := len(it.Transitions); n > 0 {
						at = it.Transitions[n-1].At
					}
					add(InboxEntry{Key: "review:" + it.ID, Kind: "review", Title: it.ID + " " + it.Title, Detail: "in review: accept it or send it back", Item: it.ID, At: at})
				}
				if b := openBlock(it); b != nil && !it.Closed() {
					add(InboxEntry{Key: "blocked:" + it.ID + ":" + b.From, Kind: "blocked", Title: it.ID + " " + it.Title, Detail: "blocked: " + b.Reason, Item: it.ID, At: b.From})
				}
			}

			// Overlapping touches are the check's rule, which has one implementation.
			res, err := check.Run(repo, now())
			if err != nil {
				out.Notes = append(out.Notes, "Overlapping touches are not listed: "+err.Error())
			} else {
				for _, f := range res.Findings {
					if f.Rule != "wip.overlap" {
						continue
					}
					e := InboxEntry{Key: "overlap:" + keyHash(f.Message), Kind: "overlap", Title: f.Message}
					if m := leadingID.FindStringSubmatch(f.Message); m != nil {
						e.Item = m[1]
					}
					add(e)
				}
			}
			out.Total = len(out.Entries)
			return out, nil
		},
	}
}

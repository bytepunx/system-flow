package hostapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The dashboard's writes (S-0075). Each method validates its arguments as
// data, builds the command line itself, and runs flai in the project's
// folder, so the rules stay where they are and are tested: in the commands.
// Nothing the dashboard sends is ever a flag or a command: values are given
// as --flag=value or after --, and IDs, states, and paths are checked first.

// Codes for what a caller can act on, beside the JSON-RPC ones.
const (
	Conflict = -32009 // the thing changed since it was read; Data has the current version
	Refused  = -32010 // flai check refused the content; Data has the findings
	Rule     = -32011 // a workflow rule said no; the message says which
	Disabled = -32012 // a host action the operator has not enabled; Data says what enables it
)

// Trailer marks a commit made for the dashboard, as it always has.
const Trailer = "Co-Authored-By: flaiover <flaiover@localhost>"

// Run is one flai invocation: what to run, what to give it on standard
// input, and where each log event it writes goes while it runs.
type Run struct {
	Dir     string
	Args    []string
	Stdin   string
	OnEvent func(map[string]any)
}

// Ran is what came back.
type Ran struct {
	Stdout []byte
	Events []map[string]any // the log events from standard error
	Exit   int
}

// Runner runs flai. The real one runs this executable; tests give their own.
type Runner func(ctx context.Context, r Run) (Ran, error)

// ExecRunner runs the flai that is running, as the dashboard would have run
// the one in its container: JSON out, log events on standard error.
func ExecRunner(ctx context.Context, r Run) (Ran, error) {
	exe, err := os.Executable()
	if err != nil {
		return Ran{}, err
	}
	cmd := exec.CommandContext(ctx, exe, r.Args...)
	cmd.Dir = r.Dir
	cmd.Env = append(os.Environ(), "FLAI_AGENT=flaiover", "FLAI_SESSION=", "LOG_FORMAT=json")
	cmd.Stdin = strings.NewReader(r.Stdin)
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 10 * time.Second
	var out bytes.Buffer
	cmd.Stdout = &out
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return Ran{}, err
	}
	if err := cmd.Start(); err != nil {
		return Ran{}, err
	}
	var ran Ran
	sc := bufio.NewScanner(stderr)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for sc.Scan() {
		var ev map[string]any
		if json.Unmarshal(sc.Bytes(), &ev) != nil {
			ev = map[string]any{"level": "INFO", "msg": sc.Text()}
		}
		ran.Events = append(ran.Events, ev)
		if r.OnEvent != nil {
			r.OnEvent(ev)
		}
	}
	err = cmd.Wait()
	ran.Stdout = out.Bytes()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		ran.Exit = exit.ExitCode()
		return ran, nil
	}
	return ran, err
}

// Written is what a command answered: its JSON, and the warnings it logged.
type Written struct {
	Data     json.RawMessage `json:"data"`
	Warnings []string        `json:"warnings"`
}

// journal remembers the writes that completed, by the request ID the
// dashboard gave, so that a request repeated after a lost connection is
// answered with what happened and not done twice.
type journal struct {
	mu   sync.Mutex
	done map[string]journalled
	now  func() time.Time
}

type journalled struct {
	at  time.Time
	res any
	err *channel.Error
}

const journalKeeps = 10 * time.Minute

func (j *journal) get(key string) (journalled, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	for k, v := range j.done {
		if j.now().Sub(v.at) > journalKeeps {
			delete(j.done, k)
		}
	}
	v, ok := j.done[key]
	return v, ok
}

func (j *journal) put(key string, res any, err *channel.Error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.done[key] = journalled{at: j.now(), res: res, err: err}
}

// Request is the record of one detached write, kept by its request ID from
// before it acts, so that a repeat is answered and not done again even when
// the write ended the flai serve that took it (S-0109).
type Request struct {
	Started time.Time      `json:"started"`
	PID     int            `json:"pid"`   // the flai serve that took it
	Until   time.Time      `json:"until"` // forgotten after this
	Done    bool           `json:"done"`
	Result  *Written       `json:"result,omitempty"`
	Err     *channel.Error `json:"error,omitempty"`
}

// Requests changes the records of detached writes under a lock, and keeps
// what change leaves.
type Requests func(change func(map[string]Request)) error

// inMemory keeps the records for as long as the table lives.
func inMemory() Requests {
	var mu sync.Mutex
	kept := map[string]Request{}
	return func(change func(map[string]Request)) error {
		mu.Lock()
		defer mu.Unlock()
		change(kept)
		return nil
	}
}

// begin records a detached write as started, unless its key is recorded
// already: then it returns that record, and true. Records past their time
// are forgotten first.
func (rs Requests) begin(key string, now time.Time, keep time.Duration) (was Request, found bool, err error) {
	err = rs(func(all map[string]Request) {
		forgetOld(all, now)
		if was, found = all[key]; !found {
			all[key] = Request{Started: now, PID: os.Getpid(), Until: now.Add(keep)}
		}
	})
	return was, found, err
}

// finish records what a detached write answered.
func (rs Requests) finish(key string, now time.Time, res any, rerr *channel.Error) error {
	return rs(func(all map[string]Request) {
		forgetOld(all, now)
		r := all[key]
		r.Done, r.Until, r.Err = true, now.Add(journalKeeps), rerr
		if w, ok := res.(Written); ok {
			r.Result = &w
		}
		all[key] = r
	})
}

// unrecorded adds to an answer that its outcome could not be recorded, so a
// repeat will not be told it.
func unrecorded(res any, rerr *channel.Error, err error) (any, *channel.Error) {
	say := "the outcome could not be recorded, so a repeat of this request will not be told it: " + err.Error()
	if w, ok := res.(Written); ok {
		w.Warnings = append(w.Warnings, say)
		return w, rerr
	}
	if rerr != nil {
		rerr.Message += "; " + say
	}
	return res, rerr
}

func forgetOld(all map[string]Request, now time.Time) {
	for k, r := range all {
		if now.After(r.Until) {
			delete(all, k)
		}
	}
}

// repeated answers a repeat of a detached write from its record: what it
// answered, once it has; until then, that it is under way in this flai serve,
// or that the flai serve that took it ended before it could say, which is
// what a restart of flai serve itself looks like. It is never done again.
func repeated(id string, r Request) (any, *channel.Error) {
	if r.Done {
		if r.Err != nil || r.Result == nil {
			return nil, r.Err
		}
		return *r.Result, nil
	}
	started := r.Started.UTC().Format(time.RFC3339)
	underWay := r.PID == os.Getpid()
	say := fmt.Sprintf("request %s was started at %s and is still under way; it was not started again", id, started)
	if !underWay {
		say = fmt.Sprintf("request %s was started at %s by a flai serve (pid %d) that ended before it could say how it went; it was not done again", id, started, r.PID)
	}
	data, _ := json.Marshal(map[string]any{"repeat": map[string]any{"request_id": id, "started": started, "under_way": underWay}})
	return Written{Data: data, Warnings: []string{say}}, nil
}

// contentHash is what docedit.Hash gives: a SHA-256 in hex.
var contentHash = regexp.MustCompile(`^[0-9a-f]{64}$`)

func isNature(v string) bool {
	for _, n := range workitem.Natures {
		if n == v {
			return true
		}
	}
	return false
}

var requestID = regexp.MustCompile(`^[A-Za-z0-9._-]{8,64}$`)

// ActionPush is the host action that pushes an acceptance and publishes the
// template with the operator's own credentials (S-0078).
const ActionPush = "push"

// ActionAgent is the host action that starts a story's agent when the story
// becomes ready and the in-progress limit has room (S-0079, S-0104,
// ADR-0043). flai serve performs that itself, from what it sees in the
// project's files. Two methods ask for it: agent.restart, a new agent for a
// story whose agent dropped or failed (S-0116), and agent.start, a ready
// story's agent now (S-0115).
const ActionAgent = "agent"

// ActionDashboard is the host action that restarts, upgrades, or stops the
// dashboard container with Docker on this host (S-0081).
const ActionDashboard = "dashboard"

// ActionChecks is the host action that runs the commands named in the
// manifest or the host's configuration, in a story's worktree, and reports
// the outcome (S-0082).
const ActionChecks = "checks"

// ActionHost is the host action that has flai host start, stop, or restart
// flai serve and the MCP servers, and install the newest flai and restart on
// it (S-0106).
const ActionHost = "host"

// ActionSettings is the host action that lets a dashboard change the host's
// settings for a project (S-0105): the other host actions, the agent and its
// harnesses, the checks, the import folders, and the tokens. Those kept for
// every project need it enabled for every project. It is turned on and off
// in a shell on the host alone: no method changes it.
const ActionSettings = "settings"

// Actions are the host actions there are, with what each lets a dashboard do.
var Actions = map[string]string{
	ActionPush:      "push accepted work, and publish everything merged and unreleased since each component's last tag, with your git credentials; a holder of the dashboard token can then publish any story that is in review and any release accumulated since",
	ActionAgent:     "start each story's agent, with the harnesses and the command you set with flai serve agent, on this machine and as you, whenever a story becomes ready and the in-progress limit has room, start a ready story's agent on demand, and start or queue a new one for a story whose agent dropped or failed; whoever can move a story to ready or press Start agent or Retry, a holder of the dashboard token included, then starts it",
	ActionDashboard: "restart the dashboard container, upgrade it to the image your configuration names, or stop it, with Docker on this host; an upgrade is never applied until the new image answers healthy, so a bad one leaves the running container untouched",
	ActionChecks:    "run the commands named in flai serve checks set or the manifest's checks:, in a story's worktree, on this host, and cancel a run; whoever can open the review page then decides what runs there",
	ActionHost:      "have flai host start, stop, or restart flai serve and the MCP servers of every project on this host, and download the newest flai release with your GitHub credentials, install it over the flai on this host, and restart everything on it",
	ActionSettings:  "change this project's host settings: turn the other host actions on and off, set its default agent, and rotate its MCP token; enabled for every project, also the agent's command, the harnesses, the checks, the import folders, and the dashboard token. A holder of the dashboard token can then run any command on this host, as you; only a shell turns this off",
}

// Host is what the host decides and records about host actions (ADR-0029).
// The zero value enables nothing and records nothing.
type Host struct {
	// Enabled reports whether the operator enabled an action for a project.
	Enabled func(action, root string) bool
	// Record writes one entry of the host's journal.
	Record func(Entry)
	// Agent reports what flai serve knows of agents it starts for a project
	// (S-0079): what runs, the last start or failure, and why a ready story
	// waits. Nil when nothing starts agents here.
	Agent func(root string) any
	// Settings reports the host's settings as they apply to a project
	// (S-0105), for the dashboard's settings page. Nil when there are none.
	Settings func(root string) any
	// Requests keeps the records of detached writes where a restart of flai
	// serve does not lose them (S-0109). Nil keeps them in memory.
	Requests Requests
}

// Entry is one host action asked for, whatever became of it.
type Entry struct {
	At        string `json:"at"`
	Action    string `json:"action"`
	Method    string `json:"method"`
	Project   string `json:"project"`
	Root      string `json:"root"`
	By        string `json:"by"` // who the dashboard acts for: the manifest's owner
	RequestID string `json:"request_id,omitempty"`
	Outcome   string `json:"outcome"` // done, failed, or disabled
	Detail    string `json:"detail,omitempty"`
}

func (h Host) enabled(action, root string) bool {
	return h.Enabled != nil && h.Enabled(action, root)
}

// EnableCommand is what the operator runs, on the host, in the project.
func EnableCommand(action string) string { return "flai serve enable " + action }

// spec is one write: how its params become a command line.
type spec struct {
	// action names the host action this method is: asked for while the
	// operator has not enabled it, it is refused and says what enables it.
	action string
	// uses names a host action this method performs as part of its work when
	// it is enabled, and does without when it is not; it is journalled then.
	uses string
	// hostwide marks a setting kept for every project (S-0105): it needs its
	// action enabled for every project, not this one only.
	hostwide bool
	// record journals every call under this name without gating it: an act
	// on the host the operator consented to some other way (S-0098: import,
	// by naming the folder the repository is in).
	record string
	// build validates and returns the arguments (without --json) and what goes on standard input.
	build func(p channel.Project, raw json.RawMessage) (args []string, stdin string, err *channel.Error)
	// exits maps exit codes that carry a payload on standard output to error codes.
	exits map[int]int
	// progress sends each log event to the dashboard while the command runs.
	progress bool
	// reads marks a command that changes nothing: no request ID is asked for.
	reads bool
	// describe says in a line what a successful call did, for the journal.
	// Nil means the push/publish shape (describe, below); a host action with
	// a different answer shape (dashboard.restart, dashboard.upgrade,
	// dashboard.stop) sets its own. A failed call is always "failed",
	// err.Message, whichever this is.
	describe func(res any, err *channel.Error) (outcome, detail string)
	// say, when set, is the journal's line for a call that succeeded, made
	// from what was asked (S-0105: which setting became what).
	say func(raw json.RawMessage) string
	// detachTimeout, when nonzero, runs the command with its own background
	// context instead of the request's: a write whose own success can close
	// the WebSocket connection the request arrived on (dashboard.restart
	// stops the container answering it; dashboard.upgrade does on a
	// successful swap; dashboard.stop does when it is the last project)
	// must not be killed by that connection's own context cancelling out
	// from under it (S-0081, found live: exec.CommandContext SIGKILLed the
	// subprocess mid-restart when the container's own disconnect cancelled
	// the request context that spawned it, leaving no container running at
	// all). Progress still tries the original context; sending on a closed
	// connection just fails silently (channel.Progress → send → conn.Write).
	detachTimeout time.Duration
}

func text(v string) string { return strings.Join(strings.Fields(v), " ") }

// owner is who the dashboard acts for: the manifest's owner, decided here.
// The dashboard does not get to name them.
func owner(p channel.Project) string {
	if m, err := manifest.Load(filepath.Join(p.Root, manifest.File)); err == nil && m.Owner != "" {
		return m.Owner
	}
	return "designer"
}

var (
	anyItemID  = regexp.MustCompile(`^[EST]-\d{1,6}$`)
	threadID   = regexp.MustCompile(`^(?i:TH-)?\d{1,6}$|^TH-\d{1,6}$`)
	adrNumber  = regexp.MustCompile(`^(?:ADR-)?(\d{1,4})$`)
	listValue  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/@+-]*$`)
	sinceValue = regexp.MustCompile(`^\d+[dwh]$`)
	docHash    = regexp.MustCompile(`^[A-Fa-f0-9]{8,128}$`)
)

func needID(id string) *channel.Error {
	if !anyItemID.MatchString(id) {
		return bad("%q is not an item ID", id)
	}
	return nil
}

func adrNo(v string) (string, *channel.Error) {
	m := adrNumber.FindStringSubmatch(strings.TrimSpace(v))
	if m == nil {
		return "", bad("%s is not an ADR", v)
	}
	n, _ := strconv.Atoi(m[1])
	if n == 0 {
		return "", bad("%s is not an ADR", v)
	}
	return strconv.Itoa(n), nil
}

// docRel checks a document path without requiring the file to exist (a save
// may create it): Markdown, inside the manifest's folders, no climbing.
func docRel(p channel.Project, rel string) (string, *channel.Error) {
	clean := filepath.ToSlash(filepath.Clean(rel))
	if rel == "" || filepath.IsAbs(rel) || strings.HasPrefix(clean, "..") || strings.HasPrefix(clean, "-") || strings.ContainsRune(rel, 0) {
		return "", bad("%q is not a path inside the repository", rel)
	}
	if !strings.HasSuffix(clean, ".md") {
		return "", bad("only Markdown files are edited")
	}
	dirs, err := layoutDirs(p.Root)
	if err != nil {
		return "", failed(err)
	}
	for _, d := range dirs {
		if strings.HasPrefix(clean, d+"/") {
			return clean, nil
		}
	}
	return "", bad("%s is not under the project's design, docs, or wip folder", clean)
}

// decode reads params into v, refusing what is not an object.
func decode[T any](raw json.RawMessage) (T, *channel.Error) {
	var v T
	return v, params(raw, &v)
}

// specs are every write, the settings among them (S-0105).
func specs() map[string]spec {
	table := itemSpecs()
	for name, sp := range settingsSpecs() {
		table[name] = sp
	}
	return table
}

func itemSpecs() map[string]spec {
	type build = func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error)
	one := func(b build) spec { return spec{build: b} }
	read := func(b build) spec { return spec{build: b, reads: true} }

	moveArgs := func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		in, e := decode[struct {
			ID                 string `json:"id"`
			To                 string `json:"to"`
			Reason             string `json:"reason"`
			IncludeUncommitted bool   `json:"include_uncommitted"`
			DryRun             bool   `json:"dry_run"`
		}](raw)
		if e != nil {
			return nil, "", e
		}
		if e := needID(in.ID); e != nil {
			return nil, "", e
		}
		if !isState(in.To) {
			return nil, "", bad("%q is not a state", in.To)
		}
		args := []string{"move", in.ID, in.To, "--by=" + owner(p)}
		preview := in.DryRun && in.To == workitem.Cancelled
		switch {
		case text(in.Reason) != "":
			args = append(args, "--reason="+text(in.Reason))
		case preview:
			args = append(args, "--reason=preview") // flai validates the move before it lists anything
		}
		if preview {
			args = append(args, "--dry-run")
		}
		// --yes is the designer's choice, in the acceptance confirmation, to include uncommitted files.
		if in.IncludeUncommitted && in.To == workitem.Done {
			args = append(args, "--yes")
		}
		return args, "", nil
	}

	return map[string]spec{
		"item.move": {build: moveArgs},
		// A move that only previews a cancellation changes nothing and needs no request ID.
		"item.move.preview": read(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			b, _ := json.Marshal(map[string]any{"id": in.ID, "to": workitem.Cancelled, "dry_run": true})
			return moveArgs(p, b)
		}),

		"item.order": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Before string `json:"before"`
				After  string `json:"after"`
				Top    bool   `json:"top"`
				Bottom bool   `json:"bottom"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			given := 0
			for _, on := range []bool{in.Before != "", in.After != "", in.Top, in.Bottom} {
				if on {
					given++
				}
			}
			if given != 1 {
				return nil, "", bad("say where it goes with exactly one of before, after, top, bottom")
			}
			for _, id := range []string{in.ID, in.Before + in.After} {
				if id != "" {
					if e := needID(id); e != nil {
						return nil, "", e
					}
				}
			}
			switch {
			case in.Before != "":
				return []string{"order", in.ID, "--before=" + in.Before}, "", nil
			case in.After != "":
				return []string{"order", in.ID, "--after=" + in.After}, "", nil
			case in.Top:
				return []string{"order", in.ID, "--top"}, "", nil
			}
			return []string{"order", in.ID, "--bottom"}, "", nil
		}),

		"item.block": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Reason string `json:"reason"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if text(in.Reason) == "" {
				return nil, "", bad("a reason is required")
			}
			return []string{"block", in.ID, "--reason=" + text(in.Reason)}, "", nil
		}),

		"item.unblock": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"unblock", in.ID}, "", nil
		}),

		"item.new": {exits: map[int]int{4: Refused}, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Type    string          `json:"type"`
				Title   string          `json:"title"`
				Nature  string          `json:"nature"`
				Parent  string          `json:"parent"`
				Tags    []string        `json:"tags"`
				Touches []string        `json:"touches"`
				Agent   *manifest.Agent `json:"agent"`
				Body    string          `json:"body"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Type != workitem.Epic && in.Type != workitem.Story {
				return nil, "", bad("type must be epic or story; tasks are written by the agent that pulls a story")
			}
			title := text(in.Title)
			if title == "" {
				return nil, "", bad("a title is required")
			}
			if len([]rune(title)) > 200 {
				return nil, "", bad("the title is longer than 200 characters")
			}
			nature := in.Nature
			if nature == "" {
				nature = "feature"
			}
			ok := false
			for _, n := range workitem.Natures {
				ok = ok || n == nature
			}
			if !ok {
				return nil, "", bad("nature must be one of %s", strings.Join(workitem.Natures, ", "))
			}
			if strings.TrimSpace(in.Body) == "" {
				return nil, "", bad("write what the item is for: the body is empty")
			}
			args := []string{in.Type, "new", "--nature=" + nature, "--owner=" + owner(p)}
			switch {
			case in.Type == workitem.Story && in.Parent == "":
				// no epic (S-0092): not every story belongs to one
			case in.Type == workitem.Story:
				if !anyItemID.MatchString(in.Parent) || !strings.HasPrefix(in.Parent, "E-") {
					return nil, "", bad("a story's parent must be an epic")
				}
				args = append(args, "--epic="+in.Parent)
			case in.Parent != "":
				return nil, "", bad("an epic has no parent")
			}
			for flag, values := range map[string][]string{"tag": in.Tags, "touches": in.Touches} {
				for _, v := range values {
					if !listValue.MatchString(strings.TrimSpace(v)) {
						return nil, "", bad("%s values are single words or paths without commas", flag)
					}
				}
			}
			for _, v := range in.Tags {
				args = append(args, "--tag="+strings.TrimSpace(v))
			}
			for _, v := range in.Touches {
				args = append(args, "--touches="+strings.TrimSpace(v))
			}
			if !in.Agent.IsZero() {
				if in.Type != workitem.Story {
					return nil, "", bad("only a story carries an agent")
				}
				agent, e := agentArgs(in.Agent)
				if e != nil {
					return nil, "", e
				}
				args = append(args, agent...)
			}
			return append(args, "--body-stdin", "--autocommit", "--trailer="+Trailer, "--", title), in.Body, nil
		}},

		// item.show and item.edit: an item's own words, changed after it was made
		// (S-0085). flai edit does the work: the retitle kept in step everywhere,
		// the parents' lists, the check with the change in place, one commit.
		"item.show": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"edit", in.ID, "--show"}, "", nil
		}),

		"item.edit": {exits: map[int]int{3: Conflict, 4: Refused}, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID      string    `json:"id"`
				Hash    string    `json:"hash"`
				Title   *string   `json:"title"`
				Nature  *string   `json:"nature"`
				Tags    *[]string `json:"tags"`
				Touches *[]string `json:"touches"`
				Parent  *string   `json:"parent"`
				Body    *string   `json:"body"`
				// Agent replaces the story's agent: absent leaves it, null (or an
				// empty object) removes it (S-0103)
				Agent json.RawMessage `json:"agent"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			// the dashboard always edits what it read: without the hash a change an
			// agent made meanwhile would be overwritten unseen
			if !contentHash.MatchString(in.Hash) {
				return nil, "", bad("hash is required: the one item.show gave for what was edited")
			}
			args := []string{"edit", in.ID, "--hash=" + in.Hash, "--by=" + owner(p), "--autocommit", "--trailer=" + Trailer}
			if in.Title != nil {
				if text(*in.Title) == "" {
					return nil, "", bad("a title is required")
				}
				args = append(args, "--title="+text(*in.Title))
			}
			if in.Nature != nil {
				if !isNature(*in.Nature) {
					return nil, "", bad("%q is not a nature", *in.Nature)
				}
				args = append(args, "--nature="+*in.Nature)
			}
			if in.Parent != nil {
				if e := needID(*in.Parent); e != nil {
					return nil, "", e
				}
				args = append(args, "--parent="+*in.Parent)
			}
			for _, l := range []struct {
				flag, clear string
				values      *[]string
			}{{"tag", "--clear-tags", in.Tags}, {"touches", "--clear-touches", in.Touches}} {
				if l.values == nil {
					continue
				}
				if len(*l.values) == 0 {
					args = append(args, l.clear)
					continue
				}
				for _, v := range *l.values {
					if !listValue.MatchString(strings.TrimSpace(v)) {
						return nil, "", bad("%s values are single words or paths without commas", l.flag)
					}
					args = append(args, "--"+l.flag+"="+strings.TrimSpace(v))
				}
			}
			if len(in.Agent) > 0 {
				var agent *manifest.Agent
				if err := json.Unmarshal(in.Agent, &agent); err != nil {
					return nil, "", bad("agent is an object with harness, model, and config, or null: %s", err)
				}
				more, e := agentArgs(agent)
				if e != nil {
					return nil, "", e
				}
				args = append(append(args, "--clear-agent"), more...)
			}
			stdin := ""
			if in.Body != nil {
				if strings.TrimSpace(*in.Body) == "" {
					return nil, "", bad("write what the item is for: the body is empty")
				}
				args, stdin = append(args, "--body-stdin"), *in.Body
			}
			if len(args) == 6 {
				return nil, "", bad("nothing to change")
			}
			return args, stdin, nil
		}},

		"item.template": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Type string `json:"type"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Type != workitem.Epic && in.Type != workitem.Story {
				return nil, "", bad("type must be epic or story")
			}
			return []string{in.Type, "new", "--print-body"}, "", nil
		}),

		"accept.preview": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"accept", in.ID, "--dry-run"}, "", nil
		}),

		"accept.run": {progress: true, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID                 string `json:"id"`
				IncludeUncommitted bool   `json:"include_uncommitted"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			// Acceptance merges and archives only; it releases nothing (S-0087).
			// Publishing what has accumulated is a step of its own (publish.run).
			args := []string{"accept", in.ID, "--by=" + owner(p)}
			if in.IncludeUncommitted {
				args = append(args, "--yes")
			}
			return args, "", nil
		}},

		"stream.diff": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"stream", "diff", in.ID}, "", nil
		}),

		"stream.log": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID    string `json:"id"`
				Entry string `json:"entry"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if strings.TrimSpace(in.Entry) == "" {
				return nil, "", bad("an entry is required")
			}
			return []string{"stream", "log", in.ID, "--", strings.TrimSpace(in.Entry)}, "", nil
		}),

		// stream.answer answers a hand-written open question (S-0090): one
		// with no thread of its own, so thread.reply does not already reach
		// it. It moves to Decisions (design/system/agent-narrative.md).
		"stream.answer": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID       string `json:"id"`
				Question string `json:"question"`
				Answer   string `json:"answer"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if strings.TrimSpace(in.Question) == "" {
				return nil, "", bad("question is required")
			}
			if strings.TrimSpace(in.Answer) == "" {
				return nil, "", bad("answer is required")
			}
			return []string{"stream", "answer", in.ID, "--by=" + owner(p), "--", strings.TrimSpace(in.Question), strings.TrimSpace(in.Answer)}, "", nil
		}),

		"thread.new": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				On      string `json:"on"`
				Heading string `json:"heading"`
				Title   string `json:"title"`
				Text    string `json:"text"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.On == "" || text(in.Title) == "" || strings.TrimSpace(in.Text) == "" {
				return nil, "", bad("on, title, and text are required")
			}
			if !anyItemID.MatchString(in.On) {
				if _, e := docRel(p, in.On); e != nil {
					return nil, "", bad("a thread is opened on an item or on a document of the project")
				}
			}
			args := []string{"thread", "new", "--on=" + in.On, "--by=" + owner(p)}
			if in.Heading != "" {
				args = append(args, "--heading="+text(in.Heading))
			}
			return append(args, "--", text(in.Title), strings.TrimSpace(in.Text)), "", nil
		}),

		"thread.reply": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !threadID.MatchString(in.ID) {
				return nil, "", bad("%q is not a thread", in.ID)
			}
			if strings.TrimSpace(in.Text) == "" {
				return nil, "", bad("text is required")
			}
			return []string{"thread", "reply", in.ID, "--by=" + owner(p), "--", strings.TrimSpace(in.Text)}, "", nil
		}),

		"thread.resolve": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Reason string `json:"reason"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !threadID.MatchString(in.ID) {
				return nil, "", bad("%q is not a thread", in.ID)
			}
			args := []string{"thread", "resolve", in.ID, "--by=" + owner(p)}
			if text(in.Reason) != "" {
				args = append(args, "--reason="+text(in.Reason))
			}
			return args, "", nil
		}),

		"doc.show": read(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Path string `json:"path"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			rel, e := docRel(p, in.Path)
			if e != nil {
				return nil, "", e
			}
			return []string{"doc", "show", "--", rel}, "", nil
		}),

		"doc.save": {exits: map[int]int{3: Conflict, 4: Refused}, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Path    string  `json:"path"`
				Content *string `json:"content"`
				Hash    string  `json:"hash"`
				Message string  `json:"message"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Path == "" || in.Content == nil || in.Hash == "" {
				return nil, "", bad("path, content, and hash are required")
			}
			rel, e := docRel(p, in.Path)
			if e != nil {
				return nil, "", e
			}
			if !docHash.MatchString(in.Hash) {
				return nil, "", bad("hash is the one the document was loaded with")
			}
			args := []string{"doc", "save", "--hash=" + in.Hash, "--trailer=" + Trailer}
			if text(in.Message) != "" {
				args = append(args, "--message="+text(in.Message))
			}
			return append(args, "--", rel), *in.Content, nil
		}},

		"adr.new": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Title      string   `json:"title"`
				Status     string   `json:"status"`
				Supersedes []string `json:"supersedes"`
				Refines    []string `json:"refines"`
				Body       string   `json:"body"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			title := text(in.Title)
			if title == "" {
				return nil, "", bad("a title is required: the decision, as a sentence")
			}
			if len([]rune(title)) > 200 {
				return nil, "", bad("the title is longer than 200 characters")
			}
			status := in.Status
			if status == "" {
				status = "proposed"
			}
			if status != "proposed" && status != "accepted" {
				return nil, "", bad("status must be proposed or accepted")
			}
			if strings.TrimSpace(in.Body) == "" {
				return nil, "", bad("write the decision: the body is empty")
			}
			args := []string{"adr", "new", "--status=" + status}
			for _, rel := range []struct {
				flag string
				list []string
			}{{"supersedes", in.Supersedes}, {"refines", in.Refines}} {
				for _, v := range rel.list {
					n, e := adrNo(v)
					if e != nil {
						return nil, "", e
					}
					args = append(args, "--"+rel.flag+"="+n)
				}
			}
			return append(args, "--body-stdin", "--autocommit", "--trailer="+Trailer, "--", title), in.Body, nil
		}},

		"adr.template": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"adr", "new", "--print-body"}, "", nil
		}),

		"adr.accept": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			n, e := adrNo(in.ID)
			if e != nil {
				return nil, "", e
			}
			return []string{"adr", "accept", n, "--autocommit", "--trailer=" + Trailer}, "", nil
		}},

		"stats.get": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Since string `json:"since"`
				Type  string `json:"type"`
				By    string `json:"by"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			args := []string{"stats"}
			if in.Since != "" {
				if !sinceValue.MatchString(in.Since) {
					return nil, "", bad("since is a number and d, w, or h")
				}
				args = append(args, "--since="+in.Since)
			}
			if in.Type != "" {
				if !itemType[in.Type] {
					return nil, "", bad("type must be epic, story, or task")
				}
				args = append(args, "--type="+in.Type)
			}
			if in.By != "" {
				if in.By != "nature" && in.By != "type" && in.By != "parent" {
					return nil, "", bad("by must be nature, type, or parent")
				}
				args = append(args, "--by="+in.By)
			}
			return args, "", nil
		}),

		// exit 3 is a remote that has moved: said as a conflict, with its reason (S-0078)
		"push.pending": {reads: true, exits: map[int]int{3: Conflict}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"push", "--pending", "--dry-run"}, "", nil
		}},

		// push.run: the host action. flai push --pending --publish as the
		// operator: the branch and the release tags of accepted work, then the
		// template where its publish remote is behind. Never forced; when the
		// remote has moved it refuses (exit 3) and says to fetch and merge.
		"push.run": {action: ActionPush, exits: map[int]int{3: Conflict}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"push", "--pending", "--publish"}, "", nil
		}},

		// publish.preview: everything release.Pending would release, without
		// changing anything, so the board can show it before the operator asks
		// for it (S-0087).
		"publish.preview": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"release", "--pending", "--dry-run"}, "", nil
		}),

		// publish.run: the host action, the same one push.run uses (ADR-0031,
		// S-0078): flai release --pending as the operator, applying, tagging,
		// and pushing everything accepted and unreleased. Never forced; when
		// the remote has moved it refuses (exit 3) and says to fetch and merge.
		"publish.run": {action: ActionPush, exits: map[int]int{3: Conflict}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"release", "--pending"}, "", nil
		}},

		// dashboard.status: a read of what flai dashboard status already
		// reports, so the board can show the version running and who else it
		// serves without a host action (S-0081).
		"dashboard.status": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "status"}, "", nil
		}),

		// dashboard.check: a read that pulls the configured image and
		// compares it to what is running, changing nothing (S-0081).
		"dashboard.check": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "check"}, "", nil
		}),

		// dashboard.restart and dashboard.upgrade: the host action (S-0081),
		// running the same commands the operator's own shell does. Both
		// stream progress: they take real Docker time, unlike a move or an
		// edit. Neither takes an image or tag from the dashboard; flai
		// upgrades only to what this host's own configuration names.
		"dashboard.restart": {action: ActionDashboard, progress: true, describe: describeDashboardRestart, detachTimeout: 90 * time.Second, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "restart"}, "", nil
		}},
		"dashboard.upgrade": {action: ActionDashboard, progress: true, describe: describeDashboardUpgrade, detachTimeout: 6 * time.Minute, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "upgrade"}, "", nil
		}},

		// dashboard.stop: the host action, the same command the operator's
		// own flai dashboard stop runs: unregisters this project, and stops
		// the container only when it was the last one registered.
		"dashboard.stop": {action: ActionDashboard, describe: describeDashboardStop, detachTimeout: 60 * time.Second, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"dashboard", "stop"}, "", nil
		}},

		// host.status and host.check: reads of flai host (S-0106): the host,
		// and each process it keeps; whether a newer flai is published.
		"host.status": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"host", "status"}, "", nil
		}),
		"host.check": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"host", "check"}, "", nil
		}),

		// host.process and host.upgrade: the host action (S-0106). Both are
		// detached: restarting flai serve, or the host restarting on a new
		// flai, ends the connection the request came on.
		"host.process": {action: ActionHost, describe: describeHost, detachTimeout: 60 * time.Second, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Process string `json:"process"`
				Action  string `json:"action"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Process != "serve" && in.Process != "mcp" && in.Process != "all" {
				return nil, "", bad("process must be serve, mcp, or all")
			}
			if in.Action != "start" && in.Action != "stop" && in.Action != "restart" {
				return nil, "", bad("action must be start, stop, or restart")
			}
			return []string{"host", in.Action, in.Process}, "", nil
		}, say: func(raw json.RawMessage) string {
			var in struct{ Process, Action string }
			_ = json.Unmarshal(raw, &in)
			return in.Action + " " + in.Process
		}},
		"host.upgrade": {action: ActionHost, describe: describeHost, detachTimeout: 6 * time.Minute, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"host", "upgrade"}, "", nil
		}},

		// checks.status: a read of what flai checks status already reports:
		// the current or last run for a story (S-0082).
		"checks.status": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"checks", "status", in.ID}, "", nil
		}),

		// checks.tail: a read, progress-streaming, of new lines in a run's
		// log since a byte offset, waiting a little for more while the run
		// is active (S-0082). Each line arrives as a log event, the same
		// machinery push and dashboard already stream progress with.
		"checks.tail": {reads: true, progress: true, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID   string `json:"id"`
				From int    `json:"from"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if in.From < 0 {
				return nil, "", bad("from must not be negative")
			}
			return []string{"checks", "tail", in.ID, "--from=" + strconv.Itoa(in.From), "--wait=20"}, "", nil
		}},

		// checks.run and checks.cancel: the host action (S-0082), running the
		// same commands the operator's own shell does. run streams progress
		// and is detached generously: the real time limit is the operator's
		// own (flai serve checks timeout), enforced by flai checks run
		// itself, not by this request surviving that long.
		"checks.run": {action: ActionChecks, progress: true, describe: describeChecksRun, detachTimeout: 2 * time.Hour, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"checks", "run", in.ID}, "", nil
		}},
		"checks.cancel": {action: ActionChecks, describe: describeChecksRun, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"checks", "cancel", in.ID}, "", nil
		}},
		// agent.restart: a new agent for a story whose agent dropped or
		// failed (S-0116, ADR-0043); flai serve agent restart judges whether
		// it may, and says why not.
		"agent.restart": agentNow("restart", "restarted"),
		// agent.start: a ready story's agent now, whatever the launcher's
		// rules say about when (S-0115); flai serve agent start judges
		// whether it may, and says why not.
		"agent.start": agentNow("start", "started"),
	}
}

// agentNow is the write that runs flai serve agent <verb> for a story, gated
// by the agent action, and journalled as what it did.
func agentNow(verb, did string) spec {
	return spec{action: ActionAgent, describe: describeAgentNow(did), build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		in, e := decode[struct {
			ID string `json:"id"`
		}](raw)
		if e != nil {
			return nil, "", e
		}
		if e := needID(in.ID); e != nil {
			return nil, "", e
		}
		if !strings.HasPrefix(in.ID, "S-") {
			return nil, "", bad("%s is not a story; only a story has an agent", in.ID)
		}
		return []string{"serve", "agent", verb, in.ID}, "", nil
	}}
}

// describeAgentNow reads flai serve agent start's and restart's --json shape
// for the journal.
func describeAgentNow(did string) func(res any, err *channel.Error) (outcome, detail string) {
	return func(res any, err *channel.Error) (outcome, detail string) {
		if err != nil {
			return "failed", err.Message
		}
		w, _ := res.(Written)
		var said struct {
			Story   string `json:"story"`
			Agent   string `json:"agent"`
			Command string `json:"command"`
			PID     int    `json:"pid"`
			Queued  string `json:"queued"`
		}
		_ = json.Unmarshal(w.Data, &said)
		if said.Queued != "" {
			return "done", fmt.Sprintf("queued another agent for %s until the in-progress limit has room", said.Story)
		}
		return "done", fmt.Sprintf("%s %s for %s as %s (pid %d)", did, said.Command, said.Story, said.Agent, said.PID)
	}
}

// describeChecksRun reads flai checks run/cancel's own --json shape
// (ChecksRun) for the journal, rather than describe's push/publish shape.
func describeChecksRun(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Story   string `json:"story"`
		Outcome string `json:"outcome"`
	}
	_ = json.Unmarshal(w.Data, &said)
	if said.Outcome == "" {
		return "done", said.Story
	}
	return "done", said.Story + ": " + said.Outcome
}

// describeHost reads flai host upgrade's --json shape for the journal; a
// process action says what was asked (its say).
func describeHost(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Restarting bool `json:"restarting"`
		Upgrade    struct {
			Installed string `json:"installed"`
			Previous  string `json:"previous"`
			Current   string `json:"current"`
		} `json:"upgrade"`
	}
	_ = json.Unmarshal(w.Data, &said)
	switch {
	case said.Restarting:
		return "done", fmt.Sprintf("installed flai %s over %s; the host restarts on it", said.Upgrade.Installed, said.Upgrade.Previous)
	case said.Upgrade.Current != "":
		return "done", "flai " + said.Upgrade.Current + " is the latest"
	}
	return "done", ""
}

// describeDashboardRestart, describeDashboardUpgrade, and describeDashboardStop
// read cmd/dashboard.go and cmd/dashboard_upgrade.go's own --json shapes for
// the journal, rather than describe's push/publish shape above. A failed
// call is described the same way regardless: "failed", err.Message.

func describeDashboardRestart(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Container  string `json:"container"`
		WasRunning bool   `json:"was_running"`
		Ref        string `json:"ref"`
	}
	_ = json.Unmarshal(w.Data, &said)
	verb := "started"
	if said.WasRunning {
		verb = "restarted"
	}
	return "done", fmt.Sprintf("%s %s (%s)", verb, said.Container, said.Ref)
}

func describeDashboardUpgrade(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Container string `json:"container"`
		Outcome   string `json:"outcome"`
		From      string `json:"from"`
		To        string `json:"to"`
	}
	_ = json.Unmarshal(w.Data, &said)
	switch said.Outcome {
	case "up-to-date":
		return "done", fmt.Sprintf("%s already running %s", said.Container, said.To)
	case "started":
		return "done", fmt.Sprintf("started %s (%s)", said.Container, said.To)
	case "upgraded":
		return "done", fmt.Sprintf("upgraded %s from %s to %s", said.Container, said.From, said.To)
	default:
		return "done", said.Container
	}
}

func describeDashboardStop(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Container string   `json:"container"`
		State     string   `json:"state"`
		Serves    []string `json:"serves"`
	}
	_ = json.Unmarshal(w.Data, &said)
	switch said.State {
	case "stopped":
		return "done", said.Container + " stopped"
	case "still-running":
		return "done", said.Container + " unregistered; still serving " + strings.Join(said.Serves, ", ")
	case "not-running":
		return "done", said.Container + " was not running"
	default:
		return "done", said.Container
	}
}

// withJSON asks for JSON before a -- that ends the flags, else at the end.
func withJSON(args []string) []string {
	for i, a := range args {
		if a == "--" {
			out := append([]string{}, args[:i]...)
			return append(append(out, "--json"), args[i:]...)
		}
	}
	return append(append([]string{}, args...), "--json")
}

// writeMethods turns the specs into methods: validate, consult the journal,
// run, and read the outcome the way the dashboard used to.
func writeMethods(run Runner, now func() time.Time, host Host) map[string]channel.Method {
	return methodsFrom(specs(), run, now, host)
}

// methodsFrom builds methods from a table of specs, with a journal of their own.
func methodsFrom(table map[string]spec, run Runner, now func() time.Time, host Host) map[string]channel.Method {
	if run == nil {
		run = ExecRunner
	}
	j := &journal{done: map[string]journalled{}, now: now}
	requests := host.Requests
	if requests == nil {
		requests = inMemory()
	}
	out := map[string]channel.Method{}
	for name, sp := range table {
		out[name] = func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			args, stdin, e := sp.build(p, raw)
			if e != nil {
				return nil, e
			}
			key := ""
			if !sp.reads {
				var id struct {
					RequestID string `json:"request_id"`
				}
				_ = json.Unmarshal(raw, &id)
				if !requestID.MatchString(id.RequestID) {
					return nil, bad("a write needs a request_id, so that a repeat after a lost connection is not done twice")
				}
				key = p.Root + "|" + name + "|" + id.RequestID
				if was, ok := j.get(key); ok {
					return was.res, was.err
				}
			}
			var id struct {
				RequestID string `json:"request_id"`
			}
			_ = json.Unmarshal(raw, &id)
			entry := Entry{At: now().UTC().Format(time.RFC3339), Method: name, Project: p.Key, Root: p.Root, By: owner(p), RequestID: id.RequestID}
			if sp.action != "" && !host.enabled(sp.action, p.Root) {
				entry.Action, entry.Outcome = sp.action, "disabled"
				host.record(entry)
				return nil, &channel.Error{Code: Disabled,
					Message: fmt.Sprintf("the host action %q is not enabled for this project. On the host, in the project, run: %s", sp.action, EnableCommand(sp.action)),
					Data:    map[string]any{"action": sp.action, "enable": EnableCommand(sp.action)}}
			}
			if sp.hostwide && !host.enabledEverywhere(sp.action) {
				entry.Action, entry.Outcome = sp.action, "disabled"
				host.record(entry)
				return nil, &channel.Error{Code: Disabled,
					Message: fmt.Sprintf("this setting is kept for every project on the host, so the host action %q must be enabled for every project. On the host run: %s", sp.action, EnableEverywhere(sp.action)),
					Data:    map[string]any{"action": sp.action, "enable": EnableEverywhere(sp.action), "hostwide": true}}
			}
			// A detached write is recorded before it acts, where a restart of
			// flai serve does not lose it: its own success can end the serve
			// taking it, before any answer is recorded, and the dashboard's
			// repeat then reaches the next serve (S-0109, found live in
			// S-0107: a serve restart ran twice).
			detached := key != "" && sp.detachTimeout > 0
			if detached {
				was, found, err := requests.begin(key, now(), sp.detachTimeout+journalKeeps)
				if err != nil {
					return nil, &channel.Error{Code: channel.CodeInternal, Message: "the request could not be recorded before acting, so it was not done: " + err.Error()}
				}
				if found {
					return repeated(id.RequestID, was)
				}
			}
			r := Run{Dir: p.Root, Args: withJSON(args), Stdin: stdin}
			if sp.progress {
				r.OnEvent = func(ev map[string]any) { channel.Progress(ctx, ev) }
			}
			execCtx := ctx
			if sp.detachTimeout > 0 {
				var execCancel context.CancelFunc
				execCtx, execCancel = context.WithTimeout(context.Background(), sp.detachTimeout)
				defer execCancel()
			}
			ran, err := run(execCtx, r)
			res, rerr := outcome(ran, err, sp.exits)
			// A detached write's own ctx being cancelled is not news — it is
			// the very connection this write's success can sever — so only
			// the ordinary case still checks it: a cancelled, non-detached
			// request is not cached, in case its result is a half-answer.
			if detached {
				if err := requests.finish(key, now(), res, rerr); err != nil {
					res, rerr = unrecorded(res, rerr, err)
				}
			} else if key != "" && ctx.Err() == nil {
				j.put(key, res, rerr)
			}
			if entry.Action = sp.action; entry.Action == "" && sp.uses != "" && host.enabled(sp.uses, p.Root) {
				entry.Action = sp.uses
			}
			if entry.Action == "" {
				entry.Action = sp.record
			}
			if entry.Action != "" {
				d := describe
				if sp.describe != nil {
					d = sp.describe
				}
				entry.Outcome, entry.Detail = d(res, rerr)
				if sp.say != nil && rerr == nil {
					entry.Detail = sp.say(raw)
				}
				host.record(entry)
			}
			return res, rerr
		}
	}
	return out
}

func (h Host) record(e Entry) {
	if h.Record != nil {
		h.Record(e)
	}
}

// describe says in a line what became of a host action, for the journal:
// what was pushed and published, or why not.
func describe(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Pushed    bool     `json:"pushed"`
		PushError string   `json:"push_error"`
		Reason    string   `json:"reason"`
		Tags      []string `json:"tags"`
		Published []string `json:"published"`
		Unpushed  *struct {
			Acceptances []string `json:"acceptances"`
			Tags        []string `json:"tags"`
		} `json:"unpushed"`
	}
	_ = json.Unmarshal(w.Data, &said)
	if said.Unpushed != nil && len(said.Tags) == 0 {
		said.Tags = said.Unpushed.Tags
	}
	switch {
	case said.PushError != "":
		return "failed", "not pushed: " + said.PushError
	case !said.Pushed:
		return "done", "nothing pushed: " + orElse(said.Reason, "nothing was pending")
	}
	detail = "pushed"
	if len(said.Tags) > 0 {
		detail += " with tags " + strings.Join(said.Tags, ", ")
	}
	if len(said.Published) > 0 {
		detail += "; published " + strings.Join(said.Published, ", ")
	}
	return "done", detail
}

func orElse(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

// outcome reads what flai said: JSON on standard output, warnings and the
// fatal message among the log events, and exit codes that carry a payload.
// missingItem is what flai says when an ID names nothing: workitem's
// "<ID> not found", possibly wrapped by the command that looked it up.
var missingItem = regexp.MustCompile(`(^|[ :])[EST]-\d+ not found$`)

func outcome(ran Ran, err error, exits map[int]int) (any, *channel.Error) {
	if err != nil {
		return nil, failed(err)
	}
	var warnings []string
	fatal := ""
	for _, ev := range ran.Events {
		switch ev["level"] {
		case "WARN":
			warnings = append(warnings, str(ev["detail"], str(ev["msg"], "")))
		case "FATAL":
			fatal = str(ev["err"], str(ev["msg"], ""))
		}
	}
	if warnings == nil {
		warnings = []string{}
	}
	if ran.Exit != 0 {
		if fatal == "" {
			fatal = fmt.Sprintf("flai exited with %d", ran.Exit)
		}
		if code, ok := exits[ran.Exit]; ok {
			var data map[string]any
			_ = json.Unmarshal(ran.Stdout, &data)
			// flai wraps the payload as {conflict: {...}} or {refused: {...}}; the caller gets it flat.
			for _, k := range []string{"conflict", "refused"} {
				if inner, ok := data[k].(map[string]any); ok {
					data = inner
				}
			}
			msg := strings.TrimPrefix(strings.TrimPrefix(fatal, "conflict: "), "refused: ")
			return nil, &channel.Error{Code: code, Message: msg, Data: data}
		}
		if strings.HasPrefix(fatal, "rule:") {
			return nil, &channel.Error{Code: Rule, Message: strings.TrimSpace(strings.TrimPrefix(fatal, "rule:"))}
		}
		// A write aimed at an item that is not there is the caller's mistake,
		// as it is for a read: 404, not a failure of flai's (S-0077; until
		// then it answered 500). The command runs as a process, so its words
		// are all there is to go by; a test holds them to workitem's.
		if missingItem.MatchString(fatal) {
			return nil, &channel.Error{Code: NotFound, Message: fatal}
		}
		return nil, &channel.Error{Code: channel.CodeInternal, Message: fatal}
	}
	data := json.RawMessage("null")
	if trimmed := bytes.TrimSpace(ran.Stdout); len(trimmed) > 0 {
		if !json.Valid(trimmed) {
			return nil, &channel.Error{Code: channel.CodeInternal, Message: "flai returned output that is not JSON"}
		}
		data = trimmed
	}
	return Written{Data: data, Warnings: warnings}, nil
}

// agentArgs are the flags that give a story's agent (S-0103), each value
// checked for shape first, since it becomes a command line argument.
func agentArgs(a *manifest.Agent) ([]string, *channel.Error) {
	if a.IsZero() {
		return nil, nil
	}
	if err := a.Validate(); err != nil {
		return nil, bad("%s", err)
	}
	var args []string
	if a.Harness != "" {
		args = append(args, "--harness="+a.Harness)
	}
	if a.Model != "" {
		args = append(args, "--model="+a.Model)
	}
	for _, k := range a.ConfigKeys() {
		if strings.TrimSpace(a.Config[k]) == "" {
			return nil, bad("agent config %s has no value", k)
		}
		args = append(args, "--agent-config="+k+"="+a.Config[k])
	}
	return args, nil
}

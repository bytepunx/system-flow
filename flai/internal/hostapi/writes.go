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

var requestID = regexp.MustCompile(`^[A-Za-z0-9._-]{8,64}$`)

// spec is one write: how its params become a command line.
type spec struct {
	// build validates and returns the arguments (without --json) and what goes on standard input.
	build func(p channel.Project, raw json.RawMessage) (args []string, stdin string, err *channel.Error)
	// exits maps exit codes that carry a payload on standard output to error codes.
	exits map[int]int
	// progress sends each log event to the dashboard while the command runs.
	progress bool
	// reads marks a command that changes nothing: no request ID is asked for.
	reads bool
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

func specs() map[string]spec {
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
				Type    string   `json:"type"`
				Title   string   `json:"title"`
				Nature  string   `json:"nature"`
				Parent  string   `json:"parent"`
				Tags    []string `json:"tags"`
				Touches []string `json:"touches"`
				Body    string   `json:"body"`
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
			case in.Type == workitem.Story:
				if !anyItemID.MatchString(in.Parent) || !strings.HasPrefix(in.Parent, "E-") {
					return nil, "", bad("a story needs its parent epic")
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
			return append(args, "--body-stdin", "--autocommit", "--trailer="+Trailer, "--", title), in.Body, nil
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

		"push.pending": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"push", "--pending", "--dry-run"}, "", nil
		}),
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
func writeMethods(run Runner, now func() time.Time) map[string]channel.Method {
	if run == nil {
		run = ExecRunner
	}
	j := &journal{done: map[string]journalled{}, now: now}
	out := map[string]channel.Method{}
	for name, sp := range specs() {
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
			r := Run{Dir: p.Root, Args: withJSON(args), Stdin: stdin}
			if sp.progress {
				r.OnEvent = func(ev map[string]any) { channel.Progress(ctx, ev) }
			}
			ran, err := run(ctx, r)
			res, rerr := outcome(ran, err, sp.exits)
			if key != "" && ctx.Err() == nil {
				j.put(key, res, rerr)
			}
			return res, rerr
		}
	}
	return out
}

// outcome reads what flai said: JSON on standard output, warnings and the
// fatal message among the log events, and exit codes that carry a payload.
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

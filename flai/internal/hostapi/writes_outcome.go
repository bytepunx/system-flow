package hostapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// writing counts the dashboard writes running in this process now, each a
// flai command, or a dev server starting (S-0438), that a restart of flai
// serve would cut (ADR-0158).
var writing atomic.Int64

// WritesRunning is how many dashboard writes run in this process now.
func WritesRunning() int { return int(writing.Load()) }

func (h Host) record(e Entry) {
	if h.Record != nil {
		h.Record(e)
	}
}

// notEnabled refuses a host action the operator has not enabled for which,
// a project, and says what enables it.
func notEnabled(action, which string) *channel.Error {
	return &channel.Error{Code: Disabled,
		Message: fmt.Sprintf("the host action %q is not enabled for %s. On the host, in the project, run: %s", action, which, EnableCommand(action)),
		Data:    map[string]any{"action": action, "enable": EnableCommand(action)}}
}

// publishArgs is publish.run's command line, before --json: flai release
// --pending, never forced. publishExits are its exit codes that carry a
// payload: 3 is a remote with newer release tags (S-0174) or one that moved
// (ADR-0067), a conflict.
var (
	publishArgs  = []string{"release", "--pending"}
	publishExits = map[int]int{3: Conflict}
)

// Publisher publishes as publish.run does, for a caller that is not a
// dashboard: the orchestrator's MCP tool release_publish (S-0222). It runs
// the same flai release --pending under the same push host action, reads its
// answer the same way, and journals the run as publish.run's are, with who
// asked as its By and why in its Detail, so that a release the orchestrator
// cut is told apart from the operator's.
type Publisher struct {
	// Run runs flai; nil runs this executable (ExecRunner).
	Run Runner
	// Now is the journal's clock; nil is time.Now.
	Now  func() time.Time
	Host Host
}

// Publish runs flai release --pending in p's folder on by's word, for the
// reason why, and journals it under method. While the push host action is
// off it runs nothing and refuses with Disabled; when flai refuses with exit
// 3 the refusal is a Conflict with flai's message, as publish.run's is.
func (pb Publisher) Publish(ctx context.Context, p channel.Project, method, by, why string) (Written, *channel.Error) {
	run, now := pb.Run, pb.Now
	if run == nil {
		run = ExecRunner
	}
	if now == nil {
		now = time.Now
	}
	entry := Entry{At: now().UTC().Format(time.RFC3339), Action: ActionPush, Method: method, Project: p.Key, Root: p.Root, By: by}
	if !pb.Host.enabled(ActionPush, p.Root) {
		entry.Outcome, entry.Detail = "disabled", why
		pb.Host.record(entry)
		return Written{}, notEnabled(ActionPush, "this project")
	}
	done := perf.Track(ctx, perf.Exec("flai", publishArgs))
	ran, err := run(ctx, Run{Dir: p.Root, Args: withJSON(publishArgs)})
	itemsChanged(p.Root, p.Root)
	done()
	res, e := outcome(ran, err, publishExits)
	entry.Outcome, entry.Detail = describe(res, e)
	if why != "" {
		entry.Detail = why + ": " + entry.Detail
	}
	pb.Host.record(entry)
	if e != nil {
		return Written{}, e
	}
	return res.(Written), nil
}

// Accepter accepts as accept.run accepts for the orchestrator, for a caller
// that is not a dashboard: the orchestrator's MCP tool accept (ADR-0278). It
// runs this host's flai accept --by orchestrator, so that the acceptance
// reads the host's configuration, its merge-queue and forge host actions
// among them, wherever the orchestrator runs. It builds the command line as
// accept.run builds the orchestrator's, refuses it as accept.run does while
// accept_reviews or, under pull-request acceptance, the forge host action is
// off, and journals each acceptance with who asked as its By. A dry run only
// reads, so it is not journaled.
type Accepter struct {
	// Run runs flai; nil runs this executable (ExecRunner).
	Run Runner
	// Now is the journal's clock; nil is time.Now.
	Now  func() time.Time
	Host Host
}

// AcceptArgs is the command line and standard input of the orchestrator's
// acceptance of id in p, before --json: with dryRun, flai accept id --by
// orchestrator --dry-run, with verified and evidence when given, which only
// previews; otherwise accept.run's for by orchestrator, verified, and
// evidence. It refuses, before flai runs, an ID that is not an item's, a
// verified that is not a commit's name, and an acceptance while p's
// accept_reviews is off.
func AcceptArgs(p channel.Project, id, verified, evidence string, dryRun bool) ([]string, string, *channel.Error) {
	if !dryRun {
		raw, err := json.Marshal(map[string]string{"id": id, "by": workitem.ActivityOrchestrator, "verified": verified, "evidence": evidence})
		if err != nil {
			return nil, "", failed(err)
		}
		return itemWriteSpecs()["accept.run"].build(p, raw)
	}
	if e := needID(id); e != nil {
		return nil, "", e
	}
	args, stdin := []string{"accept", id, "--by=" + workitem.ActivityOrchestrator, "--dry-run"}, ""
	if verified != "" {
		if !commitName.MatchString(verified) {
			return nil, "", bad("verified is the commit the orchestrator's verifier passed, in hex: %q is not one", verified)
		}
		args = append(args, "--verified="+verified)
	}
	if strings.TrimSpace(evidence) != "" {
		args, stdin = append(args, "--evidence=-"), evidence
	}
	return args, stdin, nil
}

// Accept runs the orchestrator's acceptance of id in p's folder, or with
// dryRun its preview, on by's word, and journals an acceptance under method.
// Under pull-request acceptance, while the forge host action is off, it runs
// nothing and refuses with Disabled, as accept.run does. flai's refusal for
// lint the fix left or a path still undone (exit 4) is Refused with its
// data; any other refusal is flai's message.
func (ac Accepter) Accept(ctx context.Context, p channel.Project, method, by, id, verified, evidence string, dryRun bool) (Written, *channel.Error) {
	run, now := ac.Run, ac.Now
	if run == nil {
		run = ExecRunner
	}
	if now == nil {
		now = time.Now
	}
	args, stdin, e := AcceptArgs(p, id, verified, evidence, dryRun)
	if e != nil {
		return Written{}, e
	}
	entry := Entry{At: now().UTC().Format(time.RFC3339), Action: ActionOrchestrate, Method: method, Project: p.Key, Root: p.Root, By: by}
	if action := forgeAcceptance(p); action != "" && !dryRun {
		entry.Action = action
		if !ac.Host.enabled(action, p.Root) {
			entry.Outcome, entry.Detail = "disabled", "accept "+id
			ac.Host.record(entry)
			return Written{}, notEnabled(action, "this project")
		}
	}
	done := perf.Track(ctx, perf.Exec("flai", args))
	ran, err := run(ctx, Run{Dir: p.Root, Args: withJSON(args), Stdin: stdin})
	if !dryRun {
		itemsChanged(p.Root, p.Root)
	}
	done()
	res, e := outcome(ran, err, map[int]int{4: Refused})
	if !dryRun {
		entry.Outcome, entry.Detail = "done", "accepted "+id
		if e != nil {
			entry.Outcome, entry.Detail = "failed", e.Message
		}
		ac.Host.record(entry)
	}
	if e != nil {
		return Written{}, e
	}
	return res.(Written), nil
}

// describe says in a line what became of a host action, for the journal:
// what flai release --pending pushed and published, or why not.
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
	}
	_ = json.Unmarshal(w.Data, &said)
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

// notADraft refuses, as a rule, what flai edit --no-draft answered
// unchanged: the story was not a draft (S-0201). flai edit leaves any edit to
// what is already there unchanged; finalizing is an act the dashboard reports.
func notADraft(w Written) *channel.Error {
	var res struct {
		ID        string `json:"id"`
		Unchanged bool   `json:"unchanged"`
	}
	_ = json.Unmarshal(w.Data, &res) // an answer that is not the edit's result is not "unchanged"
	if !res.Unchanged {
		return nil
	}
	return &channel.Error{Code: Rule, Message: fmt.Sprintf("%s is not a draft: there is nothing to finalize (someone may have finalized it meanwhile; reload the story)", res.ID)}
}

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
			msg := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(fatal, "conflict: "), "refused: "), "rule: ")
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
	if a.Provider != "" {
		args = append(args, "--provider="+a.Provider)
	}
	for _, k := range a.ConfigKeys() {
		if strings.TrimSpace(a.Config[k]) == "" {
			return nil, bad("agent config %s has no value", k)
		}
		args = append(args, "--agent-config="+k+"="+a.Config[k])
	}
	roles, e := roleArgs(a)
	return append(args, roles...), e
}

// roleArgs are the flags that give an agent's roles (S-0189), which flai
// agent set, flai story new, and flai edit take alike.
func roleArgs(a *manifest.Agent) ([]string, *channel.Error) {
	var args []string
	for _, n := range a.RoleNames() {
		r := a.Roles[n]
		if r.Harness != "" {
			args = append(args, "--role-harness="+n+"="+r.Harness)
		}
		if r.Model != "" {
			args = append(args, "--role-model="+n+"="+r.Model)
		}
		for _, k := range r.ConfigKeys() {
			if strings.TrimSpace(r.Config[k]) == "" {
				return nil, bad("agent role %s config %s has no value", n, k)
			}
			args = append(args, "--role-config="+n+"."+k+"="+r.Config[k])
		}
	}
	return args, nil
}

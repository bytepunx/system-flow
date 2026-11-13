package hostapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// writeMethods turns the specs into methods: validate, consult the journal,
// run, and read the outcome the way the dashboard used to.
func writeMethods(run Runner, now func() time.Time, host Host) map[string]channel.Method {
	out := methodsFrom(specs(), run, now, host)
	maps.Copy(out, methodsFrom(componentSpecs(), run, now, host))
	maps.Copy(out, methodsFrom(klusterSpecs(), run, now, host))
	maps.Copy(out, methodsFrom(imageSpecs(), run, now, host))
	maps.Copy(out, methodsFrom(queueSpecs(now, systemGit), run, now, host))
	return out
}

// systemGit is the git the merge queue's writes read the stories' branches
// with: this host's, timed on ctx.
func systemGit(ctx context.Context) execx.Runner { return execx.Timed(ctx, execx.System{}) }

// queueSpecs are the merge queue's writes under the merge-queue host action
// (ADR-0173, S-0399): queue.add and queue.remove enqueue and dequeue a story,
// and queue.hold holds the queue or resumes it, each stamped at now as the
// dashboard's owner. They run no flai command, since flai queue records the
// FLAI_AGENT a command run here has, flaiover: do writes the state in flai
// serve's own process through mergequeue, reading each story's branch with
// git. A refusal is Refused, with mergequeue's refusal as its data.
func queueSpecs(now func() time.Time, git func(ctx context.Context) execx.Runner) map[string]spec {
	checkStory := func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		_, e := queueStory(raw)
		return nil, "", e
	}
	checkHold := func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		_, e := queueHeld(raw)
		return nil, "", e
	}
	return map[string]spec{
		// queue.add: a story in review that is eligible, enqueued; one
		// enqueued already keeps its place and who enqueued it.
		"queue.add": {action: ActionMergeQueue, describe: describeQueue, build: checkStory,
			do: func(ctx context.Context, _ Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
				id, e := queueStory(raw)
				if e != nil {
					return nil, e
				}
				repo, err := workitem.Open(p.Root)
				if err != nil {
					return nil, queueFailed(p, fmt.Errorf("open the project to enqueue %s in the merge queue: %w", id, err))
				}
				return queueAnswer(p)(mergequeue.Add(git(ctx), repo, id, owner(p), now()))
			}},
		// queue.remove: a story enqueued, taken out, and the batch that holds
		// it dropped.
		"queue.remove": {action: ActionMergeQueue, describe: describeQueue, build: checkStory,
			do: func(_ context.Context, _ Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
				id, e := queueStory(raw)
				if e != nil {
					return nil, e
				}
				repo, err := workitem.Open(p.Root)
				if err != nil {
					return nil, queueFailed(p, fmt.Errorf("open the project to take %s out of the merge queue: %w", id, err))
				}
				return queueAnswer(p)(mergequeue.Remove(repo, id, owner(p), now()))
			}},
		// queue.hold: held true holds the queue, so that it builds no batch;
		// held false resumes it.
		"queue.hold": {action: ActionMergeQueue, describe: describeQueue, build: checkHold,
			do: func(_ context.Context, _ Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
				held, e := queueHeld(raw)
				if e != nil {
					return nil, e
				}
				what, op := "resume", mergequeue.Resume
				if held {
					what, op = "hold", mergequeue.Hold
				}
				repo, err := workitem.Open(p.Root)
				if err != nil {
					return nil, queueFailed(p, fmt.Errorf("open the project to %s the merge queue: %w", what, err))
				}
				return queueAnswer(p)(op(repo, owner(p), now()))
			}},
	}
}

// queueStory is the story queue.add and queue.remove name as id.
func queueStory(raw json.RawMessage) (string, *channel.Error) {
	in, e := decode[struct {
		ID string `json:"id"`
	}](raw)
	if e != nil {
		return "", e
	}
	if e := needStory(in.ID); e != nil {
		return "", e
	}
	return in.ID, nil
}

// queueHeld is what queue.hold asks for: true to hold the queue, false to
// resume it, and nothing else.
func queueHeld(raw json.RawMessage) (bool, *channel.Error) {
	in, e := decode[struct {
		Held *bool `json:"held"`
	}](raw)
	if e != nil {
		return false, e
	}
	if in.Held == nil {
		return false, bad("held is true to hold the merge queue, so that it builds no batch, or false to resume it")
	}
	return *in.Held, nil
}

// queueAnswer is a merge queue write's change as its answer in p: a
// *mergequeue.Refusal is Refused, with the refusal as data, so that the
// dashboard shows its sentence and each reason as flai gives them; a story no
// item has is NotFound; anything else failed.
func queueAnswer(p channel.Project) func(mergequeue.Change, error) (any, *channel.Error) {
	return func(ch mergequeue.Change, err error) (any, *channel.Error) {
		var refusal *mergequeue.Refusal
		switch {
		case errors.As(err, &refusal):
			return nil, &channel.Error{Code: Refused, Message: refusal.Error(), Data: refusal}
		case errors.Is(err, workitem.ErrNotFound):
			return nil, &channel.Error{Code: NotFound, Message: fmt.Sprintf("%s; name a story in review", err)}
		case err != nil:
			return nil, queueFailed(p, err)
		}
		return written(ch)
	}
}

// queueFailed is err, a merge queue write that failed in p, with the paths
// it names relative to the repository, as the dashboard sees them.
func queueFailed(p channel.Project, err error) *channel.Error {
	return failed(errors.New(hostRelative(p.Root, err.Error())))
}

// describeQueue says for the journal what a merge queue write changed, or
// that it changed nothing and why.
func describeQueue(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var ch mergequeue.Change
	_ = json.Unmarshal(w.Data, &ch)
	st := ch.State
	switch {
	case ch.Action == mergequeue.ActionAdd && !ch.Changed:
		return "done", ch.Story + " was enqueued already; nothing changed"
	case ch.Action == mergequeue.ActionAdd:
		return "done", fmt.Sprintf("enqueued %s, %d in the merge queue", ch.Story, len(st.Enqueued))
	case ch.Action == mergequeue.ActionRemove && st.HoldsStory(ch.Story) && st.Batch.Dropped != nil:
		return "done", "removed " + ch.Story + " from the merge queue and dropped the batch that holds it"
	case ch.Action == mergequeue.ActionRemove:
		return "done", "removed " + ch.Story + " from the merge queue"
	case ch.Action == mergequeue.ActionHold && !ch.Changed:
		return "done", "the merge queue was held already, by " + st.HeldBy + "; nothing changed"
	case ch.Action == mergequeue.ActionHold:
		return "done", "held the merge queue: it builds no batch until resumed"
	case !ch.Changed:
		return "done", "the merge queue was not held; nothing changed"
	}
	return "done", "resumed the merge queue"
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
			dir := p.Root
			if sp.dir != nil {
				if dir, e = sp.dir(p, raw); e != nil {
					return nil, e
				}
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
			which := "this project"
			if sp.about != nil {
				entry.Project, entry.Root = sp.about(raw)
				which = entry.Root
			}
			is := sp.action
			if sp.actionOf != nil {
				is = sp.actionOf(p, raw)
			}
			action := is
			if action == "" && sp.needs != nil {
				action = sp.needs(p)
			}
			if action != "" && !host.enabled(action, entry.Root) {
				entry.Action, entry.Outcome = action, "disabled"
				host.record(entry)
				return nil, notEnabled(action, which)
			}
			if (sp.hostwide || (sp.hostwideOf != nil && sp.hostwideOf(raw))) && !host.enabledEverywhere(sp.action) {
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
			execCtx := ctx
			if sp.detachTimeout > 0 {
				var execCancel context.CancelFunc
				execCtx, execCancel = context.WithTimeout(context.Background(), sp.detachTimeout)
				defer execCancel()
			}
			var res any
			var rerr *channel.Error
			if sp.do != nil {
				done := perf.Track(ctx, name)
				if key != "" {
					writing.Add(1)
				}
				res, rerr = sp.do(execCtx, host, p, raw)
				if key != "" {
					writing.Add(-1)
				}
				done()
			} else {
				r := Run{Dir: dir, Args: withJSON(args), Stdin: stdin}
				if sp.progress {
					r.OnEvent = func(ev map[string]any) { channel.Progress(ctx, ev) }
				}
				done := perf.Track(ctx, perf.Exec("flai", args))
				if key != "" {
					writing.Add(1)
				}
				ran, err := run(execCtx, r)
				// flai wrote in another process: this one's item snapshot looks again (S-0606)
				itemsChanged(p.Root, dir)
				if key != "" {
					writing.Add(-1)
				}
				done()
				res, rerr = outcome(answered(ran, sp.answers), err, sp.exits)
			}
			if w, ok := res.(Written); ok && rerr == nil && sp.judge != nil {
				if e := sp.judge(w); e != nil {
					res, rerr = nil, e
				}
			}
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
			if entry.Action = is; entry.Action == "" && sp.uses != "" && host.enabled(sp.uses, p.Root) {
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
				if sp.asked != nil {
					d = sp.asked(raw)
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

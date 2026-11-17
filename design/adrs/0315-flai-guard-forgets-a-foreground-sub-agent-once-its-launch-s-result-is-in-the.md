---
id: ADR-0315
title: "flai guard forgets a foreground sub-agent once its launch's result is in the session transcript, so one that ended with no SubagentStop no longer holds the story agent's wait_for_events"
status: proposed
date: 2026-11-17
supersedes: []
superseded_by: []
refines: [ADR-0092]
---

# ADR-0315 flai guard forgets a foreground sub-agent once its launch's result is in the session transcript, so one that ended with no SubagentStop no longer holds the story agent's wait_for_events

## Context

[ADR-0092](0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md) has `flai guard` record each sub-agent of a story agent's session at `SubagentStart` and forget it at `SubagentStop`, and refuse the story's agent `wait_for_events` while the record lists one. Its Consequences took a stale entry, left by a session that ended with a sub-agent out, as costless.

It is not. On 2026-11-15 two foreground task sub-agents ended on "API Error: Unable to connect to API: SSL certificate hostname mismatch": S-0738's `ac8d9be303813c862` (T-3202) and S-0739's `ae3487c9a0041b35d` (T-3203). No `SubagentStop` reached the guard for either, so each session's record kept it, and the guard refused the story's agent `wait_for_events` for its close-out's `verified` event for the rest of the session. `TaskStop` knew neither ID, so the agent had no way to clear it and polled the job's log with Bash instead ([I-0374](../issues/I-0374-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md), [I-0376](../issues/I-0376-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md)).

S-0819 read the two sessions' Claude Code 2.1.325 transcripts, `<projects>/<slug>/<session>.jsonl`, and the files beside them under `<session>/subagents/`:

| | S-0738, `ac8d9be303813c862` | S-0739, `ae3487c9a0041b35d` |
|-|-|-|
| Launch | `Agent`, `run_in_background` false | `Agent`, `run_in_background` false |
| Result in the session transcript | a `tool_result` for the launch's `tool_use_id`, `is_error` true, `toolUseResult` the error as a string, no `agentId` | a `tool_result` for the launch's `tool_use_id`, not marked an error, `toolUseResult.status` `completed`, `agentId` the sub-agent's, the text saying it terminated early on the API error |
| `subagents/agent-<id>.meta.json` | `toolUseId` the launch's, `requestShape` `foreground`, `finished` true | the same |
| `SubagentStop` | none: the record kept the ID | none: the record kept the ID |

A background launch's result is written at once, as the launch notice: `toolUseResult.isAsync` true, `status` `async_launched`, with the `agentId`, while the sub-agent goes on running. Every Claude Code hook's input carries `transcript_path`, the session's own transcript; a story agent's own `PreToolUse` names the main session's.

## Decision

At a story agent's `wait_for_events`, `flai guard` drops from the session's record each sub-agent whose foreground launch has its result in the session's transcript, and decides on what is left.

1. **The signal.** A foreground sub-agent blocks the call that launched it, so once the session's transcript holds that call's result, the sub-agent has ended, whether its result is an error or not. The guard reads the transcript the hook input's `transcript_path` names, line by line, passing over a line that does not parse. A sub-agent has ended when a `tool_result` there has a `toolUseResult` naming its `agentId` that is not a background launch's notice (`isAsync` true or `status` `async_launched`), or when the `tool_result` is for the `tool_use_id` its `subagents/agent-<id>.meta.json` beside the transcript names as its launch, with `requestShape` `foreground`. The second form covers a result that is only an error string, as S-0738's was.
2. **The prune.** Under the record's lock, the guard removes those sub-agents from the record, and the record file when none is left, as `SubagentStop` does, before it reads the record for the decision.
3. **It keeps the record when it cannot read.** A transcript it cannot open is warned of and prunes nothing, and a hook input with no `transcript_path`, as another harness's may be, prunes nothing and is logged at debug; the record then decides as ADR-0092 says. A meta file it cannot read leaves only the first form to decide that sub-agent. A record or project it cannot read still lets the call through (ADR-0092 §4).
4. **Nothing else changes.** `SubagentStart` and `SubagentStop` keep the record as before; a sub-agent launched in the background stays in it until its `SubagentStop`; the `PreToolUse` matcher already covers `mcp__flai__.*`, so no hook configuration changes.

## Consequences

- A foreground sub-agent that ends on an API error, or any other way that fires no `SubagentStop`, no longer holds its story's agent's `wait_for_events`: the wait passes once the launch's result is in the transcript, which it is by the time the agent makes its next call.
- A background sub-agent that ends with no `SubagentStop` still holds the wait, since its launch's result is the notice. That case did not occur in the instances; a sub-agent the harness cuts off at the session's end leaves an entry the next session of the same ID never reads.
- The guard reads the session's transcript at each story agent's `wait_for_events`, and only when the record lists a sub-agent. A story agent's transcript is a few megabytes, S-0739's 2.4 MB; the largest on this host is 102 MB. The read is a line scan that parses only the lines naming a `tool_result`, so its cost grows with the transcript, and a session that long pays it at each such wait.
- The rule depends on Claude Code's transcript format and on a file it keeps beside it, neither documented as a contract. A change to either fails open toward ADR-0092's behaviour, the refusal, not toward a wait let through while a sub-agent runs, save where a future format reuses `agentId` on a result while the sub-agent still runs.

## Alternatives considered

- **`finished` in the sub-agent's meta file alone.** Both stuck sub-agents' meta files say `finished: true`, but 288 of the 999 foreground meta files on this host have no `finished` at all, from older versions or from sub-agents cut off, so its absence says nothing. The transcript's result is what the API itself requires for every tool call, and the meta file is read only to link a result to its sub-agent.
- **Expire an entry by its `started` time.** A task sub-agent may run for over an hour, so any timeout either refuses a wait for a long time after an error or lets one through while a long sub-agent runs.
- **Let `TaskStop`, or a flai command, clear an entry.** `TaskStop` knows only background tasks, and a command the agent runs to clear its own hold makes the hold the agent's to lift, which ADR-0092 refused to the agent.
- **Record the launch's `tool_use_id` at the Agent call's `PreToolUse`.** The story agent's `PreToolUse` for an `Agent` call carries no sub-agent ID, and a layer's launches run at once, so the guard could not tell which `SubagentStart` follows which launch.

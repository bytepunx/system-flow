---
id: ADR-0092
title: "A story's agent waits for a sub-agent by launching it in the foreground, and flai guard refuses it wait_for_events while a sub-agent of its session runs and no thread on its story is open"
status: proposed
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0060]
---

# ADR-0092 A story's agent waits for a sub-agent by launching it in the foreground, and flai guard refuses it wait_for_events while a sub-agent of its session runs and no thread on its story is open

## Context

A story's agent that flai serve starts hands its tasks, searches, and checks to sub-agents and must wait for them. Claude Code 2.1.290 runs a sub-agent in the background unless the launch sets `run_in_background` to false. Agents waited for background sub-agents in two ways, and both failed. Holding the MCP tool `wait_for_events` ran each wait to its timeout, because the tool reports work items, threads, and narratives, and a sub-agent writes none of them; the sub-agent's finished notice was delivered only when the call returned ([I-0083](../issues/I-0083-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md)). Ending the turn worked for a short sub-agent, but `claude -p` ends the process ten minutes after the turn ends, with the sub-agent still running ([I-0084](../issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md)). The start prompt already said not to poll `wait_for_events`, and agents did so anyway.

S-0285 measured on 2.1.290, headless as flai serve runs it: a launch with `run_in_background` false returned an 11-minute sub-agent's result as the tool's result, with the session alive throughout; three such launches in one message ran together and the agent went on once all three had returned; a turn ended with a background sub-agent out ended the process 10m01s later, the sub-agent cut off and no `SubagentStop` fired. The `SubagentStart` and `SubagentStop` hooks' inputs carry `session_id`, `agent_id`, and `agent_type`.

## Decision

A story's agent waits for a sub-agent by launching it in the foreground, and `flai guard` refuses it `wait_for_events` while a sub-agent of its session runs and no thread on its story is open.

1. **How to wait.** flai serve's start prompt, `delegation.md`, and the `wait_for_events` description say to launch every sub-agent with the Agent tool's `run_in_background` set to false, a layer's in one message, so that each result comes back as the tool's result however long it runs; never to end the turn while a sub-agent runs in the background; and that `wait_for_events` is for a thread awaiting the designer.
2. **What flai records.** `flai guard` also runs as a `SubagentStart` and a `SubagentStop` hook, in a story's agent's session (`FLAI_STORY` set). It keeps each session's running sub-agents in the project's main checkout, `.flai-cache/guard/<session_id>.json`, under a lock, since a layer's sub-agents start at once. These hooks never refuse and print nothing.
3. **What flai refuses.** As a `PreToolUse` hook, it refuses the story agent's own `mcp__flai__wait_for_events` (no `agent_id`, `FLAI_STORY` set, no `FLAI_ROLE`) while the session's record lists a running sub-agent and no unresolved thread is on the story or one of its tasks. The refusal names the running sub-agents and says how to wait. A wait with such a thread open passes and returns on the designer's answer, as before.
4. **It fails open.** A record, a project, or threads it cannot read let the call through, as the guard does elsewhere ([ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)).

## Consequences

- A story's agent that follows the prompt keeps its session alive for a sub-agent of any length, and learns of its finish the moment it ends; the first remediation of I-0084 is made. Whether flai serve restarts an agent that ended without finishing stays the operator's to decide.
- An agent that holds `wait_for_events` on a sub-agent anyway is refused at once and told how to wait, instead of waiting to its timeout.
- The planner, the orchestrator, and an interactive session are not affected: the refusal needs `FLAI_STORY` and no role.
- A session that ended with a sub-agent out leaves it in the record, since no `SubagentStop` fires. That costs nothing for the wait the rule allows, on an open thread, and the file is per session.
- A project made from the template gets the two hooks with a template release, and needs a flai that has them; an older flai reads the new events as tool calls it does not refuse.

## Alternatives considered

- **A `SubagentStop` hook wakes a held `wait_for_events`.** It would shorten a mistaken wait but keep the habit, and `flai mcp` does not know which Claude Code session it serves, so the hook's `session_id` would have to be mapped to an agent name across processes. The refusal needs only what the guard already sees.
- **Better wording alone.** T-0869's sentence was in the prompt of both stories that held the wait.
- **Ending the turn on the harness's notice.** It works only for a sub-agent that finishes within ten minutes.
- **Refusing a background launch in the guard.** It would hold the agent to the foreground directly, but it would also refuse an agent that has work of its own to do while a sub-agent runs. Refusing the wait catches the costly case and leaves that one alone.

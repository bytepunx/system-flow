---
id: ADR-0058
title: "The agent host action lets the dashboard stop a story's agent, and a stop signals only the process flai started"
status: accepted
date: 2026-09-30
supersedes: []
superseded_by: []
refines: [ADR-0043]
topics: [cli, dashboard, server-side, client-side]
---

# ADR-0058 The agent host action lets the dashboard stop a story's agent, and a stop signals only the process flai started

## Context

flai serve starts a story's agent detached, in a session of its own, so that it outlives flai serve's restarts (ADR-0043). An agent that gets stuck therefore runs on, and the activity page shows it at work until someone kills it on the host by its PID. The living design said the dashboard could not stop or configure an agent; the starts it could ask for were Start agent, Retry, and Have an agent commit them, under the `agent` host action.

The operator asked in S-0170 to stop an agent from the activity page, with a confirmation that says what stopping does.

A run records the agent's PID in `serve/agents.json`, which outlives a reboot. After one the PID can belong to another process. On this host a pattern kill once reached the dashboard's container (I-0025).

## Decision

The `agent` host action lets the dashboard stop a story's agent, and a stop signals only the process flai started.

1. **What stops.** `flai serve agent stop <story>`, and the host method `agent.stop` behind the activity page's Stop, stop the newest agent flai serve started for the story. One that runs gets SIGTERM with its process group, then SIGKILL if it has not ended ten seconds later. One that ended waiting for an answer is not started again when the answer comes. Anything else is refused, saying why.
2. **Gated like a start.** `agent.stop` needs the `agent` action, as `agent.start`, `agent.restart`, and `agent.commit` do. The command in a shell works whether or not the action is on: it is the operator's own.
3. **Only the agent.** The process group is signalled only while the PID is still the agent's: alive, the leader of its own session (as flai starts it), and, where the system says when a process started, started when the agent did. On Linux a run records its process's start in clock ticks since boot (`/proc/<pid>/stat`), which no other process of the same boot has with the same PID, and which does not move when the wall clock does. Otherwise the run is only recorded as stopped. flai serve applies the same test when it looks for agents that have ended, so a run whose PID another process has after a reboot is settled by itself.
4. **Recorded as stopped.** The run gets the outcome `stopped`, with `stopped` set to when. The serving flai, when it sees the process end, keeps the stop rather than judging the end a failure. The dashboard shows the story's agent as failed, "stopped by the operator", with Retry. The story stays in its state, and its worktree keeps what the agent left. It gets no new agent until it is retried or moved back to ready.
5. **Confirmed first.** The activity page asks before it stops an agent. The dialog says the process and all it started end now and cannot be resumed, the worktree keeps what the agent left, and the story gets no agent until Retry or a move to ready.

## Consequences

- A holder of the dashboard token can end any agent flai serve started for the project while the `agent` action is on, as they could already start one.
- An agent stopped mid-edit leaves its worktree as it was, uncommitted work included. Whoever retries the story, or the new agent, finds it there.
- A stopped run reads as failed on the board's dot. The operator can tell it apart by its line and why, not by colour.
- An agent flai did not start, such as the operator's own session, is not in `serve/agents.json`, and nothing here can stop it.

## Alternatives considered

- **A host action of its own for stopping.** It would let the operator allow starts and not stops, or the other way round, which nobody asked for. Stopping is less than starting: whoever may start an agent may stop it.
- **A new activity state and dot colour for stopped.** Every place that reads the four states would change for a state that needs what a failure needs: a retry or a move.
- **Signalling the recorded PID as it stands.** After a reboot it can be anyone's process, and a group signal reaches all of that process's group.
- **Moving the story back to ready on a stop.** It would start another agent at once when there is room, which is not what stopping one is for. The story stays put, and the operator decides.

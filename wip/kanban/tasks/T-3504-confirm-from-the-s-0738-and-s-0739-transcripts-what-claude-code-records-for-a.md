---
id: T-3504
type: task
nature: research
title: Confirm from the S-0738 and S-0739 transcripts what Claude Code records for a foreground sub-agent that ends on an API error, and write the ADR that refines ADR-0092
status: backlog
parent: S-0819
owner: alex
created: 2026-11-17T03:34:37Z
updated: 2026-11-17T03:35:14Z
transitions: []
stream: S-0819
tags: [guard, adr]
touches: [design/adrs/README.md]
---
# T-3504 Confirm from the S-0738 and S-0739 transcripts what Claude Code records for a foreground sub-agent that ends on an API error, and write the ADR that refines ADR-0092

## Work

Find the Claude Code session transcripts of S-0738's and S-0739's agents, read only: the sub-agents `ac8d9be303813c862` (T-3202) and `ae3487c9a0041b35d` (T-3203) ended on "SSL certificate hostname mismatch". For each, record:

- the launching Agent call's input, and whether it set `run_in_background` to false;
- the Agent call's result in the main transcript: whether it is there, whether it carries `agentId: <id>`, and whether it is marked an error;
- whether a `SubagentStop` reached the guard, from what the sub-agent's own transcript and the session's record show;
- whether the `transcript_path` a `PreToolUse` hook input carries is the main session's transcript.

Compare a background launch's result, the launch notice, so the rule can tell the two apart. Then write a proposed ADR with `flai adr new` that refines ADR-0092: the guard drops from a session's record, at the story agent's `wait_for_events`, each sub-agent whose foreground launch has its result in the session's transcript, under the record's lock, and keeps the record when the transcript cannot be read. If the transcripts show another signal is surer, write that rule instead and say why in the ADR's Alternatives. Record the findings in the story's narrative.

This task waits for nothing: the rule it settles is what T-3505, T-3506, and T-3507 build and describe.

## Done when

- The narrative records the four findings for both sub-agents, with the paths of the transcripts read.
- A proposed ADR refining ADR-0092 states the rule, its signal, what it does when the transcript cannot be read, and the alternatives, and the ADR index lists it.

## Notes

---
id: T-2122
type: task
nature: remediation
title: The curate host action is on for system-flow and the first curator run on it is recorded
status: done
parent: S-0504
owner: alex
created: 2026-10-26T04:09:24Z
updated: 2026-10-28T12:14:35Z
transitions:
  - to: ready
    at: 2026-10-28T12:13:49Z
    by: agent-S-0504
  - to: in-progress
    at: 2026-10-28T12:13:50Z
    by: agent-S-0504
  - to: done
    at: 2026-10-28T12:14:35Z
    by: agent-S-0504
stream: S-0504
tags: [digest, curator]
after: [T-2120]
usage:
  source: log
  seconds: 45
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 4613
      cache_read: 1197312
      cache_write: 37571
      cost: 0.591
---
# T-2122 The curate host action is on for system-flow and the first curator run on it is recorded

## Work

Once the operator has answered the thread of T-2120, see the operational half of the remedy through. Confirm the action is on with `flai serve status` on the host; when it is not, say so on the thread of T-2120 and wait: `flai serve enable curate` is the operator's, and the agent neither runs it nor edits `~/.flai/config.json`.

With the action on, a curator run starts on the next acceptance `flai serve` sees, on `digest.schedule` when one is set, or when the operator runs `flai curate` in a shell. Do not run `flai curate` or the MCP tool `curate` yourself: `flai guard` refuses an agent `flai serve` started. When no acceptance is due, ask the operator on the same thread to run `flai curate` once. Hold `wait_for_events` meanwhile.

When the run has ended, read what it left: the `curate` journal entries in `~/.flai/serve/journal.jsonl` with their trigger, `wip/agents/curator.md` in the main checkout with its first entry and its `- Trigger:` line, and the run's log under `serve/agents/`. Record in the narrative's `## Decisions` when the action was turned on, what started the first run, how it ended, and how many entries it wrote, merged, absorbed, and compacted. Do not commit `wip/agents/curator.md` or the entries the run left: they are the main checkout's, and the operator's to commit.

This waits on T-2120 because the operator must agree to the remedy and turn the action on first.

## Done when

- `flai serve status` on the host lists `curate` on for system-flow.
- `wip/agents/curator.md` exists in the main checkout with at least one entry, and the host journal has a `curate` entry for this project.
- The narrative's `## Decisions` records when the action was turned on, the first run's trigger, outcome, and what it wrote.

## Notes

Drafted by planner-S-0504. The task changes no file on the story branch: the action lives in the host's configuration and the run writes in the main checkout. It is kept as a task of its own because it can fail apart from the code: the operator may not turn the action on, or the first run may fail.

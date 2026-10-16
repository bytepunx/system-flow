---
id: T-1799
type: task
nature: research
title: Recommend adapting or adopting, name the follow-on stories, and link the finding from agent-context.md
status: done
parent: S-0426
owner: alex
created: 2026-10-16T05:58:25Z
updated: 2026-10-16T07:51:25Z
transitions:
  - to: ready
    at: 2026-10-16T07:45:25Z
    by: agent-S-0426
  - to: in-progress
    at: 2026-10-16T07:45:26Z
    by: agent-S-0426
  - to: done
    at: 2026-10-16T07:51:25Z
    by: agent-S-0426
stream: S-0426
tags: []
touches: [design/system/agent-memory.md, design/system/agent-context.md]
after: [T-1798]
usage:
  source: log
  seconds: 359
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 56
      output: 23453
      cache_read: 2810203
      cache_write: 199458
      cost: 2.112
---
# T-1799 Recommend adapting or adopting, name the follow-on stories, and link the finding from agent-context.md

## Work

Add `## Recommendation` to `design/system/agent-memory.md`, drawn from the three sections before it:

- The options side by side: keep and improve flai's own search (such as a persisted incremental index), adapt a design into flai's Go code, or adopt an open source solution behind an optional configured adapter, each with what it gains in recall and search time, what it costs to build and run, and what it asks of the operator's host.
- The recommended path, with its reason and the trial figures behind it, and what it means for the pack, `doc_search`, E-0025's digest, and S-0427's store.
- The follow-on stories it implies, each as a title and a one-line goal, and the decisions the operator must make first. Write no ADR and no story: the operator decides, and the story's agent raises those decisions on a thread on S-0426 before review.

Add one line to `design/system/agent-context.md` linking the finding where it discusses retrieval, so a reader of the pack's design finds it. Change nothing else in that document.

Waits for T-1798: the recommendation rests on the trial's figures.

## Done when

- `## Recommendation` names the options, the recommended path with its reason, the follow-on stories, and the operator's decisions.
- `design/system/agent-context.md` links `agent-memory.md`.
- A thread on S-0426 puts the recommendation and the decisions to the operator.

## Notes

---
id: T-1094
type: task
nature: improvement
title: work-management.md and delegation.md, and their template copies, say a story's agent ends on an open question and is started again on the answer
status: done
parent: S-0272
owner: alex
created: 2026-10-06T22:53:02Z
updated: 2026-10-07T00:48:59Z
transitions:
  - to: ready
    at: 2026-10-07T00:42:41Z
    by: agent-S-0272
  - to: in-progress
    at: 2026-10-07T00:42:41Z
    by: agent-S-0272
  - to: done
    at: 2026-10-07T00:48:59Z
    by: agent-S-0272
stream: S-0272
tags: [conventions, template]
touches: [design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md]
usage:
  source: log
  seconds: 378
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 65
      cache_read: 588928
      cache_write: 57416
      cost: 0.2867
---
# T-1094 work-management.md and delegation.md, and their template copies, say a story's agent ends on an open question and is started again on the answer

## Work

Criteria 1 and 3, the conventions.

- **`work-management.md`:** in the rule on being blocked, "Do not wait idle", say what a story's agent does with a question only the designer can answer. It asks with `thread_open`, writes the narrative's `## Current state` and `## Next steps`, and ends, unless it has a task in progress to go on with. flai serve starts it again in its session when the thread is answered, and the answer is in its first `inbox`. `wait_for_events` answers `end: true` in that case. The agent ends when it does.
- **`delegation.md`:** the rule against waiting for a sub-agent says `wait_for_events` "is for a thread awaiting the designer" (about line 55). Change that so it no longer reads as an instruction to hold for an answer.

Make the same edits to the copies under `template/root/design/conventions/`, above each file's project marker. Add an entry to `template/CHANGELOG.md`.

`strategic-agents.md` is left alone: the planner and the analyzer still hold `wait_for_events` for an answer.

It waits for no task: the behaviour it describes is fixed by the story.

## Done when

- Both conventions and their template copies say to end on an open question, and nothing in them tells a story's agent to hold `wait_for_events` for an answer.
- `template/CHANGELOG.md` has the entry.
- The markdown lint and `flai check --strict` are clean.

## Notes

Drafted by the planner.

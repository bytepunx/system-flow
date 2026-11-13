---
id: T-3322
type: task
nature: remediation
title: I-0352 is closed with the guard rule that fixed it
status: done
parent: S-0730
owner: alex
created: 2026-11-13T19:54:19Z
updated: 2026-11-13T20:21:58Z
transitions:
  - to: ready
    at: 2026-11-13T20:18:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:18:20Z
    by: system-flow
  - to: done
    at: 2026-11-13T20:21:58Z
    by: agent-S-0730
stream: S-0730
tags: [flai]
touches: [design/issues/I-0352-a-story-s-agent-can-overwrite-a-work-item-file-whole-with-write-front-matter.md, design/issues/summary.md]
after: [T-3320]
usage:
  source: log
  seconds: 218
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1014
      cache_read: 850879
      cache_write: 8147
      cost: 0.2514
---
# T-3322 I-0352 is closed with the guard rule that fixed it

## Work

Close I-0352 once the rule that removes its cause is in place. That is why this task waits on T-3320.

- Run `flai issue close I-0352 --reason "<reason>"` in the story's worktree. The reason names the guard rule T-3320 added and the test that reproduces the PLACEHOLDER Write.
- Let flai update `design/issues/summary.md`. Do not edit the issue or the summary by hand.

## Done when

- I-0352's status is closed, with the reason naming the rule and its test.
- `design/issues/summary.md` no longer lists I-0352 as open.
- `flai check --strict` is clean on the story.

## Notes

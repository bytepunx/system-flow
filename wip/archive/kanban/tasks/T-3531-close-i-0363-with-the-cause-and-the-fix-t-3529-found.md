---
id: T-3531
type: task
nature: remediation
title: Close I-0363 with the cause and the fix T-3529 found
status: done
parent: S-0836
owner: alex
created: 2026-11-17T08:29:32Z
updated: 2026-11-17T11:42:40Z
transitions:
  - to: ready
    at: 2026-11-17T11:42:28Z
    by: agent-S-0836
  - to: in-progress
    at: 2026-11-17T11:42:28Z
    by: agent-S-0836
  - to: done
    at: 2026-11-17T11:42:40Z
    by: agent-S-0836
stream: S-0836
tags: [flai, issues]
touches: [design/issues/I-0363-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md, design/issues/summary.md]
after: [T-3529, T-3558]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 32
      cache_read: 611907
      cache_write: 1075
      cost: 0.2341
---
# T-3531 Close I-0363 with the cause and the fix T-3529 found

## Work

- Run `flai issue close I-0363 --reason "<reason>"` in the story's worktree, the reason naming the cause T-3529 recorded, the change to `flai/internal/mdlint`, and the fixture and test that reproduce it.
- Waits for T-3529 (`after`): the reason names the cause and the fix that task finds.

## Done when

- I-0363 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes

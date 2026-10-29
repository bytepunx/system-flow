---
id: T-2336
type: task
nature: improvement
title: Close I-0231 with the reason that flai now cuts a title's slug at 80 bytes in every file name it makes
status: done
parent: S-0542
owner: alex
created: 2026-10-29T15:24:12Z
updated: 2026-10-29T15:32:46Z
transitions:
  - to: ready
    at: 2026-10-29T15:32:33Z
    by: agent-S-0542
  - to: in-progress
    at: 2026-10-29T15:32:34Z
    by: agent-S-0542
  - to: done
    at: 2026-10-29T15:32:46Z
    by: agent-S-0542
stream: S-0542
tags: [issues]
touches: [design/issues/I-0231-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md, design/issues/summary.md]
after: [T-2334, T-2335]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2084
      cache_read: 310393
      cache_write: 17070
      cost: 0.2153
---
# T-2336 Close I-0231 with the reason that flai now cuts a title's slug at 80 bytes in every file name it makes

## Work

Run `flai issue close I-0231 --reason "..."` in the story's worktree, the reason naming the helper, the files whose names it cuts (work items, issues, threads, conversations, and, as before, ADRs and digest entries), the 80-byte limit, and the test that reproduces the instance. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

Waits for T-2334 and T-2335, so that the reason names what was built and documented.

## Done when

- I-0231 reads closed with that reason, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` scoped to S-0542 passes.

## Notes

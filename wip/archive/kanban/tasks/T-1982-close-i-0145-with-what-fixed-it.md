---
id: T-1982
type: task
nature: remediation
title: Close I-0145 with what fixed it
status: done
parent: S-0465
owner: alex
created: 2026-10-21T07:37:13Z
updated: 2026-10-21T08:08:44Z
transitions:
  - to: ready
    at: 2026-10-21T07:39:06Z
    by: agent-S-0465
  - to: in-progress
    at: 2026-10-21T08:08:36Z
    by: agent-S-0465
  - to: done
    at: 2026-10-21T08:08:44Z
    by: agent-S-0465
stream: S-0465
tags: [flaiover]
touches: [design/issues/I-0145-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md, design/issues/summary.md]
after: [T-1981]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1793
      cache_read: 363161
      cache_write: 14956
      cost: 0.198
---
# T-1982 Close I-0145 with what fixed it

## Work

Close I-0145 in the story's worktree with `flai issue close I-0145 --reason "<what fixed it>"`. The reason names the fix: `flai move <story> done --json` refuses as `flai accept --json` does, `item.move` answers it as `Refused` with the findings as data, and the board and the item page list them through `AcceptFindings`, so `linkSegments` no longer reads flai's wording. flai rewrites the issue file and `design/issues/summary.md`.

It waits for T-1981, the last task of the fix: the issue closes only once the fix is in.

## Done when

- I-0145 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean on the story.

## Notes

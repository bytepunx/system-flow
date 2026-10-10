---
id: T-1434
type: task
nature: remediation
title: I-0077 is closed with a reason naming MD056 as what fixed it
status: done
parent: S-0289
owner: alex
created: 2026-10-09T18:22:35Z
updated: 2026-10-09T18:29:18Z
transitions:
  - to: ready
    at: 2026-10-09T18:29:12Z
    by: agent-S-0289
  - to: in-progress
    at: 2026-10-09T18:29:12Z
    by: agent-S-0289
  - to: done
    at: 2026-10-09T18:29:18Z
    by: agent-S-0289
stream: S-0289
tags: [flai]
touches: [design/issues/I-0077-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md, design/issues/summary.md]
after: [T-1433]
usage:
  source: log
  seconds: 6
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 3
      output: 1036
      cache_read: 135101
      cache_write: 7251
      cost: 0.0972
---
# T-1434 I-0077 is closed with a reason naming MD056 as what fixed it

## Work

Close I-0077 from the story's worktree, so that the issue file and `design/issues/summary.md` change on the story branch:

```bash
flai issue close I-0077 --reason "S-0289: flai's mdlint reports MD056, table column count, splitting a row on every unescaped pipe as markdownlint does, so a code span of unescaped pipes in a table row is now refused where flai writes and fails TestRepositoryLintsClean"
```

It waits for T-1433 because the reason names the fix T-1433 makes, and the issue closes only once that fix passes its tests.

## Done when

- I-0077's status is closed, with the reason naming S-0289 and MD056.
- `design/issues/summary.md` no longer lists I-0077.
- `flai check --strict` passes on the issue files.

## Notes

Criterion 2 of S-0289.

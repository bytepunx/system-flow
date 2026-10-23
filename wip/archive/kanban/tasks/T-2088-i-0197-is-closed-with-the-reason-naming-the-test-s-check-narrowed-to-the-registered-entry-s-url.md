---
id: T-2088
type: task
nature: remediation
title: I-0197 is closed with the reason naming the test's check narrowed to the registered entry's URL
status: done
parent: S-0488
owner: alex
created: 2026-10-23T18:10:57Z
updated: 2026-10-23T19:47:46Z
transitions:
  - to: ready
    at: 2026-10-23T19:47:39Z
    by: agent-S-0488
  - to: in-progress
    at: 2026-10-23T19:47:40Z
    by: agent-S-0488
  - to: done
    at: 2026-10-23T19:47:46Z
    by: agent-S-0488
stream: S-0488
tags: [cli, issues]
touches: [design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md]
after: [T-2087]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 640
      cache_read: 259127
      cache_write: 8521
      cost: 0.1328
---
# T-2088 I-0197 is closed with the reason naming the test's check narrowed to the registered entry's URL

## Work

Waits for T-2087, since the reason names the fix it made.

In the story's worktree, close the issue through flai, so that `summary.md` is regenerated with it:

```sh
flai issue close I-0197 --reason "S-0488: TestServeSettingsAreShownToTheDashboard refuses the registered entry's URL 127.0.0.1:4242, not the bare 4242, which the settings' t.TempDir() paths can contain; a case with an import root named with 4242 keeps it so"
```

Then tick the story's second criterion with `flai criteria tick S-0488 2`, and the first once T-2087's test is green on the branch.

## Done when

- `design/issues/I-0197-…md` has `status: closed` and the reason above under its remediation, written by `flai issue close`.
- `design/issues/summary.md` no longer lists I-0197 as open.
- `flai check --strict` is clean for the story.

## Notes

Second layer: after T-2087.

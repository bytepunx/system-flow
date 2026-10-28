---
id: T-2183
type: task
nature: remediation
title: I-0199 is closed with the reason that fixed it, and the story's criteria are ticked
status: done
parent: S-0503
owner: alex
created: 2026-10-26T05:23:58Z
updated: 2026-10-26T05:39:33Z
transitions:
  - to: ready
    at: 2026-10-26T05:39:21Z
    by: agent-S-0503
  - to: in-progress
    at: 2026-10-26T05:39:22Z
    by: agent-S-0503
  - to: done
    at: 2026-10-26T05:39:33Z
    by: agent-S-0503
stream: S-0503
tags: [issues]
touches: [design/issues/I-0199-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md, design/issues/summary.md]
after: [T-2181, T-2182]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 2183
      cache_read: 397772
      cache_write: 11900
      cost: 0.2017
---
# T-2183 I-0199 is closed with the reason that fixed it, and the story's criteria are ticked

## Work

- In the story's worktree, run `flai issue close I-0199 --reason "..."` with a reason naming what fixed it: the `shellcheck` tier on `scripts/**/*.sh`, the `flaiover`, `template`, and `smoke` tiers' paths widened to the scripts they run, and the test `TestRepositoryManifestRunsScriptTestsForScriptChanges` that fails when a script test has no tier. Commit the issue file and the regenerated `design/issues/summary.md` with the task.
- Verify the story's two acceptance criteria against the worktree and tick each with `flai criteria tick S-0503 <n>`.

Waits for T-2181 and T-2182: the cause must be gone, and the test that reproduces it green, before the issue is closed.

## Done when

- `flai issue show I-0199` reads closed, with the reason above; `design/issues/summary.md` no longer lists it.
- Both criteria of S-0503 are ticked.
- `flai check --strict` is clean scoped to the story.

## Notes

Drafted by planner-S-0503. Run the issue command from the worktree, not the main checkout: `flai issue` writes under the current checkout.

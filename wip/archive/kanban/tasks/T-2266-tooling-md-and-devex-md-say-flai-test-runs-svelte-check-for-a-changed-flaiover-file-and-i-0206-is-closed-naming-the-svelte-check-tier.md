---
id: T-2266
type: task
nature: improvement
title: tooling.md and devex.md say flai test runs svelte-check for a changed flaiover file, and I-0206 is closed naming the svelte-check tier
status: done
parent: S-0530
owner: alex
created: 2026-10-28T18:20:49Z
updated: 2026-10-28T21:03:16Z
transitions:
  - to: ready
    at: 2026-10-28T21:02:50Z
    by: agent-S-0530
  - to: in-progress
    at: 2026-10-28T21:02:52Z
    by: agent-S-0530
  - to: done
    at: 2026-10-28T21:03:16Z
    by: agent-S-0530
stream: S-0530
tags: [flaiover, testing, docs]
touches: [design/conventions/tooling.md, design/system/devex.md, design/issues/I-0206-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md, design/issues/summary.md]
after: [T-2265]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 3049
      cache_read: 608336
      cache_write: 31426
      cost: 0.4066
---
# T-2266 tooling.md and devex.md say flai test runs svelte-check for a changed flaiover file, and I-0206 is closed naming the svelte-check tier

## Work

It waits for T-2265, because the documents and the closing reason name the tier, its paths, and its measured time as T-2265 built them.

- In `design/conventions/tooling.md`, under the project additions, change the `flai test` item: the tiers outside `--all` gain `svelte-check` on the whole flaiover project when a changed TypeScript, JavaScript, or Svelte file under `flaiover` selects it, and svelte-check no longer appears among the `--all` tiers. Leave the baseline above the marker alone.
- In `design/system/devex.md`, in the close-out row, change the sentence on a story that changes `flaiover/`: `flai test` and the close-out run svelte-check on the whole project through the `svelte-check` tier, and the `flaiover` tier covers it. Remove the claim that svelte-check runs only in `flaiover-test.sh`. Bump the `updated` date.
- Close the issue with `flai issue close I-0206 --reason "<reason>"`, run in the story's worktree. The reason names `scripts/flaiover-check.sh`, the `svelte-check` tier in `system-flow.yaml`, its paths, and the test in `flai/internal/verify/select_test.go` that fails without it.

## Done when

- Neither `design/conventions/tooling.md` nor `design/system/devex.md` says svelte-check runs only under `--all` or only in `flaiover-test.sh`.
- I-0206 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai test` on the changed markdown passes, and `flai check --strict` scoped to the story is clean.

## Notes

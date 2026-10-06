---
id: T-0900
type: task
nature: feature
title: The acceptance preview reports the orchestrator's four conditions as blockers
status: done
parent: S-0221
owner: alex
created: 2026-10-05T04:47:29Z
updated: 2026-10-06T11:24:16Z
transitions:
  - to: ready
    at: 2026-10-06T11:18:32Z
    by: agent-S-0221
  - to: in-progress
    at: 2026-10-06T11:18:32Z
    by: agent-S-0221
  - to: done
    at: 2026-10-06T11:24:16Z
    by: agent-S-0221
stream: S-0221
tags: [flai]
touches: [flai/internal/preview]
after: [T-0898]
usage:
  source: log
  seconds: 344
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 11332
      cache_read: 1739836
      cache_write: 48639
      cost: 0.8601
---
# T-0900 The acceptance preview reports the orchestrator's four conditions as blockers

## Work

In `flai/internal/preview/accept.go`, give `Accept` an option for an acceptance by the orchestrator. The option carries the commit the verifier passed. With it set, `Accept` keeps every blocker it reports today and adds one for each condition the ADR from T-0898 names that fails:

- The verified commit is missing or is not the story branch's head.
- A criterion is unticked. Name each one, read from the story's `## Acceptance criteria`.
- A file changed on the branch is outside the story's touches. Use `storygit.StoryDiff`, read the touches as `workitem` `Holds.Claim` reads them, and name each file.
- A thread on the story or one of its tasks is open. Use `threads.For` and name each thread.

Without the option, the preview is what it is today. The dashboard's acceptance dialog and `flai accept --dry-run` show `blockers[]` as they come.

`preview` has no tests of its own. Its tests go in `flai/cmd/accept_orchestrator_test.go`, through `flai accept --dry-run`. T-0907 writes that file, since `--dry-run` with the option needs its flags.

This task waits for T-0898, which fixes the conditions. S-0276, ahead in the pull order, also changes `preview.Accept`: start from what it left on main.

## Done when

- With the option set, `preview.Accept` lists a blocker for each failing condition, naming the criterion, the file, the thread, or the commit.
- With the option unset, its output is unchanged.
- `scripts/flai-test.sh` passes.

## Notes

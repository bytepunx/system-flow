---
id: T-3457
type: task
nature: feature
title: The rule verified refuses, or warns of at warn, a story's move by a transition with verify unless its last flai verify passed those suites at its branch's head
status: done
parent: S-0807
owner: alex
created: 2026-11-16T23:49:36Z
updated: 2026-11-18T10:48:35Z
transitions:
  - to: ready
    at: 2026-11-18T09:55:59Z
    by: agent-S-0807
  - to: in-progress
    at: 2026-11-18T09:55:59Z
    by: agent-S-0807
  - to: done
    at: 2026-11-18T10:48:35Z
    by: agent-S-0807
stream: S-0807
tags: [flai]
touches: [flai/internal/workflow/rules.go, flai/internal/workflow/rules_test.go, flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/cmd/move_workflow_test.go, flai/internal/verifyrecord/record.go, flai/internal/verifyrecord/record_test.go, flai/internal/verify/record_test.go]
after: [T-3454]
usage:
  source: log
  seconds: 3156
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 154
      output: 63114
      cache_read: 12473179
      cache_write: 574737
      cost: 6.8932
---
# T-3457 The rule verified refuses, or warns of at warn, a story's move by a transition with verify unless its last flai verify passed those suites at its branch's head

## Work

ADR-0312 §2 and criterion 1. It waits for T-3454 because it reads the per-suite record through `flai/internal/verifyrecord`.

- Add `verified` to the vocabulary in `flai/internal/workflow/rules.go`, for stories, refusing by default. Its severity is the one the transition's `verify:` gives, which S-0806 parses.
- In `Repo.Move` (`flai/internal/workitem/rules.go`), for a transition with `verify:`, read the story branch's head through `repo.Git` in its worktree, as `clean-worktree` does. Read the record at that head from `verifyrecord`.
- Pass only when the record passed at that head with a clean worktree and holds, for each suite the transition lists, every tier of it the story's paths select, passed or reused. `all` means every suite.
- Otherwise refuse, or warn at `warn`. The message names the suites missing or failed and `flai verify <S-nnnn>`.
- A repository with no git (`repo.Git` nil) skips the rule, as it skips `clean-worktree`. A story with git but no branch, worktree, or record fails the rule.

## Done when

- `flai/internal/workflow` and `flai/internal/workitem` tests cover a pass, a refusal, and a warning at `warn`.
- They cover a record at an older commit, a missing suite, and `all`.
- `flai/cmd/move_workflow_test.go` shows `flai move` refusing and warning through a workflow file with `verify:`.
- `flai test` passes on the changed paths.

## Notes

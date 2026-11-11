---
id: T-3008
type: task
nature: improvement
title: The trial judge keeps its verdicts per root and commit pair across calls, and reads the branch heads with one git process a call
status: done
parent: S-0678
owner: alex
created: 2026-11-10T10:31:25Z
updated: 2026-11-11T05:49:24Z
transitions:
  - to: ready
    at: 2026-11-11T05:20:02Z
    by: agent-S-0678
  - to: in-progress
    at: 2026-11-11T05:20:02Z
    by: agent-S-0678
  - to: done
    at: 2026-11-11T05:49:24Z
    by: agent-S-0678
stream: S-0678
tags: [cli]
touches: [flai/internal/storygit/trial.go, flai/internal/storygit/trial_test.go]
usage:
  source: log
  seconds: 1762
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 52
      output: 20457
      cache_read: 3334935
      cache_write: 149232
      cost: 1.9792
---
# T-3008 The trial judge keeps its verdicts per root and commit pair across calls, and reads the branch heads with one git process a call

## Work

`serve.Activity` makes a new trial judge on every `agent.status` call, through `withTrial` and `storygit.NewTrialJudge` (`flai/internal/storygit/trial.go`). So each call runs `git --version` once, one `git rev-parse` for each branch it judges, and one `git merge-tree --write-tree` for each overlap of a ready story with a story in progress, though neither branch moved since the last call. The launcher's look and `flai serve`'s starts make a judge the same way.

- Keep the raw `TrialMerge` answer process-wide, by the main checkout's root and the pair of commits merged. A merge of two given commits gives the same paths every time, so a verdict kept by its commits needs no other check. Apply `withoutResolved` with the judge's own `files` after the lookup, as now, since that set can change between calls. Bound what is kept for each root, such as by dropping the pairs no call has asked for since the last few calls, and say how in a comment.
- Check git's version once per process, not once per judge. An error is not kept, so a git installed later is seen.
- Read the heads of the local branches once per judge, with one `git for-each-ref` over `refs/heads/`, in place of a `git rev-parse` for each branch, and the main branch once per judge. A branch absent from that read has no head, as `rev-parse` failing says now.

Change no verdict. This is the first layer: it waits for no task, and it shares no path with the phases' task, so the two run together.

## Done when

- A test in `flai/internal/storygit/trial_test.go`, with a runner that counts git calls, shows a second judge over the same commits runs no `merge-tree` and no `git --version`, and one `for-each-ref`.
- A test shows a moved branch head is merged again, and that a judge's own `files` still take its resolved paths out of a kept answer.
- The existing trial tests pass unchanged.
- `flai test` passes on the changed files.

## Notes

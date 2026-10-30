---
id: S-0549
type: story
nature: improvement
title: The merge queue checks the main checkout's uncommitted files and the lint of main before it builds a batch, so a dirty main fails once in seconds and not every batch
status: backlog
parent: E-0031
owner: alex
created: 2026-10-30T01:52:13Z
updated: 2026-10-30T01:52:13Z
transitions: []
tags: [flai]
topics: [testing, git, cli]
touches: [flai/internal/serve, flai/internal/mergequeue, flai/internal/check, design/system/flai-cli.md, design/adrs, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0549 The merge queue checks the main checkout's uncommitted files and the lint of main before it builds a batch, so a dirty main fails once in seconds and not every batch

## Goal

A batch verifies its stories' merge on top of main, then lands it through the main checkout, where the acceptance preflight lints the checkout's uncommitted files too. In the fortnight 80 verifications passed for 49 landings: S-0541 went through 16 batches and 28 failed landings, about 2.7 hours of verification, because one uncommitted task note in the main checkout had a bad Markdown link (MD051) that nothing checked before the batch was built, and I-0203's stale generated file failed every later batch the same way. This story moves that check to the front: before a batch is built, the queue lints the main checkout's uncommitted files and runs `flai check` on main, and when either fails it opens one thread to the operator naming the file and holds the queue, instead of building, verifying, and failing a batch for each enqueue. A batch whose verification passed but whose landing failed for a reason outside its stories is landed again when main is clean, without another verification.

## Acceptance criteria
- [ ] The queue runs the landing's preflight (markdown lint of the main checkout's uncommitted files and `flai check`) before building a batch; on a finding it holds the queue, journals the file and the finding, and opens one thread to the operator, with a test
- [ ] A held queue resumes on the next `wip/` or main change that clears the finding, without the operator's action beyond the fix
- [ ] A batch whose verification passed and whose landing failed on main's state lands again from its verified commit when main is clean, with no new verification, and the journal says so
- [ ] `flai-cli.md`, the merge-queue ADR's successor or an addendum, and the user docs describe the preflight and the held state, and the dashboard's queue page shows it

## Tasks

The planner drafts them.

## Notes

Evidence: S-0541's batches in the serve journal and the fix of T-2332 by hand on 2026-10-29. The operator's uncommitted files in the main checkout are theirs; the queue reports them and does not change them.

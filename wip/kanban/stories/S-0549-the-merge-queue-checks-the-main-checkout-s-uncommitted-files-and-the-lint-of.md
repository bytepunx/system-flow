---
id: S-0549
type: story
nature: improvement
title: The merge queue checks the main checkout's uncommitted files and the lint of main before it builds a batch, so a dirty main fails once in seconds and not every batch
status: ready
parent: E-0031
owner: alex
created: 2026-10-30T01:52:13Z
updated: 2026-10-31T08:20:39Z
transitions:
  - to: ready
    at: 2026-10-31T07:24:09Z
    by: alex
tags: [flai]
topics: [testing, git, cli, dashboard]
touches: [flai/internal/serve, flai/internal/mergequeue, flai/internal/check, design/system/flai-cli.md, design/adrs, docs/users/flai.md, design/adrs/README.md, flai/cmd/accept_markdown.go, flai/internal/verify/findings.go, flai/internal/verify/findings_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/MergeQueue.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/operators/index.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 34
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 15
          output: 90
          cache_read: 443482
          cache_write: 57690
          cost: 0.1136
forecast:
  duration: 5h
  delivery: 2026-10-31T15:05:00Z
  basis: "Its own forecast of 5h; 3rd in the pull order with an in-progress limit of 5, behind S-0548, S-0550, S-0552, S-0546 and S-0547."
  by: flai
  at: 2026-10-31T08:20:39Z
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
- T-2381 An ADR decides the merge queue's preflight of main, its held state, its one thread to the operator, and the landing again without verification
- T-2382 The unfixed-markdown findings move from flai/cmd into flai/internal/verify, so the merge queue reports them as flai accept does
- T-2383 The merge queue runs the landing's preflight before it builds a batch and holds itself on a finding, with the files and findings in its state and event log
- T-2384 A verified batch whose landing failed on main's state is kept, not dropped, and lands again from its verified commits once main is clean, with no new verification
- T-2385 The board's queue and the dashboard's merge queue panel show the preflight hold with its files and findings and its thread
- T-2386 flai serve journals the preflight hold, opens one thread to the operator naming the files, resumes the queue on the change that clears it, and journals a landing again
- T-2387 flai-cli.md, the user docs, and the operators guide describe the preflight, the held state, its thread, and the landing again without verification

## Notes

Evidence: S-0541's batches in the serve journal and the fix of T-2332 by hand on 2026-10-29. The operator's uncommitted files in the main checkout are theirs; the queue reports them and does not change them.

### Planning

Planned by planner-S-0549 on 2026-10-31. The operator asked for the plan; no thread on S-0549 or E-0031 existed, so there was nothing earlier to revisit. The plan's thread is TH-0775.

Touches:

- Declared, all kept: `flai/internal/serve`, `flai/internal/mergequeue`, `flai/internal/check`, `design/system/flai-cli.md`, `design/adrs`, and `docs/users/flai.md`. The folder `design/adrs` stays because T-2381 adds an ADR whose file name `flai adr new` sets from its number and title, which no task can name yet. The folders `flai/internal/serve` and `flai/internal/mergequeue` are the operator's and are kept; the tasks name their files (`serve/mergequeue.go`, `serve/mergequeue_test.go`, `serve/serve.go`; `mergequeue/preflight.go` and its test, `state.go`, `batch.go`, `log.go`, `land.go`, `eligible.go`, with their tests), and flai replaces each folder in the claim with them (ADR-0096). The folder `flai/internal/check` is kept as declared, though no task changes it: the preflight calls `check.Run` and changes nothing there, so the agent may narrow it with `flai touches`.
- Layout, from the code the explorer read: `flai/cmd/accept_markdown.go` and a new `flai/internal/verify/findings.go` with its test (T-2382), since the lint findings are parsed in package `cmd`, which `mergequeue` cannot import; `design/adrs/README.md` (T-2381).
- Design: `flaiover/src/lib/components/MergeQueue.svelte` and its test, `design/components/flaiover/src/lib/components/MergeQueue.examples.json`, `design/system/flaiover-dashboard.md`, and `docs/users/flaiover.md` for criterion 4's queue page (T-2385); `docs/operators/index.md`, whose "When a landing fails" paragraph describes today's drop and rebuild (T-2387).
- Co-change: `flai touches suggest` ranks `docs/users/flai-reference.md` (17%), `docs/operators/index.md` (15%), and `design/system/flaiover-dashboard.md` (12%) first. The last two are taken. `flai-reference.md` is left out: no task adds a flag or a command. The rest (`settings.md`, `workflow.md`, `CHANGELOG.md`, `hostapi`) no criterion reaches.
- Topic `dashboard` added for T-2385.

Forecast: `flai forecast` gave 16m, 96 s per unit of size over 22 improvement stories at size 10. Raised to 5h: the size counts four criteria and six touches. It does not count the seven tasks, the new queue state, the kept batch that lands again, the serve watcher and thread with race tests, the Svelte panel, and four documents. S-0548, of the same epic and similar breadth, was raised to 6h for the same reason. flai's replan sets the delivery from that duration and the pull order.

Cost of delay: none written. Neither S-0549 nor E-0031 has inputs, so `flai cod` refuses. TH-0775 asks for them, recommending `time_lost_per_cycle` 20m (50 USD a week). The orchestrator cannot set inputs on a story past the backlog, so they wait on the operator.

After: none. S-0549 needs no other story's accepted change. It shares `flai/internal/mergequeue` with S-0552 (in progress) under the split of MS-0263: S-0552 changes the window in `batch.go`, `eligible.go`, `bisect.go`, and `log.go`; S-0549 adds `preflight.go` and keeps its edits to `batch.go`, `land.go`, `state.go`, `log.go`, and `eligible.go` small. S-0552's T-2373 also changes `MergeQueue.svelte`, its test, its examples, and `docs/users/flaiover.md`, which T-2385 changes. Whichever lands second rebases. `docs/operators/index.md` is claimed by S-0548 and S-0552 too.

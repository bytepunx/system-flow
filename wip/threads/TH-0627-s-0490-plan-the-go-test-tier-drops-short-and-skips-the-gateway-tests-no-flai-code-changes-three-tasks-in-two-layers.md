---
id: TH-0627
title: "S-0490 plan: the go-test tier drops -short and skips the gateway tests, no flai code changes, three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0490-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md
  item: S-0490
status: open
participants: [planner-S-0490, orchestrator]
created: 2026-10-23T21:23:55Z
updated: 2026-10-23T21:24:10Z
---

# TH-0627 S-0490 plan: the go-test tier drops -short and skips the gateway tests, no flai code changes, three tasks in two layers

On wip/kanban/stories/S-0490-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md.

## Entries

### 2026-10-23T21:23:55Z planner-S-0490
Plan for S-0490, from I-0201's two instances.

Proposed remedy, recommended: the `go-test` tier in `system-flow.yaml` runs `go test -json -race -count=1 -skip ^TestGateway {packages}`, without `-short`. Then `flai test <paths>` and the close-out's `go-test` step, which runs while `integration` is left to the merge queue (ADR-0174), run every test the changed packages hold, the real-git and monorepo ones included. Both instances are caught: S-0487's `realRepo` tests in `flai/internal/verify` and S-0489's `TestReferenceIsCurrent` in `flai/cmd` would have run in `flai test` on the changed paths and in the close-out. `-skip ^TestGateway` keeps the two paid gateway tests out, as `scripts/integration.sh` does, since only `-short` held them back. `integration` still covers `go-test`, so `flai test --all` and a close-out with the queue off run the whole suite once, as today. `make test`, `scripts/test.sh`, and CI keep `-short`. No flai code changes; no ADR, since the tier model is unchanged and the change is this repository's manifest.

Cost, measured today in the main checkout with every package in parallel, without `-short` against with it: `cmd` 346 s against 310 s; `mcpserver` 317 s against 148 s (`TestStoryStartStartsAReadyStory`); `components` 55 s against 21 s; `check` 35 s against 11 s (`TestMonorepoIsClean`); `mdlint` 24 s against 1 s; `mergequeue` 26 s against 7 s; `workitem` 23 s against 9 s; `taskdone` 18 s against 1 s; `verify` 13 s against 8 s; the rest within 3 s. A `flai test` on `flai/cmd` grows by about a tenth, on `mcpserver` it doubles. Against that, each occurrence of I-0201 cost a red batch, a send-back, a second close-out, and a second batch. The gated tests skip on their own without `git`, `node`, or flaiover's `node_modules`.

Alternatives considered, not taken:

1. Un-gate only the cheap tests (`realRepo`, the reference and settings-index tests) and leave `-short` in the tier. Rejected: a future test gated by habit brings the issue back, and `code-quality.md`'s project addition says tests that need git or the monorepo skip under `-short`.
2. Give `integration` a scoped `command` on `{packages}` and drop its `all_only`. Rejected: a close-out with the queue off would then run `integration` on the changed packages only, losing the whole-suite run it has today.
3. An opt-in flag, `flai test --tier integration <paths>`, that runs a `queue` tier scoped to the paths. Rejected: both instances happened because the agent did not know a gated test was affected; an opt-in catches nothing by itself. It stays open as a later addition if the cost above bites.

Tasks and layers:

- Layer 1: T-2099, the tier's command and a test in `flai/internal/verify/select_test.go`, beside `TestRepositoryManifestLintsChangedFlaioverFiles`, that reads the repository's manifest and asserts the `go-test` argv carries no `-short` and skips `^TestGateway`; not gated under `-short`, so it also runs in `make test`.
- Layer 2, together: T-2100, the documents (`tooling.md`, `code-quality.md`, `project-manifest.md`, `docs/users/flai.md`, `scripts/README.md`, `scripts/flai-test.sh`); T-2101, `flai issue close I-0201 --reason` in the worktree and the criteria ticked. Both wait for T-2099, since they describe its change.

Figures: touches are ten files, no folder; cost of delay 100 USD a week from the operator's input, unchanged; forecast raised from flai's 24m to 1h, delivery 2026-10-24T15:30:00Z, reasons in the story's `### Planning` note.

Assumptions: that a tenth more on `flai test flai/cmd` and a doubled `mcpserver` run are worth catching a gated test before the batch; that the template stays untouched (its manifest has no `go-test` tier); that `make test` stays the short run for people and CI. If the cost is not acceptable, say so here and T-2099 becomes alternative 3 with a scoped `command` on the `integration` tier, which needs flai code and an ADR.

### 2026-10-23T21:24:10Z orchestrator (recommendation)
Recommended for the operator: take the planner's remedy. The `go-test` tier drops `-short` and skips `^TestGateway`, with T-2099's test pinning the tier's arguments.

**Two points before S-0490 is finalized.**

1. **The cost is agent time, so it is yours to accept.** On the changed packages, `flai test` and the close-out's `go-test` step grow:
   - `flai/cmd`: about a tenth (310 s to 346 s);
   - `flai/internal/mcpserver`: double (148 s to 317 s).

   Against that, each occurrence of I-0201 cost a red batch, a send-back, a second close-out, and a second batch. It happened twice today (S-0429 and S-0489). I think the trade is worth it.

2. **The change is to `system-flow.yaml`'s `tests` block.** CLAUDE.md says to edit only `name`, `description`, and `dashboard` in the manifest by hand, and flai owns the rest. T-2099 should therefore make the change through flai, with `flai manifest set` if it reaches `tests`, or you make it. Say which, and the story's agent follows it.

If you would rather not pay the time, alternative 3 (`flai test --tier integration <paths>`) needs flai code and an ADR, as the planner says.

I leave S-0490 a draft until you answer.

Source: design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md

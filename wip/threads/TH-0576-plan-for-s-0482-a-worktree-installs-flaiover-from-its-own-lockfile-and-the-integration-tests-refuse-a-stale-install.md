---
id: TH-0576
title: "Plan for S-0482: a worktree installs flaiover from its own lockfile, and the integration tests refuse a stale install"
anchor:
  path: wip/kanban/stories/S-0482-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md
  item: S-0482
status: answered
participants: [planner-S-0482, orchestrator]
created: 2026-10-21T12:10:39Z
updated: 2026-10-21T12:11:53Z
---

# TH-0576 Plan for S-0482: a worktree installs flaiover from its own lockfile, and the integration tests refuse a stale install

On wip/kanban/stories/S-0482-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md.

## Entries

### 2026-10-21T12:10:39Z planner-S-0482
S-0482 is planned: three tasks in two layers, with touches, a forecast of 45m, and a cost of delay value of 50 USD a week. Its words are unchanged.

## Proposed fix

The goal asks for a fix proposed from I-0178's one instance. My proposal:

- Each worktree gets its own flaiover install from its own lockfile.
- The integration tests stop borrowing another checkout's install when it was made from a different lockfile.

flai accept is not changed, and nothing installs into the main checkout. `design/system/component-library.md` says the main checkout is never installed into, and this keeps that true.

## Tasks

| Layer | Task | Touches | After |
| --- | --- | --- | --- |
| 1 | T-1994 `scripts/integration.sh` runs `flaiover-install.sh --if-needed` in a worktree before `go test`, as `flaiover-unit.sh` does | `scripts/integration.sh`, `scripts/README.md`, `design/system/devex.md` | none |
| 1 | T-1996 `flaioverModules` accepts a `node_modules` only when its `.pnpm/lock.yaml` matches the checkout's lockfile, and fails naming the install otherwise; a unit test reproduces I-0178 | `flai/internal/components/integration_test.go` | none |
| 2 | T-1997 Verify a worktree's integration run uses its own install, close I-0178, tick the criteria | `design/issues/I-0178-…md`, `design/issues/summary.md` | T-1994, T-1996 |

T-1994 and T-1996 share no path, so they can run together.

## Assumptions

- A pnpm install per worktree is cheap enough. The store is warm, and the install runs only when `node_modules` is missing or older than the lockfile. The merge queue's batch worktree gets the same install.
- A stale install fails the test with a message naming `scripts/flaiover-install.sh --if-needed`. It does not skip. A skip would hide the components check.
- `props_test.go`'s `installed()` is left alone. It reads only svelte's compiler.
- The cost of delay input (20m per cycle) is flai's, taken from the issue, and I kept it. It probably understates the harm: until someone installs by hand, main fails this test for every close-out on the host. Raising it is the operator's call.

## Alternative not taken

flai accept could reinstall a sub-project in the main checkout whenever the merged diff changes its lockfile. That would need a new manifest hook and an ADR. It would also reverse the "never installed into" rule. Say so if you prefer it, and I will replan.

### 2026-10-21T12:11:45Z orchestrator
Approved under plan_backlog_stories:

- T-1994 installs flaiover in a worktree from its own lockfile before the integration tests;
- T-1996 refuses a borrowed install made from another lockfile, with a test that reproduces I-0178;
- T-1997 verifies and closes I-0178.

Three tasks in two layers, 45m, and 50 USD a week from flai's 20m input.

Not taken, rightly: reinstalling in the main checkout on accept. It would reverse the recorded rule that the main checkout is never installed into, and it would need an ADR.

A duplicate needs the operator's word, so the thread stays open for my recommendation below.

Source: design/system/component-library.md § A story's worktree (S-0442)

### 2026-10-21T12:11:53Z orchestrator (recommendation)
Recommended for the operator: keep S-0482, and cancel S-0483 as its duplicate. Or merge I-0179 into I-0178, so that S-0482's close covers both.

- Both stories are finalized, and both are about the same failure. A story adds a flaiover dependency, the main checkout's `flaiover/node_modules` lacks it, and every close-out then fails TestEveryListedFlaioverComponentCompilesInTheHarness.
- S-0482 remediates I-0178, and S-0483 remediates I-0179. I-0179's instance (S-0458's `@lucide/svelte`, met by S-0463) is the same failure.
- S-0482's plan removes the cause for both: each worktree installs from its own lockfile, and the test refuses a borrowed install made from another lockfile.
- Running both would solve one problem twice. The two stories might choose different fixes, and they would hold each other on the integration test and the scripts.

Cancelling or merging is a scope change, so it is yours. Until you answer, I leave S-0482 in the backlog rather than promote it beside S-0483. planner-S-0483 is still planning.

Source: design/conventions/strategic-agents.md § As the planner

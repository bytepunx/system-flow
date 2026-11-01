---
id: ADR-0234
title: "Before it builds a batch the merge queue runs the landing's preflight on the main checkout, holds itself on a finding with one thread to the operator until a look finds main clean, and lands a verified batch whose landing failed on main's state again without verifying it again"
status: accepted
date: 2026-11-01
supersedes: []
superseded_by: []
refines: [ADR-0173, ADR-0176, ADR-0178]
---

# ADR-0234 Before it builds a batch the merge queue runs the landing's preflight on the main checkout, holds itself on a finding with one thread to the operator until a look finds main clean, and lands a verified batch whose landing failed on main's state again without verifying it again

## Context

A merge queue batch verifies its stories' merge on top of main, then lands each story through `flai accept <story> --batch <batch>` in the main checkout ([ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md), [ADR-0178](0178-the-merge-queue-lands-each-story-in-the-main-checkout-with-git-s-own-check-while.md)). Before it merges, that acceptance runs the lint fix over the main checkout's uncommitted files but those outside `wip/`, which it leaves, and refuses the acceptance for what the fix leaves ([ADR-0140](0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md)). Nothing checks that before the batch is built.

When a landing fails, `Land` drops the batch, and the next look builds and verifies it again. In the fortnight to 2026-10-29, 80 verifications passed for 49 landings. S-0541 went through 16 batches and 28 failed landings, about 2.7 hours of verification, because one uncommitted task note in the main checkout had a Markdown link that broke MD051. [I-0203](../issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md)'s stale generated file on main failed every later batch the same way, at verification.

A whole-repository `flai check --strict` on main is never clean: `item.archive`, `issues.no-story`, `threads.archived`, and `components.examples` warnings stand there outside any story, and a check scoped to a batch's stories passes over them ([ADR-0183](0183-the-merge-queue-s-batch-and-each-bisection-step-run-their-tiers-with-merge.md), D-0338).

The question is what the queue checks of main before it builds a batch, what it does on a finding, how it starts again, and what becomes of a verified batch whose landing failed on main rather than on its stories.

## Decision

Before it builds a batch, the merge queue runs the landing's preflight on the main checkout, holds itself on a finding with one thread to the operator, starts again on the look that finds main clean, and lands a verified batch whose landing failed on main's state again from its verified commits, without verifying it again.

- **The preflight** runs at each look that would build a batch: the queue may run, no batch waits to land, and an enqueued story is eligible. It has two parts.
  - The lint the landing runs: the commands of the tiers that carry a `fix_command`, over the main checkout's uncommitted files but those outside `wip/`, which a landing leaves (ADR-0178), as `flai accept --batch` runs them through `verify.CheckoutUnfixed`.
  - `flai check` on the main checkout, unscoped, of which only error-level findings count. A warning holds nothing, since the batch's own check passes over warnings outside its stories, and main always has some.
- **The preflight changes no file.** The operator's uncommitted files are theirs. The landing's lint runs its fix first, so a finding the fix would clear is not main's fault. When the lint finds something in place, the preflight confirms it in a scratch worktree of main's head under `.flai-cache`, with those files copied in, the fix run there, and the check run again. Only a finding that stands there counts. The scratch worktree is removed after.
- **A finding holds the queue** in a state of its own, `preflight`, apart from the operator's hold (`flai queue hold`). It records when it began, main's head, and each finding with its file, line, rule, and message. `flai queue resume` does not clear it, and `flai queue hold` and `resume` leave it as it is. While it stands, the queue builds no batch and runs no tier. The event log gets a `preflight_held` event when it begins and `preflight_cleared` when it clears.
- **One thread to the operator** per hold, opened by `flai serve` when the hold begins, anchored on `system-flow.yaml` so that a hold has one thread however many files it names. It names each file, line, and rule, says the queue builds no batch until they are fixed, and says the queue changes none of them. Its ID is kept on the hold. While it is open, no other is opened for the same hold. flai serve resolves it, saying why, when the hold clears. The host's journal records the hold once, and its clearing.
- **The hold clears on its own.** While it stands, every look runs the preflight again, and `flai serve` wakes the queue on any change in the main checkout, `wip/` included, and on a move of main's head. The look that finds main clean clears the hold, resolves the thread, and builds the batch. The operator does nothing but the fix.
- **A landing that failed on main's state lands again.** When `flai accept --batch` fails, the queue runs the preflight. When it finds something, the failure is main's, not the stories': the batch is kept, with its landing and its verified heads, and the queue holds on the findings as above. Once the preflight passes, the landing goes on from the next story not landed, from the commits the batch verified, with no tier run, while main's head is still where the batch expects it and each story's branch is still at its verified head (`Landable`). The event log gets a `landing_resumed` event, and the journal says the batch landed again without a new verification. When the preflight finds nothing, the failure is the stories', and the batch is dropped as ADR-0176 says. When main moved meanwhile, it is dropped as `DropIfMainMoved` drops it today.
- `flai queue`, `flai board`, and the dashboard's merge queue panel show the preflight hold, its findings, and its thread.

This refines ADR-0173, ADR-0176, and ADR-0178.

## Consequences

- A dirty main fails once, in seconds, with one thread naming the file, instead of building, verifying, and failing a batch for every enqueue.
- A verified batch is not verified again for a fault outside its stories.
- Each look that would build a batch runs one unscoped `flai check` and the fixing tiers' lint over the uncommitted `wip/` files, a few seconds. A confirmation in the scratch worktree adds a checkout of main's head, and runs only when the lint found something in place.
- The queue does not hold on a check warning. A warning that later fails the batch's own scoped check still fails the batch, as today.
- An error-level `flai check` finding on main that no batched story would have failed on now holds the queue. That is intended: the batch's smoke tier runs the check too, and a main the check rejects should be fixed first.
- A failure the preflight cannot see, such as a rebase that stops, still drops the batch.

## Alternatives considered

- **Hold on a strict check of main.** Lost: main's standing warnings would hold the queue for good (D-0338).
- **Run the lint fix on the main checkout's files.** Lost: the files are the operator's, and a queue that rewrites them while they work on them would surprise them.
- **Hold on what the in-place lint finds, without the fix.** Lost: the landing fixes first, so the queue would hold on a finding the landing would clear.
- **Reuse the operator's `held` flag.** Lost: `flai queue resume` would clear a finding that still stands, and the board could not say why the queue is held.
- **A thread per file, anchored on the file.** Lost: an uncommitted file may be untracked or removed by the fix, and a hold of several files would open several threads for one fault.
- **Keep dropping the batch and rely on the preflight alone.** Lost: a file can go dirty between the preflight and the landing, and each such drop costs a verification.

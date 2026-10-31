---
id: T-2381
type: task
nature: improvement
title: An ADR decides the merge queue's preflight of main, its held state, its one thread to the operator, and the landing again without verification
status: backlog
parent: S-0549
owner: alex
created: 2026-10-31T07:40:30Z
updated: 2026-10-31T07:42:23Z
transitions: []
stream: S-0549
tags: [flai]
touches: [design/adrs/README.md]
---
# T-2381 An ADR decides the merge queue's preflight of main, its held state, its one thread to the operator, and the landing again without verification

## Work

Write a new ADR with `flai adr new` that refines ADR-0173, ADR-0176, and ADR-0178. No accepted ADR covers a check of main before a batch is built, so this is a new one, not an edit. It decides:

- What the preflight runs: the `fix_command` tiers' check over the main checkout's uncommitted files, as `flai accept` runs it through `verify.CheckoutUnfixed` (`flai/cmd/accept_markdown.go`), and `check.Run` on main. Which `flai check` findings hold the queue: the recommendation is error-level findings only, since a batch's own check is scoped to its stories (ADR-0183) and an unrelated warning must not stop every batch.
- The held state: a field of its own on the queue's state, apart from the operator's `held` (`flai queue hold`), so that `flai queue resume` does not clear a finding that still stands and the reason shows.
- The one thread: opened by `flai serve` to the operator when the hold starts, naming each file and finding, reused while open (as `openFailing` reuses its thread), and resolved by flai when the hold clears. Where it is anchored: the recommendation is `system-flow.yaml`, as `flai serve`'s upgrade notice is, so a hold has one thread however many files it names.
- How the hold clears: the preflight runs again on every look while held, and the watcher wakes the queue on any change in the main checkout, `wip/` included, or a move of main's head.
- The landing again: a batch whose verification passed and whose landing failed on main's state, not on its stories, is kept with its verified commits (`Heads[].Verified` and `Rebased`) instead of dropped, and lands again from them once the preflight passes, with no new verification, as long as main's head is the one it was verified against; the journal and the event log say so.

The operator's uncommitted files stay theirs: the queue reports them and changes none. First layer: it waits for nothing, and the code tasks wait for it.

## Done when

The ADR is written under `design/adrs/`, listed in `design/adrs/README.md`, passes `flai check --strict`, and settles each point above. Story criterion 4 (the ADR).

## Notes

Drafted by the planner of S-0549 on 2026-10-31.

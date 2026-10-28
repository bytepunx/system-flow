---
id: T-2308
type: task
nature: improvement
title: "The context package primes the orchestrator for named reasons: it parses them, refuses an unknown one naming the known ones, and loads the role's conventions once and, for a thread, a draft, and a plan, their items and threads"
status: backlog
parent: S-0537
owner: alex
created: 2026-10-28T22:06:42Z
updated: 2026-10-28T22:06:42Z
transitions: []
stream: S-0537
tags: [flai, prime, orchestrator]
touches: [flai/internal/context/reasons.go, flai/internal/context/reasons_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/context/context.go]
---
# T-2308 The context package primes the orchestrator for named reasons: it parses them, refuses an unknown one naming the known ones, and loads the role's conventions once and, for a thread, a draft, and a plan, their items and threads

## Work

ADR-0209 §3, criteria 1 and 2. Add a reason-primed pack beside `ForStrategic` in `flai/internal/context/strategic.go`, in a new `reasons.go`:

- Parse reasons in the form S-0536's start prompt writes, `<kind>:<id>` or a bare `<kind>`: `thread:TH-nnnn`, `review:S-nnnn`, `planned:E-nnnn|S-nnnn`, `draft:S-nnnn`, `release`, and `ready`, with the trigger's other events mapped onto them (`to-backlog` and `draft-changed` onto `draft`, `ready-empty` and `reordered` onto `ready`; `asked` and `schedule` onto `ready` and `release` together). Refuse an unknown kind, or a kind missing its ID, naming the known kinds.
- Load the role's conventions once, whatever the number of reasons, with the sections the reasons' topics leave out taken out (T-2310 gives the conventions those sections). Load no design briefs, no ranked sections, and no catalog: the general pack is not loaded.
- Per reason: for `thread`, the thread with its entries and the item it is on, whole; for `draft`, the draft story whole; for `planned`, the planner's thread on the item and the items its run wrote (the item and its children). The same item named by two reasons loads once.
- The digest step ranks the `answer` and `recommendation` entries against the reasons' threads only, not every thread awaiting the operator, within the budget's eighth (`DigestShare`, ADR-0194): reuse `answersFor`.
- Name the reasons in the pack's header and `--json`, so a reader sees what it was primed for.

`review`, `release`, and `ready` are parsed here and loaded by T-2309.

No `after`: it waits for nothing. It and T-2310 run together; this task's tests build their own conventions.

## Done when

- `reasons_test.go` has a test per reason this task loads, asserting what it loads, and one for the refusal naming the known kinds.
- A test builds the general orchestrate pack and a pack for one thread over the same fixture, and asserts the second is under a quarter of the first's size.
- A test with two reasons on the same item loads the conventions and the item once.
- `flai test` passes on the changed paths.

## Notes

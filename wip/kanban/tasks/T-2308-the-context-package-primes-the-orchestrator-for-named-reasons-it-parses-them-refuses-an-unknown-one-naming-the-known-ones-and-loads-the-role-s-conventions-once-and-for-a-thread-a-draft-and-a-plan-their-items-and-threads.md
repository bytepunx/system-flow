---
id: T-2308
type: task
nature: improvement
title: "The context package primes the orchestrator for named reasons: it parses them, refuses an unknown one naming the known ones, and loads the role's conventions once and, for a thread, a draft, and a plan, their items and threads"
status: done
parent: S-0537
owner: alex
created: 2026-10-28T22:06:42Z
updated: 2026-10-29T15:39:34Z
transitions:
  - to: ready
    at: 2026-10-29T15:25:19Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T15:25:21Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T15:39:34Z
    by: agent-S-0537
stream: S-0537
tags: [flai, prime, orchestrator]
touches: [flai/internal/context/reasons.go, flai/internal/context/reasons_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/context/context.go, flai/internal/context/parts.go]
usage:
  source: log
  seconds: 853
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 94
      output: 435
      cache_read: 6576124
      cache_write: 185647
      cost: 2.323
---
# T-2308 The context package primes the orchestrator for named reasons: it parses them, refuses an unknown one naming the known ones, and loads the role's conventions once and, for a thread, a draft, and a plan, their items and threads

## Work

ADR-0209 §3 and ADR-0217, criteria 1 and 2. In a new `reasons.go`, beside `ForStrategic` in `flai/internal/context/strategic.go`, add a reason-primed pack:

- Parse reasons as ADR-0217 §1 gives them: `<kind>:<id>` or a bare `<kind>`, the kinds being `thread:TH-nnnn`, `review:S-nnnn`, `planned:E-nnnn|S-nnnn`, `draft:S-nnnn`, `release`, and `ready`. Refuse an unknown kind, a kind missing its ID or given one it does not take, and an ID of the wrong type, naming the known kinds. The mapping of the trigger's events onto these kinds is flai serve's (T-2313), not the parser's.
- Load the role's conventions once, however many reasons there are, filtered by ADR-0217 §2. A section whose effective topics are `all` alone is kept. A section whose topics are `all` and other words is kept only when one of those words is a kind of the run's reasons, or `orchestration`. Every other pack keeps using `Matches` unchanged. Export the words the rule reads, the reason kinds and the built-in strategic roles' topics, for T-2316's `flai check` vocabulary.
- Per reason: for `thread`, the thread with its entries and the item it is on, whole; for `draft`, the draft story whole; for `planned`, the item planned, the items the planner wrote under it (its children), and the planner's open thread on it. An item named by two reasons loads once.
- Digest step: rank the `answer` and `recommendation` entries against the reasons' threads only, the threads of `thread` reasons and the planner threads `planned` loads. Do this within the budget's eighth (ADR-0194), reusing `answersFor`.
- No design briefs, no ranked sections, no open-issues table, and no catalog (ADR-0217 §5). Name the reasons in the pack's header and `--json`. The pack fits the budget and cuts into parts as every pack does.

`review`, `release`, and `ready` are parsed here and loaded by T-2309, through a hook this task leaves for it.

No `after`. It and T-2310 run together, since they share no path. This task's tests build their own conventions.

## Done when

- `reasons_test.go` has a test per reason this task loads, asserting what it loads. It also has a test for the refusal naming the known kinds, and one for the section rule: `all` alone kept, `all` plus a reason kept only for that reason, `all` plus `orchestration` kept for any reason.
- A test, skipped under `-short`, builds the general orchestrate pack over this repository. It builds a pack for one thread over the same conventions and design, with a fixture thread on a fixture story, and asserts the second is under a quarter of the first's size, printing both sizes.
- A test with two reasons on the same item loads the conventions and the item once.
- `flai test` passes on the changed paths.

## Notes

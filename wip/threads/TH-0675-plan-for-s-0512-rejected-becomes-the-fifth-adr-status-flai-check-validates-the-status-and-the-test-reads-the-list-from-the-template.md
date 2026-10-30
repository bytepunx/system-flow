---
id: TH-0675
title: "Plan for S-0512: rejected becomes the fifth ADR status, flai check validates the status, and the test reads the list from the template"
anchor:
  path: wip/kanban/stories/S-0512-search-test-ts-lists-adrs-with-the-supersession-chain-fails-on-main-since-adr-0202-s-status-became-rejected-a-status-outside-the-four-the-test-and-the-adr-template-name.md
  item: S-0512
status: resolved
participants: [planner-S-0512, orchestrator, flai serve]
created: 2026-10-28T08:29:03Z
updated: 2026-10-30T08:31:23Z
---

# TH-0675 Plan for S-0512: rejected becomes the fifth ADR status, flai check validates the status, and the test reads the list from the template

On wip/kanban/stories/S-0512-search-test-ts-lists-adrs-with-the-supersession-chain-fails-on-main-since-adr-0202-s-status-became-rejected-a-status-outside-the-four-the-test-and-the-adr-template-name.md.

## Entries

### 2026-10-28T08:29:03Z planner-S-0512
The story is enriched: 16 file touches, forecast 1h (flai forecast gave 28m; raised because the story spans an ADR, two Go packages with tests, the template, a vitest change, and four documents), cost of delay 12.50 USD a week from the operator's input of 5m per cycle. The reasons are under `### Planning` in the story's Notes.

Proposed solution, from I-0208's one instance: S-0501 already added `rejected` to the test, so the test passes on main today, but nothing validates an ADR's status, so the cause stays. The plan makes `rejected` the fifth ADR status, a proposed decision the operator declined in favour of another, immutable once set as an accepted one is; gives `flai check` the finding `adr.status` on a status outside the five, so that an unknown status stops a close-out before any vitest tier fails; and makes the test read its list from `design/adrs/0000-template.md`'s status comment, the one source.

Tasks and layers:

- Layer 1, together: T-2189 records the ADR; T-2190 adds one list of the five statuses in a package `check` and `docedit` can import (`adr` imports both) and the `adr.status` finding with a test; T-2191 names `rejected` in both ADR template copies and `documentation-standard.md`, with a template changelog entry and patch bump.
- Layer 2, together: T-2192 (after T-2190) makes `flai doc` refuse a change to a rejected ADR but its topics, built from the shared list; T-2193 (after T-2191) makes `search.test.ts` read its allowed statuses from the template's comment.
- Layer 3: T-2194 (after T-2189, T-2190, T-2191, T-2192) documents `adr.status` and the five statuses in `flai-cli.md` and `docs/users/flai.md`.
- Layer 4: T-2195 (after T-2192, T-2193, T-2194) closes I-0208 with `flai issue close --reason` and ticks the criteria.

Assumptions:

1. `rejected` is wanted as a status, since ADR-0202 already carries it on main (commit a554b18ff, in favour of ADR-0203), rather than recording a rejection as a superseding ADR. If the operator would rather keep four statuses, T-2189, T-2191, T-2192, and T-2194 change shape: say so here and the story's agent replans.
2. `flai adr new --status` keeps taking only `proposed` and `accepted`; a rejection is written on a proposed ADR, which is editable, and no `flai adr reject` command is added.
3. A rejected ADR is immutable like an accepted one. The alternative, editable like a proposed one, would let a declined decision be rewritten after the fact.
4. The new ADR's file is not in the touches: `flai adr new` allocates its number and `flai task done` widens the touches with it (ADR-0144).

Nothing to split, merge, or drop: the story has no prior tasks.

### 2026-10-28T08:29:26Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan with assumptions 1 to 3 as written. Make `rejected` the fifth ADR status, immutable once set, written by hand on a proposed ADR, with no `flai adr reject`.

## The plan fits

- **Tasks:** seven tasks in four layers. Layer 1 (T-2189, T-2190, T-2191) shares no path. T-2192 waits on T-2190's shared list, and T-2193 on T-2191's template comment. T-2194 documents what the layers before it build, and T-2195 closes I-0208 last.
- **Touches:** 16 files, every task at file level. The new ADR's file is widened by `flai task done` (ADR-0144).
- **Figures:** forecast 1h, kept over flai's 28m with the reason stated. Value 12.50 USD a week, from the 5m-a-cycle input the issue step gave it.

## Why it is yours

- A fifth ADR status that `flai check` enforces, and that `flai doc` refuses to change, is a tool-enforced rule. Under `decisions.md` § Rules it needs an ADR, and no recorded decision yet says a rejected ADR exists or is immutable.
- **Why I recommend it:** you set ADR-0202 to `rejected` yourself (a554b18ff, in favour of ADR-0203). The plan records what you already did, and makes the check catch the next unknown status before a vitest tier does. Making a rejected ADR immutable keeps a declined decision from being rewritten after the fact, as `decisions.md` keeps an accepted one.
- **The alternative:** keep four statuses and record a rejection as a superseding ADR. ADR-0202 then needs its status changed back, and T-2189, T-2191, T-2192 and T-2194 change shape.

The thread stays open for that answer. S-0512 stays a draft until then.

Source: design/conventions/decisions.md § Rules

### 2026-10-30T08:30:11Z flai serve
Confirmed the recommendation of 2026-10-28T08:29:26Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/conventions/decisions.md § Rules

### 2026-10-30T08:31:06Z orchestrator
The plan is approved, as confirmed after 48h under `orchestration.permissions.confirm_after`. Assumptions 1 to 3 stand: `rejected` becomes the fifth ADR status, immutable once set, written by hand on a proposed ADR, with no `flai adr reject`. T-2189 records that decision as an ADR. The seven tasks, 16 touches, the 1h forecast, and the 12.50 USD a week value fit the story. I resolve this thread and finalize S-0512, which `flai promote --drafts` lists as complete.

Source: design/conventions/decisions.md § Rules

### 2026-10-30T08:31:23Z orchestrator
Resolved: Plan approved: the orchestrator's recommendation was confirmed by flai serve after 48h (confirm_after); rejected becomes the fifth ADR status, immutable, recorded by T-2189.

---
id: TH-0611
title: "Plan of S-0429 revisited: five tasks in four layers stand, forecast 50m, cost of delay 76.92 USD a week, and two proposals"
anchor:
  path: wip/kanban/stories/S-0429-flai-digest-harvest-lists-what-the-archive-and-the-threads-hold-since-a-date-and-a-one-off-seed-fills-the-digest-from-the-operator-s-answers-and-the-closed-issues-so-there-is-something-to-prime-with.md
  item: S-0429
status: resolved
participants: [planner-S-0429, orchestrator]
created: 2026-10-23T16:44:37Z
updated: 2026-10-23T16:44:50Z
---

# TH-0611 Plan of S-0429 revisited: five tasks in four layers stand, forecast 50m, cost of delay 76.92 USD a week, and two proposals

On wip/kanban/stories/S-0429-flai-digest-harvest-lists-what-the-archive-and-the-threads-hold-since-a-date-and-a-one-off-seed-fills-the-digest-from-the-operator-s-answers-and-the-closed-issues-so-there-is-something-to-prime-with.md.

## Entries

### 2026-10-23T16:44:37Z planner-S-0429
planner-S-0429 revisited S-0429 on 2026-10-23, after E-0025's cost of delay input was set on TH-0476 and with S-0428 in progress. Nothing needs an answer; the two proposals at the end are the operator's to take or leave.

## Tasks, order, and layers

Unchanged from planner-E-0025's draft: no task was added, split, merged, or dropped, and no title, `## Work`, or `## Done when` was changed.

- Layer 1: T-1813 (`answer.go`, the operator-answer rule S-0433 reuses) and T-1814 (`sources.go`, narrative decisions and closed issues). No path in common, so they run together.
- Layer 2: T-1815 (`harvest.go`, `seed.go`), after both.
- Layer 3: T-1816 (`flai digest harvest` and `seed` in `flai/cmd/digest.go`, `flai-reference.md` regenerated), after T-1815.
- Layer 4: T-1817 (the seed run once here, entries committed and counted, `flai-cli.md`, `docs/users/flai.md`, `design/digest/README.md`), after T-1816.

## What changed

- Forecast 45m to 50m. flai gives 36m for size 20; the seed over 610 threads and 168 closed issues, and fixing what the check reports on its entries, is work the size does not count. Delivery 2026-10-24T06:15Z, flai's play-out moved by the raise.
- Cost of delay value 76.92 USD a week, from `flai cod`: 50m of E-0025's 6h30m, of its 600 USD a week.
- Touches: the 15 declared files stand, plus the `design/digest` folder on the story and T-1817, because the seed's file names are known only once it runs. The co-change list (operator docs, dashboard docs, the ADR index) was not added: no setting, dashboard change, or ADR is in scope.
- Notes added to T-1813, T-1814, and T-1815 with what the code has today: `Thread.Entries`, `Entry.Confirms`, `Thread.Opener`, `manifest.Manifest.Owner`, `topics.OfStory`; `workitem.ReadNarrative` and `NarrativeSection`, `Item.CompletedAt`; `threads.List` returning resolved threads too (there is no threads archive). The `### Planning` section of S-0429 was rewritten.

## Assumptions

1. S-0428 lands `flai/internal/digest` with a writer and a lister (the draft calls them `digest.New` and `digest.List`) and `flai digest new|edit|list|show`, as ADR-0191 on its branch says. S-0429 is `after: S-0428`, so its agent takes the names the package has then.
2. The operator is the manifest's `owner` (TH-0476 question 5, confirmed).
3. Decisions are harvested but never seeded (S-0429's notes); the curator takes them in.

## Proposals

1. **T-1814 `## Work`, closed issues.** An issue file has no cause or fix heading: its sections are `## Description`, `## Instances`, `## Remediation`, and `issues.Close` appends `Closed <ts>: <reason>`. I propose the bullet read: "each issue closed on or after `since`, as a `cause` candidate whose cause is its `## Description`, whose fix is the `Closed` line's reason (or `## Remediation` when the reason names none), whose date is that line's timestamp, and whose topics and paths are those of the story the reason names". The task's Notes already say this; the rewording is yours or the story's agent's, since planner-E-0025 wrote the task.
2. **The template's digest README.** Criterion 5 names this repository's `design/digest/README.md`. If the template's copy, `template/root/design/digest/README.md` (which S-0428 creates), should carry the same sentence about the seed, add it and `template/CHANGELOG.md` to T-1817's touches; I did not, since it widens the story to the template component.

### 2026-10-23T16:44:47Z orchestrator
Approved: the plan fits S-0429.

- **Tasks and layers.** Five tasks in four layers cover the five criteria. T-1813 and T-1814 share no path.
- **Touches.** The 15 touches name files. The `design/digest` folder is kept because the seed's file names are known only once it runs. That is a stated reason.
- **Forecast.** 50m, raised from flai's 36m. The reason given is the seed run over 610 threads and 168 closed issues.
- **Cost of delay.** 76.92 USD a week, from `flai cod` on E-0025's input.

**Proposal 1 is taken.** T-1814's closed-issue bullet takes the cause from `## Description` and the fix from the `Closed` line's reason, with `## Remediation` as the fallback. That is how issue files are written, and it changes no scope. The task's Notes already say it, and the story's agent rewords the Work when it reviews the tasks.

**Proposal 2 is left out.** Criterion 5 names this repository's README. Adding the template's copy and `template/CHANGELOG.md` would widen the story to the template component, and its scope stays as written. `docs/users/flai.md` describes the seed for every project.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T16:44:50Z orchestrator
Resolved: Plan approved under plan_backlog_stories; proposal 1 taken for the story's agent, proposal 2 left out.

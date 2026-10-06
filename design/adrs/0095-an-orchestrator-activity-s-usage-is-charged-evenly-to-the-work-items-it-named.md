---
id: ADR-0095
title: "An orchestrator activity's usage is charged evenly to the work items it named and the items above them, and what no item carries is its kind's project strategic total"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0079, ADR-0083]
topics: [cli, orchestration, planning, analysis]
---

# ADR-0095 An orchestrator activity's usage is charged evenly to the work items it named and the items above them, and what no item carries is its kind's project strategic total

## Context

ADR-0083 charges a planner activity's usage to the item it planned and every item above it, under `usage.strategic` per kind, and leaves room for the orchestrator and the analyzer. The orchestrator (ADR-0087) runs for the whole project and makes many decisions in one run: it promotes and orders stories, answers threads, accepts stories, and publishes. Each activity it logs, through `activity_log` or as its run ends (ADR-0079), names the items it touched, or none. Without a charge, all of its cost stays in `wip/agents/orchestrator.md` and none on the items it was spent for. The designer asked (S-0226) that each activity be charged to the items the decision concerned, or to a project total when it concerned none, so that none of its cost is lost, and that `flai stats` report that total per kind beside the per-item figures.

## Decision

An orchestrator activity's usage is split evenly between the work items it named and summed up the hierarchy from each, and what a kind's activity document accrued that no item carries is that kind's project strategic total, which `flai stats` reports beside what the items carry.

1. **Charged to the items it named.** When flai serve logs an orchestrator activity, the usage apportioned to its span, as ADR-0083 apportions a planner's, is split evenly between the work items the entry names: each epic, story, and task that exists, archived ones included, each named once. Anything else it names, such as a thread or an issue, takes no share. Each share's tokens and seconds are whole numbers, the remainder going to the first items named, so that the shares add up to the whole; each share's cost is the cost over the count. Each share is added under the `orchestrator` entry of its item and of every item above it, up to its epic, at once, as ADR-0083 sums a planner's charge. A task and its story both named each take a share, and the story carries the task's too.
2. **The project strategic total.** What a kind's activity document accrued that no item carries is that kind's project strategic total: its `accrued_cost` and `accrued_seconds` less the `strategic` cost and seconds of that kind on the items at the top of the hierarchy, those with no parent, which carry every charge made below them. It holds the activities that named no work item, the planner's activities with no planned item, the analyzer's activities until S-0227 charges them, those logged before their kind was charged to items, and any charge that failed. It is worked out where it is reported, not written: nothing can make it drift from the document.
3. **Reported per kind.** `flai stats` reports, for each strategic agent kind beside its document's totals, what the items carry (`items`) and the project total (`project`), cost and seconds each, so that the two add up to the document's totals, to a hundredth of a cent. `design/system/metrics.md` defines them.
4. **The planner unchanged.** A planner activity is still charged to the item its run planned (ADR-0083), not to the items its entry names.

## Consequences

- None of the orchestrator's cost is lost: it is on the items it concerned or in its project total, and the two add up to `wip/agents/orchestrator.md`.
- An orchestrator decision named on many items puts a small share on each; the share is an apportionment, marked estimated like every strategic entry.
- Costs on items are rounded to a hundredth of a cent each, so the project total can differ from the exact remainder by that much per share; it is reported to a hundredth of a cent and never below zero.
- The project total of the planner holds what it spent before S-0225 charged it to items, so `flai stats` accounts for it with no migration.
- An item named by mistake is charged; the orchestrator names the items a decision concerned, as its agent definition says.
- `activity_log` and `flai activity log` say what the activity charged, the items or the project total.

## Alternatives considered

- **Write the project total in the activity document's front matter**, accrued as an activity charging no item is logged. It would need a migration for the activities logged before it, raise `flai.minimum` because front matter is parsed strictly, and could drift from the document's totals whenever a charge failed.
- **Record the project total on the project manifest or an item of its own.** The manifest is configuration, and an item would appear on the board and in every listing.
- **Charge the whole cost to each item named.** The items would add up to more than was spent.
- **Charge only the first item named.** The others would carry nothing of a decision that concerned them all.

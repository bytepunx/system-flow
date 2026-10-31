---
id: T-2388
type: task
nature: improvement
title: "Write the ADRs and the design: a trial merge replaces the overlap hold, after names a read dependency with its reason, and the in-progress share is a metric"
status: backlog
parent: S-0551
owner: alex
created: 2026-10-31T07:42:42Z
updated: 2026-10-31T07:42:42Z
transitions: []
stream: S-0551
tags: [flai, design]
touches: [design/adrs/README.md, design/system/workflow.md, design/system/strategic-agents.md, design/system/metrics.md]
---
# T-2388 Write the ADRs and the design: a trial merge replaces the overlap hold, after names a read dependency with its reason, and the in-progress share is a metric

## Work

Write two ADRs with `flai adr new`, and the design they change. It waits for nothing: the other tasks build what it records, and the story's goal states the decision already.

- An ADR that refines [ADR-0046](../../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md). A ready story whose claim overlaps an in-progress story's is started when a trial merge of the in-progress story's branch at the ready story's base, main's head, is clean, and held only when it is not, with the conflicting paths in the hold's reason and the serve journal. Say what holds when the trial merge cannot run (git older than 2.38, no branch yet): the overlap hold as today. Say whether a story in review still holds nothing (ADR-0096). `after` stays honoured, but names only a story whose output this story reads, with the reason in the story's Notes naming that story; `flai check` advises of an `after` with no reason.
- An ADR for the metric, since `design/system/metrics.md` changes only with one: the share of the window with n stories in progress at once, for n from 0 to the in-progress limit, and the ready wait as queue time at p50 and p85 ([ADR-0190](../../../design/adrs/0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md)).
- `design/system/workflow.md`: the trial merge in place of the overlap hold, and the `after` rule at its section on `after`.
- `design/system/strategic-agents.md`: the planner's `after` rule, near its line on `after`.
- `design/system/metrics.md`: the new aggregate and its JSON key, under `## Aggregates` or `### Claims and touches`.
- `design/adrs/README.md`: both entries, as `flai adr new` writes them.

## Done when

- Both ADRs are written and indexed, and the first names ADR-0046 as the decision it refines.
- `workflow.md`, `strategic-agents.md`, and `metrics.md` say what the ADRs decide.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes

ADR-0046 is accepted and immutable: the new ADR refines it by reference and does not edit it.

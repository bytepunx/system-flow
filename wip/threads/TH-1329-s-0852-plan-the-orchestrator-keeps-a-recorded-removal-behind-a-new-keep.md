---
id: TH-1329
title: "S-0852 plan: the orchestrator keeps a recorded removal behind a new keep_removals permission, eight tasks in six layers"
anchor:
  path: wip/kanban/stories/S-0852-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md
  item: S-0852
status: resolved
participants: [planner-S-0852, claude-concierge]
created: 2026-11-27T23:59:14Z
updated: 2026-11-28T00:01:51Z
---

# TH-1329 S-0852 plan: the orchestrator keeps a recorded removal behind a new keep_removals permission, eight tasks in six layers

On wip/kanban/stories/S-0852-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md.

## Entries

### 2026-11-27T23:59:14Z planner-S-0852
Plan for S-0852 (I-0415). Eight tasks in six layers. There was no earlier thread on S-0852 to answer, so this is the story's first plan.

## Recommendation for you to confirm first

ADR-0155 makes `--keep-removal` the operator's alone. It rejected letting the orchestrator keep a removal. I recommend refining it as follows:

- Add `orchestration.permissions.keep_removals`. It is off by default and needs `accept_reviews`.
- With it on, the orchestrator may keep a recorded removal through `flai accept --by orchestrator --keep-removal`, the MCP tool `accept`, or `queue_add`.
- It may keep only a path that the story's narrative names under `## Decisions`, and it must give a reason. The narrative entry is the story's agent's own record that the removal is a rewrite.
- The reason is written under `### Removals kept` with `by: orchestrator`.

With the permission off, nothing changes. If you would rather keep ADR-0155 as it stands, say so: the code tasks then drop, and the story closes I-0415 as a recorded decision.

## Tasks and layers

1. T-3618: an ADR proposing the rule above, listed in `design/adrs/README.md`.
2. T-3619: the manifest reads the permission, and the settings index lists it. Waits for T-3618.
3. T-3620: `flai accept --by orchestrator` takes the keep, with a test that reproduces I-0415 (S-0811's rewrite). Waits for T-3619.
4. These two share no file and run together. Both wait for T-3620.
   - T-3621: the MCP tool `accept` takes `keep_removals`, and the guard's text changes.
   - T-3622: `queue_add` records the keeps, eligibility leaves out kept paths, and the batch landing passes them.
5. These two run together. Both wait for T-3621 and T-3622.
   - T-3624: the orchestrator's role file (baseline and both copies) and the strategic-agents convention.
   - T-3626: the design, the user and operator guides, and `template/CHANGELOG.md`.
6. T-3629: close I-0415. Waits for T-3624 and T-3626.

## Assumptions

- The queue is in scope. This project batches, and S-0811 was ineligible there on `undone_main`, so a fix to `accept` alone would not have moved it.
- The dashboard's `accept.run` stays operator-only.
- The ADR number 0322 in T-3618's touches is a guess.
- Forecast is 1h13m, as `flai forecast` gives it.
- Cost of delay is 22200 USD a week, from your 148h input. I left the input as it is.

### 2026-11-28T00:01:51Z claude-concierge
For the operator, who delegated completing drafts and queueing ready work in this session: the plan stands, with keep_removals off by default so ADR-0155's rule holds until the operator turns the permission on. The operator was told this refines ADR-0155 and can still cancel the story or reject the ADR at review.

Resolved: plan reviewed and the story finalized on the operator's delegation

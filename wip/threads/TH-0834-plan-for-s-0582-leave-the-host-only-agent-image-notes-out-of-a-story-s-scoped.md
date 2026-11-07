---
id: TH-0834
title: "Plan for S-0582: leave the host-only agent.image notes out of a story's scoped check"
anchor:
  path: wip/kanban/stories/S-0582-flai-check-finds-agent-image-outside-the-story-at-close-out.md
  item: S-0582
status: resolved
participants: [planner-S-0582, orchestrator]
created: 2026-11-07T00:31:59Z
updated: 2026-11-07T00:33:04Z
---

# TH-0834 Plan for S-0582: leave the host-only agent.image notes out of a story's scoped check

On wip/kanban/stories/S-0582-flai-check-finds-agent-image-outside-the-story-at-close-out.md.

## Entries

### 2026-11-07T00:31:59Z planner-S-0582
S-0582 is planned: four tasks in three layers, touches, a 45m forecast, and a cost of delay of 25 USD a week.

## Tasks

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-2493 | none | ADR: the scoped check leaves out the three `agent.image` findings only the host clears |
| 2 | T-2494 | T-2493 | `image.go` marks them; `ScopeToStories` drops them unless the story changes `system-flow.yaml` or the layer; tests reproduce I-0264 |
| 2 | T-2495 | T-2493 | `docs/users/flai.md`, `flai-cli.md`, `workflow.md`, both `work-management.md`, `template/CHANGELOG.md` |
| 3 | T-2496 | T-2494, T-2495 | `flai issue close I-0264 --reason` |

T-2494 and T-2495 share no file, so they run together.

## Assumptions

- The fix follows ADR-0171 (S-0446) for `components.examples`: leave the finding out of the scoped run, keep it in an unscoped `flai check`.
- The not-built warning is left out too, not only the two no-docker notes: on a host without the layer image, it would be recorded the same way.
- The warning for a layer that is not a file stays: a story can fix it.
- The ADR is numbered 0255 in T-2493's touch; the story keeps the folder touch `design/adrs` in case another story takes that number first.
- No story `after`: S-0582 reads nothing another open story writes.

## Figures

- Forecast 45m: `flai forecast` gave 31m; raised to S-0446's 47m, the same shape.
- Cost of delay 25 USD a week, from the 10m input set on TH-0833.

Nothing to split, merge, or drop.

### 2026-11-07T00:32:54Z orchestrator
Approved as planned. The plan fits S-0582:

- **Tasks:** T-2493 puts the fix in an ADR before any code, as the goal asks. T-2494 builds it and adds a test that reproduces I-0264, which is criterion 1. T-2496 closes I-0264, which is criterion 2. T-2495 updates the docs and the close-out rule, as `documentation.md` asks.
- **Layers:** T-2494 and T-2495 share no file, so they can run together after T-2493.
- **Touches:** the folder touch `design/adrs` is named under `### Planning`, as ADR-0188 asks, and it covers the ADR number another story might take first.
- **Figures:** a 45m forecast, based on S-0446, the same shape of fix. A value of 25 USD a week from the 10m input set on TH-0833.

The fix follows ADR-0171's precedent for `components.examples`. No `after` is needed.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T00:33:04Z orchestrator
Resolved: Plan approved under plan_backlog_stories: four tasks, T-2493 to T-2496, in three layers; touches, a 45m forecast, and a value of 25 USD a week all fit S-0582.

---
id: TH-0673
title: "Plan for S-0510: five tasks in three layers, figures written"
anchor:
  path: wip/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md
  item: S-0510
status: resolved
participants: [planner-S-0510, orchestrator]
created: 2026-10-28T08:06:46Z
updated: 2026-10-28T08:07:03Z
---

# TH-0673 Plan for S-0510: five tasks in three layers, figures written

On wip/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md.

## Entries

### 2026-10-28T08:06:46Z planner-S-0510
Revisited the plan planner-E-0029 wrote with the epic. The five tasks, their touches, and their `after` stand; the cost of delay value is now written from the operator's input on E-0029 (TH-0657).

Tasks and layers:

1. Layer 1, run together, no path in common: T-2164 (`agent.image.layer` read by `flai/internal/manifest`, with its rows in `project-manifest.md` and `settings.md`, the template's comment, and the changelog line) and T-2165 (`scripts/agent-image.sh --layer`, `make agent-images`, this repository's layer `flai/images/system-flow/Dockerfile`, and `system-flow.yaml` naming it). The manifest reader ignores unknown keys, so T-2165 can name the key before T-2164 lands.
2. Layer 2, run together: T-2171 (the `agent.image` check in `flai/internal/check/image.go`) and T-2172 (the container adapter runs the layer's tag in place of the harness image), each after T-2164 for the key.
3. Layer 3: T-2176 (one Go-changing story of this repository worked end to end on `harness: container` with the layer, an issue per limitation met, and the five documents), after T-2165, T-2171, and T-2172.

Figures: touches unchanged, 24 files and the folder `design/issues` (kept for the issue files T-2176's `flai issue new` writes, which no task can name); a second `flai touches suggest` on the 25 touches names only paths outside the story, each left out with its reason under `### Planning`. Forecast kept at 2h30m over flai's 56m, delivery 2026-10-30T02:47Z from flai's play-out. Cost of delay 27.78 USD a week, this story's share of E-0029's 150 by forecast duration, as `flai cod` gives it.

Assumptions:

- `after` names S-0508 alone: it lands the adapter, and through its own `after` S-0507 and S-0505, which land pilot, the base images, `scripts/agent-image.sh`, and `make agent-images`. Nothing else is named: the pull order sets the rest.
- `scope.go` needs no change for the new check: its findings fall on `system-flow.yaml`, which the generic scope rule marks inside or outside by the branch's diff; the `docker`-missing note is advisory and fails nothing.
- `TestSettingsIndexIsComplete` walks the `Manifest` struct's yaml tags, so T-2164's row in `settings.md` is enough and `manifest/settings.go`, the `flai manifest set` catalog, gains no entry: the key is not settable from the dashboard.
- The story T-2176 works on `harness: container` is one the story's agent picks when it gets there, small and Go-changing, with the host's container settings from S-0508 already set; it is not named here.
- S-0508's planner proposed (TH-0666) that the container harness skip ADR-0106's `claude` version check, with the image's version left to S-0510 or an issue. This story does not take it up: reading `claude --version` from the image is the adapter's, and an issue from S-0508 would make its own story.

No task to split, merge, or drop.

### 2026-10-28T08:06:57Z orchestrator
Approved: the plan fits S-0510.

- **Tasks:** five tasks in three layers. T-2164 and T-2165 share no path, and neither do T-2171 and T-2172. Both of those wait on T-2164 for the key. T-2176 waits for the three it proves and documents.
- **Touches:** file-level on every task. The folder `design/issues` is kept for the issue files T-2176 writes, and the story's `### Planning` names it.
- **Figures:** forecast 2h30m, kept over flai's 56m with the reason stated. Value 27.78 USD a week, flai's share of E-0029's 150 USD input.
- **Order:** `after: [S-0508]` names the story whose adapter it needs on main.
- **ADR-0106's version check:** leaving it to S-0508 is right. That question is TH-0666's, and it waits on the operator there.

Nothing waits on an answer, so I resolve the thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T08:07:03Z orchestrator
Resolved: Plan approved under plan_backlog_stories: five tasks in three layers, file-level touches, forecast 2h30m, value 27.78 USD a week; nothing waits on an answer.

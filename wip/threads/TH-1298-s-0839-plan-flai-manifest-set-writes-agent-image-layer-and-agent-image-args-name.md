---
id: TH-1298
title: "S-0839 plan: flai manifest set writes agent.image.layer and agent.image.args.<NAME>"
anchor:
  path: wip/kanban/stories/S-0839-no-flai-command-writes-agent-image-args-so-a-project-s-layer-build-arguments.md
  item: S-0839
status: resolved
participants: [planner-S-0839, orchestrator]
created: 2026-11-17T12:32:47Z
updated: 2026-11-17T12:35:11Z
---

# TH-1298 S-0839 plan: flai manifest set writes agent.image.layer and agent.image.args.<NAME>

On wip/kanban/stories/S-0839-no-flai-command-writes-agent-image-args-so-a-project-s-layer-build-arguments.md.

## Entries

### 2026-11-17T12:32:47Z planner-S-0839
**Plan.** S-0839 has four tasks in three layers. The remedy it proposes, from I-0364's one instance: `flai manifest set agent.image.args.<NAME>=<value>` and `agent.image.layer=<path>`, with `--unset`. The manifest's existing image validation checks them: no `BASE`, an `ARG` name, one line, and no args without a layer.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3574: the two settings in `SettingFor`, and their `EditSettings` tests | `flai/internal/manifest/settings.go`, `settings_test.go` | none |
| 2 | T-3577: the help and an example, a command test that reproduces I-0364, and the regenerated reference | `flai/cmd/manifest.go`, `manifest_test.go`, `docs/users/flai-reference.md` | T-3574 |
| 2 | T-3581: the settings index, `project-manifest.md`, `flai-cli.md`, `flai.md`, the template's manifest comment, and its changelog | six documents | T-3574 |
| 3 | T-3582: `flai issue close I-0364 --reason` | the issue and `summary.md` | T-3577, T-3581 |

T-3577 and T-3581 share no path, so they run together.

**Figures.**

- Forecast: 1h, delivery 2026-11-17T21:10Z. flai gave 28m; it is raised for four task closes, a regenerated reference, and a close-out of about 10m.
- Cost of delay: 25 USD a week, from flai's inputs, unchanged.
- Tags: `flai` and `template`.

**Assumptions.**

1. `agent.image.layer` is in scope: the issue names it as the same exception, and args are refused without a layer.
2. The dashboard's settings page, `projects[].install` (D-0593), and `flai agent set`/`clear` are out of scope. The `tests` tiers the issue cites through D-0449 are written by `flai manifest set` already.
3. No `after` on another story: S-0839 reads no other open story's output.

Proposed: nothing to split, merge, or drop.

### 2026-11-17T12:34:32Z orchestrator
Approved. The plan fits S-0839.

- T-3574 adds `agent.image.layer` and `agent.image.args.<NAME>` to `flai manifest set` and `--unset`, checked by the manifest's existing image validation.
- T-3577 adds a command test that reproduces I-0364. That meets criterion 1.
- T-3582 closes I-0364. That meets criterion 2.
- T-3581 covers the settings index and the user docs, as `documentation.md` asks for a changed setting. It also covers the template's comment and its changelog.
- The touches are 13 files and no folder. Each task's touches are inside them.
- Layer 2's two tasks share no path, so they can run together.
- The forecast of 1h and the value of 25 USD a week are flai's figures, with the planner's stated reason for the forecast.

Source: design/conventions/documentation.md § Project additions

### 2026-11-17T12:35:11Z orchestrator
Resolved: The orchestrator approved S-0839's plan under plan_backlog_stories.

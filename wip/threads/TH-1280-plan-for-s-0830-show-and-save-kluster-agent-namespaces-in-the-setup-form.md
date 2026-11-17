---
id: TH-1280
title: "Plan for S-0830: show and save kluster.agent_namespaces in the setup form"
anchor:
  path: wip/kanban/stories/S-0830-the-dashboard-s-kluster-setup-form-does-not-show-kluster-agent-namespaces.md
  item: S-0830
status: resolved
participants: [planner-S-0830, orchestrator]
created: 2026-11-17T08:29:45Z
updated: 2026-11-17T08:33:40Z
---

# TH-1280 Plan for S-0830: show and save kluster.agent_namespaces in the setup form

On wip/kanban/stories/S-0830-the-dashboard-s-kluster-setup-form-does-not-show-kluster-agent-namespaces.md.

## Entries

### 2026-11-17T08:29:45Z planner-S-0830
The plan fixes the form in flaiover alone. flai already sends the setting to the dashboard, and already saves and checks it.

## Proposed fix

`kluster.status` already lists `kluster.agent_namespaces`: kind `number`, whole, default `2`, its value a JSON number. `flai manifest set` and `settings.manifest` already write it and refuse a bad value. The dashboard drops it in three places:

- `PLAIN_KEYS` in `flaiover/src/lib/components/DeploySettings.svelte` leaves it out.
- `typed()` shows a value that is not text as empty.
- `KlusterSetting` in `flaiover/src/lib/server/repo.ts` has no `number` kind.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3526 The kluster setup form shows and saves kluster.agent_namespaces | `DeploySettings.svelte`, its test, `repo.ts`, `DeploySettings.examples.json` | none |
| 2 | T-3528 The design and the guides say the setup form sets kluster.agent_namespaces, and I-0356 is closed | `design/system/kluster.md`, `docs/operators/settings.md`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` | T-3526 |

T-3528 waits for T-3526 because it describes the field as T-3526 builds it. T-3528 also closes I-0356, which is criterion 2.

## Figures

- Forecast 25m, as `flai forecast` gives it. It stands.
- Cost of delay 25.00 USD a week, from the inputs flai set from I-0356. It stands.

## Assumptions

- No Go change is needed. The hostapi tests already list the key in `kluster.status`'s settings.
- A numeric input that saves the trimmed text is enough. flai parses the text and refuses anything that is not a whole number of 0 or more. The form shows that refusal beside the field, as it does for the other fields.
- The Deployments page's Setup pane reuses `DeploySettings`, so its own test file is not touched.
- No split, merge, or drop is proposed. Both tasks are small.

### 2026-11-17T08:32:32Z orchestrator
Approved. The plan fits S-0830:

- T-3526 shows and saves the setting, with tests that fail without it, which meets criterion 1. T-3528 brings the four documents in step and closes I-0356, which meets criterion 2. T-3528 waits for T-3526 because it describes the field T-3526 builds.
- flai already sends, saves, and checks `kluster.agent_namespaces`, so a fix in flaiover alone is the right scope.
- Every touch is a file one task names. There is no folder touch and no `after`.
- The 25m forecast comes from `flai forecast` and the 25.00 USD a week value from `flai cod`.

I am finalizing S-0830 under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T08:33:40Z orchestrator
Resolved: Plan approved under plan_backlog_stories.

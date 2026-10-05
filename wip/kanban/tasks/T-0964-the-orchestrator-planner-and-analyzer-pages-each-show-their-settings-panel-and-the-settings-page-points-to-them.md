---
id: T-0964
type: task
nature: feature
title: The orchestrator, planner, and analyzer pages each show their settings panel, and the Settings page points to them
status: backlog
parent: S-0229
owner: alex
created: 2026-10-05T05:47:02Z
updated: 2026-10-05T05:47:02Z
transitions: []
stream: S-0229
tags: [dashboard]
touches: [flaiover/src/routes/workflow/orchestrator, flaiover/src/routes/workflow/analyzer, flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts]
after: [T-0963]
---
# T-0964 The orchestrator, planner, and analyzer pages each show their settings panel, and the Settings page points to them

## Work

The first criterion puts each agent's settings where its page is. The orchestrator and analyzer pages are S-0228's, at `/workflow/orchestrator` and `/workflow/analyzer`; the planner page is S-0259's, at `/workflow/planner`.

- `/workflow/orchestrator`: a **Settings** section with T-0963's `StrategicSettings` for `orchestration`: the permissions, each with its sentence and its risk, then `policy`, then `release` (its kind, thresholds, theme, and `whole_epics`).
- `/workflow/planner`: the same for `planning`: `agent`, `replan`, `schedule`, `hour_rate`, `cycle`, `default_duration`.
- `/workflow/analyzer`: the same for `analysis`: `agent`, `schedule`.
- Each page asks `settings.get` with its other data and again after a save, so what it shows is what the manifest now holds.
- `SettingsPanel.svelte`: the read-only "Planning again, for this project" section says the keys are edited on the Planner page and links to it, instead of telling the operator to edit `system-flow.yaml` by hand.

This task waits for T-0963, for the component it places. It assumes S-0228 is done, so the orchestrator and analyzer pages exist; if they do not, stop and ask on the story.

## Done when

- each page's test renders its panel with the block's keys from a fixture `settings.get`
- each page's test saves a key, and one renders the panel read-only while `settings` is off
- the Settings page's test finds the link to the Planner page
- `npm test` and `npm run check` in `flaiover` pass

## Notes

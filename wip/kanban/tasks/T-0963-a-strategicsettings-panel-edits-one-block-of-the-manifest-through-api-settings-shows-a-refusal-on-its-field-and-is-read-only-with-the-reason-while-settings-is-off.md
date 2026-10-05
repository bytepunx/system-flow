---
id: T-0963
type: task
nature: feature
title: A StrategicSettings panel edits one block of the manifest through /api/settings, shows a refusal on its field, and is read-only with the reason while settings is off
status: backlog
parent: S-0229
owner: alex
created: 2026-10-05T05:46:51Z
updated: 2026-10-05T05:46:51Z
transitions: []
stream: S-0229
tags: [dashboard]
touches: [flaiover/src/lib/settings.ts, flaiover/src/routes/api/settings/+server.ts, flaiover/src/routes/api/settings/settings.test.ts, flaiover/src/lib/components/StrategicSettings.svelte, flaiover/src/lib/components/StrategicSettings.svelte.test.ts]
after: [T-0961]
---
# T-0963 A StrategicSettings panel edits one block of the manifest through /api/settings, shows a refusal on its field, and is read-only with the reason while settings is off

## Work

One component serves the three pages, so the save, the refusal, and the read-only state are written and tested once.

- `flaiover/src/lib/settings.ts`: add `manifest` to `SETTINGS_KINDS` (a project setting, not in `HOSTWIDE`), and types for `settings.get`'s `strategic` block from T-0961.
- `flaiover/src/routes/api/settings/+server.ts`: pass `settings.manifest`'s refusal through as 422 with each `{field, reason}`, and its `Disabled` as 403, as the route does for the other kinds.
- `flaiover/src/lib/components/StrategicSettings.svelte`: given a block (`orchestration`, `planning`, or `analysis`) and the `strategic` answer, it shows each key with its sentence, and a permission with its risk beside it; an input by kind (a switch, a choice, a number, a duration, a cron expression or `daily`, an agent's harness, model, and config); and the default when unset. **Save** posts only the keys that changed, an emptied one as an unset. A refusal keeps what was typed and shows each reason under its field. While the `settings` action is off every input is disabled, and the panel says why and gives `flai serve enable settings`.

This task waits for T-0961, for the write and the read it calls.

## Done when

- a component test saves a changed key and posts only that key
- a component test shows a refusal's reason under its field and keeps the input
- a component test renders the panel read-only, with the reason and the command, while `settings` is off
- a route test passes a 422 with its fields and a 403 through
- `npm test` and `npm run check` in `flaiover` pass

## Notes

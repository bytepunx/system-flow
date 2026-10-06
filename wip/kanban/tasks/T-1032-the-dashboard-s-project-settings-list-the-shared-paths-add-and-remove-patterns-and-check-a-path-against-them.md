---
id: T-1032
type: task
nature: feature
title: The dashboard's project settings list the shared paths, add and remove patterns, and check a path against them
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:16:27Z
updated: 2026-10-06T12:16:41Z
transitions: []
stream: S-0295
tags: [flaiover]
touches: [flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/server/agent.ts]
after: [T-1031]
---
# T-1032 The dashboard's project settings list the shared paths, add and remove patterns, and check a path against them

## Work

Add a shared paths section to the project settings in `flaiover/src/lib/components/SettingsPanel.svelte`, with its kind and types in `flaiover/src/lib/settings.ts`. Add `settings.shared` and the check method T-1031 adds to the allowed methods in `flaiover/src/lib/server/agent.ts`. The section:

- lists the patterns from `settings.get`, each with a remove button;
- adds a pattern from a text field, showing flai's refusal inline when the pattern is invalid or already listed;
- explains the glob dialect in one line beside the field: `*` within a folder, `**` across folders, and a plain path covering everything under it. This makes it easy for the operator, as TH-0180 asks;
- checks a typed path, or a story ID, through T-1031's check method, and shows which pattern matched or that none did.

Gate the writes by the settings action like the panel's other project settings, read-only when the action is off. The check stays available, since it changes nothing.

This task waits for T-1031, whose methods it calls.

## Done when

- `SettingsPanel.svelte.test.ts` covers listing, adding, removing, a refused pattern shown inline, the read-only state with the action off, and a path and a story checked.
- `npm run check`, `npm run lint`, and the flaiover tests pass, as `scripts/` runs them for flaiover.

## Notes

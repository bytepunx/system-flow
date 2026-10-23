---
id: T-2097
type: task
nature: improvement
title: Add a Sub-projects section to the dashboard's Projects configuration page that lists, adds, edits, and removes a project's sub-projects through settings.subproject
status: done
parent: S-0491
owner: alex
created: 2026-10-23T20:40:52Z
updated: 2026-10-23T22:07:02Z
transitions:
  - to: ready
    at: 2026-10-23T21:50:50Z
    by: agent-S-0491
  - to: in-progress
    at: 2026-10-23T21:50:51Z
    by: agent-S-0491
  - to: done
    at: 2026-10-23T22:07:02Z
    by: agent-S-0491
stream: S-0491
tags: [dashboard]
touches: [flaiover/src/lib/settings.ts, flaiover/src/lib/settingspages.ts, flaiover/src/lib/settingspages.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/routes/api/settings/+server.ts, flaiover/src/routes/api/settings/settings.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-2094]
usage:
  source: log
  seconds: 971
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 85
      output: 34893
      cache_read: 5223930
      cache_write: 203288
      cost: 2.8593
---
# T-2097 Add a Sub-projects section to the dashboard's Projects configuration page that lists, adds, edits, and removes a project's sub-projects through settings.subproject

## Work

The Projects page under Configuration (`settingspages.ts`, slug `projects`) shows the section `projects`, the projects flai serves on the host, and `tokens`. Add a section `sub-projects` between them, headed "Sub-projects of this project", so the two are not confused:

- `settingspages.ts`: `'sub-projects'` in `SETTINGS_SECTIONS` and in the Projects page's `sections`; `settingspages.test.ts` follows (the sections count and the page's anchors).
- `settings.ts`: the kind `subproject` in `SETTINGS_KINDS`, per project (not in `HOSTWIDE`), and the type of a sub-project as `project.info` answers it (`repo.ts` § `Manifest.projects`: name, path, kind, tags); `+server.ts` passes the kind through as it passes `shared` and `manifest`, and `settings.test.ts` shows it.
- `SettingsPanel.svelte`: the section lists each sub-project with its name, path, kind, and tags; a form adds one (name, path defaulting to the name, kind as a text field with the kinds flai knows suggested, tags), saying that a folder the project lacks is created with a README; a row edits its keys or removes the entry, with the same confirm the shared paths use; each call is `change('sub-projects', 'subproject', {op, name, ...})`, and a refusal shows flai's text verbatim as the other sections do. The section is read-only, with the reason shown, while the `settings` host action is off for the project, as the shared paths section is.
- Tests in `SettingsPanel.svelte.test.ts` follow the shared-paths block: the list from `project.info`, an add that posts `op: add`, an edit, a remove, a refusal shown, and the read-only state.
- `docs/users/flaiover.md` § Settings › Projects gains the section; `design/system/flaiover-dashboard.md` § Writes lists the new write beside the settings item.

It waits for T-2094, whose method it calls; it runs beside the orchestrator task, which touches other files.

## Done when

- The Projects configuration page shows a Sub-projects section listing the manifest's `projects[]` entries with name, path, kind, and tags.
- Adding, editing, and removing a sub-project posts `kind: subproject` with `op` `add`, `set`, or `remove` to `/api/settings`, which writes `settings.subproject`, and a refusal is shown with flai's text, each shown by a test.
- With the `settings` host action off, the section is read-only and says why.
- `pnpm test` in `flaiover` passes, `scripts/flai.sh test` passes for the paths changed, and `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` describe the section.

## Notes

Precedent: the shared paths section of `SettingsPanel.svelte` (S-0295) and its tests.

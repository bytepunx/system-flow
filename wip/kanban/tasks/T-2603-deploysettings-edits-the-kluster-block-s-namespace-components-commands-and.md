---
id: T-2603
type: task
nature: feature
title: DeploySettings edits the kluster block's namespace, components, commands, and timeout, and asks an agent to write the scripts
status: ready
parent: S-0595
owner: alex
created: 2026-11-07T22:07:32Z
updated: 2026-11-07T22:11:01Z
transitions:
  - to: ready
    at: 2026-11-07T22:10:11Z
    by: agent-S-0595
stream: S-0595
tags: [dashboard]
touches: [flaiover/src/lib/components/DeploySettings.svelte, flaiover/src/lib/components/DeploySettings.svelte.test.ts, design/components/flaiover/src/lib/components/DeploySettings.examples.json, flaiover/src/lib/settings.ts, flaiover/src/lib/settings.test.ts]
after: [T-2601]
---
# T-2603 DeploySettings edits the kluster block's namespace, components, commands, and timeout, and asks an agent to write the scripts

## Work

ADR-0263 §1 and §5. Extract the Klusters page's step 4, "Deploying stories" (`flaiover/src/routes/ops/klusters/+page.svelte`, its `DEPLOY_KEYS`, `settingField`, `save('deploy', …)`, and `writeWithAgent`), into a component `flaiover/src/lib/components/DeploySettings.svelte`, without changing either page yet. Props: `settings` (the `kluster.status` settings list), `settingsOn`, `settingsEnable`, `deployable`, what the agent story's epic picker needs, and a `saved` callback taking the new status. It edits `kluster.namespace`, `kluster.components`, `kluster.build`, `kluster.deploy`, `kluster.teardown`, and `kluster.timeout` and saves them through `POST /api/settings` `{kind:'manifest', set, unset}` with each refusal beside its field. `kluster.components` is edited as rows, a name and its build command, one argument a line, with add, remove, and move up/down, saved as the JSON list; `build` and `components` are offered as a choice (one command, or one per component), since a block takes one. With `settings` off it shows the one host step (`settingsEnable`) and disables the fields. Keep **Write these with an agent** as it is, its story now also naming the components. Put the parse and print helpers for the components value in `flaiover/src/lib/settings.ts` beside `argv` and `lines`, tested in `settings.test.ts`. UI from the design tokens only. Write its examples file. Waits for T-2601, which adds the `components` kind to the types.

## Done when

- `DeploySettings.svelte.test.ts` covers saving commands, saving components as the JSON list, removing a component, switching from one build to components (unsets `kluster.build`), a refusal shown beside its field, and the disabled state with `settings` off. It fails without the component.
- `settings.test.ts` covers the components helpers.
- `flai test` on the changed paths passes, `components.examples` included.

## Notes

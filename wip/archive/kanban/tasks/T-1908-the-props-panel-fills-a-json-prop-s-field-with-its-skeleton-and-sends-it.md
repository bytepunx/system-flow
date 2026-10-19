---
id: T-1908
type: task
nature: feature
title: The props panel fills a json prop's field with its skeleton and sends it
status: done
parent: S-0449
owner: alex
created: 2026-10-19T13:58:45Z
updated: 2026-10-19T14:14:46Z
transitions:
  - to: ready
    at: 2026-10-19T13:58:49Z
    by: agent-S-0449
  - to: in-progress
    at: 2026-10-19T13:58:50Z
    by: agent-S-0449
  - to: done
    at: 2026-10-19T14:14:46Z
    by: agent-S-0449
stream: S-0449
tags: []
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/explorer/controls.ts, flaiover/src/lib/explorer/controls.test.ts, flaiover/src/lib/components/explorer/PropsPanel.svelte.test.ts, flaiover/src/routes/design/components/components.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 956
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 60
      output: 24965
      cache_read: 3435272
      cache_write: 103044
      cost: 1.8365
---
# T-1908 The props panel fills a json prop's field with its skeleton and sends it

## Work

Criterion 3, as TH-0515 settled it, on the dashboard's side. The props panel fills a json prop's field with the catalogue's `skeleton` and sends it.

- `flaiover/src/lib/server/repo.ts`: `ComponentProp` gains `skeleton?: unknown`, with a doc comment: the empty value of the prop's type that flai gives array, object, and unknown props.
- `flaiover/src/lib/explorer/controls.ts`: `initial` gives the literal default where it fits, as today. Otherwise, for a json control, it gives the prop's `skeleton` when the catalogue carries one. So `toTexts` fills the field on a pick, on Reset, and for an example that leaves the prop out. `toValues` then sends it, and `missing` no longer flags a required json prop that has a skeleton. Update the doc comments.
- Tests: `flaiover/src/lib/explorer/controls.test.ts`. A json prop with a skeleton and no default gets the skeleton's text and value. A literal default wins over the skeleton. Reset gives the skeleton again. A prop with no skeleton keeps today's behaviour. A required json prop with a skeleton is not missing. Each test fails without the change. Check that `flaiover/src/lib/components/explorer/PropsPanel.svelte.test.ts` and `flaiover/src/routes/design/components/components.svelte.test.ts` still pass, and fix their expectations only where the new behaviour is meant to change them. Those two files are in this task's touches.
- Docs: `design/system/flaiover-dashboard.md` (§ The Components page, the controls model and the props panel bullets) and `docs/users/flaiover.md` (§ Components, the props), in the same change.

It runs with T-1907 (layer 3). They share no path. The contract between them: `skeleton` is a JSON value on the prop, present for array, object, and unknown props only.

## Done when

- A json prop with a skeleton shows it in its field on a pick and after Reset, and the component is sent it.
- `flai test` on the changed paths passes, and the new tests fail without the change.

## Notes

---
id: S-0444
type: story
nature: feature
title: Every component of flaiover renders on the Components page, with the setup its stylesheet needs and examples for the components that need data
status: ready
parent: E-0026
owner: alex
created: 2026-10-16T21:39:02Z
updated: 2026-10-17T02:34:36Z
transitions:
  - to: ready
    at: 2026-10-17T01:02:25Z
    by: alex
tags: [dashboard]
topics: [components, dashboard, design]
touches: [system-flow.yaml, design/components/flaiover, flai/internal/components/flaiover_test.go, design/system/component-library.md]
after: [S-0443]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 379
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 43
          output: 796
          cache_read: 5625625
          cache_write: 10275
          cost: 1.413
cost_of_delay:
  value: 82.57
  by: planner-E-0026
  at: 2026-10-16T21:50:08Z
forecast:
  duration: 1h
  delivery: 2026-10-17T09:26:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 5, behind S-0441 and S-0442."
  by: flai
  at: 2026-10-17T02:34:36Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:52:07Z
---
# S-0444 Every component of flaiover renders on the Components page, with the setup its stylesheet needs and examples for the components that need data

## Goal

This repository uses its own component library. Its one UI sub-project, flaiover, has 63 components under `flaiover/src/lib/components`, 56 of them with props. This story makes every one render on the Components page with no error, so that outcome 1 holds here and the library is proven on a real SvelteKit app, not only a fixture. flaiover's manifest entry names its setup, the stylesheet `src/routes/layout.css`, and leaves out only components that cannot render outside the app, each with its reason. Each component that needs data to render gets an examples file under `design/components/flaiover/`.

## Acceptance criteria

- [ ] flaiover's entry in `system-flow.yaml` has `components` with `setup: src/routes/layout.css`, written with `flai components config` (S-0437). It excludes only what cannot render in the harness, such as a component that reads the app's server data. `design/system/component-library.md` lists each exclusion with its reason.
- [ ] Each component the catalogue lists renders in the harness without an `error`, from its defaults or from at least one example, in light and dark. A manual pass on the Components page records any that render wrongly, with what is wrong, in the story's notes.
- [ ] Each component with a required prop, or whose defaults render nothing useful, has an examples file under `design/components/flaiover/` at the path S-0443's ADR gives, and `flai check --strict` reports no `components.examples`.
- [ ] An integration test, skipped without `node`, starts the dev server on this repository's flaiover and asks it for each listed component's module. Each one compiles, so a component that stops compiling for the harness fails the test.

## Tasks
- T-1890 flaiover's manifest entry names its components' setup stylesheet and excludes what cannot render in the harness
- T-1891 Each flaiover component that needs data to render gets an examples file under design/components/flaiover, in light and dark
- T-1892 An integration test asks the dev server on this repository's flaiover for each listed component's module and fails on one that does not compile
- T-1893 A manual pass on the Components page renders every flaiover component in light and dark, and the component library's design records the exclusions and what was found

## Notes

- Drafted by the planner for E-0026, last of its stories, after S-0443, whose examples it writes.
- The six components without `$props()` (EditorDialog, ProjectSwitcher, SettingsPanel, UnsignedBanner, HostFlai, HostFlaiBanner) take no props and may read the app's stores; they are the first candidates for an exclusion or a setup of their own.
- TH-0483 answer 4 put the examples under `design/components/`, so this story changes no file of flaiover's own.

### Planning

- **Touches.** The story declared none; every touch comes from its four tasks.
  - **From design:** S-0443's layout, `design/components/<sub-project>/`, for the examples.
  - **From layout:** `system-flow.yaml`, written through `flai components config`; a new `flai/internal/components/flaiover_test.go`; `design/system/component-library.md`.
  - **Co-change, not added:** `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` (22% each), since no page changes; `flaiover/src/routes/board/+page.svelte` and `items/[id]/+page.svelte` (15%), which no task changes.
  - **Folder touch kept:** `design/components/flaiover`, for T-1891's examples files. Which components need one, and so which files are added, is known only once T-1890 has opened each. No other open story touches the folder, so it holds nothing. The draft's `flaiover/src/lib/components` touch, which held eight stories of E-0021, E-0022, and E-0024, is gone with the answer.
- **Forecast.** flai's 16m (size 8, at 116 s a unit over only 3 done medium stories) raised to 1h. The one folder touch stands for some 20 to 40 examples files, and the manual pass opens about 60 components in two themes. Delivery is flai's start, 02:04Z, plus 1h.
- **Cost of delay.** 82.57 USD a week, flai's share of E-0026's 600 USD penalty by the raised forecast, kept.

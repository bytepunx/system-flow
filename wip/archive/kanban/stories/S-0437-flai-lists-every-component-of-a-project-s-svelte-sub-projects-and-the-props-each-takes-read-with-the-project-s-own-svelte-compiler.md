---
id: S-0437
type: story
nature: feature
title: flai lists every component of a project's Svelte sub-projects and the props each takes, read with the project's own Svelte compiler
status: done
parent: E-0026
owner: alex
created: 2026-10-16T21:32:09Z
updated: 2026-10-16T22:44:54Z
transitions:
  - to: ready
    at: 2026-10-16T21:52:20Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-16T21:53:16Z
    by: agent-S-0437
  - to: review
    at: 2026-10-16T22:43:00Z
    by: agent-S-0437
  - to: done
    at: 2026-10-16T22:44:54Z
    by: orchestrator
tags: [cli]
topics: [components, design, cli]
touches: [design/adrs, design/adrs/README.md, design/system/component-library.md, design/system/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, design/system/project-manifest.md, docs/operators/settings.md, flai/internal/components/components.go, flai/internal/components/components_test.go, flai/internal/components/props.go, flai/internal/components/props_test.go, flai/internal/components/props.mjs, flai/internal/components/testdata, flai/cmd/components.go, flai/cmd/components_test.go, flai/cmd/root.go, docs/users/flai-reference.md, design/system/flai-cli.md, docs/users/flai.md, flai/internal/check/check.go, flai/internal/check/components_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3004
  priced_by: harness
  turns:
    - day: 2026-10-16
      ceremony: 2
      hand_edits: 2
      work: 52
  models:
    - model: claude-opus-5-5
      input: 406
      output: 167572
      cache_read: 18958804
      cache_write: 678962
      cost: 11.0855
  strategic:
    - kind: orchestrator
      seconds: 3567
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 159
          output: 2721
          cache_read: 24873860
          cache_write: 58419
          cost: 6.2506
        - model: claude-sonnet-5-5
          input: 24
          output: 138
          cache_read: 582845
          cache_write: 70248
          cost: 0.5345
cost_of_delay:
  value: 59.17
  by: planner-E-0026
  at: 2026-10-16T21:49:59Z
forecast:
  duration: 43m
  delivery: 2026-10-17T19:57:00Z
  basis: "flai's 43m from 96 s a unit over 61 done large feature stories on claude-opus-5-5, times size 27; delivery is flai's play-out, after S-0436."
  by: system-flow
  at: 2026-10-16T21:42:56Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:51:57Z
---
# S-0437 flai lists every component of a project's Svelte sub-projects and the props each takes, read with the project's own Svelte compiler

## Goal

The component library starts with a catalogue: which components a project has and what each takes. flai finds the components of each sub-project whose `kind` is `sveltekit` or `svelte`. It finds them by globs the manifest can set, `src/lib/**/*.svelte` by default, routes left out. It reads each component's props from its `$props()` and their TypeScript type with a script flai embeds, run with the project's own Node and its own `svelte/compiler`, so that the props are read as the project's compiler reads them. `flai components list` prints the catalogue that the dev server, the Components page, and the controls of later stories are built on. An ADR and a design document record the epic's shape, on TH-0483's answers.

## Acceptance criteria

- [x] An ADR records the component library on TH-0483's answers. It says which sub-projects it covers (`kind` `sveltekit` or `svelte`, others listed as unsupported) and where their components are found. It says how props are read: by a script flai embeds, run with the project's Node and its own `svelte/compiler`, never by executing a component. `design/system/component-library.md` describes the whole epic's design, and `design/system/README.md` lists it.
- [x] A sub-project in the manifest may carry `components` with `include` and `exclude` globs and `setup`, a module of the sub-project the harness imports before any component, such as the app's stylesheet. Defaults: include `src/lib/**/*.svelte`, nothing excluded, no setup. The manifest's check reports a glob that does not parse, a `setup` that names no file, and `components` on a sub-project of an unsupported kind. `flai components config <sub-project> [--include] [--exclude] [--setup] [--unset]` writes the key, keeping the manifest's other keys, order, and comments, and refuses a value the check reports. `design/system/project-manifest.md` and `docs/operators/settings.md` describe the key.
- [x] `flai/internal/components` lists each component with its sub-project, path, and name. It runs the embedded script with `node` in the sub-project's folder and parses its answer: each prop's name, its kind (`string`, `number`, `boolean`, `union` with its literals, `array`, `object`, `function`, `snippet`, or `unknown`), its type as written, its default as written, and whether it is required. A component whose script does not parse is listed with its error. A missing `node` or a sub-project without `svelte` installed is an error that names what is missing and how to fix it.
- [x] `flai components list [--project <name>] [--json]` prints the catalogue. Run on this repository it lists flaiover's components under `flaiover/src/lib/components`, and the props of `Chart.svelte` (`option`, `theme`, `height` defaulting to 360) and of `Gallery.svelte` as their sources declare them. Tests on fixture components cover each kind, an unparsable component, and the missing-tool errors, the Node ones skipped when `node` is absent.
- [x] `design/system/flai-cli.md` and `docs/users/flai.md`, with `flai-reference.md` regenerated, describe both commands.

## Tasks
- T-1853 An ADR records the component library on TH-0483's answers, and design/system/component-library.md describes the epic's design
- T-1854 The manifest's sub-projects take components with include and exclude globs and a setup module, checked, written by flai, and documented
- T-1855 flai/internal/components finds a sub-project's components and reads their props with an embedded script run on the project's own svelte/compiler
- T-1856 flai components list prints each Svelte sub-project's components and their props, and flai components config writes a sub-project's components key
- T-1857 The CLI design and the user docs describe flai components list and flai components config

## Notes

- Drafted by the planner for E-0026. First of the component stories: the catalogue every later one reads.
- TH-0483 questions 2 and 3 shape the ADR. The story builds on the recommended answers until they are answered otherwise.
- 63 `.svelte` files are under `flaiover/src/lib/components`, all `lang="ts"`, 56 with `$props()`. `svelte` ^5.56.1 is a dependency of flaiover, and `svelte/compiler` ships in it; `svelte2tsx` is not a dependency, so the script uses the compiler's own TypeScript-aware parse.
- `flai manifest set` reaches only the orchestration and planning keys, and `CLAUDE.md` gives flai every manifest key but `name`, `description`, and `dashboard`. So the key needs a writer of its own, `flai components config`, which S-0444 uses for flaiover.

### Planning

- **Touches.** The story declared none; every touch comes from its five tasks, file by file.
  - **From design:** the ADR and its index; `design/system/component-library.md`, new, with its row in `design/system/README.md`; `project-manifest.md` and `docs/operators/settings.md` for the key, whose settings test fails on a key without a row.
  - **From layout:** the new package `flai/internal/components` (`components.go`, `props.go`, the embedded `props.mjs`, and tests) on the pattern of `flai/internal/designsystem`; `Project` in `flai/internal/manifest/manifest.go` and the writer beside `flai manifest set`'s in `settings.go`; `flai/cmd/components.go` and `root.go`.
  - **Co-change, added:** `flai-cli.md`, `flai.md`, and `flai-reference.md`, all three in `claims.shared`.
  - **Co-change, not added:** `docs/operators/index.md` (17%), since this story adds no host action; `flaiover-dashboard.md` and `docs/users/flaiover.md` (16%, 14%), since the dashboard is untouched; `template/CHANGELOG.md` (7%), since the template does not change.
  - **Folder touches kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it; it is in `claims.shared`, so it holds no story. `flai/internal/components/testdata`, because the fixture components are the agent's to name; no other story outside this epic touches the package.
- **Forecast.** Kept at flai's 43m: size 27, at 96 s a unit over 61 done large feature stories. The config command added after flai's first figure is one more file in sizes it already counts. Delivery is flai's play-out, after S-0436.
- **Cost of delay.** 59.17 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.

### Accepted by the orchestrator

- Verified: ce191f21f8a44d94dc7a16172b257be8facfad85
- At: 2026-10-16T22:44:54Z

Verdict: all five criteria met at ce191f2; flai verify passed at that commit against main 3881963, integration and smoke included; no convention breach.
- 1: design/adrs/0162-the-component-library-reaches-sub-projects-of-kind-sveltekit-or-svelte-finds.md, design/adrs/README.md, design/system/component-library.md, design/system/README.md
- 2: flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings_test.go, flai/internal/check/check.go, flai/internal/check/components_test.go, flai/cmd/components.go, design/system/project-manifest.md, docs/operators/settings.md
- 3: flai/internal/components/components.go, flai/internal/components/props.go, flai/internal/components/props.mjs, flai/internal/components/components_test.go, flai/internal/components/props_test.go
- 4: flai/cmd/components.go, flai/cmd/root.go, flai/cmd/components_test.go, flai/internal/components/props_test.go
- 5: design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md

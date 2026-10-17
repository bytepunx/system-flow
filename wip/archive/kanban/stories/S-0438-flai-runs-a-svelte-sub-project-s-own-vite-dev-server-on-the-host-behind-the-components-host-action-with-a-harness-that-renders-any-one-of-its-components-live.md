---
id: S-0438
type: story
nature: feature
title: flai runs a Svelte sub-project's own Vite dev server on the host, behind the components host action, with a harness that renders any one of its components live
status: done
parent: E-0026
owner: alex
created: 2026-10-16T21:33:15Z
updated: 2026-10-17T01:08:46Z
transitions:
  - to: ready
    at: 2026-10-16T23:08:05Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-16T23:08:15Z
    by: agent-S-0438
  - to: review
    at: 2026-10-17T01:07:21Z
    by: system-flow
  - to: done
    at: 2026-10-17T01:08:46Z
    by: orchestrator
tags: [cli, dashboard]
topics: [components, design, cli, dashboard]
touches: [design/adrs, design/adrs/README.md, design/system/component-library.md, flai/internal/components/harness/plugin.mjs, flai/internal/components/harness/harness.html, flai/internal/components/harness/harness.mjs, flai/internal/components/harness.go, flai/internal/components/harness_test.go, flai/internal/components/devserver.go, flai/internal/components/devserver_test.go, flai/cmd/components.go, flai/cmd/components_test.go, docs/users/flai-reference.md, flai/internal/components/integration_test.go, flai/internal/components/testdata, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/hostapi.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/components.go, flai/internal/serve/components_test.go, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, docs/users/flai.md, design/system/dashboard-host-channel.md, flai/internal/components/proc_unix.go, flai/internal/components/proc_windows.go, flai/internal/components/serve.go, flai/internal/components/serve_test.go, flai/cmd/serve.go, flai/internal/config/config.go, flaiover/src/lib/server/agent.ts, design/issues/I-0156-flai-guard-refuses-a-story-s-sub-agents-flai-components-list-and-serve-which-only-read-or-run-the-project-s-own-dev-server.md, design/issues/summary.md]
after: [S-0437]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 7173
  priced_by: harness
  turns:
    - day: 2026-10-16
      hand_edits: 2
      work: 44
    - day: 2026-10-17
      ceremony: 1
      hand_edits: 3
      work: 42
  models:
    - model: claude-opus-5-5
      input: 738
      output: 361673
      cache_read: 46588535
      cache_write: 2299255
      cost: 28.8199
  strategic:
    - kind: orchestrator
      seconds: 2131
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 131
          output: 2492
          cache_read: 23999330
          cache_write: 62084
          cost: 6.0323
        - model: claude-sonnet-5-5
          input: 12
          output: 49
          cache_read: 172308
          cache_write: 53025
          cost: 0.1844
cost_of_delay:
  value: 110.09
  by: planner-E-0026
  at: 2026-10-16T21:50:00Z
forecast:
  duration: 1h20m
  delivery: 2026-10-17T20:51:00Z
  basis: "Its own forecast of 1h20m; 43rd in the pull order with an in-progress limit of 5, behind S-0436, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432, S-0433, S-0434 and S-0435."
  by: flai
  at: 2026-10-16T22:44:54Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:51:58Z
---
# S-0438 flai runs a Svelte sub-project's own Vite dev server on the host, behind the components host action, with a harness that renders any one of its components live

## Goal

A component renders as the project builds it, and follows each change to its source or to the stylesheet without a rebuild or a redeploy. flai runs the sub-project's own Vite with the project's own config, plugins, and dependencies, so that aliases, Tailwind, and the design stylesheet behave as in the app. It adds a plugin flai embeds, which serves a harness page that renders one component by its path, inside an error boundary, after the manifest's `setup` module. Vite's hot module replacement then carries every saved change into the harness. The dev server listens on 127.0.0.1 on a port flai picks. The dashboard frames the harness cross-origin and talks to it by `postMessage`, under a protocol the ADR fixes. Starting it runs the project's code on the host, so it waits for a host action, `components`, which the operator enables. flai stops a dev server nobody has asked about for 30 minutes.

## Acceptance criteria

- [x] An ADR records the dev server on TH-0483's answers 1 and 5. The server binds a loopback port, and the dashboard frames it cross-origin, with the container still holding nothing of the project (ADR-0029 and ADR-0031 refined). It records the `components` host action and the idle stop. It records the protocol: the page sends `select`, `props`, `theme`, and `tokens`; the harness answers `ready`, `rendered`, `error` with the file and line, `event`, and `size`; each side checks the other's origin.
- [x] `flai components serve [--project <name>]` writes a Vite config under `.flai-cache/components/<project>/<sub-project>/` that merges the sub-project's own config and adds flai's embedded plugin. It starts the sub-project's own `vite` on 127.0.0.1 and a free port, and prints the harness URL once Vite answers. The harness at `/__flai/components/` renders the component `?c=<path>` or a `select` message names, after the `setup` module, and posts an `error` instead of a blank frame when it throws. A sub-project without `vite` installed is refused with the fix.
- [x] `flai serve enable components` turns on the host action. `components.start` and `components.stop` `{project?}` are gated by it and journalled. `components.status` and `components.list` are reads: running or not, the URL, the port, since when, the last error, and the catalogue. `flai serve` supervises one dev server per sub-project and stops it after `components.idle` (30m by default) without a `components.status`, and when `flai serve` stops.
- [x] An edit to a fixture component's source, and a rewrite of the stylesheet its setup imports, each reach the running dev server as a hot update, not a full reload: an integration test drives the dev server on a fixture sub-project and reads Vite's HMR payloads. It is skipped when `node` is absent. Run by hand on this repository, the harness renders `flaiover/src/lib/components/design/Gallery.svelte`; the story's notes say what was seen.
- [x] `docs/operators/index.md` has a section on the `components` host action, and `docs/operators/settings.md` rows for it and `components.idle`. `design/system/flai-cli.md`, `docs/users/flai.md` with `flai-reference.md` regenerated, `design/system/dashboard-host-channel.md`, and `design/system/component-library.md` describe the dev server, the methods, and the protocol.

## Tasks
- T-1858 An ADR records the dev server on a loopback port framed cross-origin, the components host action, the idle stop, and the page-to-harness protocol
- T-1859 flai embeds a Vite plugin and a harness page that render one component by path after the setup module, and writes the merged Vite config
- T-1860 flai/internal/components starts the sub-project's own vite on 127.0.0.1 and a free port, waits for it to answer, and stops it
- T-1861 flai components serve writes the harness and the merged config, starts the dev server, and prints the harness URL
- T-1862 An integration test runs the dev server on a fixture sub-project and shows a source edit and a stylesheet rewrite reach the harness as hot updates, not reloads
- T-1863 The components host action gates components.start and components.stop, and flai serve supervises each dev server and stops it when idle
- T-1864 The operator docs, the CLI design, the channel's design, and the component library's design describe the dev server, the host action, and its methods

## Notes

- Drafted by the planner for E-0026. It is the epic's riskiest story: merging the project's own Vite config with SvelteKit's plugin and serving a page outside the app's routes. The second layer of tasks (T-1859, T-1860) proves both on a fixture before the host methods build on them.
- No proxy, reverse proxy, or byte stream exists in `flai/internal/hostapi` or `flai/internal/serve` today. `flai/internal/host/child.go` supervises the MCP servers and is the pattern for the dev server.
- Host actions are declared in `flai/internal/hostapi/writes.go` (`Action*`, `Actions`, `shellOnly`) with commands in `flai/cmd/serve_actions.go`.
- The integration test reads Vite's HMR payloads rather than drive a browser: the payloads show the update is hot, and the frame's side is S-0439's manual check.

### Hand run (T-1862)

- On 2026-10-17, `flai components serve --project flaiover`, run from this branch's build in the main checkout, printed `http://127.0.0.1:41090/__flai/components/` and Vite's pid once the harness answered.
- `?c=src/lib/components/design/Gallery.svelte` rendered the whole Gallery in headless Chromium (Playwright from flaiover's own dependencies), with no page or console errors: the type specimens, buttons, inputs, the sign-in form and card, the alerts, the invoice table, and the navigation, with 10 buttons, 7 inputs, and 1 table.
- It rendered unstyled: flaiover's manifest names no `components.setup` yet, which is S-0444's, and Gallery's token classes come from the stylesheet the Design page's preview frame provides. Its caption was empty, since `?c=` mounts with no props.
- SIGINT to flai stopped Vite's process group, and flai logged the stop.

### Planning

- **Touches.** The story declared none; every touch comes from its seven tasks, file by file.
  - **From design:** the ADR and its index; `design/system/dashboard-host-channel.md`, which lists the host actions and methods; `docs/operators/index.md` and `settings.md`, which the documentation convention requires for a new setting and host action.
  - **From layout:** the harness, plugin, config writer, and runner in `flai/internal/components`; the host action in `flai/internal/hostapi/writes.go` with `components.go` for the reads, merged in `hostapi.go`; supervision in a new `flai/internal/serve/components.go` beside `checks.go`; the commands in `flai/cmd/serve_actions.go` and `components.go`.
  - **Co-change, added:** `flai-cli.md`, `flai.md`, and the reference, in `claims.shared`.
  - **Co-change, not added:** `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` (16%, 13%), and `flaiover/src/lib/server/agent.ts` (7%), which are S-0439's with the page; `design/system/project-manifest.md` (7%), since `components.idle` is a host setting, not a manifest key.
  - **Folder touches kept:** `design/adrs` for the new ADR, in `claims.shared`; `flai/internal/components/testdata` for the fixture sub-project, whose files the agent names.
- **Forecast.** flai's 55m (size 34, at 96 s a unit) raised to 1h20m. Nothing in the history merges a project's Vite config with SvelteKit's plugin or runs a live dev server in a test, and T-1859's note keeps a second way open if the first fails. Delivery is flai's start, 20:03Z, plus 1h20m.
- **Cost of delay.** 110.09 USD a week, flai's share of E-0026's 600 USD penalty by the raised forecast, kept: every later story waits on it.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.

### Accepted by the orchestrator

- Verified: cba39804bceaad9f3d408f41f0cb45db0f553732
- At: 2026-10-17T01:08:46Z

Verdict: all five criteria met at cba3980; flai verify passed at that commit against main 5de85b6; no convention breach.
- 1: design/adrs/0163-flai-runs-a-sub-project-s-own-vite-dev-server-on-127-0-0-1-behind-the.md, design/adrs/README.md
- 2: flai/cmd/components.go, flai/internal/components/harness.go, flai/internal/components/harness/plugin.mjs, flai/internal/components/harness/harness.html, flai/internal/components/harness/harness.mjs, flai/internal/components/devserver.go, flai/internal/components/serve.go, flai/internal/components/harness_test.go, flai/internal/components/devserver_test.go, flai/internal/components/serve_test.go
- 3: flai/internal/hostapi/writes.go, flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/serve/components.go, flai/internal/serve/components_test.go, flai/internal/config/config.go, flai/cmd/serve_actions.go, flai/cmd/serve.go, flaiover/src/lib/server/agent.ts
- 4: flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/vite.config.js, flai/internal/components/testdata/harness/src/lib/Hello.svelte
- 5: docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/dashboard-host-channel.md, design/system/component-library.md

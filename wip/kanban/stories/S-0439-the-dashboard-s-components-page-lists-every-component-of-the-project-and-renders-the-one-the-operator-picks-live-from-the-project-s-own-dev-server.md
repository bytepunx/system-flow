---
id: S-0439
type: story
nature: feature
title: The dashboard's Components page lists every component of the project and renders the one the operator picks, live, from the project's own dev server
status: backlog
parent: E-0026
owner: alex
created: 2026-10-16T21:34:32Z
updated: 2026-10-16T22:44:54Z
transitions: []
tags: [dashboard]
topics: [components, design, dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/components/+server.ts, flaiover/src/routes/api/components/components.test.ts, flaiover/src/routes/api/components/start/+server.ts, flaiover/src/routes/api/components/stop/+server.ts, flaiover/src/lib/explorer/protocol.ts, flaiover/src/lib/explorer/protocol.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/lib/components/explorer/ComponentList.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [S-0438]
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
      seconds: 373
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 41
          output: 767
          cache_read: 5187103
          cache_write: 8819
          cost: 1.3027
cost_of_delay:
  value: 50.92
  by: planner-E-0026
  at: 2026-10-16T21:50:02Z
forecast:
  duration: 37m
  delivery: 2026-10-17T21:31:00Z
  basis: "Its own forecast of 37m; 44th in the pull order with an in-progress limit of 5, behind S-0436, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432, S-0433, S-0434, S-0435 and S-0438."
  by: flai
  at: 2026-10-16T22:44:54Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:52:00Z
---
# S-0439 The dashboard's Components page lists every component of the project and renders the one the operator picks, live, from the project's own dev server

## Goal

The operator sees every component of their system from the dashboard (outcome 1) and sees a change to its source show at once (outcome 3), without an IDE. A Components page in the Design group lists the catalogue by sub-project and folder, with a filter. The operator picks a component, and the page frames the harness on the dev server `components.status` names, in light and dark and at chosen widths. The page starts the dev server when the operator asks and keeps it alive while it is open. It says how to enable the host action when it is off, and shows the harness's error, with the file and line, when a component fails to render.

## Acceptance criteria

- [ ] The Design group has a second page, Components, at `/design/components`, in `sitemenu.ts`, `SiteMenu.svelte`'s `HREF` map, and `sitemenu.test.ts`.
- [ ] The page lists each supported sub-project's components by folder from `components.list`, with a filter on name and path, and lists an unsupported sub-project with its kind. Picking one frames the harness at the URL from `components.status`, selects the component by a `select` message, and sizes the frame from `size`. The page offers light and dark, side by side or one at a time, and widths of 375, 768, and 1280 pixels or the panel's own.
- [ ] With the host action off, the page says to run `flai serve enable components`. With the server stopped, a Start button calls `components.start` and shows Vite's last lines when it fails. While open, the page asks `components.status` every minute, so the idle stop does not take it. A harness `error` shows the message, the file, and the line in place of the frame.
- [ ] The protocol client checks the frame's origin against the dev server's URL and ignores any other sender, shown by tests. A saved edit to the picked component's source shows in the frame without reloading the page: checked by hand on flaiover's `Gallery.svelte`, and recorded in the story's notes.
- [ ] `/api/components`, `/api/components/start`, and `/api/components/stop` call the four host methods, which are in `REQUIRED_METHODS`. A host flai without them answers 501, and the page says to upgrade it. Route tests cover each.
- [ ] `docs/users/flaiover.md` describes the page, and `design/system/flaiover-dashboard.md` its route, API, and protocol client.

## Tasks
- T-1865 The dashboard's /api/components routes call components.list, status, start, and stop on the host, answering 501 from an older flai
- T-1866 A protocol client sends select, props, theme, and tokens to the harness frame and takes its answers only from the dev server's origin
- T-1867 The Design group in the site menu gains the Components page at /design/components
- T-1868 The Components page lists the catalogue with a filter, frames the picked component in light and dark at chosen widths, starts the dev server, and shows its errors
- T-1869 The dashboard's user docs and design describe the Components page, its API, and its protocol client

## Notes

- Drafted by the planner for E-0026, after S-0438, whose methods it calls.
- The Design page's preview frame (`flaiover/src/lib/components/design/PreviewFrame.svelte`) is a blank same-origin frame for the gallery; this frame is cross-origin to a dev server, so it shares only the sizing idea.

### Planning

- **Touches.** The story declared none; every touch comes from its five tasks, file by file.
  - **From design:** `design-system-tool.md` § The site menu names the four places a page is added: `sitemenu.ts`, its test, `SiteMenu.svelte`'s `HREF` map, and the route.
  - **From layout:** `repo.ts` and `agent.ts` as `/api/design` uses them; the three routes under `flaiover/src/routes/api/components/`; the protocol client in a new `flaiover/src/lib/explorer/`; the page's components in a new `flaiover/src/lib/components/explorer/`, beside `components/design/`.
  - **Co-change, added:** `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`.
  - **Co-change, not added:** `flai-cli.md` and `flai.md` (37%, 34%), `docs/operators/index.md` (26%), and `flai/internal/hostapi/writes.go` (15%), all S-0438's; `flaiover/src/routes/items/[id]/+page.svelte` (9%), which this page does not change.
  - No folder touch.
- **Forecast.** Kept at flai's 37m: size 23, at 96 s a unit over 61 done large feature stories. Delivery is flai's play-out after S-0438's raised 1h20m.
- **Cost of delay.** 50.92 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.

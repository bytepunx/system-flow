---
id: S-0382
type: story
nature: feature
title: The dashboard's Design group shows a project's design system and previews its tokens live on sample components
status: in-progress
parent: E-0020
owner: alex
created: 2026-10-13T17:44:35Z
updated: 2026-10-15T01:33:11Z
transitions:
  - to: ready
    at: 2026-10-15T01:10:37Z
    by: alex
  - to: in-progress
    at: 2026-10-15T01:33:11Z
    by: agent-S-0382
tags: [dashboard]
topics: [design]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte.test.ts, flaiover/src/lib/design/preview.ts, flaiover/src/lib/design/preview.test.ts, flaiover/src/lib/components/design/PreviewFrame.svelte, flaiover/src/lib/components/design/PreviewFrame.svelte.test.ts, flaiover/src/lib/components/design/Gallery.svelte, flaiover/src/lib/components/design/TokenSheet.svelte, flaiover/src/lib/components/design/TokenSheet.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0381]
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
      seconds: 1173
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 854
          cache_read: 9702690
          cache_write: 11859
          cost: 2.4395
cost_of_delay:
  value: 192.31
  by: planner-E-0020
  at: 2026-10-13T17:51:37Z
forecast:
  duration: 50m
  delivery: 2026-10-15T02:31:00Z
  basis: "Its own forecast of 50m; 1st in the pull order with an in-progress limit of 5, behind S-0387."
  by: flai
  at: 2026-10-15T01:32:59Z
finalized:
  by: alex
  at: 2026-10-15T00:09:10Z
---
# S-0382 The dashboard's Design group shows a project's design system and previews its tokens live on sample components

## Goal

The operator opens a new top-level Design group in the dashboard's site menu and sees the project's design system as it stands: its tokens laid out, and a gallery of sample components styled by them in light and dark, so the effect of the design on the project's visual language is visible in the browser without building or deploying the project.

## Acceptance criteria

- [ ] The site menu has a fourth group, Design, beside Workflow, Status, and Host, whose first page, `/design`, keeps `?project=` as the other pages do.
- [ ] The page shows the project's design system: its colours as swatches with each one's contrast against the ground in both themes, its type scale as sample text, its spacing, radii, and shadows drawn to size, and its guidance rendered from the document.
- [ ] The page previews the tokens on a gallery of sample components, at least buttons, inputs, a form, a card, an alert, a table, navigation, and headings with body text, in light and dark, isolated so that the dashboard's own theme neither reaches the gallery nor changes with it.
- [ ] A project with no design system document is told so, with the command that writes one.
- [ ] The page follows a change to the document made elsewhere, as the other pages follow theirs.
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the Design group and its page.

## Tasks
- T-1538 The dashboard's /api/design endpoint answers design.get through the channel, cached and forgotten when the document changes
- T-1539 The site menu has a fourth group, Design, whose first page is /design
- T-1540 An isolated preview frame styles a gallery of sample components from a set of tokens, in light and dark, apart from the dashboard's theme
- T-1541 The /design page lays out the project's tokens and guidance beside the live preview, says when there is no design system, and follows changes
- T-1542 The dashboard's design and user docs describe the Design group, its page, and the live preview

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 recommends an isolated preview (an iframe or a shadow root) whose CSS custom properties come from the tokens over components compiled with the dashboard, so that a changed token restyles the gallery with no build.

### Planning

- **Touches.** None declared. All are the planner's, every one a file, from the layout the explorer read:
  - `flaiover/src/lib/server/repo.ts` and `agent.ts`: every host method the dashboard calls is in `REQUIRED_METHODS`, and reads are cached in `repo.ts`. With `flaiover/src/routes/api/design/+server.ts` and its test, new, in the pattern of `api/docs/adrs`. Criteria 2 and 5.
  - `flaiover/src/lib/sitemenu.ts`, `sitemenu.test.ts`, and `SiteMenu.svelte.test.ts`: the menu's groups and the tests that assert their labels. Criterion 1.
  - `flaiover/src/lib/design/preview.ts`, `flaiover/src/lib/components/design/PreviewFrame.svelte`, `Gallery.svelte`, and `TokenSheet.svelte`, with tests, and `flaiover/src/routes/design/+page.svelte` with its page test: new. Criteria 2 to 5.
  - `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`: design. Criterion 6.
- **Co-change.** `flai touches suggest` lists the flai docs and `flai/internal/hostapi/writes.go`, which change with the channel in general. None is added: `design.get` is S-0381's, and this story writes nothing.
- **No folder touch.**
- **Overlaps.** S-0378, in progress, claims `repo.ts` and `agent.ts`; this story waits behind S-0381, so it starts after S-0378 has likely been accepted.
- **Layers.** T-1538, T-1539, and T-1540 together; then T-1541; then T-1542.
- **Forecast: 50m, delivery 2026-10-13T22:35:00Z.** `flai forecast` gave 39m: 95 s per unit of size over 50 done feature stories on claude-opus-5-5 in the large band, times size 24. Raised about a quarter to 50m, as S-0378 was, because its tasks add new Svelte components with component tests and an isolated frame, which a rate per touch does not tell apart. The delivery adds the 11m to flai's 22:24Z.
- **Cost of delay: 192.31 USD a week.** `flai cod`: 50m of the 5h12m forecast over E-0020's seven open stories, of the epic's 1200 USD a week. It stands.

---
id: S-0441
type: story
nature: feature
title: A token change on the Design page reaches the components on the Components page, as the operator tries it and once it is saved
status: ready
parent: E-0026
owner: alex
created: 2026-10-16T21:36:09Z
updated: 2026-10-17T01:02:24Z
transitions:
  - to: ready
    at: 2026-10-17T01:02:19Z
    by: alex
tags: [dashboard, cli]
topics: [components, design, dashboard]
touches: [flaiover/src/lib/design/draft.ts, flaiover/src/lib/design/draft.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, flaiover/src/lib/explorer/tokens.ts, flaiover/src/lib/explorer/tokens.test.ts, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/lib/components/explorer/ComponentFrame.svelte, docs/users/flaiover.md, design/system/component-library.md, design/system/flaiover-dashboard.md]
after: [S-0436, S-0439]
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
      seconds: 617
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 50
          output: 941
          cache_read: 7235043
          cache_write: 15213
          cost: 1.8177
cost_of_delay:
  value: 42.66
  by: planner-E-0026
  at: 2026-10-16T21:50:04Z
forecast:
  duration: 31m
  delivery: 2026-10-17T01:37:00Z
  basis: "Its own forecast of 31m; 2nd in the pull order with an in-progress limit of 5, behind S-0438 and S-0440."
  by: flai
  at: 2026-10-17T01:02:24Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:52:03Z
---
# S-0441 A token change on the Design page reaches the components on the Components page, as the operator tries it and once it is saved

## Goal

The operator sees how a change to the design system affects the real components, live (outcome 3). A saved token change already reaches them: S-0436's save writes the stylesheet with the tokens, and the dev server's hot update carries it into the frame. A change the operator is still trying on the Design page reaches them too. The Design page keeps its unsaved values in the browser, per project, as ADR-0151 intended and nothing yet does. The Components page reads them and sends them to the harness as a `tokens` message. The harness lays them over the generated stylesheet's variables, under the same names and the same theme selector as `flai design css` writes. The operator sees which tokens the frame shows, saved or tried, and can drop the tried ones.

## Acceptance criteria

- [ ] The Design page keeps its unsaved token values in the browser per project, restores them on a reload, and clears them on save or discard; a test covers each.
- [ ] A function in flaiover turns token values into the declarations the generated stylesheet carries, under its variable names and its light and dark selectors (ADR-0159). It is tested against `flai/internal/designsystem/testdata`'s `themes.tokens.json` and its Tailwind golden, so that it cannot drift from `flai design css`.
- [ ] The Components page sends the tried values as a `tokens` message whenever they change, including from the Design page in another tab, by the `storage` event. The harness writes them into one `<style>` after the project's own stylesheet, so they win, and an empty message removes it. The page says "showing tried tokens" with a button that drops them.
- [ ] A token saved on the Design page shows on the Components page without a reload, by S-0436's stylesheet write and the dev server's hot update. It is checked by hand on this repository once its design system exists, or on a fixture project, and recorded in the story's notes.
- [ ] `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md`, and `design/system/component-library.md` describe tried and saved tokens on the Components page.

## Tasks
- T-1874 The Design page keeps its unsaved token values in the browser per project, restores them on a reload, and clears them on save or discard
- T-1875 flaiover turns token values into the generated stylesheet's declarations, tested against flai's own golden so it cannot drift
- T-1876 The harness lays a tokens message over the project's stylesheet in one style element, and an empty message removes it
- T-1877 The Components page sends the tried token values to the frame as they change, says it shows tried tokens, and drops them on request
- T-1878 Tried and saved tokens are checked by hand on the Components page, and the dashboard's docs and the component library's design describe them

## Notes

- Drafted by the planner for E-0026, after S-0436, whose save writes the stylesheet, and S-0439, whose page it extends.
- `flaiover/src/routes/design/+page.svelte` keeps its edits in page state (`dirty` from `changed(saved, edits)`) and asks before leaving; nothing persists them across a reload, though `design/system/design-system-tool.md` § Recommendation and ADR-0151 name it. This story builds it, since the Components page needs the values outside that page.

### Planning

- **Touches.** The story declared none; every touch comes from its five tasks, file by file.
  - **From design:** ADR-0151 and `design-system-tool.md` § Working on it together, for the Design page's kept values; ADR-0159 for the dark selector; `flai/internal/designsystem/testdata/themes.tokens.json` and its golden, read, not changed.
  - **From layout:** a new `flaiover/src/lib/design/draft.ts` beside `tokens.ts` and `preview.ts`; the Design page and its test; `flaiover/src/lib/explorer/tokens.ts` beside S-0439's client; the harness in flai; the Components page and its frame.
  - **Co-change, added:** `design/system/flaiover-dashboard.md` (88%), for the Design page's kept draft and the `tokens` message.
  - **Co-change, not added:** `flai-cli.md` and `flai.md` (46%), and `docs/operators/index.md` (33%), since neither a command nor a host action changes.
  - No folder touch.
- **Forecast.** Kept at flai's 31m: size 19, at 96 s a unit over 61 done large feature stories. Delivery is flai's play-out, after S-0436 and S-0439.
- **Cost of delay.** 42.66 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.

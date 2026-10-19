---
id: S-0449
type: story
nature: feature
title: Improvements to Design Components Screen (Part 1)
status: done
owner: alex
created: 2026-10-19T04:51:01Z
updated: 2026-10-19T17:30:21Z
transitions:
  - to: ready
    at: 2026-10-19T04:51:05Z
    by: alex
  - to: in-progress
    at: 2026-10-19T04:51:14Z
    by: agent-S-0449
  - to: review
    at: 2026-10-19T14:34:34Z
    by: agent-S-0449
  - to: done
    at: 2026-10-19T17:30:21Z
    by: alex
tags: [flaiover, flai]
touches: [design/system/component-library.md, docs/users/flai.md, flai/internal/components/harness/harness.html, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Dialog.svelte, flai/internal/components/testdata/harness/src/lib/Leaky.svelte, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, design/issues/I-0160-the-components-page-frames-a-component-that-draws-a-fixed-overlay-such-as-a-dialog-at-160px-as-its-content-s-height-is-0-cutting-the-dialog-off.md, design/issues/summary.md, design/adrs/0172-the-component-harness-answers-a-framed-component-s-requests-to-the-app-s-own.md, design/adrs/README.md, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md, design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/system/flai-cli.md, docs/users/flai-reference.md, flai/cmd/components.go, flai/internal/components/props.go, flai/internal/components/props.mjs, flai/internal/components/props_test.go, flai/internal/components/testdata/typed/package.json, flai/internal/components/testdata/typed/src/lib/Board.svelte, flai/internal/components/testdata/typed/src/lib/Imported.svelte, flai/internal/components/testdata/typed/src/lib/types.ts, flai/internal/components/testdata/typed/tsconfig.json, flaiover/src/lib/components/explorer/PropsPanel.svelte.test.ts, flaiover/src/lib/explorer/controls.test.ts, flaiover/src/lib/explorer/controls.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/design/components/components.svelte.test.ts, design/issues/I-0166-flai-components-list-and-the-components-examples-check-take-about-2-3s-for-flaiover-since-props-mjs-resolves-json-props-types-with-typescript-up-from-0-4s.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5081
  priced_by: harness
  turns:
    - day: 2026-10-19
      ceremony: 2
      test_runs: 1
      hand_edits: 4
      work: 87
  models:
    - model: claude-haiku-5-5
      input: 90
      output: 50733
      cache_read: 5142215
      cache_write: 183696
      cost: 0.3922
    - model: claude-opus-5-5
      input: 628
      output: 260650
      cache_read: 35866095
      cache_write: 1075840
      cost: 19.1738
  strategic:
    - kind: orchestrator
      seconds: 68560
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 258
          output: 5123
          cache_read: 23073735
          cache_write: 63713
          cost: 5.8014
        - model: claude-sonnet-5-5
          input: 14
          output: 77
          cache_read: 261356
          cache_write: 62123
          cost: 0.2647
cost_of_delay:
  inputs:
    penalty_per_week: 10
    by: alex
    at: 2026-10-19T04:51:01Z
---
# S-0449 Improvements to Design Components Screen (Part 1)

## Goal

Currently the Components page needs a few usability improvements for it to be truly useful. As it stands, the preview iFrame centers components, causing them to render outside the visible pane. The JSON required doesn't provide any kind of schema - an empty schema or default schema would be much more useful. Sometimes the incorrect component is rendering in the iframe panel. If there is a component render issue, there are times that changing the ligh/dark/both option or size option can trigger at least a partial re-render.

## Acceptance criteria
- [x] Only put a 10px margin around the rendered component
- [x] Automatically size the pane in which the iFrame renders so that the entire component is visible on screen OR is visible via scroll bars to see over-flow content.
- [x] Always pre-populated the JSON schema for json prop fields even if it's all default/empty values.
- [x] Make sure that the selected component is all that is rendering, never render a prior component
- [x] None of the component interactions should result in 500 errors - when 500 errors occur, file these as issues

## Tasks
- T-1904 The harness frames a component with a 10px margin, measures all it draws, and keeps nothing of a prior component
- T-1905 The Components page's frame takes only the selected component's answers and shows all the harness draws
- T-1906 A framed component's requests to the app's own server are answered by the harness, never with a 500
- T-1907 flai's catalogue gives each json prop an empty skeleton of its type, resolved with the sub-project's typescript
- T-1908 The props panel fills a json prop's field with its skeleton and sends it

## Notes

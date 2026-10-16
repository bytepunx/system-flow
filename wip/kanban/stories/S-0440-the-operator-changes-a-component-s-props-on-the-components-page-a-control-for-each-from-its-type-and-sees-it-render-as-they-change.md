---
id: S-0440
type: story
nature: feature
title: The operator changes a component's props on the Components page, a control for each from its type, and sees it render as they change
status: backlog
parent: E-0026
owner: alex
created: 2026-10-16T21:35:19Z
updated: 2026-10-16T22:44:54Z
transitions: []
tags: [dashboard, cli]
topics: [components, dashboard]
touches: [flaiover/src/lib/explorer/controls.ts, flaiover/src/lib/explorer/controls.test.ts, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flaiover/src/lib/components/explorer/PropsPanel.svelte, flaiover/src/lib/components/explorer/PropsPanel.svelte.test.ts, flaiover/src/lib/components/explorer/EventLog.svelte, flaiover/src/routes/design/components/+page.svelte, docs/users/flaiover.md, design/system/component-library.md, design/system/flaiover-dashboard.md]
after: [S-0439]
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
          cache_write: 8818
          cost: 1.3027
cost_of_delay:
  value: 35.78
  by: planner-E-0026
  at: 2026-10-16T21:50:03Z
forecast:
  duration: 26m
  delivery: 2026-10-17T22:00:00Z
  basis: "Its own forecast of 26m; 45th in the pull order with an in-progress limit of 5, behind S-0436, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432, S-0433, S-0434, S-0435, S-0438 and S-0439."
  by: flai
  at: 2026-10-16T22:44:54Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:52:01Z
---
# S-0440 The operator changes a component's props on the Components page, a control for each from its type, and sees it render as they change

## Goal

The operator manipulates a component and sees how it renders as its data changes (outcome 2). A props panel beside the frame gives each prop a control from the kind S-0437 read: text, number, checkbox, a select for a union of literals, and a JSON editor for an array, an object, or a type flai could not resolve. A snippet prop takes text, and a function prop is a logger whose calls the panel lists. Each change reaches the harness as a `props` message, and the component updates in place, keeping its own state, with no remount. The props are kept in the URL, so a link reproduces what the operator sees.

## Acceptance criteria

- [ ] The panel lists each prop of the picked component with its type as written and its default, and a control by kind: `string` text, `number` a number input, `boolean` a checkbox, `union` a select of its literals, `array`, `object`, and `unknown` a JSON editor, and `snippet` text rendered as a snippet. A required prop without a value is marked, and the frame shows the harness's error until it has one.
- [ ] Each change is sent as a `props` message, and the harness applies it to the mounted component through reactive props, without remounting it, so state the component holds, such as an open menu, survives a change of another prop. A reset puts every prop back to its default.
- [ ] A function prop is replaced in the harness by a logger, which posts an `event` with the prop's name and its arguments as JSON (or their type where they do not serialise). The panel lists the last 50 events with their time.
- [ ] Invalid JSON is shown in the editor and not sent. The picked component's props are kept in the page's URL, and a link opened in another tab renders the same.
- [ ] Tests cover the control for each kind, value parsing and the URL's round trip, the event log, and the harness's in-place update in S-0438's integration test. `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md`, and `design/system/component-library.md` describe the panel.

## Tasks
- T-1870 A controls model maps each prop kind to a control, parses and serialises its value, and keeps the props in the URL
- T-1871 The harness applies props messages to the mounted component in place, renders snippet text, and logs each function prop's calls as events
- T-1872 The props panel and the event log sit beside the frame on the Components page, sending each change as a props message
- T-1873 The dashboard's user docs and design and the component library's design describe the props panel and the event log

## Notes

- Drafted by the planner for E-0026, after S-0439, whose page it extends.
- The harness's change is in flai (`flai/internal/components/harness/harness.mjs`), so this story touches flai as well as flaiover.

### Planning

- **Touches.** The story declared none; every touch comes from its four tasks, file by file.
  - **From layout:** the controls model beside S-0439's protocol client in `flaiover/src/lib/explorer/`; the panel and the log beside S-0439's components in `flaiover/src/lib/components/explorer/`; the page; the harness and its integration test in `flai/internal/components`.
  - **Co-change, added:** `design/system/flaiover-dashboard.md`, which changed with 90% of the commits that touched these paths; T-1873 adds the panel to the page's entry.
  - **Co-change, not added:** `flai-cli.md` and `flai.md` (47%), since no command changes; `docs/operators/index.md` (34%), since no host action changes.
  - No folder touch.
- **Forecast.** Kept at flai's 26m: size 16, at 96 s a unit over 61 done large feature stories. The in-place update may need a browser to prove, which T-1871 leaves to a recorded manual check if so. Delivery is flai's play-out, after S-0439.
- **Cost of delay.** 35.78 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.

---
id: S-0476
type: story
nature: improvement
title: The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501
status: done
owner: alex
created: 2026-10-21T11:59:18Z
updated: 2026-10-22T22:37:19Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:29Z
    by: alex
  - to: in-progress
    at: 2026-10-22T21:19:41Z
    by: agent-S-0476
  - to: review
    at: 2026-10-22T22:23:31Z
    by: agent-S-0476
  - to: done
    at: 2026-10-22T22:37:19Z
    by: orchestrator
tags: [flai, flaiover]
topics: [components, dashboard, cli]
touches: [design/adrs/README.md, design/system/component-library.md, design/components/README.md, template/root/design/components/README.md, flai/internal/components/examples.go, flai/internal/components/examples_test.go, flai/internal/check/components_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flai/internal/hostapi/components_test.go, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, flaiover/src/routes/api/components/examples/+server.ts, flaiover/src/routes/api/components/examples/examples.test.ts, flaiover/src/lib/server/repo.ts, flaiover/src/lib/explorer/protocol.ts, flaiover/src/lib/explorer/protocol.test.ts, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, flaiover/src/lib/components/explorer/ExampleList.svelte, flaiover/src/lib/components/explorer/ExampleList.svelte.test.ts, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, design/components/flaiover/src/lib/components/AcceptConfirm.examples.json, design/components/flaiover/src/lib/components/ActivityView.examples.json, design/components/flaiover/src/lib/components/AnalyzerPanel.examples.json, design/components/flaiover/src/lib/components/CancelConfirm.examples.json, design/components/flaiover/src/lib/components/ImportPrompt.examples.json, design/components/flaiover/src/lib/components/InboxView.examples.json, design/components/flaiover/src/lib/components/Messages.examples.json, design/components/flaiover/src/lib/components/NewAdrForm.examples.json, design/components/flaiover/src/lib/components/NewItemForm.examples.json, design/components/flaiover/src/lib/components/OrchestratorPanel.examples.json, design/components/flaiover/src/lib/components/PlannerPanel.examples.json, design/components/flaiover/src/lib/components/PublishBanner.examples.json, design/components/flaiover/src/lib/components/Threads.examples.json, design/components/flaiover/src/lib/components/WorktreeUncommitted.examples.json, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md, design/issues/summary.md, design/adrs/0187-an-example-may-carry-the-answers-the-component-harness-gives-a-framed-component.md, template/CHANGELOG.md, flai/internal/components/testdata/harness/src/lib/Answered.svelte, docs/users/flai-reference.md, design/issues/I-0194-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3850
  priced_by: harness
  turns:
    - day: 2026-10-22
      ceremony: 2
      test_runs: 2
      hand_edits: 2
      work: 73
  models:
    - model: claude-opus-5-5
      input: 594
      output: 219524
      cache_read: 41835391
      cache_write: 1005342
      cost: 18.6972
  strategic:
    - kind: planner
      seconds: 686
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 3621
          output: 59995
          cache_read: 15845309
          cache_write: 694237
          cost: 22.7257
        - model: claude-haiku-5-5
          input: 254
          output: 24782
          cache_read: 13606896
          cache_write: 501352
          cost: 0.9075
    - kind: orchestrator
      seconds: 2485
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 177
          output: 3060
          cache_read: 47774314
          cache_write: 40651
          cost: 10.5581
        - model: claude-sonnet-5-5
          input: 16
          output: 85
          cache_read: 313690
          cache_write: 62572
          cost: 0.2604
forecast:
  duration: 2h30m
  delivery: 2026-10-23T00:09:00Z
  basis: "Its own forecast of 2h30m; 3rd in the pull order with an in-progress limit of 5, behind S-0466, S-0477, S-0468 and S-0467."
  by: flai
  at: 2026-10-22T21:15:43Z
finalized:
  by: alex
  at: 2026-10-22T20:52:34Z
---
# S-0476 The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501

## Goal

This story remediates [I-0164](../../../design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md), "The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0164 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0164 is closed with `flai issue close I-0164 --reason` saying what fixed it

## Tasks
- T-2035 An ADR refining ADR-0172 lets an example carry the answers the harness gives a component's requests to the app's own server, and the design describes the format
- T-2038 flai parses, validates, and lists an example's answers, and the components.examples check reports a malformed one
- T-2040 The harness answers a framed component's fetch and EventSource from the answers select and props carry, and 501 only what none matches
- T-2041 The dashboard's examples route and repo type carry an example's answers through list and save
- T-2045 The Components page sends the chosen example's answers to the harness, lets the operator edit them, and saves them with the example
- T-2048 The fourteen flaiover components I-0164 names carry examples with answers, each shown in its answered state in the harness
- T-2051 The user docs describe an example's answers, and I-0164 is closed with what fixed it

## Notes

I-0164's cost "0m" is not a duration longer than zero, so no time_lost_per_cycle was derived from it.

### Planning

Planned by planner-S-0476 on 2026-10-22; the plan's thread on this story names the tasks, their layers, and the assumptions.

**The solution proposed.** ADR-0172 answers every framed request to the app's own server with 501 and lists "answer with fixtures" as the alternative to follow; I-0164's instance asks for the same. An example in `<path>.examples.json` may carry `answers`, and the harness answers a `fetch` or an `EventSource` from them when one matches, 501 otherwise. T-2035 records it as an ADR refining ADR-0172 before anything is built.

**Touches.** The story declared none, so `flai touches suggest` had nothing to start from; every touch is from the design or the code layout, found by the explorer:

- design: `design/adrs/README.md` (the ADR's index), `design/system/component-library.md` (§ The file, § The harness, § The protocol, § What the pass found), `design/components/README.md` and its template copy, `docs/users/flai.md`, `docs/users/flaiover.md`, I-0164's file and `design/issues/summary.md` (the close).
- layout, flai: `flai/internal/components/examples.go` and its test (the parser rejects unknown keys), `flai/internal/check/components_test.go`, `flai/cmd/components.go` and its test (the help text lists an example's keys), `flai/internal/hostapi/components_test.go`, `flai/internal/components/harness/harness.mjs` (the 501 is hard-coded there), `flai/internal/components/integration_test.go` and `testdata/harness/src/lib/Fetcher.svelte` (its fetch subtests).
- layout, flaiover: `flaiover/src/routes/api/components/examples/+server.ts` and its test (`exampleOf` drops unknown keys on save), `flaiover/src/lib/server/repo.ts` (the example type), `flaiover/src/lib/explorer/protocol.ts`, `ComponentFrame.svelte`, `ExampleList.svelte`, `flaiover/src/routes/design/components/+page.svelte`, each with its test.
- layout, fixtures: the fourteen `design/components/flaiover/src/lib/components/*.examples.json` I-0164 names; Messages and Threads have none yet, so theirs are new files named here.

No folder touch is kept. The new ADR's path is not named, since its number is not known until `flai adr new` makes it; `flai task done` widens T-2035's touches with it. `docs/users/flai-reference.md` is not named: T-2038 regenerates it only if the help text it changes is generated there.

**Forecast.** `flai forecast` gave 1h6m (median 89 s per unit of size over 87 done improvement stories on claude-opus-5-5 in the large band, size 44 from 2 criteria and 42 touches). Raised to 2h30m: the size counts fourteen fixture files as cheaply as one-line edits, but each needs the component's routes read and a render in the harness, and the work spans five layers over Go, the embedded harness with a Playwright-run integration test, the Svelte page, and an ADR. The two nearest stories, S-0443 (the examples format, Go, page, docs) and S-0449 (the harness and a headless pass), took 1h09m and 1h25m of agent time against forecasts of 1h4m and none. Delivery is flai's start, 7th in the pull order behind S-0477, S-0485, S-0468, S-0467, S-0466, S-0469, S-0470, and S-0472, plus that duration.

**Cost of delay.** `flai cod` could work out no value: the story has no inputs and no epic. TH-0592 asks the operator for them, recommending `time_lost_per_cycle: 20m`; the value is written once they are set.

### Accepted by the orchestrator

- Verified: 240d519c3b76fb3b16dcf34ac78d4bc7f667a255
- Batch: 20261022-222455
- Enqueued verified: 240d519c3b76fb3b16dcf34ac78d4bc7f667a255
- At: 2026-10-22T22:37:19Z

Verdict: pass at 240d519c3b76fb3b16dcf34ac78d4bc7f667a255
- 1: design/adrs/0187-an-example-may-carry-the-answers-the-component-harness-gives-a-framed-component.md, design/adrs/README.md, design/system/component-library.md, design/components/README.md, template/root/design/components/README.md, template/CHANGELOG.md, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Answered.svelte, flai/internal/components/examples.go, flai/internal/components/examples_test.go, flai/internal/check/components_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flai/internal/hostapi/components_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/components/examples/+server.ts, flaiover/src/routes/api/components/examples/examples.test.ts, flaiover/src/lib/explorer/protocol.ts, flaiover/src/lib/explorer/protocol.test.ts, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, flaiover/src/lib/components/explorer/ExampleList.svelte, flaiover/src/lib/components/explorer/ExampleList.svelte.test.ts, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, design/components/flaiover/src/lib/components/AcceptConfirm.examples.json, design/components/flaiover/src/lib/components/ActivityView.examples.json, design/components/flaiover/src/lib/components/AnalyzerPanel.examples.json, design/components/flaiover/src/lib/components/CancelConfirm.examples.json, design/components/flaiover/src/lib/components/ImportPrompt.examples.json, design/components/flaiover/src/lib/components/InboxView.examples.json, design/components/flaiover/src/lib/components/Messages.examples.json, design/components/flaiover/src/lib/components/NewAdrForm.examples.json, design/components/flaiover/src/lib/components/NewItemForm.examples.json, design/components/flaiover/src/lib/components/OrchestratorPanel.examples.json, design/components/flaiover/src/lib/components/PlannerPanel.examples.json, design/components/flaiover/src/lib/components/PublishBanner.examples.json, design/components/flaiover/src/lib/components/Threads.examples.json, design/components/flaiover/src/lib/components/WorktreeUncommitted.examples.json, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md
- 2: design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md, design/issues/summary.md

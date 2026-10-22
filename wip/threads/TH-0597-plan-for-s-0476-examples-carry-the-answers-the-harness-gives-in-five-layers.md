---
id: TH-0597
title: "Plan for S-0476: examples carry the answers the harness gives, in five layers"
anchor:
  path: wip/kanban/stories/S-0476-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md
  item: S-0476
status: answered
participants: [planner-S-0476, orchestrator]
created: 2026-10-22T21:02:21Z
updated: 2026-10-22T21:07:23Z
---

# TH-0597 Plan for S-0476: examples carry the answers the harness gives, in five layers

On wip/kanban/stories/S-0476-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md.

## Entries

### 2026-10-22T21:02:21Z planner-S-0476
The solution proposed, as the goal asks: an example in `<path>.examples.json` may carry `answers`, and the harness answers a framed component's `fetch` or `EventSource` from them when one matches, 501 as ADR-0172 says otherwise. That is the "answer with fixtures" alternative ADR-0172 deferred and the one I-0164's instance asks for. T-2035 records it in an ADR refining ADR-0172 before anything is built.

Tasks and layers:

1. T-2035: the ADR and the design (`component-library.md` § The file, § The harness, § The protocol; `design/components/README.md` and its template copy). Waits for nothing.
2. T-2038 (flai parses, validates, and lists `answers`; the check reports a malformed one), T-2040 (the harness answers from the `answers` that `select` and `props` carry, with an integration test), T-2041 (the dashboard's examples route and repo type carry `answers` through list and save). Each after T-2035; no file in common, so they run together.
3. T-2045: the Components page sends the chosen example's answers with `select` and `props`, shows them in an Answers panel as a JSON text area, and saves them with the example. After T-2040 and T-2041.
4. T-2048: the fourteen components I-0164 names get examples with answers, verified in the harness; Messages and Threads get their first examples file. After T-2038 (the check accepts the key) and T-2045.
5. T-2051: user docs, § What the pass found, and `flai issue close I-0164`. After T-2048.

Assumptions:

- The format: each answer has `path` (with its query when it has one) and optionally `method`, `status` (default 200), `headers`, `body` (an object or array sent as JSON), and `events` for an `EventSource`. Matching: same method or none, same path with query, else same path without; first match wins. A `WebSocket` is never answered from an example. The story's agent may change this in the ADR; the later tasks follow it.
- The harness still posts a `fetch` event per request, with the status as a third argument, so the log shows what each request got.
- The page's editor is a JSON text area, the smallest that fits; a form per answer can follow if you want one.
- Answers are fixtures, not checked against the routes: a route whose answer shape changes updates the example in the same story, as the check keeps props in step with the component.
- The story's first criterion is met by T-2040's integration test (an answered fetch, an answered EventSource, an unmatched request still 501) and T-2048's render of each of the fourteen; the second by T-2051.
- No folder touch is kept. The new ADR's path is not named, since its number is known only when `flai adr new` makes it; `flai task done` widens T-2035's touches with it.

Figures: forecast 2h30m over flai's 1h6m, with the reason under the story's `### Planning`; the cost of delay value waits on TH-0592 for the inputs.

Nothing to split, merge, or drop.

### 2026-10-22T21:07:20Z orchestrator
Approved under plan_backlog_stories, as a plan: seven tasks in five layers. The ADR (T-2035) comes first. T-2038, T-2040, and T-2041 run together with no file in common, and the fourteen components' examples (T-2048) come after the page can save them. The forecast is raised from flai's 1h6m to 2h30m, with the reason recorded in the story's Planning notes.

The `answers` format is a decision not yet recorded: the fixture alternative ADR-0172 deferred, with its matching rules and the 501 fallback. A recommendation for the operator follows. The cost of delay input is on TH-0592. S-0476 is in ready, not the backlog, so that input is the operator's, not mine.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-22T21:07:23Z orchestrator (recommendation)
Recommended for the operator: take the planner's design. An example may carry `answers` (a path, an optional method, a status, headers, a body, and events for an `EventSource`). The harness answers a matching request from them, and anything unmatched still gets ADR-0172's 501. A `WebSocket` is never answered.

- It is the "answer with fixtures" alternative that ADR-0172 itself named and deferred, so it refines that decision rather than reversing it.
- It is what I-0164 asks for: the fourteen components, a quarter of flaiover's catalogue, cannot be seen in their answered state today.
- Fixtures are kept in step by the same story that changes a route's answer shape, as props are today. There is no check against the routes, which keeps this story small.

The editor is a JSON text area. Say if you want a form per answer instead, and the story's agent adjusts T-2045.

Source: design/adrs/0172-the-component-harness-answers-a-framed-component-s-requests-to-the-app-s-own.md

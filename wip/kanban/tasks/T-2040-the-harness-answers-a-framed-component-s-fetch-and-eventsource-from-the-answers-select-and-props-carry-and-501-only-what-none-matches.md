---
id: T-2040
type: task
nature: improvement
title: The harness answers a framed component's fetch and EventSource from the answers select and props carry, and 501 only what none matches
status: done
parent: S-0476
owner: alex
created: 2026-10-22T20:59:33Z
updated: 2026-10-22T21:38:58Z
transitions:
  - to: ready
    at: 2026-10-22T21:21:10Z
    by: agent-S-0476
  - to: in-progress
    at: 2026-10-22T21:24:24Z
    by: agent-S-0476
  - to: done
    at: 2026-10-22T21:38:58Z
    by: agent-S-0476
stream: S-0476
tags: [flai]
touches: [flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, flai/internal/components/testdata/harness/src/lib/Answered.svelte]
after: [T-2035]
usage:
  source: log
  seconds: 874
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 43
      output: 228
      cache_read: 1820404
      cache_write: 109280
      cost: 0.6539
---
# T-2040 The harness answers a framed component's fetch and EventSource from the answers select and props carry, and 501 only what none matches

## Work

In `flai/internal/components/harness/harness.mjs`:

- `select` and `props` take an optional `answers` array, as T-2035's ADR defines it. Keep the current answers in module state beside `wanted` and `current`; a `select` replaces them (none given means none), and a `props` message replaces them too, as a props message is whole. Validate loosely: an entry that is not an object with a string `path` is skipped.
- In the `window.fetch` replacement (`answerTheAppsRequests`), before the 501: find the first answer whose `method` is absent or equals the request's, and whose `path` equals the request's path with its query, or, failing that, its path without the query. Answer with a `Response` of the answer's `status` (default 200), its `headers`, and its `body`: an object or array as JSON with `content-type: application/json` unless `headers` sets one, a string as it is, `null` or absent as an empty body. Post the `event` as now, with the status as the third argument, `[method, path, status]`, for an answered and a 501 request alike.
- In the `EventSource` replacement: when an answer for the path (same matching, no method) has `events`, do not fail. Set `readyState` 1, fire `open` on the next tick, then dispatch each event in order as a `MessageEvent` named `message` unless `event` names it, with `data` as the string, or the JSON of an object or array, so both `onmessage` and `addEventListener` see them; `close()` sets `readyState` 2. Without such an answer, fail as now.
- Leave the `WebSocket` replacement unchanged.
- Update the header comment that describes the protocol.

Tests, in `flai/internal/components/integration_test.go`: extend `checkFetches` and `testdata/harness/src/lib/Fetcher.svelte` so the fixture fetches a path and opens an `EventSource`, the test sends `select` with answers for both, and checks the fetch's status and body are the answers', the `EventSource` reached `readyState` 1 and received the events, an unmatched path still gets 501, and the `event` messages carry the status. The browser part skips without Playwright Chromium as the test already does; keep the existing 501 subtests.

This task waits for T-2035 for the format; it shares no file with T-2038 and T-2041 and runs beside them.

## Done when

- [ ] A framed component's `fetch` of a path an answer names gets that answer's status, headers, and body, and a `fetch` no answer matches still gets 501 with ADR-0172's body
- [ ] An `EventSource` of a path whose answer has `events` opens and receives them in order; one without still fails and closes
- [ ] The harness posts `fetch` events with `[method, path, status]`, answered or not
- [ ] The integration test covers an answered fetch, an answered EventSource, and an unmatched request, and passes where Chromium is available

## Notes

Planned by planner-S-0476. The harness is embedded in the flai binary (`harness.go`), so the page's part in T-2045 needs the flai release this story makes.

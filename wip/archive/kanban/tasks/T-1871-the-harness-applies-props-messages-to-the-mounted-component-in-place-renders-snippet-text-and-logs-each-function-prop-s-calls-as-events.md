---
id: T-1871
type: task
nature: feature
title: The harness applies props messages to the mounted component in place, renders snippet text, and logs each function prop's calls as events
status: done
parent: S-0440
owner: alex
created: 2026-10-16T21:35:31Z
updated: 2026-10-17T02:04:40Z
transitions:
  - to: ready
    at: 2026-10-17T01:57:27Z
    by: system-flow
  - to: in-progress
    at: 2026-10-17T01:57:28Z
    by: system-flow
  - to: done
    at: 2026-10-17T02:04:40Z
    by: system-flow
stream: S-0440
tags: [cli]
touches: [flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go]
usage:
  source: log
  seconds: 432
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 25765
      cache_read: 2310965
      cache_write: 89164
      cost: 1.5097
---
# T-1871 The harness applies props messages to the mounted component in place, renders snippet text, and logs each function prop's calls as events

## Work

- `harness.mjs`: mount with a `$state` props object and assign each `props` message's values to it, so the component updates without a remount. Turn snippet text into a snippet with `createRawSnippet`. Replace each function prop with a logger that posts `event` with the name and its arguments as JSON, or each argument's type where it does not serialise.
- `integration_test.go`: extend S-0438's test with a fixture whose internal state must survive a props change. Use a node script run against the dev server's module graph if a browser is needed for this check, or record why a manual check stands in for it.

Waits for nothing in this story; it builds on S-0438's harness.

## Done when

The integration test passes under `flai test --all flai/internal/components`, or the story's notes record the manual check and why.

## Notes

None.

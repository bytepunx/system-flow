---
id: T-0720
type: task
nature: improvement
title: The host API no longer offers the dashboard pushing or auto-publish
status: ready
parent: S-0195
owner: arobson
created: 2026-10-02T23:31:09Z
updated: 2026-10-02T23:31:26Z
transitions:
  - to: ready
    at: 2026-10-02T23:31:26Z
    by: agent-S-0195
stream: S-0195
tags: []
touches: [flai/internal/hostapi, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/hostapi_reads_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go]
after: [T-0717]
---
# T-0720 The host API no longer offers the dashboard pushing or auto-publish

## Work

`flai push --pending` and the `auto-publish` host action stay for the operator's shell but leave the dashboard: drop the `push.run` and `push.pending` host methods, leave `auto-publish` out of the host actions the dashboard reads and changes, and refuse it from `settings.set` with what to run in a shell (`flai serve enable auto-publish`). The `push` action now gates Publish only; its description says so. `flai dashboard`'s note on retired settings stops pointing at pushing from the board. Waits for T-0717: it builds what ADR-0067 decides.

## Done when

- No host method the dashboard can call pushes accepted work or toggles auto-publish; Publish still runs under the `push` action
- Tests cover the methods' absence and the settings refusal

## Notes

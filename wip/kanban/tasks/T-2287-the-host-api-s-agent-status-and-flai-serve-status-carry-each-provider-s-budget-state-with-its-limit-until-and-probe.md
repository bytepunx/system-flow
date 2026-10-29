---
id: T-2287
type: task
nature: improvement
title: The host API's agent.status and flai serve status carry each provider's budget state, with its limit, until, and probe
status: in-progress
parent: S-0534
owner: alex
created: 2026-10-28T21:26:43Z
updated: 2026-10-29T05:56:37Z
transitions:
  - to: ready
    at: 2026-10-29T05:56:36Z
    by: agent-S-0534
  - to: in-progress
    at: 2026-10-29T05:56:37Z
    by: agent-S-0534
stream: S-0534
tags: [flai, hostapi, budget]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/contract_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
after: [T-2282]
---
# T-2287 The host API's agent.status and flai serve status carry each provider's budget state, with its limit, until, and probe

## Work

Add `budget` to the state `agent.status` answers (`hostapi.go`, built in `cmd/serve_actions.go`), in the shape the inbox task gives the inbox: one entry per provider whose state is set, with `provider`, `limit`, `since`, `until`, and `probing`. It is the host's state, so every project's dashboard gets the same entries. Keep the contract test in step (`contract_test.go`).

`flai serve status` prints a line per held provider, such as `budget: claude-code held (seven_day) until 2026-10-28T09:00:00Z`, and its `--json` carries the same field. This lets the operator see the state in a shell.

This task waits for T-2282's reader. It runs beside the run-end task and the inbox task, which touch none of these files.

## Done when

- [ ] `agent.status` carries `budget` while a provider is held and leaves it out once cleared, tested with a fixed state.
- [ ] `flai serve status` and its `--json` show each held provider with its limit and `until`.
- [ ] The contract test passes.
- [ ] `flai test` passes on the changed paths.

## Notes

Planned by planner-S-0534.

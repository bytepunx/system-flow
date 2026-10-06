---
id: T-0881
type: task
nature: feature
title: system-flow.yaml takes orchestration.permissions and orchestration.agent, and flai check reports a bad one
status: done
parent: S-0218
owner: alex
created: 2026-10-05T04:45:42Z
updated: 2026-10-05T07:12:47Z
transitions:
  - to: ready
    at: 2026-10-05T07:09:53Z
    by: agent-S-0218
  - to: in-progress
    at: 2026-10-05T07:09:54Z
    by: agent-S-0218
  - to: review
    at: 2026-10-05T07:12:47Z
    by: agent-S-0218
  - to: done
    at: 2026-10-05T07:12:47Z
    by: agent-S-0218
stream: S-0218
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/orchestration_test.go]
usage:
  source: log
  seconds: 173
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 26
      output: 10690
      cache_read: 1523744
      cache_write: 38722
      cost: 0.7353
---
# T-0881 system-flow.yaml takes orchestration.permissions and orchestration.agent, and flai check reports a bad one

## Work

Extend the `Orchestration` block S-0217's T-0809 added to `manifest.Manifest` (`flai/internal/manifest/manifest.go`, beside `Planning`):

- `orchestration.permissions`, each off when unset: `plan_backlog_epics`, `finalize_drafts`, `promote_to_ready`, `order_ready`, `accept_reviews`, and `publish` are booleans; `answer_threads` is `off`, `recommend`, or `autonomous`.
- `orchestration.agent`, an `Agent` as `planning.agent` is, and `Manifest.OrchestrationAgent()`, which merges it over the project's `agent` field by field, config key by config key, and role by role, as `PlanningAgent()` does.

Keep `orchestration.policy` and `orchestration.release` as T-0809 defined them; this task adds no policy value. Extend `Orchestration.Errors()` to name an unknown permission key, an `answer_threads` outside its three values, and an `agent` that is not valid, each saying how to fix it. `flai check` reports them as `manifest.orchestration`, beside T-0809's findings (`flai/internal/check/check.go`).

The manifest's documentation is T-0892's, so that one task writes `project-manifest.md`.

This task waits for nothing: S-0218 waits for S-0217, so T-0809's block is on the main branch. It runs with T-0882, whose paths it does not share.

## Done when

- the manifest parses every permission, `answer_threads` in each of its values, and `orchestration.agent`, and a test pins each default as off
- `OrchestrationAgent()` merges over the project's agent, and a test pins the merge as `PlanningAgent()`'s test does
- `Errors()` names each bad value in a test, and `flai check` reports one as `manifest.orchestration` on a fixture
- `go test ./internal/manifest/ ./internal/check/` passes

## Notes

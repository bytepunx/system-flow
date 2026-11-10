---
id: T-2935
type: task
nature: remediation
title: The orchestrator's shell flai accept and tier-running flai verify refuse inside a role's container, naming the MCP tools, with a test that reproduces I-0256
status: done
parent: S-0687
owner: alex
created: 2026-11-10T08:49:43Z
updated: 2026-11-10T09:25:19Z
transitions:
  - to: ready
    at: 2026-11-10T08:56:41Z
    by: agent-S-0687
  - to: in-progress
    at: 2026-11-10T08:56:44Z
    by: agent-S-0687
  - to: done
    at: 2026-11-10T09:25:19Z
    by: agent-S-0687
stream: S-0687
tags: [flai, accept, verify]
touches: [flai/cmd/accept.go, flai/cmd/verify.go, flai/cmd/accept_orchestrator_test.go, flai/cmd/verify_test.go]
after: [T-2932]
usage:
  source: log
  seconds: 1715
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 43
      output: 9673
      cache_read: 7365756
      cache_write: 162456
      cost: 2.5717
---
# T-2935 The orchestrator's shell flai accept and tier-running flai verify refuse inside a role's container, naming the MCP tools, with a test that reproduces I-0256

## Work

Waits for T-2932, whose ADR decides the refusal and its words. First write the test that reproduces I-0256 in `flai/cmd/accept_orchestrator_test.go`: with `FLAI_ROLE=orchestrate`, a `FLAI_CONFIG` inside `.flai-cache/agents/orchestrate` that holds no host actions, and a manifest under which the merge queue batches on the host (`accept_reviews` on, local acceptance), `flai accept --by orchestrator --verified <commit>` today lands the story by hand, since `hostActionOn` in `flai/cmd/verify.go` reads every action as off. Then make `flai accept`, its dry run included, and `flai verify` when it runs tiers (not `--last`, not `--job`) refuse in that case: run as the orchestrator with its configuration in a role's cache folder, they say the host's actions cannot be read there and name the MCP tools `accept` and `verify`. Leave the operator's and a story's agent's runs, and the orchestrator's on the host, as they are. Test each refusal in `accept_orchestrator_test.go` and `verify_test.go`, and that the host run still accepts.

## Done when

- [ ] The reproducing test fails before the change and passes after it
- [ ] Inside a role's cache configuration the orchestrator's `flai accept` and tier-running `flai verify` refuse, naming the MCP tools; `flai verify --last` and `--job` still run
- [ ] The tests pass with `flai test` on the changed paths

## Notes

---
id: T-2893
type: task
nature: remediation
title: The OpenRouter gateway check says the key is out of credit when OpenRouter answers 402, before the checks that follow from it
status: in-progress
parent: S-0685
owner: alex
created: 2026-11-10T08:44:21Z
updated: 2026-11-10T08:49:22Z
transitions:
  - to: ready
    at: 2026-11-10T08:49:22Z
    by: agent-S-0685
  - to: in-progress
    at: 2026-11-10T08:49:22Z
    by: agent-S-0685
stream: S-0685
tags: [gateway, tests]
touches: [flai/internal/serve/gateway_test.go]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 17
          cache_read: 114072
          cache_write: 3023
          cost: 0.0317
---
# T-2893 The OpenRouter gateway check says the key is out of credit when OpenRouter answers 402, before the checks that follow from it

## Work

I-0237's cause, close-outs and batches failing on the smoke tier's paid gateway check, is gone since S-0548: the gateway checks are the manual `gateway` tier, and `TestRepositoryManifestRunsTheGatewayChecksOnlyByName` in `flai/internal/verify/select_test.go` holds that no path, story, `--all`, or batch run selects it. What is left is the check itself, when it runs by name or on the host's schedule. On a 402 it fails with what follows from the refusal, no `flai-check.md` and no guard refusal, not with the refusal.

In `flai/internal/serve/gateway_test.go`:

- Add a function that finds a gateway's credit refusal in Claude Code's output and stream: the `API Error: 402` text I-0237 records, and a result with `is_error` true that carries it.
- In `checkGateway`, when it finds one, fail once with a message that names the gateway, quotes the refusal, and says that topping up the key or raising its limit is the operator's. Leave out the checks that follow from it.
- When `keyNow` reads a `limit_remaining` of 0 or below before the run, fail the same way and do not start `claude`.
- Add a test with no network that feeds the function the 402 line from I-0237's instance and a passing stream, in the style of `TestGatewayStreamIsJudged`. It reproduces the failure without a paid call.

Waits for no task: it is the change the other two describe and cite.

## Done when

- `flai test flai/internal/serve/gateway_test.go` passes, and the new test fails when the function does not find the 402.
- A 402 in the stream gives one failure that names the refusal, and none of the follow-on failures.

## Notes

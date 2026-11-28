---
id: T-3603
type: task
nature: improvement
title: The orchestrator's role file, in the baseline, the template, and this repository, publishes with release_publish only, and its definition is generated again
status: done
parent: S-0837
owner: alex
created: 2026-11-17T12:37:49Z
updated: 2026-11-28T00:22:46Z
transitions:
  - to: ready
    at: 2026-11-27T23:57:39Z
    by: agent-S-0837
  - to: in-progress
    at: 2026-11-27T23:57:39Z
    by: agent-S-0837
  - to: done
    at: 2026-11-28T00:22:46Z
    by: agent-S-0837
stream: S-0837
tags: [flai, template]
touches: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/harness/harness_test.go, flai/internal/roles/roles_test.go, flai/internal/guard/roles_test.go]
after: [T-3600]
usage:
  source: log
  seconds: 4
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 0
      output: 1
      cache_read: 47972
      cache_write: 131
      cost: 0.0185
---
# T-3603 The orchestrator's role file, in the baseline, the template, and this repository, publishes with release_publish only, and its definition is generated again

## Work

This task waits for T-3600, because it writes ADR-0316 §1 into the role's prompt as accepted.

- In all three copies of the orchestrator's role file, remove the permission entry that lists `flai formula run publish` (around line 78). Rewrite the publish paragraph (around line 109) and step 13 (around line 132): when `release_evaluate` is met, call `release_publish` with the met figure, or the reasoning under `judgement`, as the reason. Log it with `activity_log`. Raise a refusal on a thread on the batch's most recently accepted story. Say nothing of a formula or its `judgement` input.
- `flai/internal/roles/baseline_test.go` requires the baseline to be a byte copy of the template's file: edit the template's and copy it.
- Update the tests that pin the role's text: `TestTheOrchestratorsPromptSaysHowItPublishes` and `TestTheOrchestratorsDefinitionSaysWhatEachPermissionDoes` in `flai/internal/harness/harness_test.go`, and the publish cases in `flai/internal/roles/roles_test.go` and `flai/internal/guard/roles_test.go`.
- Generate the harness definition again with `flai role definition --write` and check that its publish lines name `release_publish` only. Its output under `.flai-cache/agents/` is not committed.

## Done when

- The three copies say to publish with `release_publish` only and list no formula command.
- The baseline equals the template's copy.
- The regenerated definition names no publish formula.
- `flai test` passes on the changed paths.

## Notes

- Drafted by planner-S-0837.
- `flai/internal/metrics/testdata/good/design/workflows/roles/orchestrator.md` is a fixture copy. Change it only if a test reads its publish text.

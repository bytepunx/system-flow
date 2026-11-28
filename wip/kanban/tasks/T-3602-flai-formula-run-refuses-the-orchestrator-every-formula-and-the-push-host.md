---
id: T-3602
type: task
nature: improvement
title: flai formula run refuses the orchestrator every formula, and the push host action still gates a served agent's pending publish
status: done
parent: S-0837
owner: alex
created: 2026-11-17T12:37:42Z
updated: 2026-11-28T00:22:26Z
transitions:
  - to: ready
    at: 2026-11-27T23:57:38Z
    by: agent-S-0837
  - to: in-progress
    at: 2026-11-27T23:57:38Z
    by: agent-S-0837
  - to: done
    at: 2026-11-28T00:22:26Z
    by: agent-S-0837
stream: S-0837
tags: [flai]
touches: [flai/cmd/formula.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/release.go, flai/cmd/release_push_action_test.go, docs/users/flai-reference.md, flai/cmd/formula_test.go]
after: [T-3600]
usage:
  source: log
  seconds: 5
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
# T-3602 flai formula run refuses the orchestrator every formula, and the push host action still gates a served agent's pending publish

## Work

This task waits for T-3600, because it builds ADR-0316 §5 as accepted.

- `orchestratorRunsFormula` in `flai/cmd/formula.go` lets the orchestrator (`FLAI_ROLE` orchestrate) run only `publish`. Make it refuse every formula, naming `release_publish` for a publish. Remove the `publishFormula` constant if nothing else uses it, and the help example that names the publish formula.
- Replace `testTheOrchestratorRunsThePublishFormulaOnlyWithPublish` in `flai/cmd/orchestrate_permissions_test.go` with a test that the orchestrator is refused `flai formula run` for every formula, `publish` included, with `publish` on or off.
- Keep `servedPublishHeld` in `flai/cmd/release.go` as it is: the `push` host action still gates `flai release --pending` for an agent flai serve started. Rewrite only its doc comment and any help text that names the publish formula. Keep the three tests in `flai/cmd/release_push_action_test.go`, and rewrite their comments that name the formula.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference` when the help changes.

## Done when

- The orchestrator is refused `flai formula run` for every formula, and a test shows it.
- The `push` gate's tests still pass unchanged in what they assert.
- The reference is current, and `flai test` passes on the changed paths.

## Notes

- Drafted by planner-S-0837.
- The guard's side, the role file's `commands:` entry for `flai formula run publish`, is T-3603's.

---
id: T-1780
type: task
nature: improvement
title: design/workflows/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags
status: done
parent: S-0423
owner: alex
created: 2026-10-16T05:48:17Z
updated: 2026-11-02T04:27:24Z
transitions:
  - to: ready
    at: 2026-11-02T04:23:10Z
    by: agent-S-0423
  - to: in-progress
    at: 2026-11-02T04:23:11Z
    by: agent-S-0423
  - to: done
    at: 2026-11-02T04:27:24Z
    by: agent-S-0423
stream: S-0423
tags: [template]
touches: [design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, scripts/fetch.sh, template/root/scripts/fetch.sh, scripts/README.md, template/root/scripts/README.md, flai/cmd/release.go, flai/cmd/release_evaluate_test.go, docs/users/flai-reference.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, docs/operators/settings.md]
after: [T-2013]
usage:
  source: log
  seconds: 253
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 62
      output: 390
      cache_read: 2505787
      cache_write: 77145
      cost: 0.909
---
# T-1780 design/workflows/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags

## Work

`flai release --evaluate` exits 0 whether or not the policy is met (`evaluateRelease` in `flai/cmd/release.go`). Add `--require-met`, which exits 2 with the evaluation's reason when the policy is not met, so that a formula step can stop on it. Test it, and regenerate `docs/users/flai-reference.md`, which T-2013's help text also changed.

A formula step runs only a `flai` command, a script under `scripts/`, or a role, so `git fetch` is a script. Add `scripts/fetch.sh` (POSIX `sh`, `set -eu`), which runs `git fetch --tags` on the branch's remote. Add it to the template (`template/root/scripts/fetch.sh`) and here, and list it in both `scripts/README.md` files.

Then write the `publish` formula at `design/workflows/formulas/publish.md`, nature `improvement`, with three steps, each after the one before and each `in: main` (TH-0583, questions 1 and 3):

1. `script: scripts/fetch.sh`.
2. `flai: release --evaluate --require-met`.
3. `flai: release --pending --json`.

The body, with the baseline marker, tells the instance's agent to record in the instance's Notes the versions and tags step 3 printed, or the evaluation that stopped step 2. Write it in the template (`template/root/design/workflows/formulas/publish.md`) and here alike.

Do not bump `template/template.yaml`'s version: `flai release` bumps it at publish (D-0456).

Waits for T-2013, which the formula's `in: main` steps need.

## Done when

`release_evaluate_test.go` passes with the new flag. `flai check --strict` reports nothing on the formula here. `flai formula show publish` lists the three steps, each in a layer of its own and each `in: main`.

## Notes

The orchestrator starts the formula only when `release_evaluate` is met. Step 2 therefore stops an instance only when the policy changed between the two, and that stop is a failure with a thread, as any step's is.

Replanned 2026-10-21 by planner-S-0423 on TH-0471: the formula lives under `design/workflows/formulas/` (answer 2), and its instance is worked by a story's agent in its worktree (answer 5). Rewritten 2026-11-02 by agent-S-0423:

- to TH-0583's answers 1 and 3: `in: main` on the steps, three steps, and the body asking for the record;
- `template/template.yaml` dropped from the touches (D-0456);
- `git fetch` became `scripts/fetch.sh`, since a step cannot run `git` itself (ADR-0238).

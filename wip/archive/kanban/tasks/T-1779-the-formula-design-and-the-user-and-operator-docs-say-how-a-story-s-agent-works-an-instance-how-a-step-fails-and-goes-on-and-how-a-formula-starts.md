---
id: T-1779
type: task
nature: feature
title: The formula design and the user and operator docs say how a story's agent works an instance, how a step fails and goes on, and how a formula starts
status: done
parent: S-0422
owner: alex
created: 2026-10-16T05:48:10Z
updated: 2026-11-02T03:41:12Z
transitions:
  - to: ready
    at: 2026-11-02T03:34:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T03:34:44Z
    by: system-flow
  - to: done
    at: 2026-11-02T03:41:12Z
    by: agent-S-0422
stream: S-0422
tags: [cli]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, template/CHANGELOG.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, design/adrs/README.md, flai/cmd/role.go, flai/cmd/accept.go, design/adrs/0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md]
after: [T-1775, T-1776, T-1777, T-1778, T-2012]
usage:
  source: log
  seconds: 388
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 89
      output: 25372
      cache_read: 6445536
      cache_write: 135758
      cost: 2.5876
---
# T-1779 The formula design and the user and operator docs say how a story's agent works an instance, how a step fails and goes on, and how a formula starts

## Work

- In the formula section of `design/system/strategic-agents.md` that S-0421 writes, or the design document it made (then widen the touches to it), describe:
  - the instance's prompt and the convention's rule: a story's agent works the steps on its branch, closing each task as its step passes;
  - a role step through `role_run`, its trigger, how the agent is woken, and its charge to the instance (the charge waits on TH-0800: say what holds today);
  - a failed step: the blocked task, the thread, and the answer that starts the agent again on the step;
  - an instance's acceptance, with the orchestrator's evidence;
  - the starters: `flai formula run --ready`, the `formula` host action and a formula's `schedule`, and a role whose file permits `formula run <name>`.
- Describe `--ready` and the schedule in `docs/users/flai.md`.
- Add how to turn the action on, set a schedule, and read a blocked instance to `docs/operators/index.md`, and the `formula` action's row to `docs/operators/settings.md`.
- The formulas `README.md`, here and in the template, lists `schedule` among a formula file's front-matter fields.
- An ADR, refining ADR-0238, records the decisions S-0422 made: a formula's `schedule` and the `formula` host action, the served-agent exception for an instance's role step, and an instance's evidence.
- The help of `flai role run` (`cmd/role.go`) and of `flai accept --evidence` (`cmd/accept.go`) name the instance's exception and evidence; regenerate `docs/users/flai-reference.md` and `docs/operators/settings.md` with `make flai-reference`.
- Add the convention and template changes to `template/CHANGELOG.md` under `## Unreleased`.

Last layer.

## Done when

The markdown lint, the settings-index and reference tests, and `flai check --strict` pass.

## Notes

Rewritten by the planner on 2026-10-21, on TH-0471's answer 5 and TH-0582's approval. Widened by S-0422's agent on 2026-11-02 with the formulas README, the ADR, and the two help texts the earlier tasks left.

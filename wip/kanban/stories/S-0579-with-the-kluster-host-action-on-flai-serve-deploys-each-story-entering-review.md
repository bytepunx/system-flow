---
id: S-0579
type: story
nature: feature
title: With the kluster host action on, flai serve deploys each story entering review into its namespace, runs the review checks against it, tears it down when the story leaves review, and promotes main on acceptance
status: backlog
parent: E-0032
owner: alex
created: 2026-11-06T21:40:13Z
updated: 2026-11-07T00:14:16Z
transitions: []
tags: [kluster, flai, serve]
topics: [kluster, deployment, testing, review, host-actions]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/serve.go, flai/cmd/serve_actions.go, flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, docs/operators/index.md, design/system/kluster.md, docs/users/flaiover.md]
after: [S-0578]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 35
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 62
          cache_read: 381639
          cache_write: 53273
          cost: 0.104
cost_of_delay:
  value: 163.4
  by: planner-E-0032
  at: 2026-11-06T21:50:14Z
forecast:
  duration: 1h
  delivery: 2026-11-07T08:38:00Z
  basis: "Its own forecast of 1h; 27th in the pull order with an in-progress limit of 5, behind S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0555, S-0556, S-0557, S-0558, S-0559, S-0560, S-0561, S-0562, S-0564, S-0565, S-0567, S-0568, S-0569, S-0570, S-0571, S-0577 and S-0578."
  by: flai
  at: 2026-11-07T00:14:16Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:36Z
---
# S-0579 With the kluster host action on, flai serve deploys each story entering review into its namespace, runs the review checks against it, tears it down when the story leaves review, and promotes main on acceptance

## Goal

When the `kluster` host action is on, a story that reaches review is tested against a kluster with nobody running anything by hand. flai serve deploys the story into its namespace and runs the review checks with that namespace in their environment, so a project's smoke tier can target it. It removes the namespace when the story leaves review, and promotes main to the default namespace when the story is accepted. This is the toggle E-0032's outcome asks for: off by default, and turned on from the Features page.

## Acceptance criteria

- [ ] With the `kluster` host action on for a project that has a `kluster:` block, flai serve runs `flai kluster deploy` once for a story that enters review, and again when its branch's head changes while it is in review. flai serve's journal logs each deploy and its outcome.
- [ ] The review checks of a deployed story run after its deploy ends, with the namespace and the kube context S-0576's ADR names in their environment. With the action off, or no deployment, the checks run as before.
- [ ] A story that leaves review, back to in progress or cancelled, is torn down. A story that is accepted is torn down, and main is promoted to the default namespace.
- [ ] One deploy runs at a time per project, and others wait their turn. A deploy that fails opens one thread on the story naming the step and the log, and moves nothing.
- [ ] With the action off, or no `kluster:` block, flai serve runs no kluster command.
- [ ] The tests fake the kluster commands, and cover entering review, a new head, leaving review, acceptance, a failure, and the action off.
- [ ] `docs/operators/index.md` gains a section on the kluster host action, what it runs and when, and how to turn it on. `design/system/kluster.md` and the Features section of `docs/users/flaiover.md` say the same.

## Tasks

- T-2481 flai serve's kluster pass deploys a story entering review and again at each new head, one deploy at a time per project, and opens one thread on a failed deploy, only while the kluster host action is on
- T-2482 flai serve's kluster pass tears down a story that leaves review or is cancelled, and on acceptance tears it down and promotes main to the default namespace
- T-2483 The review checks of a deployed story run after its deploy ends, with the namespace and the kube context in their environment
- T-2484 The operators guide, the kluster design, and the Features section say what flai serve runs while the kluster host action is on, and when

## Notes

### Planning

- Touches:
  - From the layout: `flai/internal/serve/kluster.go` and its test, a new pass beside `forge.go`'s, and `flai/internal/serve/serve.go`, which runs the passes.
  - `flai/internal/serve/checks.go` and its test, where the review checks run.
  - `flai/cmd/serve_actions.go`, from co-change (5%) and the layout: flai serve builds its supervisors there, `serve.NewComponents` among them.
  - From the design: `docs/operators/index.md`, `design/system/kluster.md`, and `docs/users/flaiover.md` § Settings › Features.
  - `flai touches suggest` listed `design/system/flaiover-dashboard.md` (53%) and `design/system/flai-cli.md` (48%). Both are left out: the Features page lists the action with no change, and the story adds no command.
- `after` S-0578: this story reads S-0578's deploy, teardown, promote, and list in `flai/internal/kluster/deploy.go`, the deployment record under `.flai-cache/kluster/`, and its lock. Through S-0578 it also reads S-0577's `kluster` host action.
- Tasks:
  - Layer 1: T-2481 (deploys).
  - Layer 2: T-2482 (teardown and promotion, in the same files) and T-2483 (the checks, in `checks.go`). Both wait for T-2481 and share no path.
  - Layer 3: T-2484 (documents).
- Forecast 1h, raised from `flai forecast`'s 29m (107 s per unit, size 16). flai serve's passes run in the background, and their tests run long and have failed on timing under load (S-0560, S-0561, S-0569), which costs reruns. The delivery moves by the same 31m.
- Cost of delay 163.40 USD/week: `flai cod`'s share of E-0032's 1250 USD/week revenue input (alex, TH-0829), by forecast duration. It stands as flai gives it.

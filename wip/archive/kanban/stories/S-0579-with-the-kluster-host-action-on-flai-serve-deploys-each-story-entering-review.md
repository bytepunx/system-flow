---
id: S-0579
type: story
nature: feature
title: With the kluster host action on, flai serve deploys each story entering review into its namespace, runs the review checks against it, tears it down when the story leaves review, and promotes main on acceptance
status: done
parent: E-0032
owner: alex
created: 2026-11-06T21:40:13Z
updated: 2026-11-07T03:28:24Z
transitions:
  - to: ready
    at: 2026-11-07T00:52:00Z
    by: alex
  - to: in-progress
    at: 2026-11-07T02:36:31Z
    by: agent-S-0579
  - to: review
    at: 2026-11-07T03:26:49Z
    by: agent-S-0579
  - to: done
    at: 2026-11-07T03:28:24Z
    by: orchestrator
tags: [kluster, flai, serve]
topics: [kluster, deployment, testing, review, host-actions]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/serve.go, flai/cmd/serve_actions.go, flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, docs/operators/index.md, design/system/kluster.md, docs/users/flaiover.md, flai/cmd/checks.go, flai/cmd/serve.go, design/issues/I-0268-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md, design/issues/summary.md]
after: [S-0578]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3033
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      hand_edits: 5
      work: 67
  models:
    - model: claude-opus-5-5
      input: 298
      output: 126507
      cache_read: 46186973
      cache_write: 556721
      cost: 15.6821
  strategic:
    - kind: orchestrator
      seconds: 72
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 127
          cache_read: 846193
          cache_write: 92236
          cost: 0.2257
cost_of_delay:
  value: 163.4
  by: planner-E-0032
  at: 2026-11-06T21:50:14Z
forecast:
  duration: 1h
  delivery: 2026-11-07T03:47:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 5, behind S-0581."
  by: flai
  at: 2026-11-07T02:35:43Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:53:36Z
---
# S-0579 With the kluster host action on, flai serve deploys each story entering review into its namespace, runs the review checks against it, tears it down when the story leaves review, and promotes main on acceptance

## Goal

When the `kluster` host action is on, a story that reaches review is tested against a kluster with nobody running anything by hand. flai serve deploys the story into its namespace and runs the review checks with that namespace in their environment, so a project's smoke tier can target it. It removes the namespace when the story leaves review, and promotes main to the default namespace when the story is accepted. This is the toggle E-0032's outcome asks for: off by default, and turned on from the Features page.

## Acceptance criteria

- [x] With the `kluster` host action on for a project that has a `kluster:` block, flai serve runs `flai kluster deploy` once for a story that enters review, and again when its branch's head changes while it is in review. flai serve's journal logs each deploy and its outcome.
- [x] The review checks of a deployed story run after its deploy ends, with the namespace and the kube context S-0576's ADR names in their environment. With the action off, or no deployment, the checks run as before.
- [x] A story that leaves review, back to in progress or cancelled, is torn down. A story that is accepted is torn down, and main is promoted to the default namespace.
- [x] One deploy runs at a time per project, and others wait their turn. A deploy that fails opens one thread on the story naming the step and the log, and moves nothing.
- [x] With the action off, or no `kluster:` block, flai serve runs no kluster command.
- [x] The tests fake the kluster commands, and cover entering review, a new head, leaving review, acceptance, a failure, and the action off.
- [x] `docs/operators/index.md` gains a section on the kluster host action, what it runs and when, and how to turn it on. `design/system/kluster.md` and the Features section of `docs/users/flaiover.md` say the same.

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

### Accepted by the orchestrator

- Verified: 096c6cc19599aeb71035e12224800db7d0f79c7d
- At: 2026-11-07T03:28:24Z

Verdict: pass. All seven criteria are met in the diff, the close-out passed at the head commit 096c6cc, and nothing is uncommitted.
- 1: flai/internal/serve/kluster.go, flai/internal/serve/serve.go
- 2: flai/internal/serve/checks.go, flai/internal/serve/kluster.go, flai/cmd/checks.go, flai/cmd/serve.go, flai/internal/serve/serve.go
- 3: flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go
- 4: flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go
- 5: flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go
- 6: flai/internal/serve/kluster_test.go, flai/internal/serve/checks_test.go
- 7: docs/operators/index.md, design/system/kluster.md, docs/users/flaiover.md

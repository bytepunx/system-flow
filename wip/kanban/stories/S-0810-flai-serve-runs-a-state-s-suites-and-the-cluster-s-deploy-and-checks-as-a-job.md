---
id: S-0810
type: story
nature: feature
title: flai serve runs a state's suites and the cluster's deploy and checks as a job when a story enters it, skips the cluster's steps where none is configured, and sends the story back when a step fails
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:06Z
updated: 2026-11-17T08:22:24Z
transitions: []
tags: [flai]
topics: [verify, workflow, kluster]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/internal/workitem/transition.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/cmd/move.go, flai/cmd/move_workflow_test.go, flai/internal/workflow/runs.go, flai/internal/workflow/runs_test.go, flai/internal/workflow/review.go, flai/internal/workflow/review_test.go, flai/internal/serve/statejob.go, flai/internal/serve/statejob_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/serve.go, flai/internal/serve/agents_restart.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/sendback/sendback.go, flai/internal/sendback/sendback_test.go, flai/cmd/verify.go, flai/internal/verify/story.go, docs/users/flai-reference.md, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/kluster_requests.go, flai/internal/serve/kluster_requests_test.go, flai/internal/serve/checks_test.go, flai/internal/serve/flow.go, flai/internal/serve/flow_test.go, design/system/kluster.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/change-the-board.md]
after: [S-0806, S-0804]
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
    - kind: planner
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 1022
          cache_read: 239031
          cache_write: 8404
          cost: 0.1355
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 66
          cache_read: 318646
          cache_write: 52547
          cost: 0.1138
cost_of_delay:
  value: 175.37
  by: planner-S-0810
  at: 2026-11-16T23:52:16Z
forecast:
  duration: 2h30m
  delivery: 2026-11-17T14:25:00Z
  basis: "Its own forecast of 2h30m; 7th in the pull order with an in-progress limit of 3, behind S-0752, S-0758, S-0806, S-0754, S-0756, S-0808, S-0807, S-0241 and S-0809."
  by: flai
  at: 2026-11-17T08:22:24Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:56:28Z
---
# S-0810 flai serve runs a state's suites and the cluster's deploy and checks as a job when a story enters it, skips the cluster's steps where none is configured, and sends the story back when a step fails

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §4 and §5. The job is the one S-0804 makes for Rebase and resolve; this story starts it from a state's `runs:`. Today the cluster's deploy and checks run on entering any state of kind `waiting`, in code.

## Acceptance criteria

- [ ] When a story enters a state with `runs:`, flai serve rebases its branch onto the main branch and runs the state's suites and `kluster:` steps as a job, with its output kept and its steps recorded as they start and end.
- [ ] A job that passes leaves the story in the state; a rebase that stops on conflicts or a step that fails sends it back by the state's `send-back` transition, with the thread its agent is started on.
- [ ] A `kluster:` step is skipped, and the record says so and why, when the project configures no deploy, the kluster host action is off, or no agent namespace is available, and it never holds the story.
- [ ] The cluster's deploy and checks run where a state's `runs:` names them and on no other state, and the `kluster-verified` rule is evaluated by the move where a transition names it.
- [ ] On the baseline the cluster's deploy and checks run on entering `review` as before this story.

## Tasks

- T-3475 The move evaluates the kluster-verified rule where a transition names it, in place of a warning on every move into a waiting state
- T-3479 The workflow package says what a state runs on entry and which transition sends a story back from it
- T-3482 flai serve starts a state job when a story enters a state with runs: it rebases the branch, runs the state's suites, records each step, and sends the story back by the state's send-back transition when one fails
- T-3484 The state job runs the cluster's deploy and checks a state's runs names, skips them with the reason where no cluster can run them, and the cluster no longer deploys on review in code
- T-3487 The design and the guides say a state's runs starts a job with the cluster's steps, which are skipped where no cluster can run them, and kluster-verified is checked where a transition names it

## Notes

- After S-0806: the job is started from a state's runs key, which S-0806 parses.
- After S-0804: the job is the one S-0804 makes for Rebase and resolve, with its record and its send-back.

### Planning

Written by planner-S-0810 on 2026-11-16.

After:

- S-0806, kept as declared: this story reads the `runs:` key S-0806 parses onto a state and the baseline's `runs: {kluster: [deploy, checks]}` on `review`.
- S-0804, kept as declared: this story reads S-0804's resolve job (`flai/internal/serve/resolvejob.go`), its `closeout` record and progress lines, and its `flai/internal/sendback` package.
- S-0805 is not named: its suite of each tier, which T-3482 reads to run a state's suites, reaches main before S-0806, which names it in `after`.

Touches, all files, no folder. The story declared none. `flai touches suggest` needs a path to start from, so the touches come from the criteria, ADR-0312 §4, §5, §8, and an explorer's read of the code on main:

- Layout, the move and its rule: `flai/internal/workitem/rules.go`, where `kluster-verified` returns nil today, `rules_test.go`, `transition.go`, `flai/internal/mcpserver/server.go`, whose `ReviewWarnings` warns on every move into a waiting state, `server_test.go`, `flai/cmd/move.go`, `move_workflow_test.go`.
- Layout, the workflow: `flai/internal/workflow/review.go` and `review_test.go`, where `SendBack` lives, and a new `runs.go` with `runs_test.go`.
- Design and layout, the job: new `flai/internal/serve/statejob.go` and `statejob_test.go`; `resolvejob.go`, `flai/internal/sendback/sendback.go` and its test, and `flai/internal/closeout/closeout.go` and its test, which S-0804 writes; `serve.go`, which wires the pass; `agents_restart.go`, which starts the agent on a send-back thread.
- Design, the suites: `flai/cmd/verify.go` and `flai/internal/verify/story.go`, for a `flai verify` limited to a state's suites, and `docs/users/flai-reference.md`, which is generated from the flags.
- Layout, the cluster: `flai/internal/serve/kluster.go`, whose `lookReview` deploys on `AwaitsAcceptance` today, `kluster_test.go`, `kluster_requests.go`, `kluster_requests_test.go`, `checks_test.go`, `flow.go`, `flow_test.go`.
- Design, the documents that say the cluster deploys on review: `design/system/kluster.md`, `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/runbooks/change-the-board.md`.
- Co-change: `flai touches suggest` from these paths lists `docs/operators/settings.md`, `design/system/flaiover-dashboard.md`, and `flai/cmd/accept.go` first, at 17% or less. None is added: this story adds no setting, the dashboard's view of the job is S-0814's, and S-0804 moves the send-back out of `accept.go`.

Forecast: 2h30m, delivery 2026-11-17T09:47:00Z. `flai forecast` gives 1h16m, the median 110 s a unit times size 41 (5 criteria, 36 touches). Doubled, since three of E-0034's four serve and check stories took 1.9 to 2.7 times their forecasts: S-0743 took 2h02m on a 45m forecast, S-0755 1h30m on 47m, and S-0741 1h19m on 29m; S-0742 took 1h06m on 1h20m. Delivery is flai's 08:33Z plus the 1h14m added; the stories ahead in the pull order may run long too.

Cost of delay: 175.37 USD a week, as `flai cod` gives it. The story has no inputs of its own. Its value is its share of E-0034's 1500.00 USD a week, from the operator's 10h time lost per cycle: 2h30m of the 21h23m forecast over the epic's 20 open stories without inputs. Taken as is: it follows from the forecast above.

Task layers:

1. T-3475 (the `kluster-verified` rule in the move) and T-3479 (the workflow accessors): no `after`, no path in common.
2. T-3482 (the state job, its rebase, suites, record, and send-back), after T-3479, whose accessors it reads.
3. T-3484 (the `kluster:` steps and skips in the job, and no deploy on review in code), after T-3482, whose `statejob.go` it extends.
4. T-3487 (the design and the guides), after T-3475 and T-3484, whose behaviour it describes.

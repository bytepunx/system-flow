---
id: S-0663
type: story
nature: feature
title: The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it
status: ready
parent: E-0035
owner: alex
created: 2026-11-10T08:09:11Z
updated: 2026-11-13T09:21:34Z
transitions:
  - to: ready
    at: 2026-11-13T07:26:40Z
    by: alex
tags: [dashboard]
topics: [kluster, deployment]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/kluster_forward_test.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/kluster/checks/+server.ts, flaiover/src/routes/api/kluster/checks/checks.test.ts, flaiover/src/routes/api/kluster/deployments.test.ts, flaiover/src/routes/api/kluster/forward/forward.test.ts, flaiover/src/lib/components/DeploymentRow.svelte, flaiover/src/lib/components/DeploymentRow.svelte.test.ts, design/components/flaiover/src/lib/components/DeploymentRow.examples.json, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, design/system/kluster.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md]
after: [S-0660, S-0661]
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
      seconds: 83
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 6930
          cache_read: 1331725
          cache_write: 106206
          cost: 0.9966
    - kind: orchestrator
      seconds: 152
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 194
          cache_read: 1132810
          cache_write: 83048
          cost: 0.3311
        - model: claude-sonnet-5-5
          input: 6
          output: 822
          cache_read: 76393
          cache_write: 29498
          cost: 0.0798
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: alex
    at: 2026-11-10T09:37:49Z
  value: 5
  by: planner-S-0663
  at: 2026-11-10T09:41:12Z
forecast:
  duration: 1h30m
  delivery: 2026-11-13T14:07:00Z
  basis: "Its own forecast of 1h30m; 8th in the pull order with an in-progress limit of 3, behind S-0761, S-0763, S-0760, S-0745, S-0759, S-0751, S-0757, S-0740 and S-0661."
  by: flai
  at: 2026-11-13T09:21:34Z
finalized:
  by: orchestrator
  at: 2026-11-10T09:44:45Z
---
# S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it

## Goal

The dashboard's Deployments page shows the namespace of each story in progress beside those in review, and the operator deploys, checks, and tears down from it. The operator sees what agents have running in the cluster and can clear it.

## Acceptance criteria

- [ ] The Deployments page lists each story namespace with its story's column, its commit, whether that is the branch's head, its last deploy and checks run with their outcomes and logs, and who asked for each
- [ ] The operator deploys a story in progress, runs its checks, and tears it down from the page, each through flai, behind the `kluster` host action
- [ ] The page shows how many story namespaces live against the limit, and which request was refused for it
- [ ] `flai kluster deployments` and its `--json` report the same
- [ ] A port-forward to a story in progress starts and stops as one to a story in review does
- [ ] The component examples under `design/components/` are in step with the components changed
- [ ] `design/system/kluster.md`, `design/system/flaiover-dashboard.md`, and `docs/users/flaiover.md` say all of it

## Tasks
- T-2868 flai kluster deployments and its records list each story namespace with its column, head, last deploy and checks run, who asked, and the count against the limit
- T-2869 kluster.deployments answers each story namespace's column, head, last deploy and checks run, and who asked, the stories in progress not deployed, and the count against the limit
- T-2870 The operator deploys, checks, and tears down a story in progress through kluster.deploy, a new kluster.checks, and kluster.teardown, and forwards into its namespace
- T-2871 flaiover's client reads the new deployments fields and calls kluster.checks through a route of its own
- T-2872 DeploymentRow shows one story namespace's column, commit and head, last deploy and checks run, and who asked, with its actions and its examples
- T-2873 The Deployments page lists the stories in progress beside those in review, acts on them through DeploymentRow, and shows the count against the limit and the refused request
- T-2874 kluster.md, flaiover-dashboard.md, flaiover.md, and flai.md say what the Deployments page, kluster.deployments, kluster.checks, and flai kluster deployments do for a story in progress

## Notes

- S-0581 built the page.

### Planning

Planned by planner-S-0663 on 2026-11-10, in two runs. The first wrote the touches, the forecast, and the tasks; the orchestrator approved them on TH-0986. The second wrote the cost of delay value once the operator gave the input on TH-0981, and changed nothing else: the seven tasks still fit the story's outcome.

Touches, file by file, with no folder touch. Each came from the design (`design/system/kluster.md`), co-change (`flai touches suggest` from the page, the host methods, the records, the command, and `kluster.md`: 24 commits), the criteria, or the code layout:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/kluster/deploy.go`, `deploy_test.go` | design, co-change | `Record` and `Deployments` gain who asked, the story's column, and the count against the limit |
| `flai/internal/serve/kluster.go`, `kluster_test.go` | layout | Where S-0660's request is refused over the limit; the refusal is kept for the list |
| `flai/cmd/kluster.go`, `kluster_test.go`, `docs/users/flai-reference.md` | criterion 4, co-change | `flai kluster deployments` and `--json`; the reference is generated from the command's help |
| `flai/internal/hostapi/kluster.go`, `kluster_test.go` | design, co-change | `kluster.deployments` lists stories in progress, `kluster.deploy` takes one, and a new write `kluster.checks` |
| `flai/internal/hostapi/kluster_forward_test.go` | criterion 5 | Forwards do not check the story's column; a test holds that for a story in progress |
| `flai/internal/hostapi/writes_test.go` | layout | The table of writes and their host action gains `kluster.checks` |
| `flaiover/src/lib/server/repo.ts`, `agent.ts` | co-change | The types, the client methods, and the allow-list for `kluster.checks` |
| `flaiover/src/routes/api/kluster/checks/+server.ts`, `checks/checks.test.ts` | layout | A route for the new write, as `forward/` has |
| `flaiover/src/routes/api/kluster/deployments.test.ts`, `forward/forward.test.ts` | co-change | The routes' answers with a story in progress |
| `flaiover/src/lib/components/DeploymentRow.svelte`, its test, `design/components/flaiover/src/lib/components/DeploymentRow.examples.json` | criterion 6 | One namespace's row, drawn out of the page, with its examples |
| `flaiover/src/routes/ops/deployments/+page.svelte`, `deployments.svelte.test.ts` | criteria 1 to 3 and 5, co-change | The page |
| `design/system/kluster.md`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` | criterion 7 | What the page and the methods do |
| `docs/users/flai.md` | co-change, `documentation.md` | A changed command's user documentation changes with it |

`after`:

- S-0660: S-0663 reads its deploy request for a story in progress, the namespace kept until acceptance or cancellation, the limit on story namespaces and its refusal, and who asked, in the namespace's record.
- S-0661: S-0663 reads its checks run against a story in progress's namespace, and the record of that run with its commit, its outcome, and who asked.
- S-0659 is not named: its ADR reaches S-0663 through S-0660, which builds it.

Figures:

- Forecast 1h30m, adjusted, and it stands. `flai forecast` gives 1h: 109 s per unit of size over 118 done large feature stories on claude-opus-5-5, times size 33 (7 criteria, 26 touches). The stories that last changed this page, its host methods, and their docs took longer than that rate gives: S-0581 38m, S-0595 64m, S-0586 2h32m. This story spans all three layers, as S-0595 and S-0586 did, so 1h30m. flai keeps the delivery from the pull order.
- Cost of delay 5.00 USD a week, as `flai cod` gives it: the operator's own input on S-0663, `time_lost_per_cycle: 2m`, at 150 USD an hour over one 168h cycle a week. The operator set it on 2026-11-10 after confirming on TH-0981 the epic's input, `revenue_per_week: 1250` on E-0035. A story's own input takes it out of the epic's share, so S-0663 is valued apart from its siblings, which take their shares of the epic (S-0659 263.16, S-0661 986.84) or their own inputs (S-0660 1250, S-0662 30). Not adjusted: the input is the operator's, and the figure follows from it with no judgement of mine.

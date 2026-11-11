---
id: S-0662
type: story
nature: feature
title: The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:10Z
updated: 2026-11-11T05:24:55Z
transitions: []
tags: [cli, dashboard, template]
topics: [kluster, deployment, review]
touches: [flai/internal/kluster/evidence.go, flai/internal/kluster/evidence_test.go, flai/cmd/move.go, flai/cmd/move_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/review_test.go, flai/internal/verify/story.go, flai/internal/verify/kluster.go, flai/internal/verify/kluster_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flai/internal/preview/orchestrator.go, flai/internal/preview/evidence_test.go, flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, flai/internal/context/story.go, flai/internal/context/context.go, flai/internal/context/context_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/workflow.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, template/CHANGELOG.md, design/system/kluster.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md]
after: [S-0661]
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
      seconds: 134
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 26
          output: 13399
          cache_read: 889246
          cache_write: 75747
          cost: 1.0519
    - kind: orchestrator
      seconds: 166
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 189
          cache_read: 1044801
          cache_write: 85611
          cost: 0.3074
cost_of_delay:
  inputs:
    penalty_per_week: 30
    by: alex
    at: 2026-11-10T09:39:38Z
  value: 30
  by: planner-S-0662
  at: 2026-11-10T09:46:48Z
forecast:
  duration: 2h30m
  delivery: 2026-11-11T17:42:00Z
  basis: "Its own forecast of 2h30m; 15th in the pull order with an in-progress limit of 3, behind S-0678, S-0686, S-0691, S-0693, S-0721, S-0694, S-0681, S-0718, S-0696, S-0697, S-0698, S-0700, S-0702, S-0705, S-0241, S-0659, S-0660 and S-0661."
  by: flai
  at: 2026-11-11T05:24:55Z
finalized:
  by: orchestrator
  at: 2026-11-10T09:50:14Z
---
# S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first

## Goal

The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there before review. A criterion that needs the running system reaches review ticked, or unticked with a reason the operator can weigh.

## Acceptance criteria

- [ ] While the `kluster` host action is on and a `deploy` is set, `flai move <story> review` and `item_move` warn when the story's namespace is not deployed at its branch's head, or its checks did not pass there, naming the commands that do both; they do not refuse
- [ ] Their warning for unticked acceptance criteria says, in such a project, whether the namespace was deployed and checked at the head
- [ ] The close-out records the namespace's last deploy and checks run, each with its commit and outcome, and `flai verify --last` and the MCP `item_get` show them
- [ ] The dashboard's review page shows the same beside the criteria, with a link to each log
- [ ] The orchestrator's acceptance reads the record as evidence for a criterion, as `design/system/strategic-agents.md` says
- [ ] The work-management convention, in `design/system`, in `design/conventions`, and in `template/`, tells a story's agent in a project with a kluster to deploy, check, and tick before review, and to leave unticked only what it could not verify, with the reason
- [ ] A story's prime pack in such a project names the deploy and checks commands
- [ ] `docs/users/flai.md` and `docs/users/flaiover.md` say all of it

## Tasks
- T-2875 kluster.Evidence reads a story namespace's last deploy and checks run and says whether both passed at the branch's head
- T-2877 The work-management convention tells a story's agent in a kluster project to deploy, check, and tick before review
- T-2878 A story's prime pack in a kluster project names the deploy and checks commands
- T-2880 flai move to review and item_move warn when the namespace was not deployed and checked at the head, and the unticked-criteria warning says which
- T-2881 The close-out records the namespace's last deploy and checks run with their commits and outcomes, and flai verify --last prints them
- T-2883 The orchestrator's acceptance takes the namespace's deploy and checks at the verified commit as evidence for a criterion
- T-2885 The MCP item_get shows a story's last deploy and checks run in its namespace, with their commits and outcomes
- T-2886 The dashboard's review page shows the namespace's last deploy and checks run beside the criteria, with a link to each log
- T-2888 kluster.md, flai-cli.md, flaiover-dashboard.md, and the users' guides say what the move, the close-out, the review page, and the pack say of the namespace

## Notes

- `Repo.Move` in `flai/internal/workitem/rules.go` makes the warning for unticked criteria today (S-0282).

### Planning

Written by the planner, 2026-11-10; revised the same day once the operator's input was in.

**After.** S-0661 stays in `after`. S-0662 reads its checks run record, with its commit, outcome, and who asked (S-0661's criterion 7), and names its checks command in the warning and the pack. It also reads S-0660's deploy command and deploy record, and S-0659's ADR on what is recorded before review and where the move, the close-out, and an acceptance read it. S-0661's own `after` names both, so they are not repeated here. Naming S-0659 again would make it a hub of three (`wip.after-hub`) and add no wait.

**Touches.** `flai touches suggest` had nothing to start from on the first run: the story declared no touches. On the second run, from the 40 touches below, nothing it lists that the story lacks reaches 18%. Every touch is a file; no folder is kept.

| Touch | Source |
|-------|--------|
| `flai/internal/kluster/evidence.go`, `evidence_test.go` (new) | layout: `KlusterEnv` in `flai/internal/serve/kluster.go` already judges a passed record at the branch head, but `mcpserver`, `verify`, and `preview` should not import `serve`; `kluster` is importable by all of them |
| `flai/cmd/move.go`, `move_test.go`; `flai/internal/mcpserver/server.go`, `review_test.go` | criteria 1 and 2; the notes' pointer: `Repo.Move` builds the warning, and the CLI (`move.go`) and `itemMove` in `server.go` return it. `workitem` cannot see kluster, so the callers add the sentence |
| `flai/internal/verify/story.go`, `kluster.go`, `kluster_test.go` (new); `flai/cmd/verify.go`, `verify_test.go` | criterion 3: `verify.SaveReport` writes the report the close-out keeps, and `printLastVerify` prints `--last` |
| `flai/internal/mcpserver/server.go`, `server_test.go` | criterion 3: `ItemOut` carries `pull_request`, beside which the record goes |
| `flai/internal/hostapi/writes_host.go`, `writes_test.go`; `flaiover/src/lib/review.ts`, `review.test.ts`, `Review.svelte`, `Review.svelte.test.ts` | criterion 4: `verify.status` feeds the review page's verification section |
| `flai/internal/preview/orchestrator.go`, `evidence_test.go`; `flai/cmd/accept.go`, `accept_orchestrator_test.go` | criterion 5: `ParseEvidence` and `unevidenced` take only changed files, or a formula step's task, as a criterion's evidence |
| `design/system/strategic-agents.md`; `design/conventions/strategic-agents.md` and its template copy; the orchestrator's role file, its template copy, and `flai/internal/roles/baseline/orchestrator.md` | criterion 5 names the design; the role and the convention tell the orchestrator how to write evidence, and `baseline_test.go` holds the embedded copy to the template's |
| `design/conventions/work-management.md`, its template copy, `design/system/workflow.md`, `template/CHANGELOG.md` | criterion 6; a template convention change takes a changelog line |
| `flai/internal/context/story.go`, `context.go`, `context_test.go` | criterion 7: `ForStory` builds the pack and `Pack.Header` prints it; nothing in it reads the manifest yet |
| `design/system/kluster.md`, `flai-cli.md`, `flaiover-dashboard.md`; `docs/users/flai.md`, `flaiover.md` | criterion 8, and E-0035's design; `kluster.md` says today that the close-out does not use the cluster |

Left out on purpose:

- `docs/users/flai-reference.md` (co-change, 17%): no flag is planned. The story's agent adds it with `flai touches` if a flag comes.
- `docs/operators/index.md` (co-change, 16%): the operator's side of the namespace is S-0663's Deployments page.

**Forecast.** `flai forecast` gives 1h27m: 109 s a unit over 118 large feature stories, size 48 (8 criteria, 40 touches). Raised to 2h30m: nine tasks in four layers, and three of them read record shapes that S-0659, S-0660, and S-0661 have not written yet. S-0628, of a like spread, was planned at 2h with six tasks. The delivery, 2026-11-10T19:37:00Z, is flai's play-out of the pull order (11th, in-progress limit 5), 18:34, plus the 1h03m added.

**Cost of delay.** 30 USD a week, as `flai cod S-0662` works it out from the operator's `penalty_per_week: 30` on this story (alex, 2026-11-10T09:39:38Z; TH-0983). Not adjusted: the operator set the input on this story directly, which takes S-0662 out of the share of E-0035's `revenue_per_week: 1250`, and a penalty fits a story that saves review time and earns no revenue of its own.

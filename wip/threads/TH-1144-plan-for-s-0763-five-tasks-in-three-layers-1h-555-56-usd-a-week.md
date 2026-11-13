---
id: TH-1144
title: "Plan for S-0763: five tasks in three layers, 1h, 555.56 USD a week"
anchor:
  path: wip/kanban/stories/S-0763-each-deploy-gives-a-story-s-namespace-a-kubeconfig-held-to-it-and-a-story-s.md
  item: S-0763
status: resolved
participants: [planner-S-0763, orchestrator]
created: 2026-11-13T09:04:50Z
updated: 2026-11-13T09:06:38Z
---

# TH-1144 Plan for S-0763: five tasks in three layers, 1h, 555.56 USD a week

On wip/kanban/stories/S-0763-each-deploy-gives-a-story-s-namespace-a-kubeconfig-held-to-it-and-a-story-s.md.

## Entries

### 2026-11-13T09:04:50Z planner-S-0763
S-0763 is planned: five tasks in three layers, touches widened from 10 to 17 files, forecast 1h, cost of delay 555.56 USD a week. Nothing here needs your answer unless you disagree with an assumption below.

## Tasks and layers

| Layer | Task | What it builds | Waits for |
|-------|------|----------------|-----------|
| 1 | T-3276 | The kubeconfig step after the deploy step passes, `Cluster.AgentEnv`, and its removal at teardown (criteria 1, 2) | Nothing |
| 1 | T-3277 | `flai kluster log`, with the `kluster.log` reader moved from `hostapi/kluster.go` into `flai/internal/kluster/log.go` so both read alike (criterion 4) | Nothing |
| 1 | T-3278 | The guard rules for `exec` and `log`, the `kluster` program, and `.flai-cache/kluster/kubeconfig` (criterion 5) | Nothing |
| 2 | T-3279 | `flai kluster exec` (criterion 3) | T-3276 for `AgentEnv`; T-3277 for `flai/cmd/kluster.go` and the generated reference |
| 3 | T-3280 | kluster.md, flai-cli.md, flai.md, and the settings index (criterion 6) | All four |

## Assumptions

- The kubeconfig step is part of the `deploy` step, not a recorded step of its own. Its output goes to the deploy log, and a failure fails the deploy at `deploy`. This keeps the record's step names, `kluster.log`'s step list, and the Deployments page unchanged.
- The step runs for a story's namespace only, never for a promotion into the default namespace.
- A failed teardown keeps the kubeconfig, as it keeps the record.
- `flai kluster exec` runs the command as an argument list, never through a shell, and exits with the command's exit code. When the kubeconfig is absent but the record exists, it warns and runs anyway.
- No change in `flai/internal/serve`: its jobs call `DeployBy` and `TeardownBy`. No MCP tool: the criteria name none.
- An agent in a container reaches neither the kubeconfig nor the logs this way, since `.flai-cache/kluster` is not mounted. T-3280 says so in the design; this story builds no mount.
- The guard refuses `.flai-cache/kluster/kubeconfig` and the `kluster` program to every agent and sub-agent, strategic agents included, as ADR-0292 §2 lists. The operator's own session is not refused.
- `docs/operators/settings.md` is a touch because `scripts/flai-reference.sh` writes its flag index with the reference, and its kluster variables row names the commands that give them.

## For the operator to see

- S-0763 is 31st in the pull order, after S-0661, S-0662, and S-0663, which all wait for it through S-0661's `after`. Until it moves up, those three stay held. The order is yours; I changed nothing.
- S-0661's T-2864 (the guard refuses another story's namespace) overlaps T-3278. S-0661 waits for S-0763, so T-2864 can build on T-3278 and keep to `env`. S-0661's T-2865 should read `Cluster.AgentEnv` from T-3276. Neither is S-0763's to change.
- `flai check` warns `wip.after-reason` on S-0661: it names S-0763 in `after` but its Notes do not say why.

### 2026-11-13T09:06:24Z orchestrator
Approved: S-0763's plan fits the story.

- The five tasks cover the six criteria: T-3276 covers 1 and 2, T-3277 covers 4, T-3278 covers 5, T-3279 covers 3, and T-3280 covers 6.
- The layers hold. T-3276, T-3277, and T-3278 share no touch. T-3279 waits for T-3277 on `flai/cmd/kluster.go`, its test, and the generated reference, and for T-3276 on `Cluster.AgentEnv`. T-3280 waits for all four, `docs/operators/settings.md` among their files.
- Every task's touches are within the story's 17 file touches.
- The forecast of 1h and the value of 555.56 USD a week come from flai, from the operator's input on E-0035.

One correction, from ADR-0292 §2. Your assumption says an agent in a container reaches neither the kubeconfig nor the logs. The ADR says an agent in a container reaches the logs only through `flai kluster log`. So T-3280 must not say a container agent cannot read the logs. When S-0763's agent starts T-3277, it should check whether `flai kluster log` can read `kluster.log` through the host where `.flai-cache/kluster` is not mounted. If that widens the story, it should ask on S-0763 rather than build it unasked. The ADR does not open the kubeconfig to a container agent, so that part of the assumption stands.

S-0661's `wip.after-reason` warning and the pull order are not this plan's to change.

Source: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md § 2. What the agent may run against its namespace

### 2026-11-13T09:06:38Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: the tasks, touches, and figures fit S-0763. The assumption about container agents and the logs is corrected from ADR-0292 §2.

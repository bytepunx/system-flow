---
id: D-0602
title: "Does S-0660 also build flai kluster env, exec, log, and the namespace's kubeconfig: Recommended: keep S-0660 to its nine criteria, as you propose, with one change to where the split-out parts go"
kind: answer
sources: [TH-1138, S-0660]
topics: [kluster, deployment, flai, cli, go, code, all]
paths: [flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/manifest/manifest_kluster.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/issues/I-0355-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md, design/issues/summary.md, flai/internal/hostapi/kluster_test.go]
date: 2026-11-13
status: absorbed
absorbed_by: "design/system/kluster.md#what-the-agent-runs-against-its-namespace-s-0660-s-0661-s-0763"
asked_again: []
cited: []
---

# D-0602 Does S-0660 also build flai kluster env, exec, log, and the namespace's kubeconfig: Recommended: keep S-0660 to its nine criteria, as you propose, with one change to where the split-out parts go

The operator, alex, answered TH-1138 at 2026-11-13T07:57:17Z.

What was asked:

```text
ADR-0292's Consequences assign to S-0660 more than its nine criteria name: besides the requests, the queue entries, the `kluster` event, and `kluster.agent_namespaces`, it lists the namespace's own kubeconfig (a `flai-agent` ServiceAccount, RoleBinding to `edit`, 24h token), `flai kluster env`, `flai kluster exec`, and `flai kluster log`, with the guard rules for them (§2). `design/system/kluster.md` marks those as S-0660's too. The story's criteria and the planner's six tasks cover none of them.
```

The answer:

```text
**Recommended: keep S-0660 to its nine criteria, as you propose, with one change to where the split-out parts go.** A new backlog story under E-0035 builds the namespace's kubeconfig (the `flai-agent` ServiceAccount, the RoleBinding to `edit`, the 24h token, made after `deploy` passes and deleted at teardown), `flai kluster exec`, and `flai kluster log`, with their guard rules. `kluster.md` names that story in place of S-0660 for those parts. S-0661 gains it in `after`.

**Why `env` is not in the new story:** S-0661 already builds `flai kluster env` (T-2865, its criterion 3, "in the way the ADR decides"). Two stories should not build the same command.

**Why S-0661 must wait for the new story:** what `env` gives includes `KUBECONFIG`. The operator confirmed on TH-1103 (D-0599, choice 5) that only the namespace-scoped kubeconfig holds the boundary, because the guard cannot read inside the agent's scripts. If `env` landed first, it would have only flai's whole kubeconfig to hand over.

**Why keep S-0660 to its criteria:** the criteria are the acceptance contract. The requests are useful alone: the deploy runs through flai serve with flai's kubeconfig, and the agent holds no credential. ADR-0292's Consequences say what each story builds. `kluster.md`, which is living design, can move these parts to another story without changing the ADR's decision.

**Alternative:** add them to S-0660 as two more tasks and four criteria, as you list.

Go on with the criteria's work. This changes only whether the two tasks are added at the end.

Source: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md § Consequences
```

What came of it: S-0660's narrative records that the operator confirmed the recommendation. S-0660 kept to its nine criteria. S-0763 (backlog, E-0035, after S-0660) builds the namespace's kubeconfig, `flai kluster exec`, and `flai kluster log`. `flai kluster env` stays S-0661's, and S-0661 gained S-0763 in `after`. `design/system/kluster.md` § What the agent runs against its namespace now names them so.

Sources: TH-1138; S-0660 (`wip/archive/agents/S-0660.md`, Decisions).

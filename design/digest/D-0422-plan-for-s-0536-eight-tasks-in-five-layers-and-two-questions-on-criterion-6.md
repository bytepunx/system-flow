---
id: D-0422
title: "Plan for S-0536: eight tasks in five layers, and two questions on criterion 6: confirm recommendation"
kind: recommendation
sources: [TH-0717, S-0536, S-0537]
topics: [orchestration, cli, template, conventions, mcp, dashboard, flai, go, code, all]
paths: [flai/internal/serve, flai/internal/mcpserver, flai/internal/harness, flai/internal/manifest, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/project-manifest.md, docs, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/check/roles_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/orchestrate_timeout.go, flai/internal/serve/orchestrate_timeout_test.go, flai/internal/serve/serve.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/orchestrator_end_test.go, design/system/agent-narrative.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/roles/baseline/orchestrator.md, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/role.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: [TH-0719]
cited: []
---

# D-0422 Plan for S-0536: eight tasks in five layers, and two questions on criterion 6: confirm recommendation

The operator, alex, confirmed the planner's recommendation on TH-0717 at 2026-10-29T13:52:32Z ("confirm recommendation"). It was not overruled.

What was decided for S-0536, and why:

1. A story's acceptance criterion that needs the host's flai serve to run the story's own change, such as a day-long run of the orchestrator under a new trigger, cannot be verified before review. The host's flai serve runs the installed release, not the story branch's build, so the change reaches the host only once the story is accepted and published. Such a run moves to a follow-up story after the story, and the criterion keeps only what the agent can verify on the branch (here, ADR-0209's status).
2. A folder touch on `design/adrs` is dropped when no task writes an ADR: while the story is in progress it claims every ADR and holds any ready story that adds one (ADR-0096).
3. `flai stats` does not report context tokens per action. A measurement that needs them takes the cost and activities from `flai stats` and the tokens from the session logs, and says which came from where. Adding the figure to `flai stats` changes `metrics.md`, which needs an ADR.

The recommendation the operator confirmed:

```text
1. Drop `design/adrs` from S-0536's touches. No task writes an ADR, since ADR-0209 is already accepted.
2. Keep criterion 6 to ADR-0209's status, and move the day-long run into a follow-up story after S-0536. The host's flai serve runs the installed release, not the branch's build. So the new trigger can only run here once S-0536 is accepted and published, and a criterion the story's agent cannot verify before review would block its acceptance.
3. T-2307, or its follow-up, takes cost and activities from `flai stats` and the context tokens from the session logs, saying which came from where.

Source: design/conventions/work-management.md § Rules
```

S-0536's agent then acted on it, as its narrative's Decisions record: "`design/adrs` is dropped from the story's touches (TH-0717, proposal 1). No task writes an ADR, and a story's agent may narrow its own claim (work-management.md)", and "criterion 6 is narrowed to ADR-0209's status, and T-2307 is cancelled. The day-long run under the new trigger becomes a follow-up story under E-0030, recorded in S-0536's Notes." S-0536 was accepted on 2026-10-29.

Source: TH-0717; S-0536 (wip/archive/agents/S-0536.md, Decisions).

TH-0719 asked the same for S-0537's criterion 5, a day's run of reason-primed packs on the host, and the planner's recommendation named this answer: "answer criterion 5 as TH-0717 answers S-0536's criterion 6." The operator confirmed it at 2026-10-29T16:27:25Z (D-0423). S-0537's narrative records: "T-2315 is cancelled. Criterion 5 is removed. The day of runs joins S-0536's follow-up under E-0030." S-0537 was accepted on 2026-10-29.

Source: TH-0719; S-0537 (wip/archive/agents/S-0537.md, Decisions).

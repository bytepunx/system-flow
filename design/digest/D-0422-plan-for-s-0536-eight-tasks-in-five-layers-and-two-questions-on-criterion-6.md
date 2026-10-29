---
id: D-0422
title: "Plan for S-0536: eight tasks in five layers, and two questions on criterion 6: confirm recommendation"
kind: recommendation
sources: [TH-0717]
topics: [orchestration, cli, template, conventions, mcp, dashboard, flai, go, code, all]
paths: [flai/internal/serve, flai/internal/mcpserver, flai/internal/harness, flai/internal/manifest, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/project-manifest.md, docs, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/check/roles_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/orchestrate_timeout.go, flai/internal/serve/orchestrate_timeout_test.go, flai/internal/serve/serve.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/orchestrator_end_test.go, design/system/agent-narrative.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/roles/baseline/orchestrator.md, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/role.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0422 Plan for S-0536: eight tasks in five layers, and two questions on criterion 6: confirm recommendation

The operator, alex, answered TH-0717 at 2026-10-29T13:52:32Z.

What was asked:

```text
S-0536 is planned: eight tasks in five layers, with touches, a forecast of 5h, and a cost of delay of 300.00 USD a week (the 2h input set on TH-0716).
```

The recommendation it overruled:

```text
**Recommended:** take all three proposals.

1. **Drop `design/adrs` from S-0536's touches.** No task writes an ADR, since ADR-0209 is already accepted. While the story is in progress, the folder touch claims every ADR and holds S-0537 and any story that adds one (ADR-0096).
2. **Keep criterion 6 to ADR-0209's status, and move the day-long run into a follow-up story after S-0536.** The host's flai serve runs the installed release, not the branch's build. So the new trigger can only run here once S-0536 is accepted and published, and a criterion the story's agent cannot verify before review would block its acceptance. The follow-up is one story under E-0030: a day of runs, with the reasons, the cost, and the context tokens recorded.
3. **T-2307, or its follow-up, takes cost and activities from `flai stats` and the context tokens from the session logs, saying which came from where.** Adding tokens per action to `flai stats` changes `metrics.md`, which needs an ADR, and that is beyond this story.

**Or:** keep criterion 6 as written and run the branch's build on the host for a day before review. Say so here, and T-2307 stays as drafted.

Source: design/conventions/work-management.md § Rules
```

The operator's answer:

```text
confirm recommendation
```

Source: TH-0717.

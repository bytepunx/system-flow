---
id: D-0450
title: "S-0420 criterion 2: role.stop and role.start do not exist; the orchestrator alone has a hold: Recommended: take the agent's recommendation"
kind: answer
sources: [TH-0797]
topics: [dashboard, planning, orchestration, cli, template, conventions, flai, go, flaiover, sveltekit, code, all]
paths: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, flai/internal/hostapi/contract_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/agents.go, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flaiover/src/lib/components/RolePage.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/StrategicSettings.svelte, flaiover/src/lib/components/StrategicSettings.svelte.test.ts, flaiover/src/lib/roles.test.ts, flaiover/src/lib/roles.ts, flaiover/src/lib/settings.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/sitemenu.ts, "flaiover/src/routes/api/roles/[name]/roles.test.ts", "flaiover/src/routes/api/roles/[name]/+server.ts", flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, "flaiover/src/routes/workflow/[role]/+page.svelte", flai/internal/serve/agents_test.go, flaiover/src/lib/server/agent.ts, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, docs/operators/settings.md, flaiover/src/lib/strategic.ts, flaiover/src/lib/strategic.test.ts, "flaiover/src/routes/api/agent-stream/[story]/+server.ts", "flaiover/src/routes/api/agent-stream/[story]/agent-stream.test.ts", flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts, design/system/project-manifest.md, design/system/flai-cli.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0450 S-0420 criterion 2: role.stop and role.start do not exist; the orchestrator alone has a hold: Recommended: take the agent's recommendation

The operator, alex, answered TH-0797 at 2026-11-02T00:40:41Z.

What was asked:

```text
Criterion 2 says the role page offers **Run**, **Stop**, and **Start** "each through `role.run`, `role.stop`, and `role.start`". Only `role.run` exists. S-0418 dropped `role.stop`, `role.start`, and the `continuous` trigger (ADR-0209 §1). The only hold today is the orchestrator's, through `orchestrate.stop` and `orchestrate.start` (S-0228).
```

The answer:

```text
**Recommended:** take the agent's recommendation. Reword S-0420 criterion 2 as follows. The role page offers **Run** for any role through `role.run`. It offers **Stop** and **Start** only for a role that has a hold, shown where `agent.status` says the role can be held. Today that is the orchestrator alone, through `orchestrate.stop` and `orchestrate.start`. S-0420's agent then ticks the criterion on that reading.

## Why

- **The ADR keeps the only hold.** ADR-0209 § Decision point 1 makes every role's trigger a list of events and rejects `continuous`. Point 7 keeps `flai serve orchestrate stop|start` as the one hold. D-0425 (TH-0749) dropped `role stop|start` from S-0418 for this reason. Criterion 2 was written before that and names methods that no longer exist.
- **A per-role hold is a decision not yet recorded.** It would need a hold that every event and schedule start obeys, `flai serve role stop|start`, and two host methods. That belongs in its own story, with an ADR if wanted. It does not fit S-0420's claim.
- **The page is ready for a hold later.** The agent is building Stop/Start to show wherever `agent.status` reports `held`. A later per-role hold needs no page change.

Rewording a finalized story's criterion changes its scope, so this is a recommendation for you rather than an answer.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § decision
```

Source: TH-0797.

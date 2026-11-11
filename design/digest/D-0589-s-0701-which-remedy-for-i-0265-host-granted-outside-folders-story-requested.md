---
id: D-0589
title: "S-0701: which remedy for I-0265 — host-granted outside folders, story-requested within a host list, or a prompt warning only: Recommended: remedy A, the host grants outside folders per harness, as…"
kind: answer
sources: [TH-1025]
topics: [flai, cli, go, code, all]
paths: [design/adrs, design/adrs/README.md, design/system/agent-adapters.md, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_agent.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md, flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, docs/operators/index.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md, design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md, design/issues/summary.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: [TH-1025]
---

# D-0589 S-0701: which remedy for I-0265 — host-granted outside folders, story-requested within a host list, or a prompt warning only: Recommended: remedy A, the host grants outside folders per harness, as…

The operator, alex, answered TH-1025 at 2026-11-11T09:00:54Z.

What was asked:

```text
The fix decides what a story's agent may reach on your host, which ADR-0038 makes yours alone ("a story names only the tunables its adapter checks"). I need your choice before I draft S-0701's tasks.
```

The answer:

```text
**Recommended:** remedy A, the host grants outside folders per harness, as the plan is drafted.

- ADR-0038 makes what an agent may reach on your host yours alone. A keeps it so: only `flai serve agent harness --outside` widens it, on the host.
- B adds a story field that anyone who can edit a story in the dashboard writes. It still needs your host list, so it is A plus about half again the work.
- C changes nothing a story can do. I-0265's two instances in a day would recur with E-0032's kluster stories (D-0480).
- Allow T-3059's one scratch `claude -p` on the cheapest model to confirm that a sub-agent inherits `--add-dir`. It costs a few cents, and criterion 1 rests on it.

This is host configuration and a paid call, so it is yours alone.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md
```

Source: TH-1025.

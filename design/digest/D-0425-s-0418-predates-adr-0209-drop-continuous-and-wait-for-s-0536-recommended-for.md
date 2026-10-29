---
id: D-0425
title: "S-0418 predates ADR-0209: drop continuous and wait for S-0536: Recommended for the operator: answer A"
kind: answer
sources: [TH-0749]
topics: [cli, planning, orchestration, template, conventions, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/system/strategic-agents.md, design/system/dashboard-host-channel.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/role_continuous.go, flai/internal/serve/role_continuous_test.go, flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go, flai/internal/serve/role.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flaiover/src/lib/server/agent.ts, template/CHANGELOG.md, design/system/flai-cli.md, design/issues/I-0232-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done-message-from-another-story-each-needing-only-a-settled-reply.md, design/issues/summary.md]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0755]
---

# D-0425 S-0418 predates ADR-0209: drop continuous and wait for S-0536: Recommended for the operator: answer A

The operator, alex, answered TH-0749 at 2026-10-29T16:35:26Z.

What was asked:

```text
**Recommendation:** answer A. I rewrite S-0418's criteria and tasks to fit ADR-0209 and ADR-0211, move S-0418 back to ready with `after: S-0536`, and work it once S-0536 is accepted.
```

The answer:

```text
**Recommended for the operator:** answer A. Re-scope S-0418 to ADR-0209 and ADR-0211, move it back to ready with `after: S-0536`, and let it be worked once S-0536 is accepted.

## Why A

- **The accepted ADR settles criterion 3.** ADR-0209 § Decision point 1 says "`flai check` rejects `continuous` on any role". A criterion that keeps a `continuous` role running, with `flai role stop|start`, contradicts it, and S-0536's check would fail on it. Option B needs a superseding ADR and undoes S-0536's check, for no use the story names.
- **One queue, not two.** Point 2 makes starts coalesced through a reason queue. S-0536 builds that queue in `orchestrate_reasons.go`, so a project role's event trigger reuses it rather than building a second.
- **Waiting is cheaper than rebasing.** S-0536 changes `serve.go`, `roles.go`, `cmd/role.go`, and the docs S-0418 touches. Working S-0418 against today's main means a heavy rebase once S-0536 lands, and S-0536 is in review awaiting only your acceptance (TH-0748).
- **C loses work already planned.** The `roles.<name>.schedule` key and its check are still needed, and A keeps them.

## Why it is yours

Rewriting a finalized story's criteria and cancelling two of its tasks is a change of scope. I promoted S-0418 at 14:42Z after S-0536's move to review freed its hold, before this conflict was known.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision
```

Source: TH-0749.

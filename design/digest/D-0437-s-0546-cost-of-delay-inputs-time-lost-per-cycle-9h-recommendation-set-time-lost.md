---
id: D-0437
title: "S-0546 cost of delay inputs: time lost per cycle 9h: Recommendation: set `time_lost_per_cycle: 9h` on S-0546 itself, as the planner recommends"
kind: answer
sources: [TH-0773]
topics: [testing, mcp, cli, git, orchestration, planning, flai, go, template, conventions, code, all]
paths: [flai/internal/serve, flai/internal/mcpserver, flai/cmd, scripts/close-out.sh, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md, design/adrs/README.md, design/system/metrics.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/serve.go, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/internal/mcpserver/server.go, flai/internal/metrics/closeouts.go, flai/internal/metrics/closeouts_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/root/scripts/close-out.sh, docs/users/flai-reference.md, docs/operators/index.md, design/conventions/code-quality.md, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/code-quality.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0437 S-0546 cost of delay inputs: time lost per cycle 9h: Recommendation: set `time_lost_per_cycle: 9h` on S-0546 itself, as the planner recommends

The operator, alex, answered TH-0773 at 2026-10-31T08:59:18Z.

What was asked:

```text
**Recommendation:** set `time_lost_per_cycle: 9h` on S-0546: `flai edit S-0546 --time-lost-per-cycle 9h`.
```

The answer:

````text
**Recommendation:** set `time_lost_per_cycle: 9h` on S-0546 itself, as the planner recommends:

```bash
flai edit S-0546 --time-lost-per-cycle 9h
```

## Why 9h on the story

- **The goal measures the loss.** Close-outs re-run after the 600 s tool timeout cost 18.5 of the fortnight's 113 session-hours, about 9h a week. That is agent time lost, not revenue or a penalty.
- **Not on E-0031.** The epic's stories each measure a different loss: S-0551's ready wait, S-0546's re-run close-outs. A share from one epic figure would blur those. Inputs on the epic would also stop flai from taking inputs on any of its stories, S-0551 among them (TH-0772).

## Why this needs you

S-0546 is ready. You moved it there at 2026-10-31T07:24:17Z. flai lets the orchestrator set cost of delay inputs only on a backlog story. It refused this for S-0551 on TH-0772 just now, and for S-0427 on TH-0474 (D-0316). So the input is yours to set. A confirmation of this entry does not write the input. Running the command does.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
````

Source: TH-0773.

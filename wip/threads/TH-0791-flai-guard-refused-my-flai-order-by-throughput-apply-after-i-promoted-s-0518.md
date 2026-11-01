---
id: TH-0791
title: "flai guard refused my flai order --by throughput --apply after I promoted S-0518: apply it, or check the guard"
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: resolved
participants: [orchestrator]
created: 2026-11-01T23:50:39Z
updated: 2026-11-01T23:50:48Z
---

# TH-0791 flai guard refused my flai order --by throughput --apply after I promoted S-0518: apply it, or check the guard

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-11-01T23:50:39Z orchestrator
**Recommended:** run this in the main checkout, so the ready column is in throughput order:

```bash
flai order --by throughput --apply
```

**What happened:**

- I promoted S-0518 to ready under `promote_to_ready`. `flai promote --candidates` listed it at rank 1 (1h15m, throughput).
- Under `order_ready`, I then ran `flai order --by throughput --apply --json`. flai guard refused it with: "the orchestrator never does it, whatever its permissions: it orders the ready column by its policy, with flai order --by <policy> --apply".
- The same call passed at 2026-11-01T23:45:52Z, in my previous run. So the guard's verdict on it changed between the two runs.
- `order_by_policy` puts S-0518 at position 5, after S-0424 (1h10m). It sits at position 9 now, last in ready.

**Why you:** I do not retry a refusal or work around it. If the guard should allow the call, its change is a story.

### 2026-11-01T23:50:48Z orchestrator
Resolved: Moot: agent-S-0518 pulled S-0518 into in-progress at 23:50:15Z, so it left ready, and the remaining eight ready stories are already in throughput order (order_by_policy matches the board). The guard's refusal of flai order --apply is logged in wip/agents/orchestrator.md; it needs a look only if it recurs.

---
id: TH-0731
title: Main fails integration at TestADefinedRoleIsHeldToItsFile, so every merge-queue batch fails until it is fixed
anchor:
  path: flai/internal/guard/roles_test.go
status: open
participants: [orchestrator]
created: 2026-10-29T03:25:12Z
updated: 2026-10-29T03:25:12Z
---

# TH-0731 Main fails integration at TestADefinedRoleIsHeldToItsFile, so every merge-queue batch fails until it is fixed

On flai/internal/guard/roles_test.go.

## Entries

### 2026-10-29T03:25:12Z orchestrator
**Recommended:** fix the test's expected wording on main through a remediation story. Until it lands, hold the merge queue with `flai queue hold`. Batch 20261029-031456 has already sent S-0519 back on this failure (TH-0730), and S-0532 and S-0538 are enqueued behind it.

## What fails

`go test ./internal/guard -run TestADefinedRoleIsHeldToItsFile` fails on main at roles_test.go:139. It fails the same way in the queue's integration tier:

```text
the role reviewer cannot move S-0001 to review: it may make it only under to-backlog (on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog), or to-ready (a move to ready and nowhere else).
```

The test expects the prefix `to-backlog (a move to backlog and nowhere else)`.

## Why

- The guard's refusal is right; only its wording changed.
- S-0535 landed through the queue at 01:07Z and gave the to-backlog condition its call and trigger wording (ADR-0211).
- S-0416's test was written and verified before that, at 6c25aad7. You accepted S-0416 by hand at 02:01Z, outside the queue, so the integration tier never ran on the two together. This is the case S-0529 describes.
- flai 1.75.0, which I released with S-0416, carries the failing test but no behaviour fault.

## What fixes it

- Change the expected reason in `flai/internal/guard/roles_test.go` for the `item_move S-0001 review` and `flai move S-0001 review` cases to S-0535's wording, or to a shorter prefix both share.
- The fix needs a story, which I cannot create. S-0519's agent, started on TH-0730, may also take it as a task on its own branch if you prefer. It is outside S-0519's touches.

I take no other action on the queue: holding and dequeuing are yours.

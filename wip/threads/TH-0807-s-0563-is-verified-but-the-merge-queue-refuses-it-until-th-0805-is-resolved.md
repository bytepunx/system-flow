---
id: TH-0807
title: "S-0563 is verified but the merge queue refuses it until TH-0805 is resolved: resolve TH-0805"
anchor:
  path: wip/kanban/stories/S-0563-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md
  item: S-0563
status: resolved
participants: [orchestrator, alex]
created: 2026-11-02T04:45:43Z
updated: 2026-11-02T04:55:47Z
---

# TH-0807 S-0563 is verified but the merge queue refuses it until TH-0805 is resolved: resolve TH-0805

On wip/kanban/stories/S-0563-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md.

## Entries

### 2026-11-02T04:45:43Z orchestrator
**Recommended:** resolve TH-0805, because its answer is built. Run `flai thread resolve TH-0805 --reason "option 1 built in S-0563 (ADR-0241)"`. Then I enqueue S-0563 on my next run.

**What blocks it.** `queue_add` refused S-0563 with this message:

```text
thread_open: thread TH-0805 is not resolved on S-0563 or its tasks: resolve each before the queue takes it
```

TH-0805 was opened by agent-S-0563. You answered "option 1" at 2026-11-02T04:27:14Z, and the story built it. Its agent ended without resolving the thread. I cannot resolve a thread I did not open.

**Nothing else blocks it.**
- The close-out passed every step, integration and smoke included, at a55f8867df2a, the head of story/S-0563.
- The verifier matched both criteria to changed files:
  - Criterion 1: `flai/internal/context/context.go` and the tests in `brief_test.go` and `prime_test.go`, with ADR-0241, `agent-context.md`, `flai-cli.md`, and the prime docs.
  - Criterion 2: I-0250 is closed, with its Remediation section and `design/issues/summary.md`.

S-0518 waits for this fix to reach main (TH-0804).

### 2026-11-02T04:55:47Z alex
Resolved.

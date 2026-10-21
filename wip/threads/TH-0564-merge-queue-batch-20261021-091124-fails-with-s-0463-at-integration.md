---
id: TH-0564
title: Merge queue batch 20261021-091124 fails with S-0463 at integration
anchor:
  path: wip/kanban/stories/S-0463-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md
  item: S-0463
status: answered
participants: [flai serve, agent-S-0463]
created: 2026-10-21T09:23:55Z
updated: 2026-10-21T09:24:23Z
---

# TH-0564 Merge queue batch 20261021-091124 fails with S-0463 at integration

On wip/kanban/stories/S-0463-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md.

## Entries

### 2026-10-21T09:23:55Z flai serve
Merge queue batch `20261021-091124` fails with S-0463 at `integration`. The queue verified S-0463 alone on main at `2b383acc20fb`, at `92c4b9f89cbd`, the copy of its verified head `71f8ab9aaa0d` rebased in the batch, and `integration` failed:

```text
flai/internal/serve/agents_test.go:1348 TestASettledReplyStartsTheAgentWaitingOnItOnce: started before the reply: [{At:2026-10-21T09:17:59Z Action:agent Method:serve.agent Project:t Root:/tmp/TestASettledReplyStartsTheAgentWaitingOnItOnce69115204/001 By:flai serve RequestID: Outcome:done Detail:started stub-agent (command) again for S-0001 as builder-S-0001, a message on MS-0001 (pid 1889218); log /tmp/TestASettledReplyStartsTheAgentWaitingOnItOnce69115204/003/serve/agents/t-S-0001-20261021T091759Z.log}]
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261021-091124/step-1`.

flai moved S-0463 back to in-progress and took it out of the merge queue. Fix it on `story/S-0463`, run flai verify S-0463, move it to review, and enqueue it again with flai queue add S-0463.

### 2026-10-21T09:24:23Z agent-S-0463
S-0463's change did not cause this failure. `TestASettledReplyStartsTheAgentWaitingOnItOnce` in `flai/internal/serve` is the flake I-0158 records; S-0463 changes only `flai/internal/host/host_test.go` and `design/issues`. I bumped I-0158 (count 4; S-0463 met it twice). I am closing out again and enqueueing S-0463. The flake now sends stories back out of the merge queue, so I-0158 is worth a remediation story at acceptance.

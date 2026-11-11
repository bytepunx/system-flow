---
id: T-3028
type: task
nature: remediation
title: flai accept refuses a story with the blocker kluster while the action is on, and --without-kluster lands it on the operator's word
status: done
parent: S-0694
owner: alex
created: 2026-11-10T10:33:07Z
updated: 2026-11-11T14:22:53Z
transitions:
  - to: ready
    at: 2026-11-11T09:13:58Z
    by: agent-S-0694
  - to: in-progress
    at: 2026-11-11T09:13:59Z
    by: agent-S-0694
  - to: done
    at: 2026-11-11T14:22:53Z
    by: agent-S-0694
stream: S-0694
tags: [flai, kluster, accept]
touches: [flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/move.go, flai/cmd/accept_kluster_test.go, flai/internal/hostapi/reads.go]
after: [T-3024]
usage:
  source: log
  seconds: 1731
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 144
      output: 970
      cache_read: 8576828
      cache_write: 250619
      cost: 3.2679
---
# T-3028 flai accept refuses a story with the blocker kluster while the action is on, and --without-kluster lands it on the operator's word

## Work

This is ADR-0254 §5a's third bullet. It covers every acceptance: the operator's, the orchestrator's (the `accept` MCP tool runs `flai accept --by orchestrator` on the host), and the merge queue's (`flai accept --batch`).

- In `flai/internal/preview`, add `BlockKluster` beside `BlockQueueTiers`. Add a `kluster` gate beside `queueTiersGate` in `AcceptWith`, and fields such as `KlusterOn` and `ChecksOn` to `AcceptOptions`. The gate calls T-3024's `kluster.AcceptGate` with the story branch's head, and its blocker names the last job and its outcome. With the action off, or no `kluster:` block, the gate passes.
- In `flai/cmd/accept.go`, set those fields from the host actions `kluster` and `checks`, as `hostActionOn` in `verify.go` reads them.
- Add `--without-kluster` beside `--without-queue-tiers` in `addAcceptFlags`, which `flai move` shares. Write its help line in `move.go` too. It is the operator's alone: refuse it with `--batch` and `--by orchestrator`, as `--without-queue-tiers` is refused. Record its use in the story's Notes, as `withoutQueueTiersNotes` does.
- Pass the same fields from the dashboard's `accept.preview` in `hostapi/reads.go` and from the orchestrator's enqueue preview in `mcpserver/queue.go`, so both list the blocker.

This task waits for T-3024, whose gate it calls. It runs beside T-3025 and T-3029, since none of their files overlap.

## Done when

- `accept_kluster_test.go` shows that each case below is refused with `kluster` while the action is on: the operator's acceptance, `--by orchestrator`, and `--batch`. The cases are a namespace never deployed, a deploy at an older head, a failed deploy, and checks failed at the head.
- It shows that each passes once the deploy and the checks passed at the head, or with checks off.
- It shows that `--without-kluster` lands the operator's acceptance and writes a line in the story's Notes, and is refused with `--batch` and with `--by orchestrator`.
- `preview/accept_test.go` covers the gate with the action off and on.
- `flai test` passes on the eight files.

## Notes

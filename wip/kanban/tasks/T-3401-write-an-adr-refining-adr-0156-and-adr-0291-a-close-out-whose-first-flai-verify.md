---
id: T-3401
type: task
nature: remediation
title: "Write an ADR refining ADR-0156 and ADR-0291: a close-out whose first flai verify stops at sync syncs and verifies again within its three rounds"
status: done
parent: S-0798
owner: alex
created: 2026-11-16T01:43:28Z
updated: 2026-11-16T02:17:44Z
transitions:
  - to: ready
    at: 2026-11-16T02:08:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T02:08:25Z
    by: system-flow
  - to: done
    at: 2026-11-16T02:17:43Z
    by: agent-S-0798
stream: S-0798
tags: [close-out, merge-queue]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 558
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 85
      cache_read: 1504836
      cache_write: 20553
      cost: 0.578
---
# T-3401 Write an ADR refining ADR-0156 and ADR-0291: a close-out whose first flai verify stops at sync syncs and verifies again within its three rounds

## Work

Write the ADR with `flai adr new`, refining ADR-0156 and ADR-0291, and accept it with `flai adr accept`. ADR-0156 §1 syncs only when the last check, `--sync-only`, stops at sync. ADR-0291 §4 tells the operator and the orchestrator that asking for the close-out again syncs the branch onto main. Today it does not when the branch was behind main from the start.

- Context: I-0383's instance. The close-out jobs of S-0739, S-0742, and S-0747 were 37 to 51 commits behind main and each ended in under a second at sync. Each left a failed report as the story's last, and the queue then refused all three as `verify_failed` (TH-1226, TH-1227, TH-1228).
- Decision: when a round's first step, `flai verify --record-issues`, stops at sync, the close-out commits what is outstanding and runs `flai stream sync`. When the sync ends clean, it verifies again within the same three rounds, and stops at `flai stream sync` on a conflict. The script and the template's do this, and so does a container story's close-out on the host (ADR-0274). ADR-0291 §4's words, and the `verify_behind_main` blocker's, then hold as written.
- Alternatives considered: keep a passing report as the story's last when a close-out stops at sync before any tier ran; or have the queue's `verify_failed` tell a stop at sync apart. Say why each is not needed once the close-out syncs: a stop at sync then lasts only while the sync conflicts, and the story needs its agent then.

The ADR's number is the next free one at the time, 0309 when S-0796's ADR-0308 has landed. The file's name follows the title flai is given. Say the name in the task's Notes when it differs from the touches. `flai task done` widens the touches to what is committed.

This task waits for none: it records the decision the story's goal states.

## Done when

- The ADR is accepted, refines ADR-0156 and ADR-0291, and `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass on it.

## Notes

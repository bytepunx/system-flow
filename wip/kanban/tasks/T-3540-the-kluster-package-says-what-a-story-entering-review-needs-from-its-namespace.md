---
id: T-3540
type: task
nature: remediation
title: "The kluster package says what a story entering review needs from its namespace's record: nothing, the checks alone, or a deploy"
status: in-progress
parent: S-0835
owner: alex
created: 2026-11-17T08:30:46Z
updated: 2026-11-28T00:26:25Z
transitions:
  - to: ready
    at: 2026-11-28T00:26:25Z
    by: agent-S-0835
  - to: in-progress
    at: 2026-11-28T00:26:25Z
    by: agent-S-0835
stream: S-0835
tags: [flai]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
---
# T-3540 The kluster package says what a story entering review needs from its namespace's record: nothing, the checks alone, or a deploy

## Work

Waits for no task. The review pass decides today with `kluster.AcceptGate` in `flai/internal/kluster/deploy.go`, which answers only passed or not, so a record whose deploy passed at the head and whose checks failed or never ran reads as a deploy to make (I-0362).

1. Add one judgement beside `AcceptGate`, such as `kluster.ReviewNeeds(git, root, m, story, head, checksOn)`, that answers one of three: nothing, when the deploy passed at the head and, while `checks` is on, its checks too; the checks alone, when the deploy passed at the head and `checks` is on but the checks failed, never ran, or ran at a commit not at the head; a deploy otherwise (no record, a deploy not at the head, a failed deploy, or a record git cannot judge). At the head is `kluster.AtHead`'s, as ADR-0292 §3 says for this row.
2. Keep `AcceptGate` as it is, or have it read the same judgement, so the gate and the review pass cannot disagree.
3. Table tests in `deploy_test.go` beside `TestAcceptGate`: each of the three answers, with checks on and off, and a deploy at a commit followed only by `design`, `docs`, and `wip` changes that answers the checks alone.

## Done when

- `flai test flai/internal/kluster` passes.
- A test shows a record whose deploy passed at the head and whose checks failed there answers the checks alone, and the same record with `checks` off answers nothing.

## Notes

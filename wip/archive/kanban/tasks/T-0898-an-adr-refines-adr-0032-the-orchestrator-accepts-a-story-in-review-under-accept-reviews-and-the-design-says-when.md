---
id: T-0898
type: task
nature: feature
title: "An ADR refines ADR-0032: the orchestrator accepts a story in review under accept_reviews, and the design says when"
status: done
parent: S-0221
owner: alex
created: 2026-10-05T04:47:21Z
updated: 2026-10-06T11:18:01Z
transitions:
  - to: ready
    at: 2026-10-06T11:16:19Z
    by: agent-S-0221
  - to: in-progress
    at: 2026-10-06T11:16:20Z
    by: agent-S-0221
  - to: done
    at: 2026-10-06T11:18:01Z
    by: agent-S-0221
stream: S-0221
tags: [flai]
touches: [design/adrs, design/system/workflow.md, design/system/strategic-agents.md]
usage:
  source: log
  seconds: 101
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 40
      output: 13402
      cache_read: 2057600
      cache_write: 57522
      cost: 1.0171
---
# T-0898 An ADR refines ADR-0032: the orchestrator accepts a story in review under accept_reviews, and the design says when

## Work

Write the ADR with `flai adr new`, refining ADR-0032, which leaves acceptance to the operator alone. Decide and record:

- Under `orchestration.permissions.accept_reviews` (S-0218), the orchestrator may accept a story in review. Its done transition is recorded `by: orchestrator`.
- The four conditions, each a blocker when it fails:
  - the verifier's run passed at the story branch's head commit;
  - every acceptance criterion is ticked;
  - every file the branch changes is under the story's touches, read as a claim reads them;
  - no thread on the story or its tasks is open.
- How flai knows the verifier passed. The plan thread recommends this: the orchestrator runs its verifier sub-agent on the story's worktree, then gives `flai accept` the commit it verified and the evidence. flai refuses when that commit is not the branch head.
- The evidence: the verified commit, the verifier's verdict, and for each criterion the changed files that meet it. It is written to the accepted story's Notes and to the orchestrator's activity log.
- A story whose criteria the orchestrator cannot check against the diff is never accepted. The orchestrator leaves it in review with a thread saying what is missing.
- With the permission off, the guard refuses the orchestrator's `flai accept` and names `accept_reviews`.

Then update the design to match:

- In `design/system/workflow.md` § Transitions and who makes them, the review-to-done row names the orchestrator under `accept_reviews` and links the ADR.
- In `design/system/strategic-agents.md`, a section under the orchestrator on acceptance gives the conditions, the evidence, and the refusal.

The other tasks wait for this one because it fixes the conditions and the evidence format they implement.

## Done when

- The ADR is accepted in `design/adrs` with `refines` naming ADR-0032, and `design/adrs/README.md` lists it.
- `workflow.md` and `strategic-agents.md` describe the orchestrator's acceptance and link the ADR.
- `flai check --strict` and `scripts/lint-md.sh` pass.

## Notes

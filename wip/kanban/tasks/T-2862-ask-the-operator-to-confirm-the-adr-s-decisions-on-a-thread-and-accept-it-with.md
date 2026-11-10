---
id: T-2862
type: task
nature: feature
title: Ask the operator to confirm the ADR's decisions on a thread, and accept it with their answer
status: backlog
parent: S-0659
owner: alex
created: 2026-11-10T08:21:19Z
updated: 2026-11-10T08:21:19Z
transitions: []
stream: S-0659
tags: [cli]
touches: [design/adrs/README.md, design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md]
after: [T-2861]
---
# T-2862 Ask the operator to confirm the ADR's decisions on a thread, and accept it with their answer

## Work

Open one thread on S-0659 naming the ADR's choices an operator could reasonably make otherwise, each with its recommendation first: the limit on story namespaces and its default, what a request over it gets, how long the agent waits behind the queue, and when a failure reaches the operator. The cluster is shared infrastructure on the operator's machine (ADR-0254 §7), so these are theirs to confirm. End the session on the question as `work-management.md` says.

On the answer, change the ADR and `design/system/kluster.md` where it differs, then accept the ADR with `flai adr accept`, which dates it and writes `design/adrs/README.md` again. Where the ADR supersedes ADR-0254 rather than refining it, check that ADR-0254's `superseded_by` names it. It waits for T-2861 so the operator confirms the ADR and the living design together, and because both write `design/adrs/README.md`.

## Done when

- The operator's answer is on the thread and the ADR says what it decided.
- The ADR's status is accepted, and `design/adrs/README.md` lists it so.
- ADR-0254 names it in `superseded_by`, or the ADR names the sections of ADR-0254 it amends.

## Notes

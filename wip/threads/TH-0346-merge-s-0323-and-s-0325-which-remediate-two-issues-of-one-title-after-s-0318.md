---
id: TH-0346
title: Merge S-0323 and S-0325, which remediate two issues of one title, after S-0318?
anchor:
  path: wip/kanban/stories/S-0323-flai-check-finds-narrative-state-outside-the-story-at-close-out.md
  item: S-0323
status: resolved
participants: [orchestrator, alex]
created: 2026-10-08T00:11:36Z
updated: 2026-10-08T04:29:26Z
---

# TH-0346 Merge S-0323 and S-0325, which remediate two issues of one title, after S-0318?

On wip/kanban/stories/S-0323-flai-check-finds-narrative-state-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T00:11:36Z orchestrator
Recommendation: cancel S-0325 as a duplicate of S-0323, and widen S-0323 to cover I-0111 as well as I-0109, with `after: S-0318`.

- S-0323 (I-0109) and S-0325 (I-0111) have the same title: "flai check finds `narrative.state` outside the story at close-out". Both issues' instances are the same kind of finding, on another open story's narrative.
- planner-S-0318 proposed this on TH-0345. S-0318 now has an approved plan: a scoped check leaves out a `markdown.*` finding on another open story's narrative. The merged story would extend that leave-out to every finding on another open story's narrative, which covers `narrative.state`.
- Planned apart, the two would each write an ADR and the same change to `flai/internal/check/scope.go`.

What I need from you:

- Merging is a change of scope, so it is yours to make. Cancel S-0325 with `flai move S-0325 cancelled --reason "duplicate of S-0323"`. Then add I-0111 to S-0323's goal and a second close criterion for it.
- Meanwhile I will plan S-0323 when `flai plan --candidates` reaches it, and I will not plan S-0325.
- Alternative: keep both. Then each gets its own planner run and ADR, and the second rebases onto the first.

### 2026-10-08T04:27:10Z orchestrator
Follow-up: S-0325 is cancelled, and S-0323 has a third criterion for I-0111. One slip in it: criterion 3 reads "I-0111 is closed with `flai issue close I-0109 --reason`". The command should name I-0111. Criteria are yours to edit, so please correct it.

S-0323's plan has no task yet to close I-0111. planner-S-0323 proposed one on TH-0351: a layer-3 task after T-1315, touching the I-0111 file and `design/issues/summary.md`. The story's agent can add it when it starts S-0323. S-0323 waits for S-0318 either way.

### 2026-10-08T04:27:26Z alex
325 was cancelled, 323's criteria was expanded to include I-0111

### 2026-10-08T04:29:26Z orchestrator
Resolved: Merged as recommended: S-0325 cancelled, and S-0323 widened with criterion 3 for I-0111. Criterion 3's command still reads I-0109 where I-0111 is meant, for S-0323's agent to correct with the operator

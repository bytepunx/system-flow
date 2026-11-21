---
id: I-0410
title: The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right
class: efficiency
status: open
count: 1
cost: 20m
first_reported: 2026-11-21T18:38:23Z
last_reported: 2026-11-21T18:38:23Z
updated: 2026-11-21T18:38:23Z
---

# I-0410 The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right

## Description
The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right

## Instances

### 2026-11-21T18:38:23Z
Story: S-0811.
S-0811 passed review at 1c58c4b, but accept's dry run named undoes_main in design/system/workflow.md and docs/operators/runbooks/change-the-board.md. The one line removed from each was main's 'nothing else reads the keys yet' sentence, which the branch rewrote on its rebase to add the publish, and which cannot be kept word for word because it lists the publish as still to come. The orchestrator's accept tool takes no --keep-removal, and flai stream restore would bring the stale sentence back beside the new one, so it opened TH-1321 asking the operator to run flai accept --keep-removal by hand. Meanwhile S-0812 was accepted, S-0811 was sent back to rebase over it (TH-1322), and the new rebase gives the same two removals again. With accept_reviews on, a removal the story's narrative records as a rewrite, and that the orchestrator checked, could be kept by the orchestrator with its reason logged.

## Remediation

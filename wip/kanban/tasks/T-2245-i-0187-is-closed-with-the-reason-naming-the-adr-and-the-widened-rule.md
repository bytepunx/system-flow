---
id: T-2245
type: task
nature: remediation
title: I-0187 is closed with the reason naming the ADR and the widened rule
status: backlog
parent: S-0523
owner: alex
created: 2026-10-28T14:33:05Z
updated: 2026-10-28T14:33:05Z
transitions: []
stream: S-0523
tags: [issues, orchestrator, cost-of-delay]
touches: [design/issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md, design/issues/summary.md]
after: [T-2244]
---
# T-2245 I-0187 is closed with the reason naming the ADR and the widened rule

## Work

Fourth layer: waits for T-2244, so that the issue closes only once the behaviour, the ADR, and the documents are all in place. In the story's worktree, run `flai issue close I-0187 --reason "<reason>"` with a reason that names the new ADR and says what fixed it: `flai guard` admits the orchestrator's `item_edit` and `flai edit` cost of delay inputs on an epic as on a story, and `Repo.OrchestratorSetsCostOfDelay` admits an epic or a story in any state but done or cancelled, with or without inputs, inputs only, under `plan_backlog_stories`. The command rewrites the issue file and `design/issues/summary.md`. Then tick the story's two criteria with `flai criteria tick S-0523 <n>`, each once verified.

## Done when

- `design/issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md` has status closed and the reason.
- `design/issues/summary.md` no longer lists I-0187 as open.
- Both acceptance criteria of S-0523 are ticked, and `flai check --strict` is clean on the story.

## Notes

Drafted by planner-S-0523. Run the issue command in the story's worktree: `flai issue close` writes under the checkout it runs in.

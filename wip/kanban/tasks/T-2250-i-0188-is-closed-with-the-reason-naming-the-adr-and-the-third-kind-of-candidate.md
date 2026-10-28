---
id: T-2250
type: task
nature: remediation
title: I-0188 is closed with the reason naming the ADR and the third kind of candidate
status: backlog
parent: S-0524
owner: alex
created: 2026-10-28T14:42:56Z
updated: 2026-10-28T14:42:56Z
transitions: []
stream: S-0524
tags: [flai, issues]
touches: [design/issues/I-0188-flai-plan-candidates-lists-no-epic-or-story-planned-before-so-the-orchestrator-cannot-ask-the-planner-to-plan-one-again.md, design/issues/summary.md]
after: [T-2249]
---
# T-2250 I-0188 is closed with the reason naming the ADR and the third kind of candidate

## Work

Fourth layer: waits for T-2249, so that the fix it names is built and recorded in the design before the issue closes.

In the story's worktree, run `flai issue close I-0188 --reason "<reason>"`, the reason naming the ADR T-2246 wrote and saying that `flai plan --candidates` now lists an epic or story planned before that the operator has written on since its newest planner run, with the thread, so the orchestrator can ask the planner for it under `plan_backlog_epics` or `plan_backlog_stories` (T-2247, T-2248). The command rewrites the issue's file and `design/issues/summary.md`; commit both with `flai task done`. Then tick the story's criteria you verified with `flai criteria tick S-0524 <n>`: criterion 1 on T-2247's reproducing test, criterion 2 on the close.

## Done when

- I-0188's status is `closed`, its file records the reason, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` scoped to S-0524 passes, and both acceptance criteria of S-0524 are ticked.

## Notes

`design/issues` is in the manifest's shared paths, so the two files hold no ready story (ADR-0096). Run the issue command from the story's worktree, since it writes under the current checkout.

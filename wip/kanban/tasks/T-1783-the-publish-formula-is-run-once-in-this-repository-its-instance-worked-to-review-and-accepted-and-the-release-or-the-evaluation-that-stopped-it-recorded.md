---
id: T-1783
type: task
nature: improvement
title: The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:44Z
updated: 2026-10-21T13:26:52Z
transitions: []
stream: S-0423
tags: [cli]
after: [T-1781, T-1782, T-2014]
---
# T-1783 The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded

## Work

With the `formula` and `push` host actions on, start the formula: run `flai formula run publish --ready` in the main checkout, or let the orchestrator start it under `publish` when `release_evaluate` is met. `flai serve` starts the instance's agent, which works the steps as its tasks say, the main-checkout steps from the main checkout, closes each task as its step passes, and moves the instance to review. Ask the operator on this story's thread to accept it, or let the orchestrator accept it under `accept_reviews`. Record in this story's Notes:

- the instance's ID;
- each step's exit;
- the versions and tags it published, or the evaluation that stopped it when nothing was pending.

Fix what the run finds in a task of its own before this one is done.

Waits for the `push` gate, the orchestrator's rule, and `release_publish`'s refusal, which the run exercises.

## Done when

The instance is accepted and archived, and this story's Notes record the run.

## Notes

The run changes no file of this story: the instance is a story of its own, whose agent's commits are its own, so the task touches nothing.

Rewritten 2026-10-21 by planner-S-0423 on TH-0471's answer 5, as approved on TH-0583: a story's agent works the instance, in place of `flai serve`.

---
id: T-1783
type: task
nature: improvement
title: The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:44Z
updated: 2026-10-16T05:48:44Z
transitions: []
stream: S-0423
tags: [cli]
after: [T-1781, T-1782]
---
# T-1783 The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded

## Work

With the `formula` and `push` host actions on, run `flai formula run publish` in the main checkout, or let the orchestrator start it under `publish`. Watch `flai serve` work the instance to review. Ask the operator on this story's thread to accept it, or let the orchestrator accept it under `accept_reviews`. Record in this story's Notes:

- the instance's ID;
- each step's exit;
- the versions and tags it published, or the evaluation that stopped it when nothing was pending.

Fix what the run finds in a task of its own before this one is done.

Waits for both the publish step and the orchestrator's rule, which the run exercises.

## Done when

The instance is accepted and archived, and this story's Notes record the run.

## Notes

The run changes no file of this story: the instance's items are written in the main checkout by `flai serve`, so the task touches nothing.

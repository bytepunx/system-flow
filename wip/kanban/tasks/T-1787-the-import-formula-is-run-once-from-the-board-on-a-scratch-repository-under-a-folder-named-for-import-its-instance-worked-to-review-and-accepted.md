---
id: T-1787
type: task
nature: improvement
title: The import formula is run once from the board on a scratch repository under a folder named for import, its instance worked to review and accepted
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:49:05Z
updated: 2026-10-16T05:49:05Z
transitions: []
stream: S-0424
tags: [cli, dashboard]
after: [T-1785]
---
# T-1787 The import formula is run once from the board on a scratch repository under a folder named for import, its instance worked to review and accepted

## Work

Make a scratch git repository under a folder the host serves for import (`flai serve import add`), outside this repository. Import it from the board's Import. Watch `flai serve` work the instance to review, and ask the operator on this story's thread to accept it, or let the orchestrator accept it. Record the instance's ID and each step's exit in this story's Notes. Remove the scratch repository afterwards only with the operator's word.

Waits for the import formula.

## Done when

The instance is accepted and archived, and this story's Notes record the run.

## Notes

The run changes no file of this story, so the task touches nothing.

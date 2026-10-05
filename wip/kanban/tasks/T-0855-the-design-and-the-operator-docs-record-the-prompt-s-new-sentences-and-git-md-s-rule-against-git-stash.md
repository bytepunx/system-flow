---
id: T-0855
type: task
nature: improvement
title: The design and the operator docs record the prompt's new sentences and git.md's rule against git stash
status: backlog
parent: S-0268
owner: alex
created: 2026-10-05T00:29:37Z
updated: 2026-10-05T00:29:37Z
transitions: []
stream: S-0268
tags: [docs, flai]
touches: [design/system/flai-cli.md, docs/operators/index.md]
after: [T-0853, T-0854]
---
# T-0855 The design and the operator docs record the prompt's new sentences and git.md's rule against git stash

## Work

In `design/system/flai-cli.md`, find the paragraph on the `flai serve agent` row that begins "Its prompt (`harness.Prompt`". Add a "Since S-0268" sentence there. It says the story agent is told to make independent edits and commands in one turn, and to move a task it has just written to ready and in-progress in one command. It also says that `git.md` now tells the agent how to show a test fails without the change, and not to use `git stash` in a worktree.

In `docs/operators/index.md`, the Harnesses bullet describes what the prompt tells the agent. Add the same in a clause, in the operator's words.

This task waits for T-0853 and T-0854, so that it describes the prompt and the rule as they were written.

## Done when

- `design/system/flai-cli.md` describes the prompt as it now reads, naming S-0268.
- `docs/operators/index.md` says the same for the operator.
- The markdown lint passes.

## Notes

Drafted by the planner for S-0268.

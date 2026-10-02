---
id: ADR-0069
title: "Agents commit each task on the story branch after flai stream sync, and never rebase or merge it by hand"
status: accepted
date: 2026-10-02
supersedes: []
superseded_by: []
refines: [ADR-0019]
---

# ADR-0069 Agents commit each task on the story branch after flai stream sync, and never rebase or merge it by hand

## Context

ADR-0019 put each story on `story/S-nnnn` in a worktree under `.flai-cache/worktrees/` and made `flai stream sync`, which rebases the branch onto `main`, the agent's step at every task transition. The git convention did not say so: it named branches `<type>-<story-id>-<slug>`, which flai never creates, said to commit only when a story moves to review and again at acceptance, and the designer's draft of 2026-10-02 said to rebase from `main` by hand before review. The designer decided on 2026-10-02 that the convention says what flai does.

## Decision

Agents use `flai stream sync` for a story branch's git operations and never rebase or merge it by hand. Each task is committed on `story/S-nnnn` when it is done: the agent runs `flai stream sync`, resolves any conflicts it reports in the worktree, runs the tests for what the task changed, and commits. Before moving the story to review the agent syncs again, resolves conflicts, and commits what is outstanding. This refines ADR-0019; its branches, worktrees, `wip/` in the main checkout, and acceptance stand.

## Consequences

- `design/conventions/git.md` and the template's copy say this, and drop the `<type>-<story-id>-<slug>` branch rule.
- `flai stream sync` must leave the worktree in a state an agent can resolve: it reports each conflicting path and how to continue, and refuses to sync over uncommitted changes rather than losing them.
- `harness.Prompt` for `claude-code`, the work-management convention, and the close-out script agree with it.
- Commits per task make a story's history longer and its review easier to follow task by task.

## Alternatives considered

- Commit only when the story moves to review: conflicts with the designer's edits surface late and large, which is what ADR-0019 set out to avoid.
- Let agents rebase by hand: two ways to do one thing, and a hand rebase does not know about the main checkout's `wip/`.

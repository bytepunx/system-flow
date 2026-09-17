---
id: ADR-0019
title: Agents work each story on a branch in a worktree while wip stays on main
status: accepted
date: 2026-09-17
supersedes: []
superseded_by: []
---

# ADR-0019 Agents work each story on a branch in a worktree while wip stays on main

## Context

Until now all work happened on `main` (git.md project addition). With the designer editing documents in the dashboard while agents implement stories, two writers touch the same files. The operator's preference is that agents branch and that collisions are resolved by rebase and merge, and that front matter can say what a story is working on.

## Decision

`flai stream open` creates `story/S-nnnn` from `main` and a git worktree for it under `.flai-cache/worktrees/`. Code, design, and docs changes for the story land on that branch. `wip/` stays in the main checkout: flai commands run from a worktree write items, narratives, and the board into the main checkout, so the board the dashboard mounts is always current. `flai stream sync` rebases the branch onto `main`; agents run it at every task transition so conflicts with the designer's edits surface small and early. `flai accept` rebases, merges (fast-forward when possible), removes the worktree and branch, then tags and pushes as before, and refuses while a rebase has conflicts.

Stories and tasks may carry `touches`, a list of paths or component names. `flai check` warns when two in-progress items overlap, the board and the dashboard show it, and the MCP server answers who touches a path.

## Consequences

- Commits for a story move to its branch; the git convention's story-landing rule is unchanged in spirit and its project addition is rewritten.
- The dashboard's own edits commit on `main`, which is what makes them visible to `flai stream sync`.
- A story's diff against `main` is well defined, which S-0041 uses for review.
- Worktrees live under `.flai-cache/`, keeping the agent's write footprint inside the repository.

## Alternatives considered

- Everything on the branch, including `wip/`: the board goes stale until merge and the dashboard would have to read every branch; rejected.
- Everyone on `main` with dirty-tree refusals: no isolation and no reviewable diff; rejected by the operator.
- Locks instead of `touches`: a lock nobody releases blocks work; a warning plus a visible badge keeps humans in charge.

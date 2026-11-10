---
id: D-0560
title: "flai issue new on a story branch can take an issue ID another story recorded and folded away meanwhile, since a folded issue's file is removed; S-0636's I-0326 collided with S-0622's and was recorded again as I-0333"
kind: cause
sources: [S-0636, I-0333, TH-0961]
topics: [issues, flai, git]
paths: [design/issues, design/issues/I-0333-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0560 flai issue new on a story branch can take an issue ID another story recorded and folded away meanwhile, since a folded issue's file is removed; S-0636's I-0326 collided with S-0622's and was recorded again as I-0333

S-0636 recorded its 300 ms shortfall as `I-0326` on `story/S-0636` (225420b8b, 2026-11-09 16:09). S-0622 had recorded a different `I-0326` on its own branch (880a14deb) and folded it into I-0323 a second later (3f444437d), which removed the file. Both reached `main` after S-0636's base, so S-0636's worktree had no `I-0326` when it numbered its issue, and the reuse surfaced only on `main`.

The orchestrator caught it on TH-0961 (D-0556): "on `main`, `I-0326` was already used and retired ... Make sure `flai check --strict` accepts it after the rebase on `main`. If it does not, re-record the issue under the next free ID and name that ID in the reworded criterion and in server-performance.md." From S-0636's Decisions (`wip/archive/agents/S-0636.md`):

```text
The shortfall was first recorded as I-0326, but S-0622 had already used and retired that ID on `main`. It is recorded again as I-0333, and the branch drops the I-0326 file.
```

So an issue ID a story records can collide with one a parallel story used and folded away, which leaves no file behind. Before a story names an issue ID it recorded in its criteria or design, it checks `git log --all -- 'design/issues/I-nnnn-*'` for the ID, and re-records under the next free ID when another story used it. D-0210 was the same kind of collision for ADR numbers.

Sources: S-0636 (`wip/archive/agents/S-0636.md`, Decisions); I-0333; TH-0961.

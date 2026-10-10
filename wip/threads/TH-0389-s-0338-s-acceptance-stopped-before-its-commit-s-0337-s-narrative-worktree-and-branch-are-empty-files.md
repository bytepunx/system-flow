---
id: TH-0389
title: "S-0338's acceptance stopped before its commit: S-0337's narrative, worktree, and branch are empty files"
anchor:
  path: wip/archive/kanban/stories/S-0338-a-held-card-says-when-its-holding-story-s-agent-was-asked-about-the-hold-and-a-story-started-on-a-share-names-the-paths-it-shares.md
  item: S-0338
status: open
participants: [orchestrator]
created: 2026-10-08T20:55:51Z
updated: 2026-10-08T20:55:51Z
---

# TH-0389 S-0338's acceptance stopped before its commit: S-0337's narrative, worktree, and branch are empty files

On wip/archive/kanban/stories/S-0338-a-held-card-says-when-its-holding-story-s-agent-was-asked-about-the-hold-and-a-story-started-on-a-share-names-the-paths-it-shares.md.

## Entries

### 2026-10-08T20:55:51Z orchestrator
**Recommendation:** clear S-0337's empty start, then complete S-0338's acceptance yourself with `flai accept S-0338 --by alex`. flai's help says the orchestrator may not complete a failed acceptance.

## What happened

At 20:55Z I accepted S-0338 under `accept_reviews`, at verified commit 2d114ffc. The verify had passed at the head, the verifier matched all five criteria, and the dry-run had no blockers. flai accept did steps 0 to 2: it merged (main is now at 2d114ffc), moved S-0338 to done, archived it with its four tasks and narrative, and closed MS-0042. Then it failed before the commit:

```text
FATAL command failed err="/home/alex/git/bytepunx/system-flow/wip/agents/S-0337.md: no front matter"
```

The main checkout now holds that acceptance uncommitted (`git status`: the archive moves, board.md, E-0018, S-0337's story file, MS-0042, plus the new MS-0043). Any flai write that commits wip on main may sweep these into an unrelated commit.

## The cause: S-0337's start at 11:04Z left only empty files

Each of these is 0 bytes, timestamped 11:04Z, when S-0337 moved to in-progress:

- `wip/agents/S-0337.md`, untracked.
- `.git/refs/heads/story/S-0337`, so `git log story/S-0337` fails with `bad object`.
- `.flai-cache/worktrees/S-0337/.git` and all 4147 files in that worktree (`invalid gitfile format`).

Nothing else in the main checkout is empty. This looks like the host stopped mid-write; my previous run also ended at 11:04:30Z.

## What fixes it

1. Remove S-0337's empty narrative, branch ref, and worktree. Move S-0337 back to ready so flai serve starts it fresh, or restart it.
2. Run `flai accept S-0338 --by alex`. Per `flai accept --help`, it finishes an item that is done and archived but not committed: it commits, then tells the open stories which paths changed.

I will not retry the acceptance. I am publishing nothing until it is committed.

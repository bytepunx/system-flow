---
id: S-0157
type: story
nature: improvement
title: Pending releases are worked out without starting a git process per accepted item
status: backlog
parent: E-0012
owner: alex
created: 2026-09-29T07:00:28Z
updated: 2026-09-29T07:00:28Z
transitions: []
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/release]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0157 Pending releases are worked out without starting a git process per accepted item

## Goal

Every board, from the dashboard or MCP, works out which accepted items are not yet published by starting 25 git processes one after another (13 `diff-tree`, 6 `log`, 6 `tag` on this repository): 75 ms each time, and more as items are accepted between releases. Cause 2 of `design/system/server-performance.md`. Work it out with a few git processes, or once per change of `HEAD` and the tags, and keep the answer.

## Acceptance criteria
- [ ] `release.pending` in a `board.get` event takes under 15 ms on this repository, with at most 3 `exec.git` steps, or none while `HEAD` and the tags are unchanged.
- [ ] `release.PendingIDs` answers exactly what it answers today, for a repository with pending items in several components, none, and a release just cut: a behaviour test compares them.

## Tasks

## Notes

Measured by S-0152: `release.pending=75.4 exec.git.diff-tree=34.7x13 exec.git.log=27.8x6 exec.git.tag=11.3x6`. Proposed: one `git log --name-only` over the range since the oldest component tag instead of a `diff-tree` per commit, and a cache keyed by `git rev-parse HEAD` and the tag list.

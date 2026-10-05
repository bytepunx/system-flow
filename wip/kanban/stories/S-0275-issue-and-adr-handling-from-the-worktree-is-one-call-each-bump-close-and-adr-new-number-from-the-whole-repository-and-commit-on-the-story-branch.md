---
id: S-0275
type: story
nature: improvement
title: "Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:33Z
updated: 2026-10-05T01:35:33Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0275 Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch

## Goal

Recording friction and decisions costs several turns each: `flai issue bump` or `new`, then `git add` and `commit`, then `flai touches` to claim the issue file; `flai adr new` the same, and ADR and issue numbers taken from the worktree alone collide between parallel branches (S-0245, S-0250, S-0252 fix the numbering). `flai issue bump|new|close` and `flai adr new`, run in a story's worktree with `--commit`, number from the whole repository, commit the files they wrote on the story branch with the story's prefix, and widen the story's touches to them, in one call; over MCP they are `issue_bump`, `issue_new`, `issue_close`, and `adr_new`, and the host channel has the same, so the dashboard can file an issue against a story.

## Acceptance criteria
- [ ] `flai issue bump`, `new`, and `close` and `flai adr new` take `--commit`, which commits what they wrote on the story branch and widens the story's touches to it, in one call, as text and `--json`
- [ ] The same operations exist over MCP and on the host channel
- [ ] The conventions, the harness prompt, `design/system/flai-cli.md`, and the user guide send the agent to the one call

## Tasks

## Notes

Waits for S-0250 and S-0252, which make the numbering safe across branches; this story adds only the single-call shape.

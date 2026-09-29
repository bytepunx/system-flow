---
id: S-0164
type: story
nature: feature
title: Reviews need real diffs, not just line counts
status: backlog
parent: E-0014
owner: alex
created: 2026-09-29T21:10:42Z
updated: 2026-09-29T21:10:42Z
transitions: []
tags: [dashboard, cli]
topics: [client-side, server-side]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
---
# S-0164 Reviews need real diffs, not just line counts

## Goal

Each file diff currently shows the path and the lines added/removed.

Add to each of these a toggle control to expand and collapse the detail line to reveal a control that shows the removed lines (- sign in the margin, dim red background behind each line with a brighter red outline around consecutive lines removed) and the added lines (+ sign in the margin, dim green background behind each line with a brighter green outline around consecutive lines added)

## Acceptance criteria
- [ ] a toggle control on the diff summary line reveals or hides the diff panel
- [ ] a `-` in the margin of every subtracted line
- [ ] a `+` in the margin of every added line
- [ ] a dim red background behind removed lines and a brighter red outline surrounding consecutive removed lines
- [ ] a dim green background behind added lines and a brighter green outline surrounding consecutive added lines

## Tasks

## Notes

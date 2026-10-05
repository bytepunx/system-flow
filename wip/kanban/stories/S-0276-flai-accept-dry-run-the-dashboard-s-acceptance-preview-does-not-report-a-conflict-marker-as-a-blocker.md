---
id: S-0276
type: story
nature: remediation
title: flai accept --dry-run, the dashboard's acceptance preview, does not report a conflict marker as a blocker
status: backlog
owner: alex
created: 2026-10-05T03:48:28Z
updated: 2026-10-05T03:48:28Z
transitions: []
tags: []
touches: [flai/internal/preview, flai/cmd/branch.go, design/system/flai-cli.md, docs/users/flai.md]
after: [S-0253]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0276 flai accept --dry-run, the dashboard's acceptance preview, does not report a conflict marker as a blocker

## Goal

The acceptance preview names every reason the acceptance would refuse, so that the operator sees a story branch carrying a conflict marker before pressing Accept rather than when the acceptance refuses it.

## Acceptance criteria
- [ ] `flai accept --dry-run`, and the dashboard's preview that runs it, lists a blocker naming each file and line on the story branch that carries a conflict marker, as `flai accept` refuses it (S-0253's `refuseConflictMarkers`).
- [ ] The preview and the acceptance find markers by one check, `flai/internal/conflictmark` over the files the branch adds or changes, so they cannot disagree.
- [ ] A test shows the preview reporting the blocker for a branch with a marker and none for a clean branch.
- [ ] `design/system/flai-cli.md` and `docs/users/flai.md` say the preview reports it.

## Tasks

## Notes

Raised in S-0253: its acceptance gate refuses a branch with a conflict marker, but `flai/internal/preview` was outside that story, so the preview still offers Accept. The operator asked for this story in S-0253's Decisions (2026-10-05).

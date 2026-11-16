---
id: S-0798
type: story
nature: remediation
title: A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible
status: backlog
owner: alex
created: 2026-11-16T00:41:49Z
updated: 2026-11-16T00:41:49Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: flai
    at: 2026-11-16T00:41:49Z
---
# S-0798 A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible

## Goal

This story remediates [I-0383](../../../design/issues/I-0383-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md), "A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible". The issue recommends this solution:

The close-out a story is asked for when its branch is behind main syncs first: `scripts/close-out.sh` runs `flai stream sync` when its first `flai verify` stops at sync, as it already does when main moves during the run (ADR-0156), within the same three rounds, and stops with the conflicts when the sync does. Until then a close-out that stops at sync before any tier ran should not replace a passing report as the story's last, or the queue's `verify_failed` should tell that case apart. Bring ADR-0291 §4's words and the `verify_behind_main` blocker's in step with what the close-out does.

## Acceptance criteria
- [ ] The cause I-0383 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0383 is closed with `flai issue close I-0383 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0383. time_lost_per_cycle 2h carried over from I-0383's Impact section. Its Impact time_lost_per_cycle was taken rather than the 2h derived from its cost and count.

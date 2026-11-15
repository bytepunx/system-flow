---
id: S-0792
type: story
nature: remediation
title: Every story has to go back to be rebased and retested
status: backlog
owner: alex
created: 2026-11-15T07:37:37Z
updated: 2026-11-15T07:37:37Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12h
    by: alex
    at: 2026-11-15T07:37:37Z
---
# S-0792 Every story has to go back to be rebased and retested

## Goal

Recently every single story in review is getting the following error requiring rebase and re-verify. This is adding fat to much time.

Error: “verify does not cover the main branch it would land on: its last flai verify …”

We must correct this immediately.

## Acceptance criteria
- [ ] Every story in review should not have to rebase and reverify for every story that is approved before it.
- [ ] The only time we need to reverify is when changes between a story accepted overlap with the same change set any given story in review such that the behavior of the PR is undoubtedly impacted by the changes in main that will be rebased on.

## Tasks

## Notes

---
id: S-0517
type: story
nature: remediation
title: A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder
status: backlog
owner: alex
created: 2026-10-28T12:52:47Z
updated: 2026-10-28T12:52:47Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-28T12:52:47Z
---
# S-0517 A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder

## Goal

This story remediates [I-0215](../../../design/issues/I-0215-a-story-s-worktree-mounted-alone-into-a-container-is-not-a-git-repository-since-its-git-file-names-the-main-checkout-s-git-worktrees-folder.md), "A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0215 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0215 is closed with `flai issue close I-0215 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0215. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:27:05Z, 0 days before this story; under one cycle counts as one).

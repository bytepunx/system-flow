---
id: S-0481
type: story
nature: improvement
title: pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY
status: backlog
owner: alex
created: 2026-10-21T12:00:00Z
updated: 2026-10-21T12:00:00Z
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
    at: 2026-10-21T12:00:00Z
---
# S-0481 pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY

## Goal

This story remediates [I-0177](../../../design/issues/I-0177-pnpm-add-run-bare-in-a-story-worktree-installs-from-the-default-store-so-the-next-flaiover-install-sh-aborts-with-err-pnpm-aborted-remove-modules-dir-no-tty.md), "pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0177 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0177 is closed with `flai issue close I-0177 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0177. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T08:25:41Z, 0.1 days before this story; under one cycle counts as one).

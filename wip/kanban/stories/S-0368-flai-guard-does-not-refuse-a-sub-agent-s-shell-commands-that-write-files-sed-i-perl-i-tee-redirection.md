---
id: S-0368
type: story
nature: improvement
title: flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)
status: backlog
owner: alex
created: 2026-10-12T16:52:41Z
updated: 2026-10-12T16:52:41Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0368 flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)

## Goal

This story remediates [I-0052](../../../design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md), "flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0052 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0052 is closed with `flai issue close I-0052 --reason` saying what fixed it

## Tasks

## Notes

I-0052's cost "0m" is not a duration longer than zero, so no time_lost_per_cycle was derived from it.

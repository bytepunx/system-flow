---
id: S-0666
type: story
nature: improvement
title: This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would
status: backlog
owner: alex
created: 2026-11-10T08:34:38Z
updated: 2026-11-10T08:34:38Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0666 This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would

## Goal

This story remediates [I-0275](../../../design/issues/I-0275-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md), "This repository's guard hooks let a call through when scripts/flai.sh exits other than 2, 126, or 127, as a failed go build inside the agent container would". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0275 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0275 is closed with `flai issue close I-0275 --reason` saying what fixed it

## Tasks

## Notes

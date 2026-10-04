---
id: S-0248
type: story
nature: remediation
title: "flai check splits an item's front-matter errors on \"; \", so a message that holds one becomes two findings, the second on line 1"
status: ready
owner: alex
created: 2026-10-03T17:49:41Z
updated: 2026-10-04T21:42:42Z
transitions:
  - to: ready
    at: 2026-10-04T21:42:42Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0248 flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1

## Goal

This story remediates [I-0061](../../../design/issues/I-0061-flai-check-splits-an-item-s-front-matter-errors-on-so-a-message-that-holds-one-becomes-two-findings-the-second-on-line-1.md), "flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0061 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0061 is closed with `flai issue close I-0061 --reason` saying what fixed it

## Tasks

## Notes

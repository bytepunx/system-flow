---
id: S-0654
type: story
nature: improvement
title: A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses
status: backlog
owner: alex
created: 2026-11-10T05:45:23Z
updated: 2026-11-10T05:45:23Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0654 A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses

## Goal

This story remediates [I-0328](../../../design/issues/I-0328-a-checked-edit-of-a-work-item-writes-checks-and-restores-it-with-no-lock-so-a.md), "A checked edit of a work item writes, checks, and restores it with no lock, so a reader in between sees an edit the check refuses". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0328 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0328 is closed with `flai issue close I-0328 --reason` saying what fixed it

## Tasks

## Notes

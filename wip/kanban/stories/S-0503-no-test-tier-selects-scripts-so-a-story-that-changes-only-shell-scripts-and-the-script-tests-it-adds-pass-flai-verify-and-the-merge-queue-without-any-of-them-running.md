---
id: S-0503
type: story
nature: improvement
title: "No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running"
status: backlog
owner: alex
created: 2026-10-26T04:01:45Z
updated: 2026-10-26T04:01:45Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0503 No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running

## Goal

This story remediates [I-0199](../../../design/issues/I-0199-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md), "No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0199 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0199 is closed with `flai issue close I-0199 --reason` saying what fixed it

## Tasks

## Notes

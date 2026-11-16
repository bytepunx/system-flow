---
id: S-0800
type: story
nature: remediation
title: flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story
status: backlog
owner: alex
created: 2026-11-16T01:39:43Z
updated: 2026-11-16T01:39:43Z
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
    time_lost_per_cycle: 45m
    by: flai
    at: 2026-11-16T01:39:43Z
---
# S-0800 flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story

## Goal

This story remediates [I-0385](../../../design/issues/I-0385-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-the-branch.md), "flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story". The issue recommends this solution:

The undone check does not count a line of main's that carries no change of its own. Proposed, for the story to settle in an ADR refining ADR-0155: a line main added is undone only when the branch removes it where main added it, by comparing the hunks' positions in main's file rather than the file's lines as one multiset; or, more narrowly, a line that holds only brackets, braces, or white space is never counted alone. The revert ADR-0155 was written for (I-0124) must still be caught.

## Acceptance criteria
- [ ] The cause I-0385 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0385 is closed with `flai issue close I-0385 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0385. time_lost_per_cycle 45m: 45m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-16T01:39:37Z, 0 days before this story; under one cycle counts as one).

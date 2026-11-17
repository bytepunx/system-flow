---
id: I-0403
title: "A refusal and a command's help name the baseline's states where a project's workflow adds one: has-task says before it goes to review, flai board limit lists ready, in-progress, or review"
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-17T11:35:25Z
last_reported: 2026-11-17T11:35:25Z
updated: 2026-11-17T11:35:25Z
---

# I-0403 A refusal and a command's help name the baseline's states where a project's workflow adds one: has-task says before it goes to review, flai board limit lists ready, in-progress, or review

## Description
A refusal and a command's help name the baseline's states where a project's workflow adds one: has-task says before it goes to review, flai board limit lists ready, in-progress, or review

## Instances

### 2026-11-17T11:35:25Z
Story: S-0756.
S-0756's end-to-end test puts has-task on the move from in-progress into an added state, testing. The refusal (flai/internal/workitem/rules.go, the has-task rule) says 'a story needs at least one task before it goes to review' when the refused move is into testing; it should name the transition's target. flai board limit's help (flai/cmd/board.go) says 'Set the WIP limit of ready, in-progress, or review', though it takes any state whose limit is true, testing included. Wording only: both moves and the limit behave as the file says.

## Remediation

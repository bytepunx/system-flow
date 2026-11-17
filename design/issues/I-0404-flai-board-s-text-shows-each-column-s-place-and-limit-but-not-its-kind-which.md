---
id: I-0404
title: flai board's text shows each column's place and limit but not its kind, which only --json and board.get give
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-17T11:35:28Z
last_reported: 2026-11-17T11:35:28Z
updated: 2026-11-17T11:35:28Z
---

# I-0404 flai board's text shows each column's place and limit but not its kind, which only --json and board.get give

## Description
flai board's text shows each column's place and limit but not its kind, which only --json and board.get give

## Instances

### 2026-11-17T11:35:28Z
Story: S-0756.
S-0756's end-to-end test checks that the board shows an added state, testing, with its kind and WIP limit. flai board --json and board.get carry kind active; the plain flai board prints 'testing (0/2)' in its place with no kind, as it does for every column. Not particular to an added state, but an operator reading the text cannot tell an active column from a waiting one.

## Remediation

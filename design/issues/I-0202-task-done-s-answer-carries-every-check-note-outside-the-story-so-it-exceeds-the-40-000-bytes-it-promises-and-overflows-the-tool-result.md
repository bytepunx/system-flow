---
id: I-0202
title: task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result
class: efficiency
status: closed
count: 3
cost: 2m
first_reported: 2026-10-23T21:50:41Z
last_reported: 2026-10-24T00:18:48Z
updated: 2026-10-24T00:18:48Z
---

# I-0202 task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result

## Description
task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result

## Instances

### 2026-10-23T21:50:41Z
Story: S-0491.
S-0491: each of four task_done calls answered about 90,500 characters, over Claude Code's tool output limit, so the result went to a file I had to read with jq. The bulk is check.notes: 399 notes outside the story (382 doc.updated, 17 components.examples). The 40,000-byte fit trims the inbox's changes and thread entries but not the check's notes. Fix: cap or count the notes in task_done's answer, as the inbox's changes are counted in changes_omitted.

### 2026-10-24T00:04:39Z
Story: S-0490.
S-0490's first task_done (T-2099) answered 92,534 characters, 400 check notes outside the story, so Claude Code saved it to a file and the agent read it back with jq.

### 2026-10-24T00:18:48Z
Story: S-0486.
task_done T-2074 over MCP answered 94,567 characters, with 399 check notes outside the story, and overflowed the tool result. The commit, move, and check had gone through, which reading the saved file with jq showed. Later tasks were closed with flai task done --json on the host instead.

## Remediation

Story S-0492 remediates this issue, created from it at 2026-10-23T23:38:52Z.
Closed 2026-10-23T23:59:14Z: S-0492: flai task done's check step lists every finding in the story but only the first ten notes outside it, counting the rest in check.notes_omitted (ADR-0198), so a check with hundreds of notes no longer takes the answer over 40,000 bytes; TestNewCheckListsAtMostMaxNotesAndCountsTheRest reproduces S-0491's 399 notes.

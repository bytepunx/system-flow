---
id: I-0202
title: task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result
class: efficiency
status: open
count: 1
cost: 2m
first_reported: 2026-10-23T21:50:41Z
last_reported: 2026-10-23T21:50:41Z
updated: 2026-10-23T21:50:41Z
---

# I-0202 task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result

## Description
task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result

## Instances

### 2026-10-23T21:50:41Z
Story: S-0491.
S-0491: each of four task_done calls answered about 90,500 characters, over Claude Code's tool output limit, so the result went to a file I had to read with jq. The bulk is check.notes: 399 notes outside the story (382 doc.updated, 17 components.examples). The 40,000-byte fit trims the inbox's changes and thread entries but not the check's notes. Fix: cap or count the notes in task_done's answer, as the inbox's changes are counted in changes_omitted.

## Remediation

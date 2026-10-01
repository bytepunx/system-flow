---
id: I-0043
title: A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint
class: defect
status: open
count: 3
cost: 3m
first_reported: 2026-09-24T09:19:30Z
last_reported: 2026-10-01T07:46:11Z
updated: 2026-10-01T07:46:11Z
---

# I-0043 A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint

## Description
A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint

## Instances

### 2026-09-24T09:19:30Z
S-0116, 2026-09-24: thread_reply then thread_resolve on TH-0010 at 09:00:12Z wrote two '### 2026-09-24T09:00:12Z system-flow' headings; make lint-md failed MD024 in the main checkout. Folded the Resolved line into the reply's entry by hand. flai should merge or disambiguate same-second entries by one author.

### 2026-09-28T22:05:36Z
S-0134, 2026-09-28: thread_reply then thread_resolve on TH-0028 at 22:04:56Z wrote two same-second headings by agent-S-0134; MD024 in the main checkout. Folded the Resolved line into the reply's entry by hand.

### 2026-10-01T07:46:11Z
S-0173, 2026-10-01: make lint-md on main fails MD024 at wip/threads/TH-0035-...md:59, two same-second agent-S-0153 headings; left as found, not this story's file.

## Remediation

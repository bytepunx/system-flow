---
id: I-0043
title: A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint
class: defect
status: open
count: 2
cost: 4m
first_reported: 2026-09-24T09:19:30Z
last_reported: 2026-09-28T22:05:36Z
updated: 2026-09-28T22:05:36Z
---

# I-0043 A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint

## Description
A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint

## Instances

### 2026-09-24T09:19:30Z
S-0116, 2026-09-24: thread_reply then thread_resolve on TH-0010 at 09:00:12Z wrote two '### 2026-09-24T09:00:12Z system-flow' headings; make lint-md failed MD024 in the main checkout. Folded the Resolved line into the reply's entry by hand. flai should merge or disambiguate same-second entries by one author.

### 2026-09-28T22:05:36Z
S-0134, 2026-09-28: thread_reply then thread_resolve on TH-0028 at 22:04:56Z wrote two same-second headings by agent-S-0134; MD024 in the main checkout. Folded the Resolved line into the reply's entry by hand.

## Remediation

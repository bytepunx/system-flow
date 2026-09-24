---
id: I-0043
title: A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-09-24T09:19:30Z
last_reported: 2026-09-24T09:19:30Z
updated: 2026-09-24T09:19:30Z
---

# I-0043 A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint

## Description
A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint

## Instances

### 2026-09-24T09:19:30Z
S-0116, 2026-09-24: thread_reply then thread_resolve on TH-0010 at 09:00:12Z wrote two '### 2026-09-24T09:00:12Z system-flow' headings; make lint-md failed MD024 in the main checkout. Folded the Resolved line into the reply's entry by hand. flai should merge or disambiguate same-second entries by one author.

## Remediation

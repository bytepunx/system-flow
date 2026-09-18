---
id: I-0011
title: flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI
class: defect
status: open
count: 2
cost: 10m
first_reported: 2026-09-17T06:15:37Z
last_reported: 2026-09-18T17:11:45Z
updated: 2026-09-18T17:11:45Z
---

# I-0011 flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI

## Description
flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI

## Instances

### 2026-09-17T06:15:37Z
Three stream log calls in one close-out chain landed in the same second; the narrative got three identical timestamp headings and system-flow-check failed on MD024. Fixed by making LogStream append under an existing same-second heading.

### 2026-09-18T17:11:45Z
Second occurrence, in a different command: two flai issue bump calls in one second wrote duplicate instance headings in I-0016 and failed the markdown lint in CI. The narrative fix did not cover issue instances. I also pushed that chore commit without running scripts/lint-md.sh.

## Remediation

---
id: I-011
title: flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-17T06:15:37Z
last_reported: 2026-09-17T06:15:37Z
updated: 2026-09-17T06:15:37Z
---

# I-011 flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI

## Description
flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI

## Instances

### 2026-09-17T06:15:37Z
Three stream log calls in one close-out chain landed in the same second; the narrative got three identical timestamp headings and system-flow-check failed on MD024. Fixed by making LogStream append under an existing same-second heading.

## Remediation

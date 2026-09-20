---
id: I-0027
title: Work items written in the main checkout reach CI without a markdown lint
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-20T13:03:19Z
last_reported: 2026-09-20T13:03:19Z
updated: 2026-09-20T13:03:19Z
---

# I-0027 Work items written in the main checkout reach CI without a markdown lint

## Description
Work items written in the main checkout reach CI without a markdown lint

## Instances

### 2026-09-20T13:03:19Z
2026-09-20: second time today. A task body of S-0077 had 'GIT_ or PROJECT_DIR', which markdownlint reads as emphasis (MD037); the system-flow check failed on main for the S-0076 and S-0077 acceptance commits. Earlier the same day E-0007's bullets failed MD004. wip/ is written in the main checkout (ADR-0019), so the lint a story runs in its worktree never sees it, and flai check does not lint markdown. Remedy to consider: lint wip/ in the main checkout before a story goes to review, or have flai task new refuse a body that the project's markdownlint rules reject.

## Remediation

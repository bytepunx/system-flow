---
id: I-0027
title: Work items written in the main checkout reach CI without a markdown lint
class: defect
status: open
count: 2
cost: 8m
first_reported: 2026-09-20T13:03:19Z
last_reported: 2026-09-26T07:22:40Z
updated: 2026-09-26T07:22:40Z
---

# I-0027 Work items written in the main checkout reach CI without a markdown lint

## Description
Work items written in the main checkout reach CI without a markdown lint

## Instances

### 2026-09-20T13:03:19Z
2026-09-20: second time today. A task body of S-0077 had `GIT_` beside `PROJECT_DIR` unquoted, which markdownlint reads as emphasis (MD037); the system-flow check failed on main for the S-0076 and S-0077 acceptance commits. Earlier the same day E-0007's bullets failed MD004. wip/ is written in the main checkout (ADR-0019), so the lint a story runs in its worktree never sees it, and flai check does not lint markdown. Remedy to consider: lint wip/ in the main checkout before a story goes to review, or have flai task new refuse a body that the project's markdownlint rules reject.

### 2026-09-26T07:22:40Z
S-0122's acceptance (32f53ab) committed wip/archive/agents/S-0122.md and its archived story with a double blank line (MD012), and TH-0012 with emphasis as a heading (MD036); the system-flow check run on main (36225691919) failed at Lint markdown after flai push --pending on 2026-09-26. The archive is not the agent's to edit, so main stays red until the operator fixes or allows it.

## Remediation

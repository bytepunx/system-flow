---
id: I-0027
title: Work items written in the main checkout reach CI without a markdown lint
class: defect
status: open
count: 10
cost: 6m
first_reported: 2026-09-20T13:03:19Z
last_reported: 2026-09-29T23:48:59Z
updated: 2026-09-29T23:48:59Z
---

# I-0027 Work items written in the main checkout reach CI without a markdown lint

## Description
Work items written in the main checkout reach CI without a markdown lint

## Instances

### 2026-09-20T13:03:19Z
2026-09-20: second time today. A task body of S-0077 had `GIT_` beside `PROJECT_DIR` unquoted, which markdownlint reads as emphasis (MD037); the system-flow check failed on main for the S-0076 and S-0077 acceptance commits. Earlier the same day E-0007's bullets failed MD004. wip/ is written in the main checkout (ADR-0019), so the lint a story runs in its worktree never sees it, and flai check does not lint markdown. Remedy to consider: lint wip/ in the main checkout before a story goes to review, or have flai task new refuse a body that the project's markdownlint rules reject.

### 2026-09-26T07:22:40Z
S-0122's acceptance (32f53ab) committed wip/archive/agents/S-0122.md and its archived story with a double blank line (MD012), and TH-0012 with emphasis as a heading (MD036); the system-flow check run on main (36225691919) failed at Lint markdown after flai push --pending on 2026-09-26. The archive is not the agent's to edit, so main stays red until the operator fixes or allows it.

### 2026-09-26T07:34:08Z
S-0126: make lint-md in the main checkout found six findings in files no agent linted before they were written there: S-0124's story (MD009, MD012), S-0127's story heading (MD026), and TH-0012 (MD036, already in TH-0017). None in S-0126's files; left to their owners.

### 2026-09-27T04:08:25Z
S-0134 (2026-09-27): make smoke failed on six lint errors in wip files main already carried: the S-0122 archive (TH-0017's fix was never made) and now S-0124's archived narrative (three leftover lines, MD026/MD029). Fixed on main in 88e345f under TH-0017's answer A.

### 2026-09-29T07:21:07Z
Found by S-0154's smoke run: wip/threads/TH-0035 (written in the main checkout during S-0153) fails markdownlint MD024 at line 59 (two entries headed with the same timestamp and author), so scripts/lint-md.sh and make smoke fail on main. Left for the thread's owner; S-0154 does not touch it.

### 2026-09-29T20:27:08Z
S-0162's smoke tier stopped at markdown lint: wip/threads/TH-0035 on main has MD024 (duplicate heading from two thread entries logged in the same second by one agent). A thread written through flai in the main checkout reached main without a lint.

### 2026-09-29T20:59:54Z
S-0163's smoke tier: the markdown lint still fails on wip/threads/TH-0035 on main (MD024, two entries headed with the same second and author), the third story in a day to meet it. The tier also stops a step earlier, at flai check --strict, on four warnings that are the operator's to clear: E-0003, E-0010, and E-0012 are done and not archived, and TH-0032 is answered on an archived story. Nothing of S-0163 is in either.

### 2026-09-29T21:21:56Z
S-0164's smoke tier: the markdown lint still fails on wip/threads/TH-0035 on main (MD024), the fourth story in a day to meet it. The tier stops a step earlier, at flai check --strict, on five warnings that are the operator's to clear: E-0003, E-0010, E-0011, and E-0012 are done and not archived, and TH-0032 is answered on an archived story. Nothing of S-0164 is in either.

### 2026-09-29T22:51:11Z
S-0166: the story's title, written in the main checkout, ended in a period, so its heading failed MD026 in scripts/lint-md.sh; retitled without it. TH-0035 fails MD024 (two log headings alike), left to its owner.

### 2026-09-29T23:48:59Z
S-0167's smoke tier: markdown lint fails on wip/threads/TH-0035 (MD024, two log headings alike), a thread written in the main checkout and committed to main unlinted.

## Remediation

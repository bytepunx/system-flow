---
id: I-0117
title: A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix
class: blocker
status: open
count: 1
cost: 30m
first_reported: 2026-10-08T04:12:50Z
last_reported: 2026-10-08T04:12:50Z
updated: 2026-10-08T04:12:50Z
---

# I-0117 A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

## Description
A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

## Instances

### 2026-10-08T04:12:50Z
Story: S-0324.
S-0324 makes flai's MD034 report a bare `www.` literal. `TestRepositoryLintsClean` in flai/internal/mdlint lints every markdown file the story branch holds, `wip/` included. A story branch holds `wip/` only as main last committed it. The close-out's integration tier failed twice on bare `www.` lines in that copy. The first time it was four lines (TH-0354), the second two lines of TH-0355's title, which S-0315's acceptance had just committed. Each time the lines were fixed in the main checkout within minutes, but the story could not go on until a later commit on main, such as an acceptance, took the fix. The installed flai (1.38.1) lacks the rule until this fix is published, so agents keep writing such lines. markdownlint-cli2 reports the same lines, so `make lint-md` on main fails on them too.

## Remediation

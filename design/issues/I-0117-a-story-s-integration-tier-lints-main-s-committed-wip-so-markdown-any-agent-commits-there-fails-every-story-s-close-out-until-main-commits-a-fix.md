---
id: I-0117
title: A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix
class: blocker
status: open
count: 3
cost: 17m
first_reported: 2026-10-08T04:12:50Z
last_reported: 2026-10-08T05:31:34Z
updated: 2026-10-08T05:31:34Z
---

# I-0117 A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

## Description
A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

## Instances

### 2026-10-08T04:12:50Z
Story: S-0324.
S-0324 makes flai's MD034 report a bare `www.` literal. `TestRepositoryLintsClean` in flai/internal/mdlint lints every markdown file the story branch holds, `wip/` included. A story branch holds `wip/` only as main last committed it. The close-out's integration tier failed twice on bare `www.` lines in that copy. The first time it was four lines (TH-0354), the second two lines of TH-0355's title, which S-0315's acceptance had just committed. Each time the lines were fixed in the main checkout within minutes, but the story could not go on until a later commit on main, such as an acceptance, took the fix. The installed flai (1.38.1) lacks the rule until this fix is published, so agents keep writing such lines. markdownlint-cli2 reports the same lines, so `make lint-md` on main fails on them too.

### 2026-10-08T05:00:50Z
Story: S-0318.
S-0318's close-out failed its integration tier on TestRepositoryLintsClean: `wip/agents/orchestrator.md` line 1529, the orchestrator's publish summary for S-0324, holds a bare `www.` literal (MD034), committed on main by cad49eee (publish flai 1.39.5). The tier's tail did not name the test (I-0113); the outside note on the same line pointed to it. Fixed in the main checkout by wrapping it in a code span, but S-0318 cannot pass until a main commit takes the fix, and every other story's close-out fails meanwhile.

### 2026-10-08T05:31:34Z
Story: S-0318.
After main took the orchestrator.md fix, S-0318's close-out, now on a flai with the `www.` rule, stopped at check on a bare `www.` in TH-0360's title and its mirror in wip/agents/S-0318.md. The installed flai 1.39.3 behind the MCP server wrote that title without complaint. Quoted both by hand in the main checkout.

## Remediation

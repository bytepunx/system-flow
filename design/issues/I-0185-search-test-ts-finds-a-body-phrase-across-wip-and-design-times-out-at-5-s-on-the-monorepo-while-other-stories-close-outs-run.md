---
id: I-0185
title: search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-10-21T12:58:18Z
last_reported: 2026-10-21T12:58:18Z
updated: 2026-10-26T04:36:02Z
---

# I-0185 search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run

## Description
search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run

## Instances

### 2026-10-21T12:58:18Z
Story: S-0461.
S-0461's close-out stopped at the flaiover tier: src/lib/server/search.test.ts > adrs on the monorepo > finds a body phrase across wip and design timed out in 5000ms, the other 2000 tests passing, while four other stories worked on the host. S-0461 changes no server or search code; the test searches the whole repository with vitest's default timeout.

## Remediation

Story S-0501 remediates this issue, created from it at 2026-10-26T04:01:35Z.
Closed 2026-10-26T04:36:02Z: S-0501. The cause: each monorepo test in flaiover/src/lib/server/search.test.ts runs a one-shot flai hostapi that builds the whole repository's search index, measured at 1.8 to 2.0 s wall and 1.9 s of CPU for 4,845 markdown files, and 4.8 to 6.1 s each with 24 at once on a 24-core host, over vitest's 5 s default; flai serve keeps the index, the test cannot. The fix gives the 'adrs on the monorepo' block a 30 s timeout with the figures in a comment, as writes.test.ts did in S-0116. No test reproduces it: the timeout depends on host load, which a unit test cannot set. Whether the one-shot build should be made cheaper is asked on TH-0654.

---
id: I-0185
title: search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-21T12:58:18Z
last_reported: 2026-10-21T12:58:18Z
updated: 2026-10-26T04:01:35Z
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

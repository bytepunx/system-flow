---
id: I-0211
title: ADR-0202's status rejected is not an ADR status flaiover's search test allows, so the flaiover tier fails on main for every story that selects it
class: blocker
status: closed
count: 1
cost: 15m
first_reported: 2026-10-26T05:43:18Z
last_reported: 2026-10-26T05:43:18Z
updated: 2026-10-26T05:43:36Z
---

# I-0211 ADR-0202's status rejected is not an ADR status flaiover's search test allows, so the flaiover tier fails on main for every story that selects it

## Description
ADR-0202's status rejected is not an ADR status flaiover's search test allows, so the flaiover tier fails on main for every story that selects it

## Instances

### 2026-10-26T05:43:18Z
Story: S-0503.
S-0503's close-out stopped at the flaiover tier, which the story's script changes select (I-0199's fix): `adrs on the monorepo > lists ADRs with the supersession chain` in flaiover/src/lib/server/search.test.ts:81 expects every ADR's status in proposed, accepted, superseded, deprecated, the set design/adrs/0000-template.md documents. Commit a554b18ff on main (2026-10-25, "docs: Rejected in favor of 0203") set ADR-0202's status to rejected. Nothing in flai or flaiover validates the status otherwise, so flai check passes; but the flaiover tier now fails on main for any story, or merge queue batch, whose paths select it.

## Remediation
Closed 2026-10-26T05:43:36Z: Duplicate of I-0208, recorded on S-0501's branch, which adds rejected to search.test.ts's statuses; S-0503's close-out runs again once S-0501 is accepted.

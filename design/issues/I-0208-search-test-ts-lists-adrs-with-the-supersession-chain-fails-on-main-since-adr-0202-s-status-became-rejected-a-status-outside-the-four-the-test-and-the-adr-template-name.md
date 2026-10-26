---
id: I-0208
title: search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-26T04:35:14Z
last_reported: 2026-10-26T04:35:14Z
updated: 2026-10-26T16:03:03Z
---

# I-0208 search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name

## Description
search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name

## Instances

### 2026-10-26T04:35:14Z
Story: S-0501.
S-0501's flai test on search.test.ts failed at 'adrs on the monorepo > lists ADRs with the supersession chain': expected [ { ADR-0202 } ] to deeply equal []. Commit a554b18ff on main ("docs: Rejected in favor of 0203") set ADR-0202's status to rejected; the test allowed only proposed, accepted, superseded, and deprecated, as template/root/design/adrs/0000-template.md's status comment lists them, and flai check accepts rejected without a finding. Every story whose change selects that test fails its vitest tier until it is fixed. S-0501 adds rejected to the test's list. Left open: whether rejected is an ADR status the template and documentation-standard.md should name, and flai check should validate the status at all.

## Remediation

Story S-0512 remediates this issue, created from it at 2026-10-26T16:03:03Z.

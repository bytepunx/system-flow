---
id: I-0304
title: flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out
class: efficiency
status: open
count: 1
cost: 8m
first_reported: 2026-11-07T11:41:27Z
last_reported: 2026-11-07T11:41:27Z
updated: 2026-11-07T11:41:27Z
---

# I-0304 flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out

## Description
flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out

## Instances

### 2026-11-07T11:41:27Z
Story: S-0587.
S-0587 changed the gate of settings.agent in flai/internal/hostapi/writes.go so that it is per project unless everywhere is true. `flai test` on the hostapi and flaiover paths passed. The close-out's flaiover tier then failed at flaiover/src/lib/server/writes.test.ts:578, which runs the built flai and still expected a 403 without --all-projects. Its vitest `related` run picks tests by the TypeScript import graph, so no change under flai/internal/hostapi selects it. It cost one close-out round: 2m29s, plus the fix and a new request.

## Remediation

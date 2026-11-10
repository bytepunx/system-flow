---
id: I-0269
title: repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it
class: defect
status: open
count: 3
cost: 6m
first_reported: 2026-11-07T04:07:16Z
last_reported: 2026-11-08T02:50:04Z
updated: 2026-11-10T08:54:28Z
---

# I-0269 repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it

## Description
repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it

## Instances

### 2026-11-07T04:07:16Z
Story: S-0557.
S-0557's close-out stopped at the flaiover tier on one test: flaiover/src/lib/server/repo.test.ts › Repo on the monorepo › reads every item with a consistent status timed out at 5000ms. The other 2351 tests passed. The test reads every work item of the monorepo, archive included, with vitest's default 5s timeout, while five stories' agents ran on the host. S-0557 changes only ProjectSwitcher.svelte, its unit test, an e2e file, and I-0180, none of which the test reads beyond one issue file. The close-out was asked for again.

### 2026-11-07T11:04:52Z
Story: S-0586.
S-0586's second close-out (job S-0586-20261107T110011Z, at bc07452) stopped at the flaiover tier: repo.test.ts › Repo on the monorepo › reads every item with a consistent status timed out at 5000ms. S-0586 does not change repo.test.ts's subject, and the same tier passed at b509acf half an hour earlier.

### 2026-11-08T02:50:04Z
Story: S-0607.
S-0607, T-2656: the flaiover-flai tier failed once on `flaiover/src/lib/server/repo.test.ts:94` ("Repo on the monorepo reads every item with a consistent status", STACK_TRACE_ERROR), on a host at load about 90 on 24 cores. The change, in `flai/internal/hostapi/reads.go`, does not touch what that test reads. The same tier passed in the next run and in the close-out.

## Remediation

Story S-0695 remediates this issue, created from it at 2026-11-10T08:54:28Z.

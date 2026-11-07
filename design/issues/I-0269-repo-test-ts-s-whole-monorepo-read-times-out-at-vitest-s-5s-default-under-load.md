---
id: I-0269
title: repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it
class: defect
status: open
count: 1
cost: 6m
first_reported: 2026-11-07T04:07:16Z
last_reported: 2026-11-07T04:07:16Z
updated: 2026-11-07T04:07:16Z
---

# I-0269 repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it

## Description
repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it

## Instances

### 2026-11-07T04:07:16Z
Story: S-0557.
S-0557's close-out stopped at the flaiover tier on one test: flaiover/src/lib/server/repo.test.ts › Repo on the monorepo › reads every item with a consistent status timed out at 5000ms. The other 2351 tests passed. The test reads every work item of the monorepo, archive included, with vitest's default 5s timeout, while five stories' agents ran on the host. S-0557 changes only ProjectSwitcher.svelte, its unit test, an e2e file, and I-0180, none of which the test reads beyond one issue file. The close-out was asked for again.

## Remediation

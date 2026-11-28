---
id: T-3650
type: task
nature: remediation
title: Make the dashboard's budget banner say the login failed and how to renew it
status: backlog
parent: S-0855
owner: alex
created: 2026-11-28T00:00:02Z
updated: 2026-11-28T00:00:02Z
transitions: []
stream: S-0855
tags: [flaiover]
touches: [flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/BudgetBanner.svelte, flaiover/src/lib/components/BudgetBanner.svelte.test.ts]
after: [T-3647]
---
# T-3650 Make the dashboard's budget banner say the login failed and how to renew it

## Work

Waits for T-3647, since it reads the kind that task keeps on the budget state, as the host API's status gives it.

- Add the kind to `BudgetHold` in `flaiover/src/lib/hostflai.svelte.ts`.
- In `BudgetBanner.svelte`, say that a login hold's login failed, and name the command that renews it, not "out of budget". A usage-limit hold keeps its words.
- Add cases to `BudgetBanner.svelte.test.ts` for a login hold, with a probe and without one.

This task runs together with T-3649: the two share no path.

## Done when

- The flaiover tests and lint pass on the changed files.
- A usage-limit hold's banner reads as it did before.

## Notes

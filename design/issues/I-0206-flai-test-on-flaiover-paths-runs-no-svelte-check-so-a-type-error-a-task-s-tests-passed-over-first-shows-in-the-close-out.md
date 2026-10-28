---
id: I-0206
title: flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out
class: efficiency
status: closed
count: 1
cost: 10m
first_reported: 2026-10-24T22:32:29Z
last_reported: 2026-10-24T22:32:29Z
updated: 2026-10-28T21:03:10Z
---

# I-0206 flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out

## Description
flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out

## Instances

### 2026-10-24T22:32:29Z
Story: S-0434.
S-0434's T-1845 passed flai test on flaiover/src/lib/viz and flaiover/src/routes/charts (flaiover-lint, vitest), but its normalise wrote an object literal whose spread overwrote every key, which svelte-check reports as TS2783. Only the close-out's flaiover tier runs svelte-check, so the error stopped the close-out after a 6-minute go-test tier, and the fix needed a second close-out. A svelte-check scoped to the changed files, or the flaiover tier when a .ts or .svelte file under flaiover changes, would catch it at the task.

## Remediation

Story S-0530 remediates this issue, created from it at 2026-10-28T12:54:31Z.
Closed 2026-10-28T21:03:10Z: S-0530: system-flow.yaml's tests gain a svelte-check tier, outside --all and after vitest, that runs scripts/flaiover-check.sh (pnpm check --incremental, svelte-kit sync then svelte-check on the whole flaiover project) when a changed flaiover/src .ts, .js, or .svelte file, a root flaiover .ts or .js file, flaiover/tsconfig.json, flaiover/package.json, or the script itself selects it, so flai test finds a type error such as TS2783 at the task, in 8 to 10 s warm, before the close-out. The flaiover tier covers it, so --all and the close-out run svelte-check once. TestRepositoryManifestTypeChecksChangedFlaioverFiles in flai/internal/verify/select_test.go fails without the tier.

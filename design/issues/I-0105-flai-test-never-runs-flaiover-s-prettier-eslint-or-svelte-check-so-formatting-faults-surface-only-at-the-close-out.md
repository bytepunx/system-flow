---
id: I-0105
title: flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-10-07T07:23:32Z
last_reported: 2026-10-07T07:23:32Z
updated: 2026-10-07T18:59:50Z
---

# I-0105 flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out

## Description
flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out

## Instances

### 2026-10-07T07:23:32Z
Story: S-0212.
The manifest's `flaiover` tier (`scripts/flaiover-test.sh`: prettier, eslint, svelte-check) is `all_only`, so `flai test` on changed flaiover paths runs vitest alone. Four task sub-agents of S-0212 formatted by hand and still left prettier findings in three files, which stopped the close-out at the flaiover step and cost a second verifier run. A per-file prettier --check tier selected by `flaiover/**` paths would catch them at each task.

## Remediation

Story S-0319 remediates this issue, created from it at 2026-10-07T18:59:50Z.

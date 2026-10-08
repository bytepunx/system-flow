---
id: I-0105
title: flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out
class: efficiency
status: closed
count: 2
cost: 8m
first_reported: 2026-10-07T07:23:32Z
last_reported: 2026-10-07T20:41:40Z
updated: 2026-10-08T05:55:03Z
---

# I-0105 flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out

## Description
flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out

## Instances

### 2026-10-07T07:23:32Z
Story: S-0212.
The manifest's `flaiover` tier (`scripts/flaiover-test.sh`: prettier, eslint, svelte-check) is `all_only`, so `flai test` on changed flaiover paths runs vitest alone. Four task sub-agents of S-0212 formatted by hand and still left prettier findings in three files, which stopped the close-out at the flaiover step and cost a second verifier run. A per-file prettier --check tier selected by `flaiover/**` paths would catch them at each task.

### 2026-10-07T20:41:40Z
Story: S-0329.
S-0329's close-out stopped at the flaiover tier on prettier in four files (HostAgentNotice.svelte, SettingsPanel.svelte, StrategicAgentPanel.svelte, charts.test.ts) that three task sub-agents had changed and tested with flai test, which ran only vitest for them.

## Remediation

Story S-0319 remediates this issue, created from it at 2026-10-07T18:59:50Z.
Closed 2026-10-08T05:55:03Z: S-0319: system-flow.yaml's tests gain a `flaiover-lint` tier, outside `--all` and before vitest, that runs `scripts/flaiover-lint.sh` (prettier `--check --ignore-unknown` and eslint `--no-warn-ignored`) on the changed files under `flaiover/`, so `flai test` finds a formatting or lint fault in a flaiover file before the close-out. svelte-check takes no file list and stays in the `--all` `flaiover` tier. `TestRepositoryManifestLintsChangedFlaioverFiles` in `flai/internal/verify/select_test.go` fails without the tier.

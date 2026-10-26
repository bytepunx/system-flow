---
id: TH-0661
title: "Plan for S-0503: shellcheck tier, script tests selected by the tiers that run them, and a cost of delay input to set"
anchor:
  path: wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md
  item: S-0503
status: resolved
participants: [planner-S-0503, orchestrator]
created: 2026-10-26T05:25:15Z
updated: 2026-10-26T05:25:33Z
---

# TH-0661 Plan for S-0503: shellcheck tier, script tests selected by the tiers that run them, and a cost of delay input to set

On wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md.

## Entries

### 2026-10-26T05:25:15Z planner-S-0503
The plan is written: five tasks in five layers, touches on 17 files, forecast 1h. One input is missing, the cost of delay, asked below.

## Solution proposed from I-0199's instance

S-0481 changed `scripts/flaiover-install.sh` and added `scripts/flaiover-install-test.sh`; no tier selected either, so `flai verify` and the merge queue's batch ran neither shellcheck nor the test. Two changes close the gap:

1. A `shellcheck` tier in `system-flow.yaml`, not `all_only`, selecting `scripts/**/*.sh`, running a new `scripts/shellcheck.sh` (`shellcheck -x -P SCRIPTDIR`) on the changed scripts, so every script change is linted in `flai test` and `flai verify`.
2. The `all_only` tiers that run scripts select them: `flaiover` adds `scripts/flaiover-*.sh`, `template` adds `scripts/template-test.sh`, `smoke` adds `scripts/**`. The batch runs the tiers its paths select (ADR-0174), so a scripts-only story then runs its script tests there and in `flai verify`.

A Go test in `flai/internal/verify/select_test.go` reproduces the instance against the real manifest and fails for any `scripts/*-test.sh` no tier runs.

## Tasks and layers

| Layer | Task | After |
|-------|------|-------|
| 1 | T-2179 `scripts/install-tools.sh` installs a pinned ShellCheck v0.10.0 into `bin/`; `design/tech/shellcheck.md` | none |
| 2 | T-2180 `scripts/shellcheck.sh`; fix the seven findings in four scripts; `make shellcheck`; `flai-test.sh` runs it | T-2179 |
| 3 | T-2181 the `shellcheck` tier and the widened paths in `system-flow.yaml`; docs in tooling.md, scripts/README.md, project-manifest.md | T-2180 |
| 4 | T-2182 `TestRepositoryManifestRunsScriptTestsForScriptChanges` | T-2181 |
| 5 | T-2183 `flai issue close I-0199`, tick the criteria | T-2181, T-2182 |

## Assumptions

- ShellCheck pinned at v0.10.0, the host's version; installed by release download with checksum, since it is neither a Go module nor an npm package.
- `scripts/**/*.sh` selects `scripts/flai.sh` as `flai/**/*.go` selects `flai/main.go`; T-2181 checks and adds `scripts/*.sh` if not.
- The template's `template/root/scripts/*.sh` and its manifest stay out of scope: its `test` tier selects `**`.
- Forecast raised from flai's 29m to 1h: five sequential layers, a download, lint fixes in four scripts, a test shown failing and passing.

## Question: cost of delay inputs

S-0503 has no cost of delay inputs and no epic to take a share of, so `flai cod` cannot work out a value. Recommended: `time_lost_per_cycle: 5m` (value 12.5, by this project's scale, as for the other remediation stories of its size). The gap costs each scripts-only story a regression that reaches main unseen, rare so far (one instance in S-0481), but E-0029's container stories will add scripts. Alternatives: `3m` (7.5) if you weigh its rarity, `10m` (25) if a shell regression on main is the concern. Set it with `flai edit S-0503 --time-lost-per-cycle 5m` or `item_edit` `cost_of_delay`, and I write the value.

### 2026-10-26T05:25:29Z orchestrator
Approved, and the input is set.

- **Set on S-0503:** `time_lost_per_cycle: 5m`, the figure you recommend. I set the input only. Write the value with `flai cod`. Nothing on the thread or the story argues for 3m or 10m: there is one instance so far, and E-0029's scripts are not yet drafted into work.
- **The remedy fits.** I-0199's instance was S-0481's scripts-only change, which no tier ran. A `shellcheck` tier selecting `scripts/**/*.sh` closes it. So does having `flaiover`, `template` and `smoke` select the scripts they run, which the queue's batch runs by ADR-0174. T-2182's manifest test reproduces the instance.
- **The tasks fit.** There are five tasks in five sequential layers. T-2179 and T-2180 share `Makefile` and `scripts/README.md` but run one after the other. There are 17 file touches and no folder. Pinning ShellCheck in `design/tech` follows the documentation convention.
- **The forecast fits.** It is 1h, raised from flai's 29m with the reason stated.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T05:25:33Z orchestrator
Resolved: Plan approved under plan_backlog_stories, and S-0503's input set: time_lost_per_cycle 5m, as the planner recommended.

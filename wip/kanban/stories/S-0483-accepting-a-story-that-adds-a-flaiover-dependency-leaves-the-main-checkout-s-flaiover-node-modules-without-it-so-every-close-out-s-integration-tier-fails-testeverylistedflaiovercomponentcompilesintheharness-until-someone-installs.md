---
id: S-0483
type: story
nature: remediation
title: Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs
status: in-progress
owner: alex
created: 2026-10-21T12:00:17Z
updated: 2026-10-21T13:28:19Z
transitions:
  - to: ready
    at: 2026-10-21T13:04:10Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T13:04:17Z
    by: agent-S-0483
tags: [flai, cli]
topics: [release, testing]
touches: [design/adrs/README.md, design/system/project-manifest.md, design/system/flai-cli.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, system-flow.yaml, docs/operators/settings.md, flai/cmd/accept.go, flai/cmd/accept_install.go, flai/cmd/accept_install_test.go, docs/users/flai.md, design/system/devex.md, design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md, design/issues/summary.md, design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md, flai/internal/check/check.go, flai/internal/check/install_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/shared.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2161
  priced_by: harness
  turns:
    - day: 2026-10-21
      hand_edits: 2
      work: 48
  models:
    - model: claude-opus-5-5
      input: 248
      output: 70021
      cache_read: 12574280
      cache_write: 633793
      cost: 7.7524
  strategic:
    - kind: orchestrator
      seconds: 52
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 170
          cache_read: 3203296
          cache_write: 6002
          cost: 0.7086
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-21T12:00:17Z
  value: 37.5
  by: planner-S-0483
  at: 2026-10-21T12:12:36Z
forecast:
  duration: 1h
  delivery: 2026-10-22T07:23:00Z
  basis: "Its own forecast of 1h; 39th in the pull order with an in-progress limit of 5, behind S-0461, S-0471, S-0473, S-0479, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432, S-0433, S-0434, S-0435, S-0466, S-0467, S-0468, S-0469, S-0470, S-0472, S-0476, S-0477, S-0478 and S-0481."
  by: flai
  at: 2026-10-21T12:51:45Z
finalized:
  by: alex
  at: 2026-10-21T12:03:49Z
---
# S-0483 Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs

## Goal

This story remediates [I-0179](../../../design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md), "Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0179 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0179 is closed with `flai issue close I-0179 --reason` saying what fixed it

## Tasks
- T-1998 An ADR decides that flai accept runs a sub-project's install command in the main checkout when the accepted change touched its lock file
- T-2000 The manifest takes a sub-project's install command, flai check validates it, and flaiover's names scripts/flaiover-install.sh --if-needed
- T-2003 flai accept runs the install command of each sub-project whose lock file the accepted change touched, in the main checkout, with a test that reproduces I-0179
- T-2004 The user and developer docs say that an acceptance runs a sub-project's install when it changed the lock file
- T-2006 Verify an acceptance that changes flaiover's lock file leaves the main checkout's node_modules current, and close I-0179

## Notes

Cost of delay inputs set by flai from I-0179. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T08:44:04Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed fix: `flai accept` runs a sub-project's `install` command, a new manifest key, in the main checkout when the accepted change touched that sub-project's lock file or `package.json`. flaiover's is `scripts/flaiover-install.sh --if-needed`, which uses the pnpm store `scripts/env.sh` sets (I-0177). S-0482 fixes the story worktree's side: its integration tier installs its own `node_modules`. This story keeps the main checkout current.

Touches, by source:

- Design, from I-0179, I-0178, and the practice of recording acceptance steps in an ADR: `design/adrs/README.md`, `design/system/project-manifest.md`, `design/system/flai-cli.md`.
- Layout, from the explorer's read of the code:
  - `flai/internal/manifest/manifest.go` and `manifest_test.go`, which hold `Project`.
  - `flai/cmd/accept.go`, whose `acceptItem` every acceptance mode passes through.
  - The new `flai/cmd/accept_install.go` and `accept_install_test.go`, beside `accept_fix_test.go`.
  - `system-flow.yaml`, for flaiover's `install`.
  - `docs/operators/settings.md`, since `settings_doc_test.go` needs a row for each manifest key.
  - `docs/users/flai.md` and `design/system/devex.md`.
- Declared by the criteria: the I-0179 file.
- Co-change, from `flai touches suggest` on the I-0179 file: `design/issues/summary.md`.
- Left out of what `flai touches suggest` listed: `harness.mjs`, `component-library.md`, `scripts/README.md`, `smoke.sh`, and `Fetcher.svelte`. They changed with the harness test. This plan does not change that test; S-0482 changes it and `scripts/integration.sh`.
- No folder touch is kept. The ADR's own file has no number until `flai adr new` gives it one, and `flai task done` adds it to T-1998's touches.

Figures:

- Forecast 1h, adjusted from flai's 24m. flai's figure is 89 s per unit of size over 33 done remediation stories in the large band, times size 16. The sample's remediations mostly change one function or one test. This one adds an ADR, a manifest key with its check and settings row, and an acceptance step tested through a git merge fixture. Delivery moves by the 36m difference, from 2026-10-22T06:34:00Z to 07:10:00Z.
- Cost of delay 37.50 USD a week, from `flai cod`: 15m lost per 168h cycle at 150 USD an hour. It stands. I-0178 records the same cause at 20m for S-0482, so the joint cost is higher, but each story carries its own issue's input.

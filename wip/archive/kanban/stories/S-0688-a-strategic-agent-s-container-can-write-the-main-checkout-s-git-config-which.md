---
id: S-0688
type: story
nature: improvement
title: A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host
status: done
owner: alex
created: 2026-11-10T08:37:15Z
updated: 2026-11-11T02:28:57Z
transitions:
  - to: ready
    at: 2026-11-10T11:00:29Z
    by: alex
  - to: in-progress
    at: 2026-11-10T19:29:44Z
    by: agent-S-0688
  - to: review
    at: 2026-11-10T19:50:12Z
    by: agent-S-0688
  - to: in-progress
    at: 2026-11-11T01:42:29Z
    by: flai
  - to: review
    at: 2026-11-11T01:50:13Z
    by: agent-S-0688
  - to: done
    at: 2026-11-11T02:28:57Z
    by: alex
tags: [flai, container, security]
touches: [design/adrs, design/adrs/README.md, flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/serve/container_test.go, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, design/issues/I-0257-a-strategic-agent-s-container-can-write-the-main-checkout-s-git-config-which.md, design/issues/summary.md, flai/internal/serve/container.go, design/system/agent-adapters.md, design/issues/I-0346-a-trial-merge-that-finds-two-open-stories-editing-different-sentences-of-one.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1704
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 4
      work: 43
    - day: 2026-11-11
      ceremony: 3
      work: 46
  models:
    - model: claude-opus-5-5
      input: 236
      output: 40718
      cache_read: 31121341
      cache_write: 755247
      cost: 12.9169
  strategic:
    - kind: orchestrator
      seconds: 125
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 27
          output: 180
          cache_read: 1378745
          cache_write: 162221
          cost: 0.4235
        - model: claude-sonnet-5-5
          input: 10
          output: 101
          cache_read: 54536
          cache_write: 69692
          cost: 0.0936
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T08:37:15Z
  value: 22.5
  by: planner-S-0688
  at: 2026-11-10T10:34:24Z
forecast:
  duration: 45m
  delivery: 2026-11-10T20:27:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 3, behind S-0666, S-0667 and S-0708."
  by: flai
  at: 2026-11-10T19:26:06Z
finalized:
  by: alex
  at: 2026-11-10T11:00:24Z
---
# S-0688 A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host

## Goal

This story remediates [I-0257](../../../design/issues/I-0257-a-strategic-agent-s-container-can-write-the-main-checkout-s-git-config-which.md), "A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0257 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0257 is closed with `flai issue close I-0257 --reason` saying what fixed it

## Tasks
- T-3045 Record in an ADR refining ADR-0247 that a strategic agent's container mounts .git/config read-only
- T-3049 roleStart mounts the main checkout's .git/config read-only, with tests that a strategic container's config is not writable
- T-3052 The design and the guides say a strategic container mounts .git/config read-only
- T-3055 Close I-0257 with the read-only .git/config mount as its fix

## Notes

The analyzer found I-0257: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0257. time_lost_per_cycle 9m: 5m per occurrence × 2 occurrences ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:06:40Z, 7.6 days before this story).

Strategic usage carried over from I-0257: analyzer 42s, 0.8156 USD.

### Planning

Proposed remedy, the first of the two I-0257's first instance weighs: a strategic agent's container mounts the main checkout's `.git/config` read-only over the read-write checkout, as it mounts `.git/hooks`, recorded in a new ADR that refines ADR-0247. ADR-0247 kept `.git/config` writable only because the orchestrator's `flai accept` deletes a merged branch's section from it. ADR-0278 now runs that acceptance on the host through the MCP tool `accept`, and the shell's `flai accept` refuses in a role's container (`flai/cmd/accept.go`, `orchestratorInRoleContainer` in `flai/cmd/verify.go`). Nothing else a role runs there writes `.git/config`: `flai subproject --autocommit` and `flai formula run --autocommit` only add and commit, and the other `git branch -D` sites (`flai/cmd/branch.go`, `flai/cmd/accept_remote.go`, `flai/internal/mergequeue/batch.go`) run on the host. The second remedy, a check of `.git/config` after each strategic run, is proposed on the plan's thread as an addition, not planned.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `design/adrs` | design | The new ADR refining ADR-0247; its file name and number are not known until `flai adr new` writes it. Kept as a folder touch for that one file; the tasks name `design/adrs/README.md`, which narrows the claim (ADR-0096). |
| `design/adrs/README.md` | design | The ADR index `flai adr new` regenerates. |
| `flai/internal/harness/container.go` | layout | `roleStart` builds the strategic container's mounts. |
| `flai/internal/harness/container_test.go` | layout | `TestAStrategicAgentsContainerMountsTheCheckout` asserts those mounts. |
| `flai/internal/serve/container_test.go` | layout | `roleContainerStarted`, the helper the four roles' container start tests share, checks the hooks mount. |
| `design/system/strategic-agents.md` | design | § In a container (S-0509): its mount table, the write-gap paragraph, and the I-0257 line. |
| `design/system/flai-cli.md` | design | The `roleStart` section's mount table. |
| `docs/operators/index.md` | design | The strategic container section's mount table and `.git/config` note. |
| `docs/users/flai.md` | design | The strategic container mount table. |
| `design/issues/I-0257-…md`, `design/issues/summary.md` | declared by criterion 2 | `flai issue close` writes both. |

`flai touches suggest` listed nothing, since the story declared no touches.

No `after`: S-0665, in progress, changes `container.go`, its test, `flai-cli.md`, `docs/users/flai.md`, and `strategic-agents.md` too, to record a run's mounts, but this story reads nothing it writes. The trial merge judges the overlap.

Forecast: with the touches written, `flai forecast` gives 23m (103 s per unit of size over 155 large-band improvement stories, size 13). Raised to 45m from comparable container mount stories: S-0517 took 39m and S-0520 49m, and this one adds an ADR and four documents to a one-line mount change. Delivery is flai's 17:14 plus the 22m added, 17:36.

Cost of delay value: 22.50 USD a week, as `flai cod` works it out from the input flai set from I-0257 (9m lost per 168h cycle at 150 USD an hour). It stands: the time lost is the issue's own measured cost, and no revenue or penalty input exists.
- 2026-11-11T01:42:29Z: moved to in-progress: the rebase of story/S-0688 onto main stopped on conflicts in design/system/agent-adapters.md, docs/operators/index.md, which the acceptance of S-0688 by alex cannot settle, so it aborted the rebase; its agent settles it in the worktree and syncs it with flai stream sync S-0688

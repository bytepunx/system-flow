---
id: S-0679
type: story
nature: improvement
title: "flai check finds `item.heading` outside the story at close-out"
status: review
owner: alex
created: 2026-11-10T08:34:56Z
updated: 2026-11-10T20:06:59Z
transitions:
  - to: ready
    at: 2026-11-10T10:58:07Z
    by: alex
  - to: in-progress
    at: 2026-11-10T19:36:19Z
    by: agent-S-0679
  - to: review
    at: 2026-11-10T20:06:59Z
    by: agent-S-0679
tags: [flai, check, container]
touches: [flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/check/itemwrite.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, design/adrs/README.md, design/issues/I-0316-flai-check-finds-item-heading-outside-the-story-at-close-out.md, design/issues/summary.md, scripts/agent-container.sh, design/adrs/0283-a-close-out-s-record-step-confirms-each-finding-outside-the-story-on-a-work.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1854
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 1
      work: 42
  models:
    - model: claude-opus-5-5
      input: 238
      output: 62229
      cache_read: 29732153
      cache_write: 713122
      cost: 11.6569
  strategic:
    - kind: orchestrator
      seconds: 54
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 74
          cache_read: 456020
          cache_write: 55315
          cost: 0.1407
        - model: claude-sonnet-5-5
          input: 6
          output: 51
          cache_read: 41856
          cache_write: 42397
          cost: 0.0635
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: orchestrator
    at: 2026-11-10T10:55:49Z
forecast:
  duration: 1h
  delivery: 2026-11-10T20:47:00Z
  basis: "Its own forecast of 1h; 3rd in the pull order with an in-progress limit of 3, behind S-0666, S-0667, S-0708, S-0688 and S-0677."
  by: flai
  at: 2026-11-10T19:26:06Z
finalized:
  by: alex
  at: 2026-11-10T10:58:03Z
---
# S-0679 flai check finds `item.heading` outside the story at close-out

## Goal

This story remediates [I-0316](../../../design/issues/I-0316-flai-check-finds-item-heading-outside-the-story-at-close-out.md), "flai check finds `item.heading` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0316 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0316 is closed with `flai issue close I-0316 --reason` saying what fixed it

## Tasks
- T-3014 Confirm a close-out's findings on work items under the checked-write lock before recording them
- T-3017 Mount the main checkout's checked-write lock file in the close-out's container
- T-3018 Record in an ADR, the design, and the guides that a close-out confirms work-item findings under the checked-write lock
- T-3021 Close I-0316 with what fixed it

## Notes

### Planning

The cause, from the one instance. S-0599's close-out found S-0636's body missing `## Tasks` at 2026-11-09T07:36:46Z. S-0636 was created at 07:36:47Z, and its first commit (`1a0ee3c7b`, S-0599's acceptance) has every section. The close-out's check read S-0636 while a checked creation was between its write and its keep or removal. [ADR-0272](../../../design/adrs/0272-a-checked-creation-holds-a-checked-write-lock-from-its-write-until-it-keeps-or.md) puts only flai serve's look under the checked-write lock, and says other readers read without it. The planner proposes this fix: `recordOutside` in `flai/cmd/check.go`, which `flai check --record-issues` and `flai verify --record-issues` both call, confirms each finding outside the story on a work item under `LockCheckedWrites` with `check.RunItemWrite`, and records only what still holds. The default agent runs in a container, where the close-out's verify runs (ADR-0274), so that container mounts the main checkout's `.flai-cache/items-checked.lock` too. The plan's thread is TH-1030.

Touches, by where each came from:

| Touch | Source |
|---|---|
| `flai/cmd/check.go` | layout: `recordOutside`, where findings outside the story are recorded |
| `flai/cmd/check_test.go` | co-change (31%): the test that reproduces the cause |
| `flai/internal/check/itemwrite.go` | layout: `RunItemWrite`, which the confirmation runs, if it needs a helper |
| `flai/internal/harness/container.go` | layout: `ContainerCommand`, the close-out container's mounts |
| `flai/internal/harness/container_test.go` | co-change (14%): the mount's test |
| `docs/users/flai-reference.md` | co-change (35%): regenerated if `--record-issues`'s help changes |
| `docs/users/flai.md` | co-change (24%): the `--record-issues` paragraph and the container mounts |
| `docs/operators/runbooks/agent-container.md` | design: the close-out container's mounts |
| `design/system/flai-cli.md` | co-change (22%): the record step's behaviour |
| `design/adrs/README.md` | co-change (8%): the index the new ADR is listed in |
| `design/issues/I-0316-flai-check-finds-item-heading-outside-the-story-at-close-out.md` | declared by the goal: the issue closed |
| `design/issues/summary.md` | co-change: written again when an issue closes |

`flai touches suggest S-0679` needed seed paths, since the story declared no touches. It was given `flai/cmd/check.go`, `flai/internal/check/itemwrite.go`, `flai/internal/workitem/itemlock.go`, and `flai/internal/harness/container.go`. The new ADR's file has no name yet, and `flai task done` widens the touches with it.

The forecast: `flai forecast` gave 24m once the touches were set (103 s per unit over 155 done large improvement stories, size 14), and 2m before. It is raised to 1h because the story makes two code changes in two packages, needs a concurrency test, writes an ADR, and runs the integration and smoke tiers a change under `flai/` selects at close-out. The delivery is flai's 15:10Z moved 36m later to match.

The cost of delay: `flai cod` cannot work out a value, since the story has no inputs and no epic. The planner asked the operator for `time_lost_per_cycle` on TH-1030, recommending 2m, and sets no value until it is answered.

No `after`: the story reads nothing another open story writes.

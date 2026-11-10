---
id: S-0710
type: story
nature: improvement
title: "flai check finds `adr.index` outside the story at close-out"
status: backlog
owner: alex
created: 2026-11-10T09:56:02Z
updated: 2026-11-10T11:01:50Z
transitions: []
tags: [flai, adr]
touches: [flai/internal/adr/adr.go, flai/internal/adr/adr_test.go, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0336-flai-check-finds-adr-index-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: orchestrator
    at: 2026-11-10T10:55:46Z
forecast:
  duration: 45m
  delivery: 2026-11-11T06:46:00Z
  basis: "Its own forecast of 45m; 38th in the pull order with an in-progress limit of 3, behind S-0665, S-0667, S-0708, S-0674, S-0704, S-0666, S-0675, S-0676, S-0677, S-0678, S-0679, S-0681, S-0688, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0686, S-0689, S-0690, S-0691, S-0692, S-0693, S-0694, S-0695, S-0696, S-0697, S-0698, S-0699, S-0700, S-0701, S-0702, S-0703, S-0705, S-0706, S-0707 and S-0709."
  by: flai
  at: 2026-11-10T11:01:50Z
---
# S-0710 flai check finds `adr.index` outside the story at close-out

## Goal

This story remediates [I-0336](../../../design/issues/I-0336-flai-check-finds-adr-index-outside-the-story-at-close-out.md), "flai check finds `adr.index` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0336 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0336 is closed with `flai issue close I-0336 --reason` saying what fixed it

## Tasks
- T-2980 Commit an ADR index the committed ADR files render
- T-2983 Say in the design and the user guide what an ADR commit records of the index
- T-2986 Close I-0336 with what fixed it

## Notes

### Planning

The cause, found from the instances and main's history. Both instances (S-0680 at 09:08:17Z, S-0685 at 09:11:58Z) report rows 0101 and 0122 of `design/adrs/README.md`. On main, `859f16842` (`docs: ADR-0103 accepted`) changed row 0101 and `4a247f7a9` (`docs: ADR-0123 accepted`) changed row 0122, while the ADR-0101 and ADR-0122 files stayed `proposed` in HEAD until `b4e28372c` (S-0664's acceptance, 09:11:32Z) committed them from the main checkout's uncommitted files ([ADR-0248](../../../design/adrs/0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md)). `adr.Accept`, `adr.New`, and `adr.SetTopics` write the index from the working tree's ADR files (`renderIndex`, `adrindex.Write`), and `Result.commit` then commits the index with the one ADR file it wrote. So when another ADR file is left uncommitted, here because its own `flai adr accept --autocommit` from the dashboard wrote it but did not commit it, the commit records an index that the committed ADR files do not render. Every worktree made from that main then finds `adr.index` outside its story. This planner assumes the fix: the commit records the index as HEAD's ADR files, with the ones the commit records, render it, and the working tree keeps the index its own files render. [ADR-0228](../../../design/adrs/0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md) § 1 is the rule the index follows. The plan's thread is TH-1021.

Touches, by where each came from:

| Touch | Source |
|---|---|
| `flai/internal/adr/adr.go` | layout: `Result.commit` and `renderIndex`, where the cause is |
| `flai/internal/adr/adr_test.go` | layout: the test that reproduces the cause |
| `design/system/flai-cli.md` | design: the `flai adr` row says what `--autocommit` and `--commit` commit |
| `docs/users/flai.md` | design: "Record a decision" says what the commit holds |
| `design/issues/I-0336-flai-check-finds-adr-index-outside-the-story-at-close-out.md` | declared by the goal: the issue closed |
| `design/issues/summary.md` | co-change: written again when an issue closes |

`flai touches suggest S-0710` listed nothing, since the story declared no touches. A helper that renders the index from git's tree may belong in `flai/internal/adrindex/adrindex.go`; it is left out until a task needs it, and `flai task done` widens the touches then.

The forecast: `flai forecast` gave 15m once the touches were set (107 s per unit over 30 done medium improvement stories, size 8), and 2m before. It is raised to 45m because a change under `flai/` runs the integration and smoke tiers in its close-out and needs a git-level test, as S-0685's 45m and S-0680's 45m planned. The delivery is flai's 16:34Z moved 30m later to match.

The cost of delay: `flai cod` cannot work out a value, since the story has no inputs and no epic. The planner asked the operator for `time_lost_per_cycle` on TH-1021, recommending 2m, and sets no value until it is answered.

No `after`: the story reads nothing another open story writes.

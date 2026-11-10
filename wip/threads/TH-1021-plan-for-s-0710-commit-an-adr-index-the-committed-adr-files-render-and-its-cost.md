---
id: TH-1021
title: "Plan for S-0710: commit an ADR index the committed ADR files render, and its cost of delay"
anchor:
  path: wip/kanban/stories/S-0710-flai-check-finds-adr-index-outside-the-story-at-close-out.md
  item: S-0710
status: open
participants: [planner-S-0710]
created: 2026-11-10T10:30:02Z
updated: 2026-11-10T10:30:02Z
---

# TH-1021 Plan for S-0710: commit an ADR index the committed ADR files render, and its cost of delay

On wip/kanban/stories/S-0710-flai-check-finds-adr-index-outside-the-story-at-close-out.md.

## Entries

### 2026-11-10T10:30:02Z planner-S-0710
Question first. S-0710 has no cost of delay inputs and no epic, so `flai cod S-0710` cannot work out a value.

Recommended: `time_lost_per_cycle: 2m`, the figure S-0680 carries. Each instance costs an agent a finding outside its story to read and an issue instance to record at close-out, and two close-outs met it.

Alternatives: 5m, as S-0713 carries, if you count the operator's time to trace a finding outside the story; or no input, leaving the story without a value.

## The cause

Both instances (S-0680, S-0685) name rows 0101 and 0122 of `design/adrs/README.md`. On main, `859f16842` (`docs: ADR-0103 accepted`) changed row 0101, and `4a247f7a9` (`docs: ADR-0123 accepted`) changed row 0122. The ADR-0101 and ADR-0122 files stayed `proposed` in HEAD until S-0664's acceptance (`b4e28372c`, 09:11:32Z) committed the main checkout's uncommitted files. `flai adr accept`, `new`, and `topics` write the index from the working tree's ADR files, and then commit it with only the one ADR file they wrote. So an ADR left uncommitted puts its row in another ADR's commit.

## Tasks and layers

| Layer | Task | Touches | after |
|---|---|---|---|
| 1 | T-2980 Commit an ADR index the committed ADR files render, with a git-level test that reproduces the cause | `flai/internal/adr/adr.go`, `flai/internal/adr/adr_test.go` | none |
| 2 | T-2983 Say in the design and the user guide what an ADR commit records of the index | `design/system/flai-cli.md`, `docs/users/flai.md` | T-2980 |
| 2 | T-2986 Close I-0336 with what fixed it | the I-0336 file, `design/issues/summary.md` | T-2980 |

T-2983 and T-2986 share no path, so they run together.

## Figures

- Forecast: 45m, delivery 2026-11-10T17:04:00Z. `flai forecast` gave 15m (107 s per unit, size 8). It is raised because a change under `flai/` runs the integration and smoke tiers at close-out, as S-0680 and S-0685 planned. The delivery is moved 30m later to match.
- Cost of delay: none until the question above is answered.

## Assumptions

- The fix: the commit records the index that HEAD's ADR files, with the files the commit records, render. The working tree keeps the index its own files render. It covers `--autocommit` and the worktree's `--commit` alike. The index stays generated as ADR-0228 § 1 says.
- Why ADR-0101's and ADR-0122's own dashboard commits did not land is not known. The fix keeps main's index consistent whatever leaves an ADR file uncommitted. Showing a `commit_error` on the dashboard would be a separate story; I propose none here.
- No new ADR. The change refines how ADR-0228's generated index is committed, and the `flai adr` row in `design/system/flai-cli.md` records it. If the story's agent finds it is a decision, it adds one, and `flai task done` widens the touches.
- A helper that renders the index from git's tree may go in `flai/internal/adrindex/adrindex.go`. It is left out of the touches until T-2980 needs it.
- No `after` on another story.

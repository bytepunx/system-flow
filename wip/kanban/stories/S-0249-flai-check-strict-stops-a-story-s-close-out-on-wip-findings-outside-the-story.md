---
id: S-0249
type: story
nature: improvement
title: flai check --strict stops a story's close-out on wip/ findings outside the story
status: ready
owner: alex
created: 2026-10-03T18:03:58Z
updated: 2026-10-05T00:25:39Z
transitions:
  - to: ready
    at: 2026-10-05T00:13:03Z
    by: alex
tags: [flai, template]
topics: [cli, conventions, template]
touches: [flai/internal/check, flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/issues, scripts/close-out.sh, scripts/check.sh, template/root/scripts/close-out.sh, template/root/scripts/check.sh, template/CHANGELOG.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/adrs, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 7m
    by: planner-S-0249
    at: 2026-10-05T00:17:24Z
  value: 17.5
  by: planner-S-0249
  at: 2026-10-05T00:19:21Z
forecast:
  duration: 1h
  delivery: 2026-10-05T01:33:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, behind S-0252 and S-0259."
  by: flai
  at: 2026-10-05T00:25:39Z
---
# S-0249 flai check --strict stops a story's close-out on wip/ findings outside the story

## Goal

This story remediates [I-0057](../../../design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md), "flai check --strict stops a story's close-out on wip/ findings outside the story". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0057 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0057 is closed with `flai issue close I-0057 --reason` saying what fixed it

## Tasks
- T-0840 flai check --story reports findings outside the story as notes that --strict passes over
- T-0841 flai check --story --record-issues opens or bumps one issue per rule for the findings outside the story
- T-0844 The close-out scripts, here and in the template, scope flai check to the story and record what lies outside it
- T-0845 An ADR, the design, the conventions, and the user guide say that a close-out reports findings outside the story as notes recorded in issues
- T-0847 S-0249's own close-out ends clean past main's findings, and I-0057 is closed

## Notes

### Planning

The proposed solution follows I-0057's 19 instances and the operator's answer on TH-0056, "notes are preferable when its a finding outside the story", with flai recording each one as an issue.

- `flai check --story S-nnnn` treats a finding as inside the story when its path is one of the following: the story's file or its tasks' files, its narrative, its threads, or a path in its diff against main. A `wip.overlap` with another open story is never inside the story.
- Findings outside the story are notes the run passes over, errors included: main's `issues.duplicate-id` stopped S-0200 and S-0203. The unscoped check, as CI runs it, gates everything as before.
- `--record-issues` opens or bumps one issue per rule, and writes each story's instance once.
- `scripts/close-out.sh` exports the story, and `scripts/check.sh` passes both flags when it is set. That covers both paths to the check: the smoke tier and the direct call. `TestMonorepoIsClean` is scoped the same way.

Touches, none declared before. `flai touches suggest S-0249` listed nothing until it was given starting paths. I gave it `flai/internal/check`, `flai/cmd/check.go`, `scripts/close-out.sh`, and `flai/internal/issues`, and its co-change list covers 54 of 809 commits.

- `flai/internal/check`, `flai/cmd/check.go`, `flai/cmd/check_test.go`: layout. These are the check package and its command, where the scope, the outside marker, and the reproduction tests go.
- `flai/internal/issues`: layout. `New` and `Bump` live there, and the lookup by rule goes there.
- `scripts/close-out.sh`, `scripts/check.sh`: layout. The close-out reaches the check through `check.sh` directly, and through `smoke.sh` when the story touches `flai`.
- `template/root/scripts/close-out.sh`, `template/root/scripts/check.sh`: design. They are the template's copies, and convention changes land in the template in the same story.
- `template/CHANGELOG.md`: co-change. 2 of 2 recent commits to the template's close-out changed it.
- `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`: co-change (5 of 54) and design. They hold the close-out rule.
- `design/system/flai-cli.md`, `docs/users/flai.md`: co-change, 19 and 20 of 54. They describe `flai check`.
- `docs/users/flai-reference.md`: layout. It is regenerated from the help text.
- `design/system/continuous-improvement.md`: co-change (4 of 54) and design. It describes how issues are recorded.
- `design/adrs`: design. The new ADR refines ADR-0073, and `design/adrs/README.md` co-changes (7 of 54).
- `design/issues`, with I-0057 and `summary.md`: design and co-change (11 of 54). The second criterion closes I-0057, and the close-out writes the issues it records.

Forecast: flai gave 4m with no touches, and 30m once the touches were declared (89 s per unit of size times size 20). I raised it to 1h. The work is five tasks in four layers: the scope, the issue recording, two scripts and their template copies, an ADR with three documents and two conventions, and a close-out rehearsal. S-0252, a comparable flai remediation, was planned at 40m for three tasks. The delivery of 2026-10-05T01:45Z allows for S-0249 being held behind S-0252, in progress, which shares `flai/internal/issues`, `design/issues`, `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/users/flai-reference.md`.

Cost of delay: 17.50 USD a week, as `flai cod` works it out from `time_lost_per_cycle: 7m`: 7m per 168h cycle at 150 USD an hour. The input is I-0057's recorded cost, which the operator chose on TH-0112. It stands unadjusted. The issue's instances come about 58 a week, so the figure understates the delay, but the operator chose it knowing the frequency, and it ranks S-0249 the same way flai ranks the remediations it makes from issues.

Tags `flai` and `template`: `flai check` asked for one, because the story delivers to both components. Topics `cli`, `conventions`, and `template`, which the tasks reach.

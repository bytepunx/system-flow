---
id: S-0446
type: story
nature: improvement
title: "flai check finds `components.examples` outside the story at close-out"
status: review
owner: alex
created: 2026-10-17T16:03:19Z
updated: 2026-10-18T03:41:55Z
transitions:
  - to: ready
    at: 2026-10-18T02:53:21Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-18T02:54:45Z
    by: agent-S-0446
  - to: review
    at: 2026-10-18T03:41:55Z
    by: agent-S-0446
tags: [flai, check, components]
topics: [components, close-out]
touches: [flai/internal/check/check.go, flai/internal/check/components_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/system/component-library.md, design/system/flai-cli.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, design/adrs, design/adrs/README.md, design/issues/I-0161-flai-check-finds-components-examples-outside-the-story-at-close-out.md, design/issues/summary.md, design/system/workflow.md, flai/cmd/check_stats_test.go, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2858
  priced_by: harness
  turns:
    - day: 2026-10-18
      hand_edits: 2
      work: 52
  models:
    - model: claude-opus-5-5
      input: 224
      output: 63541
      cache_read: 11944155
      cache_write: 547689
      cost: 6.9742
  strategic:
    - kind: orchestrator
      seconds: 866
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 77
          output: 1683
          cache_read: 29772089
          cache_write: 52399
          cost: 7.4768
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: orchestrator
    at: 2026-10-18T02:32:43Z
  value: 25
  by: planner-S-0446
  at: 2026-10-18T02:39:01Z
forecast:
  duration: 45m
  delivery: 2026-10-18T21:36:00Z
  basis: "flai forecast's 26m (86 s per unit of size times size 18, 76 done large improvement stories), raised to 45m because S-0277 and S-0280, each one scoped leave-out kind, took 32m and 22m, and this story adds a check change and an ADR to one; delivery moved by the 19m added."
  by: planner-S-0446
  at: 2026-10-18T02:39:01Z
finalized:
  by: orchestrator
  at: 2026-10-18T02:39:45Z
---
# S-0446 flai check finds `components.examples` outside the story at close-out

## Goal

This story remediates [I-0161](../../../design/issues/I-0161-flai-check-finds-components-examples-outside-the-story-at-close-out.md), "flai check finds `components.examples` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0161 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0161 is closed with `flai issue close I-0161 --reason` saying what fixed it

## Tasks
- T-1899 A sub-project without svelte installed gets one components.examples note, not a warning for each examples file
- T-1900 A story's scoped check leaves out the no-svelte components.examples note outside the story, so the close-out records no issue for it
- T-1901 I-0161 is closed with what fixed it

## Notes

### Planning

The cause, from both instances (S-0392 and S-0393): a story worktree has no `flaiover/node_modules`. `components.ReadProps` returns `ErrSvelteMissing`, and `exampleProps` in `flai/internal/check/check.go` warns on each of about 48 examples files. `ScopeToStory` in `flai/internal/check/scope.go` does not leave `components.examples` out, so `recordOutside` in `flai/cmd/check.go` records them in I-0161. Nothing installs flaiover before the check: `scripts/flaiover-install.sh --if-needed` runs only in the test tiers, after it.

Proposed fix, in two tasks:

- T-1899: one advisory note for each sub-project without svelte, as a host without node gets already.
- T-1900: a story's scoped check leaves that note out when the story claims nothing under the examples folder.

Rejected:

- Leaving out every `components.examples` finding outside the story: it would hide real format errors.
- Installing flaiover before the check: it gives the check a side effect on every close-out.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/check/check.go`, `components_test.go` | layout: the explorer found the finding at `exampleProps` and its tests |
| `flai/internal/check/scope.go`, `scope_test.go` | layout: `ScopeToStory` holds the left-out kinds |
| `flai/cmd/check.go`, `check_test.go` | layout: `--story` help, `recordOutside`, and the `--record-issues` tests |
| `design/system/component-library.md`, `design/adrs/README.md`, and the `design/adrs` folder | design: ADR-0164 §5 states the per-file warning this changes |
| `design/system/flai-cli.md`, `docs/users/flai.md` | design: the `flai check` row and the user docs of the check and its scope |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | design: the close-out rule lists the left-out kinds, in the project and in the template |
| `design/issues/I-0161-…md`, `design/issues/summary.md` | goal: criterion 2 closes the issue |

The story declared no touches. `flai touches suggest` from these paths gave no file worth adding: its top co-changes are broad docs (`docs/operators/index.md`, `design/system/flaiover-dashboard.md`) at 11% or less, and the earlier I-0073 and I-0078 issue files.

One folder touch was kept: `design/adrs`. The new ADR's file name is not known until `flai adr new` makes it. The folder is in `claims.shared`, so it holds no other story.

Figures:

- Forecast 45m, up from flai's 26m. S-0277 (`threads.archived`) and S-0280 (`item.archive`) each added one scoped leave-out kind and took 32m and 22m. This story adds a check change and an ADR to one. Delivery 2026-10-18T21:36:00Z is flai's 21:17 plus the 19m added.
- Cost of delay 25 USD a week, as `flai cod` gives. The input is 10m lost per 168h cycle at 150 USD an hour. The orchestrator set it on TH-0507, at the figure recommended there.

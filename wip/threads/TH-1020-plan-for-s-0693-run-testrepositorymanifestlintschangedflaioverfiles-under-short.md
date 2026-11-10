---
id: TH-1020
title: "Plan for S-0693: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0693-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md
  item: S-0693
status: open
participants: [planner-S-0693]
created: 2026-11-10T10:29:43Z
updated: 2026-11-10T10:29:43Z
---

# TH-1020 Plan for S-0693: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers

On wip/kanban/stories/S-0693-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md.

## Entries

### 2026-11-10T10:29:43Z planner-S-0693
Plan for S-0693, remediating I-0327.

Proposed fix: drop the `testing.Short()` gate from `TestRepositoryManifestLintsChangedFlaioverFiles` in `flai/internal/verify/select_test.go`. It reads only the manifest and the file tree, needs no git, and its three repository-manifest neighbours already run under `-short`. The go-test tier then catches a tier it must list, not the close-out's integration tier.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2982 Run TestRepositoryManifestLintsChangedFlaioverFiles under -short | `flai/internal/verify/select_test.go` | none |
| 2 | T-2985 Close I-0327 with what fixed it | the I-0327 file, `design/issues/summary.md` | T-2982 |

## Figures

- Forecast 45m, adjusted from flai's 2m to count the close-out's integration and smoke tiers. Delivery 2026-11-10T14:51Z.
- Cost of delay 25 USD a week, from the operator's 10m per cycle input, unchanged.

## Assumptions

- The thread the orchestrator's prompt named ("asked") does not exist, and S-0693 had no thread and no earlier plan. I planned it from scratch.
- The scope is this one test. Three other tests skip under `-short` with the same message: `TestPendingOnThisRepositoryAnswersAsTheProcessesDid` (release), `TestRoundTripRepositoryItems` (workitem), and `TestShippedRoleFilesAreClean` (check). I have not checked whether they need git. If you want them reviewed too, I would add a task to S-0693 or draft a separate story; say which.
- No change to `design/conventions/tooling.md`: its list of the tests that skip under `-short` does not name this one.

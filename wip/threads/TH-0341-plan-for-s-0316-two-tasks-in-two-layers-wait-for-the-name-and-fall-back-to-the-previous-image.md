---
id: TH-0341
title: "Plan for S-0316: two tasks in two layers, wait for the name and fall back to the previous image"
anchor:
  path: wip/kanban/stories/S-0316-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md
  item: S-0316
status: resolved
participants: [planner-S-0316, orchestrator]
created: 2026-10-07T23:55:38Z
updated: 2026-10-07T23:55:48Z
---

# TH-0341 Plan for S-0316: two tasks in two layers, wait for the name and fall back to the previous image

On wip/kanban/stories/S-0316-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md.

## Entries

### 2026-10-07T23:55:38Z planner-S-0316
S-0316 is planned: 8 touches, a 30m forecast, and a cost of delay value of 30 USD a week. It stays a draft for you to finalize.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1295: reproduce the lingering name with a `fakeRunner` stand-in, wait for the name to be released after `docker stop` (bounded, then `docker rm -f`), and start the previous image again when the new one fails | `flai/cmd/dashboard_upgrade.go`, `flai/cmd/dashboard_test.go`, `docs/users/flai-reference.md` | none |
| 2 | T-1296: operator docs and design say what the upgrade does now; `flai issue close I-0094` | `docs/operators/index.md`, `docs/operators/runbooks/update.md`, `design/system/flai-cli.md`, the I-0094 file, `design/issues/summary.md` | T-1295 |

T-1296 waits because the docs describe the limit and the failure report that T-1295 settles.

## Assumptions

- The fix waits for the name and adds the rollback. The rename direction is not taken: the swap's code already proves the image under a temporary name, and Docker's rename of a running container still needs the name to be free.
- `flai dashboard restart` has the same stop-then-run race, so T-1295 fixes it too. That widens the story slightly past the issue's wording.
- `docs/users/flai.md`, `design/system/flaiover-dashboard.md`, and `flai/cmd/dashboard.go` are co-change suggestions left out. The story's agent widens the touches if it needs them.
- The forecast is raised from flai's 16m to 30m. The fake runner needs new behaviour, and two commands plus a rollback change.
- The value is 30 USD a week, from flai's inputs. The issue's text counts four occurrences against the two in its `count`, so the inputs understate the loss. Raising them is your call.

No input is missing. Nothing here needs an answer before you finalize.

### 2026-10-07T23:55:47Z orchestrator
Approved.

- T-1295 removes the cause with a stand-in that holds the name after `stop` (criterion 1). T-1296 closes I-0094 (criterion 2).
- The remedy is the direction I-0094 itself recommends: wait for the name with a limit, and start the previous image again on failure. No recorded decision changes.
- Fixing `flai dashboard restart`'s identical stop-then-run race in the same file stays within "the cause no longer occurs".
- The 30m forecast and the 30 USD a week value stand.

For the operator: the issue's text counts four occurrences against the `count` of two that flai's 12m input rests on. Raising the input is yours, with `flai edit S-0316 --time-lost-per-cycle <duration>`.

Source: design/issues/I-0094-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md

### 2026-10-07T23:55:48Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories

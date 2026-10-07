---
id: TH-0328
title: "S-0236 plan: eight tasks in five layers, forecast raised to 1h30m"
anchor:
  path: wip/kanban/stories/S-0236-flai-dashboard-resolves-the-image-from-the-signed-digest-list-and-runs-it-by-digest.md
  item: S-0236
status: resolved
participants: [planner-S-0236, orchestrator]
created: 2026-10-07T22:55:44Z
updated: 2026-10-07T22:55:52Z
---

# TH-0328 S-0236 plan: eight tasks in five layers, forecast raised to 1h30m

On wip/kanban/stories/S-0236-flai-dashboard-resolves-the-image-from-the-signed-digest-list-and-runs-it-by-digest.md.

## Entries

### 2026-10-07T22:55:44Z planner-S-0236
S-0236 now has eight backlog tasks in five layers, file-level touches, a 1h30m forecast, and a cost of delay of 3.95 USD a week.

## Tasks and layers

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-1239 | `selfupgrade` lists the flaiover GitHub releases with their assets and downloads an asset | none |
| 1 | T-1240 | new `flai/internal/dashboard`: parse, verify, and cache the digest list by version | none |
| 2 | T-1241 | resolve `dashboard.tag` to a release and a by-digest reference; refuse a release with no list it can verify | T-1239, T-1240 |
| 2 | T-1242 | `flai dashboard versions` marks releases that are not deployable; `--published` refuses them | T-1239, T-1240 |
| 3 | T-1243 | `flai dashboard` runs by digest or refuses; `--image` and `--build` reported unsigned; `status` shows the digest; restart keeps it | T-1241 |
| 3 | T-1245 | Updates page shows a release that is not deployable and offers no deploy | T-1242 |
| 4 | T-1244 | `check` and `upgrade` compare and upgrade by digest | T-1243 (both change `dashboard.go` and `dashboard_test.go`) |
| 5 | T-1246 | settings, update runbook, `flai.md`, reference, `flai-cli.md`, `flaiover-dashboard.md`, `docker.md` | T-1244, T-1245 |

## Figures

- **Forecast:** 1h30m, up from 1h15m. `flai forecast` gives 58m. This story has eight tasks against S-0233's five, which is forecast at 1h15m.
- **Cost of delay:** 3.95 USD a week, as `flai cod` gives it, up from 3.29: a 1h30m share of the 9h30m over E-0015's open stories.
- **Other stories:** their shares were worked out from the earlier total, so the epic's other values are now slightly out of date until they are planned again.

## Assumptions

- S-0233 has an ECDSA verifier in `flai/internal/selfupgrade` that takes its trusted keys through `Options`, plus a test key pair under `testdata/` (T-1226). T-1240 reuses both rather than adding keys of its own.
- S-0234 records the digest list's format in `design/system/release-signing.md`, and its releases are tagged `flaiover/vX.Y.Z` with `flaiover_<version>.digests` and its signature as assets (T-1231, T-1232).
- "Beside `flai serve`'s state" means under the directory `serve.DirFor` gives. Reading the cache verifies again, and S-0238 looks it up by digest.
- The host API and flaiover's `/api/dashboard` pass `flai dashboard versions --json` through unchanged, so only `HostPanel.svelte` changes on the dashboard side.
- A restart record from an older flai, holding an `image:tag` reference, still restarts.

## Proposals

None: no task to split, merge, or drop. No input from you is needed.

### 2026-10-07T22:55:51Z orchestrator
Approved. The eight tasks cover the six criteria:

- T-1241 and T-1243 cover criterion 1.
- T-1243 (status) and T-1244 cover criterion 2.
- T-1240 covers criterion 3.
- T-1243 covers criterion 4.
- The tests across T-1239 to T-1244, and T-1246, cover criterion 5.
- T-1242 and T-1245 cover criterion 6.

No two tasks of a layer share a path. T-1244 waiting on T-1243 for `dashboard.go` follows the layering rule. Both folder touches narrow to the files the tasks name. The 1h30m forecast and the 3.95 USD a week value stand.

The other E-0015 values that the new total leaves stale will be refreshed when each of those stories is planned. S-0237, S-0238, and S-0239 are next in `flai plan --candidates`.

Source: wip/kanban/stories/S-0236-flai-dashboard-resolves-the-image-from-the-signed-digest-list-and-runs-it-by-digest.md

### 2026-10-07T22:55:52Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories

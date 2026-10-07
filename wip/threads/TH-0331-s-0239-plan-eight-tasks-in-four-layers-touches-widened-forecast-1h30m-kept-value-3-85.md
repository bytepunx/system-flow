---
id: TH-0331
title: "S-0239 plan: eight tasks in four layers, touches widened, forecast 1h30m kept, value 3.85"
anchor:
  path: wip/kanban/stories/S-0239-dashboard-allow-unsigned-lets-a-development-build-connect-shown-on-every-page-and-in-every-status.md
  item: S-0239
status: resolved
participants: [planner-S-0239, orchestrator]
created: 2026-10-07T23:11:03Z
updated: 2026-10-07T23:11:11Z
---

# TH-0331 S-0239 plan: eight tasks in four layers, touches widened, forecast 1h30m kept, value 3.85

On wip/kanban/stories/S-0239-dashboard-allow-unsigned-lets-a-development-build-connect-shown-on-every-page-and-in-every-status.md.

## Entries

### 2026-10-07T23:11:03Z planner-S-0239
S-0239 now has eight backlog tasks in four layers. Its touches gained 12 files, and every declared touch is kept. The forecast stays at 1h30m, and the cost of delay value is now 3.85 USD a week, the figure `flai cod` gives. The reasons are under `### Planning` in its Notes. Its goal and criteria are unchanged.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1259 the setting in the configuration and the manifest, plus its `settings.md` row | none | `config.go`, `manifest.go`, their tests, `docs/operators/settings.md` |
| 1 | T-1260 the dashboard accepts an unsigned flai only with `FLAIOVER_ALLOW_UNSIGNED`, and records the unsigned side | none | `agent.ts`, `release.ts`, their tests |
| 2 | T-1261 `flai dashboard --allow-unsigned`, the setting, or `--build` passes the variable; the reference is regenerated | T-1259 | `dashboard.go`, `dashboard_upgrade.go`, `dashboard_test.go`, `flai-reference.md`, `settings.md` |
| 2 | T-1262 flai dials an unlisted container and accepts an unstamped dashboard under the allowance, with one `warn` per connection | T-1259 | `serve/serve.go`, `channel/channel.go`, their tests |
| 2 | T-1263 the connection list says "unsigned allowed" per project | T-1260 | `api/projects/+server.ts` and its test, `project.svelte.ts`, `ProjectSwitcher.svelte` and its test |
| 3 | T-1265 the three statuses say "unsigned allowed" with the side | T-1261, T-1262 | `cmd/serve.go`, `cmd/host.go`, `cmd/dashboard.go`, their tests |
| 3 | T-1264 a banner on every page names the unsigned side | T-1263 | `UnsignedBanner.svelte` (new) and its test, `+layout.svelte` |
| 4 | T-1266 contributor, operator, user, template, and design docs | T-1264, T-1265 | ten documents |

No two tasks in the same layer touch the same path.

## Assumptions

- **S-0237 and S-0238 are on main first.** The story waits for both. T-1260 and T-1262 relax S-0237's stamp checks in `release.ts`, `agent.ts`, and `channel.go`. T-1262 lets S-0238's image gate pass an unlisted image. If those stories put the checks in other files, the story's agent widens the touches. Each task's Notes says where.
- **The manifest overrides the host.** A project's `dashboard.allow_unsigned`, when set, wins over the host configuration, in either direction. With no manifest value, the host's value applies.
- **`--allow-unsigned` lasts for one start.** It does not write the configuration. A restart by `flai dashboard upgrade`, `restart`, or `flai host` keeps what that start was given.
- **The credential proof always runs.** The allowance lets only a missing or invalid stamp through, and only an unlisted image. The credential is still checked.
- **The dashboard names its side on its own.** Each side knows whether its own stamp is valid and whether the peer's is. The banner needs no change to the protocol.
- **The banner cannot be dismissed.** ADR-0070 says that allowed is never silent.
- **The folder touch on `flaiover/src/lib/components` is kept.** The banner component is new, and its name is my guess. The config and manifest folder touches are kept as declared. In the claim, all three narrow to the files their tasks name.
- **This repository's `.flai-cache/config.json` is git-ignored.** T-1266 documents the setting there. Nothing commits it.

## Proposed changes

None. Eight tasks fit one story at 1h30m. Splitting it into a flai side and a dashboard side would not shorten the chain, because the docs need both.

### 2026-10-07T23:11:10Z orchestrator
Approved. The eight tasks cover the five criteria:

- T-1259 and T-1261 cover criterion 1.
- T-1260, T-1261, and T-1262 cover criterion 2.
- T-1262, T-1263, T-1264, and T-1265 cover criterion 3.
- T-1266 covers criterion 4.
- The tests of T-1260 to T-1265 cover criterion 5.

No two tasks of a layer share a path. The three folder touches narrow to the files the tasks name. The 1h30m forecast and the 3.85 USD a week value stand.

The manifest overriding the host fits the goal's "settable per project in the manifest". The banner that cannot be dismissed follows ADR-0070's "allowed is never silent".

On release: S-0237 is held back from publishing until this story is accepted, as settled on TH-0329. The two go out together.

Source: wip/kanban/stories/S-0239-dashboard-allow-unsigned-lets-a-development-build-connect-shown-on-every-page-and-in-every-status.md

### 2026-10-07T23:11:11Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories

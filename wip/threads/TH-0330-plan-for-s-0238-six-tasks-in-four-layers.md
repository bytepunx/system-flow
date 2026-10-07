---
id: TH-0330
title: "Plan for S-0238: six tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0238-flai-measures-the-container-s-image-through-docker-before-it-dials-and-refuses-one-no-signed-list-names.md
  item: S-0238
status: resolved
participants: [planner-S-0238, orchestrator]
created: 2026-10-07T23:07:00Z
updated: 2026-10-07T23:07:08Z
---

# TH-0330 Plan for S-0238: six tasks in four layers

On wip/kanban/stories/S-0238-flai-measures-the-container-s-image-through-docker-before-it-dials-and-refuses-one-no-signed-list-names.md.

## Entries

### 2026-10-07T23:07:00Z planner-S-0238
S-0238 is planned: six tasks in four layers. Its touches, forecast, and cost of delay are updated.

Tasks:

| Layer | Task | Changes | After |
|-------|------|---------|-------|
| 1 | T-1253 the image check: `docker inspect` and `docker image inspect` through `execx.Runner`, judged against S-0236's verified lists | `flai/internal/serve/imagecheck.go` and its test | |
| 1 | T-1254 a gate the channel client asks before each dial; a refusal skips the dial, records why, logs once, waits a minute | `flai/internal/channel/channel.go` and its test | |
| 2 | T-1255 `flai serve` connects the check to the gate and records the verdict in its status | `flai/internal/serve/serve.go` and its test | T-1253, T-1254 |
| 3 | T-1256 `flai serve status` and `flai host status` show the verdict | `flai/cmd/serve.go`, `flai/cmd/host.go`, their tests | T-1255 |
| 3 | T-1257 `flai dashboard status` shows the verdict | `flai/cmd/dashboard.go` and its test | T-1255 |
| 4 | T-1258 user, operator, and design docs | `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/index.md`, `design/system/flai-cli.md`, `design/system/release-signing.md` | T-1256, T-1257 |

Figures:

- Forecast: 1h15m, raised from 1h. `flai forecast` gives 40m. S-0237 has six tasks with code on both sides and is planned at 1h30m. This story is flai alone, so 1h15m.
- Cost of delay: 3.21 USD a week, up from 2.63. It follows the new forecast's share of E-0015's 25 USD a week.

Assumptions:

1. The check needs a hook in `flai/internal/channel/channel.go`, because `channel.Client.Run` reconnects on its own. I added that file and its test to the touches. S-0237 changes the same file, so while either story is in progress, the other is held. I think that is right, since both change the dial loop's backoff.
2. An unlisted image waits a minute between looks, as on `4409`. The error is logged once per digest. T-1255 also looks again when the container's ID or image ID changes. I read "another look when the container changes" as satisfied by that.
3. "Not local" means three things: Docker is not on `PATH`, no `flaiover` container runs here, or the dashboard's URL is not on this host. Each falls back to the stamp check and says so in status.
4. S-0236's `flai/internal/dashboard` package provides the list cache and a resolver that fetches and verifies a version's list. T-1253 reads both and changes neither. The running version comes from the container's `FLAIOVER_VERSION`.
5. T-1258 settles the **to check** in `release-signing.md`: which digest `RepoDigests` records after a pull by digest.

I propose no split, merge, or drop of tasks.

### 2026-10-07T23:07:07Z orchestrator
Approved. The six tasks cover the five criteria:

- T-1253 and T-1255 cover criterion 1.
- T-1254, T-1255, T-1256, and T-1257 cover criterion 2.
- T-1253 and T-1256 cover criterion 3, with the "not local" cases of assumption 3.
- The tests of T-1253 to T-1257 cover criterion 4.
- T-1258 covers criterion 5.

No two tasks of a layer share a path. `channel.go` overlapping S-0237 is right: S-0238 runs after S-0236, which already orders it behind S-0237's dependencies. The 1h15m forecast and the 3.21 USD a week value stand.

One correction for the story's agent, on assumption 5: S-0234's T-1231 settles which digest `RepoDigests` records (S-0234 criterion 4), and S-0238 runs after it. T-1258 should cite that finding in `release-signing.md`, not check it again.

Source: wip/kanban/stories/S-0234-the-flaiover-image-s-digest-list-is-signed-in-ci-and-published-on-a-flaiover-github-release.md

### 2026-10-07T23:07:08Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories

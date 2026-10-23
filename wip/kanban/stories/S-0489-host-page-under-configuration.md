---
id: S-0489
type: story
nature: feature
title: Host page under configuration
status: in-progress
owner: alex
created: 2026-10-23T17:57:14Z
updated: 2026-10-23T19:07:12Z
transitions:
  - to: ready
    at: 2026-10-23T17:57:21Z
    by: alex
  - to: in-progress
    at: 2026-10-23T18:07:10Z
    by: agent-S-0489
tags: [flai, flaiover, dashboard]
topics: [dashboard, cli]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/lib/server/allowed-hosts.ts, flaiover/src/lib/server/allowed-hosts.test.ts, flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, flaiover/server.js, flaiover/src/lib/settings.ts, flaiover/src/lib/components/HostAddresses.svelte, flaiover/src/lib/components/HostAddresses.svelte.test.ts, flaiover/src/routes/host/+page.svelte, flaiover/src/routes/api/settings/+server.ts, flaiover/src/routes/api/settings/settings.test.ts, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, design/components/flaiover/src/lib/components/HostAddresses.examples.json, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flaiover.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/adrs/0193-the-dashboard-accepts-a-browser-origin-naming-one-of-the-hosts-in-dashboard.md, design/adrs/README.md, flaiover/src/lib/server/agent.ts]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3761
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-23
      ceremony: 2
      hand_edits: 1
      work: 58
  models:
    - model: claude-haiku-5-5
      input: 48
      output: 190
      cache_read: 1795939
      cache_write: 132061
      cost: 0.1352
    - model: claude-opus-5-5
      input: 440
      output: 2382
      cache_read: 25838671
      cache_write: 894281
      cost: 9.0319
  strategic:
    - kind: orchestrator
      seconds: 670
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 174
          output: 2928
          cache_read: 34107903
          cache_write: 30426
          cost: 7.5086
forecast:
  duration: 2h
  delivery: 2026-10-23T20:22:00Z
  basis: "flai forecast gave 1h1m (median 107 s per unit over 86 done feature stories on claude-opus-5-5 in the large band, size 34); doubled because the three layers are serial, two of five tasks wait on two others each, and each task carries a Go or vitest cycle of its own; delivery is the start at 18:08Z plus 2h times the cycle factor 1.12."
  by: planner-S-0489
  at: 2026-10-23T18:10:52Z
---
# S-0489 Host page under configuration

## Goal

Bring back the host page and place it under configuration. In addition, add a way to change what host headers are accepted so that flai can be hosted from one or more custom domains without getting a host header rejection.

## Acceptance criteria
- [x] The Host page is available from a link in the Configuration sub-menu in the nav rail
- [x] The host page includes a way to update configuration that should accept different host headers so that flai can serve traffic through tunnels/domains

## Tasks
- T-2082 The nav rail's Configuration group links the Host page
- T-2083 The dashboard's server accepts a request whose Origin names one of FLAIOVER_ALLOWED_HOSTS
- T-2084 flai keeps dashboard.allowed_hosts and starts the container with FLAIOVER_ALLOWED_HOSTS, kept across restarts and reported by status and settings.get
- T-2085 The Host page edits the accepted hosts through settings.dashboard_hosts, behind the settings host action, and says a restart applies them
- T-2086 The guides say how a tunnel or custom domain is admitted: the Host page's Addresses area, the key, and the restart that applies it

## Notes

### Planning

Planned by planner-S-0489 on 2026-10-23; the plan is in TH-0621, the cost of delay inputs are asked for in TH-0619.

Touches. The story declared none, and `flai touches suggest` refuses a story with none, so the layout and the design gave the first set, and co-change the rest:

- Layout, the rail (T-2082): `flaiover/src/lib/sitemenu.ts` and its test, `NavRail.svelte.test.ts`.
- Layout, the origin check (T-2083): `flaiover/src/lib/server/events-socket.ts` and its test, the new `allowed-hosts.ts` and its test, `flaiover/server.js`.
- Layout, the host configuration and the container (T-2084): `flai/internal/config/config.go`, `flai/cmd/dashboard.go`, `dashboard_upgrade.go`, `dashboard_watch.go`, `serve_actions.go`, each with its test; `docs/operators/settings.md`, whose index test enforces a row per key and variable; `docs/users/flai.md`.
- Layout, the write and the page (T-2085): `flai/internal/hostapi/settings.go` and `writes_test.go`, `flaiover/src/lib/settings.ts`, the new `HostAddresses.svelte` with its test and its examples file under `design/components/`, `flaiover/src/routes/host/+page.svelte`, `flaiover/src/routes/api/settings/+server.ts` and its test.
- Design: `design/system/flaiover-dashboard.md` (the rail table, the Host page, the writes), `docs/users/flaiover.md`, `docs/operators/index.md` (the tunnel expectation, the security posture).
- Co-change: `design/system/flai-cli.md`, changed by 45% of the commits that changed the files above; it describes `flai dashboard` and the host API's methods, which T-2084 and T-2085 extend. `docs/users/flai-reference.md` (23%) is generated from flags, and no flag is added, so it is left out; `design/issues/summary.md` and `design/adrs/README.md` are regenerated, not touched.
- No folder touch: every file is named, the two new components among them.

Forecast. `flai forecast` gave 1h1m from a median of 107 s per unit of size over 86 done feature stories on claude-opus-5-5 in the large band, size 34 (2 criteria, 32 touches). Set to 2h: the three layers are serial, T-2085 waits on T-2083 and T-2084 and T-2086 on T-2082 and T-2085, so little of the work overlaps, and each task carries a Go or vitest cycle of its own. Delivery is the start at 18:08Z plus 2h times the cycle factor 1.12, 20:22Z.

Cost of delay. `flai cod` refuses a value: the story has no inputs and no epic. The value is written once the operator or the orchestrator sets the inputs asked for in TH-0619.

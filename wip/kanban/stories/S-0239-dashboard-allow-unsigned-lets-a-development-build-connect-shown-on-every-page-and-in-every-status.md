---
id: S-0239
type: story
nature: feature
title: dashboard.allow_unsigned lets a development build connect, shown on every page and in every status
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-08T10:28:37Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [flai/internal/config, flai/internal/manifest, flai/cmd/dashboard.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/components, docs/operators/settings.md, docs/users/flaiover.md, flaiover/src/lib/server/agent.test.ts, flai/cmd/serve.go, flai/cmd/host.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/contributors/index.md, template/root/docs/operators/index.md.tmpl, template/CHANGELOG.md, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flaiover/src/routes/+layout.svelte, flaiover/src/routes/api/projects/+server.ts, flaiover/src/routes/api/projects/projects.test.ts, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, flai/cmd/serve_test.go, flai/cmd/host_test.go, flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts, flaiover/src/lib/project.svelte.ts, docs/operators/index.md, design/system/flaiover-dashboard.md, design/system/release-signing.md]
after: [S-0237, S-0238]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 238
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 54
          output: 854
          cache_read: 7040272
          cache_write: 20250
          cost: 1.7397
cost_of_delay:
  value: 3.85
  by: planner-S-0239
  at: 2026-10-07T23:10:45Z
forecast:
  duration: 1h30m
  delivery: 2026-10-08T16:57:00Z
  basis: "Its own forecast of 1h30m; 11th in the pull order with an in-progress limit of 5, behind S-0232, S-0338, S-0342, S-0337, S-0343, S-0233, S-0234, S-0235, S-0236, S-0237 and S-0238."
  by: flai
  at: 2026-10-08T10:28:37Z
---
# S-0239 dashboard.allow_unsigned lets a development build connect, shown on every page and in every status

## Goal

Unsigned builds are allowed only on purpose: `dashboard.allow_unsigned` in the host's flai configuration, settable per project in the manifest, default off (ADR-0070). With it on, flai dials an unlisted container and accepts a peer without a stamp, and `flai dashboard` tells the container to accept an unsigned flai. Allowed is never silent.

## Acceptance criteria
- [ ] `dashboard.allow_unsigned` exists in the configuration and the manifest, default false, with a row in `docs/operators/settings.md` and the generated flag table; `flai dashboard --allow-unsigned` sets it for one start.
- [ ] With it on, or with `--build`, `flai dashboard` starts the container with `FLAIOVER_ALLOW_UNSIGNED=1`; the dashboard accepts a flai without a valid stamp only when that variable is set, and flai dials an unlisted container and accepts a dashboard without a stamp only when the setting is on.
- [ ] Allowed is shown: `flai serve` logs one `warn` per connection naming the unsigned side, `flai serve status`, `flai dashboard status`, and `flai host status` say "unsigned allowed" with the side, the dashboard shows a banner on every page naming which side is unsigned, and the connection list says it per project; a signed pair with the setting on shows nothing.
- [ ] This repository's own contributor documentation says to set the allowance in `.flai-cache/config.json`, since it builds both from source, and the setting is covered by the template's operator documentation.
- [ ] Tests on both sides cover allowed and refused, each side unsigned, and the banner.

## Tasks
- T-1259 dashboard.allow_unsigned exists in the host configuration and the manifest, default false, with its settings row
- T-1260 The dashboard accepts a flai without a valid stamp only when FLAIOVER_ALLOW_UNSIGNED is set, and records which side is unsigned
- T-1261 flai dashboard --allow-unsigned, and the setting or --build, start the container with FLAIOVER_ALLOW_UNSIGNED=1
- T-1262 With the allowance in force, flai serve dials an unlisted container, accepts a dashboard without a stamp, warns once per connection, and records the unsigned side
- T-1263 The dashboard's connection list says "unsigned allowed" per project, with the side
- T-1264 The dashboard shows a banner on every page naming the unsigned side while unsigned is allowed
- T-1265 flai serve status, flai dashboard status, and flai host status say "unsigned allowed" with the side
- T-1266 The contributor, operator, user, template, and design documentation describe dashboard.allow_unsigned and how it is shown

## Notes

### Planning

Tasks, in four layers:

| Layer | Tasks | Why |
|-------|-------|-----|
| 1 | T-1259 (the setting), T-1260 (the dashboard's handshake) | No path in common; the dashboard reads only its variable, not flai's setting |
| 2 | T-1261 (`flai dashboard`) and T-1262 (flai's handshake and dial) after T-1259; T-1263 (connection list) after T-1260 | Each reads what its layer-1 task adds; no path in common |
| 3 | T-1265 (three statuses) after T-1261 and T-1262; T-1264 (banner) after T-1263 | The statuses show the side T-1262 records and change `flai/cmd/dashboard.go` after T-1261; the banner reads T-1263's project list |
| 4 | T-1266 (docs) after T-1264 and T-1265 | Describes the banner and the statuses in the words they print |

Touches:

- Declared, all kept: `flai/internal/config`, `flai/internal/manifest`, `flai/cmd/dashboard.go`, `flaiover/src/lib/server/agent.ts` and its test, `flaiover/src/lib/components`, `docs/operators/settings.md`, `docs/users/flaiover.md`, `flai/cmd/serve.go`, `flai/cmd/host.go`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/project-manifest.md`, `design/system/flai-cli.md`, `docs/contributors/index.md`, `template/root/docs/operators/index.md.tmpl`, `template/CHANGELOG.md`, `flai/internal/serve/serve.go` and its test, `flaiover/src/routes/+layout.svelte`, `flaiover/src/routes/api/projects/+server.ts` and its test.
- Layout, added by this run: `flai/internal/channel/channel.go` and its test, where S-0237 checks the dashboard's stamp and the allowance lets it through; `flaiover/src/lib/server/release.ts` and its test, where S-0237 checks flai's stamp; `flai/cmd/dashboard_upgrade.go`, which builds the `docker run` arguments, and `flai/cmd/dashboard_test.go`; `flai/cmd/serve_test.go` and `flai/cmd/host_test.go`, the tests of the statuses; `flaiover/src/lib/project.svelte.ts`, the client's project list the switcher and the banner read.
- Design, added by this run: `design/system/flaiover-dashboard.md` for the banner and the variable; `design/system/release-signing.md` § The development case, as built; `docs/operators/index.md`, beside S-0237's and S-0238's sections on the stamp and the image check.
- Co-change: `flai touches suggest` listed `design/system/flaiover-dashboard.md` (24%) and `docs/operators/index.md` (15%), added above. Not added: `flai/internal/hostapi/writes.go` (9%), `flai/cmd/serve_actions.go` (6%), and the harness and MCP files, since no host method, host action, or tool changes.
- Folder touches kept, as declared: `flai/internal/config` and `flai/internal/manifest`, where T-1259 names `config.go`, `manifest.go`, and their tests; `flaiover/src/lib/components`, where T-1264 adds a new banner component whose name is a guess, `UnsignedBanner.svelte`, and T-1263 changes `ProjectSwitcher.svelte`. In the claim, each narrows to the files its tasks name (ADR-0096).

Forecast: 1h30m, kept. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gives 1h8m after this run's touches, from a median of 104 s per unit of size over 33 done large-band feature stories, times size 39.
- S-0237 has six tasks in three layers with code on both sides and is planned at 1h30m. This story has eight in four layers on both sides, each smaller: a setting threaded through what S-0237 and S-0238 built. So 1h30m stands over flai's 1h8m.

Cost of delay: 3.85 USD a week, as `flai cod` gives it, down from 3.95 after S-0238's forecast rose. It is this story's 1h30m share of the 9h45m forecast over E-0015's eight open stories, applied to the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure closes only when the chain is done, so a share by work fits.

---
id: S-0236
type: story
nature: feature
title: flai dashboard resolves the image from the signed digest list and runs it by digest
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-08T08:53:22Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/cmd/dashboard.go, flai/internal/selfupgrade, flai/internal/dashboard, docs/operators/settings.md, docs/users/flai.md, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_versions.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_test.go, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/docker.md, flai/cmd/dashboard_versions_test.go, flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, docs/users/flaiover.md, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, flai/internal/dashboard/digests.go, flai/internal/dashboard/digests_test.go, flai/internal/dashboard/cache.go, flai/internal/dashboard/cache_test.go, flai/internal/dashboard/resolve.go, flai/internal/dashboard/resolve_test.go, flai/cmd/dashboard_watch_test.go, docs/operators/runbooks/update.md, docs/users/flai-reference.md]
after: [S-0233, S-0234]
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
      seconds: 325
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 730
          cache_read: 4306919
          cache_write: 24150
          cost: 1.0673
cost_of_delay:
  value: 3.95
  by: planner-S-0236
  at: 2026-10-07T22:55:27Z
forecast:
  duration: 1h30m
  delivery: 2026-10-08T13:34:00Z
  basis: "Its own forecast of 1h30m; 14th in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322, S-0340, S-0341, S-0288, S-0344, S-0290, S-0345, S-0338, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234 and S-0235."
  by: flai
  at: 2026-10-08T08:53:22Z
---
# S-0236 flai dashboard resolves the image from the signed digest list and runs it by digest

## Goal

`flai dashboard`, `flai dashboard check`, and `flai dashboard upgrade` choose the image from a flaiover release's signed digest list, verified with the public key in the binary, and pull and run it by digest (ADR-0070), so that only a signed release of the dashboard is started unless the operator allowed otherwise.

## Acceptance criteria
- [ ] `flai dashboard` resolves `dashboard.tag` to a flaiover release through the release client `self-upgrade` uses (`latest` is the newest release; a version is that release), downloads and verifies its digest list, and pulls and runs `ghcr.io/bytepunx/flaiover@sha256:…`; a list that is missing or does not verify refuses to start the container and says why.
- [ ] `flai dashboard check` and `flai dashboard upgrade` compare the running container's digest with the list's and upgrade by digest; `flai dashboard status` shows the digest and whether it is on a signed list.
- [ ] The verified list is cached beside `flai serve`'s state with its release version, for the dial-time check to read.
- [ ] `--image` naming another registry or `--build` run what they name, as today, and are reported as unsigned; refusing or allowing them is the next story's.
- [ ] Tests cover resolution, verification, refusal, and the cache; `docs/operators/settings.md` and `docs/users/flai.md` describe the new meaning of `dashboard.tag`.
- [ ] `flai dashboard versions`, with `--json`, marks each dashboard release with no digest list this binary can verify, such as one from before signing, and the Updates page shows it as not deployable and offers no way to deploy it.

## Tasks
- T-1239 selfupgrade lists the flaiover GitHub releases with their assets and downloads a release's asset
- T-1240 flai/internal/dashboard parses a flaiover digest list, verifies its signature with the embedded key, and caches it in flai serve's state by release version
- T-1241 flai/internal/dashboard resolves dashboard.tag to a flaiover release and its verified image digest, and refuses a release with no list it can verify
- T-1242 flai dashboard versions marks each dashboard release with no digest list this binary can verify, in its table and its JSON
- T-1243 flai dashboard runs the resolved image by digest, refuses to start without a verified list, reports --image and --build as unsigned, and status shows the digest
- T-1244 flai dashboard check and upgrade compare the running container's digest with the verified list and upgrade by digest
- T-1245 The Updates page shows a dashboard release with no verifiable digest list as not deployable and offers no way to deploy it
- T-1246 The user, operator, and design documentation say that dashboard.tag names a signed flaiover release run by digest, and what a refusal means

## Notes

### Planning

Tasks, in five layers:

1. T-1239, the flaiover release listing and asset download in `selfupgrade`; T-1240, the digest list's parsing, verification, and cache in the new `flai/internal/dashboard`.
2. T-1241, the resolver, after both; T-1242, `flai dashboard versions`' deployable mark and the `--published` refusal, after both.
3. T-1243, `flai dashboard`, `status`, and the restart record by digest, after T-1241; T-1245, the Updates page, after T-1242.
4. T-1244, `check` and `upgrade` by digest, after T-1243, since both change `dashboard.go` and `dashboard_test.go`.
5. T-1246, the documentation, after T-1244 and T-1245.

Touches:

- Declared: `flai/cmd/dashboard.go`, `flai/internal/selfupgrade`, `flai/internal/dashboard`, `docs/operators/settings.md`, `docs/users/flai.md`.
- Layout: `flai/cmd/dashboard_upgrade.go`, which holds `check` and `upgrade`; `flai/cmd/dashboard_versions.go`, which lists dashboard releases and holds `requirePublishedDashboard` (ADR-0117); `flai/cmd/dashboard_watch.go`, whose record keeps the chosen release through a restart (ADR-0118), and its test; `flai/cmd/dashboard_test.go`, with `fakeRunner`.
- Layout: `flai/internal/selfupgrade/selfupgrade.go` and its test, the release client, which lists only flai releases with assets and flaiover tags without them; `flai/internal/dashboard/digests.go`, `cache.go`, and `resolve.go`, each with its test, predicted new files of the new package.
- Design: `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md` § Deploying a chosen release, and `design/tech/docker.md` § Run, which say how the image is chosen and run; `docs/operators/runbooks/update.md`, which says `latest` follows the main branch, no longer true here.
- Co-change: `docs/users/flai-reference.md`, at 24%, generated from the help text the commands may change.
- Criterion 6, added on TH-0313: `flai/cmd/dashboard_versions_test.go`; `flaiover/src/lib/components/HostPanel.svelte` and its test, which list and deploy dashboard releases on the Updates page; `docs/users/flaiover.md`, which describes that page.
- Not added: `flai/internal/hostapi/writes.go` and `flaiover/src/routes/api/dashboard/+server.ts`, which pass `flai dashboard versions --json` through unchanged; `flai/cmd/host.go`, since no criterion changes `flai host status`.
- Folder touches kept, as declared: `flai/internal/selfupgrade` and `flai/internal/dashboard`. The tasks name every file below them, which replace them in the story's claim (ADR-0096).

Forecast: 1h30m, raised from 1h15m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gives 58m, from a median of 104 s per unit of size over 33 large-band feature stories, times a size of 33.
- Done feature stories with 4 to 7 criteria took a median of about 1h of agent time. S-0233, of like scope with five tasks, is forecast at 1h15m. This one has eight tasks in five layers, a new package, five dashboard subcommands, a Svelte page, and seven documents, so 1h30m.
- Delivery: flai played this story out at 02:01 with its own 58m; 2026-10-08T02:33Z adds the 32m more this forecast gives.

Cost of delay: 3.95 USD a week, as `flai cod` gives it after the forecast changed, up from 3.29. It is this story's 1h30m share of the 9h30m forecast over E-0015's eight open stories, applied to the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure closes only when the chain is done, so a share by work fits.

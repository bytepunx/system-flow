---
id: S-0236
type: story
nature: feature
title: flai dashboard resolves the image from the signed digest list and runs it by digest
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-07T22:10:05Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/cmd/dashboard.go, flai/internal/selfupgrade, flai/internal/dashboard, docs/operators/settings.md, docs/users/flai.md, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_versions.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_test.go, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/docker.md, flai/cmd/dashboard_versions_test.go, flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, docs/users/flaiover.md]
after: [S-0233, S-0234]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 1h15m
  delivery: 2026-10-08T06:08:00Z
  basis: "Its own forecast of 1h15m; 6th in the pull order with an in-progress limit of 3, behind S-0333, S-0332, S-0232, S-0233, S-0234 and S-0235."
  by: flai
  at: 2026-10-07T22:10:05Z
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

## Notes

### Planning

Touches:

- Declared: `flai/cmd/dashboard.go`, `flai/internal/selfupgrade`, `flai/internal/dashboard`, `docs/operators/settings.md`, `docs/users/flai.md`.
- Layout: `flai/cmd/dashboard_upgrade.go`, which holds `check` and `upgrade`; `flai/cmd/dashboard_versions.go`, which lists dashboard releases (ADR-0117); `flai/cmd/dashboard_watch.go`, whose restart keeps the chosen release (ADR-0118); `flai/cmd/dashboard_test.go`.
- Design: `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md` § Deploying a chosen release, and `design/tech/docker.md` § Run, which say how the image is chosen and run.
- Criterion 6, added on TH-0313: `flai/cmd/dashboard_versions_test.go`; `flaiover/src/lib/components/HostPanel.svelte` and its test, which list and deploy dashboard releases on the Updates page; `docs/users/flaiover.md`, which describes that page.
- Folder touches kept, as declared: `flai/internal/selfupgrade`, which may gain a digest-list file; `flai/internal/dashboard`, a package that does not exist yet, for the resolver and its cache.

Forecast: 1h15m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 20m from 116 s per unit of size, over only 3 medium-band feature stories.
- Done feature stories of this size took a median of about 1h of agent time. This one adds resolution, verification, a cache, and four subcommands that must keep ADR-0118's chosen release, so 1h15m.
- The first delivery was played out after S-0233 and S-0234 at flai's cycle factor of 6.85.

Cost of delay: no value yet. E-0015 and its stories have no inputs; TH-0312 asks the operator for them.

---
id: S-0239
type: story
nature: feature
title: dashboard.allow_unsigned lets a development build connect, shown on every page and in every status
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-07T21:00:01Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [flai/internal/config, flai/internal/manifest, flai/cmd/dashboard.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/components, docs/operators/settings.md, docs/users/flaiover.md, flaiover/src/lib/server/agent.test.ts, flai/cmd/serve.go, flai/cmd/host.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/contributors/index.md, template/root/docs/operators/index.md.tmpl, template/CHANGELOG.md, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flaiover/src/routes/+layout.svelte, flaiover/src/routes/api/projects/+server.ts, flaiover/src/routes/api/projects/projects.test.ts]
after: [S-0237, S-0238]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 1h30m
  delivery: 2026-10-08T08:27:00Z
  basis: "Its own forecast of 1h30m; 8th in the pull order with an in-progress limit of 3, behind S-0232, S-0233, S-0234, S-0235, S-0236, S-0237 and S-0238."
  by: flai
  at: 2026-10-07T21:00:01Z
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

## Notes

### Planning

Touches:

- Declared: `flai/internal/config`, `flai/internal/manifest`, `flai/cmd/dashboard.go`, `flai/internal/serve`, `flaiover/src/lib/server/agent.ts`, `flaiover/src/routes`, `flaiover/src/lib/components`, `docs/operators/settings.md`, `docs/users/flaiover.md`.
- Layout: `flaiover/src/lib/server/agent.test.ts`; `flai/cmd/serve.go` and `flai/cmd/host.go` for the statuses; `docs/users/flai.md` and the generated `docs/users/flai-reference.md` for `--allow-unsigned`.
- Design: `design/system/project-manifest.md` for the manifest key; `design/system/flai-cli.md`.
- Criteria: `docs/contributors/index.md`, this repository's contributor documentation (criterion 4); `template/root/docs/operators/index.md.tmpl` and `template/CHANGELOG.md`, the template's operator documentation and its changelog.
- Folder touches kept, as declared: `flai/internal/config` and `flai/internal/manifest`, one file and its test each; `flaiover/src/lib/components`, where the banner may be a new component no task can name yet.
- Narrowed on TH-0313: `flai/internal/serve` to `flai/internal/serve/serve.go`, where the allowance changes what is dialled and accepted, and its test; `flaiover/src/routes` to `flaiover/src/routes/+layout.svelte`, which shows the banner on every page, and `flaiover/src/routes/api/projects/+server.ts` with its test, the per-project connection list.

Forecast: 1h30m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 27m from 114 s per unit of size, over 29 large-band feature stories.
- Done feature stories of this size took a median of about 1h of agent time. This one carries a setting through the configuration, the manifest, both components, a banner on every page, and three statuses, so 1h30m.
- The first delivery was played out after S-0238, the later of its two afters, at flai's cycle factor of 6.85.

Cost of delay: no value yet. E-0015 and its stories have no inputs; TH-0312 asks the operator for them.

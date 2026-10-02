---
id: S-0239
type: story
nature: feature
title: dashboard.allow_unsigned lets a development build connect, shown on every page and in every status
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-02T12:37:24Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [flai/internal/config, flai/internal/manifest, flai/cmd/dashboard.go, flai/internal/serve, flaiover/src/lib/server/agent.ts, flaiover/src/routes, flaiover/src/lib/components, docs/operators/settings.md, docs/users/flaiover.md]
after: [S-0237, S-0238]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

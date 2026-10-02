---
id: S-0236
type: story
nature: feature
title: flai dashboard resolves the image from the signed digest list and runs it by digest
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-02T12:37:24Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/cmd/dashboard.go, flai/internal/selfupgrade, flai/internal/dashboard, docs/operators/settings.md, docs/users/flai.md]
after: [S-0233, S-0234]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

## Tasks

## Notes

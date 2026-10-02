---
id: S-0238
type: story
nature: feature
title: flai measures the container's image through Docker before it dials and refuses one no signed list names
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-02T12:37:24Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/internal/serve, flai/cmd/serve.go, flai/cmd/dashboard.go, docs/operators, docs/users/flai.md]
after: [S-0236]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0238 flai measures the container's image through Docker before it dials and refuses one no signed list names

## Goal

Before `flai serve` dials the dashboard, and whenever the connection is opened anew, it asks Docker which image the `flaiover` container runs and compares its digest with the verified digest lists it has (ADR-0070). A digest on a list is dialled; one on none is refused, logged, and shown, so that a container restarted from another image or a moved tag is caught by something the peer does not control.

## Acceptance criteria
- [ ] `flai serve` reads the container's image digest with `docker inspect` and `docker image inspect` through the `execx.Runner`, compares it with the cached verified lists, fetching and verifying the list for the running version when it has none, and dials only when the digest is on one.
- [ ] An unlisted digest is not dialled: one `error` event with the image and digest, the state in `flai serve status`, `flai dashboard status`, and `flai host status`, and another look when the container changes.
- [ ] A dashboard flai cannot inspect, because Docker is not where `flai serve` runs or the container is not local, falls back to the stamp check alone and says so in status.
- [ ] Tests cover a listed digest, an unlisted one, a local image with no `RepoDigests`, and an uninspectable dashboard, with a fake runner.
- [ ] `docs/users/flai.md` and the operator documentation describe the check and what status shows.

## Tasks

## Notes

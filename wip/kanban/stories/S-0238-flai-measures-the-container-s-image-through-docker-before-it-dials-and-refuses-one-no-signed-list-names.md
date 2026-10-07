---
id: S-0238
type: story
nature: feature
title: flai measures the container's image through Docker before it dials and refuses one no signed list names
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-07T22:25:33Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/cmd/serve.go, flai/cmd/dashboard.go, docs/operators, docs/users/flai.md, flai/cmd/host.go, design/system/flai-cli.md, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/serve/imagecheck.go, flai/internal/serve/imagecheck_test.go]
after: [S-0236]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 2.63
  by: planner-E-0015
  at: 2026-10-07T22:13:46Z
forecast:
  duration: 1h
  delivery: 2026-10-08T07:27:00Z
  basis: "Its own forecast of 1h; 7th in the pull order with an in-progress limit of 3, behind S-0232, S-0333, S-0332, S-0233, S-0234, S-0235, S-0236 and S-0237."
  by: flai
  at: 2026-10-07T22:25:33Z
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

### Planning

Touches:

- Declared: `flai/internal/serve`, `flai/cmd/serve.go`, `flai/cmd/dashboard.go`, `docs/operators`, `docs/users/flai.md`.
- Layout: `flai/cmd/host.go`, which prints `flai host status`.
- Co-change and design: `design/system/flai-cli.md`, which describes the three statuses.
- Folder touch kept, as declared: `docs/operators`, whose runbooks may gain a page on the check.
- Narrowed on TH-0313: `flai/internal/serve` to `flai/internal/serve/serve.go`, which builds the channel client and would run the check before it dials, and its test; and `flai/internal/serve/imagecheck.go` with its test, a predicted new file for the check, which the story's agent may name otherwise.

Forecast: 1h. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 20m from 116 s per unit of size, over only 3 medium-band feature stories.
- Done feature stories of this size took a median of about 1h of agent time. This one is a Docker check before dialling, three statuses, and a fake runner, so 1h.
- The first delivery was played out after S-0236 at flai's cycle factor of 6.85.

Cost of delay: 2.63 USD a week, as `flai cod` gives it: this story's 1h share of the 9h30m forecast over E-0015's eight open stories, of the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure is closed only when the chain is done, so a share by work fits.

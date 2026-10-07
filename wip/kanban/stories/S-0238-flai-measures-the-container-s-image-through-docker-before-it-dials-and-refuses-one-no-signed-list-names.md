---
id: S-0238
type: story
nature: feature
title: flai measures the container's image through Docker before it dials and refuses one no signed list names
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-07T23:45:02Z
transitions: []
tags: [cli]
topics: [release, security]
touches: [flai/cmd/serve.go, flai/cmd/dashboard.go, docs/operators, docs/users/flai.md, flai/cmd/host.go, design/system/flai-cli.md, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/serve/imagecheck.go, flai/internal/serve/imagecheck_test.go, flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flai/cmd/serve_test.go, flai/cmd/host_test.go, flai/cmd/dashboard_test.go, docs/operators/index.md, docs/users/flai-reference.md, design/system/release-signing.md]
after: [S-0236]
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
      seconds: 294
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 614
          cache_read: 4573334
          cache_write: 22537
          cost: 1.1325
cost_of_delay:
  value: 3.21
  by: planner-S-0238
  at: 2026-10-07T23:06:46Z
forecast:
  duration: 1h15m
  delivery: 2026-10-08T04:27:00Z
  basis: "Its own forecast of 1h15m; 6th in the pull order with an in-progress limit of 3, behind S-0232, S-0310, S-0314, S-0233, S-0234, S-0235, S-0236 and S-0237."
  by: flai
  at: 2026-10-07T23:45:02Z
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
- T-1253 flai/internal/serve measures the flaiover container's image digest through Docker and judges it against the verified digest lists
- T-1254 The channel client asks a gate before each dial, and a refusal skips the dial, records why, and waits a minute
- T-1255 flai serve runs the image check before each dial, refuses an unlisted image, and records the verdict in its status
- T-1256 flai serve status and flai host status say whether the dashboard's image is signed, unsigned, or not measured
- T-1257 flai dashboard status says whether flai serve dials the running image or refused it as unsigned
- T-1258 The user, operator, and design documentation describe the dial-time image check and what each status shows

## Notes

### Planning

Tasks, in four layers:

| Layer | Tasks | Why |
|-------|-------|-----|
| 1 | T-1253 (the image check), T-1254 (a gate before each dial in the channel client) | No path in common; the check reads S-0236's `flai/internal/dashboard`, on main first |
| 2 | T-1255 (serve runs the check through the gate) after T-1253 and T-1254 | It joins the two |
| 3 | T-1256 (`flai serve status`, `flai host status`) and T-1257 (`flai dashboard status`) after T-1255 | Each shows the verdict T-1255 records; they share no path |
| 4 | T-1258 (docs) after T-1256 and T-1257 | Documents the words the statuses print |

Touches:

- Declared, all kept: `flai/cmd/serve.go`, `flai/cmd/dashboard.go`, `docs/operators`, `docs/users/flai.md`, `flai/cmd/host.go`, `design/system/flai-cli.md`, `flai/internal/serve/serve.go` and its test, `flai/internal/serve/imagecheck.go` and its test (narrowed from `flai/internal/serve` on TH-0313).
- Layout, added by this run: `flai/internal/channel/channel.go` and its test. `channel.Client.Run` dials and reconnects inside the channel package, so a check before each dial, reconnections included, needs a hook there. S-0237 changes the same file; while either story is in progress the other is held, which is right, since both change the dial loop's backoff.
- Layout, added by this run: `flai/cmd/serve_test.go`, `flai/cmd/host_test.go`, and `flai/cmd/dashboard_test.go`, the tests of the three statuses; `host_test.go` is new unless S-0237 creates it first.
- Design, added by this run: `design/system/release-signing.md`, whose § flai measures the container before it dials leaves to this story which digest `RepoDigests` records; `docs/operators/index.md`, § The connection from flai on the host or § Security posture.
- Co-change: `docs/users/flai-reference.md` (27%), generated from the help `flai dashboard status` changes. Not added: `design/system/flaiover-dashboard.md` (24%) and `docs/users/flaiover.md` (19%), since the dashboard does not change; `flai/internal/hostapi/writes.go` (8%) and `flai/cmd/serve_actions.go` (6%), since no host method or action changes.
- Folder touch kept, as declared: `docs/operators`, where T-1258 names `index.md` and the story's agent may add a runbook page on a refused image. In the claim it narrows to the file its task names (ADR-0096).

Forecast: 1h15m, raised from 1h. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gives 40m after this run's touches, from a median of 104 s per unit of size over 33 large-band feature stories, times size 23.
- S-0237 has six tasks in three layers with code on both sides and is planned at 1h30m. This story has six tasks in four layers, on flai alone: a Docker check, a gate in the channel client, serve's wiring, three statuses, and docs. So 1h15m.
- Delivery: flai played this story out at 03:26 with its own 40m; 2026-10-08T04:01Z adds the 35m more this forecast gives.

Cost of delay: 3.21 USD a week, as `flai cod` gives it after the forecast changed, up from 2.63. It is this story's 1h15m share of the 9h45m forecast over E-0015's eight open stories, applied to the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure closes only when the chain is done, so a share by work fits.

---
id: T-1265
type: task
nature: feature
title: "flai serve status, flai dashboard status, and flai host status say \"unsigned allowed\" with the side"
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:46Z
updated: 2026-10-07T23:09:46Z
transitions: []
stream: S-0239
tags: [cli]
touches: [flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go]
after: [T-1261, T-1262]
---
# T-1265 flai serve status, flai dashboard status, and flai host status say "unsigned allowed" with the side

## Work

The statuses of criterion 3. It waits for T-1262, which records the unsigned side in `serve.Status.Connections`, and for T-1261, which changes `flai/cmd/dashboard.go` and its test first.

- `flai/cmd/serve.go`: `flai serve status` says, per project whose connection has an unsigned side, "unsigned allowed" and the side (`dashboard`, `flai`, or both), beside S-0237's refusal and S-0238's image verdict. `--json` carries the side.
- `flai/cmd/host.go`: `flai host status` says the same per project.
- `flai/cmd/dashboard.go`: `flai dashboard status` says the same for the running dashboard, and says when the container was started with `FLAIOVER_ALLOW_UNSIGNED=1`.
- A signed pair, with the setting on or off, shows nothing of it in any of the three.
- Tests in `serve_test.go`, `host_test.go`, and `dashboard_test.go`: each side, both sides, and a signed pair with the setting on.

## Done when

- The three statuses say "unsigned allowed" with the side for a connection allowed unsigned, and nothing for a signed pair.
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. If the status help text changes, regenerate `docs/users/flai-reference.md` with `make flai-reference` and widen the touches to it.

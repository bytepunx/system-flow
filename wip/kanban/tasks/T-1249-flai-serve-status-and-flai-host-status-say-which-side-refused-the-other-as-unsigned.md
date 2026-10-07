---
id: T-1249
type: task
nature: feature
title: flai serve status and flai host status say which side refused the other as unsigned
status: backlog
parent: S-0237
owner: alex
created: 2026-10-07T23:00:39Z
updated: 2026-10-07T23:00:39Z
transitions: []
stream: S-0237
tags: [cli]
touches: [flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/host.go, flai/cmd/host_test.go]
after: [T-1247]
---
# T-1249 flai serve status and flai host status say which side refused the other as unsigned

## Work

The status part of criterion 2. It waits for T-1247, which adds the refusal to `channel.State`; this task prints it.

- `flai serve status` (`printServeStatus` in `flai/cmd/serve.go`): for a project whose last attempt was refused, print "the dashboard refused this flai: unsigned" or "flai refused the dashboard: unsigned", with the reason and the peer's version, in place of the bare last error. Do the same for the projects served from the folder and from the import folders. `--json` carries the fields `channel.State` gained, unchanged.
- `flai host status` (`describeHost` in `flai/cmd/host.go`): `host.Status` has no channel state today. Read `flai serve`'s state, as `readServeStatus` does, and print each refused project under the `serve` child in the same words. Leave `flai/internal/host` alone unless the state cannot be read from the command.
- `flai dashboard status` shows the part of `flai serve status` about its own project. Check that it shows the refusal too, and note in the narrative if it needs a line of its own.
- Tests: `flai/cmd/serve_test.go` covers both refusals in the plain and JSON output. A new `flai/cmd/host_test.go` covers `describeHost` printing them; no test of `flai host status` exists today.

## Done when

- `flai serve status` and `flai host status` print "the dashboard refused this flai: unsigned" and "flai refused the dashboard: unsigned" for a project refused each way, and nothing new for a connected one.
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. If the story's agent finds that `host.Status` should carry the connections itself, it widens this task's touches to `flai/internal/host/host.go` and says why in the narrative's Decisions.

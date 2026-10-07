---
id: T-1255
type: task
nature: feature
title: flai serve runs the image check before each dial, refuses an unlisted image, and records the verdict in its status
status: backlog
parent: S-0238
owner: alex
created: 2026-10-07T23:05:48Z
updated: 2026-10-07T23:05:48Z
transitions: []
stream: S-0238
tags: [cli]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
after: [T-1253, T-1254]
---
# T-1255 flai serve runs the image check before each dial, refuses an unlisted image, and records the verdict in its status

## Work

Connect T-1253's check to T-1254's gate in `flai/internal/serve/serve.go`. This task waits for both.

- Add a Docker runner to `serve.Options`, `execx.System{}` when nil, as `Git` is. Tests then give a fake one.
- The default `NewClient` gives each client a gate that runs the image check for the client's URL:
  - listed: the gate passes.
  - unlisted: the gate refuses with the image and the digest.
  - not measured: the gate passes, and the stamp check alone stands.
- Record each project's last verdict where `flai serve status` reads it, such as a field in `Status` or in its connection's `channel.State`: listed, unlisted with the image and digest, or not measured with the reason.
- While a verdict is unlisted, look again when the container's ID or image ID changes, and at least at every retry.
- A connection the dashboard ends, such as on a container restart, is checked again before it is dialled again. The gate does this.

## Done when

- `serve_test.go` covers four cases with a fake runner: a listed image is dialled, an unlisted one is not dialled and its verdict is in the status, an uninspectable dashboard is dialled with "not measured" in the status, and a container changed to a listed image is dialled on the next look.
- `flai test flai/internal/serve` passes.

## Notes

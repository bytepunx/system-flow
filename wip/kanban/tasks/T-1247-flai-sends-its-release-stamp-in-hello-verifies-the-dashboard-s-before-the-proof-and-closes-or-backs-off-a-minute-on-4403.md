---
id: T-1247
type: task
nature: feature
title: flai sends its release stamp in hello, verifies the dashboard's before the proof, and closes or backs off a minute on 4403
status: backlog
parent: S-0237
owner: alex
created: 2026-10-07T23:00:18Z
updated: 2026-10-07T23:00:18Z
transitions: []
stream: S-0237
tags: [cli]
touches: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flai/internal/channel/channeltest/dashboard.go, flai/internal/buildinfo/stamp.go, flai/internal/buildinfo/stamp_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
---
# T-1247 flai sends its release stamp in hello, verifies the dashboard's before the proof, and closes or backs off a minute on 4403

## Work

The flai side of criteria 1 and 4, and the back-off and log of criterion 2. It waits for no task of this story: the stamp in `buildinfo` and its verification are S-0235's (`flai/internal/buildinfo/stamp.go`), and the public key is S-0232's, both on main before this story starts.

- In `flai/internal/buildinfo/stamp.go`, add a check of a peer's stamp beside the one S-0235 wrote for flai's own: the signature with the given public key, and the component the statement names, without comparing its version or commit with flai's. Tests in `stamp_test.go`: a valid peer stamp, another key, a wrong component.
- In `flai/internal/channel/channel.go`, bump `Protocol` to 2. `hello` carries `release: {statement, signature}` from the client's own stamp, empty for a development build. Read `release` from the dashboard's answer and verify it, component `flaiover`, before the credential proof. A missing, unverifiable, or wrong-component stamp, including an answer from a protocol 1 dashboard without the field, closes the socket with a new `CloseUnsigned` (4403, reason `unsigned peer`).
- Read 4403 from the dashboard as a refusal of this flai. Back off `HeldBackoff` on 4403 from either side, as on 4409.
- Record the refusal in `State`: who refused (`dashboard` or `flai`), the reason (`unsigned`, `invalid signature`, `wrong component`), and the peer's version when it named one. The status task reads these fields; it waits for this task for that reason.
- Log one `error` per refusal, with the peer's version and the reason, not one per retry.
- Take the stamp and the verifier as fields of `Client`, so that tests sign with a key pair they generate. `flai/internal/serve/serve.go` sets them from `buildinfo` and the embedded key when it builds each client.
- `flai/internal/channel/channeltest/dashboard.go` answers with a stamp the test gives, none, or one of another component, and can close with 4403.
- Tests in `channel_test.go`: a signed dashboard connects; an unsigned one, a wrong component, and a protocol 1 answer without the field each end in 4403 with the reason in `State`; a 4403 from the dashboard backs off a minute. `flai/internal/serve/serve_test.go` checks that the client gets the stamp and the key.

## Done when

- `hello` carries flai's stamp, `Protocol` is 2, and flai closes with 4403 on a missing, invalid, or wrong-component stamp in the dashboard's answer, before it sends `hello.prove`.
- flai backs off a minute on 4403 from either side, logs one `error` with the peer's version and the reason, and `State` records who refused and why.
- The channel, stamp, and serve tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. The field names in `State` and the exact reason words are the story's agent's to choose; the status task prints them.

---
id: T-1252
type: task
nature: feature
title: The operator, user, and design documentation describe the stamp exchange, 4403, its promise and limit, and the connection list
status: backlog
parent: S-0237
owner: alex
created: 2026-10-07T23:01:03Z
updated: 2026-10-07T23:01:03Z
transitions: []
stream: S-0237
tags: [cli, dashboard]
touches: [docs/operators/index.md, docs/users/flaiover.md, docs/users/flai.md, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, design/system/flai-cli.md, design/system/release-signing.md]
after: [T-1249, T-1250, T-1251]
---
# T-1252 The operator, user, and design documentation describe the stamp exchange, 4403, its promise and limit, and the connection list

## Work

Criterion 5, and the design kept true. It waits for the status and dashboard tasks, so that it documents the exchange, the words, and the fields as built. Those tasks waited for the two protocol tasks, so everything is done by then.

- `docs/operators/index.md` § The connection from flai on the host: add the stamp exchange to the credential bullet, and close code 4403 beside 4409 with the minute's back-off. State the promise and its limit in the words of `design/system/release-signing.md` § Verifying the peer: it names the peer's release and refuses builds that are not one, and it does not defend against someone who holds the credential and a release binary. Say what to do when either side is refused.
- `docs/users/flaiover.md` § Host flai and § More than one project: the banner for a refused flai, and what the connection list shows per project (signed, unsigned, refused, with the reason).
- `docs/users/flai.md`: the new lines of `flai serve status` and `flai host status`.
- `design/system/flaiover-dashboard.md` § The channel to flai on the host: the handshake with `release` each way, protocol 2, 4403 `CLOSE_UNSIGNED`, and the per-project refusal and signed mark.
- `design/system/flai-cli.md`, the `flai serve` row: the back-off on 4403 and the status words. `design/system/dashboard-host-channel.md`: one entry for S-0237 in its list of what landed. `design/system/release-signing.md` § A release stamp in `hello`: what was built, where it differs from the design.
- Bump `updated` on each document changed.

## Done when

- `docs/operators/index.md` states the promise and its limit as `release-signing.md` § Verifying the peer does, and names 4403.
- `docs/users/flaiover.md` describes what the connection list shows, and `docs/users/flai.md` the status lines.
- The design documents describe the handshake as built, and the markdown lint passes with `flai test` on the changed paths.

## Notes

Drafted by the planner. S-0235's docs task also edits `docs/operators/index.md`, and is on main before this story starts.

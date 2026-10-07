---
id: T-1258
type: task
nature: feature
title: The user, operator, and design documentation describe the dial-time image check and what each status shows
status: backlog
parent: S-0238
owner: alex
created: 2026-10-07T23:06:10Z
updated: 2026-10-07T23:06:10Z
transitions: []
stream: S-0238
tags: [cli]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, design/system/flai-cli.md, design/system/release-signing.md]
after: [T-1256, T-1257]
---
# T-1258 The user, operator, and design documentation describe the dial-time image check and what each status shows

## Work

Document the check as built. This task waits for T-1256 and T-1257, so it describes the words the three statuses print.

- `docs/users/flai.md`: where it describes `flai serve status`, `flai dashboard status`, and `flai host status`, say what the image line means. Cover signed, "the dashboard runs an unsigned image", and "image not measured", and what to do for each.
- `docs/users/flai-reference.md`: regenerate it with `make flai-reference` after T-1257 changes the help.
- `docs/operators/index.md`, under § The connection from flai on the host or § Security posture, cover three points:
  - What the check measures: Docker, not the peer, says what the container runs.
  - When it runs: before each dial and on each reconnection.
  - When it falls back to the stamp alone: no Docker where `flai serve` runs, or a dashboard on another host.
- `design/system/flai-cli.md`: the three statuses' new line, as built.
- `design/system/release-signing.md` § flai measures the container before it dials: settle its **to check** with what T-1253 found about which digest `RepoDigests` records.

## Done when

- Each document above says what the command prints, and `updated` is bumped where front matter carries it.
- `flai test` on the changed documents passes the markdown lint, and the generated reference matches the help.

## Notes

---
id: MS-0071
title: "S-0233 holds S-0237 on overlap alone: may S-0237 start beside it?"
from: S-0237
to: S-0233
about: [docs/users/flaiover.md, flai/cmd/host.go]
status: open
participants: [flai, agent-S-0233, agent-S-0237]
created: 2026-10-12T16:53:40Z
updated: 2026-10-12T17:16:33Z
shares:
  - holder: S-0233
    held: S-0237
    paths: [docs/users/flaiover.md, flai/cmd/host.go]
    split: "In flai/cmd/host.go, S-0233 changes only `flai host upgrade` and `flai host versions`: their help text and printHostVersions, which marks a release whose signature this flai cannot verify. S-0237 changes the rest, such as the host's dial and its hello. In docs/users/flaiover.md, S-0233 changes only the Updates page's list of flai releases; S-0237 changes what the page says about the release stamp and an unsigned peer."
    by: agent-S-0233
    at: 2026-10-12T16:54:20Z
  - holder: S-0233
    held: S-0237
    paths: [design/system/release-signing.md]
    split: "S-0233 changes only the sections on verifying a release before it is installed (Today › How a release is verified before it is installed, and what Signing releases says of self-upgrade and install.sh). S-0237 changes the Verifying the peer sections."
    by: agent-S-0233
    at: 2026-10-12T17:04:07Z
---

# MS-0071 S-0233 holds S-0237 on overlap alone: may S-0237 start beside it?

Between S-0237 and S-0233, about `docs/users/flaiover.md`, `flai/cmd/host.go`.

## Entries

### 2026-10-12T16:53:40Z flai S-0237
S-0233 holds S-0237 on overlap alone: may S-0237 start beside it?

S-0237, "flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403", is ready, and its claim overlaps S-0233's on `docs/users/flaiover.md`, `flai/cmd/host.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The channel's handshake carries each side's release stamp; each verifies the other's with the public key and the component it names before the credential proof, and a missing, unverifiable, or wrong stamp closes the connection with `4403 unsigned peer` (ADR-0070). The refusal is visible on both sides.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0233 --remove <path>`, or on the task that names the path, if S-0233 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0071 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0237 can start.
3. Reply with `message_reply` saying why the hold stands: S-0237 waits until S-0233 moves to review, is cancelled, or is sent back.

### 2026-10-12T16:54:20Z agent-S-0233 S-0233
S-0233 shares `docs/users/flaiover.md`, `flai/cmd/host.go` with S-0237, split so:

> In flai/cmd/host.go, S-0233 changes only `flai host upgrade` and `flai host versions`: their help text and printHostVersions, which marks a release whose signature this flai cannot verify. S-0237 changes the rest, such as the host's dial and its hello. In docs/users/flaiover.md, S-0233 changes only the Updates page's list of flai releases; S-0237 changes what the page says about the release stamp and an unsigned peer.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-12T17:04:07Z agent-S-0233 S-0233
S-0233 shares `design/system/release-signing.md` with S-0237, split so:

> S-0233 changes only the sections on verifying a release before it is installed (Today › How a release is verified before it is installed, and what Signing releases says of self-upgrade and install.sh). S-0237 changes the Verifying the peer sections.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-12T17:04:35Z agent-S-0237 S-0237
S-0237 keeps to both splits. In flai/cmd/host.go it changes only `flai host status` and the host's dial/hello, never `host upgrade`, `host versions`, or printHostVersions. In docs/users/flaiover.md it changes only what the page says about the release stamp, an unsigned peer, and the connection list, not the Updates page's list. In design/system/release-signing.md it changes only the Verifying the peer sections.

### 2026-10-12T17:16:33Z agent-S-0233 S-0233
T-1228 of S-0233 changed paths S-0237's claim covers.

T-1228, flai self-upgrade --list and flai host versions mark each release this binary cannot verify, and flai host upgrade reports the refusal, committed 7c66e9b on story/S-0233, `feat: [S-0233] self-upgrade --list and host versions mark releases this flai cannot verify`, changing `flai/cmd/host.go`. It reaches the main branch when S-0233 is accepted; `git show 7c66e9b` shows it until then. Reply here if it breaks your work, or adjust to it early.

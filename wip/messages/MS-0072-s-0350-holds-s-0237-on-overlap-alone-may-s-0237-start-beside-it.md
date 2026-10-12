---
id: MS-0072
title: "S-0350 holds S-0237 on overlap alone: may S-0237 start beside it?"
from: S-0237
to: S-0350
about: [docs/operators/index.md, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [flai, agent-S-0350, agent-S-0237, orchestrator]
created: 2026-10-12T16:53:40Z
updated: 2026-10-12T17:19:46Z
shares:
  - holder: S-0350
    held: S-0237
    paths: [docs/operators/index.md]
    split: S-0350 adds one or two sentences to docs/operators/index.md saying a provider's key variable must be set in the environment flai host or flai serve runs in, near what it says about agents and providers. S-0237 changes what it says about release stamps and the 4403 refusal. Neither edits the other's lines.
    by: agent-S-0350
    at: 2026-10-12T16:53:57Z
---

# MS-0072 S-0350 holds S-0237 on overlap alone: may S-0237 start beside it?

Between S-0237 and S-0350, about `docs/operators/index.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-12T16:53:40Z flai S-0237
S-0350 holds S-0237 on overlap alone: may S-0237 start beside it?

S-0237, "flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403", is ready, and its claim overlaps S-0350's on `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The channel's handshake carries each side's release stamp; each verifies the other's with the public key and the component it names before the credential proof, and a missing, unverifiable, or wrong stamp closes the connection with `4403 unsigned peer` (ADR-0070). The refusal is visible on both sides.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0350 --remove <path>`, or on the task that names the path, if S-0350 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0072 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0237 can start.
3. Reply with `message_reply` saying why the hold stands: S-0237 waits until S-0350 moves to review, is cancelled, or is sent back.

### 2026-10-12T16:53:57Z agent-S-0350 S-0350
S-0350 shares `docs/operators/index.md` with S-0237, split so:

> S-0350 adds one or two sentences to docs/operators/index.md saying a provider's key variable must be set in the environment flai host or flai serve runs in, near what it says about agents and providers. S-0237 changes what it says about release stamps and the 4403 refusal. Neither edits the other's lines.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-12T17:04:36Z agent-S-0237 S-0237
S-0237 keeps to the split: in docs/operators/index.md it changes only what it says about release stamps and the 4403 refusal, and leaves S-0350's provider-key sentences alone.

### 2026-10-12T17:06:50Z agent-S-0350 S-0350
T-1380 of S-0350 changed paths S-0237's claim covers.

T-1380, The user, operator, and design documents say what a provider sets at start, what it refuses, and where its key's variable must be, committed 9b1766c on story/S-0350, `docs: [S-0350] say what a provider sets at start, what it refuses, and where its key's variable must be`, changing `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`. It reaches the main branch when S-0350 is accepted; `git show 9b1766c` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T17:19:46Z orchestrator
Closed: S-0350 was accepted

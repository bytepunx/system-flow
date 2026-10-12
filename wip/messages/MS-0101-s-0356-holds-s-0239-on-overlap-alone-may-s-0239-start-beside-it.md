---
id: MS-0101
title: "S-0356 holds S-0239 on overlap alone: may S-0239 start beside it?"
from: S-0239
to: S-0356
about: [docs/contributors/index.md, docs/operators/index.md]
status: open
participants: [flai, agent-S-0356]
created: 2026-10-12T20:52:11Z
updated: 2026-10-12T21:06:13Z
shares:
  - holder: S-0356
    held: S-0239
    paths: [docs/contributors/index.md, docs/operators/index.md]
    split: "S-0356 changes only the § Gateway checks section of docs/contributors/index.md (its LiteLLM rows and run example) and, in docs/operators/index.md, adds only a LiteLLM provider entry beside the bullet on a provider's key under the agent/provider material. S-0239 owns everything else in both files, including any dashboard.allow_unsigned text and the sections on the dashboard image and the connection from flai on the host. Where both edit near the same lines, the second to sync keeps both additions."
    by: agent-S-0356
    at: 2026-10-12T21:06:13Z
---

# MS-0101 S-0356 holds S-0239 on overlap alone: may S-0239 start beside it?

Between S-0239 and S-0356, about `docs/contributors/index.md`, `docs/operators/index.md`.

## Entries

### 2026-10-12T20:52:11Z flai S-0239
S-0356 holds S-0239 on overlap alone: may S-0239 start beside it?

S-0239, "dashboard.allow_unsigned lets a development build connect, shown on every page and in every status", is ready, and its claim overlaps S-0356's on `docs/contributors/index.md`, `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Unsigned builds are allowed only on purpose: `dashboard.allow_unsigned` in the host's flai configuration, settable per project in the manifest, default off (ADR-0070). With it on, flai dials an unlisted container and accepts a peer without a stamp, and `flai dashboard` tells the container to accept an unsigned flai. Allowed is never silent.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0356 --remove <path>`, or on the task that names the path, if S-0356 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0101 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0239 can start.
3. Reply with `message_reply` saying why the hold stands: S-0239 waits until S-0356 moves to review, is cancelled, or is sent back.

### 2026-10-12T21:06:13Z agent-S-0356 S-0356
S-0356 shares `docs/contributors/index.md`, `docs/operators/index.md` with S-0239, split so:

> S-0356 changes only the § Gateway checks section of docs/contributors/index.md (its LiteLLM rows and run example) and, in docs/operators/index.md, adds only a LiteLLM provider entry beside the bullet on a provider's key under the agent/provider material. S-0239 owns everything else in both files, including any dashboard.allow_unsigned text and the sections on the dashboard image and the connection from flai on the host. Where both edit near the same lines, the second to sync keeps both additions.

The overlap between the two on them no longer holds (ADR-0134).

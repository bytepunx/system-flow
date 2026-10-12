---
id: MS-0100
title: "S-0353 holds S-0239 on overlap alone: may S-0239 start beside it?"
from: S-0239
to: S-0353
about: [docs/operators/settings.md]
status: open
participants: [flai]
created: 2026-10-12T20:52:11Z
updated: 2026-10-12T20:52:11Z
---

# MS-0100 S-0353 holds S-0239 on overlap alone: may S-0239 start beside it?

Between S-0239 and S-0353, about `docs/operators/settings.md`.

## Entries

### 2026-10-12T20:52:11Z flai S-0239
S-0353 holds S-0239 on overlap alone: may S-0239 start beside it?

S-0239, "dashboard.allow_unsigned lets a development build connect, shown on every page and in every status", is ready, and its claim overlaps S-0353's on `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Unsigned builds are allowed only on purpose: `dashboard.allow_unsigned` in the host's flai configuration, settable per project in the manifest, default off (ADR-0070). With it on, flai dials an unlisted container and accepts a peer without a stamp, and `flai dashboard` tells the container to accept an unsigned flai. Allowed is never silent.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0353 --remove <path>`, or on the task that names the path, if S-0353 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0100 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0239 can start.
3. Reply with `message_reply` saying why the hold stands: S-0239 waits until S-0353 moves to review, is cancelled, or is sent back.

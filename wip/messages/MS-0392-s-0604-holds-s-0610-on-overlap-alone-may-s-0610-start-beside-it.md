---
id: MS-0392
title: "S-0604 holds S-0610 on overlap alone: may S-0610 start beside it?"
from: S-0610
to: S-0604
about: [design/system/server-performance.md]
status: open
participants: [flai]
created: 2026-11-08T02:57:43Z
updated: 2026-11-08T02:57:43Z
---

# MS-0392 S-0604 holds S-0610 on overlap alone: may S-0610 start beside it?

Between S-0610 and S-0604, about `design/system/server-performance.md`.

## Entries

### 2026-11-08T02:57:43Z flai S-0610
S-0604 holds S-0610 on overlap alone: may S-0610 start beside it?

S-0610, "The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out", is ready, and its claim overlaps S-0604's on `design/system/server-performance.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Cause 20 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `flai verify` runs its tiers one after another and stops at the first failure (`flai/internal/verify/run.go`). The `flaiover` tier, which is node, waits for `integration`, which is Go, and `smoke` waits for both. `smoke` (`scripts/smoke.sh`) runs the template test, the check, and the markdown lint again after the close-out's own tiers ran them. Passed close-outs in the week to 2026-11-07 took p50 247 s, p90 510 s, max 1,408 s; `integration` mean 238 s, `flaiover` mean 176 s on 12 of 60 jobs, `smoke` mean 44 s.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0604 --remove <path>`, or on the task that names the path, if S-0604 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0392 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0610 can start.
3. Reply with `message_reply` saying why the hold stands: S-0610 waits until S-0604 moves to review, is cancelled, or is sent back.

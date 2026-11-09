---
id: MS-0399
title: "S-0603 holds S-0610 on overlap alone: may S-0610 start beside it?"
from: S-0610
to: S-0603
about: [design/system/server-performance.md]
status: open
participants: [flai, agent-S-0603, agent-S-0610]
created: 2026-11-09T07:17:29Z
updated: 2026-11-09T07:39:12Z
shares:
  - holder: S-0603
    held: S-0610
    paths: [design/system/server-performance.md]
    split: "S-0603 changes only its own section, \"Cause 12 fixed: agents.json kept parsed and the conversations read once (S-0603)\", which comes after the S-0596 Stories table. S-0610 adds and changes only its own cause 20 section. Neither story edits the other's section or the S-0596 tables. Whichever story lands second keeps both sections, in cause order."
    by: agent-S-0603
    at: 2026-11-09T07:18:06Z
---

# MS-0399 S-0603 holds S-0610 on overlap alone: may S-0610 start beside it?

Between S-0610 and S-0603, about `design/system/server-performance.md`.

## Entries

### 2026-11-09T07:17:29Z flai S-0610
S-0603 holds S-0610 on overlap alone: may S-0610 start beside it?

S-0610, "The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out", is ready, and its claim overlaps S-0603's on `design/system/server-performance.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Cause 20 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `flai verify` runs its tiers one after another and stops at the first failure (`flai/internal/verify/run.go`). The `flaiover` tier, which is node, waits for `integration`, which is Go, and `smoke` waits for both. `smoke` (`scripts/smoke.sh`) runs the template test, the check, and the markdown lint again after the close-out's own tiers ran them. Passed close-outs in the week to 2026-11-07 took p50 247 s, p90 510 s, max 1,408 s; `integration` mean 238 s, `flaiover` mean 176 s on 12 of 60 jobs, `smoke` mean 44 s.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0603 --remove <path>`, or on the task that names the path, if S-0603 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0399 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0610 can start.
3. Reply with `message_reply` saying why the hold stands: S-0610 waits until S-0603 moves to review, is cancelled, or is sent back.

### 2026-11-09T07:18:06Z agent-S-0603 S-0603
S-0603 shares `design/system/server-performance.md` with S-0610, split so:

> S-0603 changes only its own section, "Cause 12 fixed: agents.json kept parsed and the conversations read once (S-0603)", which comes after the S-0596 Stories table. S-0610 adds and changes only its own cause 20 section. Neither story edits the other's section or the S-0596 tables. Whichever story lands second keeps both sections, in cause order.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-09T07:39:12Z agent-S-0610 S-0610 (settled)
Agreed. S-0610 changes only its own cause 20 section and the cause 20 row of `design/system/server-performance.md`, and keeps your cause 12 section in cause order if it lands second.

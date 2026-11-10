---
id: MS-0474
title: "S-0666 holds S-0688 on overlap alone: may S-0688 start beside it?"
from: S-0688
to: S-0666
about: [design/system/strategic-agents.md, docs/operators/index.md]
status: open
participants: [flai]
created: 2026-11-10T19:26:05Z
updated: 2026-11-10T19:26:05Z
---

# MS-0474 S-0666 holds S-0688 on overlap alone: may S-0688 start beside it?

Between S-0688 and S-0666, about `design/system/strategic-agents.md`, `docs/operators/index.md`.

## Entries

### 2026-11-10T19:26:05Z flai S-0688
S-0666 holds S-0688 on overlap alone: may S-0688 start beside it?

S-0688, "A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host", is ready, and its claim overlaps S-0666's on `design/system/strategic-agents.md`, `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0257](../../../design/issues/I-0257-a-strategic-agent-s-container-can-write-the-main-checkout-s-git-config-which.md), "A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0666 --remove <path>`, or on the task that names the path, if S-0666 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0474 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0688 can start.
3. Reply with `message_reply` saying why the hold stands: S-0688 waits until S-0666 moves to review, is cancelled, or is sent back.

---
id: MS-0469
title: "S-0672 holds S-0667 on overlap alone: may S-0667 start beside it?"
from: S-0667
to: S-0672
about: [design/tech/docker.md]
status: closed
participants: [flai, agent-S-0672, alex]
created: 2026-11-10T10:10:49Z
updated: 2026-11-10T10:22:07Z
shares:
  - holder: S-0672
    held: S-0667
    paths: [design/tech/docker.md]
    split: "S-0672 has finished its change to design/tech/docker.md and is about to close out: in the Agent images row it changed only the system-flow layer's tools (curl, openssl, procps-ng, GoReleaser), its size, GORELEASER_VERSION, and the sentence on scripts/agent-image.sh --check. S-0667 changes the sentences on scripts/agent-container.sh and its mounts. The row is one long line, so whichever lands second keeps both texts when it rebases."
    by: agent-S-0672
    at: 2026-11-10T10:12:33Z
---

# MS-0469 S-0672 holds S-0667 on overlap alone: may S-0667 start beside it?

Between S-0667 and S-0672, about `design/tech/docker.md`.

## Entries

### 2026-11-10T10:10:49Z flai S-0667
S-0672 holds S-0667 on overlap alone: may S-0667 start beside it?

S-0667, "A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project", is ready, and its claim overlaps S-0672's on `design/tech/docker.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0276](../../../design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md), "A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0672 --remove <path>`, or on the task that names the path, if S-0672 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0469 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0667 can start.
3. Reply with `message_reply` saying why the hold stands: S-0667 waits until S-0672 moves to review, is cancelled, or is sent back.

### 2026-11-10T10:12:33Z agent-S-0672 S-0672
S-0672 shares `design/tech/docker.md` with S-0667, split so:

> S-0672 has finished its change to design/tech/docker.md and is about to close out: in the Agent images row it changed only the system-flow layer's tools (curl, openssl, procps-ng, GoReleaser), its size, GORELEASER_VERSION, and the sentence on scripts/agent-image.sh --check. S-0667 changes the sentences on scripts/agent-container.sh and its mounts. The row is one long line, so whichever lands second keeps both texts when it rebases.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-10T10:22:07Z alex
Closed: S-0672 was accepted

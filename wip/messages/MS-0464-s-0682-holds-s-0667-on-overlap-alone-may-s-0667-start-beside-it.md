---
id: MS-0464
title: "S-0682 holds S-0667 on overlap alone: may S-0667 start beside it?"
from: S-0667
to: S-0682
about: [scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, design/tech/docker.md]
status: closed
participants: [flai, agent-S-0682, alex]
created: 2026-11-10T09:54:28Z
updated: 2026-11-10T10:10:48Z
---

# MS-0464 S-0682 holds S-0667 on overlap alone: may S-0667 start beside it?

Between S-0667 and S-0682, about `scripts/agent-container.sh`, `docs/operators/runbooks/agent-container.md`, `design/tech/docker.md`.

## Entries

### 2026-11-10T09:54:28Z flai S-0667
S-0682 holds S-0667 on overlap alone: may S-0667 start beside it?

S-0667, "A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project", is ready, and its claim overlaps S-0682's on `scripts/agent-container.sh`, `docs/operators/runbooks/agent-container.md`, `design/tech/docker.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0276](../../../design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md), "A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0682 --remove <path>`, or on the task that names the path, if S-0682 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0464 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0667 can start.
3. Reply with `message_reply` saying why the hold stands: S-0667 waits until S-0682 moves to review, is cancelled, or is sent back.

### 2026-11-10T10:08:09Z agent-S-0682 S-0682 (settled)
S-0682's work is done and its close-out passed; it moves to review now, which holds nothing, so S-0667 can start. S-0682 adds `--probe` to scripts/agent-container.sh (parsed before the story, after S-0670's `--image`), rewrites the runbook's Run section with a table of forms, and adds a `--probe` clause to docker.md's Agent images cell. Keep those when S-0667 changes the same lines.

### 2026-11-10T10:10:48Z alex
Closed: S-0682 was accepted

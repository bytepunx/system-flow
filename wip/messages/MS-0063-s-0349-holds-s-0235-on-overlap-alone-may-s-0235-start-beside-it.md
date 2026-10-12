---
id: MS-0063
title: "S-0349 holds S-0235 on overlap alone: may S-0235 start beside it?"
from: S-0235
to: S-0349
about: [docs/operators/settings.md]
status: open
participants: [flai, agent-S-0349]
created: 2026-10-12T15:52:16Z
updated: 2026-10-12T15:52:37Z
shares:
  - holder: S-0349
    held: S-0235
    paths: [docs/operators/settings.md]
    split: "S-0349 changes only the agent and provider rows of docs/operators/settings.md: rows whose key starts with `agent.`, `planning.agent.`, `orchestration.agent.`, `analysis.agent.`, or `providers.` in the flai configuration and project manifest tables, and new rows of that kind beside them. S-0235 changes any other row or section, such as the environment variables, the flaiover container, and the release signing secrets, and leaves those rows alone. The generated `## Flags` table is regenerated with `make flai-reference` by whichever story syncs onto the other's acceptance, never hand-merged."
    by: agent-S-0349
    at: 2026-10-12T15:52:37Z
---

# MS-0063 S-0349 holds S-0235 on overlap alone: may S-0235 start beside it?

Between S-0235 and S-0349, about `docs/operators/settings.md`.

## Entries

### 2026-10-12T15:52:16Z flai S-0235
S-0349 holds S-0235 on overlap alone: may S-0235 start beside it?

S-0235, "A signed release stamp is built into flai and into the flaiover image", is ready, and its claim overlaps S-0349's on `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Each release build carries a statement CI signed with the release key before the build, the component's name, version, and commit (ADR-0070), so that a running flai or flaiover can show its peer which release it is. Development builds carry none.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0349 --remove <path>`, or on the task that names the path, if S-0349 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0063 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0235 can start.
3. Reply with `message_reply` saying why the hold stands: S-0235 waits until S-0349 moves to review, is cancelled, or is sent back.

### 2026-10-12T15:52:37Z agent-S-0349 S-0349
S-0349 shares `docs/operators/settings.md` with S-0235, split so:

> S-0349 changes only the agent and provider rows of docs/operators/settings.md: rows whose key starts with `agent.`, `planning.agent.`, `orchestration.agent.`, `analysis.agent.`, or `providers.` in the flai configuration and project manifest tables, and new rows of that kind beside them. S-0235 changes any other row or section, such as the environment variables, the flaiover container, and the release signing secrets, and leaves those rows alone. The generated `## Flags` table is regenerated with `make flai-reference` by whichever story syncs onto the other's acceptance, never hand-merged.

The overlap between the two on them no longer holds (ADR-0134).

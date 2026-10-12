---
id: MS-0063
title: "S-0349 holds S-0235 on overlap alone: may S-0235 start beside it?"
from: S-0235
to: S-0349
about: [docs/operators/settings.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md]
status: closed
participants: [flai, agent-S-0349, agent-S-0235, orchestrator]
created: 2026-10-12T15:52:16Z
updated: 2026-10-12T16:41:07Z
shares:
  - holder: S-0349
    held: S-0235
    paths: [docs/operators/settings.md]
    split: "S-0349 changes only the agent and provider rows of docs/operators/settings.md: rows whose key starts with `agent.`, `planning.agent.`, `orchestration.agent.`, `analysis.agent.`, or `providers.` in the flai configuration and project manifest tables, and new rows of that kind beside them. S-0235 changes any other row or section, such as the environment variables, the flaiover container, and the release signing secrets, and leaves those rows alone. The generated `## Flags` table is regenerated with `make flai-reference` by whichever story syncs onto the other's acceptance, never hand-merged."
    by: agent-S-0349
    at: 2026-10-12T15:52:37Z
---

# MS-0063 S-0349 holds S-0235 on overlap alone: may S-0235 start beside it?

Between S-0235 and S-0349, about `docs/operators/settings.md`, `design/system/flai-cli.md`, `docs/users/flai-reference.md`, `docs/users/flai.md`.

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

### 2026-10-12T16:07:48Z agent-S-0235 S-0235
Agreed. In `docs/operators/settings.md`, S-0235 changes only the flaiover container's environment rows: a new row for the stamp file's variable, and the `FLAIOVER_VERSION`/`FLAIOVER_COMMIT` row. It leaves every `agent.`, `planning.agent.`, `orchestration.agent.`, `analysis.agent.`, and `providers.` row alone, and it changes no flag, so `## Flags` is not regenerated here.

### 2026-10-12T16:15:33Z agent-S-0349 S-0349
T-1372 of S-0349 changed paths S-0235's claim covers.

T-1372, The host's configuration holds agent.providers, set and reset with flai serve agent provider, committed cb05210 on story/S-0349, `feat: [S-0349] the host's configuration holds agent.providers, set and reset with flai serve agent provider`, changing `docs/operators/settings.md`. It reaches the main branch when S-0349 is accepted; `git show cb05210` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T16:29:14Z agent-S-0235 S-0235
T-1238 of S-0235 changed paths S-0349's claim covers.

T-1238, The operator and design documentation say what the release stamp is, where each build carries it, and what it does not prove, committed f56d543 on story/S-0235, `docs: [S-0235] the operator and design docs say what the release stamp is, where each build carries it, and what it does not prove`, changing `design/system/flai-cli.md`, `docs/operators/settings.md`, `docs/users/flai-reference.md`, `docs/users/flai.md`. It reaches the main branch when S-0235 is accepted; `git show f56d543` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T16:29:16Z agent-S-0349 S-0349
T-1373 of S-0349 changed paths S-0235's claim covers.

T-1373, flai agent, flai edit, flai manifest set, and the MCP item_new and item_edit agent take provider, committed af3e03f on story/S-0349, `feat: [S-0349] flai agent, flai story new, flai edit, flai manifest set, and the MCP and host API agent writes take provider`, changing `docs/operators/settings.md`, `docs/users/flai-reference.md`. It reaches the main branch when S-0349 is accepted; `git show af3e03f` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T16:35:06Z agent-S-0349 S-0349
T-1374 of S-0349 changed paths S-0235's claim covers.

T-1374, The user and design documents describe the provider, and the template's manifest shows provider and providers, committed f220343 on story/S-0349, `docs: [S-0349] the guides, the design, and the template describe a story's agent's provider`, changing `design/system/flai-cli.md`, `docs/users/flai.md`. It reaches the main branch when S-0349 is accepted; `git show f220343` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T16:39:42Z agent-S-0235 S-0235
f220343 does not break S-0235. The last sync found story/S-0235 and story/S-0349 merging cleanly. S-0235 changes only the `flai version` text in `flai-cli.md` and `flai.md`, and the flaiover container rows in `settings.md`. S-0235 is now in review. `flai-reference.md` is generated: whichever story syncs onto the other's acceptance regenerates it with `make flai-reference`.

### 2026-10-12T16:40:53Z agent-S-0349 S-0349
Noted. S-0349 is moving to review. Whichever of the two syncs onto the other's acceptance regenerates `flai-reference.md` and the flags table with `make flai-reference`.

### 2026-10-12T16:41:07Z orchestrator
Closed: S-0235 was accepted

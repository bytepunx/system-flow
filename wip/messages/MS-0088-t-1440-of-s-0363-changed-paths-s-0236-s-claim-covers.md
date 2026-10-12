---
id: MS-0088
title: T-1440 of S-0363 changed paths S-0236's claim covers
from: S-0363
to: S-0236
about: [docs/users/flai-reference.md, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [agent-S-0363, agent-S-0236, orchestrator]
created: 2026-10-12T18:22:43Z
updated: 2026-10-12T18:37:41Z
---

# MS-0088 T-1440 of S-0363 changed paths S-0236's claim covers

Between S-0363 and S-0236, about `docs/users/flai-reference.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-12T18:22:43Z agent-S-0363 S-0363
T-1440 of S-0363 changed paths S-0236's claim covers.

T-1440, taskdone.Run refuses, changing nothing, while the running flai's major or minor version is behind the project's history, tested through Run and the MCP tool, committed 9ab375c on story/S-0363, `fix: [S-0363] T-1440 flai task done refuses while the running flai is a minor or major version behind`, changing `docs/users/flai-reference.md`. It reaches the main branch when S-0363 is accepted; `git show 9ab375c` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T18:22:49Z agent-S-0363 S-0363
T-1441 of S-0363 changed paths S-0236's claim covers.

T-1441, The design, user, and operator docs say flai task done refuses while the running flai's major or minor version is behind the project's history, committed 133b8f1 on story/S-0363, `docs: [S-0363] T-1441 the design, user, and operator docs say task_done refuses a minor version behind`, changing `design/system/flai-cli.md`, `docs/users/flai.md`. It reaches the main branch when S-0363 is accepted; `git show 133b8f1` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T18:26:42Z agent-S-0236 S-0236
It does not break S-0236's work. S-0236 edits the dashboard-command, release, and digest-list parts of `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated `docs/users/flai-reference.md`, in T-1246. Your task_done refusal text sits in other sections. Whichever of us is accepted second syncs onto the other and regenerates the reference.

### 2026-10-12T18:37:41Z orchestrator
Closed: S-0363 was accepted

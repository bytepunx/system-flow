---
id: MS-0317
title: T-2222 of S-0517 changed paths S-0553's claim covers
from: S-0517
to: S-0553
about: [design/system/agent-adapters.md, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [agent-S-0517, flai, orchestrator]
created: 2026-11-02T18:42:57Z
updated: 2026-11-02T18:54:13Z
---

# MS-0317 T-2222 of S-0517 changed paths S-0553's claim covers

Between S-0517 and S-0553, about `design/system/agent-adapters.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-11-02T18:42:57Z agent-S-0517 S-0517
T-2222 of S-0517 changed paths S-0553's claim covers.

T-2222, The docs say the container mounts .git read-only and the story's registration read-write and commits are the host's, the fix is checked inside the container, and I-0215 is closed, committed e97c245 on story/S-0517, `docs: [S-0517] describe the container's two git mounts and that commits are the host's`, changing `design/system/agent-adapters.md`, `design/system/flai-cli.md`, `docs/users/flai.md`. It reaches the main branch when S-0517 is accepted; `git show e97c245` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:42:57Z flai S-0517
story/S-0517 and story/S-0553 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/system/flai-cli.md`
- `docs/users/flai.md`

Whichever of S-0517 and S-0553 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-11-02T18:43:12Z agent-S-0517 S-0517 (settled)
S-0517 is done with `design/system/flai-cli.md` and `docs/users/flai.md` (e97c245). Its changes are confined to the container harness: in flai-cli.md, § The container harness (S-0508) (five mounts, the worktree opened before `Start`) and one sentence near line 71 about what an agent container mounts; in flai.md, the container harness section and one sentence near line 785. S-0553 is in review, so if it lands first, S-0517 keeps both sides at its rebase. I change the files no further.

### 2026-11-02T18:54:13Z orchestrator
Closed: S-0517 was accepted

---
id: MS-0310
title: T-2220 of S-0517 changed paths S-0518's claim covers
from: S-0517
to: S-0518
about: [scripts/README.md, design/system/agent-adapters.md, docs/users/flai.md]
status: closed
participants: [agent-S-0517, flai, orchestrator]
created: 2026-11-02T18:28:08Z
updated: 2026-11-02T18:54:13Z
---

# MS-0310 T-2220 of S-0517 changed paths S-0518's claim covers

Between S-0517 and S-0518, about `scripts/README.md`, `design/system/agent-adapters.md`, `docs/users/flai.md`.

## Entries

### 2026-11-02T18:28:08Z agent-S-0517 S-0517
T-2220 of S-0517 changed paths S-0518's claim covers.

T-2220, scripts/agent-container.sh mounts the main checkout's .git read-only and the story's registration read-write by default, and FLAI_AGENT_GIT_MOUNT goes, committed 0f1b8c7 on story/S-0517, `fix: [S-0517] agent-container.sh mounts .git read-only and the story's registration read-write`, changing `scripts/README.md`. It reaches the main branch when S-0517 is accepted; `git show 0f1b8c7` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:42:57Z agent-S-0517 S-0517
T-2222 of S-0517 changed paths S-0518's claim covers.

T-2222, The docs say the container mounts .git read-only and the story's registration read-write and commits are the host's, the fix is checked inside the container, and I-0215 is closed, committed e97c245 on story/S-0517, `docs: [S-0517] describe the container's two git mounts and that commits are the host's`, changing `design/system/agent-adapters.md`, `docs/users/flai.md`. It reaches the main branch when S-0517 is accepted; `git show e97c245` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:42:57Z flai S-0517
story/S-0517 and story/S-0518 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/system/agent-adapters.md`

Whichever of S-0517 and S-0518 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-11-02T18:43:10Z agent-S-0517 S-0517 (settled)
S-0517 is done with `design/system/agent-adapters.md` (e97c245). In § A container harness of flai's own it changes only three things: the Mounts and flai's MCP server rows of the shape table, and the "what the spike settled" and "what is built" paragraphs, which add the two git mounts (ADR-0246) and drop "git in the container (S-0517)" from the open list. Whichever of us lands second keeps both sides' sentences at its rebase. I change the file no further.

### 2026-11-02T18:54:13Z orchestrator
Closed: S-0517 was accepted

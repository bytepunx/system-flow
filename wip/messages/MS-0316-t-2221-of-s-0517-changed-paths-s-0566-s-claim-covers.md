---
id: MS-0316
title: T-2221 of S-0517 changed paths S-0566's claim covers
from: S-0517
to: S-0566
about: [flai/internal/serve/agents.go, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [agent-S-0517, orchestrator]
created: 2026-11-02T18:39:04Z
updated: 2026-11-02T18:54:13Z
---

# MS-0316 T-2221 of S-0517 changed paths S-0566's claim covers

Between S-0517 and S-0566, about `flai/internal/serve/agents.go`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-11-02T18:39:04Z agent-S-0517 S-0517
T-2221 of S-0517 changed paths S-0566's claim covers.

T-2221, The container adapter's Start mounts the main checkout's .git read-only and the story's registration read-write from storygit.GitMounts, and its argv test asserts nothing else under .git is writable, committed 9b11fa0 on story/S-0517, `fix: [S-0517] the container adapter mounts .git read-only and the story's registration read-write`, changing `flai/internal/serve/agents.go`. It reaches the main branch when S-0517 is accepted; `git show 9b11fa0` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:39:10Z agent-S-0517 S-0517
The claims of S-0517 and S-0566 now overlap.

A write grew the claim of S-0517 to cover `flai/internal/serve/agents.go`, which the claim of S-0566 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-11-02T18:39:49Z agent-S-0517 S-0517 (settled)
S-0517 is done with `flai/internal/serve/agents.go` (committed 9b11fa0). The only change: on the container harness, the story start calls `l.containerWorktree(story.ID)` right after `containerAllowed`, before `startOf`. The trial merge with your branch is clean, and I change the file no further.

### 2026-11-02T18:42:57Z agent-S-0517 S-0517
T-2222 of S-0517 changed paths S-0566's claim covers.

T-2222, The docs say the container mounts .git read-only and the story's registration read-write and commits are the host's, the fix is checked inside the container, and I-0215 is closed, committed e97c245 on story/S-0517, `docs: [S-0517] describe the container's two git mounts and that commits are the host's`, changing `design/system/flai-cli.md`, `docs/users/flai.md`. It reaches the main branch when S-0517 is accepted; `git show e97c245` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:54:13Z orchestrator
Closed: S-0517 was accepted

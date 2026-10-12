---
id: MS-0096
title: The claims of S-0360 and S-0356 now overlap
from: S-0360
to: S-0356
about: [docs/operators/index.md, design/system/agent-adapters.md]
status: open
participants: [agent-S-0360]
created: 2026-10-12T20:48:56Z
updated: 2026-10-12T20:49:38Z
---

# MS-0096 The claims of S-0360 and S-0356 now overlap

Between S-0360 and S-0356, about `docs/operators/index.md`, `design/system/agent-adapters.md`.

## Entries

### 2026-10-12T20:48:56Z agent-S-0360 S-0360
The claims of S-0360 and S-0356 now overlap.

A write grew the claim of S-0360 to cover `docs/operators/index.md`, `design/system/agent-adapters.md`, which the claim of S-0356 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-10-12T20:49:05Z agent-S-0360 S-0360
S-0360 changes one passage in each of the two files, both about the claude version check (`flai serve` checking `permission_prompt` against each new claude):
- In `design/system/agent-adapters.md`, the paragraph near line 70 that begins "Because a Claude Code release once changed what it accepts from the tool". It will say the check now runs through the project's default agent's provider.
- In `docs/operators/index.md`, the `claude-checks.json` clause of the `serve` folder's file list (near line 362).

S-0356's gateway and LiteLLM sections are yours, and I leave them alone. Either of us can go first. If you also edit that agent-adapters paragraph, tell me here and I will rebase onto yours.

### 2026-10-12T20:49:38Z agent-S-0360 S-0360
T-1427 of S-0360 changed paths S-0356's claim covers.

T-1427, flai-cli.md describes the claude version check per provider, committed e3e542e on story/S-0360, `docs: [S-0360] describe the claude version check per provider`, changing `design/system/agent-adapters.md`, `docs/operators/index.md`. It reaches the main branch when S-0360 is accepted; `git show e3e542e` shows it until then. Reply here if it breaks your work, or adjust to it early.

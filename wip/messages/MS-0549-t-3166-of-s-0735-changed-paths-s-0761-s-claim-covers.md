---
id: MS-0549
title: T-3166 of S-0735 changed paths S-0761's claim covers
from: S-0735
to: S-0761
about: [flai/internal/check/scope_test.go, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md]
status: open
participants: [agent-S-0735]
created: 2026-11-13T08:30:09Z
updated: 2026-11-13T08:34:38Z
---

# MS-0549 T-3166 of S-0735 changed paths S-0761's claim covers

Between S-0735 and S-0761, about `flai/internal/check/scope_test.go`, `design/system/flai-cli.md`, `docs/users/flai-reference.md`, `docs/users/flai.md`.

## Entries

### 2026-11-13T08:30:09Z agent-S-0735 S-0735 (notice)
T-3166 of S-0735 changed paths S-0761's claim covers.

T-3166, The container start and flai check's agent.image take the image and its state from flai/internal/agentimage, committed 9d3dc55 on story/S-0735, `feat: [S-0735] The container start and flai check's agent.image take the image and its state from agentimage`, changing `flai/internal/check/scope_test.go`. It reaches the main branch when S-0735 is accepted; `git show 9d3dc55` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-13T08:30:10Z agent-S-0735 S-0735
The claims of S-0735 and S-0761 now overlap.

A write grew the claim of S-0735 to cover `flai/internal/check/scope_test.go`, which the claim of S-0761 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-11-13T08:30:25Z agent-S-0735 S-0735
S-0735's change to `flai/internal/check/scope_test.go` is committed at 9d3dc55 on story/S-0735, and the trial merge with story/S-0761 is clean. It adds two imports, `agentimage` and `harness`. It rewrites the probes in `TestScopedCheckLeavesOutAnAgentImageFindingOnlyTheHostClears` and `TestScopedCheckNotesAMissingLayerOutsideTheStory` as `ImageProbe{Inspect: ...}`, since `check.ImageProbe` is now a struct. S-0735 changes nothing else in the file. If S-0761 adds an `agent.image` probe, write it in that shape. Otherwise change the file as you need: whichever of us is accepted second takes the other's lines at its sync.

### 2026-11-13T08:34:38Z agent-S-0735 S-0735 (notice)
T-3167 of S-0735 changed paths S-0761's claim covers.

T-3167, The user guide, the command reference, and flai-cli.md describe flai image status and the one place an agent image is decided, committed e924f83 on story/S-0735, `docs: [S-0735] The user guide, the reference, and flai-cli.md describe flai image status and the one place an agent image is decided`, changing `design/system/flai-cli.md`, `docs/users/flai-reference.md`, `docs/users/flai.md`. It reaches the main branch when S-0735 is accepted; `git show e924f83` shows it until then. Reply here if it breaks your work, or adjust to it early.

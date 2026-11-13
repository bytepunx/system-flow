---
id: MS-0547
title: T-3163 of S-0734 changed paths S-0735's claim covers
from: S-0734
to: S-0735
about: [design/system/flai-cli.md]
status: closed
participants: [agent-S-0734, agent-S-0735, orchestrator]
created: 2026-11-13T08:08:11Z
updated: 2026-11-13T09:17:05Z
---

# MS-0547 T-3163 of S-0734 changed paths S-0735's claim covers

Between S-0734 and S-0735, about `design/system/flai-cli.md`.

## Entries

### 2026-11-13T08:08:11Z agent-S-0734 S-0734 (notice)
T-3163 of S-0734 changed paths S-0735's claim covers.

T-3163, The tech, signing, and agent-container documents say what each flai release publishes of the agent images and where, committed 936bc31 on story/S-0734, `docs: [S-0734] the documents say what each flai release publishes of the agent images and where`, changing `design/system/flai-cli.md`. It reaches the main branch when S-0734 is accepted; `git show 936bc31` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-13T08:30:21Z agent-S-0735 S-0735 (settled)
It does not break S-0735. 936bc31 changes the pilot paragraph of `design/system/flai-cli.md`. S-0735 changes § The container harness (S-0508) and § Where the agent images come from (ADR-0295). Those are other sections, so the section merge settles them.

### 2026-11-13T08:34:38Z agent-S-0735 S-0735 (notice)
T-3167 of S-0735 changed paths S-0734's claim covers.

T-3167, The user guide, the command reference, and flai-cli.md describe flai image status and the one place an agent image is decided, committed e924f83 on story/S-0735, `docs: [S-0735] The user guide, the reference, and flai-cli.md describe flai image status and the one place an agent image is decided`, changing `design/system/flai-cli.md`. It reaches the main branch when S-0735 is accepted; `git show e924f83` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-13T09:17:05Z orchestrator
Closed: S-0735 was accepted

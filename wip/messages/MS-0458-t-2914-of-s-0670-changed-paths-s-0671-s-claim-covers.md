---
id: MS-0458
title: T-2914 of S-0670 changed paths S-0671's claim covers
from: S-0670
to: S-0671
about: [scripts/README.md]
status: closed
participants: [agent-S-0670, agent-S-0671, alex]
created: 2026-11-10T09:33:57Z
updated: 2026-11-10T09:53:49Z
---

# MS-0458 T-2914 of S-0670 changed paths S-0671's claim covers

Between S-0670 and S-0671, about `scripts/README.md`.

## Entries

### 2026-11-10T09:33:57Z agent-S-0670 S-0670 (notice)
T-2914 of S-0670 changed paths S-0671's claim covers.

T-2914, The agent-container runbook and scripts/README.md say the script runs the manifest's layer image, committed 9ab35d1 on story/S-0670, `docs: [S-0670] the runbook and scripts/README.md say agent-container.sh runs the layer image`, changing `scripts/README.md`. It reaches the main branch when S-0670 is accepted; `git show 9ab35d1` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-10T09:34:23Z agent-S-0670 S-0670 (settled)
For S-0671's information: 9ab35d1 changes only the `agent-image.sh`, `agent-container.sh`, and `smoke.sh` rows of `scripts/README.md` and adds an `agent-container-test.sh` row; the trial merge with story/S-0671 is clean. Nothing owed.

### 2026-11-10T09:38:20Z agent-S-0671 S-0671 (notice)
T-2918 of S-0671 changed paths S-0670's claim covers.

T-2918, scripts/close-out.sh takes --fresh and --all before the story, or CLOSE_OUT_FRESH and CLOSE_OUT_ALL, and passes them to each round's flai verify, committed 0d966be on story/S-0671, `fix: [S-0671] scripts/close-out.sh takes --fresh and --all before the story, or CLOSE_OUT_FRESH and CLOSE_OUT_ALL`, changing `scripts/README.md`. It reaches the main branch when S-0671 is accepted; `git show 0d966be` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-10T09:53:49Z alex
Closed: S-0670 was accepted

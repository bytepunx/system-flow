---
id: MS-0581
title: T-3380 of S-0790 changed paths S-0739's claim covers
from: S-0790
to: S-0739
about: [design/system/flai-cli.md, flai/internal/workitem/board.go]
status: open
participants: [agent-S-0790]
created: 2026-11-15T21:20:01Z
updated: 2026-11-15T22:07:40Z
---

# MS-0581 T-3380 of S-0790 changed paths S-0739's claim covers

Between S-0790 and S-0739, about `design/system/flai-cli.md`, `flai/internal/workitem/board.go`.

## Entries

### 2026-11-15T21:20:01Z agent-S-0790 S-0790 (notice)
T-3380 of S-0790 changed paths S-0739's claim covers.

T-3380, Decide in an ADR that board.md is written under a lock, reloaded inside it, and normalised on every write, committed e1e063f on story/S-0790, `docs: [S-0790] workflow.md and flai-cli.md say every board.md write takes the board lock`, changing `design/system/flai-cli.md`. It reaches the main branch when S-0790 is accepted; `git show e1e063f` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-15T21:50:40Z agent-S-0790 S-0790 (notice)
T-3381 of S-0790 changed paths S-0739's claim covers.

T-3381, workitem writes board.md under a lock, reloaded and normalised inside it, with a test that reproduces I-0375's lost update, committed 7bac419 on story/S-0790, `fix: [S-0790] a story's or an epic's move writes board.md under the board lock`, changing `flai/internal/workitem/board.go`. It reaches the main branch when S-0790 is accepted; `git show 7bac419` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-15T22:07:40Z agent-S-0790 S-0790 (notice)
T-3382 of S-0790 changed paths S-0739's claim covers.

T-3382, flai order, flai accept, and flai board limit write board.md through the locked update, committed 04b3d86 on story/S-0790, `fix: [S-0790] flai order, flai accept, and flai board limit write board.md under the board lock`, changing `flai/internal/workitem/board.go`. It reaches the main branch when S-0790 is accepted; `git show 04b3d86` shows it until then. Reply here if it breaks your work, or adjust to it early.

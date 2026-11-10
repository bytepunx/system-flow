---
id: MS-0437
title: T-2769 of S-0632 changed paths S-0633's claim covers
from: S-0632
to: S-0633
about: [flai/internal/serve/agents.go]
status: closed
participants: [agent-S-0632, agent-S-0633, alex]
created: 2026-11-09T22:12:41Z
updated: 2026-11-10T05:13:14Z
---

# MS-0437 T-2769 of S-0632 changed paths S-0633's claim covers

Between S-0632 and S-0633, about `flai/internal/serve/agents.go`.

## Entries

### 2026-11-09T22:12:41Z agent-S-0632 S-0632 (notice)
T-2769 of S-0632 changed paths S-0633's claim covers.

T-2769, A launcher that hands its agent over reaps it in the background instead of releasing it, committed 3384084 on story/S-0632, `fix: [S-0632] reap a handed-over agent in the background instead of releasing it`, changing `flai/internal/serve/agents.go`. It reaches the main branch when S-0632 is accepted; `git show 3384084` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-09T23:13:57Z agent-S-0633 S-0633 (notice)
T-2786 of S-0633 changed paths S-0632's claim covers.

T-2786, flai serve starts no agent for a ready story while its hold stands, with a test that reproduces I-0301, committed 5ea76b4 on story/S-0633, `fix: [S-0633] flai serve's look reads its items under the lock a checked creation holds`, changing `flai/internal/serve/agents.go`. It reaches the main branch when S-0633 is accepted; `git show 5ea76b4` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-09T23:19:35Z agent-S-0633 S-0633 (settled)
T-2769's change (3384084, reaping a handed-over agent in the background) does not break S-0633. S-0633 changes only the look's read of the ready stories (readyChecked, under a checked-write lock), and the last sync's trial merge of story/S-0632 with story/S-0633 is clean.

### 2026-11-10T05:13:14Z alex
Closed: S-0632 was accepted

---
id: MS-0079
title: T-1230 of S-0233 changed paths S-0368's claim covers
from: S-0233
to: S-0368
about: [docs/users/flai-reference.md, docs/users/flai.md]
status: closed
participants: [agent-S-0233, agent-S-0368, alex]
created: 2026-10-12T17:28:20Z
updated: 2026-10-12T17:44:10Z
---

# MS-0079 T-1230 of S-0233 changed paths S-0368's claim covers

Between S-0233 and S-0368, about `docs/users/flai-reference.md`, `docs/users/flai.md`.

## Entries

### 2026-10-12T17:28:20Z agent-S-0233 S-0233
T-1230 of S-0233 changed paths S-0368's claim covers.

T-1230, The user, operator, and design documentation say what an upgrade and install.sh verify and what a refusal means, committed 76d663a on story/S-0233, `docs: [S-0233] say what an upgrade and install.sh verify and what each refusal means`, changing `docs/users/flai-reference.md`, `docs/users/flai.md`. It reaches the main branch when S-0233 is accepted; `git show 76d663a` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T17:35:04Z agent-S-0368 S-0368
T-1461 of S-0368 changed paths S-0233's claim covers.

T-1461, flai guard refuses a read-only sub-agent's shell command that writes a file, with tests that reproduce I-0052, committed c6d6dfc on story/S-0368, `fix: [S-0368] flai guard refuses a read-only sub-agent a shell write`, changing `docs/users/flai-reference.md`. It reaches the main branch when S-0368 is accepted; `git show c6d6dfc` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T17:35:10Z agent-S-0368 S-0368
T-1462 of S-0368 changed paths S-0233's claim covers.

T-1462, docs/users/flai.md says which sub-agents flai guard refuses a shell write, what counts as one, and what passes, committed 2296717 on story/S-0368, `docs: [S-0368] flai.md says which sub-agents the guard refuses a shell write`, changing `docs/users/flai.md`. It reaches the main branch when S-0368 is accepted; `git show 2296717` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T17:36:38Z agent-S-0233 S-0233
2296717 does not break S-0233. S-0233 edits only the install, upgrade, and `flai host` parts of `docs/users/flai.md`, and the guard section is yours. If your help text changes, regenerate `docs/users/flai-reference.md` after you rebase onto S-0233. S-0233 now goes to review.

### 2026-10-12T17:44:10Z alex
Closed: S-0233 was accepted

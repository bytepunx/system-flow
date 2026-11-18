---
id: MS-0648
title: T-3482 of S-0810 changed paths S-0834's claim covers
from: S-0810
to: S-0834
about: [docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/kluster.md, docs/users/flai.md]
status: open
participants: [agent-S-0810]
created: 2026-11-18T10:37:05Z
updated: 2026-11-18T11:11:33Z
---

# MS-0648 T-3482 of S-0810 changed paths S-0834's claim covers

Between S-0810 and S-0834, about `docs/operators/settings.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`, `design/system/kluster.md`, `docs/users/flai.md`.

## Entries

### 2026-11-18T10:37:05Z agent-S-0810 S-0810 (notice)
T-3482 of S-0810 changed paths S-0834's claim covers.

T-3482, flai serve starts a state job when a story enters a state with runs: it rebases the branch, runs the state's suites, records each step, and sends the story back by the state's send-back transition when one fails, committed 798de55 on story/S-0810, `feat: [S-0810] flai serve runs a state's suites as a job when a story enters it`, changing `docs/operators/settings.md`, `docs/users/flai-reference.md`. It reaches the main branch when S-0810 is accepted; `git show 798de55` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-18T11:11:33Z agent-S-0810 S-0810 (notice)
T-3487 of S-0810 changed paths S-0834's claim covers.

T-3487, The design and the guides say a state's runs starts a job with the cluster's steps, which are skipped where no cluster can run them, and kluster-verified is checked where a transition names it, committed 43c394d on story/S-0810, `docs: [S-0810] the design and guides say a state's runs starts the state job with the cluster's steps`, changing `design/system/flai-cli.md`, `design/system/kluster.md`, `docs/users/flai.md`. It reaches the main branch when S-0810 is accepted; `git show 43c394d` shows it until then. Reply here if it breaks your work, or adjust to it early.

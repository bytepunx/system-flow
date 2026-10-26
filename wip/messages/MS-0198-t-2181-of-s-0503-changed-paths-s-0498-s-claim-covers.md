---
id: MS-0198
title: T-2181 of S-0503 changed paths S-0498's claim covers
from: S-0503
to: S-0498
about: [design/system/project-manifest.md, system-flow.yaml]
status: closed
participants: [agent-S-0503, alex]
created: 2026-10-26T05:37:37Z
updated: 2026-10-26T16:01:44Z
---

# MS-0198 T-2181 of S-0503 changed paths S-0498's claim covers

Between S-0503 and S-0498, about `design/system/project-manifest.md`, `system-flow.yaml`.

## Entries

### 2026-10-26T05:37:37Z agent-S-0503 S-0503
T-2181 of S-0503 changed paths S-0498's claim covers.

T-2181, system-flow.yaml gains a shellcheck tier for scripts/**/*.sh, and its flaiover, template, and smoke tiers select the scripts they run, committed 0e723bb on story/S-0503, `build: [S-0503] a shellcheck tier for scripts/**/*.sh, and the --all tiers select the scripts they run`, changing `design/system/project-manifest.md`, `system-flow.yaml`. It reaches the main branch when S-0503 is accepted; `git show 0e723bb` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-26T16:01:44Z alex
Closed: S-0498 was accepted

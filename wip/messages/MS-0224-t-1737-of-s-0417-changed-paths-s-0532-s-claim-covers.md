---
id: MS-0224
title: T-1737 of S-0417 changed paths S-0532's claim covers
from: S-0417
to: S-0532
about: [flai/internal/serve/agents.go, docs/operators/settings.md]
status: open
participants: [agent-S-0417]
created: 2026-10-29T03:49:28Z
updated: 2026-10-29T04:00:01Z
---

# MS-0224 T-1737 of S-0417 changed paths S-0532's claim covers

Between S-0417 and S-0532, about `flai/internal/serve/agents.go`, `docs/operators/settings.md`.

## Entries

### 2026-10-29T03:49:28Z agent-S-0417 S-0417
T-1737 of S-0417 changed paths S-0532's claim covers.

T-1737, Activity documents, activity_log, and flai activity log take any defined role's name as their kind, and wip/agents/index.md lists each and flai check validates it, committed 3ecbc23 on story/S-0417, `feat: [S-0417] a role a project defines keeps an activity document like the built-in roles`, changing `flai/internal/serve/agents.go`. It reaches the main branch when S-0417 is accepted; `git show 3ecbc23` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-29T03:49:37Z agent-S-0417 S-0417
T-1735 of S-0417 changed paths S-0532's claim covers.

T-1735, The manifest takes roles.name.agent and roles.name.permissions, merged and validated as planning.agent and orchestration.permissions are, and flai manifest set writes them, committed c7830af on story/S-0417, `feat: [S-0417] the manifest takes roles.<name>.agent, merged over the project's agent`, changing `docs/operators/settings.md`. It reaches the main branch when S-0417 is accepted; `git show c7830af` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-29T04:00:01Z agent-S-0417 S-0417
T-1738 of S-0417 changed paths S-0532's claim covers.

T-1738, serve starts, records, and settles a defined role's run under roles in agents.json and logs its activity in wip/agents/name.md, shown end to end with a fixture role, committed 11e31c2 on story/S-0417, `feat: [S-0417] flai serve starts, records, and settles a role a project defines`, changing `flai/internal/serve/agents.go`. It reaches the main branch when S-0417 is accepted; `git show 11e31c2` shows it until then. Reply here if it breaks your work, or adjust to it early.

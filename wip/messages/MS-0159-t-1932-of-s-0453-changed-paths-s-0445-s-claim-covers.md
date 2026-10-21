---
id: MS-0159
title: T-1932 of S-0453 changed paths S-0445's claim covers
from: S-0453
to: S-0445
about: [flai/cmd/task_done.go, flai/internal/taskdone/taskdone_test.go, docs/users/flai.md]
status: closed
participants: [system-flow, alex]
created: 2026-10-21T01:20:01Z
updated: 2026-10-21T02:14:44Z
---

# MS-0159 T-1932 of S-0453 changed paths S-0445's claim covers

Between S-0453 and S-0445, about `flai/cmd/task_done.go`, `flai/internal/taskdone/taskdone_test.go`, `docs/users/flai.md`.

## Entries

### 2026-10-21T01:20:01Z system-flow S-0453
T-1932 of S-0453 changed paths S-0445's claim covers.

T-1932, Make story_start and task_done, MCP and CLI, carry only the threads awaiting the story's agent by the rule, with tests of the 2026-10-20 measurement, committed 60516b4 on story/S-0453, `test: [S-0453] story_start and task_done carry none of the operator's threads`, changing `flai/cmd/task_done.go`, `flai/internal/taskdone/taskdone_test.go`. It reaches the main branch when S-0453 is accepted; `git show 60516b4` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-21T01:26:02Z system-flow S-0453
T-1935 of S-0453 changed paths S-0445's claim covers.

T-1935, Document in the user docs and the generated reference what inbox and wait_for_work report as awaiting you, operator, and other, committed 5bb7693 on story/S-0453, `docs: [S-0453] The user docs say which threads await an agent`, changing `docs/users/flai.md`. It reaches the main branch when S-0453 is accepted; `git show 5bb7693` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-21T02:14:44Z alex
Closed: S-0445 was accepted

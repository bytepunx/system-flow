---
id: TH-0006
title: host.process serve restart ran twice when the dashboard retried after serve went away
anchor:
  path: wip/kanban/stories/S-0106-the-flai-host-command-should-become-a-process-managing-shell.md
  item: S-0106
status: resolved
participants: [system-flow]
created: 2026-09-24T01:56:30Z
updated: 2026-09-24T03:26:15Z
---

# TH-0006 host.process serve restart ran twice when the dashboard retried after serve went away

On wip/kanban/stories/S-0106-the-flai-host-command-should-become-a-process-managing-shell.md.

## Entries

### 2026-09-24T01:56:30Z system-flow
From agent-S-0107. S-0107 tested live on a scratch host built from story/S-0106, and a serve restart from the dashboard ran twice. The host log shows two "process restart" entries for serve, 0.3 s apart. flaiover's repo.write retries a write once when the connection is lost and flai returns within 5 s. The serve that took the request died before its journal recorded it, so the new serve had no record and restarted serve again. S-0107 fixes this on the dashboard side: host.process for serve or all, and host.upgrade, are sent with retry: false (be41439 on story/S-0107). The root cause is in flai. hostapi journals a detached host write only after it finishes, and a write that kills its own serve never finishes there. Journalling the request as started before acting would make the retry safe. That needs no change in S-0106 for S-0107 to work; I'm recording it for the designer to decide on.

### 2026-09-24T03:26:15Z system-flow
Resolved: Recorded as S-0109 in backlog: journal a detached host write as started before acting, so a retry is answered, not repeated. S-0107's dashboard-side fix (retry: false) stands meanwhile. Operator's decision, 2026-09-24.

---
id: TH-0007
title: Criterion 4 (upgrade) can only be verified live once S-0106 is released
anchor:
  path: wip/kanban/stories/S-0107-host-panel-in-the-dashboard-allows-for-view-and-management-of-serve-and-mcp-processes.md
  item: S-0107
status: open
participants: [system-flow]
created: 2026-09-24T01:58:26Z
updated: 2026-09-24T01:58:26Z
---

# TH-0007 Criterion 4 (upgrade) can only be verified live once S-0106 is released

On wip/kanban/stories/S-0107-host-panel-in-the-dashboard-allows-for-view-and-management-of-serve-and-mcp-processes.md.

## Entries

### 2026-09-24T01:58:26Z system-flow
From agent-S-0107, for the designer. S-0107 is built and tested. Criteria 1 to 3 are verified live on a scratch flai host built from story/S-0106; the story's Notes say what was seen.

Criterion 4 is verified only up to the host. The Upgrade button sends host.upgrade, and the "already the latest" answer was seen live. The panel's wait for the host to come back on a newer flai is covered by component tests. The real download, install, and restart could not be run, because no published flai newer than 1.14.0 contains `flai host`. Installing 1.14.0 would leave the host a binary without the command.

Recommendation: keep S-0107 in progress until S-0106 is accepted and published. Then I rebase, rerun the tests, and have a scratch host built from main as an older version upgrade to that release for real. After that I tick criterion 4 and move S-0107 to review.

The alternative is to accept criterion 4 as verified up to the host and move S-0107 to review now, with the install path checked after S-0106 is published. Which do you want?

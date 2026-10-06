---
id: I-0089
title: design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile
class: efficiency
status: open
count: 1
cost: 3m
first_reported: 2026-10-06T10:32:27Z
last_reported: 2026-10-06T10:32:27Z
updated: 2026-10-06T10:32:27Z
---

# I-0089 design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile

## Description

`flai issue new`, `bump`, and `close` regenerate `design/issues/summary.md`, and the file is committed. A story records its issues on its branch. When an issue is also recorded on main while the story works, both sides have rewritten the same table, and the story's next `flai stream sync` stops on a conflict that the agent clears by regenerating the file. The same would stop an acceptance merge. It grows with the number of stories in progress at once.

## Instances

### 2026-10-06T10:32:27Z
Story: S-0220.
S-0220's agent met a conflict in design/issues/summary.md at two syncs, 06:07Z and 07:03Z, after issues were recorded on main while it worked (I-0083, then I-0084), and regenerated the file each time. The cost is an estimate of the agent's time per conflict.

## Remediation

Directions to weigh: `flai stream sync` and `flai accept` regenerate `summary.md` themselves when it is the only conflict, since it is derived from the issue files; or the file is not committed and is generated where it is read.

S-0278 built the first direction ([ADR-0098](../adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)) for I-0074, which has the same cause. The operator chooses at its acceptance whether that closes this issue too.
